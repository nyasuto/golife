# Legacy Package

既存の `[][]int` ベースのライフゲーム実装との互換性を提供するパッケージ。

## 概要

このパッケージは2つの主要コンポーネントを提供します：

1. **ClassicGrid** - 従来の `[][]int` 構造によるライフゲーム実装
2. **Adapter** - ClassicGrid と新しい Universe システム間の変換

## ClassicGrid

従来の実装を維持したクラシックなライフゲーム。

### 使用例

```go
import "golife/pkg/legacy"

// グリッド作成
g := legacy.NewClassicGrid(100, 100)

// ランダム初期化
g.Randomize()

// セルの設定
g.Set(10, 10, 1)

// 1世代進める
g = g.Step()

// 生存セル数
count := g.CountLiving()
```

### 主要メソッド

- `NewClassicGrid(width, height int)` - 新しいグリッド作成
- `Randomize()` - ランダムに初期化
- `Step()` - 1世代進める
- `CountNeighbors(x, y int)` - 近傍カウント
- `CountLiving()` - 生存セル数
- `Clone()` - ディープコピー

## Adapter

新旧システム間の変換を提供。

### 使用例

```go
import (
    "golife/pkg/legacy"
    "golife/pkg/universe"
    "golife/pkg/rules"
)

// ClassicGrid → Universe2D
classicGrid := legacy.NewClassicGrid(100, 100)
classicGrid.Randomize()
u := classicGrid.ToUniverse()

// Universe2D → ClassicGrid
u2 := universe.New2D(100, 100, rules.ConwayRule{})
g := legacy.NewFromUniverse(u2)

// [][]int → ClassicGrid
cells := [][]int{
    {0, 1, 0},
    {0, 0, 1},
    {1, 1, 1},
}
g2 := legacy.FromSlice(cells)

// ClassicGrid → [][]int
slice := g2.ToSlice()
```

## パフォーマンス比較

ベンチマーク結果（100x100グリッド）:

```
BenchmarkClassicGrid_Step    101509 ns/op    92336 B/op    102 allocs/op
BenchmarkUniverse2D_Step      97610 ns/op        0 B/op      0 allocs/op
```

**Universe2D の利点:**
- **4%高速** - ダブルバッファリングによる最適化
- **メモリ割り当てゼロ** - フラット配列の再利用

## 互換性保証

- 既存の `[][]int` ベースのコードは ClassicGrid で動作
- Conway B3/S23 ルールの完全互換
- 全ての有名パターン（Glider, Blinker等）が正しく動作

## テスト

```bash
go test ./pkg/legacy/... -v
go test ./pkg/legacy/... -bench=. -benchmem
```

## 使用シーン

### 1. 既存コードの移行

```go
// 既存コード
grid := make([][]int, height)
// ...

// 新システムへの移行
classicGrid := legacy.FromSlice(grid)
universe := classicGrid.ToUniverse()
```

### 2. パフォーマンス比較

```go
// 旧実装
classic := legacy.NewClassicGrid(100, 100)
classic = classic.Step()

// 新実装
modern := universe.New2D(100, 100, rules.ConwayRule{})
modern.Step()
```

### 3. パターンの互換性確認

```go
// パターンをロード
u := universe.New2D(50, 50, rules.ConwayRule{})
glider := patterns.Glider()
glider.LoadIntoUniverse(u, 10, 10)

// 旧形式に変換して検証
classic := legacy.NewFromUniverse(u)
classic = classic.Step()
```

## 将来の展開

このパッケージは主に互換性のために存在します。新しいコードでは `pkg/universe` と `pkg/rules` の使用を推奨します。

多次元実装（3D/4D）では ClassicGrid は使用せず、直接 Universe インターフェースを使用してください。
