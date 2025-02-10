package models

type Container struct {
	APID   uint32  `bson:"apid"`   // MongoDB field for Application Process ID
	Name   string  `bson:"name"`   // Name of the container
	Fields []Field `bson:"fields"` // List of telemetry fields
}

type Field struct {
	Name      string `bson:"name"`       // Field name (e.g., "Temperature")
	StartByte int    `bson:"start_byte"` // Byte offset in the packet
	Length    int    `bson:"length"`     // Number of bytes
	DataType  string `bson:"data_type"`  // Type (e.g., "uint32", "float")
}

type XTCETelemetryDefinition struct {
	SatelliteID string      `bson:"satellite_id"`
	XTCEVersion string      `bson:"xtce_version"`
	Containers  []Container `bson:"containers"`
}
