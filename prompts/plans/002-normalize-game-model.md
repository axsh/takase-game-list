# ゲームモデル正規化実装計画（@prompts/ideas/002-normalize-game-model.md 準拠）

## ゴール

現在の `Game` モデルを正規化し、`Publisher`, `Platform`, `Series`, `Genre` を別テーブルに分割する。これにより、データの一貫性と保守性を向上させ、第2正規形（2NF）を満たすデータベース構造を実現する。

## 前提・制約

- 技術スタック: Go 1.22+, Gin, GORM, SQLite（Fundamental-Rules 準拠）
- 既存データの移行は考慮しない（新規データベースとして扱う）
- 計画文書ではコードを記述しない（Fundamental-Rules 1.計画策定要件）
- テスト駆動開発（TDD）で実装を進める（Fundamental-Rules 2.2）
- テスト命名規則と配置: `_Unit`/`_Integration`, `tests/integration/` を遵守（project-coding-rules）
- 結合テストを先に記述し、その後ユニットテストを記述する（Fundamental-Rules 2.2）

## 成果物

- データモデル: Publisher, Platform, Series, Genre の新規モデル定義
- データモデル: Game モデルの変更（外部キーとリレーションの追加）
- バリデーション: 新しいモデルに対応したバリデーション関数
- HTTPハンドラー: 正規化後のモデルに対応したハンドラーの更新
- データベース接続管理: 新しいモデルのマイグレーション対応
- テスト: 結合テストとユニットテスト（既存テストの更新含む）

## 全体構造

### ディレクトリ構造（変更点）

```
takase-game-list/
├── models/                    # データモデル
│   ├── game.go               # Gameモデル定義（変更）
│   ├── game_test.go          # Gameモデルのユニットテスト（更新）
│   ├── publisher.go           # Publisherモデル定義（新規）
│   ├── publisher_test.go      # Publisherモデルのユニットテスト（新規）
│   ├── platform.go            # Platformモデル定義（新規）
│   ├── platform_test.go       # Platformモデルのユニットテスト（新規）
│   ├── series.go              # Seriesモデル定義（新規）
│   ├── series_test.go         # Seriesモデルのユニットテスト（新規）
│   ├── genre.go               # Genreモデル定義（新規）
│   └── genre_test.go          # Genreモデルのユニットテスト（新規）
├── handlers/                  # HTTPハンドラー
│   ├── game_handler.go        # ゲーム関連のエンドポイント処理（更新）
│   ├── game_handler_test.go   # ゲームハンドラーのユニットテスト（更新）
│   ├── master_handler.go      # マスタデータ関連のエンドポイント処理（新規）
│   └── master_handler_test.go # マスタデータハンドラーのユニットテスト（新規）
├── utils/                     # ユーティリティ
│   ├── validators.go          # バリデーション関数（更新）
│   └── validators_test.go     # バリデーション関数のユニットテスト（更新）
├── tests/
│   └── integration/          # 結合テスト
│       ├── game_api_test.go   # ゲームAPIの結合テスト（更新）
│       ├── normalized_model_test.go  # 正規化モデルの結合テスト（新規）
│       └── master_api_test.go  # マスタデータAPIの結合テスト（新規）
├── db/                        # データベース関連
│   ├── connection.go          # DB接続管理、マイグレーション（更新）
│   └── seed.go                # シードデータ投入（新規）
```

### パッケージ構成と役割

#### models パッケージ

##### Publisher モデル（新規）
- **役割**: 出版社情報を管理するテーブル構造を定義
- **主要な型**: `Publisher` 構造体
- **機能**:
  - GORMモデルとして定義（ID, CreatedAt, UpdatedAt, DeletedAt）
  - Name フィールド（varchar(255), not null, unique）
  - Game との1対多リレーション定義（Games []Game）

##### Platform モデル（新規）
- **役割**: プラットフォーム情報を管理するテーブル構造を定義
- **主要な型**: `Platform` 構造体
- **機能**:
  - GORMモデルとして定義（ID, CreatedAt, UpdatedAt, DeletedAt）
  - Name フィールド（varchar(100), not null, unique）
  - Game との多対多リレーション定義（Games []Game, 中間テーブル: game_platforms）

##### Series モデル（新規）
- **役割**: シリーズ情報を管理するテーブル構造を定義
- **主要な型**: `Series` 構造体
- **機能**:
  - GORMモデルとして定義（ID, CreatedAt, UpdatedAt, DeletedAt）
  - Name フィールド（varchar(255), not null, unique）
  - Game との1対多リレーション定義（Games []Game）

##### Genre モデル（新規）
- **役割**: ジャンル情報を管理するテーブル構造を定義
- **主要な型**: `Genre` 構造体
- **機能**:
  - GORMモデルとして定義（ID, CreatedAt, UpdatedAt, DeletedAt）
  - Name フィールド（varchar(100), not null, unique）
  - Game との多対多リレーション定義（Games []Game, 中間テーブル: game_genres）

##### Game モデル（変更）
- **役割**: ゲームソフト情報を管理するテーブル構造を定義（正規化後）
- **主要な型**: `Game` 構造体
- **変更内容**:
  - `Publisher` (string) → `PublisherID` (uint, 外部キー, not null)
  - `Platform` (string) → `Platforms` ([]Platform, 多対多リレーション)
  - `Series` (*string) → `SeriesID` (*uint, 外部キー, NULL許容)
  - `Genre` (*string) → `Genres` ([]Genre, 多対多リレーション)
  - リレーションフィールドの追加（Publisher, Platforms, Series, Genres）
- **保持する機能**:
  - 必須項目: Title, ReleaseYear, PublisherID
  - 任意項目: SeriesID, Price
  - システム項目: ID, CreatedAt, UpdatedAt, DeletedAt

#### utils パッケージ

##### validators.go（更新）
- **役割**: 正規化後のモデルに対応したバリデーション関数を提供
- **変更内容**:
  - `ValidateGame` 関数の更新:
    - PublisherID の存在確認（必須）
    - SeriesID の存在確認（任意、nilでない場合のみ）
    - Platforms 配列の各IDの存在確認（必須、少なくとも1つ必要）
    - Genres 配列の各IDの存在確認（任意、空配列も許容）
  - 新規関数の追加:
    - `ValidatePublisherID`: PublisherIDの存在確認
    - `ValidateSeriesID`: SeriesIDの存在確認（nil許容）
    - `ValidatePlatformIDs`: PlatformID配列の存在確認（少なくとも1つ必要）
    - `ValidateGenreIDs`: GenreID配列の存在確認（空配列許容）

#### handlers パッケージ

##### game_handler.go（更新）
- **役割**: 正規化後のモデルに対応したHTTPリクエスト処理を実装
- **変更内容**:
  - `CreateGame`: 
    - リクエスト形式の変更（PublisherID, PlatformIDs配列, SeriesID, GenreIDs配列）
    - バリデーションの更新
    - 多対多リレーションの保存処理（Platforms, Genres）
    - レスポンスにPreloadした関連データを含める
  - `GetGames`:
    - フィルタリング条件の変更（PublisherID, PlatformIDs, SeriesID, GenreIDs）
    - Preloadを使用した関連データの取得
  - `UpdateGame`:
    - リクエスト形式の変更
    - 多対多リレーションの更新処理
    - Preloadを使用した関連データの取得
  - `SearchGames`:
    - Preloadを使用した関連データの取得
  - `GetStatistics`:
    - 集計ロジックの更新（正規化後のテーブル構造に対応）
  - `DeleteGame`:
    - 多対多リレーションの削除処理（中間テーブルのレコード削除）

#### db パッケージ

##### connection.go（更新）
- **役割**: データベース接続の初期化とマイグレーション管理
- **変更内容**:
  - `Migrate` 関数の更新:
    - 新規モデル（Publisher, Platform, Series, Genre）のマイグレーション追加
    - Game モデルのマイグレーション更新
    - 中間テーブル（game_platforms, game_genres）の自動生成確認

## 実装の論理的な順序

### フェーズ1: 結合テストの定義

結合テストを先に記述し、要求仕様を明確にする。

#### 1.1 正規化モデルの結合テスト作成

**ファイル**: `tests/integration/normalized_model_test.go`

**テスト内容**:
- データベースマイグレーションの確認:
  - Publisher, Platform, Series, Genre テーブルの存在確認
  - Game テーブルの構造変更確認（外部キー、中間テーブル）
  - 中間テーブル（game_platforms, game_genres）の存在確認
- モデル間のリレーション確認:
  - Publisher と Game の1対多関係
  - Platform と Game の多対多関係
  - Series と Game の1対多関係（NULL許容）
  - Genre と Game の多対多関係
- データ操作の確認:
  - Publisher, Platform, Series, Genre の作成
  - Game の作成（外部キーと多対多リレーションを含む）
  - Game の取得（Preloadによる関連データの取得）
  - Game の更新（多対多リレーションの更新）
  - Game の削除（中間テーブルのレコード削除確認）

#### 1.2 既存結合テストの更新

**ファイル**: `tests/integration/game_api_test.go`

**更新内容**:
- APIリクエスト/レスポンス形式の変更に対応:
  - POST /games: PublisherID, PlatformIDs配列, SeriesID, GenreIDs配列を受け取る
  - GET /games: レスポンスにPublisher, Platforms, Series, Genresを含む
  - PUT /games/:id: 更新リクエスト形式の変更
  - フィルタリング条件の変更（PublisherID, PlatformIDs等）
- バリデーションエラーの確認:
  - PublisherIDが存在しない場合
  - PlatformIDsが空配列の場合
  - PlatformIDsに存在しないIDが含まれる場合
  - SeriesIDが存在しない場合（nilでない場合）

### フェーズ2: ユニットテストの定義と実装

結合テストを仕様として、ユニットテストを記述する。

#### 2.1 新規モデルのユニットテスト

**ファイル**: 
- `models/publisher_test.go`
- `models/platform_test.go`
- `models/series_test.go`
- `models/genre_test.go`

**テスト内容**:
- 各モデルの構造体定義の確認
- GORMタグの設定確認
- JSONタグの設定確認
- リレーション定義の確認

#### 2.2 Gameモデルのユニットテスト更新

**ファイル**: `models/game_test.go`

**更新内容**:
- 新しいフィールド（PublisherID, SeriesID, Platforms, Genres）の確認
- リレーション定義の確認
- 既存テストの更新（文字列フィールドから外部キー/リレーションへ）

#### 2.3 バリデーション関数のユニットテスト更新

**ファイル**: `utils/validators_test.go`

**更新内容**:
- `ValidateGame` 関数のテスト更新
- 新規バリデーション関数のテスト追加:
  - `ValidatePublisherID` のテスト
  - `ValidateSeriesID` のテスト
  - `ValidatePlatformIDs` のテスト
  - `ValidateGenreIDs` のテスト

#### 2.4 ハンドラーのユニットテスト更新

**ファイル**: `handlers/game_handler_test.go`

**更新内容**:
- 各ハンドラー関数のテスト更新（モックを使用）
- リクエスト/レスポンス形式の変更に対応
- バリデーションエラーのテスト追加

### フェーズ3: モデルの実装

ユニットテストを通過させるために最小限の目的コードを記述する。

#### 3.1 新規モデルの実装

**実装順序**:
1. `models/publisher.go`: Publisher モデルの実装
2. `models/platform.go`: Platform モデルの実装
3. `models/series.go`: Series モデルの実装
4. `models/genre.go`: Genre モデルの実装

**実装内容**:
- 各モデルの構造体定義
- GORMタグの設定（primaryKey, not null, unique等）
- JSONタグの設定（スネークケース）
- リレーション定義（GORMのmany2manyタグ等）

#### 3.2 Gameモデルの更新

**ファイル**: `models/game.go`

**実装内容**:
- 既存フィールドの削除（Publisher, Platform, Series, Genreの文字列フィールド）
- 新規フィールドの追加:
  - PublisherID (uint, not null, index)
  - SeriesID (*uint, index, NULL許容)
  - Platforms ([]Platform, many2many)
  - Genres ([]Genre, many2many)
- リレーションフィールドの追加:
  - Publisher (Publisher, foreignKey)
  - Series (*Series, foreignKey)
- GORMタグとJSONタグの更新

### フェーズ4: バリデーション関数の実装

#### 4.1 バリデーション関数の更新

**ファイル**: `utils/validators.go`

**実装内容**:
- `ValidateGame` 関数の更新:
  - PublisherID の必須チェック
  - SeriesID の存在確認（nilでない場合のみ）
  - Platforms 配列の必須チェック（少なくとも1つ必要）
  - Genres 配列のチェック（空配列許容）
- 新規バリデーション関数の実装:
  - `ValidatePublisherID`: データベースにPublisherIDが存在するか確認
  - `ValidateSeriesID`: データベースにSeriesIDが存在するか確認（nil許容）
  - `ValidatePlatformIDs`: データベースにPlatformID配列の各IDが存在するか確認（少なくとも1つ必要）
  - `ValidateGenreIDs`: データベースにGenreID配列の各IDが存在するか確認（空配列許容）

**注意**: データベース接続が必要なバリデーション関数は、関数シグネチャに `*gorm.DB` パラメータを追加する。

### フェーズ5: ハンドラーの実装

#### 5.1 ハンドラー関数の更新

**ファイル**: `handlers/game_handler.go`

**実装内容**:
- `CreateGame`:
  - リクエスト構造体の定義（PublisherID, PlatformIDs, SeriesID, GenreIDs）
  - バリデーションの呼び出し
  - Platform と Genre の取得（ID配列から）
  - Game の作成と多対多リレーションの設定
  - Preloadを使用したレスポンスデータの取得
- `GetGames`:
  - フィルタリング条件の変更（PublisherID, PlatformIDs等）
  - Preloadを使用した関連データの取得
- `UpdateGame`:
  - リクエスト構造体の定義
  - 多対多リレーションの更新処理
  - Preloadを使用したレスポンスデータの取得
- `SearchGames`:
  - Preloadを使用した関連データの取得
- `GetStatistics`:
  - 集計ロジックの更新（正規化後のテーブル構造に対応）
- `DeleteGame`:
  - 多対多リレーションの削除処理（GORMが自動的に中間テーブルを削除）

### フェーズ6: データベース接続管理の更新

#### 6.1 マイグレーションの更新

**ファイル**: `db/connection.go`

**実装内容**:
- `Migrate` 関数の更新:
  - 新規モデル（Publisher, Platform, Series, Genre）のAutoMigrate追加
  - Game モデルのAutoMigrate更新
  - 中間テーブルはGORMが自動生成するため、明示的な処理は不要

#### 6.2 シードデータの実装

**ファイル**: `db/seed.go`

**実装内容**:
- `SeedMasterData` 関数の実装:
  - よく使われるPlatformデータの投入（例: "PC", "PlayStation 5", "Nintendo Switch", "Xbox Series X"等）
  - よく使われるPublisherデータの投入（例: "Nintendo", "Sony Interactive Entertainment"等）
  - よく使われるGenreデータの投入（例: "アクション", "RPG", "アドベンチャー"等）
  - 既存データの重複チェック（FirstOrCreateを使用）
- `main.go` での呼び出し:
  - マイグレーション後にシードデータを投入（オプション）

### フェーズ7: マスタデータ管理APIの実装

#### 7.1 マスタデータハンドラーの実装

**ファイル**: `handlers/master_handler.go`

**実装内容**:
- `GetPublishers`: 全Publisherの一覧取得
- `CreatePublisher`: Publisherの作成（Nameの重複チェック）
- `GetPlatforms`: 全Platformの一覧取得
- `CreatePlatform`: Platformの作成（Nameの重複チェック）
- `GetSeries`: 全Seriesの一覧取得
- `CreateSeries`: Seriesの作成（Nameの重複チェック）
- `GetGenres`: 全Genreの一覧取得
- `CreateGenre`: Genreの作成（Nameの重複チェック）

**バリデーション**:
- Name の必須チェック
- Name の文字列長チェック（Publisher, Series: 255文字、Platform, Genre: 100文字）
- Name の重複チェック（データベースのunique制約）

#### 7.2 ルーティングの追加

**ファイル**: `main.go`

**実装内容**:
- マスタデータ用のエンドポイント追加:
  - GET /publishers
  - POST /publishers
  - GET /platforms
  - POST /platforms
  - GET /series
  - POST /series
  - GET /genres
  - POST /genres

### フェーズ8: 初期検証

#### 8.1 ユニットテストの実行

- 各モデルのユニットテストを実行
- バリデーション関数のユニットテストを実行
- ハンドラーのユニットテストを実行
- ビルドスクリプトを単体テストのみ実行するオプションで実行
- エラーが存在する場合は、コードまたはテスト設計を見直す

### フェーズ9: 最終検証

#### 9.1 結合テストの実行

- 正規化モデルの結合テストを実行
- 既存の結合テスト（game_api_test.go）を実行
- ビルドスクリプトを結合テストのみ実行するオプションで実行
- 結合テストまでエラーがないことを確認

## テストの概要と計画

### 結合テスト（tests/integration/）

#### normalized_model_test.go（新規）

**テスト関数**:
- `TestNormalizedModelMigration_Integration`: マイグレーションの確認
- `TestPublisherGameRelation_Integration`: Publisher と Game の1対多関係の確認
- `TestPlatformGameRelation_Integration`: Platform と Game の多対多関係の確認
- `TestSeriesGameRelation_Integration`: Series と Game の1対多関係（NULL許容）の確認
- `TestGenreGameRelation_Integration`: Genre と Game の多対多関係の確認
- `TestCreateGameWithRelations_Integration`: リレーションを含むGame作成の確認
- `TestGetGameWithPreload_Integration`: Preloadによる関連データ取得の確認
- `TestUpdateGameRelations_Integration`: 多対多リレーションの更新確認
- `TestDeleteGameWithRelations_Integration`: 中間テーブルのレコード削除確認

#### game_api_test.go（更新）

**更新するテスト関数**:
- `TestCreateGame_Integration`: リクエスト/レスポンス形式の変更に対応
- `TestGetGames_Integration`: フィルタリング条件とレスポンス形式の変更に対応
- `TestUpdateGame_Integration`: リクエスト/レスポンス形式の変更に対応
- `TestSearchGames_Integration`: レスポンス形式の変更に対応
- `TestGetStatistics_Integration`: 集計ロジックの更新に対応

**新規テスト関数**:
- `TestCreateGameValidation_Integration`: バリデーションエラーの確認
  - PublisherIDが存在しない場合
  - PlatformIDsが空配列の場合
  - PlatformIDsに存在しないIDが含まれる場合
  - SeriesIDが存在しない場合（nilでない場合）

#### master_api_test.go（新規）

**テスト関数**:
- `TestGetPublishers_Integration`: Publisher一覧取得の確認
- `TestCreatePublisher_Integration`: Publisher作成の確認
- `TestCreatePublisherDuplicate_Integration`: Publisher重複作成のエラー確認
- `TestGetPlatforms_Integration`: Platform一覧取得の確認
- `TestCreatePlatform_Integration`: Platform作成の確認
- `TestCreatePlatformDuplicate_Integration`: Platform重複作成のエラー確認
- `TestGetSeries_Integration`: Series一覧取得の確認
- `TestCreateSeries_Integration`: Series作成の確認
- `TestGetGenres_Integration`: Genre一覧取得の確認
- `TestCreateGenre_Integration`: Genre作成の確認

### ユニットテスト

#### models パッケージ

**publisher_test.go（新規）**:
- `TestPublisher_Unit`: Publisher構造体の基本動作確認

**platform_test.go（新規）**:
- `TestPlatform_Unit`: Platform構造体の基本動作確認

**series_test.go（新規）**:
- `TestSeries_Unit`: Series構造体の基本動作確認

**genre_test.go（新規）**:
- `TestGenre_Unit`: Genre構造体の基本動作確認

**game_test.go（更新）**:
- `TestGame_Unit`: Game構造体の基本動作確認（更新）
- `TestGameWithRelations_Unit`: リレーションフィールドの確認（新規）
- `TestGameWithOptionalFields_Unit`: 任意項目の確認（更新）
- `TestGameWithNilOptionalFields_Unit`: nil許容フィールドの確認（更新）

#### utils パッケージ

**validators_test.go（更新）**:
- `TestValidateGame_Unit`: ValidateGame関数のテスト（更新）
- `TestValidatePublisherID_Unit`: ValidatePublisherID関数のテスト（新規）
- `TestValidateSeriesID_Unit`: ValidateSeriesID関数のテスト（新規）
- `TestValidatePlatformIDs_Unit`: ValidatePlatformIDs関数のテスト（新規）
- `TestValidateGenreIDs_Unit`: ValidateGenreIDs関数のテスト（新規）

#### handlers パッケージ

**game_handler_test.go（更新）**:
- 各ハンドラー関数のテストを更新（モックを使用）
- リクエスト/レスポンス形式の変更に対応
- バリデーションエラーのテスト追加

**master_handler_test.go（新規）**:
- `TestGetPublishers_Unit`: GetPublishers関数のテスト
- `TestCreatePublisher_Unit`: CreatePublisher関数のテスト
- `TestGetPlatforms_Unit`: GetPlatforms関数のテスト
- `TestCreatePlatform_Unit`: CreatePlatform関数のテスト
- `TestGetSeries_Unit`: GetSeries関数のテスト
- `TestCreateSeries_Unit`: CreateSeries関数のテスト
- `TestGetGenres_Unit`: GetGenres関数のテスト
- `TestCreateGenre_Unit`: CreateGenre関数のテスト

## 実装時の注意点

### データベース設計

1. **一意性制約**: Publisher, Platform, Series, Genre の Name に unique 制約を設定
2. **インデックス**: 外部キー（PublisherID, SeriesID）にインデックスを設定
3. **中間テーブル**: GORMが自動生成するため、明示的な定義は不要

### API設計

1. **リクエスト形式**: 
   - PublisherID: uint（必須）
   - PlatformIDs: []uint（必須、少なくとも1つ）
   - SeriesID: *uint（任意）
   - GenreIDs: []uint（任意、空配列も許容）

2. **レスポンス形式**:
   - Publisher, Platforms, Series, Genres を Preload して含める
   - ネストされた構造で返す

3. **バリデーション**:
   - 外部キーの存在確認はデータベースアクセスが必要
   - ハンドラー内でバリデーション関数を呼び出す際に、データベース接続を渡す

### パフォーマンス

1. **Preloadの使用**: 関連データを取得する際は、N+1問題を避けるため Preload を使用
2. **インデックスの活用**: 外部キーにインデックスを設定してクエリ性能を向上

### マスタデータ管理

1. **マスタデータの追加方法**:
   - **API経由**: POST /publishers, POST /platforms, POST /series, POST /genres を使用
   - **シードデータ**: よく使われるデータを初期投入（`db/seed.go`）
   - **手動追加**: 必要に応じてAPI経由で追加

2. **シードデータの実行タイミング**:
   - マイグレーション後に自動実行（オプション）
   - または、手動で実行（コマンドライン引数等）

3. **重複チェック**:
   - Name の unique 制約により、重複作成はエラーとなる
   - シードデータ投入時は FirstOrCreate を使用して重複を回避

4. **マスタデータの取得**:
   - ゲーム作成時に利用可能なマスタデータを取得するため、GET エンドポイントを提供
   - フロントエンドでドロップダウン等の選択肢として使用可能

### テスト

1. **モックの使用**: ユニットテストでは、データベース接続をモック化
2. **結合テスト**: 実際のデータベースを使用して、リレーションの動作を確認
3. **シードデータのテスト**: シードデータ投入のテストを追加

## 結論

この計画に従って実装を進めることで、以下の点を実現できます：

- ✅ 第2正規形（2NF）を満たすデータベース構造
- ✅ データの一貫性と整合性の向上
- ✅ メンテナンス性の向上
- ✅ テスト駆動開発（TDD）による品質保証
- ✅ 既存機能との互換性を考慮した段階的な実装

この正規化により、より保守性の高いデータベース構造を実現し、将来の拡張にも対応可能な設計となります。
