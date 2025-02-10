package parser

import "xtcedata/models"

type XTCEParser interface {
	Parse(filename string) (*models.XTCETelemetryDefinition, error)
}

//add other parsers here -- like for Cosmos, CCSDS, etc
