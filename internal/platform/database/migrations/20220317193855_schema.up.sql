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

INSERT INTO games(id, name, year, platform, genre, publisher, created_at, modified_at) 
VALUES ('cde48d01-3fed-435c-8281-64cbd2de94a3'::UUID, 'Super Ghouls ''n Ghosts', 1991, 
		'Super Nintendo', 'Platform', 'Capcom', now(), now());