package db

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore(dbURL string) (*Store, error) {
	database, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := database.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to PostgreSQL")
	return &Store{db: database}, nil
}

func (s *Store) InsertAnomaly(ctx context.Context, sourceIP, threatType, rawPayload string, severity int) error {
	query := `
		INSERT INTO ipsec_anomalies (source_ip, threat_type, severity, raw_payload)
		VALUES ($1, $2, $3, $4)
	`
	_, err := s.db.ExecContext(ctx, query, sourceIP, threatType, severity, rawPayload)
	return err
}

func (s *Store) Close() error {
	return s.db.Close()
}
