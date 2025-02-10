package metadata

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type PostgresMetaData struct {
	db *pgxpool.Pool
}

func NewPostgresMetadata(db *pgxpool.Pool) *PostgresMetaData {
	return &PostgresMetaData{db: db}
}

func (p *PostgresMetaData) InsertMetaData(satelliteID string, version string) error {
	query := `
		INSERT INTO telemetry_metadata (satellite_id, version)
		VALUES ($1, $2)
		ON CONFLICT (satellite_id, version) DO NOTHING;
	`

	_, err := p.db.Exec(context.Background(), query, satelliteID, version)
	if err != nil {
		return fmt.Errorf("error inserting metadata: %v", err)
	}

	logrus.Infof("Inserted telemetry metadata: Satellite %s, Format %s, Version %s\n", satelliteID, version)
	return nil
}

func (p *PostgresMetaData) GetMetaData(satelliteID string) (string, error) {
	var version string
	query := `
		SELECT version FROM telemetry_metadata
		WHERE satellite_id = $1
		ORDER BY uploaded_at DESC
		LIMIT 1;
	`

	err := p.db.QueryRow(context.Background(), query, satelliteID).Scan(&version)
	if err == pgx.ErrNoRows {
		return "", fmt.Errorf("no telemetry metadata found for Satellite: %s", satelliteID)
	} else if err != nil {
		return "", fmt.Errorf("error retrieving telemetry metadata: %v", err)
	}

	logrus.Infof("Latest version for Satellite %s Version: %s\n", satelliteID, version)
	return version, nil
}
