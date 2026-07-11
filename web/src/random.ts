import type { Item, Params, Shape } from "./types";

const SHAPES: Shape[] = ["rectangle", "square", "circle", "ellipse"];

export function randomFrameCount(): number {
  return randInt(5, 20);
}

export function randomParams(frameCount = randomFrameCount(), gap = 2): Params {
  const count = Math.max(5, Math.min(20, frameCount));
  const items: Item[] = [
    {
      id: "main",
      height: randInt(12, 18),
      width: randInt(11, 16),
      shape: "rectangle",
      centerpiece: true,
    },
  ];

  for (let i = 1; i < count; i++) {
    items.push({
      id: `p${String(i).padStart(2, "0")}`,
      height: randInt(5, 12),
      width: randInt(5, 10),
      shape: SHAPES[randInt(0, SHAPES.length - 1)]!,
    });
  }

  return { gap, items };
}

function randInt(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}
