CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    song_name VARCHAR(255) NOT NULL,
    release_date TEXT,
    text TEXT,
    link VARCHAR(255)
);