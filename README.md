# 蔵書管理アプリ

![Next.js](https://img.shields.io/badge/Next.js-16-000000?style=flat-square&logo=next.js&logoColor=white)
![Go](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat-square&logo=go&logoColor=white)
![Echo](https://img.shields.io/badge/Echo-v4-00ADD8?style=flat-square&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supabase-3ECF8E?style=flat-square&logo=supabase&logoColor=white)
![Vercel](https://img.shields.io/badge/Deploy-Vercel-000000?style=flat-square&logo=vercel&logoColor=white)
![Render](https://img.shields.io/badge/API-Render-46E3B7?style=flat-square&logo=render&logoColor=white)

スマホのカメラでISBNを読み取って蔵書を管理するWebアプリ。

---

## 機能

| 機能           | 説明                                                                          |
| -------------- | ----------------------------------------------------------------------------- |
| 蔵書登録       | カメラでバーコードをスキャンしてISBNを取得、国会図書館APIで書誌情報を自動補完 |
| 蔵書一覧       | 登録した本を一覧表示。国会図書館APIのサムネイルつき                           |
| ステータス管理 | 未読 / 読書中 / 読了 / 積読 などのステータスを設定                            |
| 蔵書削除       | 売却・紛失した本を管理から除外                                                |

---

## アーキテクチャ

```mermaid
graph TD
    Client["ブラウザ / スマホ"]

    subgraph Vercel
        Next["Next.js 16"]
    end

    subgraph Render
        Echo["Go + Echo API"]
    end

    subgraph Supabase
        Auth["Supabase Auth"]
        DB["PostgreSQL\n500MB free tier"]
    end

    NDL["国立国会図書館 API\n書誌情報 / サムネイル"]

    Client -->|HTTPS| Next
    Next -->|REST API| Echo
    Next -->|認証| Auth
    Echo -->|pgx / sqlc| DB
    Echo -->|書誌情報取得| NDL
```

---

## 技術スタック

### フロントエンド
- **Next.js 16** — App Router
- **Vercel** — ホスティング（無料・非商用OK）
- **Supabase Auth** — 認証
- カメラAPI（`getUserMedia`） — バーコードスキャン

### バックエンド
- **Go** — API サーバー
- **Echo v4** — HTTPフレームワーク
- **oapi-codegen** — OpenAPI仕様からサーバー・クライアントコードを自動生成
- **pgx** — PostgreSQLドライバー
- **sqlc** — SQLからGoコードを自動生成
- **Render** — ホスティング（無料・スリープあり）

### データベース / 認証
- **Supabase PostgreSQL** — データベース（無料枠 500MB）
- **Supabase Auth** — ユーザー認証

### 外部API
- **国立国会図書館 書誌情報API** — ISBN から書籍情報・サムネイルを取得

---

## ディレクトリ構成

```
bookstore/
├── web/          # Next.js フロントエンド
└── api/          # Go + Echo バックエンド
    ├── openapi/  # OpenAPI 仕様書 (oapi-codegen のソース)
    ├── handler/  # Echo ハンドラー
    ├── db/       # sqlc 生成コード
    └── main.go
```

---

## ローカル開発

### 前提
- Go 1.24+
- Node.js 22+
- Docker（Supabase ローカル環境用）

### セットアップ

```bash
# リポジトリのクローン
git clone https://github.com/yourname/bookstore.git
cd bookstore

# バックエンド
cd api
go mod tidy
go run main.go

# フロントエンド
cd web
npm install
npm run dev
```

### 環境変数

`.env.local`（フロントエンド）と `.env`（バックエンド）を各ディレクトリに作成してください。

**web/.env.local**
```
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_SUPABASE_URL=your-supabase-url
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-supabase-anon-key
```

**api/.env**
```
DATABASE_URL=postgresql://postgres:password@localhost:5432/bookstore
PORT=8080
```

---

## ライセンス

非商用利用のみ。
