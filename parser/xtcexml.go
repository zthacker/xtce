package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"xtcedata/models"
)

type XTCEParserForXML struct {
}

func NewXTCEParserForXML() *XTCEParserForXML {
	return &XTCEParserForXML{}
}

func (xt *XTCEParserForXML) Parse(filename string) (*models.XTCETelemetryDefinition, error) {
	// Read XML file
	xmlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading XTCE file: %v", err)
	}

	// Unmarshal into Go struct
	var spaceSystem models.SpaceSystem
	err = xml.Unmarshal(xmlFile, &spaceSystem)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling XTCE XML: %v", err)
	}

	// Extract satellite metadata
	satelliteID := spaceSystem.Name
	xtceVersion := spaceSystem.Header.Version

	// Build a map of parameters to store DataType and Length
	paramDetails := make(map[string]struct {
		DataType string
		Length   int
	})

	// Read parameters from ParameterSet
	for _, param := range spaceSystem.TelemetryMetaData.ParameterSet.Parameters {
		paramDetails[param.Name] = struct {
			DataType string
			Length   int
		}{
			DataType: param.DataTypeRef,
			Length:   getTypeLength(param.DataTypeRef), // Determine Length based on type
		}
	}

	var containers []models.Container
	for _, seqContainer := range spaceSystem.TelemetryMetaData.ContainerSet.Containers {
		apid := seqContainer.BaseContainer.RestrictionCriteria.Comparison.Value

		// Extract fields (parameters) for this container
		var fields []models.Field
		for _, entry := range seqContainer.EntryList.ParameterEntries {
			paramDetail, exists := paramDetails[entry.ParameterRef]
			if !exists {
				return nil, fmt.Errorf("missing parameter definition for %s", entry.ParameterRef)
			}

			fields = append(fields, models.Field{
				Name:      entry.ParameterRef,
				StartByte: entry.LocationInContainerInBits / 8,
				Length:    paramDetail.Length,                
				DataType:  paramDetail.DataType,                
			})
		}

		// Append a single container with all its fields
		containers = append(containers, models.Container{
			APID:   apid,
			Name:   seqContainer.Name,
			Fields: fields,
		})
	}

	return &models.XTCETelemetryDefinition{
		SatelliteID: satelliteID,
		XTCEVersion: xtceVersion,
		Containers:  containers,
	}, nil
}
