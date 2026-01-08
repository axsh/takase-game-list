# takase-game-list

## ビルド方法

本プロジェクトは、Windows環境でGo 1.22以上を使用してビルド・テストを実行します。

### 前提条件

- Go 1.22以上がインストールされていること
- PowerShellが利用可能であること

### ビルドスクリプトの実行

プロジェクトルートに配置された `build.ps1` を使用して、ビルドとテストを実行します。

#### 実行モード

ビルドスクリプトは以下の実行モードをサポートしています：

- `all`（デフォルト）: ビルド、ユニットテスト、結合テストを順次実行
- `unit`: ビルドとユニットテストのみ実行
- `integration`: ビルドと結合テストのみ実行
- `build-only`: ビルドのみ実行（テストは実行しない）

#### 実行例

```powershell
# デフォルト（allモード）で実行
powershell -NoProfile -ExecutionPolicy Bypass -File ./build.ps1

# ユニットテストのみ実行
powershell -NoProfile -ExecutionPolicy Bypass -File ./build.ps1 -Mode unit

# 結合テストのみ実行
powershell -NoProfile -ExecutionPolicy Bypass -File ./build.ps1 -Mode integration

# ビルドのみ実行
powershell -NoProfile -ExecutionPolicy Bypass -File ./build.ps1 -Mode build-only
```

#### 成果物

ビルドによって生成された実行可能ファイルは `bin/game-list.exe` に配置されます。

#### テストの命名規則

- ユニットテスト: テスト関数名に `_Unit` サフィックスを含むもの
- 結合テスト: テスト関数名に `_Integration` サフィックスを含むもの（`tests/integration/` フォルダ内に配置）

詳細な要件定義については、`prompts/ideas/000-build.md` を参照してください。