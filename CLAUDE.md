# Go Life Development Policy

## Project Overview

Conway's Game of Life implementation in Go with:
- Terminal UI using termbox-go
- Configurable grid size, speed, and generation count
- Famous pattern presets (Glider, Pulsar, Gosper's Glider Gun, etc.)
- 72.4% test coverage
- CI/CD pipeline with automated testing and quality checks

## 🚨 CRITICAL: PR Merge Policy

**ABSOLUTE RULE: Claude MUST NEVER merge PRs automatically**

- ✅ **ALLOWED**: Create PRs using `gh pr create`
- ✅ **ALLOWED**: Watch PR checks with `gh pr checks --watch`
- ❌ **FORBIDDEN**: Use `gh pr merge` or any merge commands
- ❌ **FORBIDDEN**: Automatic merging regardless of PR size or type
- ❌ **FORBIDDEN**: Merging even if all CI checks pass
- ✅ **REQUIRED**: Human must review and merge ALL PRs manually

**WHY**: Human review is essential for quality control, architectural decisions, and understanding changes.

## Git Workflow

### mainブランチへの直接コミット禁止

- ❌ mainブランチへの直接コミットは禁止
- ✅ 全ての変更はfeatureブランチから開始
- ✅ Pull Requestを経由してマージ

### ブランチ命名規則

- `feat/X-description` - 新機能
- `fix/X-description` - バグ修正
- `docs/X-description` - ドキュメント
- `ci/X-description` - CI/CD設定
- `refactor/X-description` - リファクタリング

### 開発フロー

```bash
# 1. featureブランチを作成
git checkout -b feat/issue-X-feature-name

# 2. 変更を実装

# 3. ローカルでテスト
make quality

# 4. コミット
git add .
git commit -m "feat: 機能の説明"

# 5. プッシュ
git push origin feat/issue-X-feature-name

# 6. Pull Request作成
gh pr create --title "feat: 機能の説明" --body "詳細..."

# 7. CI/CD通過を確認
gh pr checks --watch

# 8. ⚠️ 人間によるレビューとマージを待つ
```

## Implemented Features

### ✅ Completed Issues

1. **Issue #10: Makefile** (PR #20)
   - Comprehensive build system with quality checks
   - Targets: build, test, coverage, quality, clean, run

2. **Issue #11: Unit Tests** (PR #22)
   - 72.4% test coverage
   - 6 test suites, 15 subtests
   - Tests for randomize(), step(), edge cases, known patterns

3. **Issue #12: CI/CD Pipeline** (PR #23)
   - Multi-job workflow: Lint, Test, Build, Quality Checks
   - Codecov integration
   - Go 1.25 support

4. **Issue #13: Code Refactoring** (PR #30)
   - Extracted countNeighbors() helper function
   - Reduced step() from 70 lines to 24 lines (66% reduction)
   - Improved code readability

5. **Issue #14: Configurable Parameters** (PR #31)
   - Command-line flags: --width, --height, --speed, --generations
   - Input validation
   - Default values with constants

6. **Issue #16: Pattern Presets** (PR #32)
   - 6 famous patterns: glider, blinker, toad, beacon, pulsar, glider-gun
   - Pattern loading with center alignment
   - --pattern flag with 'list' option

### 🔄 Pending Issues

7. **Issue #15: Interactive Mode**
   - Keyboard controls for pause/resume, step, speed adjustment
   - Priority: medium

8. **Issue #17: File I/O**
   - Save/load grid state to files
   - Priority: low

9. **Issue #18: Statistics Display**
   - Show generation count, alive cells, population changes
   - Priority: low

10. **Issue #19: Colorful Display**
    - Color cells based on age
    - Priority: low

## Technical Stack

- **Language**: Go 1.25
- **UI Library**: termbox-go
- **Testing**: Go standard testing package
- **CI/CD**: GitHub Actions
- **Linting**: golangci-lint v2 (action v8)
- **Coverage**: Codecov

## Code Quality Standards

- All tests must pass
- golangci-lint must pass
- go fmt, go vet, go mod tidy checks
- Conventional commit messages
- Minimum 70% test coverage for new features