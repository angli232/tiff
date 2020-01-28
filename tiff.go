/*
Package tiff implements structures and functionality to encode and decode TIFF-like files.

To iterate through images in a multi-page TIFF and read image data:
	r, err := tiff.NewReader(f)
	if err != nil {
		// handle error
	}
	iter := r.Iter()
	for iter.Next() {
		im := iter.Image()
		width, height := im.WidthHeight()
		samplesPerPixel := im.SamplesPerPixel()
		if width * height * samplesPerPixel == 0 {
			continue
		}
		switch im.DataType() {
		case tiff.Uint8:
			buf := make([]uint8, width*height*samplesPerPixel)
			err = im.DecodeImage(buf)
		case tiff.Uint16:
			buf := make([]uint16, width*height*samplesPerPixel)
			err = im.DecodeImage(buf)
		}
	}
	if err := iter.Err(); err != nil {
		// handle error
	}

To read SubIFDs:
	if im.SubImage != nil {
		for _, subim := range im.SubImage {
			// access subim
		}
	}

To read any TIFF tag, EXIF tag or GPS tag:
	bisPerSample, ok := im.Tag[tiff.TagBitsPerSample].UintSlice()
	make, ok := im.Tag[tiff.TagMake].String()
	model, ok := im.Tag[tiff.TagModel].String()

	exposure, ok := im.Exif.Tag[tiff.TagExposureTime].Rational()
	exposureSecond := float64(exposure[0]) / float64(exposure[1])
	lensModel, ok := im.Exif.Tag[tiff.TagLensModel].String()

To encode a TIFF image with one or multiple pages:
	enc, err := tiff.NewEncoder(f)
	if err != nil {
		// handle error
	}
	// Default:
	// enc.SetByteOrder(binary.LittleEndian)
	// enc.SetVersion(tiff.VersionClassicTIFF)
	im := enc.NewImage()
	im.SetWidthHeight(400, 300)
	im.SetPixelFormat(tiff.PhotometricRGB, 3, []uint16{8, 8, 8})
	im.SetCompression(tiff.CompressionNone)

	im.SetTag(tiff.TagImageDescription, tiff.TagTypeASCII, "your image description")
	im.SetTag(tiff.TagDateTime, tiff.TagTypeASCII, "")
	im.SetTag(tiff.TagXResolution, tiff.TagTypeRational, [2]uint32{})
	im.SetTag(tiff.TagYResolution, tiff.TagTypeRational, [2]uint32{})
	im.SetTag(tiff.TagResolutionUnit, , tiff.TagTypeShort, 1)
	im.SetExifTag(tiff.ExifTagExposureTime, tiff.TagTypeRational, [2]uint32{})
	err = im.EncodeImage(buf)

	// To write another image
	im = w.NewImage()
	// ...
	err = im.EncodeImage(buf)

	w.Close()

To encode a TIFF image with sub-images (SubIFDs).
	im := w.NewImage()
	// ...
	// AddSubImage needs to be performed before calling EncodeImage()
	subim1 := im.AddSubImage()
	subim1.SetWidthHeight()
	// ...
	err = im.EncodeImage(buf)
	err = subim1.EncodeImage(buf)

	w.Close()
*/
package tiff

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strconv"
)

const maxInt = uint64((^uint(0)) >> 1)

// A FormatError reports that the input is not a valid TIFF image.
type FormatError string

func (e FormatError) Error() string {
	return "tiff: invalid format: " + string(e)
}

// An UnsupportedError reports that the input uses a valid but
// unimplemented feature.
type UnsupportedError string

func (e UnsupportedError) Error() string {
	return "tiff: unsupported feature: " + string(e)
}

var (
	ErrInvalidHeader FormatError = "invalid header"
)

type Header struct {
	ByteOrder      binary.ByteOrder
	Version        Version
	OffsetFirstIFD int64
}

func decodeHeader(r io.ReaderAt) (*Header, error) {
	buf := make([]byte, 8)
	_, err := r.ReadAt(buf, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot read header: %s", err)
	}

	h := &Header{}

	switch string(buf[0:2]) {
	case "II":
		h.ByteOrder = binary.LittleEndian
	case "MM":
		h.ByteOrder = binary.BigEndian
	default:
		return nil, ErrInvalidHeader
	}

	h.Version = Version(h.ByteOrder.Uint16(buf[2:4]))
	switch h.Version {
	case VersionClassicTIFF:
		h.OffsetFirstIFD = int64(h.ByteOrder.Uint32(buf[4:8]))
		return h, nil

	case VersionBigTIFF:
		bytesOfOffset := h.ByteOrder.Uint16(buf[4:6])
		zero := h.ByteOrder.Uint16(buf[6:8])
		if bytesOfOffset != 8 || zero != 0 {
			return nil, ErrInvalidHeader
		}
		// Read offset to first IFD
		_, err := r.ReadAt(buf, 8)
		if err != nil {
			return nil, fmt.Errorf("cannot read header: %s", err)
		}
		offsetFirstIFD := h.ByteOrder.Uint64(buf[0:8])
		if offsetFirstIFD > math.MaxInt64 {
			// Although allowed by the BigTIFF spec,
			// a offset that overflows int64 is unlikely to be valid,
			// and is unsupported by Go anyways.
			return nil, ErrInvalidHeader
		}
		h.OffsetFirstIFD = int64(offsetFirstIFD)
		return h, nil

	default:
		return nil, ErrInvalidHeader
	}
}

// encodeBytes returns encoded Header, and panics on invalid Header
func (h *Header) encodeBytes() []byte {
	switch h.Version {
	case VersionClassicTIFF:
		buf := make([]byte, 8)

		// Offset 0: Byte order indication
		if h.ByteOrder == binary.LittleEndian {
			buf[0] = 'I'
			buf[1] = 'I'
		} else {
			buf[0] = 'M'
			buf[1] = 'M'
		}

		// Offset 2: Version 42
		h.ByteOrder.PutUint16(buf[2:4], uint16(h.Version))

		// Offset 4: Offset to first IFD
		if h.OffsetFirstIFD <= 0 {
			panic("tiff: invalid OffsetFirstIFD")
		}
		if h.OffsetFirstIFD > math.MaxUint32 {
			panic("tiff: OffsetFirstIFD overflows uint32")
		}
		h.ByteOrder.PutUint32(buf[4:8], uint32(h.OffsetFirstIFD))

		return buf

	case VersionBigTIFF:
		buf := make([]byte, 16)

		// Offset 0: Byte order indication
		if h.ByteOrder == binary.LittleEndian {
			buf[0] = 'I'
			buf[1] = 'I'
		} else {
			buf[0] = 'M'
			buf[1] = 'M'
		}

		// Offset 2: Version 43
		h.ByteOrder.PutUint16(buf[2:4], uint16(h.Version))

		// Offset 4: Bytesize of offsets
		h.ByteOrder.PutUint16(buf[4:6], 8)

		// Offset 6: 0
		h.ByteOrder.PutUint16(buf[6:8], 0)

		// Offset 8: Offset to first IFD
		if h.OffsetFirstIFD <= 0 {
			panic("tiff: invalid OffsetFirstIFD")
		}
		h.ByteOrder.PutUint64(buf[8:16], uint64(h.OffsetFirstIFD))

		return buf

	default:
		panic("tiff: invalid TIFF Version")
	}
}

// Version indicates TIFF version.
type Version uint16

const (
	VersionClassicTIFF Version = 42
	VersionBigTIFF     Version = 43
)

func (v Version) String() string {
	switch v {
	case VersionClassicTIFF:
		return "ClassicTIFF"
	case VersionBigTIFF:
		return "BigTIFF"
	default:
		return strconv.Itoa(int(v))
	}
}

type Rational [2]uint32

func NewRational(a uint32, b uint32) Rational {
	return [2]uint32{a, b}
}

func (r Rational) String() string {
	return fmt.Sprintf("%d/%d", r[0], r[1])
}

type SRational [2]int32

func NewSRational(a int32, b int32) SRational {
	return [2]int32{a, b}
}

func (r SRational) String() string {
	return fmt.Sprintf("%d/%d", r[0], r[1])
}

type DataType int

const (
	InvalidDataType DataType = iota
	Uint8
	Uint16
)
