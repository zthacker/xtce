package parser

func getTypeLength(dataType string) int {
	switch dataType {
	case "uint8", "int8":
		return 1
	case "uint16", "int16":
		return 2
	case "uint32", "int32", "float":
		return 4
	case "uint64", "int64", "double":
		return 8
	default:
		return 4 // Default size (could be improved)
	}
}
