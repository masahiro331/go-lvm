package types

type Volume struct {
	LabelHeader PhysicalVolumeLabelHeader
	Header      PhysicalVolumeHeader

	MetadataArea []MetadataArea
}

type PhysicalVolumeLabelHeader struct {
	Signature     [8]byte
	SectorNumber  int64
	Checksum      int32
	DataOffset    int32
	TypeIndicator int64
}

type PhysicalVolumeHeader struct {
	PhysicalVolumeIdentifier [32]byte
	PhysicalVolumeSize       int64
	DataAreaDescriptor       []DataAreaDescriptor
	MetaDataAreaDescriptor   []DataAreaDescriptor
}

type DataAreaDescriptor struct {
	DataAreaOffset int64
	DataAreaSize   int64
}

type MetadataArea struct {
	Header   MetadataAreaHeader
	Metadata Metadata
}

type MetadataAreaHeader struct {
	Checksum               uint32
	Signature              [16]byte
	Version                int32
	MetadataAreaOffset     int64
	MetadataAreaSize       int64
	RawLocationDescriptors [4]RawLocationDescriptor
	_                      [376]byte
}

type RawLocationDescriptor struct {
	DataAreaOffset int64
	DataAreaSize   int64
	Checksum       uint32
	Flags          int32
}
