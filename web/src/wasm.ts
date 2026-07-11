// wasm_exec.js attaches Go to the global scope.
declare global {
  class Go {
    importObject: WebAssembly.Imports;
    run(instance: WebAssembly.Instance): Promise<void>;
  }

  interface Window {
    goLayout?: (paramsJSON: string) => string;
  }
}

export {};

let ready: Promise<void> | null = null;

export function initLayoutWasm(): Promise<void> {
  if (!ready) {
    ready = loadWasm();
  }
  return ready;
}

async function loadWasm(): Promise<void> {
  const go = new Go();
  const response = await fetch("/layout.wasm");
  const bytes = await response.arrayBuffer();
  const { instance } = await WebAssembly.instantiate(bytes, go.importObject);
  go.run(instance);
  if (typeof window.goLayout !== "function") {
    throw new Error("goLayout not registered");
  }
}

export function runLayout(paramsJSON: string): string {
  if (typeof window.goLayout !== "function") {
    throw new Error("WASM not loaded");
  }
  return window.goLayout(paramsJSON);
}
