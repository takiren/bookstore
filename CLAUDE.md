# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Structure

```
bookstore/
├── web/    # Next.js 16 フロントエンド (pnpm)
└── api/    # Go + Echo バックエンド
```

## Web (Next.js 16)

> **重要**: Next.js 16 はトレーニングデータと異なる破壊的変更を含む。コードを書く前に必ず `web/node_modules/next/dist/docs/` を参照すること。

```bash
cd web
pnpm dev          # 開発サーバー起動
pnpm build        # ビルド
pnpm lint         # Biome でリント
pnpm format       # Biome でフォーマット
```

- パッケージマネージャーは **pnpm**（npm/yarn 不可）
- リンター・フォーマッターは **Biome**（ESLint/Prettier 不可）。設定は `web/biome.json`
- **React Compiler** が有効（`next.config.ts`）
- **Tailwind CSS v4**（v3 とは設定・クラス名が異なる）
- App Router を使用

## API (Go + Echo)

```bash
cd api
go run ./cmd/main.go    # 開発サーバー起動
go build ./...          # ビルド
go test ./...           # テスト
go vet ./...            # 静的解析
```

予定している技術スタック:
- **Echo v4** — HTTP フレームワーク
- **oapi-codegen** — OpenAPI 仕様 → Go コード自動生成
- **pgx** — PostgreSQL ドライバー
- **sqlc** — SQL → Go コード自動生成

oapi-codegen でコード生成する場合は `api/openapi/` に仕様書を置き、生成コードは `api/gen/` に出力する想定。sqlc の生成コードは `api/db/` に置く想定。

## インフラ・外部サービス

| サービス | 用途 |
|---------|------|
| Vercel | Next.js ホスティング |
| Render | Go API ホスティング（無料枠はスリープあり） |
| Supabase | PostgreSQL + Auth |
| 国立国会図書館 API | ISBN → 書誌情報・サムネイル取得 |
