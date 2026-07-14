import {
  clearEditorState,
  loadEditorState,
  saveEditorState,
} from "./editor-storage";
import {
  type FrameEditorModel,
  type ItemsEditorTab,
  defaultEditorModel,
  editorModelToItems,
  formatItemsJSON,
  itemsToEditorModel,
} from "./frame-editor";
import { FrameEditorUI } from "./frame-editor-ui";
import type { Item } from "./types";

export class ItemsEditor {
  private readonly formPanel: HTMLElement;
  private readonly jsonPanel: HTMLElement;
  private readonly tabButtons: HTMLButtonElement[];
  private readonly jsonEl: HTMLTextAreaElement;
  private readonly formUI: FrameEditorUI;
  private items: Item[];
  private activeTab: ItemsEditorTab = "form";
  private syncing = false;
  private onChangeHandler: ((items: Item[]) => void) | null = null;
  private onErrorHandler: ((message: string) => void) | null = null;

  constructor(root: HTMLElement, jsonEl: HTMLTextAreaElement) {
    const editor = must(root.querySelector("#items-editor"), "#items-editor");
    this.formPanel = must(editor.querySelector("#panel-form"), "#panel-form");
    this.jsonPanel = must(editor.querySelector("#panel-json"), "#panel-json");
    this.tabButtons = [...editor.querySelectorAll<HTMLButtonElement>(".editor-tab")];
    this.jsonEl = jsonEl;
    this.formUI = new FrameEditorUI(must(editor.querySelector("#frame-form"), "#frame-form"));
    this.items = editorModelToItems(defaultEditorModel());

    this.formUI.onChange(() => this.onFormChanged());
    this.jsonEl.addEventListener("change", () => this.onJsonChanged());
    for (const btn of this.tabButtons) {
      btn.addEventListener("click", () => {
        const tab = btn.dataset.tab as ItemsEditorTab;
        if (tab) this.switchTab(tab);
      });
    }
  }

  onChange(handler: (items: Item[]) => void): void {
    this.onChangeHandler = handler;
  }

  onError(handler: (message: string) => void): void {
    this.onErrorHandler = handler;
  }

  getItems(): Item[] {
    return this.items;
  }

  getActiveTab(): ItemsEditorTab {
    return this.activeTab;
  }

  loadFromItems(items: Item[], preferredTab?: ItemsEditorTab): void {
    this.items = items;
    let model: FrameEditorModel;
    try {
      model = itemsToEditorModel(items);
    } catch {
      model = defaultEditorModel();
      this.items = editorModelToItems(model);
    }

    const saved = loadEditorState();
    this.activeTab = preferredTab ?? saved?.activeTab ?? "form";
    this.setActiveTab(this.activeTab);

    this.formUI.render(model);
    this.jsonEl.value = formatItemsJSON(this.items);
    this.persistEditor(model);
  }

  reset(): void {
    clearEditorState();
    const model = defaultEditorModel();
    this.items = editorModelToItems(model);
    this.activeTab = "form";
    this.setActiveTab("form");
    this.formUI.render(model);
    this.jsonEl.value = formatItemsJSON(this.items);
    this.persistEditor(model);
  }

  private switchTab(tab: ItemsEditorTab): void {
    if (tab === this.activeTab) return;

    try {
      if (tab === "json") {
        this.syncFormToItems();
      } else {
        this.syncJsonToForm();
      }
    } catch (err) {
      this.onErrorHandler?.(String(err));
      return;
    }

    this.activeTab = tab;
    this.setActiveTab(tab);
    this.persistEditor(this.readEditorModel());
  }

  private setActiveTab(tab: ItemsEditorTab): void {
    this.formPanel.hidden = tab !== "form";
    this.jsonPanel.hidden = tab !== "json";
    for (const btn of this.tabButtons) {
      const selected = btn.dataset.tab === tab;
      btn.setAttribute("aria-selected", selected ? "true" : "false");
    }
  }

  private onFormChanged(): void {
    if (this.syncing) return;
    this.syncFormToItems();
    this.emitChange();
  }

  private onJsonChanged(): void {
    if (this.syncing || this.activeTab !== "json") return;
    try {
      this.items = parseItemsJSON(this.jsonEl.value);
      const model = itemsToEditorModel(this.items);
      this.formUI.render(model);
      this.persistEditor(model);
      this.emitChange();
    } catch {
      // keep editing until valid JSON on blur
    }
  }

  private syncFormToItems(): void {
    const model = this.formUI.readModel();
    this.items = editorModelToItems(model);
    this.syncing = true;
    this.jsonEl.value = formatItemsJSON(this.items);
    this.syncing = false;
    this.persistEditor(model);
  }

  private syncJsonToForm(): void {
    this.items = parseItemsJSON(this.jsonEl.value);
    const model = itemsToEditorModel(this.items);
    this.formUI.render(model);
    this.persistEditor(model);
  }

  private readEditorModel(): FrameEditorModel {
    if (this.activeTab === "form") {
      return this.formUI.readModel();
    }
    return itemsToEditorModel(this.items);
  }

  private persistEditor(model: FrameEditorModel): void {
    saveEditorState({
      version: 1,
      activeTab: this.activeTab,
      centerpiece: model.centerpiece,
      batches: model.batches,
    });
  }

  private emitChange(): void {
    this.onChangeHandler?.(this.items);
  }

  readItemsForLayout(): Item[] {
    if (this.activeTab === "form") {
      this.syncFormToItems();
    } else {
      this.items = parseItemsJSON(this.jsonEl.value);
    }
    return this.items;
  }
}

function parseItemsJSON(text: string): Item[] {
  const items = JSON.parse(text) as Item[];
  if (!Array.isArray(items) || items.length === 0) {
    throw new Error("items must be a non-empty array");
  }
  if (!items.some((i) => i.centerpiece)) {
    throw new Error("items must include a centerpiece");
  }
  return items;
}

function must<T>(el: T | null, label: string): T {
  if (!el) throw new Error(`${label} not found`);
  return el;
}
