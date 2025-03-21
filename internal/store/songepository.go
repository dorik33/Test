package store

import (
	"context"
	"errors"

	"github.com/dorik33/Test/internal/models"
	"github.com/jackc/pgx/v5"
)

type SongRepository struct {
	store *Store
}

var (
	ErrSongNotFound = errors.New("song not found")
)

func (r *SongRepository) GetSongs(ctx context.Context, group, song string, limit, offset int) ([]models.Song, error) {
	query := "SELECT * FROM songs WHERE group_name LIKE $1 AND song_name LIKE $2 LIMIT $3 OFFSET $4;"
	rows, err := r.store.pool.Query(ctx, query, "%"+group+"%", "%"+song+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func (r *SongRepository) GetSongTextByID(ctx context.Context, songID int) (string, error) {
	query := "SELECT text FROM songs WHERE id = $1;"
	row := r.store.pool.QueryRow(ctx, query, songID)
	var text string
	err := row.Scan(&text)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrSongNotFound
		}
		return "", err
	}
	return text, nil
}

func (r *SongRepository) DeleteSong(ctx context.Context, songID int) error {
	query := "DELETE FROM songs WHERE id = $1;"
	res, err := r.store.pool.Exec(ctx, query, songID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrSongNotFound
	}
	return nil
}

func (r *SongRepository) UpdateSong(ctx context.Context, id int, song models.Song) error {
	query := "UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5 WHERE id = $6"
	res, err := r.store.pool.Exec(ctx, query, song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link, id)
	if res.RowsAffected() == 0 {
		return ErrSongNotFound
	}
	return err
}

func (r *SongRepository) AddSong(ctx context.Context, song models.Song) error {
	query := "INSERT INTO songs (group_name, song_name, release_date, text, link) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.store.pool.Exec(ctx, query, song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link)
	return err
}
