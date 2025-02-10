package extractor

import (
	"xtcedata/models"
)

type Extractor interface {
	Extract(packet []byte, container *models.Container) (map[string]interface{}, error)
}
