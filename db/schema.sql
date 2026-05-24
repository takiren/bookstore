-- ============================================================
-- 蔵書管理データベーススキーマ
-- ============================================================

-- ------------------------------------------------------------
-- updated_at 自動更新トリガー関数
-- ------------------------------------------------------------
CREATE OR REPLACE FUNCTION public.set_updated_at()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

-- ------------------------------------------------------------
-- books テーブル
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS public.books (
    id          BIGSERIAL    PRIMARY KEY,
    user_id     UUID         NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    isbn        TEXT         NOT NULL,
    is_read     BOOLEAN      NOT NULL DEFAULT FALSE,
    book_info   JSONB        NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT books_isbn_format CHECK (isbn ~ '^(\d{9}[\dX]|\d{13})$')
);

-- RLS パフォーマンス最適化: user_id にインデックスを張る
CREATE INDEX IF NOT EXISTS books_user_id_idx ON public.books (user_id);

-- updated_at トリガー
CREATE TRIGGER books_set_updated_at
    BEFORE UPDATE ON public.books
    FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

-- ------------------------------------------------------------
-- Row Level Security (RLS)
-- 認証済みユーザーは自分のレコードのみ操作可能
-- ------------------------------------------------------------
ALTER TABLE public.books ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.books FORCE ROW LEVEL SECURITY;

-- SELECT: 自分の蔵書のみ参照可
CREATE POLICY books_select ON public.books
    FOR SELECT TO authenticated
    USING ((SELECT auth.uid()) = user_id);

-- INSERT: user_id が自分自身のレコードのみ挿入可
CREATE POLICY books_insert ON public.books
    FOR INSERT TO authenticated
    WITH CHECK ((SELECT auth.uid()) = user_id);

-- UPDATE: 自分のレコードのみ更新可
CREATE POLICY books_update ON public.books
    FOR UPDATE TO authenticated
    USING  ((SELECT auth.uid()) = user_id)
    WITH CHECK ((SELECT auth.uid()) = user_id);

-- DELETE: 自分のレコードのみ削除可
CREATE POLICY books_delete ON public.books
    FOR DELETE TO authenticated
    USING ((SELECT auth.uid()) = user_id);
