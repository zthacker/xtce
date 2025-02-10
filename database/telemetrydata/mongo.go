package telemetrydata

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"xtcedata/models"
)

type MongoTelemetry struct {
	db *mongo.Database
}

func NewMongoTelemtry(client *mongo.Database) *MongoTelemetry {
	return &MongoTelemetry{db: client}
}

func (mt *MongoTelemetry) InsertXTCEContainer(xtceData *models.XTCETelemetryDefinition) error {
	collection := mt.db.Collection("xtce_containers")

	log.Printf("Inserting XTCE data into MongoDB: %+v\n", xtceData)

	filter := bson.M{"satellite_id": xtceData.SatelliteID, "xtce_version": xtceData.XTCEVersion}
	update := xtceData

	opts := options.Replace().SetUpsert(true)
	result, err := collection.ReplaceOne(context.Background(), filter, update, opts)
	if err != nil {
		log.Printf("Error inserting XTCE data: %v\n", err)
		return err
	}

	log.Printf("XTCE data inserted. Matched: %d, Modified: %d, Upserted: %v\n", result.MatchedCount, result.ModifiedCount, result.UpsertedID)
	return nil
}

func (mt *MongoTelemetry) GetContainerByAPID(satelliteID, xtceVersion string, apid uint32) (*models.Container, error) {
	var result models.XTCETelemetryDefinition

	collection := mt.db.Collection("xtce_containers")

	// Debugging: Log the exact filter used for the query
	filter := bson.M{"satellite_id": satelliteID, "xtce_version": xtceVersion}
	log.Printf("Querying MongoDB with filter: %+v\n", filter)

	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Printf("No document found for Satellite: %s, Version: %s\n", satelliteID, xtceVersion)
		return nil, nil
	} else if err != nil {
		log.Printf("MongoDB query error: %v\n", err)
		return nil, err
	}

	// Debug: Print retrieved document
	log.Printf("Retrieved XTCE Document: %+v\n", result)

	// Debug: Check if APID is stored correctly
	for _, container := range result.Containers {
		log.Printf("Checking container: %+v\n", container)
		if container.APID == apid {
			log.Printf("APID %d found\n", apid)
			return &container, nil
		}
	}

	log.Printf("Warning: APID %d not found within retrieved document\n", apid)
	return nil, nil
}
