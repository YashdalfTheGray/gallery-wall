import "@awesome.me/webawesome/dist/components/button/button.js";
import "@awesome.me/webawesome/dist/components/button-group/button-group.js";
import "@awesome.me/webawesome/dist/components/dropdown/dropdown.js";
import "@awesome.me/webawesome/dist/components/icon/icon.js";

/** Apply light/dark class from OS preference (Web Awesome uses explicit wa-light / wa-dark). */
export function syncColorScheme(): void {
  const root = document.documentElement;
  root.classList.remove("wa-light", "wa-dark");
  root.classList.add(
    window.matchMedia("(prefers-color-scheme: dark)").matches ? "wa-dark" : "wa-light",
  );
}

export function initTheme(onSchemeChange?: () => void): void {
  syncColorScheme();
  window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change", () => {
    syncColorScheme();
    onSchemeChange?.();
  });
}
