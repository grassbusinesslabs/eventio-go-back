CREATE TABLE IF NOT EXISTS public.events
(
    id              serial PRIMARY KEY,
    user_id         int NOT NULL references public.users(id),
    title           text NOT NULL,
    description     text NOT NULL,
    date            timestamp NOT NULL,
    image           text NOT NULL,
    location        text NOT NULL,
    lat             float NOT NULL,
    lon             float NOT NULL,
    created_date    timestamptz NOT NULL,
    updated_date    timestamptz NOT NULL,
    deleted_date    timestamptz NULL
);