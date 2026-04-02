CREATE TABLE employees (
    id          BIGSERIAL PRIMARY KEY,
    full_name   TEXT        NOT NULL,
    role        TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE rooms (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT        NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE chemicals (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT        NOT NULL,
    cas_number  TEXT,
    formula     TEXT,
    sds_url     TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE containers (
    id              BIGSERIAL PRIMARY KEY,
    chemical_id     BIGINT NOT NULL REFERENCES chemicals(id),
    room_id         BIGINT NOT NULL REFERENCES rooms(id),
    label_code      TEXT UNIQUE,
    quantity        DOUBLE PRECISION NOT NULL,
    unit            TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'available',
    checked_out_by  BIGINT REFERENCES employees(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);