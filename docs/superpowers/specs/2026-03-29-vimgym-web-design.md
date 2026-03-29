# VimGym Web — Design Spec

## Overview

기존 Go TUI 기반 VimGym의 Track 1-3 (86개 퍼즐)을 레트로 게임 UI 웹 버전으로 재구현한다. 서버 없이 정적 사이트로 GitHub Pages에 배포한다.

## Tech Stack

- Vite + React 19 + TypeScript
- Tailwind CSS v4
- CodeMirror 6 + @replit/codemirror-vim
- localStorage (진행도 저장)
- GitHub Pages (배포)

## Design Concept

레트로 아케이드 게임 UI를 기반으로 한다.

- **폰트**: Press Start 2P (UI), JetBrains Mono (코드 에디터)
- **팔레트**: 짙은 보라 배경 (#1a1028), 금색 액센트 (#f0c040), 보라 보조 (#3a2860), 초록 클리어 (#40c060)
- **다크 테마 고정**: 게임 분위기 유지를 위해 시스템 테마 미추종
- **픽셀 느낌**: 날카로운 모서리, 두꺼운 보더, 아케이드 HUD

## Screens

### 1. 메인 화면 (트랙 선택)

- **HUD 바**: 총 클리어 수, 총 별 수
- **트랙 카드 3개**: Foundations / Editing / Power Moves
  - 트랙 이름, 퍼즐 수, 진행률 바
  - 잠금 상태: 이전 트랙의 모든 레벨을 1개 이상 별로 클리어해야 다음 트랙 해금
- **트랙 잠금 표시**: 반투명 + 자물쇠 아이콘

### 2. 레벨 셀렉트 화면

- **5열 그리드**로 해당 트랙의 레벨들 표시
  - Track 1: Level 1-5 (29개 퍼즐이지만 5개 레벨로 그룹핑)
  - Track 2: Level 6-10 (28개 퍼즐, 5개 레벨)
  - Track 3: Level 11-15 (29개 퍼즐, 5개 레벨)
- **각 셀**: 레벨 번호 + 획득 별점 (★★★)
- **순차 잠금**: 이전 레벨의 모든 퍼즐을 클리어해야 다음 레벨 해금
- **상태 표현**:
  - 클리어: 초록 테두리 (#40c060)
  - 현재 도전 가능: 금색 테두리 (#f0c040)
  - 잠금: 반투명 (opacity 0.3)

### 3. 퍼즐 리스트 화면

레벨을 선택하면 해당 레벨의 퍼즐 목록을 보여준다.

- 퍼즐별: 이름, 태그, par, 획득 별점
- 순차 잠금: 이전 퍼즐 클리어해야 다음 퍼즐 해금

### 4. 퍼즐 플레이 화면

- **상단 HUD**: 퀘스트 설명, par 목표, 키스트로크 카운터 (실시간)
- **중앙**: CodeMirror 에디터 (Vim 모드)
  - before 텍스트를 에디터에 로드
  - 커서 위치도 퍼즐 데이터에 따라 설정
  - Vim 명령으로 직접 편집
- **하단**: 힌트 버튼 (기본 숨김, 누르면 표시), 리셋 버튼
- **목표 표시**: after 텍스트를 에디터 아래에 참고용으로 표시
- **클리어 판정**: 에디터 내용이 after 텍스트와 정확히 일치하면 자동 클리어

### 5. 클리어 화면

- **별점 결과**: par 이하 ★★★, 1.5배 이하 ★★, 클리어 ★
- **통계**: 사용한 키스트로크 수 vs par
- **버튼**: "다음 퍼즐" / "다시 도전" / "레벨로 돌아가기"

## Gamification

### 별점 시스템

| 조건 | 별점 |
|------|------|
| keystrokes ≤ par | ★★★ |
| keystrokes ≤ par × 1.5 | ★★ |
| 클리어 | ★ |

### 진행 잠금

- **퍼즐 잠금**: 이전 퍼즐 클리어 → 다음 퍼즐 해금
- **레벨 잠금**: 이전 레벨의 모든 퍼즐 클리어 → 다음 레벨 해금
- **트랙 잠금**: 이전 트랙의 모든 레벨 클리어 → 다음 트랙 해금

### 키스트로크 카운팅

모든 키 입력을 카운트한다 (기존 TUI 방식 동일). Insert 모드 내 타이핑도 포함.

## Puzzle Data

기존 JSON 파일 3개를 그대로 사용한다:

- `puzzles/track1_foundations.json` — 29 puzzles (Levels 1-5)
- `puzzles/track2_editing.json` — 28 puzzles (Levels 6-10)
- `puzzles/track3_power.json` — 29 puzzles (Levels 11-15)

### 퍼즐 구조

```json
{
  "id": "hjkl-01",
  "title": "Move Right",
  "track": 1,
  "level": 1,
  "category": "hjkl",
  "difficulty": 1,
  "before": { "text": "Hello World", "cursor": { "row": 0, "col": 0 } },
  "after": { "text": "Hello World" },
  "par": 5,
  "hint": "Use 'l' to move right",
  "optimalSolution": "lllll",
  "tags": ["hjkl"]
}
```

참고: after에 cursor가 없는 퍼즐은 텍스트 일치만으로 클리어 판정한다.

## Data Persistence (localStorage)

```typescript
interface GameState {
  puzzles: Record<string, PuzzleResult>;  // puzzleId → result
  currentTrack: number;
  currentLevel: number;
}

interface PuzzleResult {
  cleared: boolean;
  stars: 1 | 2 | 3;
  bestKeystrokes: number;
}
```

키: `vimgym-progress`

## Clear Judgment

1. 매 키 입력마다 에디터 텍스트를 after.text와 비교
2. 커서 위치도 after.cursor와 비교 (퍼즐에 커서 조건이 있는 경우)
3. 두 조건 모두 일치하면 클리어 판정
4. 키스트로크 수로 별점 산정

## Routing

React Router로 화면 전환:

- `/` — 메인 (트랙 선택)
- `/track/:trackId` — 레벨 셀렉트
- `/track/:trackId/level/:levelId` — 퍼즐 리스트
- `/track/:trackId/level/:levelId/puzzle/:puzzleId` — 퍼즐 플레이

GitHub Pages 배포를 위해 HashRouter 사용.

## Out of Scope

- 계정/로그인 시스템
- 리더보드/랭킹
- Track 4 (Vim Golf) — 고급 Vim 기능 필요
- 모바일 대응 — Vim 에디터 특성상 데스크톱 전용
- 소리/음악
