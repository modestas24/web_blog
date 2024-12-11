-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.posts (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL,
    title text NOT NULL,
    content text NOT NULL,
    
    verified boolean NOT NULL DEFAULT FALSE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

    CONSTRAINT user_fk FOREIGN KEY (user_id) REFERENCES public.users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.posts;
-- +goose StatementEnd
