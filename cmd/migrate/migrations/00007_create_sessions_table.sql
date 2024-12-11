-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.sessions (
    id text PRIMARY KEY,
    user_id bigint NOT NULL,
    expired_at timestamp(0) with time zone NOT NULL,
    
    CONSTRAINT user_fk FOREIGN KEY (user_id) REFERENCES public.users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.sessions;
-- +goose StatementEnd
