package metadata

type MetaDataDatabase interface {
	InsertMetaData(satelliteID string, version string) error
	GetMetaData(satelliteID string) (string, error)
}
