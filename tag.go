package tiff

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

// TagType indicates data type of the field.
type TagType uint16

// Field data type
const (
	TagTypeByte      TagType = 1  // uint8
	TagTypeASCII     TagType = 2  // 8-bit byte that contains a 7-bit ASCII code; the last byte must be NUL (binary zero).
	TagTypeShort     TagType = 3  // uint16
	TagTypeLong      TagType = 4  // uint32
	TagTypeRational  TagType = 5  // Two Long's: the first represents the numerator of a fraction; the second, the denominator.
	TagTypeSByte     TagType = 6  // int8
	TagTypeUndefined TagType = 7  // An 8-bit byte that may contain anything, depending on the definition of the field.
	TagTypeSShort    TagType = 8  // int16
	TagTypeSLong     TagType = 9  // int32
	TagTypeSRational TagType = 10 // Two SLong's: the first represents the numerator of a fraction, the second the denominator.
	TagTypeFloat     TagType = 11 // float32
	TagTypeDouble    TagType = 12 // float64
	TagTypeIFD       TagType = 13 // uint32 (bigTIFF)
	TagTypeLong8     TagType = 16 // uint64 (bigTIFF)
	TagTypeSLong8    TagType = 17 // int64 (bigTIFF)
	TagTypeIFD8      TagType = 18 // uint16 IFD offset (bigTIFF)
)

var typeSize = map[TagType]int{
	TagTypeByte:      1,
	TagTypeASCII:     1,
	TagTypeShort:     2,
	TagTypeLong:      4,
	TagTypeRational:  8,
	TagTypeSByte:     1,
	TagTypeUndefined: 1,
	TagTypeSShort:    2,
	TagTypeSLong:     4,
	TagTypeSRational: 8,
	TagTypeFloat:     4,
	TagTypeDouble:    8,
	TagTypeIFD:       4,
	TagTypeLong8:     8,
	TagTypeSLong8:    8,
	TagTypeIFD8:      8,
}

var typeName = map[TagType]string{
	TagTypeByte:      "Byte",
	TagTypeASCII:     "ASCII",
	TagTypeShort:     "Short",
	TagTypeLong:      "Long",
	TagTypeRational:  "Rational",
	TagTypeSByte:     "SByte",
	TagTypeUndefined: "Undefined",
	TagTypeSShort:    "SShort",
	TagTypeSLong:     "SLong",
	TagTypeSRational: "SRational",
	TagTypeFloat:     "Float",
	TagTypeDouble:    "Double",
	TagTypeIFD:       "IFD",
	TagTypeLong8:     "Long8",
	TagTypeSLong8:    "SLong8",
	TagTypeIFD8:      "IFD8",
}

func (t TagType) Size() int {
	return typeSize[t]
}

func (t TagType) String() string {
	return typeName[t]
}

type Tag struct {
	ID    TagID   // Tag identifying code
	Type  TagType // Data type of tag data
	Count int     // Number of values
	Data  []byte  // Tag data or offset to tag data

	// Filled by Encoder or Decoder
	Header      *Header // Header of the TIFF file
	OffsetEntry int64   // Offset of the tag in the TIFF file
	OffsetData  int64   // Offset of the data if the data are not in the IFD Entry
}

func NewTag(id TagID, tagType TagType, value interface{}, header *Header) (tag *Tag, err error) {
	tag = &Tag{
		ID:     id,
		Type:   tagType,
		Header: header,
	}
	switch vSlice := value.(type) {
	case string:
		return NewTag(id, tagType, []string{vSlice}, header)
	case uint16:
		return NewTag(id, tagType, []uint16{vSlice}, header)
	case uint32:
		return NewTag(id, tagType, []uint32{vSlice}, header)
	case uint64:
		return NewTag(id, tagType, []uint64{vSlice}, header)
	case int16:
		return NewTag(id, tagType, []int16{vSlice}, header)
	case int32:
		return NewTag(id, tagType, []int32{vSlice}, header)
	case int64:
		return NewTag(id, tagType, []int64{vSlice}, header)
	case float32:
		return NewTag(id, tagType, []float32{vSlice}, header)
	case float64:
		return NewTag(id, tagType, []float64{vSlice}, header)
	case Rational:
		return NewTag(id, tagType, []Rational{vSlice}, header)
	case SRational:
		return NewTag(id, tagType, []SRational{vSlice}, header)
	case int:
		return NewTag(id, tagType, []int{vSlice}, header)
	case []int:
		switch tagType {
		case TagTypeShort:
			tmp := make([]uint16, len(vSlice))
			for i, d := range vSlice {
				tmp[i] = uint16(d)
			}
			return NewTag(id, tagType, tmp, header)
		case TagTypeLong:
			tmp := make([]uint32, len(vSlice))
			for i, d := range vSlice {
				tmp[i] = uint32(d)
			}
			return NewTag(id, tagType, tmp, header)
		default:
			return nil, fmt.Errorf("type mismatch: set int to %s", tagType)
		}
	case []string:
		if tagType != TagTypeASCII {
			return nil, fmt.Errorf("type mismatch: set string to %s", tagType)
		}
		tag.Data = make([]byte, 0)
		for _, v := range vSlice {
			tag.Data = append(tag.Data, []byte(v+"\x00")...)
		}
		tag.Count = len(tag.Data)
	case []byte:
		if tagType != TagTypeByte && tagType != TagTypeUndefined {
			return nil, fmt.Errorf("type mismatch: set byte to %s", tagType)
		}
		tag.Count = len(vSlice)
		tag.Data = vSlice
	case []uint16:
		if tagType != TagTypeShort {
			return nil, fmt.Errorf("type mismatch: set uint16 to %s", tagType)
		}
		tag.Count = len(vSlice)
		tag.Data = make([]byte, tag.Count*2)
		for i := 0; i < tag.Count; i++ {
			header.ByteOrder.PutUint16(tag.Data[2*i:2*(i+1)], vSlice[i])
		}
	case []uint32:
		if tagType != TagTypeLong && tagType != TagTypeIFD {
			return nil, fmt.Errorf("type mismatch: set uint32 to %s", tagType)
		}
		tag.Count = len(vSlice)
		tag.Data = make([]byte, tag.Count*4)
		for i := 0; i < tag.Count; i++ {
			header.ByteOrder.PutUint32(tag.Data[4*i:4*(i+1)], vSlice[i])
		}
	case []uint64:
		if tagType != TagTypeLong8 && tagType != TagTypeIFD8 {
			return nil, fmt.Errorf("type mismatch: set uint64 to %s", tagType)
		}
		tag.Count = len(vSlice)
		tag.Data = make([]byte, tag.Count*8)
		for i := 0; i < tag.Count; i++ {
			header.ByteOrder.PutUint64(tag.Data[8*i:8*(i+1)], vSlice[i])
		}
	case []int16:
		if tagType != TagTypeSShort {
			return nil, fmt.Errorf("type mismatch: set int16 to %s", tagType)
		}
		tag.Count = len(vSlice)
		tag.Data = make([]byte, tag.Count*2)
		for i := 0; i < tag.Count; i++ {
			header.ByteOrder.PutUint16(tag.Data[2*i:2*(i+1)], uint16(vSlice[i]))
		}
	case []int32:
		if tagType != TagTypeSLong {
			return nil, fmt.Errorf("type mismatch: set int32 to %s", tagType)
		}
		tag.Count = len(vSlice)
		tag.Data = make([]byte, tag.Count*4)
		for i := 0; i < tag.Count; i++ {
			header.ByteOrder.PutUint32(tag.Data[4*i:4*(i+1)], uint32(vSlice[i]))
		}
	case []int64:
		if tagType != TagTypeSLong8 {
			return nil, fmt.Errorf("type mismatch: set int64 to %s", tagType)
		}
		tag.Count = len(vSlice)
		tag.Data = make([]byte, tag.Count*8)
		for i := 0; i < tag.Count; i++ {
			header.ByteOrder.PutUint64(tag.Data[8*i:8*(i+1)], uint64(vSlice[i]))
		}
	case []float32:
		if tagType != TagTypeFloat {
			return nil, fmt.Errorf("type mismatch: set float32 to %s", tagType)
		}
		tag.Count = len(vSlice)
		buf := new(bytes.Buffer)
		err = binary.Write(buf, header.ByteOrder, vSlice)
		if err != nil {
			return nil, err
		}
		tag.Data = buf.Bytes()
	case []float64:
		if tagType != TagTypeDouble {
			return nil, fmt.Errorf("type mismatch: set float64 to %s", tagType)
		}
		tag.Count = len(vSlice)
		buf := new(bytes.Buffer)
		err = binary.Write(buf, header.ByteOrder, vSlice)
		if err != nil {
			return nil, err
		}
		tag.Data = buf.Bytes()
	case []Rational:
		if tagType != TagTypeRational {
			return nil, fmt.Errorf("type mismatch: set Rational to %s", tagType)
		}
		tag.Count = len(vSlice)
		tag.Data = make([]byte, tag.Count*8)
		for i := 0; i < tag.Count; i++ {
			header.ByteOrder.PutUint32(tag.Data[8*i:8*(i+1)], vSlice[i][0])
			header.ByteOrder.PutUint32(tag.Data[8*(i+1):4*(i+2)], vSlice[i][1])
		}
	case []SRational:
		if tagType != TagTypeSRational {
			return nil, fmt.Errorf("type mismatch: set SRational to %s", tagType)
		}
		tag.Count = len(vSlice)
		tag.Data = make([]byte, tag.Count*8)
		for i := 0; i < tag.Count; i++ {
			header.ByteOrder.PutUint32(tag.Data[8*i:8*(i+1)], uint32(vSlice[i][0]))
			header.ByteOrder.PutUint32(tag.Data[8*(i+1):4*(i+2)], uint32(vSlice[i][1]))
		}
	default:
		return nil, fmt.Errorf("invalid value type")
	}
	return tag, nil
}

func (t *Tag) GoString() string {
	return formatTagValue(t.Type, t.Count, t.Data, t.Header.ByteOrder)
}

func (t *Tag) String() (value string, ok bool) {
	if t == nil {
		return "", false
	}
	if t.Type != TagTypeASCII {
		return "", false
	}
	v, err := DecodeASCII(t.Count, t.Data, t.Header.ByteOrder)
	if err != nil {
		return "", false
	}
	return v[0], true
}

func (t *Tag) StringSlice() (value []string, ok bool) {
	if t == nil {
		return nil, false
	}
	if t.Type != TagTypeASCII {
		return nil, false
	}
	v, err := DecodeASCII(t.Count, t.Data, t.Header.ByteOrder)
	if err != nil {
		return nil, false
	}
	return v, true
}

// Uint decodes Short, Long or Long8 tag value as uint, when the tag only has one value.
//
// It returns (0, false) when the tags does not exist, has zero count or multiple values, or is of other data type.
func (t *Tag) Uint() (value uint, ok bool) {
	if t == nil || t.Count != 1 {
		return 0, false
	}
	switch t.Type {
	case TagTypeShort:
		v := DecodeShort(t.Count, t.Data, t.Header.ByteOrder)
		return uint(v[0]), true
	case TagTypeLong, TagTypeIFD:
		v := DecodeLong(t.Count, t.Data, t.Header.ByteOrder)
		return uint(v[0]), true
	case TagTypeLong8, TagTypeIFD8:
		v := DecodeLong8(t.Count, t.Data, t.Header.ByteOrder)
		return uint(v[0]), true
	default:
		return 0, false
	}
}

// UintSlice decodes Short, Long or Long8 tag values as []uint.
//
// It returns ([]uint{}, false) when the tags does not exist, has zero-count, or is of other data type.
func (t *Tag) UintSlice() (value []uint, ok bool) {
	if t == nil || t.Count == 0 {
		return []uint{}, false
	}
	switch t.Type {
	case TagTypeShort:
		v := DecodeShort(t.Count, t.Data, t.Header.ByteOrder)
		value = make([]uint, len(v))
		for i := 0; i < len(v); i++ {
			value[i] = uint(v[i])
		}
		return value, true
	case TagTypeLong, TagTypeIFD:
		v := DecodeLong(t.Count, t.Data, t.Header.ByteOrder)
		value = make([]uint, len(v))
		for i := 0; i < len(v); i++ {
			value[i] = uint(v[i])
		}
		return value, true
	case TagTypeLong8, TagTypeIFD8:
		v := DecodeLong8(t.Count, t.Data, t.Header.ByteOrder)
		value = make([]uint, len(v))
		for i := 0; i < len(v); i++ {
			value[i] = uint(v[i])
		}
		return value, true
	default:
		return []uint{}, false
	}
}

type ExifTag struct {
	ID    ExifTagID // Tag identifying code
	Type  TagType   // Data type of tag data
	Count int       // Number of values
	Data  []byte    // Tag data or offset to tag data

	// Filled by Encoder or Decoder
	Header      *Header // Header of the TIFF file
	OffsetEntry int64   // Offset of the tag in the TIFF file
	OffsetData  int64   // Offset of the data if the data are not in the IFD Entry
}

func (t *ExifTag) GoString() string {
	return formatTagValue(t.Type, t.Count, t.Data, t.Header.ByteOrder)
}

type GPSTag struct {
	ID    GPSTagID // Tag identifying code
	Type  TagType  // Data type of tag data
	Count int      // Number of values
	Data  []byte   // Tag data or offset to tag data

	// Filled by Encoder or Decoder
	Header      *Header // Header of the TIFF file
	OffsetEntry int64   // Offset of the tag in the TIFF file
	OffsetData  int64   // Offset of the data if the data are not in the IFD Entry
}

type InteroperabilityTag struct {
	ID    InteroperabilityTagID // Tag identifying code
	Type  TagType               // Data type of tag data
	Count int                   // Number of values
	Data  []byte                // Tag data or offset to tag data

	// Filled by Encoder or Decoder
	Header      *Header // Header of the TIFF file
	OffsetEntry int64   // Offset of the tag in the TIFF file
	OffsetData  int64   // Offset of the data if the data are not in the IFD Entry
}

func DecodeASCII(count int, data []byte, byteOrder binary.ByteOrder) ([]string, error) {
	if len(data) != count {
		panic("decode tag: unexpected size of data")
	}
	if data[len(data)-1] != 0 {
		return nil, fmt.Errorf("decode ASCII: last byte is not 0")
	}
	return strings.Split(string(data[:len(data)-1]), "\x00"), nil
}

func DecodeSByte(count int, data []byte, byteOrder binary.ByteOrder) []int8 {
	if len(data) != count {
		panic("decode tag: unexpected size of data")
	}
	v := make([]int8, count)
	for i := 0; i < count; i++ {
		v[i] = int8(data[i])
	}
	return v
}

func DecodeShort(count int, data []byte, byteOrder binary.ByteOrder) []uint16 {
	sizeElem := TagTypeShort.Size()
	if len(data) != count*sizeElem {
		panic("decode tag: unexpected size of data")
	}
	v := make([]uint16, count)
	offset := 0
	for i := 0; i < count; i++ {
		v[i] = byteOrder.Uint16(data[offset : offset+sizeElem])
		offset += sizeElem
	}
	return v
}

func DecodeSShort(count int, data []byte, byteOrder binary.ByteOrder) []int16 {
	sizeElem := TagTypeSShort.Size()
	if len(data) != count*sizeElem {
		panic("decode tag: unexpected size of data")
	}
	v := make([]int16, count)
	offset := 0
	for i := 0; i < count; i++ {
		v[i] = int16(byteOrder.Uint16(data[offset : offset+sizeElem]))
		offset += sizeElem
	}
	return v
}

func DecodeLong(count int, data []byte, byteOrder binary.ByteOrder) []uint32 {
	sizeElem := TagTypeLong.Size()
	if len(data) != count*sizeElem {
		panic("decode tag: unexpected size of data")
	}
	v := make([]uint32, count)
	offset := 0
	for i := 0; i < count; i++ {
		v[i] = byteOrder.Uint32(data[offset : offset+sizeElem])
		offset += sizeElem
	}
	return v
}

func DecodeSLong(count int, data []byte, byteOrder binary.ByteOrder) []int32 {
	sizeElem := TagTypeSLong.Size()
	if len(data) != count*sizeElem {
		panic("decode tag: unexpected size of data")
	}
	v := make([]int32, count)
	offset := 0
	for i := 0; i < count; i++ {
		v[i] = int32(byteOrder.Uint32(data[offset : offset+sizeElem]))
		offset += sizeElem
	}
	return v
}

func DecodeRational(count int, data []byte, byteOrder binary.ByteOrder) []Rational {
	sizeElem := TagTypeRational.Size()
	if len(data) != count*sizeElem {
		panic("decode tag: unexpected size of data")
	}
	v := make([]Rational, count)
	offset := 0
	for i := 0; i < count; i++ {
		v[i] = [2]uint32{
			byteOrder.Uint32(data[offset : offset+sizeElem/2]),
			byteOrder.Uint32(data[offset+sizeElem/2 : offset+sizeElem]),
		}
		offset += sizeElem
	}
	return v
}

func DecodeSRational(count int, data []byte, byteOrder binary.ByteOrder) []SRational {
	sizeElem := TagTypeSRational.Size()
	if len(data) != count*sizeElem {
		panic("decode tag: unexpected size of data")
	}
	v := make([]SRational, count)
	offset := 0
	for i := 0; i < count; i++ {
		v[i] = [2]int32{
			int32(byteOrder.Uint32(data[offset : offset+sizeElem/2])),
			int32(byteOrder.Uint32(data[offset+sizeElem/2 : offset+sizeElem])),
		}
		offset += sizeElem
	}
	return v
}

func DecodeFloat(count int, data []byte, byteOrder binary.ByteOrder) []float32 {
	sizeElem := TagTypeFloat.Size()
	if len(data) != count*sizeElem {
		panic("decode tag: unexpected size of data")
	}
	v := make([]float32, count)
	_ = binary.Read(bytes.NewReader(data), byteOrder, &v)
	return v
}

func DecodeDouble(count int, data []byte, byteOrder binary.ByteOrder) []float64 {
	sizeElem := TagTypeDouble.Size()
	if len(data) != count*sizeElem {
		panic("decode tag: unexpected size of data")
	}
	v := make([]float64, count)
	_ = binary.Read(bytes.NewReader(data), byteOrder, &v)
	return v
}

func DecodeLong8(count int, data []byte, byteOrder binary.ByteOrder) []uint64 {
	sizeElem := TagTypeLong8.Size()
	if len(data) != count*sizeElem {
		panic("decode tag: unexpected size of data")
	}
	v := make([]uint64, count)
	offset := 0
	for i := 0; i < count; i++ {
		v[i] = byteOrder.Uint64(data[offset : offset+sizeElem])
		offset += sizeElem
	}
	return v
}

func DecodeSLong8(count int, data []byte, byteOrder binary.ByteOrder) []int64 {
	sizeElem := TagTypeSLong8.Size()
	if len(data) != count*sizeElem {
		panic("decode tag: unexpected size of data")
	}
	v := make([]int64, count)
	offset := 0
	for i := 0; i < count; i++ {
		v[i] = int64(byteOrder.Uint64(data[offset : offset+sizeElem]))
		offset += sizeElem
	}
	return v
}

// formatTagValue formats value of a tag to a string.
func formatTagValue(tagType TagType, count int, data []byte, byteOrder binary.ByteOrder) string {
	switch tagType {
	case TagTypeByte, TagTypeUndefined:
		return fmt.Sprintf("%q", string(data))
	case TagTypeASCII:
		v, err := DecodeASCII(count, data, byteOrder)
		if err != nil {
			return fmt.Sprintf("%#v", err)
		}
		s := fmt.Sprintf("%q", v)
		return s[1 : len(s)-1]
	case TagTypeShort:
		v := DecodeShort(count, data, byteOrder)
		s := fmt.Sprintf("%v", v)
		return s[1 : len(s)-1]
	case TagTypeSShort:
		v := DecodeSShort(count, data, byteOrder)
		s := fmt.Sprintf("%v", v)
		return s[1 : len(s)-1]
	case TagTypeLong, TagTypeIFD:
		v := DecodeLong(count, data, byteOrder)
		s := fmt.Sprintf("%v", v)
		return s[1 : len(s)-1]
	case TagTypeSLong:
		v := DecodeSLong(count, data, byteOrder)
		s := fmt.Sprintf("%v", v)
		return s[1 : len(s)-1]
	case TagTypeRational:
		v := DecodeRational(count, data, byteOrder)
		s := make([]string, 0, len(v))
		for _, value := range v {
			s = append(s, value.String())
		}
		return strings.Join(s, " ")
	case TagTypeSRational:
		v := DecodeSRational(count, data, byteOrder)
		s := make([]string, 0, len(v))
		for _, value := range v {
			s = append(s, value.String())
		}
		return strings.Join(s, " ")
	case TagTypeFloat:
		v := DecodeFloat(count, data, byteOrder)
		s := fmt.Sprintf("%g", v)
		return s[1 : len(s)-1]
	case TagTypeDouble:
		v := DecodeDouble(count, data, byteOrder)
		s := fmt.Sprintf("%g", v)
		return s[1 : len(s)-1]
	case TagTypeLong8, TagTypeIFD8:
		v := DecodeLong8(count, data, byteOrder)
		s := fmt.Sprintf("%v", v)
		return s[1 : len(s)-1]
	case TagTypeSLong8:
		v := DecodeSLong8(count, data, byteOrder)
		s := fmt.Sprintf("%v", v)
		return s[1 : len(s)-1]
	default:
		return "%!s(unknown type)"
	}
}
