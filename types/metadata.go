package types

import (
	"strconv"
	"strings"
)

type Type int

const (
	Unknown Type = iota
	String
	Number
	List
)

type Metadata struct {
	Blocks *[]Block `@@*`
}

type Block struct {
	Key    string   `@Ident`
	Value  *Value   `|"=" @@`
	Blocks *[]Block `| "{" @@* "}"`
}

type Value struct {
	Str  *string  `@String`
	Num  *int     `| @Number`
	List []*Value `| "[" ( @@ ( "," @@ )* )? "]"`
}

func (v *Value) Type() Type {
	switch {
	case v.Str != nil:
		return String
	case v.Num != nil:
		return Number
	case v.List != nil:
		return List
	default:
		return Unknown
	}
}

type PhysicalVolumeSection struct {
	PhysicalVolumes map[string]PhysicalVolume
}

type LogicalVolumeSection struct {
	LogicalVolumes map[string]LogicalVolume
}

type LogicalVolume struct {
	Segments     map[string]Segment
	ID           string
	Status       []string
	Flags        []string
	CreationHost string
	CreationTime int
	SegmentCount int
}

type Segment struct {
	StartExtent int
	ExtentCount int
	Type        string
	StripeCount int
	Stripes     []string
}

type MainSection struct {
	ID                    string
	SeqNo                 int
	Format                string
	Status                []string
	Flags                 []string
	ExtentSize            int
	MaxLV                 int
	MaxPV                 int
	MetadataCopies        int
	PhysicalVolumeSection PhysicalVolumeSection
	LogicalVolumeSection  LogicalVolumeSection
}

type PhysicalVolume struct {
	ID      string
	Device  string
	Status  []string
	Flags   []string
	DevSize int
	PeStart int
	PeCount int
}

func ParseMainSection(metadata *Metadata) MainSection {
	blocks := metadata.mainSection()
	if blocks == nil {
		return MainSection{}
	}

	m := MainSection{}
	walkBlocks(blocks, func(key string, b Block) {
		switch key {
		case ID:
			m.ID = *b.Value.Str
		case SeqNo:
			m.SeqNo = *b.Value.Num
		case Format:
			m.Format = *b.Value.Str
		case Status:
			if b.Value == nil {
				break
			}
			m.Status = b.Value.forceStrList()
		case Flags:
			if b.Value == nil {
				break
			}
			m.Flags = b.Value.forceStrList()
		case ExtentSize:
			m.ExtentSize = *b.Value.Num
		case MaxLV:
			m.MaxLV = *b.Value.Num
		case MaxPV:
			m.MaxPV = *b.Value.Num
		case MetadataCopies:
			m.MetadataCopies = *b.Value.Num
		case PhysicalVolumes:
			m.PhysicalVolumeSection = parsePhysicalVolumeSection(b.Blocks)
		case LogicalVolumes:
			m.LogicalVolumeSection = parseLogicalVolumeSection(b.Blocks)
		}
	})
	return m
}

func (m Metadata) mainSection() *[]Block {
	var key string
	for _, b := range *m.Blocks {
		if b.Key != "" {
			key = b.Key
		} else {
			if b.Blocks != nil && key != "" {
				return b.Blocks
			}
			key = ""
		}
	}
	return nil
}

func parseLogicalVolume(blocks *[]Block) LogicalVolume {
	lv := LogicalVolume{}
	segmentMap := map[string]Segment{}

	walkBlocks(blocks, func(key string, b Block) {
		switch key {
		case ID:
			lv.ID = *b.Value.Str
		case Status:
			if b.Value == nil {
				break
			}
			lv.Status = b.Value.forceStrList()
		case Flags:
			if b.Value == nil {
				break
			}
			lv.Flags = b.Value.forceStrList()
		case CreationTime:
			lv.CreationTime = *b.Value.Num
		case CreationHost:
			lv.CreationHost = *b.Value.Str
		case SegmentCount:
			lv.SegmentCount = *b.Value.Num
		default:
			if strings.HasPrefix(key, "segment") && b.Blocks != nil {
				segmentMap[key] = parseSegment(b.Blocks)
			}
		}
		lv.Segments = segmentMap
	})
	return lv
}

func parseLogicalVolumeSection(blocks *[]Block) LogicalVolumeSection {
	lvMap := map[string]LogicalVolume{}
	lvs := LogicalVolumeSection{
		LogicalVolumes: lvMap,
	}

	walkBlocks(blocks, func(lvName string, b Block) {
		if b.Blocks == nil {
			return
		}
		lvMap[lvName] = parseLogicalVolume(b.Blocks)
	})
	return lvs
}

func parseSegment(blocks *[]Block) Segment {
	var sg Segment
	walkBlocks(blocks, func(key string, b Block) {
		switch key {
		case StartExtent:
			sg.StartExtent = *b.Value.Num
		case ExtentCount:
			sg.ExtentCount = *b.Value.Num
		case SegmentType:
			sg.Type = *b.Value.Str
		case StripeCount:
			sg.StripeCount = *b.Value.Num
		case Stripes:
			if b.Value == nil {
				break
			}
			sg.Stripes = b.Value.forceStrList()
		}
	})
	return sg
}

func parsePhysicalVolume(blocks *[]Block) PhysicalVolume {
	pv := PhysicalVolume{}
	walkBlocks(blocks, func(key string, b Block) {
		switch key {
		case ID:
			pv.ID = *b.Value.Str
		case Device:
			pv.Device = *b.Value.Str
		case Status:
			if b.Value == nil {
				break
			}
			pv.Status = b.Value.forceStrList()
		case Flags:
			if b.Value == nil {
				break
			}
			pv.Flags = b.Value.forceStrList()
		case DevSize:
			pv.DevSize = *b.Value.Num
		case PeStart:
			pv.PeStart = *b.Value.Num
		case PeCount:
			pv.PeCount = *b.Value.Num
		}
	})
	return pv
}

func parsePhysicalVolumeSection(blocks *[]Block) PhysicalVolumeSection {
	pvMap := map[string]PhysicalVolume{}
	pvs := PhysicalVolumeSection{
		PhysicalVolumes: pvMap,
	}
	walkBlocks(blocks, func(pvName string, block Block) {
		if strings.HasPrefix(pvName, "pv") {
			if block.Blocks == nil {
				return
			}
			pvMap[pvName] = parsePhysicalVolume(block.Blocks)
		}
	})
	return pvs
}

func walkBlocks(blocks *[]Block, fn func(key string, block Block)) {
	if blocks == nil {
		return
	}

	var key string
	for _, b := range *blocks {
		if b.Key != "" {
			key = b.Key
			continue
		}
		if key != "" && (b.Value == nil && b.Blocks == nil) {
			continue
		}
		fn(key, b)
		key = ""
	}
}

func (v *Value) forceStrList() []string {
	var ret []string
	for _, item := range v.List {
		if item.Type() == String {
			ret = append(ret, *item.Str)
		} else if item.Type() == Number {
			ret = append(ret, strconv.Itoa(*item.Num))
		}
	}
	return ret
}
