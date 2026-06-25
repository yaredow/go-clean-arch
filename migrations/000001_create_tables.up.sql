CREATE TABLE IF NOT EXISTS author (
    id bigserial PRIMARY KEY,
    name varchar(200) NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS article (
    id bigserial PRIMARY KEY,
    title varchar(45) NOT NULL,
    content text NOT NULL,
    author_id bigint NOT NULL DEFAULT 0,
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    created_at timestamptz NOT NULL DEFAULT NOW()
);

INSERT INTO author (name)
    VALUES ('Iman Tumorang');

