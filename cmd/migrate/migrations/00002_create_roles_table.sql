-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.roles (
    id bigserial PRIMARY KEY,
    level int NOT NULL DEFAULT 0,
    name varchar(255) NOT NULL UNIQUE,
    description text
);

INSERT INTO public.roles (level, name, description)
VALUES 
    (1, 'user', 'user can - create posts, comments'),
    (2, 'moderator', 'moderator can - create/update/delete posts, comments'),
    (3, 'admin', 'admin');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.roles;
-- +goose StatementEnd
