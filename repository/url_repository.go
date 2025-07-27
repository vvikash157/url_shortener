package repository

import (
	"database/sql"
	"errors"
	"fmt"
)

type UrlRepository interface {
	GetCodeByLongUrl(longUrl string) (string, error)
	GetLongUrlByCode(shortCode string) (string, error)
	InsertUrl(longUrl, shortCode string) error
}

type PostgresUrlRepository struct {
	db *sql.DB
}

func NewPostgresUrlRepository(db *sql.DB) *PostgresUrlRepository {
	return &PostgresUrlRepository{db: db}
}

func (s *PostgresUrlRepository) InsertUrl(longUrl, shortCode string) error {
	_, err := s.db.Exec("INSERT INTO urls(long_url,short_url) values($1,$2)", longUrl, shortCode)
	fmt.Printf("url is inserted into DB %s", shortCode)
	return err
}

func (s *PostgresUrlRepository) GetCodeByLongUrl(longUrl string) (string, error) {
	var code string
	err := s.db.QueryRow("SELECT short_url from urls where long_url=$1", longUrl).Scan(&code)
	if err == sql.ErrNoRows {
		return "", errors.New("didn't find code for long url")
	}
	return code, err
}

func (s *PostgresUrlRepository) GetLongUrlByCode(shortCode string) (string, error) {
	var long_url string
	err := s.db.QueryRow("SELECT long_url from urls where short_url=$1", shortCode).Scan(&long_url)

	if err == sql.ErrNoRows {
		return "", errors.New("no records found of short_url")
	}
	return long_url, err
}
