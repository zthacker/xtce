package decoder

import (
	"encoding/binary"
	"xtcedata/database/metadata"
	"xtcedata/database/telemetrydata"
	"xtcedata/extractor"
)

type XTCEDecoder struct {
	extractor.Extractor
	telemetrydata.TelemetryDatabase
	metadata.MetaDataDatabase
	//in reality, we'd probably have the redis cache here, or some sort of in-memory db
}

func NewXTCEDecoder(ext extractor.Extractor, tmDB telemetrydata.TelemetryDatabase, mdDB metadata.MetaDataDatabase) *XTCEDecoder {
	return &XTCEDecoder{ext, tmDB, mdDB}
}

func (xtce *XTCEDecoder) DecodePacket(packet []byte) (map[string]interface{}, error) {
	id := packet[:7]
	apid := binary.BigEndian.Uint32(packet[8:12])
	version, err := xtce.GetMetaData(string(id))
	if err != nil {
		return nil, err
	}

	container, err := xtce.GetContainerByAPID(string(id), version, apid)
	if err != nil {
		return nil, err
	}

	decodedPacket, err := xtce.Extract(packet, container)
	if err != nil {
		return nil, err
	}

	return decodedPacket, nil
}
