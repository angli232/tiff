package tiff

import (
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"unsafe"
)

type Encoder struct {
	Header *Header

	w      io.WriteSeeker
	offset int64 // Current offset
}

func NewEncoder(w io.WriteSeeker) *Encoder {
	return &Encoder{
		Header: &Header{
			ByteOrder: binary.LittleEndian,
			Version:   VersionClassicTIFF,
		},
		w: w,
	}
}

func (enc *Encoder) SetByteOrder(byteOrder binary.ByteOrder) {
	enc.Header.ByteOrder = byteOrder
}

func (enc *Encoder) SetVersion(version Version) {
	enc.Header.Version = version
}

func (enc *Encoder) Close() error {
	return nil
}

func (enc *Encoder) NewImage() *Image {
	return &Image{
		Header: enc.Header,
		enc:    enc,
	}
}

func (im *Image) EncodeImage(buffer interface{}) error {
	//
	// Find width and height
	//
	width, height := im.WidthHeight()
	if width == 0 {
		return fmt.Errorf("invalid image width")
	}
	if height == 0 {
		return fmt.Errorf("invalid image height")
	}
	pixelCount := width * height

	//
	// Find bitDepth and samplePerPixel
	//
	samplePerPixel := im.SamplesPerPixel()
	bitDepth := im.BitDepth()
	switch bitDepth {
	case 1:
		return UnsupportedError("1-bit image is not yet supported")
	case 8:
		if buf, ok := buffer.([]uint8); ok {
			if len(buf) != pixelCount*samplePerPixel {
				return fmt.Errorf("wrong buffer size")
			}
		} else {
			return fmt.Errorf("expecting []uint8")
		}
	case 16:
		if buf, ok := buffer.([]uint16); ok {
			if len(buf) != pixelCount*samplePerPixel {
				return fmt.Errorf("wrong buffer size")
			}
		} else {
			return fmt.Errorf("expecting []uint16")
		}
	default:
		bitsPerSample := im.BitsPerSample()
		return UnsupportedError(fmt.Sprintf("BitsPerSample of %v", bitsPerSample))
	}

	//
	// PlanarConfig
	//
	planarConfig := im.PlanarConfig()
	if planarConfig != PlanarConfigContig {
		return UnsupportedError(fmt.Sprintf("PlanarConfiguration of %v", planarConfig))
	}

	// Compression
	if im.Tag[TagCompression] == nil {
		im.SetTag(TagCompression, TagTypeShort, CompressionNone)
	}

	//
	// Strip mode or tile mode
	//
	numStrips := 0
	numTiles := 0
	nx := 0
	ny := 0
	var tileWidth, tileHeight, rowsPerStrip int
	if im.Tag[TagTileWidth] != nil || im.Tag[TagTileLength] != nil {
		//
		// Tile mode
		//
		if im.Tag[TagTileWidth] != nil && im.Tag[TagTileLength] == nil {
			return fmt.Errorf("missing TileLength tag")
		}
		if im.Tag[TagTileWidth] == nil && im.Tag[TagTileLength] != nil {
			return fmt.Errorf("missing TileWidth tag")
		}
		uintTileWidth, ok := im.Tag[TagTileWidth].Uint()
		if !ok || uintTileWidth == 0 {
			return fmt.Errorf("invalid TileWidth tag")
		}
		uintTileHeight, ok := im.Tag[TagTileLength].Uint()
		if !ok || uintTileHeight == 0 {
			return fmt.Errorf("invalid TileLength tag")
		}
		tileWidth = int(uintTileWidth)
		tileHeight = int(uintTileHeight)
		nx = (int)(math.Ceil((float64)(width) / (float64)(tileWidth)))
		ny = (int)(math.Ceil((float64)(height) / (float64)(tileHeight)))
		numTiles = nx * ny
		im.SetTag(TagStripOffsets, TagTypeLong, make([]uint32, numTiles))
		im.SetTag(TagStripByteCounts, TagTypeLong, make([]uint32, numTiles))
	} else {
		//
		// Strip mode
		//
		if im.Tag[TagRowsPerStrip] != nil {
			uintRowsPerStrip, ok := im.Tag[TagRowsPerStrip].Uint()
			if !ok || uintRowsPerStrip == 0 {
				return fmt.Errorf("invalid RowsPerStrip tag")
			}
			rowsPerStrip = int(uintRowsPerStrip)
			numStrips = (int)(math.Ceil((float64)(height) / (float64)(rowsPerStrip)))
			ny = numStrips
			nx = 1
		} else {
			numStrips = 1
			ny = 1
			nx = 1
			rowsPerStrip = height
		}
		im.SetTag(TagStripOffsets, TagTypeLong, make([]uint32, numStrips))
		im.SetTag(TagStripByteCounts, TagTypeLong, make([]uint32, numStrips))
	}

	//
	// Write Header
	//
	if im.enc.offset == 0 {
		if im.Header.Version == VersionClassicTIFF {
			im.Header.OffsetFirstIFD = 8
		} else {
			im.Header.OffsetFirstIFD = 16
		}

		buf := im.Header.encodeBytes()
		_, err := im.enc.w.Write(buf)
		if err != nil {
			return err
		}
		im.enc.offset += int64(len(buf))
	}

	//
	// Write IFD tags
	//
	buf, err := im.EncodeTags(im.enc.offset)
	if err != nil {
		return err
	}
	_, err = im.enc.w.Write(buf)
	if err != nil {
		return err
	}
	im.enc.offset += int64(len(buf))

	//
	// Strip mode
	//
	if numStrips > 0 {
		stripOffsets := make([]uint32, numStrips)
		stripByteCounts := make([]uint32, numStrips)

		samplesPerStrip := rowsPerStrip * width * samplePerPixel
		for iy := 0; iy < ny; iy++ {
			offsetBuf := iy * samplesPerStrip

			// Last strip may not be a full strip, and is not padded
			if iy == ny-1 {
				samplesPerStrip = (height - iy*rowsPerStrip) * width * samplePerPixel
			}

			switch buf := buffer.(type) {
			case []uint8:
				offset, byteCount, err := im.encodeSegment(buf[offsetBuf : offsetBuf+samplesPerStrip])
				if err != nil {
					return err
				}
				stripOffsets[iy] = uint32(offset)
				stripByteCounts[iy] = uint32(byteCount)
			case []uint16:
				if im.enc.Header.ByteOrder == binary.LittleEndian {
					bufByte := (*(*[1 << 48]uint8)(unsafe.Pointer(&buf[0])))[offsetBuf*2 : (offsetBuf+samplesPerStrip)*2]
					offset, byteCount, err := im.encodeSegment(bufByte)
					if err != nil {
						return err
					}
					stripOffsets[iy] = uint32(offset)
					stripByteCounts[iy] = uint32(byteCount)
				} else {
					return UnsupportedError("big endian is not yet supported")
				}
			}
		}

		im.SetTag(TagStripOffsets, TagTypeLong, stripOffsets)
		im.SetTag(TagStripByteCounts, TagTypeLong, stripByteCounts)

		// Current workaround is to overwrite all tags
		buf, err := im.EncodeTags(im.Offset)
		if err != nil {
			return err
		}
		im.enc.w.Seek(im.Offset, io.SeekStart)
		_, err = im.enc.w.Write(buf)
		if err != nil {
			return err
		}
		im.enc.w.Seek(im.enc.offset, io.SeekStart)
	}

	//
	// Tile mode
	//
	if numTiles > 0 {
		return UnsupportedError("tile mode is not yet supported")
	}

	return nil
}

// encodeSegment handles compression of segments. This is a temporary solution, will be replaced with a Writer.
func (im *Image) encodeSegment(buf []byte) (offset int64, byteCount int, err error) {
	offset = im.enc.offset

	compression := im.Compression()
	switch compression {
	case CompressionNone:
		_, err = im.enc.w.Write(buf)
		if err != nil {
			byteCount = 0
			return
		}
		im.enc.offset += int64(len(buf))
		byteCount = len(buf)
		return
	case CompressionDeflate:
		w := zlib.NewWriter(im.enc.w)
		_, err := w.Write(buf)
		if err != nil {
			w.Close()
			return 0, 0, err
		}
		w.Close()
		currentOffset, err := im.enc.w.Seek(0, io.SeekCurrent)
		if err != nil {
			return 0, 0, err
		}
		byteCount = int(currentOffset - offset)
		im.enc.offset = currentOffset
		return offset, byteCount, nil

	default:
		return 0, 0, UnsupportedError(fmt.Sprintf("compression %d is not supported", compression))
	}
}
