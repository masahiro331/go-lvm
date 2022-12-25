package types

const (
	// Section Parameters

	// ID is volume group identifier (VG UUID)
	// Contains an ASCII string in the following format: fg1fKZ-xoHz-CfAD-yQPx-l2HL-Y7kA-9kJ9LD
	ID = "id"

	// SeqNo is metadata sequence number
	SeqNo = "seqno"

	Format = "format"

	// Status flags are contains a list of strings.
	Status = "status"

	// Flags Contains a list of strings. See section: Flags
	Flags = "flags"

	// ExtentSize is size of an extent
	// The value contains the number of sectors
	// According to [REDHAT] the sector size should be 512 bytes
	ExtentSize = "extent_size"

	// MaxLV is maximum number of logical volumes
	MaxLV = "max_lv"

	// MaxPV is Maximum number of physical volumes
	MaxPV = "max_pv"

	// MetadataCopies is number of metadata copies
	MetadataCopies = "metadata_copies"

	// Device is filename contains an ASCII string e.g. /dev/dm-0"
	Device = "device"

	// DevSize is physical volume size including non-usable space
	// The value contains the number of sectors
	// According to [REDHAT] the sector size should be 512 bytes
	DevSize = "dev_size"

	// PeStart is the start extent of the physical volume,
	// contains an offset relative to the start of the physical volume
	PeStart = "pe_start"

	// PeCount is number of (allocated) extents in the physical volume
	PeCount = "pe_count"

	CreationHost = "creation_host"
	CreationTime = "creation_time"

	// SegmentCount is number of segment subsections
	SegmentCount = "segment_count"

	// StartExtent is start extent of the segment
	// The value contains the number of extents
	// The number is relative to the start of the segment
	StartExtent = "start_extent"

	// ExtentCount is number of extents in the segment (or current logical extent)
	ExtentCount = "extent_count"

	// SegmentType is segment type
	// See. https://github.com/libyal/libvslvm/blob/main/documentation/Logical%20Volume%20Manager%20(LVM)%20format.asciidoc#segment_types
	SegmentType = "type"

	// StripeCount is number of stripes in the segment 1 => linear
	StripeCount = "stripe_count"

	// Stripes is stripes list. e.g. ["pv0", 9861]
	Stripes = "stripes"

	PhysicalVolumes = "physical_volumes"
	LogicalVolumes  = "logical_volumes"
)
