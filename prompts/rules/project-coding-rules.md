# プロジェクト固有のコーディング規約

> **注意**: 本ファイルは `Fundamental-Rules.md` の基本規範を前提としています。  
> 必ず `Fundamental-Rules.md` も併せて参照してください。

## 共通ルール（Always）

### コーディング規約

#### 言語・プロジェクト別の命名規則

##### Go 言語

以下の命名規約は `game-catalog-service` に適用されます。

**パッケージ名**
- 小文字、単数形を使用
- 例: `models`, `config`, `handlers`, `middleware`, `db`, `auth`, `utils`, `helpers`

**型名（構造体、インターフェース、型エイリアス）**
- パスカルケース（大文字で開始）
- 例: `Device`, `Config`, `User`, `Agent`, `Tag`, `Claims`, `ManageDeviceRequest`, `PostgresConnectionString`

**関数名**
- 公開関数: パスカルケース（大文字で開始）
  - 例: `GetDevices`, `AddDevice`, `Load`, `GenerateDeviceUUID`, `SetupDeviceRoutes`, `AuthMiddleware`
- 非公開関数: 小文字で開始
  - 例: `setEnvIfNotSet`, `setEnvForce`, `extractClaims`, `returnError`, `applyRequest`

**変数名**
- ローカル変数: 小文字で開始（キャメルケース）
  - 例: `device`, `devices`, `req`, `id`, `tagID`, `err`, `c`, `db`
- グローバル変数: 小文字で開始
  - 例: `conf`, `psqlDB`, `redisClient`, `jwtKey`

**定数名**
- 大文字のスネークケース（推奨）
  - 例: `MAX_RETRY_COUNT`, `DEFAULT_TIMEOUT`, `AccessTokenExpiry`

**構造体フィールド名**
- パスカルケース（大文字で開始）
- 例: `ID`, `Uuid`, `DeviceName`, `AgentID`, `CreatedAt`, `UpdatedAt`

**JSONタグ**
- スネークケース（小文字）
- 例: `json:"id"`, `json:"device_name"`, `json:"created_at"`, `json:"agent_id"`

**ファイル名**
- スネークケース（小文字）
- 例: `device.go`, `user.go`, `auth.go`, `testconfig.go`, `neo4j_test.go`, `postgres_test.go`

**テスト関数名**
- `Test` + 関数名 + `_Integration` または `_Unit` の形式
- 例: `TestNewPostgresDB_Integration`, `TestAddNode_Unit`, `TestConnectNeo4j_Integration`

**ファイル構成**
- パッケージごとにディレクトリを分離
- 例: `models/`, `handlers/`, `config/`, `middleware/`, `db/`, `auth/`, `utils/`
- テストファイルは以下のように配置
  - `tests/integration/`: 結合テスト
  - `tests/helpers/`: テストヘルパー
  - ユニットテストはテスト対象ファイルと同一のフォルダに配置

**インデント**
- タブ文字を使用（Go標準）
