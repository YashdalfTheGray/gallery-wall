import { clearEditorState } from "./editor-storage";
import { initTheme } from "./webawesome";
import { ItemsEditor } from "./items-editor";
import { editorModelToItems, defaultEditorModel } from "./frame-editor";
import { randomParams } from "./random";
import { resultToSVG } from "./svg";
import {
  applyParamsToForm,
  clearSession,
  downloadJSON,
  downloadText,
  loadSession,
  readParamsFromForm,
  saveSession,
  sessionFromParams,
} from "./storage";
import type { LayoutResponse, Params, Result, SessionState } from "./types";
import { initLayoutWasm, runLayout } from "./wasm";

interface WaButton extends HTMLElement {
  disabled: boolean;
}

const gapEl = must<HTMLInputElement>("gap");
const wallWEl = must<HTMLInputElement>("wall-width");
const wallHEl = must<HTMLInputElement>("wall-height");
const jsonEl = must<HTMLTextAreaElement>("params-json");
const previewEl = must<HTMLDivElement>("preview");
const statsEl = must<HTMLDivElement>("stats");
const statusEl = must<HTMLDivElement>("status");
const btnLayout = must<WaButton>("btn-layout");
const btnImport = must<WaButton>("btn-import");
const btnClear = must<WaButton>("btn-clear");
const btnExportToggle = must<WaButton>("btn-export-toggle");
const btnExport = must<HTMLElement>("btn-export");
const btnSvg = must<HTMLElement>("btn-svg");
const loadDropdown = must<HTMLElement>("load-dropdown");
const fileImport = must<HTMLInputElement>("file-import");

const itemsEditor = new ItemsEditor(must<HTMLElement>("items-editor-wrap"), jsonEl);

let currentResult: Result | null = null;
let currentParams: Params | null = null;

async function boot(): Promise<void> {
  initTheme(refreshPreviewTheme);

  setStatus("Loading layout engine…");
  await initLayoutWasm();

  const saved = loadSession();
  if (saved) {
    applyParamsToForm(saved.params, gapEl, wallWEl, wallHEl);
    itemsEditor.loadFromItems(saved.params.items);
    currentParams = saved.params;
    if (saved.result) {
      showResult(saved.params, saved.result);
      setStatus("Restored from local storage.");
    } else {
      setStatus("Restored params from local storage.");
    }
  } else {
    itemsEditor.loadFromItems(editorModelToItems(defaultEditorModel()));
    setStatus("Ready.");
  }

  itemsEditor.onChange(() => persistDraft());
  itemsEditor.onError((msg) => setStatus(msg, true));

  btnLayout.addEventListener("click", () => void doLayout());
  btnImport.addEventListener("click", () => fileImport.click());
  btnClear.addEventListener("click", () => clearStorage());
  fileImport.addEventListener("change", () => void importFile());

  loadDropdown.addEventListener("wa-select", (e) => {
    const item = (e as CustomEvent<{ item: HTMLElement }>).detail.item;
    if (item.id === "btn-import-menu") fileImport.click();
    if (item.id === "btn-random") loadRandom();
    if (item.id === "btn-demo") void loadDemo();
  });

  must<HTMLElement>("export-dropdown").addEventListener("wa-select", (e) => {
    const item = (e as CustomEvent<{ item: HTMLElement }>).detail.item;
    if (item.id === "btn-export") exportState();
    if (item.id === "btn-svg") exportSVG();
  });

  [gapEl, wallWEl, wallHEl].forEach((el) => {
    el.addEventListener("change", persistDraft);
  });
}

function refreshPreviewTheme(): void {
  if (currentParams && currentResult) {
    previewEl.innerHTML = resultToSVG({
      params: currentParams,
      result: currentResult,
      mode: "preview",
    });
  }
}

async function doLayout(): Promise<void> {
  try {
    const items = itemsEditor.readItemsForLayout();
    const params = readParamsFromForm(gapEl, wallWEl, wallHEl, items);
    setStatus("Laying out…");
    const raw = runLayout(JSON.stringify(params));
    const resp = JSON.parse(raw) as LayoutResponse;

    if (!resp.ok) {
      const msg = resp.error?.message ?? resp.message ?? "layout failed";
      setStatus(msg, true);
      return;
    }

    if (!resp.result) {
      setStatus("No result returned", true);
      return;
    }

    showResult(params, resp.result);
    saveSession(sessionFromParams(params, resp.result));
    setStatus(`Placed ${resp.result.items.length} frames. Saved locally.`);
  } catch (err) {
    setStatus(String(err), true);
  }
}

function showResult(params: Params, result: Result): void {
  currentParams = params;
  currentResult = result;
  setExportEnabled(true);

  const svg = resultToSVG({ params, result, mode: "preview" });
  previewEl.innerHTML = svg;

  const b = result.bounds;
  const w = b.maxX - b.minX;
  const h = b.maxY - b.minY;
  statsEl.innerHTML = `
    <dl class="stats-panel">
      <dt>Frames</dt><dd>${result.items.length}</dd>
      <dt>Gap</dt><dd>${params.gap}</dd>
      <dt>Cluster size</dt><dd>${fmt(w)} × ${fmt(h)}</dd>
    </dl>
    <dl class="stats-panel">
      <dt>Anchor</dt><dd>${result.anchor.itemId} @ (0, 0)</dd>
      <dt>Bounds.x</dt><dd>[${fmt(b.minX)}, ${fmt(b.maxX)}]</dd>
      <dt>Bounds.y</dt><dd>[${fmt(b.minY)}, ${fmt(b.maxY)}]</dd>
    </dl>
  `;
}

async function loadDemo(): Promise<void> {
  const res = await fetch("/presets/twentyfive.json");
  const params = (await res.json()) as Params;
  applyParams(params);
  setStatus(`Loaded 25-frame sample. Click Arrange frames.`);
}

function loadRandom(): void {
  const gap = Number(gapEl.value) || 2;
  const params = randomParams(undefined, gap);
  applyParams(params);
  setStatus(`Generated ${params.items.length} random frames. Click Arrange frames.`);
}

function applyParams(params: Params): void {
  applyParamsToForm(params, gapEl, wallWEl, wallHEl);
  itemsEditor.loadFromItems(params.items);
  currentParams = params;
  currentResult = null;
  previewEl.innerHTML = "";
  statsEl.innerHTML = "";
  setExportEnabled(false);
  persistDraft();
}

function clearStorage(): void {
  clearSession();
  clearEditorState();
  gapEl.value = "2";
  wallWEl.value = "0";
  wallHEl.value = "0";
  itemsEditor.reset();
  currentParams = null;
  currentResult = null;
  previewEl.innerHTML = "";
  statsEl.innerHTML = "";
  setExportEnabled(false);
  setStatus("Saved session cleared.");
}

function setExportEnabled(on: boolean): void {
  btnExportToggle.disabled = !on;
  setDisabled(btnExport, !on);
  setDisabled(btnSvg, !on);
}

function setDisabled(el: HTMLElement, disabled: boolean): void {
  if (disabled) el.setAttribute("disabled", "");
  else el.removeAttribute("disabled");
}

function exportState(): void {
  if (!currentParams) return;
  const state: SessionState = sessionFromParams(currentParams, currentResult ?? undefined);
  downloadJSON("gallery-wall-state.json", state);
}

function exportSVG(): void {
  if (!currentParams || !currentResult) return;
  const svg = resultToSVG({ params: currentParams, result: currentResult, mode: "export" });
  downloadText("gallery-wall.svg", svg, "image/svg+xml");
}

async function importFile(): Promise<void> {
  const file = fileImport.files?.[0];
  if (!file) return;
  try {
    const text = await file.text();
    const data = JSON.parse(text) as SessionState | Params;
    const params = "params" in data && data.params ? data.params : (data as Params);
    if (!params.items?.length) throw new Error("invalid JSON: missing items");
    applyParamsToForm(params, gapEl, wallWEl, wallHEl);
    itemsEditor.loadFromItems(params.items);
    currentParams = params;
    if ("result" in data && data.result) {
      showResult(params, data.result);
    } else {
      currentResult = null;
      previewEl.innerHTML = "";
      statsEl.innerHTML = "";
      setExportEnabled(false);
    }
    saveSession(sessionFromParams(params, "result" in data ? data.result : undefined));
    setStatus(`Imported ${file.name}`);
  } catch (err) {
    setStatus(String(err), true);
  } finally {
    fileImport.value = "";
  }
}

function persistDraft(): void {
  try {
    const items = itemsEditor.getItems();
    const params = readParamsFromForm(gapEl, wallWEl, wallHEl, items);
    saveSession(sessionFromParams(params, currentResult ?? undefined));
  } catch {
    // ignore invalid state while editing
  }
}

function setStatus(msg: string, isError = false): void {
  statusEl.textContent = msg;
  statusEl.classList.toggle("error", isError);
}

function fmt(n: number): string {
  return Number.isInteger(n) ? String(n) : n.toFixed(1);
}

function must<T extends HTMLElement>(id: string): T {
  const el = document.getElementById(id);
  if (!el) throw new Error(`#${id} not found`);
  return el as T;
}

void boot();
