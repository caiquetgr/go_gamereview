CREATE TABLE games (
    id UUID,
    name VARCHAR(200),
    year SMALLINT,
    platform VARCHAR(100),
    genre VARCHAR(100),
    publisher VARCHAR(200),
    created_at TIMESTAMP WITH TIME ZONE,
    modified_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (id)
);

CREATE INDEX index_game_name ON games(name);
