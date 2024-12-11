-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.verifications (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id bigint NOT NULL,
    expired_at timestamp(0) with time zone NOT NULL DEFAULT (now() + interval '24 hours'),
    
    CONSTRAINT user_fk FOREIGN KEY (user_id) REFERENCES public.users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.verifications;
-- +goose StatementEnd
