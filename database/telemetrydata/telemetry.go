package telemetrydata

import "xtcedata/models"

type TelemetryDatabase interface {
	InsertXTCEContainer(xtceData *models.XTCETelemetryDefinition) error
	GetContainerByAPID(satelliteID, xtceVersion string, apid uint32) (*models.Container, error)
}
