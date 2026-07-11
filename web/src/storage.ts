import type { Params, Result, SessionState } from "./types";
import { STORAGE_KEY } from "./types";

export function saveSession(state: SessionState): void {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
}

export function loadSession(): SessionState | null {
  const raw = localStorage.getItem(STORAGE_KEY);
  if (!raw) return null;
  try {
    const parsed = JSON.parse(raw) as SessionState;
    if (parsed.version !== 1 || !parsed.params) return null;
    return parsed;
  } catch {
    return null;
  }
}

export function clearSession(): void {
  localStorage.removeItem(STORAGE_KEY);
}

export function downloadJSON(filename: string, data: unknown): void {
  const blob = new Blob([JSON.stringify(data, null, 2)], {
    type: "application/json",
  });
  triggerDownload(filename, blob);
}

export function downloadText(filename: string, text: string, mime: string): void {
  const blob = new Blob([text], { type: mime });
  triggerDownload(filename, blob);
}

function triggerDownload(filename: string, blob: Blob): void {
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  a.click();
  URL.revokeObjectURL(url);
}

export function readParamsFromForm(
  gapEl: HTMLInputElement,
  wallWEl: HTMLInputElement,
  wallHEl: HTMLInputElement,
  items: Params["items"],
): Params {
  const params: Params = {
    gap: Number(gapEl.value) || 0,
    items,
  };
  const wallW = Number(wallWEl.value);
  const wallH = Number(wallHEl.value);
  if (wallW > 0 && wallH > 0) {
    params.wallWidth = wallW;
    params.wallHeight = wallH;
  }
  return params;
}

export function paramsToItemsOnly(params: Params): string {
  return JSON.stringify(params.items, null, 2);
}

export function sessionFromParams(params: Params, result?: Result): SessionState {
  return {
    version: 1,
    savedAt: new Date().toISOString(),
    params,
    result,
  };
}

export function applyParamsToForm(
  params: Params,
  gapEl: HTMLInputElement,
  wallWEl: HTMLInputElement,
  wallHEl: HTMLInputElement,
): void {
  gapEl.value = String(params.gap ?? 2);
  wallWEl.value = String(params.wallWidth ?? 0);
  wallHEl.value = String(params.wallHeight ?? 0);
}
