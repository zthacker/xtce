package extractor

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"xtcedata/models"
)

type XTCEExtractor struct {
}

func NewXTCEExtractor() *XTCEExtractor {
	return &XTCEExtractor{}
}

func (X *XTCEExtractor) Extract(packet []byte, container *models.Container) (map[string]interface{}, error) {
	decodedData := make(map[string]interface{})
	for _, field := range container.Fields {
		data, err := extractField(packet, field.StartByte, field.Length, field.DataType)
		if err != nil {
			return nil, err
		}
		decodedData[field.Name] = data
	}
	return decodedData, nil
}

// ExtractField reads the packet data based on byte offset, length, and data type
func extractField(packet []byte, startByte int, length int, dataType string) (interface{}, error) {
	if len(packet) < startByte+length {

		return nil, fmt.Errorf("Error: Packet too short for expected field extraction at %d (length %d)\n", startByte, length)
	}

	// Slice the relevant bytes from the packet
	dataBytes := packet[startByte : startByte+length]

	// Decode based on the expected data type
	switch dataType {
	case "uint8":
		return dataBytes[0], nil

	case "uint16":
		return binary.BigEndian.Uint16(dataBytes), nil

	case "uint32":
		return binary.BigEndian.Uint32(dataBytes), nil

	case "uint64":
		return binary.BigEndian.Uint64(dataBytes), nil

	case "int8":
		return int8(dataBytes[0]), nil

	case "int16":
		return int16(binary.BigEndian.Uint16(dataBytes)), nil

	case "int32":
		return int32(binary.BigEndian.Uint32(dataBytes)), nil

	case "int64":
		return int64(binary.BigEndian.Uint64(dataBytes)), nil

	case "float":
		bits := binary.BigEndian.Uint32(dataBytes)
		return math.Float32frombits(bits), nil

	case "double":
		bits := binary.BigEndian.Uint64(dataBytes)
		return math.Float64frombits(bits), nil

	default:
		log.Printf("Warning: Unknown data type '%s'\n", dataType)
		return nil, nil
	}
}
