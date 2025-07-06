CREATE TABLE IF NOT EXISTS secrets (
    secret_id uuid DEFAULT gen_random_uuid (),
    owner_id  uuid REFERENCES users (user_id) on delete cascade,
    name           varchar(256) not null,
    kind           integer not null,
    metadata       bytea,
    data           bytea not null,
    primary key    (name, owner_id)
)
