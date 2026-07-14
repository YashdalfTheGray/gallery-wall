import {
  type CenterpieceSpec,
  type FrameBatch,
  type FrameEditorModel,
  SHAPES,
  frameCount,
  normalizeBatch,
  normalizeCenterpiece,
  usesSingleSize,
} from "./frame-editor";
import type { Shape } from "./types";

type ChangeHandler = () => void;

export class FrameEditorUI {
  private readonly centerpieceEl: HTMLElement;
  private readonly batchListEl: HTMLElement;
  private readonly summaryEl: HTMLElement;
  private onChangeHandler: ChangeHandler | null = null;

  constructor(root: HTMLElement) {
    this.centerpieceEl = mustChild(root, ".centerpiece-fields");
    this.batchListEl = mustChild(root, ".batch-list");
    this.summaryEl = mustChild(root, ".frame-summary");
    mustChild(root, "#btn-add-batch").addEventListener("click", () => this.addBatch());
  }

  onChange(handler: ChangeHandler): void {
    this.onChangeHandler = handler;
  }

  render(model: FrameEditorModel): void {
    this.renderCenterpiece(model.centerpiece);
    this.batchListEl.innerHTML = "";
    model.batches.forEach((batch, index) => this.batchListEl.appendChild(this.createBatchRow(batch, index)));
    this.updateSummary(model);
  }

  readModel(): FrameEditorModel {
    return {
      centerpiece: this.readCenterpiece(),
      batches: this.readBatches(),
    };
  }

  private emitChange(): void {
    this.updateSummary(this.readModel());
    this.onChangeHandler?.();
  }

  private renderCenterpiece(cp: CenterpieceSpec): void {
    this.centerpieceEl.innerHTML = "";
    this.centerpieceEl.append(
      shapeCell("cp-shape", cp.shape, () => this.onCenterpieceShapeChange()),
      heightCell("cp-height", cp.height, () => this.emitChange()),
      widthCell("cp-width", cp.shape, cp.width, () => this.emitChange()),
    );
  }

  private onCenterpieceShapeChange(): void {
    const cp = this.readCenterpiece();
    this.renderCenterpiece(cp);
    this.emitChange();
  }

  private readCenterpiece(): CenterpieceSpec {
    const shape = readSelect(this.centerpieceEl, "cp-shape");
    const { height, width } = readDimensions(this.centerpieceEl, "cp", shape);
    return normalizeCenterpiece({ shape, height, width });
  }

  private readBatches(): FrameBatch[] {
    return [...this.batchListEl.querySelectorAll<HTMLElement>(".batch-row")].map((row) => {
      const index = row.dataset.index ?? "0";
      const shape = readSelect(row, `batch-shape-${index}`);
      const { height, width } = readDimensions(row, `batch-${index}`, shape);
      const count = readNumber(row, `batch-count-${index}`, 1);
      return normalizeBatch({ shape, height, width, count });
    });
  }

  private addBatch(): void {
    const model = this.readModel();
    model.batches.push({ shape: "rectangle", height: 10, width: 8, count: 1 });
    this.render(model);
    this.emitChange();
  }

  private removeBatch(index: number): void {
    const model = this.readModel();
    model.batches.splice(index, 1);
    this.render(model);
    this.emitChange();
  }

  private createBatchRow(batch: FrameBatch, index: number): HTMLElement {
    const row = document.createElement("div");
    row.className = "batch-row frame-row";
    row.dataset.index = String(index);

    row.append(
      shapeCell(`batch-shape-${index}`, batch.shape, () => this.onBatchShapeChange(row, index)),
      copiesCell(`batch-count-${index}`, batch.count, () => this.emitChange()),
      heightCell(`batch-${index}-height`, batch.height, () => this.emitChange()),
      widthCell(`batch-${index}-width`, batch.shape, batch.width, () => this.emitChange()),
    );

    const remove = document.createElement("button");
    remove.type = "button";
    remove.className = "batch-remove";
    remove.title = "Remove frame type";
    remove.textContent = "×";
    remove.addEventListener("click", () => this.removeBatch(index));
    row.append(remove);

    return row;
  }

  private onBatchShapeChange(row: HTMLElement, index: number): void {
    const shape = readSelect(row, `batch-shape-${index}`);
    const { height, width } = readDimensions(row, `batch-${index}`, shape);
    const count = readNumber(row, `batch-count-${index}`, 1);
    const replacement = this.createBatchRow({ shape, height, width, count }, index);
    row.replaceWith(replacement);
    this.emitChange();
  }

  private updateSummary(model: FrameEditorModel): void {
    const total = frameCount(model);
    const satellites = total - 1;
    const batchTypes = model.batches.length;
    this.summaryEl.textContent =
      total === 1
        ? "1 frame (centerpiece only)"
        : `${total} frames — 1 centerpiece + ${satellites} satellite${satellites === 1 ? "" : "s"} across ${batchTypes} type${batchTypes === 1 ? "" : "s"}`;
  }
}

function shapeCell(id: string, value: Shape, onChange: () => void): HTMLElement {
  const cell = document.createElement("div");
  cell.className = "frame-cell";
  const select = document.createElement("select");
  select.id = id;
  select.className = "frame-input";
  for (const shape of SHAPES) {
    const opt = document.createElement("option");
    opt.value = shape;
    opt.textContent = capitalize(shape);
    select.append(opt);
  }
  select.value = value;
  select.addEventListener("change", onChange);
  cell.append(select);
  return cell;
}

function copiesCell(id: string, value: number, onChange: () => void): HTMLElement {
  const cell = document.createElement("div");
  cell.className = "frame-cell";
  const input = document.createElement("input");
  input.type = "number";
  input.id = id;
  input.className = "frame-input frame-input-num";
  input.min = "1";
  input.max = "99";
  input.step = "1";
  input.value = String(value);
  input.addEventListener("change", onChange);
  input.addEventListener("input", onChange);
  cell.append(input);
  return cell;
}

function heightCell(id: string, value: number, onChange: () => void): HTMLElement {
  const cell = document.createElement("div");
  cell.className = "frame-cell";
  const input = document.createElement("input");
  input.type = "number";
  input.id = id;
  input.className = "frame-input frame-input-num";
  input.min = "1";
  input.max = "48";
  input.step = "1";
  input.value = String(value);
  input.addEventListener("change", onChange);
  input.addEventListener("input", onChange);
  cell.append(input);
  return cell;
}

function widthCell(id: string, shape: Shape, value: number, onChange: () => void): HTMLElement {
  const cell = document.createElement("div");
  cell.className = "frame-cell";
  const single = usesSingleSize(shape);
  const input = document.createElement("input");
  input.type = "number";
  input.id = id;
  input.className = "frame-input frame-input-num";
  input.min = "1";
  input.max = "48";
  input.step = "1";
  input.value = single ? "" : String(value);
  input.disabled = single;
  input.placeholder = single ? "—" : "";
  input.setAttribute("aria-label", "Width");
  if (!single) {
    input.addEventListener("change", onChange);
    input.addEventListener("input", onChange);
  }
  cell.append(input);
  return cell;
}

function readDimensions(root: ParentNode, prefix: string, shape: Shape): { height: number; width: number } {
  const height = readNumber(root, `${prefix}-height`, 10);
  if (usesSingleSize(shape)) {
    return { height, width: height };
  }
  return {
    height,
    width: readNumber(root, `${prefix}-width`, 8),
  };
}

function readSelect(root: ParentNode, id: string): Shape {
  const el = root.querySelector<HTMLSelectElement>(`#${CSS.escape(id)}`);
  return (el?.value ?? "rectangle") as Shape;
}

function readNumber(root: ParentNode, id: string, fallback: number): number {
  const el = root.querySelector<HTMLInputElement>(`#${CSS.escape(id)}`);
  const n = Number(el?.value);
  return Number.isFinite(n) ? n : fallback;
}

function capitalize(s: string): string {
  return s.charAt(0).toUpperCase() + s.slice(1);
}

function mustChild<T extends Element>(root: ParentNode, selector: string): T {
  const el = root.querySelector(selector);
  if (!el) throw new Error(`missing ${selector}`);
  return el as T;
}
