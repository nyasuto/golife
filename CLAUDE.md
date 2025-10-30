# Go Life Development Policy

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