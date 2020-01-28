package tiff

import (
	"fmt"
	"io"
	"sort"
)

// Image represents a Image File Directory (IFD).
type Image struct {
	Tag map[TagID]*Tag

	SubImage         []*Image
	Exif             *ExifIFD
	GPS              *GPSIFD
	Interoperability *InteroperabilityIFD

	// Filled by Encoder or Decoder
	Header     *Header     // Header of the TIFF file
	Offset     int64       // Offset of this IFD
	OffsetNext int64       // Offset of the next IFD
	rd         io.ReaderAt // Reader to access the whole TIFF file
	enc        *Encoder    // To access WriterAt and current write offset
}

// TagID returns a sorted slice of TagIDs of the Image.
func (im *Image) TagID() []TagID {
	list := make([]TagID, 0, len(im.Tag))
	for id := range im.Tag {
		list = append(list, id)
	}
	sort.Sort(byTagID(list))
	return list
}

func (im *Image) SetTag(id TagID, tagType TagType, value interface{}) error {
	if im.Tag == nil {
		im.Tag = make(map[TagID]*Tag)
	}
	tag, err := NewTag(id, tagType, value, im.Header)
	if err != nil {
		return err
	}
	im.Tag[id] = tag
	return nil
}

func (im *Image) SetWidthHeight(width int, height int) {
	im.SetTag(TagImageWidth, TagTypeLong, uint32(width))
	im.SetTag(TagImageLength, TagTypeLong, uint32(height))
}

func (im *Image) SetPixelFormat(photometric int, samplePerPixel int, bitsPerSample []int) {
	im.SetTag(TagPhotometric, TagTypeLong, uint32(photometric))
	im.SetTag(TagSamplesPerPixel, TagTypeLong, uint32(samplePerPixel))
	im.SetTag(TagBitsPerSample, TagTypeLong, bitsPerSample)
}

func (im *Image) SetCompression(compression int) {
	im.SetTag(TagCompression, TagTypeShort, uint16(compression))
}

func (im *Image) SetRowsPerStrip(rowsPerStrip int) {
	im.SetTag(TagRowsPerStrip, TagTypeLong, uint32(rowsPerStrip))
}

func (im *Image) SetTileWidthHeight(tileWidth, tileHeight int) {
	im.SetTag(TagTileWidth, TagTypeLong, uint32(tileWidth))
	im.SetTag(TagTileLength, TagTypeLong, uint32(tileHeight))
}

// WidthHeight returns the width and height of the image.
func (im *Image) WidthHeight() (width int, height int) {
	v, ok := im.Tag[TagImageWidth].Uint()
	if !ok {
		return
	}
	width = int(v)
	v, ok = im.Tag[TagImageLength].Uint()
	if !ok {
		return
	}
	height = int(v)
	return
}

// SamplesPerPixel returns the number of samples (color channels) per pixel.
func (im *Image) SamplesPerPixel() int {
	v, ok := im.Tag[TagSamplesPerPixel].Uint()
	if !ok {
		return 0
	}
	return int(v)
}

// PlanarConfig returns how the components of each pixel are stored.
func (im *Image) PlanarConfig() int {
	e := im.Tag[TagPlanarConfig]
	if e == nil {
		// Default = 1. (TIFF 6.0 Page 38)
		return 1
	}
	v, _ := e.Uint()
	return int(v)
}

// BitsPerSample returns the number of bits of each sample of a pixel.
func (im *Image) BitsPerSample() []int {
	v, ok := im.Tag[TagBitsPerSample].UintSlice()
	if !ok {
		return []int{}
	}
	r := make([]int, len(v))
	for i := 0; i < len(v); i++ {
		r[i] = int(v[i])
	}
	return r
}

// BitDepth returns bit depth of samples.
//
// It returns 0, if valid BitsPerSample tag is not found,
// or if different components have different bit depth.
func (im *Image) BitDepth() int {
	bitsPerSample, ok := im.Tag[TagBitsPerSample].UintSlice()
	if !ok || len(bitsPerSample) == 0 {
		return 0
	}

	samplePerPixel := im.SamplesPerPixel()
	if len(bitsPerSample) != samplePerPixel {
		return 0
	}

	bitDepth := int(bitsPerSample[0])
	for _, bits := range bitsPerSample {
		if int(bits) != bitDepth {
			return 0
		}
	}
	return bitDepth
}

// DataType returns the data type user should expect when decoding the image.
func (im *Image) DataType() DataType {
	bitDepth := im.BitDepth()
	switch bitDepth {
	case 1, 8:
		return Uint8
	case 16:
		return Uint16
	default:
		return InvalidDataType
	}
}

// Compression returns the compression scheme used on the image data
func (im *Image) Compression() (compression int) {
	v, ok := im.Tag[TagCompression].Uint()
	if !ok {
		// Compression does not have a default value,
		// but some tools interpret missing value as CompressionNone.
		return CompressionNone
	}
	return int(v)
}

// EncodeTags returns encoded tags with requested IFD offset.
func (im *Image) EncodeTags(offset int64) ([]byte, error) {
	im.Offset = offset

	ids := im.TagID()
	if len(ids) == 0 {
		return nil, fmt.Errorf("at least one tag is required")
	}
	if im.Header == nil {
		return nil, fmt.Errorf("missing Header")
	}

	var bytesOfNumTags int
	var bytesOfEntry int
	var byteOfOffset int
	switch im.Header.Version {
	case VersionClassicTIFF:
		bytesOfNumTags = 2
		bytesOfEntry = 12
		byteOfOffset = 4
	case VersionBigTIFF:
		bytesOfNumTags = 8
		bytesOfEntry = 20
		byteOfOffset = 8
	default:
		panic("invalid TIFF version in header")
	}

	var byteOfExtended int
	for _, tag := range im.Tag {
		if len(tag.Data) > byteOfOffset {
			byteOfExtended += len(tag.Data)
		}
	}

	buf := make([]byte, bytesOfNumTags+len(ids)*bytesOfEntry+byteOfOffset+byteOfExtended)
	bufOffset := 0
	bufOffsetOffsetNext := bytesOfNumTags + len(ids)*bytesOfEntry
	bufOffsetExtended := bufOffsetOffsetNext + byteOfOffset

	byteOrder := im.Header.ByteOrder
	if im.Header.Version == VersionClassicTIFF {
		byteOrder.PutUint16(buf[:bytesOfNumTags], uint16(len(ids)))
		byteOrder.PutUint32(buf[bufOffsetOffsetNext:bufOffsetOffsetNext+byteOfOffset], uint32(im.OffsetNext))
	} else {
		byteOrder.PutUint64(buf[:bytesOfNumTags], uint64(len(ids)))
		byteOrder.PutUint64(buf[bufOffsetOffsetNext:bufOffsetOffsetNext+byteOfOffset], uint64(im.OffsetNext))
	}
	bufOffset += bytesOfNumTags

	for _, id := range ids {
		tag := im.Tag[id]
		if im.Header.Version == VersionClassicTIFF {
			byteOrder.PutUint16(buf[bufOffset:bufOffset+2], uint16(tag.ID))
			byteOrder.PutUint16(buf[bufOffset+2:bufOffset+4], uint16(tag.Type))
			byteOrder.PutUint32(buf[bufOffset+4:bufOffset+8], uint32(tag.Count))
			if len(tag.Data) <= byteOfOffset {
				copy(buf[bufOffset+8:bufOffset+12], tag.Data)
			} else {
				pointer := offset + int64(bufOffsetExtended)
				byteOrder.PutUint32(buf[bufOffset+8:bufOffset+12], uint32(pointer))
				copy(buf[bufOffsetExtended:], tag.Data)
				bufOffsetExtended += len(tag.Data)
			}
		} else {
			byteOrder.PutUint16(buf[bufOffset:bufOffset+2], uint16(tag.ID))
			byteOrder.PutUint16(buf[bufOffset+2:bufOffset+4], uint16(tag.Type))
			byteOrder.PutUint64(buf[bufOffset+4:bufOffset+12], uint64(tag.Count))
			if len(tag.Data) <= byteOfOffset {
				copy(buf[bufOffset+12:bufOffset+20], tag.Data)
			} else {
				pointer := offset + int64(bufOffsetExtended)
				byteOrder.PutUint64(buf[bufOffset+12:bufOffset+20], uint64(pointer))
				copy(buf[bufOffsetExtended:], tag.Data)
				bufOffsetExtended += len(tag.Data)
			}
		}
		bufOffset += bytesOfEntry
	}
	return buf, nil
}

type ExifIFD struct {
	Tag map[ExifTagID]*ExifTag

	// Filled by Encoder or Decoder
	Header     *Header // Header of the TIFF file
	Offset     int64   // Offset of this IFD
	OffsetNext int64   // Offset of the next IFD
}

func (dir *ExifIFD) TagID() []ExifTagID {
	list := make([]ExifTagID, 0, len(dir.Tag))
	for id := range dir.Tag {
		list = append(list, id)
	}
	sort.Sort(byExifTagID(list))
	return list
}

type GPSIFD struct {
	Tag map[GPSTagID]*GPSTag

	// Filled by Encoder or Decoder
	Header     *Header // Header of the TIFF file
	Offset     int64   // Offset of this IFD
	OffsetNext int64   // Offset of the next IFD
}

func (dir *GPSIFD) TagID() []GPSTagID {
	list := make([]GPSTagID, 0, len(dir.Tag))
	for id := range dir.Tag {
		list = append(list, id)
	}
	sort.Sort(byGPSTagID(list))
	return list
}

type InteroperabilityIFD struct {
	Tag map[InteroperabilityTagID]*InteroperabilityTag

	// Filled by Encoder or Decoder
	Header     *Header // Header of the TIFF file
	Offset     int64   // Offset of this IFD
	OffsetNext int64   // Offset of the next IFD
}

func (dir *InteroperabilityIFD) TagID() []InteroperabilityTagID {
	list := make([]InteroperabilityTagID, 0, len(dir.Tag))
	for id := range dir.Tag {
		list = append(list, id)
	}
	sort.Sort(byInteroperabilityTagID(list))
	return list
}
