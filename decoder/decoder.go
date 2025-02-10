package decoder

type MCSDecoder interface {
	DecodePacket(packet []byte) (map[string]interface{}, error)
}
