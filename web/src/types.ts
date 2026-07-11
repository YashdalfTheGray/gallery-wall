export type Shape = "square" | "rectangle" | "circle" | "ellipse";

export interface Item {
  id: string;
  height: number;
  width: number;
  shape: Shape;
  centerpiece?: boolean;
}

export interface Params {
  gap: number;
  wallWidth?: number;
  wallHeight?: number;
  items: Item[];
}

export interface PlacedResult {
  id: string;
  centerX: number;
  centerY: number;
  x: number;
  y: number;
  width: number;
  height: number;
  shape: Shape;
  offsetFromAnchor: number;
  direction: string;
  adjacentIds?: string[];
}

export interface Bounds {
  minX: number;
  minY: number;
  maxX: number;
  maxY: number;
}

export interface Result {
  anchor: { itemId: string; centerX: number; centerY: number };
  items: PlacedResult[];
  bounds: Bounds;
}

export interface LayoutError {
  code: string;
  message: string;
  itemId?: string;
  itemIds?: string[];
}

export interface LayoutResponse {
  ok: boolean;
  result?: Result;
  error?: LayoutError;
  message?: string;
}

export interface SessionState {
  version: 1;
  savedAt: string;
  params: Params;
  result?: Result;
}

export const STORAGE_KEY = "gallery-wall-session-v1";
