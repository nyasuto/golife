# Multi-Dimensional Life Game - Design Document

## プロジェクト概要

Go言語で実装された2Dライフゲームを、多次元（2.5D → 3D → 4D）に拡張するプロジェクト。

## アーキテクチャ

### コア抽象化

```
pkg/core/
  types.go        - 次元非依存のインターフェース定義
```

#### 主要インターフェース

**Universe** - 全次元対応の宇宙インターフェース
```go
type Universe interface {
    Dimension() Dimension
    Get(coord Coord) CellState
    Set(coord Coord, state CellState)
    Step()
    Size() Coord
    Clone() Universe
    Clear()
    CountLiving() int
}
```

**Rule** - ルール定義インターフェース
```go
type Rule interface {
    Name() string
    ShouldBirth(neighborCount int) bool
    ShouldSurvive(neighborCount int, currentState CellState) bool
    NeighborWeight(distance float64) float64
}
```

**Coord** - 次元非依存の座標
```go
type Coord struct {
    X, Y, Z, W int  // 2Dではz,w=0, 3Dではw=0
}
```

### パッケージ構造

```
golife/
├── pkg/
│   ├── core/         # コアインターフェース
│   ├── universe/     # Universe実装 (2D, 2.5D, 3D, 4D)
│   ├── rules/        # ルール実装
│   ├── engine/       # シミュレーションエンジン
│   ├── patterns/     # 既知パターン
│   ├── visualizer/   # 可視化
│   └── legacy/       # 既存コード互換層
├── cmd/
│   ├── golife/       # 既存2D版（維持）
│   └── golife-nd/    # 新・多次元版
└── docs/             # ドキュメント
```

## 設計原則

### 1. 次元非依存性

全てのUniverseは同じインターフェースを実装し、次元に関わらず同じコードで操作可能。

```go
var u core.Universe
u = universe.New2D(100, 100, rules.ConwayRule{})
// または
u = universe.New3D(64, 64, 64, rules.Life3D_B6S567{})
```

### 2. 後方互換性

既存の2Dコードは`pkg/legacy`に移動し、新システムとのアダプターを提供。

### 3. パフォーマンス

- フラット配列でキャッシュ効率向上
- ダブルバッファリング
- goroutineによる並列化（Phase 2以降）

### 4. 拡張性

新しいルールや次元の追加が容易な設計。

## 実装状況

### Phase 0: アーキテクチャ設計 ✅

- [x] コアインターフェース定義
- [x] ディレクトリ構造再編成
- [x] Universe2D実装
- [x] ConwayRule実装
- [x] パターンライブラリ

### Phase 1: 2.5D実装 (予定)

- [ ] Universe25D実装
- [ ] 層間相互作用ルール
- [ ] 複数層可視化

### Phase 2: 3D実装 (予定)

- [ ] Universe3D実装
- [ ] B6/S567ルール
- [ ] 3Dグライダー
- [ ] 並列処理エンジン

### Phase 3: 4D実装 (予定)

- [ ] Universe4D実装
- [ ] B9/S7-10ルール
- [ ] スパース実装
- [ ] 4D振動子探索

## 技術的課題と解決策

### 課題1: 近傍数の爆発的増加

| 次元 | 近傍数 (Moore) |
|-----|---------------|
| 2D  | 8             |
| 3D  | 26            |
| 4D  | 80            |

**解決策:**
- ルールのスケーリング（B6/S567, B9/S7-10）
- 距離減衰型ルール
- エネルギー保存則ルール

### 課題2: メモリ使用量

| サイズ    | メモリ (uint8) |
|----------|---------------|
| 100²     | 10 KB         |
| 64³      | 256 KB        |
| 32⁴      | 64 MB         |

**解決策:**
- スパース実装（チャンクマップ）
- 生存セル周辺のみ処理

### 課題3: 可視化

**解決策:**
- 2.5D/3D: スライス表示
- 4D: W軸スライス or 4D→3D投影

## パフォーマンス目標

| 次元 | サイズ  | 目標FPS | 最適化手法 |
|-----|--------|---------|-----------|
| 2D  | 100²   | 60 FPS  | フラット配列 |
| 3D  | 128³   | 10 FPS  | 並列処理 |
| 4D  | 64⁴    | 1 FPS   | スパース |

## 参考文献

- Carter Bays (1987): "Candidates for the Game of Life in Three Dimensions"
- "Higher Dimensional Games of Life" (ResearchGate)
- ConwayLife Wiki: https://conwaylife.com/wiki/
