package go_lvm

import (
	"encoding/binary"
	"fmt"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/v2"
	"github.com/masahiro331/go-lvm/types"
	"golang.org/x/xerrors"
	"io"
)

const SectorSize = 512

type Driver interface {
	Next() (types.Volume, error)
}

func NewDriver(rs io.ReadSeeker) (Driver, error) {
	rs.Seek(SectorSize, io.SeekStart)
	vlh, err := NewPhysicalVolumeLabelHeader(rs)
	if err != nil {
		return nil, xerrors.Errorf("failed to create physical label volume header: %w", err)
	}
	vh, err := NewPhysicalVolumeHeader(rs)
	if err != nil {
		return nil, xerrors.Errorf("failed to create physical volume header: %w", err)
	}
	var v types.Volume
	v.LabelHeader = vlh
	v.Header = vh
	fmt.Printf("%+v\n", v)

	for _, descriptor := range v.Header.MetaDataAreaDescriptor {
		h, err := parseMetadataArea(rs, descriptor)
		if err != nil {
			return nil, xerrors.Errorf("failed to parse metadata area header: %w", err)
		}
		fmt.Printf("%+v\n", h)
	}

	return nil, nil
}

func NewPhysicalVolumeHeader(r io.Reader) (types.PhysicalVolumeHeader, error) {
	h := types.PhysicalVolumeHeader{}

	if err := binary.Read(r, binary.LittleEndian, &h.PhysicalVolumeIdentifier); err != nil {
		return types.PhysicalVolumeHeader{}, xerrors.Errorf("failed to read physical volume header identifier: %w", err)
	}
	if err := binary.Read(r, binary.LittleEndian, &h.PhysicalVolumeSize); err != nil {
		return types.PhysicalVolumeHeader{}, xerrors.Errorf("failed to read physical volume header size: %w", err)
	}

	var err error
	h.DataAreaDescriptor, err = parseDataAreaDescriptors(r)
	if err != nil {
		return types.PhysicalVolumeHeader{}, xerrors.Errorf("failed to parse data area descriptor: %w", err)
	}

	h.MetaDataAreaDescriptor, err = parseDataAreaDescriptors(r)
	if err != nil {
		return types.PhysicalVolumeHeader{}, xerrors.Errorf("failed to parse meta data area descriptor: %w", err)
	}
	return h, nil
}

func parseMetadataArea(r io.ReadSeeker, descriptor types.DataAreaDescriptor) (types.MetadataArea, error) {
	_, err := r.Seek(descriptor.DataAreaOffset, io.SeekStart)
	if err != nil {
		return types.MetadataArea{}, xerrors.Errorf("failed to seek to metadata area: %w", err)
	}
	var h types.MetadataArea
	if err := binary.Read(r, binary.LittleEndian, &h.Header); err != nil {
		return types.MetadataArea{}, xerrors.Errorf("failed to read metadata area header: %w", err)
	}

	for _, d := range h.Header.RawLocationDescriptors {
		if d.DataAreaSize == 0 {
			continue
		}
		offset := h.Header.MetadataAreaOffset + d.DataAreaOffset
		r.Seek(offset, io.SeekStart)
		h.Metadata, err = parseMetadata(io.LimitReader(r, d.DataAreaSize))
		if err != nil {
			return types.MetadataArea{}, xerrors.Errorf("failed to parse metadata: %w", err)
		}
	}

	return h, nil
}

var (
	Lexer = lexer.MustSimple([]lexer.SimpleRule{
		{"Comment", `(?:#|//)[^\n]*\n?`},
		{"Number", `(?:\d*\.)?\d+`},
		{"Ident", `[0-9a-zA-Z_-]+`},
		{"String", `"(\\"|[^"])*"`},
		{"Punct", `[[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
		{"Whitespace", `[ \t\n\r]+`},
	})
	parser = participle.MustBuild[types.Metadata](
		participle.Lexer(Lexer),
		participle.Elide("Comment", "Whitespace"),
		participle.UseLookahead(2),
	)
)

func parseMetadata(r io.Reader) (types.Metadata, error) {
	expr, err := parser.Parse("", r)
	if err != nil {
		return types.Metadata{}, err
	}

	return *expr, nil
}

func parseDataAreaDescriptors(r io.Reader) ([]types.DataAreaDescriptor, error) {
	var ds []types.DataAreaDescriptor
	for {
		var d types.DataAreaDescriptor
		if err := binary.Read(r, binary.LittleEndian, &d); err != nil {
			return nil, xerrors.Errorf("failed to read data area descriptor: %w", err)
		}
		if d.DataAreaOffset == 0 && d.DataAreaSize == 0 {
			break
		}
		ds = append(ds, d)
	}
	return ds, nil
}

func NewPhysicalVolumeLabelHeader(r io.Reader) (types.PhysicalVolumeLabelHeader, error) {
	h := types.PhysicalVolumeLabelHeader{}

	if err := binary.Read(r, binary.LittleEndian, &h); err != nil {
		return types.PhysicalVolumeLabelHeader{}, xerrors.Errorf("failed to read physical volume label header: %w", err)
	}
	return h, nil
}
