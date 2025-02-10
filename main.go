package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
	"xtcedata/database/metadata"
	"xtcedata/database/telemetrydata"
	"xtcedata/decoder"
	"xtcedata/extractor"
	"xtcedata/parser"
	"xtcedata/validator"
)

func main() {
	mongoDB, err := createMongo()
	if err != nil {
		logrus.Fatal(err)
	}
	mongoTelemetry := telemetrydata.NewMongoTelemtry(mongoDB)

	pg, err := createPostgres()
	if err != nil {
		logrus.Fatal(err)
	}

	pgMetaData := metadata.NewPostgresMetadata(pg)

	v := validator.NewXTCEValidatorForXML()
	p := parser.NewXTCEParserForXML()

	if err = v.Validate(os.Getenv("XTCE_FILE_PATH")); err != nil {
		logrus.Fatal(err)
	}

	xtceData, err := p.Parse(os.Getenv("XTCE_FILE_PATH"))
	if err != nil {
		logrus.Fatal(err)
	}

	if err = pgMetaData.InsertMetaData(xtceData.SatelliteID, xtceData.XTCEVersion); err != nil {
		logrus.Fatal(err)
	}

	if err = mongoTelemetry.InsertXTCEContainer(xtceData); err != nil {
		logrus.Fatal(err)
	}

	packet := []byte{
		// Satellite ID (8 bytes, assuming "SAT-001" padded)
		'S', 'A', 'T', '-', '0', '0', '1', ' ',
		// APID (4 bytes, example 1001)
		0x00, 0x00, 0x03, 0xE9,
		// Temperature (4 bytes, float32 = 25.3)
		0x41, 0xCB, 0x33, 0x33,
	}

	ext := extractor.NewXTCEExtractor()

	xtceDecoder := decoder.NewXTCEDecoder(ext, mongoTelemetry, pgMetaData)

	// Decode fields inside the selected container
	decodedData, err := xtceDecoder.DecodePacket(packet)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Infof("%+v", decodedData)

}

func createMongo() (*mongo.Database, error) {
	mongoURI := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}
	return client.Database("mcs"), nil
}

func createPostgres() (*pgxpool.Pool, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %v", err)
	}

	// Set connection timeout and pool settings
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = 30 * time.Minute

	// Create a connection pool
	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	// Check if the database is reachable
	err = db.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	logrus.Info("Successfully connected to PostgreSQL")
	return db, nil
}
