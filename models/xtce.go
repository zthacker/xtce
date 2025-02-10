package models

import "encoding/xml"

type SpaceSystem struct {
	XMLName           xml.Name          `xml:"SpaceSystem"`
	Name              string            `xml:"name,attr"`
	Header            Header            `xml:"Header"`
	TelemetryMetaData TelemetryMetaData `xml:"TelemetryMetaData"`
}

type Header struct {
	Version string `xml:"Version"`
}

type TelemetryMetaData struct {
	ParameterSet ParameterSet `xml:"ParameterSet"`
	ContainerSet ContainerSet `xml:"ContainerSet"`
}

type ParameterSet struct {
	Parameters []Parameter `xml:"Parameter"`
}

type Parameter struct {
	Name        string `xml:"name,attr"`
	DataTypeRef string `xml:"DataTypeRef"`
	Description string `xml:"Description"`
}

type ContainerSet struct {
	Containers []SequenceContainer `xml:"SequenceContainer"`
}

type SequenceContainer struct {
	Name          string        `xml:"name,attr"`
	BaseContainer BaseContainer `xml:"BaseContainer"`
	EntryList     EntryList     `xml:"EntryList"`
}

type BaseContainer struct {
	RestrictionCriteria RestrictionCriteria `xml:"RestrictionCriteria"`
}

type RestrictionCriteria struct {
	Comparison Comparison `xml:"Comparison"`
}

type Comparison struct {
	ParameterRef string `xml:"ParameterRef"`
	Value        uint32 `xml:"Value"`
}

type EntryList struct {
	ParameterEntries []ParameterRefEntry `xml:"ParameterRefEntry"`
}

type ParameterRefEntry struct {
	ParameterRef              string `xml:"ParameterRef"`
	LocationInContainerInBits int    `xml:"LocationInContainerInBits"`
}
