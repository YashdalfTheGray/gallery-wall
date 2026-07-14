import type { Item, Shape } from "./types";

export const SHAPES: Shape[] = ["rectangle", "square", "circle", "ellipse"];

export interface CenterpieceSpec {
  shape: Shape;
  height: number;
  width: number;
}

export interface FrameBatch {
  shape: Shape;
  height: number;
  width: number;
  count: number;
}

export interface FrameEditorModel {
  centerpiece: CenterpieceSpec;
  batches: FrameBatch[];
}

export type ItemsEditorTab = "form" | "json";

export interface FrameEditorPersisted {
  version: 1;
  activeTab: ItemsEditorTab;
  centerpiece: CenterpieceSpec;
  batches: FrameBatch[];
}

export function needsWidthField(shape: Shape): boolean {
  return shape === "rectangle" || shape === "ellipse";
}

export function usesSingleSize(shape: Shape): boolean {
  return shape === "square" || shape === "circle";
}

export function defaultEditorModel(): FrameEditorModel {
  return {
    centerpiece: { shape: "rectangle", height: 16, width: 14 },
    batches: [
      { shape: "rectangle", height: 10, width: 8, count: 1 },
      { shape: "rectangle", height: 10, width: 8, count: 1 },
    ],
  };
}

export function itemsToEditorModel(items: Item[]): FrameEditorModel {
  const main = items.find((i) => i.centerpiece);
  if (!main) {
    throw new Error("items must include exactly one centerpiece");
  }

  const centerpiece: CenterpieceSpec = {
    shape: main.shape,
    height: main.height,
    width: main.width,
  };

  const groups = new Map<string, FrameBatch>();
  for (const item of items) {
    if (item.centerpiece) continue;
    const key = `${item.shape}:${item.height}:${item.width}`;
    const existing = groups.get(key);
    if (existing) {
      existing.count++;
    } else {
      groups.set(key, {
        shape: item.shape,
        height: item.height,
        width: item.width,
        count: 1,
      });
    }
  }

  return {
    centerpiece,
    batches: [...groups.values()],
  };
}

export function editorModelToItems(model: FrameEditorModel): Item[] {
  const cp = model.centerpiece;
  const items: Item[] = [
    {
      id: "main",
      shape: cp.shape,
      height: cp.height,
      width: cp.width,
      centerpiece: true,
    },
  ];

  let n = 1;
  for (const batch of model.batches) {
    const count = Math.max(0, Math.floor(batch.count));
    for (let i = 0; i < count; i++) {
      items.push({
        id: `p${String(n).padStart(2, "0")}`,
        shape: batch.shape,
        height: batch.height,
        width: batch.width,
      });
      n++;
    }
  }

  return items;
}

export function normalizeBatch(batch: FrameBatch): FrameBatch {
  const height = Math.max(1, Math.floor(batch.height));
  let width = Math.max(1, Math.floor(batch.width));
  if (usesSingleSize(batch.shape)) {
    width = height;
  }
  return {
    shape: batch.shape,
    height,
    width,
    count: Math.max(1, Math.floor(batch.count)),
  };
}

export function normalizeCenterpiece(cp: CenterpieceSpec): CenterpieceSpec {
  const height = Math.max(1, Math.floor(cp.height));
  let width = Math.max(1, Math.floor(cp.width));
  if (usesSingleSize(cp.shape)) {
    width = height;
  }
  return { shape: cp.shape, height, width };
}

export function frameCount(model: FrameEditorModel): number {
  const satellites = model.batches.reduce((sum, b) => sum + Math.max(0, Math.floor(b.count)), 0);
  return 1 + satellites;
}

export function formatItemsJSON(items: Item[]): string {
  return JSON.stringify(items, null, 2);
}
