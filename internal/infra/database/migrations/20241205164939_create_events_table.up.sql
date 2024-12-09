CREATE TABLE IF NOT EXISTS public.events
(
    eventid         serial PRIMARY KEY,
    userid          int NOT NULL references public.users (id),
    tytle           text NOT NULL,
    description     text NOT NULL,
    date            timestamp NOT NULL,
    image           text NOT NULL,
    location        text NOT NULL,
    lat             float NOT NULL,
    lon             float NOT NULL,
    created_date    timestamp NOT NULL,
    updated_date    timestamp NOT NULL,
    deleted_date    timestamp NULL
);