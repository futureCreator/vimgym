# VimGym Web Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 기존 Go TUI VimGym의 Track 1-3 (86개 퍼즐)을 레트로 아케이드 게임 UI 웹 버전으로 재구현한다.

**Architecture:** Vite + React SPA. CodeMirror 6에 @replit/codemirror-vim을 얹어 Vim 에디터를 구현하고, 퍼즐 JSON을 정적 import해서 사용한다. 진행도는 localStorage에 저장. HashRouter로 GitHub Pages 배포를 지원한다.

**Tech Stack:** Vite, React 19, TypeScript, Tailwind CSS v4, CodeMirror 6, @replit/codemirror-vim, React Router 7, Vitest

---

## File Structure

```
vimgym-web/
├── index.html
├── package.json
├── tsconfig.json
├── tsconfig.app.json
├── tsconfig.node.json
├── vite.config.ts
├── postcss.config.mjs
├── public/
│   └── favicon.ico
├── src/
│   ├── main.tsx                  # React entry point
│   ├── App.tsx                   # HashRouter + routes
│   ├── index.css                 # Tailwind + retro theme + fonts
│   ├── types.ts                  # Puzzle, GameState, PuzzleResult
│   ├── data/
│   │   ├── puzzles.ts            # Puzzle loader + track/level grouping
│   │   ├── track1_foundations.json
│   │   ├── track2_editing.json
│   │   └── track3_power.json
│   ├── lib/
│   │   ├── progress.ts           # localStorage read/write
│   │   └── scoring.ts            # Star rating + clear judgment
│   ├── hooks/
│   │   └── useProgress.ts        # React hook wrapping progress store
│   ├── components/
│   │   ├── Hud.tsx               # Top HUD bar (clears, stars)
│   │   ├── Stars.tsx             # ★★★ display component
│   │   └── VimEditor.tsx         # CodeMirror + Vim mode wrapper
│   └── pages/
│       ├── TrackSelect.tsx       # Main: 3 track cards
│       ├── LevelSelect.tsx       # 5-col level grid
│       ├── PuzzleList.tsx        # Puzzle list within level
│       └── PuzzlePlay.tsx        # Play screen + clear overlay
├── tests/
│   ├── scoring.test.ts
│   ├── progress.test.ts
│   └── puzzles.test.ts
```

---

### Task 1: Project Scaffolding

**Files:**
- Create: `vimgym-web/package.json`
- Create: `vimgym-web/index.html`
- Create: `vimgym-web/vite.config.ts`
- Create: `vimgym-web/tsconfig.json`
- Create: `vimgym-web/tsconfig.app.json`
- Create: `vimgym-web/tsconfig.node.json`
- Create: `vimgym-web/postcss.config.mjs`
- Create: `vimgym-web/src/main.tsx`
- Create: `vimgym-web/src/App.tsx`
- Create: `vimgym-web/src/index.css`

**Note:** 프로젝트는 `projects/vimgym-web/`에 생성한다 (archive에서 꺼내서 활성 프로젝트로).

- [ ] **Step 1: Create package.json**

```json
{
  "name": "vimgym-web",
  "private": true,
  "version": "0.1.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc -b && vite build",
    "preview": "vite preview",
    "test": "vitest run",
    "test:watch": "vitest"
  },
  "dependencies": {
    "react": "^19.0.0",
    "react-dom": "^19.0.0",
    "react-router": "^7.0.0",
    "codemirror": "^6.0.0",
    "@codemirror/view": "^6.0.0",
    "@codemirror/state": "^6.0.0",
    "@codemirror/language": "^6.0.0",
    "@codemirror/commands": "^6.0.0",
    "@replit/codemirror-vim": "^6.0.0"
  },
  "devDependencies": {
    "@types/react": "^19.0.0",
    "@types/react-dom": "^19.0.0",
    "@vitejs/plugin-react": "^4.0.0",
    "typescript": "^5.7.0",
    "vite": "^6.0.0",
    "vitest": "^3.0.0",
    "tailwindcss": "^4.0.0",
    "@tailwindcss/postcss": "^4.0.0"
  }
}
```

- [ ] **Step 2: Create vite.config.ts**

```ts
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  base: "/vimgym-web/",
});
```

- [ ] **Step 3: Create tsconfig files**

`tsconfig.json`:
```json
{
  "files": [],
  "references": [
    { "path": "./tsconfig.app.json" },
    { "path": "./tsconfig.node.json" }
  ]
}
```

`tsconfig.app.json`:
```json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "isolatedModules": true,
    "moduleDetection": "force",
    "noEmit": true,
    "jsx": "react-jsx",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "resolveJsonModule": true
  },
  "include": ["src"]
}
```

`tsconfig.node.json`:
```json
{
  "compilerOptions": {
    "target": "ES2022",
    "lib": ["ES2023"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "isolatedModules": true,
    "moduleDetection": "force",
    "noEmit": true,
    "strict": true
  },
  "include": ["vite.config.ts"]
}
```

- [ ] **Step 4: Create postcss.config.mjs**

```js
export default {
  plugins: {
    "@tailwindcss/postcss": {},
  },
};
```

- [ ] **Step 5: Create index.html**

```html
<!doctype html>
<html lang="ko">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>VimGym</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
```

- [ ] **Step 6: Create src/index.css**

```css
@import "tailwindcss";

@theme inline {
  --color-bg: #1a1028;
  --color-bg-deep: #120c1e;
  --color-bg-card: #2a1e40;
  --color-border: #3a2860;
  --color-gold: #f0c040;
  --color-gold-shadow: #8a5000;
  --color-green: #40c060;
  --color-purple: #6a50a0;
  --color-purple-light: #b090d0;
  --color-purple-text: #c0a0e0;
  --color-text: #e0e0f0;
  --color-text-dim: #6a50a0;
  --font-pixel: "Press Start 2P", monospace;
  --font-code: "JetBrains Mono", monospace;
}

@import url("https://fonts.googleapis.com/css2?family=Press+Start+2P&family=JetBrains+Mono:wght@400;600&display=swap");

body {
  background-color: var(--color-bg);
  color: var(--color-text);
  font-family: var(--font-pixel);
  margin: 0;
}

/* CodeMirror retro overrides */
.cm-editor {
  font-family: var(--font-code) !important;
  font-size: 14px;
  background: var(--color-bg-deep) !important;
  border: 2px solid var(--color-border);
  border-radius: 6px;
}
.cm-editor .cm-content {
  color: var(--color-purple-text);
  caret-color: var(--color-gold);
}
.cm-editor .cm-cursor {
  border-left-color: var(--color-gold) !important;
}
.cm-editor .cm-activeLine {
  background: rgba(58, 40, 96, 0.3) !important;
}
.cm-editor .cm-gutters {
  background: var(--color-bg-deep) !important;
  color: var(--color-text-dim);
  border-right: 1px solid var(--color-border);
}
.cm-editor.cm-focused {
  outline: none !important;
}

/* Vim mode fat cursor */
.cm-editor .cm-fat-cursor .cm-cursor {
  background: var(--color-gold) !important;
  border: none !important;
}
.cm-editor .cm-vimMode .cm-line ::selection {
  background: rgba(240, 192, 64, 0.3) !important;
}
```

- [ ] **Step 7: Create src/main.tsx**

```tsx
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "./App";
import "./index.css";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
```

- [ ] **Step 8: Create src/App.tsx (placeholder)**

```tsx
export default function App() {
  return (
    <div className="min-h-screen flex items-center justify-center">
      <h1 className="text-gold font-pixel text-lg">VIM GYM</h1>
    </div>
  );
}
```

- [ ] **Step 9: Install dependencies and verify dev server**

```bash
cd vimgym-web && npm install
npm run dev -- --port 5180 &
sleep 3
curl -s http://localhost:5180 | head -5
kill %1
```

Expected: HTML returned with `<div id="root">`.

- [ ] **Step 10: Commit**

```bash
git add -A
git commit -m "feat: scaffold vimgym-web project with Vite + React + Tailwind"
```

---

### Task 2: Types and Puzzle Data

**Files:**
- Create: `vimgym-web/src/types.ts`
- Copy: `archive/vimgym/puzzles/track{1,2,3}_*.json` → `vimgym-web/src/data/`
- Create: `vimgym-web/src/data/puzzles.ts`
- Create: `vimgym-web/tests/puzzles.test.ts`

- [ ] **Step 1: Create src/types.ts**

```ts
export interface Cursor {
  row: number;
  col: number;
}

export interface Puzzle {
  id: string;
  title: string;
  track: number;
  level: number;
  category: string;
  difficulty: number;
  before: { text: string; cursor: Cursor };
  after: { text: string };
  par: number;
  hint: string;
  optimalSolution: string;
  solutionExplanation: string;
  tags: string[];
}

export interface PuzzleResult {
  cleared: boolean;
  stars: 1 | 2 | 3;
  bestKeystrokes: number;
}

export interface GameState {
  puzzles: Record<string, PuzzleResult>;
}

export interface Track {
  id: number;
  name: string;
  description: string;
  levels: number[];
}

export const TRACKS: Track[] = [
  {
    id: 1,
    name: "Foundations",
    description: "Motions & Navigation",
    levels: [1, 2, 3, 4, 5],
  },
  {
    id: 2,
    name: "Editing",
    description: "Insert, Delete, Change, Yank",
    levels: [6, 7, 8, 9, 10],
  },
  {
    id: 3,
    name: "Power Moves",
    description: "Text Objects, Search, Macros",
    levels: [11, 12, 13, 14, 15],
  },
];
```

- [ ] **Step 2: Copy puzzle JSON files**

```bash
mkdir -p vimgym-web/src/data
cp archive/vimgym/puzzles/track1_foundations.json vimgym-web/src/data/
cp archive/vimgym/puzzles/track2_editing.json vimgym-web/src/data/
cp archive/vimgym/puzzles/track3_power.json vimgym-web/src/data/
```

- [ ] **Step 3: Create src/data/puzzles.ts**

```ts
import type { Puzzle } from "../types";
import track1 from "./track1_foundations.json";
import track2 from "./track2_editing.json";
import track3 from "./track3_power.json";

const allPuzzles: Puzzle[] = [
  ...(track1 as Puzzle[]),
  ...(track2 as Puzzle[]),
  ...(track3 as Puzzle[]),
];

export function getAllPuzzles(): Puzzle[] {
  return allPuzzles;
}

export function getPuzzlesByTrack(trackId: number): Puzzle[] {
  return allPuzzles.filter((p) => p.track === trackId);
}

export function getPuzzlesByLevel(level: number): Puzzle[] {
  return allPuzzles.filter((p) => p.level === level);
}

export function getPuzzleById(id: string): Puzzle | undefined {
  return allPuzzles.find((p) => p.id === id);
}

export function getLevelsByTrack(trackId: number): number[] {
  const levels = new Set(
    allPuzzles.filter((p) => p.track === trackId).map((p) => p.level),
  );
  return [...levels].sort((a, b) => a - b);
}
```

- [ ] **Step 4: Write failing test**

`tests/puzzles.test.ts`:
```ts
import { describe, it, expect } from "vitest";
import {
  getAllPuzzles,
  getPuzzlesByTrack,
  getPuzzlesByLevel,
  getPuzzleById,
  getLevelsByTrack,
} from "../src/data/puzzles";

describe("puzzles", () => {
  it("loads all 86 puzzles", () => {
    expect(getAllPuzzles()).toHaveLength(86);
  });

  it("filters by track", () => {
    expect(getPuzzlesByTrack(1)).toHaveLength(29);
    expect(getPuzzlesByTrack(2)).toHaveLength(28);
    expect(getPuzzlesByTrack(3)).toHaveLength(29);
  });

  it("filters by level", () => {
    const level1 = getPuzzlesByLevel(1);
    expect(level1.length).toBe(6);
    expect(level1[0].id).toBe("hjkl-01");
  });

  it("finds puzzle by id", () => {
    const p = getPuzzleById("hjkl-01");
    expect(p).toBeDefined();
    expect(p!.title).toBe("Move Right");
    expect(p!.before.cursor.row).toBe(0);
  });

  it("returns undefined for unknown id", () => {
    expect(getPuzzleById("nonexistent")).toBeUndefined();
  });

  it("gets levels for track", () => {
    expect(getLevelsByTrack(1)).toEqual([1, 2, 3, 4, 5]);
    expect(getLevelsByTrack(2)).toEqual([6, 7, 8, 9, 10]);
    expect(getLevelsByTrack(3)).toEqual([11, 12, 13, 14, 15]);
  });

  it("every puzzle has required fields", () => {
    for (const p of getAllPuzzles()) {
      expect(p.id).toBeTruthy();
      expect(p.before.text).toBeDefined();
      expect(p.after.text).toBeDefined();
      expect(p.par).toBeGreaterThan(0);
      expect(p.hint).toBeTruthy();
    }
  });
});
```

- [ ] **Step 5: Run tests to verify they pass**

```bash
cd vimgym-web && npx vitest run tests/puzzles.test.ts
```

Expected: All 7 tests PASS.

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "feat: add types and puzzle data loader with tests"
```

---

### Task 3: Progress Store (localStorage)

**Files:**
- Create: `vimgym-web/src/lib/progress.ts`
- Create: `vimgym-web/tests/progress.test.ts`

- [ ] **Step 1: Write failing test**

`tests/progress.test.ts`:
```ts
import { describe, it, expect, beforeEach } from "vitest";
import {
  loadProgress,
  saveResult,
  getResult,
  isLevelCleared,
  isTrackCleared,
  getTotalStars,
  getTotalClears,
  resetProgress,
} from "../src/lib/progress";

beforeEach(() => {
  localStorage.clear();
});

describe("progress store", () => {
  it("returns empty state initially", () => {
    const state = loadProgress();
    expect(state.puzzles).toEqual({});
  });

  it("saves and retrieves a result", () => {
    saveResult("hjkl-01", { cleared: true, stars: 3, bestKeystrokes: 2 });
    const result = getResult("hjkl-01");
    expect(result).toEqual({ cleared: true, stars: 3, bestKeystrokes: 2 });
  });

  it("updates best result (keeps higher stars)", () => {
    saveResult("hjkl-01", { cleared: true, stars: 2, bestKeystrokes: 4 });
    saveResult("hjkl-01", { cleared: true, stars: 3, bestKeystrokes: 2 });
    expect(getResult("hjkl-01")!.stars).toBe(3);
    expect(getResult("hjkl-01")!.bestKeystrokes).toBe(2);
  });

  it("does not downgrade stars", () => {
    saveResult("hjkl-01", { cleared: true, stars: 3, bestKeystrokes: 2 });
    saveResult("hjkl-01", { cleared: true, stars: 1, bestKeystrokes: 10 });
    expect(getResult("hjkl-01")!.stars).toBe(3);
    expect(getResult("hjkl-01")!.bestKeystrokes).toBe(2);
  });

  it("returns undefined for unknown puzzle", () => {
    expect(getResult("unknown")).toBeUndefined();
  });

  it("checks level cleared (all puzzles in level cleared)", () => {
    expect(isLevelCleared(1)).toBe(false);
    // Level 1 has 6 puzzles: hjkl-01 through hjkl-06
    for (const id of ["hjkl-01", "hjkl-02", "hjkl-03", "hjkl-04", "hjkl-05", "hjkl-06"]) {
      saveResult(id, { cleared: true, stars: 1, bestKeystrokes: 10 });
    }
    expect(isLevelCleared(1)).toBe(true);
  });

  it("checks track cleared (all levels in track cleared)", () => {
    expect(isTrackCleared(1)).toBe(false);
  });

  it("counts total stars", () => {
    saveResult("hjkl-01", { cleared: true, stars: 3, bestKeystrokes: 2 });
    saveResult("hjkl-02", { cleared: true, stars: 2, bestKeystrokes: 4 });
    expect(getTotalStars()).toBe(5);
  });

  it("counts total clears", () => {
    saveResult("hjkl-01", { cleared: true, stars: 1, bestKeystrokes: 10 });
    saveResult("hjkl-02", { cleared: true, stars: 2, bestKeystrokes: 4 });
    expect(getTotalClears()).toBe(2);
  });

  it("resets progress", () => {
    saveResult("hjkl-01", { cleared: true, stars: 3, bestKeystrokes: 2 });
    resetProgress();
    expect(loadProgress().puzzles).toEqual({});
  });
});
```

- [ ] **Step 2: Run test to verify it fails**

```bash
cd vimgym-web && npx vitest run tests/progress.test.ts
```

Expected: FAIL — module not found.

- [ ] **Step 3: Implement progress store**

`src/lib/progress.ts`:
```ts
import type { GameState, PuzzleResult } from "../types";
import { getPuzzlesByLevel, getLevelsByTrack } from "../data/puzzles";

const STORAGE_KEY = "vimgym-progress";

export function loadProgress(): GameState {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) return JSON.parse(raw) as GameState;
  } catch {
    // corrupted data — start fresh
  }
  return { puzzles: {} };
}

function save(state: GameState): void {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
}

export function saveResult(puzzleId: string, result: PuzzleResult): void {
  const state = loadProgress();
  const existing = state.puzzles[puzzleId];
  if (existing) {
    // Keep best: higher stars, lower keystrokes
    if (result.stars > existing.stars || result.bestKeystrokes < existing.bestKeystrokes) {
      state.puzzles[puzzleId] = {
        cleared: true,
        stars: Math.max(existing.stars, result.stars) as 1 | 2 | 3,
        bestKeystrokes: Math.min(existing.bestKeystrokes, result.bestKeystrokes),
      };
    }
  } else {
    state.puzzles[puzzleId] = result;
  }
  save(state);
}

export function getResult(puzzleId: string): PuzzleResult | undefined {
  return loadProgress().puzzles[puzzleId];
}

export function isLevelCleared(level: number): boolean {
  const puzzles = getPuzzlesByLevel(level);
  const state = loadProgress();
  return puzzles.every((p) => state.puzzles[p.id]?.cleared);
}

export function isTrackCleared(trackId: number): boolean {
  const levels = getLevelsByTrack(trackId);
  return levels.every((l) => isLevelCleared(l));
}

export function getTotalStars(): number {
  const state = loadProgress();
  return Object.values(state.puzzles).reduce((sum, r) => sum + r.stars, 0);
}

export function getTotalClears(): number {
  const state = loadProgress();
  return Object.values(state.puzzles).filter((r) => r.cleared).length;
}

export function resetProgress(): void {
  localStorage.removeItem(STORAGE_KEY);
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
cd vimgym-web && npx vitest run tests/progress.test.ts
```

Expected: All 10 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat: add progress store with localStorage persistence"
```

---

### Task 4: Scoring Logic

**Files:**
- Create: `vimgym-web/src/lib/scoring.ts`
- Create: `vimgym-web/tests/scoring.test.ts`

- [ ] **Step 1: Write failing test**

`tests/scoring.test.ts`:
```ts
import { describe, it, expect } from "vitest";
import { calculateStars, checkClear } from "../src/lib/scoring";

describe("calculateStars", () => {
  it("returns 3 stars when keystrokes <= par", () => {
    expect(calculateStars(2, 3)).toBe(3);
    expect(calculateStars(3, 3)).toBe(3);
  });

  it("returns 2 stars when keystrokes <= par * 1.5", () => {
    expect(calculateStars(4, 3)).toBe(2);
    expect(calculateStars(4, 3)).toBe(2); // 4 <= 4.5
  });

  it("returns 1 star when cleared but above 1.5x par", () => {
    expect(calculateStars(10, 3)).toBe(1);
  });
});

describe("checkClear", () => {
  it("returns true when current text matches target", () => {
    expect(checkClear("hello world", "hello world")).toBe(true);
  });

  it("returns false when text differs", () => {
    expect(checkClear("Hello world", "hello world")).toBe(false);
  });

  it("handles multiline text", () => {
    expect(checkClear("aaa\nXbb\nccc", "aaa\nXbb\nccc")).toBe(true);
    expect(checkClear("aaa\nbbb\nccc", "aaa\nXbb\nccc")).toBe(false);
  });

  it("handles trailing newline edge case", () => {
    // CodeMirror may add trailing newline — trim both sides
    expect(checkClear("hello\n", "hello")).toBe(true);
    expect(checkClear("hello", "hello\n")).toBe(true);
  });
});
```

- [ ] **Step 2: Run test to verify it fails**

```bash
cd vimgym-web && npx vitest run tests/scoring.test.ts
```

Expected: FAIL — module not found.

- [ ] **Step 3: Implement scoring**

`src/lib/scoring.ts`:
```ts
export function calculateStars(keystrokes: number, par: number): 1 | 2 | 3 {
  if (keystrokes <= par) return 3;
  if (keystrokes <= par * 1.5) return 2;
  return 1;
}

export function checkClear(currentText: string, targetText: string): boolean {
  // Normalize trailing newlines (CodeMirror may add one)
  const normalize = (s: string) => s.replace(/\n$/, "");
  return normalize(currentText) === normalize(targetText);
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
cd vimgym-web && npx vitest run tests/scoring.test.ts
```

Expected: All 7 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "feat: add scoring logic (stars + clear judgment)"
```

---

### Task 5: useProgress Hook

**Files:**
- Create: `vimgym-web/src/hooks/useProgress.ts`

- [ ] **Step 1: Create the hook**

`src/hooks/useProgress.ts`:
```ts
import { useCallback, useSyncExternalStore } from "react";
import type { PuzzleResult } from "../types";
import {
  loadProgress,
  saveResult as saveResultToStore,
  getResult,
  isLevelCleared,
  isTrackCleared,
  getTotalStars,
  getTotalClears,
} from "../lib/progress";

let listeners: (() => void)[] = [];

function emitChange() {
  for (const listener of listeners) listener();
}

function subscribe(listener: () => void) {
  listeners = [...listeners, listener];
  return () => {
    listeners = listeners.filter((l) => l !== listener);
  };
}

function getSnapshot() {
  return loadProgress();
}

export function useProgress() {
  const state = useSyncExternalStore(subscribe, getSnapshot);

  const saveResult = useCallback((puzzleId: string, result: PuzzleResult) => {
    saveResultToStore(puzzleId, result);
    emitChange();
  }, []);

  return {
    state,
    saveResult,
    getResult,
    isLevelCleared,
    isTrackCleared,
    getTotalStars,
    getTotalClears,
  };
}
```

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "feat: add useProgress React hook"
```

---

### Task 6: Stars Component

**Files:**
- Create: `vimgym-web/src/components/Stars.tsx`

- [ ] **Step 1: Create Stars component**

`src/components/Stars.tsx`:
```tsx
interface StarsProps {
  count: 0 | 1 | 2 | 3;
  size?: "sm" | "md";
}

export default function Stars({ count, size = "sm" }: StarsProps) {
  const textSize = size === "sm" ? "text-[8px]" : "text-[12px]";
  return (
    <span className={`${textSize} tracking-wider`}>
      {[1, 2, 3].map((i) => (
        <span key={i} className={i <= count ? "text-gold" : "text-border"}>
          ★
        </span>
      ))}
    </span>
  );
}
```

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "feat: add Stars display component"
```

---

### Task 7: HUD Component

**Files:**
- Create: `vimgym-web/src/components/Hud.tsx`

- [ ] **Step 1: Create HUD component**

`src/components/Hud.tsx`:
```tsx
import { useProgress } from "../hooks/useProgress";
import { getAllPuzzles } from "../data/puzzles";

export default function Hud() {
  const { getTotalClears, getTotalStars } = useProgress();
  const total = getAllPuzzles().length;

  return (
    <div className="flex items-center justify-between px-5 py-3 bg-bg-deep border-b-[3px] border-border">
      <div className="flex flex-col items-center gap-1 font-pixel">
        <span className="text-[8px] text-gold">CLEARED</span>
        <span className="text-[11px] text-text">
          {getTotalClears()}/{total}
        </span>
      </div>
      <h1 className="text-[14px] font-pixel text-gold drop-shadow-[2px_2px_0_var(--color-gold-shadow)]">
        VIM GYM
      </h1>
      <div className="flex flex-col items-center gap-1 font-pixel">
        <span className="text-[8px] text-gold">STARS</span>
        <span className="text-[11px] text-text">{getTotalStars()}</span>
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "feat: add HUD component"
```

---

### Task 8: TrackSelect Page

**Files:**
- Create: `vimgym-web/src/pages/TrackSelect.tsx`

- [ ] **Step 1: Create TrackSelect page**

`src/pages/TrackSelect.tsx`:
```tsx
import { Link } from "react-router";
import { TRACKS } from "../types";
import { getPuzzlesByTrack } from "../data/puzzles";
import { useProgress } from "../hooks/useProgress";
import Hud from "../components/Hud";
import Stars from "../components/Stars";

export default function TrackSelect() {
  const { isTrackCleared, getResult } = useProgress();

  return (
    <div className="min-h-screen flex flex-col">
      <Hud />
      <div className="flex-1 flex flex-col items-center justify-center gap-6 p-6">
        <p className="text-[9px] font-pixel text-purple-light">
          — SELECT TRACK —
        </p>
        <div className="flex flex-col gap-4 w-full max-w-md">
          {TRACKS.map((track, i) => {
            const puzzles = getPuzzlesByTrack(track.id);
            const cleared = puzzles.filter(
              (p) => getResult(p.id)?.cleared,
            ).length;
            const prevTrackCleared = i === 0 || isTrackCleared(TRACKS[i - 1].id);
            const locked = !prevTrackCleared;

            return (
              <Link
                key={track.id}
                to={locked ? "#" : `/track/${track.id}`}
                className={`block border-[3px] rounded-lg p-5 transition-all ${
                  locked
                    ? "opacity-30 border-border cursor-not-allowed"
                    : "border-border hover:border-gold bg-bg-card"
                }`}
                onClick={(e) => locked && e.preventDefault()}
              >
                <div className="flex items-center justify-between mb-2">
                  <span className="text-[11px] font-pixel text-gold">
                    TRACK {track.id}
                  </span>
                  {locked && (
                    <span className="text-[9px] text-purple">LOCKED</span>
                  )}
                </div>
                <h2 className="text-[10px] font-pixel text-text mb-1">
                  {track.name}
                </h2>
                <p className="text-[8px] font-pixel text-purple-light mb-3">
                  {track.description}
                </p>
                <div className="h-[6px] bg-border rounded-full overflow-hidden mb-2">
                  <div
                    className="h-full bg-gold rounded-full transition-all"
                    style={{ width: `${(cleared / puzzles.length) * 100}%` }}
                  />
                </div>
                <div className="flex justify-between text-[7px] font-pixel text-purple">
                  <span>
                    {cleared}/{puzzles.length} PUZZLES
                  </span>
                  {cleared === puzzles.length && <Stars count={3} />}
                </div>
              </Link>
            );
          })}
        </div>
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Verify in browser**

```bash
cd vimgym-web && npm run dev -- --port 5180
```

Open http://localhost:5180 and confirm TrackSelect renders.

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "feat: add TrackSelect page with track cards"
```

---

### Task 9: LevelSelect Page

**Files:**
- Create: `vimgym-web/src/pages/LevelSelect.tsx`

- [ ] **Step 1: Create LevelSelect page**

`src/pages/LevelSelect.tsx`:
```tsx
import { Link, useParams } from "react-router";
import { TRACKS } from "../types";
import { getPuzzlesByLevel, getLevelsByTrack } from "../data/puzzles";
import { useProgress } from "../hooks/useProgress";
import Hud from "../components/Hud";
import Stars from "../components/Stars";

export default function LevelSelect() {
  const { trackId } = useParams<{ trackId: string }>();
  const track = TRACKS.find((t) => t.id === Number(trackId));
  const levels = track ? getLevelsByTrack(track.id) : [];
  const { isLevelCleared, getResult } = useProgress();

  if (!track) return <div className="p-6 font-pixel text-gold">NOT FOUND</div>;

  return (
    <div className="min-h-screen flex flex-col">
      <Hud />
      <div className="flex-1 flex flex-col items-center p-6 gap-6">
        <Link
          to="/"
          className="text-[8px] font-pixel text-purple hover:text-gold self-start"
        >
          ← BACK
        </Link>
        <p className="text-[9px] font-pixel text-purple-light">
          — {track.name.toUpperCase()} —
        </p>
        <div className="grid grid-cols-5 gap-3 w-full max-w-md">
          {levels.map((level, i) => {
            const puzzles = getPuzzlesByLevel(level);
            const cleared = puzzles.filter(
              (p) => getResult(p.id)?.cleared,
            ).length;
            const allCleared = cleared === puzzles.length;
            const prevCleared = i === 0 || isLevelCleared(levels[i - 1]);
            const locked = !prevCleared;
            const bestStars = allCleared
              ? Math.min(
                  ...puzzles.map((p) => getResult(p.id)?.stars ?? 0),
                ) as 0 | 1 | 2 | 3
              : (0 as const);

            return (
              <Link
                key={level}
                to={locked ? "#" : `/track/${track.id}/level/${level}`}
                className={`aspect-square flex flex-col items-center justify-center gap-1 border-[3px] rounded-lg font-pixel text-[12px] transition-all ${
                  locked
                    ? "opacity-30 border-border cursor-not-allowed"
                    : allCleared
                      ? "border-green bg-bg-card"
                      : "border-border hover:border-gold bg-bg-card"
                }`}
                onClick={(e) => locked && e.preventDefault()}
              >
                <span className="text-text">{level}</span>
                <Stars count={bestStars} />
              </Link>
            );
          })}
        </div>
        <p className="text-[7px] font-pixel text-purple">
          {track.description}
        </p>
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "feat: add LevelSelect page with 5-col grid"
```

---

### Task 10: PuzzleList Page

**Files:**
- Create: `vimgym-web/src/pages/PuzzleList.tsx`

- [ ] **Step 1: Create PuzzleList page**

`src/pages/PuzzleList.tsx`:
```tsx
import { Link, useParams } from "react-router";
import { TRACKS } from "../types";
import { getPuzzlesByLevel } from "../data/puzzles";
import { useProgress } from "../hooks/useProgress";
import Stars from "../components/Stars";
import Hud from "../components/Hud";

export default function PuzzleList() {
  const { trackId, levelId } = useParams<{
    trackId: string;
    levelId: string;
  }>();
  const track = TRACKS.find((t) => t.id === Number(trackId));
  const level = Number(levelId);
  const puzzles = getPuzzlesByLevel(level);
  const { getResult } = useProgress();

  if (!track || puzzles.length === 0) {
    return <div className="p-6 font-pixel text-gold">NOT FOUND</div>;
  }

  return (
    <div className="min-h-screen flex flex-col">
      <Hud />
      <div className="flex-1 flex flex-col items-center p-6 gap-6">
        <Link
          to={`/track/${track.id}`}
          className="text-[8px] font-pixel text-purple hover:text-gold self-start"
        >
          ← BACK
        </Link>
        <p className="text-[9px] font-pixel text-purple-light">
          — LEVEL {level} —
        </p>
        <div className="flex flex-col gap-2 w-full max-w-md">
          {puzzles.map((puzzle, i) => {
            const result = getResult(puzzle.id);
            const prevCleared =
              i === 0 || getResult(puzzles[i - 1].id)?.cleared;
            const locked = !prevCleared;

            return (
              <Link
                key={puzzle.id}
                to={
                  locked
                    ? "#"
                    : `/track/${track.id}/level/${level}/puzzle/${puzzle.id}`
                }
                className={`flex items-center justify-between p-4 border-[3px] rounded-lg font-pixel transition-all ${
                  locked
                    ? "opacity-30 border-border cursor-not-allowed"
                    : result?.cleared
                      ? "border-green bg-bg-card"
                      : "border-border hover:border-gold bg-bg-card"
                }`}
                onClick={(e) => locked && e.preventDefault()}
              >
                <div className="flex flex-col gap-1">
                  <span className="text-[10px] text-text">{puzzle.title}</span>
                  <span className="text-[7px] text-purple">
                    {puzzle.category} · PAR {puzzle.par}
                  </span>
                </div>
                <div className="flex items-center gap-3">
                  {result && (
                    <span className="text-[7px] text-purple-light">
                      {result.bestKeystrokes} keys
                    </span>
                  )}
                  <Stars count={result?.cleared ? result.stars : 0} />
                </div>
              </Link>
            );
          })}
        </div>
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "feat: add PuzzleList page"
```

---

### Task 11: VimEditor Component

**Files:**
- Create: `vimgym-web/src/components/VimEditor.tsx`

- [ ] **Step 1: Create VimEditor component**

`src/components/VimEditor.tsx`:
```tsx
import { useEffect, useRef, useCallback } from "react";
import { EditorView, keymap } from "@codemirror/view";
import { EditorState } from "@codemirror/state";
import { defaultKeymap } from "@codemirror/commands";
import { vim } from "@replit/codemirror-vim";

interface VimEditorProps {
  initialText: string;
  cursorRow: number;
  cursorCol: number;
  onChange: (text: string) => void;
  onKeystroke: () => void;
}

export default function VimEditor({
  initialText,
  cursorRow,
  cursorCol,
  onChange,
  onKeystroke,
}: VimEditorProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const viewRef = useRef<EditorView | null>(null);

  const onChangeRef = useRef(onChange);
  onChangeRef.current = onChange;
  const onKeystrokeRef = useRef(onKeystroke);
  onKeystrokeRef.current = onKeystroke;

  const reset = useCallback(() => {
    if (!viewRef.current) return;
    viewRef.current.dispatch({
      changes: {
        from: 0,
        to: viewRef.current.state.doc.length,
        insert: initialText,
      },
    });
    // Set cursor position
    const line = viewRef.current.state.doc.line(cursorRow + 1);
    const pos = line.from + cursorCol;
    viewRef.current.dispatch({
      selection: { anchor: pos },
    });
    viewRef.current.focus();
  }, [initialText, cursorRow, cursorCol]);

  useEffect(() => {
    if (!containerRef.current) return;

    const keystrokeListener = EditorView.domEventHandlers({
      keydown: () => {
        onKeystrokeRef.current();
        return false; // don't prevent default
      },
    });

    const state = EditorState.create({
      doc: initialText,
      extensions: [
        vim(),
        keymap.of(defaultKeymap),
        keystrokeListener,
        EditorView.updateListener.of((update) => {
          if (update.docChanged) {
            onChangeRef.current(update.state.doc.toString());
          }
        }),
      ],
    });

    const view = new EditorView({
      state,
      parent: containerRef.current,
    });

    viewRef.current = view;

    // Set initial cursor position
    const line = view.state.doc.line(cursorRow + 1);
    const pos = line.from + cursorCol;
    view.dispatch({ selection: { anchor: pos } });
    view.focus();

    return () => {
      view.destroy();
    };
  }, [initialText, cursorRow, cursorCol]);

  return (
    <div>
      <div ref={containerRef} className="min-h-[200px]" />
      <button
        onClick={reset}
        className="mt-2 px-4 py-2 text-[8px] font-pixel text-purple border-2 border-border rounded hover:border-gold hover:text-gold transition-colors"
      >
        RESET
      </button>
    </div>
  );
}
```

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "feat: add VimEditor component with CodeMirror + Vim mode"
```

---

### Task 12: PuzzlePlay Page

**Files:**
- Create: `vimgym-web/src/pages/PuzzlePlay.tsx`

- [ ] **Step 1: Create PuzzlePlay page**

`src/pages/PuzzlePlay.tsx`:
```tsx
import { useState, useCallback, useRef } from "react";
import { useParams, useNavigate, Link } from "react-router";
import { getPuzzleById, getPuzzlesByLevel } from "../data/puzzles";
import { useProgress } from "../hooks/useProgress";
import { calculateStars, checkClear } from "../lib/scoring";
import VimEditor from "../components/VimEditor";
import Stars from "../components/Stars";
import Hud from "../components/Hud";

export default function PuzzlePlay() {
  const { trackId, levelId, puzzleId } = useParams<{
    trackId: string;
    levelId: string;
    puzzleId: string;
  }>();
  const navigate = useNavigate();
  const puzzle = getPuzzleById(puzzleId!);
  const { saveResult } = useProgress();

  const [keystrokes, setKeystrokes] = useState(0);
  const [showHint, setShowHint] = useState(false);
  const [clearResult, setClearResult] = useState<{
    stars: 1 | 2 | 3;
    keystrokes: number;
  } | null>(null);

  const clearedRef = useRef(false);

  const handleChange = useCallback(
    (text: string) => {
      if (!puzzle || clearedRef.current) return;
      if (checkClear(text, puzzle.after.text)) {
        clearedRef.current = true;
        // +1 because the keystroke counter increments after this
        const finalKeystrokes = keystrokes + 1;
        const stars = calculateStars(finalKeystrokes, puzzle.par);
        setClearResult({ stars, keystrokes: finalKeystrokes });
        saveResult(puzzle.id, {
          cleared: true,
          stars,
          bestKeystrokes: finalKeystrokes,
        });
      }
    },
    [puzzle, keystrokes, saveResult],
  );

  const handleKeystroke = useCallback(() => {
    if (!clearedRef.current) {
      setKeystrokes((k) => k + 1);
    }
  }, []);

  if (!puzzle) {
    return <div className="p-6 font-pixel text-gold">PUZZLE NOT FOUND</div>;
  }

  // Find next puzzle
  const levelPuzzles = getPuzzlesByLevel(puzzle.level);
  const currentIndex = levelPuzzles.findIndex((p) => p.id === puzzle.id);
  const nextPuzzle = levelPuzzles[currentIndex + 1];

  const handleRetry = () => {
    clearedRef.current = false;
    setKeystrokes(0);
    setShowHint(false);
    setClearResult(null);
    // Force re-mount of editor via key change
    navigate(0);
  };

  return (
    <div className="min-h-screen flex flex-col">
      <Hud />
      <div className="flex-1 flex flex-col p-6 gap-4 max-w-2xl mx-auto w-full">
        {/* Navigation */}
        <Link
          to={`/track/${trackId}/level/${levelId}`}
          className="text-[8px] font-pixel text-purple hover:text-gold self-start"
        >
          ← BACK
        </Link>

        {/* Quest header */}
        <div className="flex items-center justify-between border-b-[3px] border-border pb-3">
          <div>
            <span className="text-[8px] font-pixel text-gold">QUEST</span>
            <h2 className="text-[10px] font-pixel text-text mt-1">
              {puzzle.title}
            </h2>
          </div>
          <div className="flex gap-6 text-[8px] font-pixel">
            <div className="flex flex-col items-center gap-1">
              <span className="text-gold">PAR</span>
              <span className="text-green">{puzzle.par}</span>
            </div>
            <div className="flex flex-col items-center gap-1">
              <span className="text-gold">KEYS</span>
              <span className="text-text">{keystrokes}</span>
            </div>
          </div>
        </div>

        {/* Editor */}
        <VimEditor
          initialText={puzzle.before.text}
          cursorRow={puzzle.before.cursor.row}
          cursorCol={puzzle.before.cursor.col}
          onChange={handleChange}
          onKeystroke={handleKeystroke}
        />

        {/* Target display */}
        <div className="bg-bg-deep border-2 border-border rounded-lg p-4">
          <span className="text-[8px] font-pixel text-gold block mb-2">
            GOAL
          </span>
          <pre className="font-code text-[13px] text-purple-text whitespace-pre-wrap">
            {puzzle.after.text}
          </pre>
        </div>

        {/* Hint */}
        <div className="flex gap-3">
          <button
            onClick={() => setShowHint(!showHint)}
            className="px-4 py-2 text-[8px] font-pixel text-purple border-2 border-border rounded hover:border-gold hover:text-gold transition-colors"
          >
            {showHint ? "HIDE HINT" : "HINT"}
          </button>
        </div>
        {showHint && (
          <div className="bg-bg-card border-2 border-border rounded-lg p-4">
            <p className="text-[9px] font-pixel text-purple-light leading-relaxed">
              {puzzle.hint}
            </p>
          </div>
        )}

        {/* Clear overlay */}
        {clearResult && (
          <div className="fixed inset-0 bg-black/80 flex items-center justify-center z-50">
            <div className="bg-bg-card border-[3px] border-gold rounded-xl p-8 flex flex-col items-center gap-4 max-w-sm mx-4">
              <span className="text-[12px] font-pixel text-gold">
                CLEAR!
              </span>
              <Stars count={clearResult.stars} size="md" />
              <div className="text-[8px] font-pixel text-purple-light text-center space-y-1">
                <p>
                  KEYSTROKES: {clearResult.keystrokes} / PAR: {puzzle.par}
                </p>
              </div>
              <div className="flex gap-3 mt-2">
                <button
                  onClick={handleRetry}
                  className="px-4 py-2 text-[8px] font-pixel text-purple border-2 border-border rounded hover:border-gold hover:text-gold transition-colors"
                >
                  RETRY
                </button>
                {nextPuzzle ? (
                  <button
                    onClick={() =>
                      navigate(
                        `/track/${trackId}/level/${levelId}/puzzle/${nextPuzzle.id}`,
                      )
                    }
                    className="px-4 py-2 text-[8px] font-pixel text-bg-deep bg-gold border-2 border-gold rounded hover:bg-green hover:border-green transition-colors"
                  >
                    NEXT →
                  </button>
                ) : (
                  <button
                    onClick={() =>
                      navigate(`/track/${trackId}/level/${levelId}`)
                    }
                    className="px-4 py-2 text-[8px] font-pixel text-bg-deep bg-gold border-2 border-gold rounded hover:bg-green hover:border-green transition-colors"
                  >
                    DONE
                  </button>
                )}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "feat: add PuzzlePlay page with clear judgment and overlay"
```

---

### Task 13: Routing

**Files:**
- Modify: `vimgym-web/src/App.tsx`

- [ ] **Step 1: Wire up all routes with HashRouter**

`src/App.tsx`:
```tsx
import { HashRouter, Routes, Route } from "react-router";
import TrackSelect from "./pages/TrackSelect";
import LevelSelect from "./pages/LevelSelect";
import PuzzleList from "./pages/PuzzleList";
import PuzzlePlay from "./pages/PuzzlePlay";

export default function App() {
  return (
    <HashRouter>
      <Routes>
        <Route path="/" element={<TrackSelect />} />
        <Route path="/track/:trackId" element={<LevelSelect />} />
        <Route path="/track/:trackId/level/:levelId" element={<PuzzleList />} />
        <Route
          path="/track/:trackId/level/:levelId/puzzle/:puzzleId"
          element={<PuzzlePlay />}
        />
      </Routes>
    </HashRouter>
  );
}
```

- [ ] **Step 2: Verify full flow in browser**

```bash
cd vimgym-web && npm run dev -- --port 5180
```

Navigate through: TrackSelect → LevelSelect → PuzzleList → PuzzlePlay. Confirm:
1. Track 1 is unlocked, Track 2/3 are locked
2. Level 1 is unlocked in Track 1
3. First puzzle in Level 1 is unlocked
4. Vim editor loads with correct before text
5. Typing Vim commands works
6. Clear judgment triggers when matching after text

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "feat: add HashRouter with all page routes"
```

---

### Task 14: GitHub Pages Deployment

**Files:**
- Create: `vimgym-web/.github/workflows/deploy.yml`
- Modify: `vimgym-web/package.json` (add deploy script)

- [ ] **Step 1: Create GitHub Actions workflow**

`.github/workflows/deploy.yml`:
```yaml
name: Deploy to GitHub Pages

on:
  push:
    branches: [main]

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: false

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 22
          cache: npm
      - run: npm ci
      - run: npm run build
      - uses: actions/upload-pages-artifact@v3
        with:
          path: dist

  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    needs: build
    steps:
      - uses: actions/deploy-pages@v4
        id: deployment
```

- [ ] **Step 2: Add 404.html for SPA routing on GitHub Pages**

Since we use HashRouter, this isn't strictly needed, but add a safety net.

Create `vimgym-web/public/404.html`:
```html
<!doctype html>
<html>
  <head>
    <script>
      // Redirect to index with hash
      window.location.href =
        window.location.origin +
        "/vimgym-web/" +
        "#" +
        window.location.pathname.replace("/vimgym-web/", "/");
    </script>
  </head>
</html>
```

- [ ] **Step 3: Build and verify**

```bash
cd vimgym-web && npm run build
ls dist/
```

Expected: `dist/` contains `index.html`, `assets/`, `404.html`.

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "feat: add GitHub Pages deployment workflow"
```

---

### Task 15: Run All Tests and Final Verification

- [ ] **Step 1: Run all tests**

```bash
cd vimgym-web && npx vitest run
```

Expected: All tests pass (puzzles, progress, scoring).

- [ ] **Step 2: Run build**

```bash
npm run build
```

Expected: Build succeeds with no errors.

- [ ] **Step 3: Preview production build**

```bash
npm run preview -- --port 5181
```

Open http://localhost:5181/vimgym-web/ and run through the full flow.

- [ ] **Step 4: Final commit**

```bash
git add -A
git commit -m "chore: verify all tests pass and production build works"
```
