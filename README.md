# takase-game-list

ゲームソフト情報を管理するWeb APIアプリケーションです。正規化されたデータベース構造により、データの一貫性と保守性を向上させています。

## 技術スタック

- **言語**: Go 1.22以上
- **Webフレームワーク**: Gin
- **ORM**: GORM
- **データベース**: SQLite
- **実行環境**: Windows

## データベース構造

本プロジェクトでは、第2正規形（2NF）を満たす正規化されたデータベース構造を採用しています。

### モデル

- **Game**: ゲームソフト情報
  - PublisherID (uint, 外部キー, 必須)
  - SeriesID (*uint, 外部キー, 任意)
  - Platforms ([]Platform, 多対多リレーション, 必須)
  - Genres ([]Genre, 多対多リレーション, 任意)

- **Publisher**: 出版社情報
  - Name (string, unique, 必須)

- **Platform**: プラットフォーム情報
  - Name (string, unique, 必須)

- **Series**: シリーズ情報
  - Name (string, unique, 必須)

- **Genre**: ジャンル情報
  - Name (string, unique, 必須)

## APIエンドポイント

### ゲームAPI

#### POST /games
ゲームを登録します。

**リクエスト例:**
```json
{
  "title": "The Legend of Zelda: Breath of the Wild",
  "release_year": 2017,
  "publisher_id": 1,
  "platform_ids": [2],
  "series_id": 1,
  "genre_ids": [2, 3],
  "price": 6980
}
```

**レスポンス例:**
```json
{
  "id": 1,
  "title": "The Legend of Zelda: Breath of the Wild",
  "release_year": 2017,
  "publisher_id": 1,
  "series_id": 1,
  "price": 6980,
  "publisher": {
    "id": 1,
    "name": "Nintendo"
  },
  "platforms": [
    {
      "id": 2,
      "name": "Nintendo Switch"
    }
  ],
  "series": {
    "id": 1,
    "name": "The Legend of Zelda"
  },
  "genres": [
    {
      "id": 2,
      "name": "RPG"
    },
    {
      "id": 3,
      "name": "アドベンチャー"
    }
  ]
}
```

#### GET /games
ゲーム一覧を取得します。

**クエリパラメータ:**
- `publisher_id`: 出版社IDでフィルタリング
- `series_id`: シリーズIDでフィルタリング
- `platform_ids`: プラットフォームID配列でフィルタリング（複数指定可能）
- `genre_ids`: ジャンルID配列でフィルタリング（複数指定可能）
- `min_year`: 最小発売年でフィルタリング
- `max_year`: 最大発売年でフィルタリング
- `sort`: ソート項目（デフォルト: `created_at`）
- `order`: ソート順（`asc` または `desc`, デフォルト: `desc`）

**例:**
```
GET /games?publisher_id=1&platform_ids=2&platform_ids=3
```

#### PUT /games/:id
ゲーム情報を更新します。

**リクエスト例:**
```json
{
  "title": "Updated Title",
  "platform_ids": [1, 2]
}
```

#### DELETE /games/:id
ゲームを削除します。

#### GET /games/search
ゲームを検索します。

**クエリパラメータ:**
- `q`: 検索キーワード（タイトル部分一致）

**例:**
```
GET /games/search?q=Zelda
```

#### GET /games/statistics
統計情報を取得します。

**レスポンス例:**
```json
{
  "total_count": 10,
  "platform_counts": {
    "PC": 3,
    "Nintendo Switch": 5,
    "PlayStation 5": 2
  },
  "publisher_counts": {
    "Nintendo": 5,
    "Square Enix": 3,
    "Sony Interactive Entertainment": 2
  },
  "genre_counts": {
    "RPG": 4,
    "アクション": 3,
    "アドベンチャー": 3
  },
  "series_counts": {
    "The Legend of Zelda": 2,
    "Final Fantasy": 3
  },
  "year_counts": {
    "2017": 2,
    "2020": 3,
    "2024": 5
  },
  "total_price": 50000,
  "average_price": 5000.0
}
```

### マスタデータAPI

#### GET /publishers
全出版社の一覧を取得します。

#### POST /publishers
出版社を登録します。

**リクエスト例:**
```json
{
  "name": "New Publisher"
}
```

#### GET /platforms
全プラットフォームの一覧を取得します。

#### POST /platforms
プラットフォームを登録します。

**リクエスト例:**
```json
{
  "name": "New Platform"
}
```

#### GET /series
全シリーズの一覧を取得します。

#### POST /series
シリーズを登録します。

**リクエスト例:**
```json
{
  "name": "New Series"
}
```

#### GET /genres
全ジャンルの一覧を取得します。

#### POST /genres
ジャンルを登録します。

**リクエスト例:**
```json
{
  "name": "New Genre"
}
```

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

## 実行方法

ビルド後、以下のコマンドでアプリケーションを起動できます：

```powershell
.\bin\game-list.exe
```

アプリケーションは `http://localhost:8080` で起動します。

起動時には、自動的にマイグレーションが実行され、シードデータ（よく使われるPlatform、Publisher、Genre）が投入されます。

## プロジェクト構造

```
takase-game-list/
├── bin/                    # ビルド成果物
├── db/                     # データベース関連
│   ├── connection.go      # DB接続管理、マイグレーション
│   └── seed.go            # シードデータ投入
├── handlers/              # HTTPハンドラー
│   ├── game_handler.go    # ゲーム関連のエンドポイント処理
│   └── master_handler.go  # マスタデータ関連のエンドポイント処理
├── models/                # データモデル
│   ├── game.go
│   ├── publisher.go
│   ├── platform.go
│   ├── series.go
│   └── genre.go
├── middleware/            # ミドルウェア
│   ├── logger.go
│   └── error_handler.go
├── tests/                 # テスト
│   └── integration/      # 結合テスト
├── utils/                 # ユーティリティ
│   └── validators.go     # バリデーション関数
├── prompts/              # 計画・アイデア文書
│   ├── ideas/
│   ├── plans/
│   └── rules/
├── main.go               # エントリーポイント
├── build.ps1             # ビルドスクリプト
└── README.md
```

## 詳細な要件定義

詳細な要件定義については、以下のファイルを参照してください：

- `prompts/ideas/000-build.md`: ビルドシステムの要件
- `prompts/ideas/001-game-management.md`: ゲーム管理機能の要件
- `prompts/ideas/002-normalize-game-model.md`: ゲームモデル正規化の要件
- `prompts/plans/002-normalize-game-model.md`: ゲームモデル正規化の実装計画