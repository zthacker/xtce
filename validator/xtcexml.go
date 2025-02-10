package validator

import (
	"bytes"
	"fmt"
	"os"
)

type XTCEValidatorForXML struct {
}

func NewXTCEValidatorForXML() *XTCEValidatorForXML {
	return &XTCEValidatorForXML{}
}

func (X XTCEValidatorForXML) Validate(filename string) error {
	xmlFile, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading XTCE file: %v", err)
	}

	// Quick check if required tags exist
	if !bytes.Contains(xmlFile, []byte("<SpaceSystem")) || !bytes.Contains(xmlFile, []byte("<TelemetryMetaData")) {
		return fmt.Errorf("invalid XTCE file: missing required elements")
	}

	return nil
}
