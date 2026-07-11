import type { FrameEditorPersisted } from "./frame-editor";
import { defaultEditorModel } from "./frame-editor";

export const EDITOR_STORAGE_KEY = "gallery-wall-editor-v1";

export function loadEditorState(): FrameEditorPersisted | null {
  const raw = localStorage.getItem(EDITOR_STORAGE_KEY);
  if (!raw) return null;
  try {
    const parsed = JSON.parse(raw) as FrameEditorPersisted;
    if (parsed.version !== 1 || !parsed.centerpiece || !Array.isArray(parsed.batches)) {
      return null;
    }
    return parsed;
  } catch {
    return null;
  }
}

export function saveEditorState(state: FrameEditorPersisted): void {
  localStorage.setItem(EDITOR_STORAGE_KEY, JSON.stringify(state));
}

export function clearEditorState(): void {
  localStorage.removeItem(EDITOR_STORAGE_KEY);
}

export function defaultEditorPersisted(): FrameEditorPersisted {
  const model = defaultEditorModel();
  return {
    version: 1,
    activeTab: "form",
    centerpiece: model.centerpiece,
    batches: model.batches,
  };
}
