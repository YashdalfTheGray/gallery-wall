import type { Params, PlacedResult, Result } from "./types";

const SCALE = 8;
const PADDING = 32;
const MEASURE_LABEL_OFFSET = 12;
const MEASURE_LABEL_HEIGHT = 10;
const EXPORT_FOOTER_GAP = 16;
const EXPORT_FOOTER_LINE_HEIGHT = 16;

export type SvgMode = "preview" | "export";

interface SvgOptions {
  params: Params;
  result: Result;
  mode?: SvgMode;
  unitsLabel?: string;
}

export function resultToSVG({
  params,
  result,
  mode = "preview",
  unitsLabel = "layout units",
}: SvgOptions): string {
  const b = result.bounds;
  const wallW = (b.maxX - b.minX) * SCALE;
  const exportMode = mode === "export";

  const items = [...result.items].sort((a, b) => {
    if (a.id === "main") return 1;
    if (b.id === "main") return -1;
    return a.id.localeCompare(b.id);
  });

  const contentBottom = contentBottomY(items, b);
  const footerLines = exportMode ? exportFooterLines(params, result, unitsLabel) : [];
  const footerHeight = footerLines.length * EXPORT_FOOTER_LINE_HEIGHT;
  const svgW = wallW + PADDING * 2;
  const svgH = exportMode
    ? contentBottom + PADDING + EXPORT_FOOTER_GAP + footerHeight
    : contentBottom + PADDING;
  const footerY = contentBottom + EXPORT_FOOTER_GAP + 12;

  const lines: string[] = [];
  lines.push(
    `<svg xmlns="http://www.w3.org/2000/svg" class="layout-svg${exportMode ? " layout-svg-export" : ""}" width="${svgW}" height="${svgH}" viewBox="0 0 ${svgW} ${svgH}">`,
  );

  if (exportMode) {
    lines.push(`<rect class="svg-canvas" width="100%" height="100%" fill="#fafafa"/>`);
  } else {
    lines.push(`<rect class="svg-canvas" width="100%" height="100%"/>`);
  }

  const ax = PADDING + (0 - b.minX) * SCALE;
  const ay = PADDING + (0 - b.minY) * SCALE;
  const guideStroke = exportMode ? `stroke="#999"` : "";

  if (params.wallWidth && params.wallHeight) {
    const wx = PADDING + (-params.wallWidth / 2 - b.minX) * SCALE;
    const wy = PADDING + (-params.wallHeight / 2 - b.minY) * SCALE;
    const wallStroke = exportMode ? `stroke="#c44"` : "";
    lines.push(
      `<rect class="svg-wall" x="${wx}" y="${wy}" width="${params.wallWidth * SCALE}" height="${params.wallHeight * SCALE}" fill="none" ${wallStroke} stroke-width="1.5" stroke-dasharray="6 4"/>`,
    );
  }

  for (const item of items) {
    lines.push(...frameShapeSVG(item, b, exportMode));
  }

  lines.push(
    `<line class="svg-guide" x1="${ax - 12}" y1="${ay}" x2="${ax + 12}" y2="${ay}" ${guideStroke} stroke-width="1" stroke-dasharray="4 3"/>`,
  );
  lines.push(
    `<line class="svg-guide" x1="${ax}" y1="${ay - 12}" x2="${ax}" y2="${ay + 12}" ${guideStroke} stroke-width="1" stroke-dasharray="4 3"/>`,
  );

  for (const item of items) {
    lines.push(...frameLabelSVG(item, b, exportMode, mainPixelRect(items, b)));
  }

  if (exportMode) {
    footerLines.forEach((line, i) => {
      lines.push(
        `<text x="${PADDING}" y="${footerY + i * EXPORT_FOOTER_LINE_HEIGHT}" fill="#333" font-family="system-ui,sans-serif" font-size="12">${esc(line)}</text>`,
      );
    });
  }

  lines.push("</svg>");
  return lines.join("\n");
}

function frameShapeSVG(item: PlacedResult, b: Result["bounds"], exportMode: boolean): string[] {
  const x = PADDING + (item.x - b.minX) * SCALE;
  const y = PADDING + (item.y - b.minY) * SCALE;
  const w = item.width * SCALE;
  const h = item.height * SCALE;
  const shapeClass = frameClass(item);
  const { fill, stroke, sw } = exportMode ? exportFrameColors(item) : { fill: "", stroke: "", sw: 0 };

  const strokeAttr = exportMode ? `stroke="${stroke}" stroke-width="${sw}"` : `stroke-width="${frameStrokeWidth(item)}"`;
  const fillAttr = exportMode ? `fill="${fill}"` : "";

  if (item.shape === "circle") {
    const cx = x + w / 2;
    const cy = y + h / 2;
    const r = Math.min(w, h) / 2 - 1;
    return [
      `<ellipse class="${shapeClass}" cx="${cx}" cy="${cy}" rx="${r}" ry="${r}" ${fillAttr} ${strokeAttr}/>`,
    ];
  }

  if (item.shape === "ellipse") {
    const cx = x + w / 2;
    const cy = y + h / 2;
    return [
      `<ellipse class="${shapeClass}" cx="${cx}" cy="${cy}" rx="${w / 2 - 1}" ry="${h / 2 - 1}" ${fillAttr} ${strokeAttr}/>`,
    ];
  }

  const rx = item.shape === "square" || item.id === "main" ? 3 : 0;
  return [
    `<rect class="${shapeClass}" x="${x}" y="${y}" width="${w}" height="${h}" rx="${rx}" ${fillAttr} ${strokeAttr}/>`,
  ];
}

function frameLabelSVG(
  item: PlacedResult,
  b: Result["bounds"],
  exportMode: boolean,
  mainRect: PixelRect | null,
): string[] {
  const x = PADDING + (item.x - b.minX) * SCALE;
  const y = PADDING + (item.y - b.minY) * SCALE;
  const w = item.width * SCALE;
  const h = item.height * SCALE;
  const out: string[] = [];

  if (w >= 28 && h >= 28) {
    const cx = x + w / 2;
    const cy = y + h / 2;
    const fs = w < 48 || h < 48 ? 8 : 10;
    const labelClass = item.id === "main" ? "frame-label frame-label-main" : "frame-label";
    const labelFill =
      exportMode && item.id === "main" ? `fill="#fff"` : exportMode ? `fill="#333"` : "";
    out.push(
      `<text class="${labelClass}" x="${cx}" y="${cy - 4}" ${labelFill} font-family="system-ui,sans-serif" font-size="${fs}" text-anchor="middle" dominant-baseline="middle">${esc(item.id)}</text>`,
    );
    out.push(
      `<text class="${labelClass}" x="${cx}" y="${cy + 8}" ${labelFill} font-family="system-ui,sans-serif" font-size="${fs - 1}" text-anchor="middle" dominant-baseline="middle">${item.width}×${item.height}</text>`,
    );
  }

  if (item.id !== "main") {
    const label = `(${fmt(item.centerX)}, ${fmt(item.centerY)})`;
    const { x: mx, y: my } = measureLabelPosition(x, y, w, h, label, mainRect);
    const measureFill = exportMode ? `fill="#666"` : "";
    out.push(
      `<text class="svg-measure" x="${mx}" y="${my}" ${measureFill} font-family="system-ui,sans-serif" font-size="8" text-anchor="middle">${esc(label)}</text>`,
    );
  }

  return out;
}

interface PixelRect {
  x: number;
  y: number;
  w: number;
  h: number;
}

function mainPixelRect(items: PlacedResult[], b: Result["bounds"]): PixelRect | null {
  const main = items.find((item) => item.id === "main");
  if (!main) return null;
  return {
    x: PADDING + (main.x - b.minX) * SCALE,
    y: PADDING + (main.y - b.minY) * SCALE,
    w: main.width * SCALE,
    h: main.height * SCALE,
  };
}

function measureLabelPosition(
  frameX: number,
  frameY: number,
  frameW: number,
  frameH: number,
  label: string,
  mainRect: PixelRect | null,
): { x: number; y: number } {
  const cx = frameX + frameW / 2;
  const belowY = frameY + frameH + MEASURE_LABEL_OFFSET;
  const aboveY = frameY - 6;
  const estHalfWidth = label.length * 2.4;

  if (!mainRect || !textHitsRect(cx, belowY, estHalfWidth, 8, mainRect)) {
    return { x: cx, y: belowY };
  }
  if (!textHitsRect(cx, aboveY, estHalfWidth, 8, mainRect)) {
    return { x: cx, y: aboveY };
  }

  const rightX = mainRect.x + mainRect.w + estHalfWidth + 4;
  if (!textHitsRect(rightX, belowY, estHalfWidth, 8, mainRect)) {
    return { x: rightX, y: belowY };
  }

  const leftX = mainRect.x - estHalfWidth - 4;
  return { x: leftX, y: belowY };
}

function textHitsRect(
  cx: number,
  cy: number,
  halfWidth: number,
  height: number,
  rect: PixelRect,
): boolean {
  const left = cx - halfWidth;
  const right = cx + halfWidth;
  const top = cy - height;
  const bottom = cy + 2;
  return left < rect.x + rect.w && right > rect.x && top < rect.y + rect.h && bottom > rect.y;
}

function exportFooterLines(params: Params, result: Result, unitsLabel: string): string[] {
  const b = result.bounds;
  const clusterW = b.maxX - b.minX;
  const clusterH = b.maxY - b.minY;
  return [
    `${result.items.length} frames · gap ${params.gap} · Cluster: ${fmt(clusterW)} × ${fmt(clusterH)} ${unitsLabel}`,
    `bounds x [${fmt(b.minX)}, ${fmt(b.maxX)}] y [${fmt(b.minY)}, ${fmt(b.maxY)}]`,
  ];
}

function contentBottomY(items: PlacedResult[], b: Result["bounds"]): number {
  let bottom = PADDING + (b.maxY - b.minY) * SCALE;
  for (const item of items) {
    if (item.id === "main") continue;
    const y = PADDING + (item.y - b.minY) * SCALE;
    const h = item.height * SCALE;
    bottom = Math.max(bottom, y + h + MEASURE_LABEL_OFFSET + MEASURE_LABEL_HEIGHT);
  }
  return bottom;
}

function frameClass(item: PlacedResult): string {
  if (item.id === "main") return "frame frame-main";
  switch (item.shape) {
    case "circle":
      return "frame frame-circle";
    case "ellipse":
      return `frame frame-ellipse-${frameOrientation(item)}`;
    case "square":
      return "frame frame-square";
    default:
      return `frame frame-rect-${frameOrientation(item)}`;
  }
}

type Orientation = "portrait" | "landscape";

function frameOrientation(item: PlacedResult): Orientation {
  return item.width > item.height ? "landscape" : "portrait";
}

function frameStrokeWidth(item: PlacedResult): number {
  return item.id === "main" ? 3 : 2;
}

function exportFrameColors(item: PlacedResult): { fill: string; stroke: string; sw: number } {
  if (item.id === "main") return { fill: "#4a90d9", stroke: "#1a4a7a", sw: 3 };
  switch (item.shape) {
    case "circle":
      return { fill: "#d4edda", stroke: "#2d6a3e", sw: 2 };
    case "ellipse":
      return frameOrientation(item) === "landscape"
        ? { fill: "#fce4ec", stroke: "#c2185b", sw: 2 }
        : { fill: "#fff3cd", stroke: "#856404", sw: 2 };
    case "square":
      return { fill: "#f0e6ff", stroke: "#5a3d8a", sw: 2 };
    default:
      return frameOrientation(item) === "landscape"
        ? { fill: "#d4e8f7", stroke: "#2b6cb0", sw: 2 }
        : { fill: "#ffe5d9", stroke: "#c45c26", sw: 2 };
  }
}

function fmt(n: number): string {
  return Number.isInteger(n) ? String(n) : n.toFixed(1);
}

function esc(s: string): string {
  return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}
