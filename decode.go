package tiff

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"unsafe"

	// Compression formats
	"compress/zlib" // Deflate
	"image"         // JPEG
	"image/jpeg"    // JPEG

	"golang.org/x/image/tiff/lzw" // LZW
)

// Iter is a iterator of Images in a TIFF file.
type Iter struct {
	index int
	dec   *Decoder
	im    *Image
	err   error
}

// All returns metadata of all remaining Images from the iterator.
func (it *Iter) All() ([]*Image, error) {
	ims := make([]*Image, 0)
	for it.Next() {
		im := it.Image()
		ims = append(ims, im)
	}
	err := it.Err()
	return ims, err
}

// Next advances the iterator to the next Image.
func (it *Iter) Next() bool {
	if it == nil {
		return false
	}
	if it.err != nil {
		// will not advance when there is an error
		return false
	}
	if it.im == nil {
		// calling Next() after the last image results in a EOF error
		it.err = io.EOF
		return false
	}
	if it.im.OffsetNext == 0 {
		// terminate the iteration at the last image
		it.im = nil
		return false
	}
	im, err := it.dec.decodeDirectories(it.im.OffsetNext)
	if err != nil {
		it.err = err
		return false
	}
	im.rd = it.im.rd         // Get the reader from the current Image
	im.Header = it.im.Header // Get Header from the current Image
	it.im = im
	it.index++
	return true
}

// Image returns the current image in the iteration. It should not be called before Next().
func (it *Iter) Image() *Image {
	if it == nil {
		return nil
	}
	if it.im == nil {
		return nil
	}
	if it.im.Offset == 0 {
		it.err = fmt.Errorf("acesss image before first call of Next()")
		return nil
	}
	return it.im
}

// Index returns the index of the current image in the file.
func (it *Iter) Index() int {
	return it.index
}

// Err returns the error that terminited the iteration.
func (it *Iter) Err() error {
	if it == nil {
		return fmt.Errorf("nil iterator")
	}
	return it.err
}

// Decoder of TIFF format.
type Decoder struct {
	Header *Header

	r io.ReaderAt
}

// NewDecoder returns a TIFF decoder.
func NewDecoder(r io.ReaderAt) (*Decoder, error) {
	header, err := decodeHeader(r)
	if err != nil {
		return nil, err
	}
	return &Decoder{
		Header: header,
		r:      r,
	}, nil
}

// Iter returns an Image iterator.
func (d *Decoder) Iter() *Iter {
	return &Iter{
		index: -1,
		dec:   d,
		im: &Image{
			Offset:     0,
			OffsetNext: d.Header.OffsetFirstIFD,
			Header:     d.Header,
			rd:         d.r,
		},
		err: nil,
	}
}

// decodeDirectories is a helper function to return IFD with SubIFD, ExifIFD, etc.
func (d *Decoder) decodeDirectories(offset int64) (im *Image, err error) {
	im, err = d.DecodeIFD(offset)
	if err != nil {
		return
	}
	offsetSubIFD, hasSubIFD := im.Tag[TagSubIFDs].UintSlice()
	if hasSubIFD && len(offsetSubIFD) > 0 {
		im.SubImage = make([]*Image, len(offsetSubIFD))
		for i := 0; i < len(offsetSubIFD); i++ {
			subIFD, err := d.DecodeIFD(int64(offsetSubIFD[i]))
			if err != nil {
				return im, fmt.Errorf("cannot decode SubIFD: %s", err)
			}
			im.SubImage[i] = subIFD
		}
	}
	offsetExifIFD, hasExif := im.Tag[TagExifIFD].Uint()
	if hasExif {
		im.Exif, err = d.DecodeExifIFD(int64(offsetExifIFD))
		if err != nil {
			return im, fmt.Errorf("cannot decode ExifIFD: %s", err)
		}
	}
	offsetGPSIFD, hasGPS := im.Tag[TagGPSIFD].Uint()
	if hasGPS {
		im.GPS, err = d.DecodeGPSIFD(int64(offsetGPSIFD))
		if err != nil {
			return im, fmt.Errorf("cannot decode GPSIFD: %s", err)
		}
	}
	offsetInteroperabilityIFD, hasInteroperability := im.Tag[TagInteroperabilityIFD].Uint()
	if hasInteroperability {
		im.Interoperability, err = d.DecodeInteroperabilityIFD(int64(offsetInteroperabilityIFD))
		if err != nil {
			return im, fmt.Errorf("cannot decode Interoperability: %s", err)
		}
	}
	return im, nil
}

// DecodeIFD decodes IFD at given offset and returns the IFD, without handling SubIFD, etc.
func (d *Decoder) DecodeIFD(offset int64) (im *Image, err error) {
	im = &Image{
		Offset: offset,
		Tag:    make(map[TagID]*Tag),
	}

	var bytesOfNumTags int
	var bytesOfEntry int
	var byteOfOffset int
	switch d.Header.Version {
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

	// Buffer for reading IFD entries
	buf := make([]byte, 0, bytesOfEntry)

	// Read number of entries
	_, err = d.r.ReadAt(buf[:bytesOfNumTags], offset)
	if err != nil {
		return nil, fmt.Errorf("cannot read IFD: %s", err)
	}
	var numEntries int
	if d.Header.Version == VersionClassicTIFF {
		numEntries = int(d.Header.ByteOrder.Uint16(buf[:bytesOfNumTags]))
	} else {
		numEntries = int(d.Header.ByteOrder.Uint64(buf[:bytesOfNumTags]))
	}

	// Read IFD Entries
	var lastTagID int

	for i := 0; i < numEntries; i++ {
		offsetEntry := offset + int64(bytesOfNumTags) + int64(i)*int64(bytesOfEntry)
		_, err = d.r.ReadAt(buf[:bytesOfEntry], offsetEntry)
		if err != nil {
			return nil, fmt.Errorf("cannot read IFD: %s", err)
		}

		t := &Tag{
			ID:   TagID(d.Header.ByteOrder.Uint16(buf[0:2])),
			Type: TagType(d.Header.ByteOrder.Uint16(buf[2:4])),

			Header:      d.Header,
			OffsetEntry: offsetEntry,
		}

		if d.Header.Version == VersionClassicTIFF {
			t.Count = int(d.Header.ByteOrder.Uint32(buf[4:8]))
		} else {
			count := d.Header.ByteOrder.Uint64(buf[4:12])
			if count > maxInt {
				return nil, fmt.Errorf("invalid IFD: count overflows int")
			}
			t.Count = int(count)
		}

		bytesOfData := t.Type.Size() * t.Count
		t.Data = make([]byte, bytesOfData)
		if bytesOfData > byteOfOffset {
			if d.Header.Version == VersionClassicTIFF {
				t.OffsetData = int64(d.Header.ByteOrder.Uint32(buf[8:12]))
			} else {
				offsetData := d.Header.ByteOrder.Uint64(buf[12:20])
				if offsetData > math.MaxInt64 {
					return nil, fmt.Errorf("invalid IFD: offset to data overflows int64")
				}
				t.OffsetData = int64(offsetData)
			}

			_, err := d.r.ReadAt(t.Data, t.OffsetData)
			if err != nil {
				return nil, fmt.Errorf("cannot read IFD: %s", err)
			}
		} else {
			if d.Header.Version == VersionClassicTIFF {
				copy(t.Data, buf[8:12])
			} else {
				copy(t.Data, buf[12:20])
			}
		}

		if int(t.ID) <= lastTagID {
			return nil, fmt.Errorf("invalid IFD: entries are not sorted by tag")
		}

		im.Tag[t.ID] = t
	}

	// Read offset to next IFD
	_, err = d.r.ReadAt(buf[:byteOfOffset], offset+int64(bytesOfNumTags)+int64(numEntries)*int64(bytesOfEntry))
	if err != nil {
		return nil, fmt.Errorf("cannot read IFD: %s", err)
	}
	if d.Header.Version == VersionClassicTIFF {
		im.OffsetNext = int64(d.Header.ByteOrder.Uint32(buf[:byteOfOffset]))
	} else {
		offsetNext := d.Header.ByteOrder.Uint64(buf[:byteOfOffset])
		if offsetNext > math.MaxInt64 {
			return nil, fmt.Errorf("invalid IFD: offset to next IFD overflows int64")
		}
		im.OffsetNext = int64(offsetNext)
	}

	return im, nil
}

// DecodeExifIFD decodes IFD at given offset and returns as ExifIFD.
func (d *Decoder) DecodeExifIFD(offset int64) (exif *ExifIFD, err error) {
	ifd, err := d.DecodeIFD(offset)
	if err != nil {
		return nil, err
	}
	// Convert types to ExifTag and ExifIFD
	exif = &ExifIFD{
		Tag: make(map[ExifTagID]*ExifTag),

		Header:     ifd.Header,
		Offset:     ifd.Offset,
		OffsetNext: ifd.OffsetNext,
	}
	for id, tag := range ifd.Tag {
		exif.Tag[ExifTagID(id)] = &ExifTag{
			ID:    ExifTagID(tag.ID),
			Type:  tag.Type,
			Count: tag.Count,
			Data:  tag.Data,

			Header:      tag.Header,
			OffsetEntry: tag.OffsetEntry,
			OffsetData:  tag.OffsetData,
		}
	}
	return exif, nil
}

// DecodeGPSIFD decodes IFD at given offset and returns as GPSIFD.
func (d *Decoder) DecodeGPSIFD(offset int64) (gps *GPSIFD, err error) {
	ifd, err := d.DecodeIFD(offset)
	if err != nil {
		return nil, err
	}
	// Convert types to GPSTag and GPSIFD
	gps = &GPSIFD{
		Tag: make(map[GPSTagID]*GPSTag),

		Header:     ifd.Header,
		Offset:     ifd.Offset,
		OffsetNext: ifd.OffsetNext,
	}
	for id, tag := range ifd.Tag {
		gps.Tag[GPSTagID(id)] = &GPSTag{
			ID:    GPSTagID(tag.ID),
			Type:  tag.Type,
			Count: tag.Count,
			Data:  tag.Data,

			Header:      tag.Header,
			OffsetEntry: tag.OffsetEntry,
			OffsetData:  tag.OffsetData,
		}
	}
	return gps, nil
}

// DecodeInteroperabilityIFD decodes IFD at given offset and returns as InteroperabilityIFD.
func (d *Decoder) DecodeInteroperabilityIFD(offset int64) (interoperability *InteroperabilityIFD, err error) {
	ifd, err := d.DecodeIFD(offset)
	if err != nil {
		return nil, err
	}
	// Convert types to InteroperabilityTagID and InteroperabilityIFD
	interoperability = &InteroperabilityIFD{
		Tag: make(map[InteroperabilityTagID]*InteroperabilityTag),

		Header:     ifd.Header,
		Offset:     ifd.Offset,
		OffsetNext: ifd.OffsetNext,
	}
	for id, tag := range ifd.Tag {
		interoperability.Tag[InteroperabilityTagID(id)] = &InteroperabilityTag{
			ID:    InteroperabilityTagID(tag.ID),
			Type:  tag.Type,
			Count: tag.Count,
			Data:  tag.Data,

			Header:      tag.Header,
			OffsetEntry: tag.OffsetEntry,
			OffsetData:  tag.OffsetData,
		}
	}
	return interoperability, nil
}

// NumStrips returns the number of strips in the file, or 0 if the files in not organized in strips.
//
// NumStrips finds the number of strips based on TagStripOffsets and TagStripByteCounts fields in metadata.
// It does not validate whether it is consistent with the number the image is supposed to have based on the size
// from TagImageLength and TagRowsPerStrip.
func (im *Image) NumStrips() int {
	// Validate the existence of tags
	offsetTag := im.Tag[TagStripOffsets]
	byteCountTag := im.Tag[TagStripByteCounts]
	if offsetTag == nil || byteCountTag == nil {
		return 0
	}
	// Validate the size of tags
	if offsetTag.Data == nil || byteCountTag.Data == nil {
		return 0
	}
	if len(offsetTag.Data) != offsetTag.Type.Size()*offsetTag.Count {
		return 0
	}
	if len(byteCountTag.Data) != byteCountTag.Type.Size()*byteCountTag.Count {
		return 0
	}
	// Make sure the two fields have the same length
	if offsetTag.Count != byteCountTag.Count {
		return 0
	}
	return offsetTag.Count
}

// NumTiles returns the number of tiles in the file, or 0 if the files in not organized in tiles.
//
// NumTiles finds the number of tiles based on TagTileOffsets and TagTileByteCounts fields in metadata.
// It does not validate whether it is consistent with the number that the image is supposed to have based on
// from TagImageWidth, TagImageLength, TagTileWidth, TagTileLength
func (im *Image) NumTiles() int {
	// Validate the existence of tags
	offsetTag := im.Tag[TagTileOffsets]
	byteCountTag := im.Tag[TagTileByteCounts]
	if offsetTag == nil || byteCountTag == nil {
		return 0
	}
	// Validate the size of tags
	if offsetTag.Data == nil || byteCountTag.Data == nil {
		return 0
	}
	if len(offsetTag.Data) != offsetTag.Type.Size()*offsetTag.Count {
		return 0
	}
	if len(byteCountTag.Data) != byteCountTag.Type.Size()*byteCountTag.Count {
		return 0
	}
	// Make sure the two fields have the same length
	if offsetTag.Count != byteCountTag.Count {
		return 0
	}
	return offsetTag.Count
}

// RawStripReader returns a Reader for reading the raw data of a strip, which may be compressed.
//
// It panics when index >= im.NumStrips().
func (im *Image) RawStripReader(index int) *io.SectionReader {
	offsets, _ := im.Tag[TagStripOffsets].UintSlice()
	byteCounts, _ := im.Tag[TagStripByteCounts].UintSlice()
	if index >= len(offsets) || index >= len(byteCounts) {
		panic("tiff: RawStripReader: index out of range")
	}
	return io.NewSectionReader(im.rd, int64(offsets[index]), int64(byteCounts[index]))
}

// RawTileReader returns a Reader for reading the raw data of a tile, which may be compressed.
//
// It panics when index >= im.NumTiles().
func (im *Image) RawTileReader(index int) *io.SectionReader {
	offsets, _ := im.Tag[TagTileOffsets].UintSlice()
	byteCounts, _ := im.Tag[TagTileByteCounts].UintSlice()
	if index >= len(offsets) || index >= len(byteCounts) {
		panic("tiff: RawTileReader: index out of range")
	}
	return io.NewSectionReader(im.rd, int64(offsets[index]), int64(byteCounts[index]))
}

// StripReader returns a Reader for reading the decompressed data of a strip.
//
// It panics when index >= im.NumStrips().
func (im *Image) StripReader(index int) (r io.ReadCloser, err error) {
	rawReader := im.RawStripReader(index)
	return im.decodeRawSegment(rawReader)
}

// TileReader returns a Reader for reading the decompressed data of a tile.
//
// It panics when index >= im.NumTiles().
func (im *Image) TileReader(index int) (r io.ReadCloser, err error) {
	rawReader := im.RawTileReader(index)
	return im.decodeRawSegment(rawReader)
}

func (im *Image) decodeRawSegment(raw *io.SectionReader) (r io.ReadCloser, err error) {
	compression := im.Compression()
	switch compression {
	case CompressionNone:
		return ioutil.NopCloser(raw), nil
	case CompressionJPEG:
		return im.decodeJPEGSegment(raw)
	case CompressionLZW:
		r = lzw.NewReader(raw, lzw.MSB, 8)
		return r, nil
	case CompressionDeflate, CompressionDeflateOld:
		return zlib.NewReader(raw)
	default:
		return nil, UnsupportedError(fmt.Sprintf("compression %d", int(compression)))
	}
}

func (im *Image) decodeJPEGSegment(raw *io.SectionReader) (r io.ReadCloser, err error) {
	var goImage image.Image

	jpegTableTag := im.Tag[TagJPEGTables]
	if jpegTableTag == nil {
		goImage, err = jpeg.Decode(raw)
		if err != nil {
			return nil, err
		}
	} else {
		const jpegSOI = "\xFF\xD8"
		const jpegEOI = "\xFF\xD9"
		if len(jpegTableTag.Data) < 4 ||
			string(jpegTableTag.Data[:2]) != jpegSOI ||
			string(jpegTableTag.Data[len(jpegTableTag.Data)-2:]) != jpegEOI {
			return nil, FormatError("invalid JPEGTable")
		}
		buf := make([]byte, 2)
		_, err = io.ReadFull(raw, buf)
		if err != nil {
			return nil, err
		}
		if string(buf) != jpegSOI {
			return nil, FormatError("invalid JPEG data: missing SOI")
		}

		rd := io.MultiReader(
			bytes.NewReader(jpegTableTag.Data[:len(jpegTableTag.Data)-2]),
			io.NewSectionReader(raw, 2, raw.Size()-2),
		)
		goImage, err = jpeg.Decode(rd)
		if err != nil {
			return nil, err
		}
	}

	width := goImage.Bounds().Dx()
	height := goImage.Bounds().Dy()
	goImageYCbCr, ok := goImage.(*image.YCbCr)
	if !ok {
		return nil, UnsupportedError("JPEG in non-YCbCr colorspace")
	}
	buf := make([]uint8, width*height*3)
	for iy := 0; iy < height; iy++ {
		for ix := 0; ix < width; ix++ {
			c := goImageYCbCr.YCbCrAt(ix, iy)
			buf[3*(iy*width+ix)+0] = c.Y
			buf[3*(iy*width+ix)+1] = c.Cb
			buf[3*(iy*width+ix)+2] = c.Cr
		}
	}

	return ioutil.NopCloser(bytes.NewReader(buf)), nil
}

// DecodeImage decodes image data and copy data to buffer.
func (im *Image) DecodeImage(buffer interface{}) error {
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

	//
	// Strip mode
	//
	numStrips := im.NumStrips()
	if numStrips > 0 {
		// Find expected number of strips
		rowsPerStrip, ok := im.Tag[TagRowsPerStrip].Uint()
		if !ok {
			return FormatError("invalid RowsPerStrip")
		}
		ny := (int)(math.Ceil((float64)(height) / (float64)(rowsPerStrip)))
		if ny != numStrips {
			return FormatError(fmt.Sprintf("expecting %d strips, found %d strips", ny, numStrips))
		}

		// Process strips
		samplesPerRow := width * samplePerPixel
		rowBuf := make([]byte, samplesPerRow*bitDepth/8)
		numRows := int(rowsPerStrip)
		destOffset := 0
		for iy := 0; iy < ny; iy++ {
			// Get image data reader
			stripReader, err := im.StripReader(iy)
			if err != nil {
				return fmt.Errorf("cannot get strip %d: %v", iy, err)
			}

			// Last strip may not be a full strip, and is not padded
			if iy == ny-1 {
				numRows = height - iy*int(rowsPerStrip)
			}

			// Process rows in the strip
			for iRow := 0; iRow < numRows; iRow++ {
				_, err := io.ReadFull(stripReader, rowBuf)
				if err != nil {
					return fmt.Errorf("cannot read strip %d row %d: %v", iy, iRow, err)
				}

				if buf, ok := buffer.([]uint8); ok {
					copy(buf[destOffset:destOffset+samplesPerRow], rowBuf)
					destOffset += samplesPerRow
				}
				if buf, ok := buffer.([]uint16); ok {
					// Interpret data as uint16
					var rowBuf16 []uint16
					if im.Header.ByteOrder == binary.LittleEndian {
						rowBuf16 = (*(*[1 << 48]uint16)(unsafe.Pointer(&rowBuf[0])))[:samplesPerRow]
						copy(buf[destOffset:destOffset+samplesPerRow], rowBuf16)
						destOffset += samplesPerRow
					} else {
						srcOffset := 0
						for i := samplesPerRow; i > 0; i-- {
							buf[destOffset] = ((uint16)(rowBuf[srcOffset]) << 8) | (uint16)(rowBuf[srcOffset+1])
							srcOffset += 2
							destOffset++
						}
					}

				}
			}

		}

		return nil
	}

	//
	// Tile mode
	//
	numTiles := im.NumTiles()
	if numTiles > 0 {
		// Find expected number of tiles
		tileHeight, ok := im.Tag[TagTileLength].Uint()
		if !ok {
			return FormatError("invalid TagTileLength")
		}
		ny := (int)(math.Ceil((float64)(height) / (float64)(tileHeight)))
		tileWidth, ok := im.Tag[TagTileWidth].Uint()
		if !ok {
			return FormatError("invalid TagTileWidth")
		}
		nx := (int)(math.Ceil((float64)(width) / (float64)(tileWidth)))

		if nx*ny != numTiles {
			return FormatError(fmt.Sprintf("expecting %d tiles, found %d tiles", nx*ny, numTiles))
		}

		// Process tiles
		samplesPerTile := int(tileWidth*tileHeight) * samplePerPixel
		tileBuf := make([]byte, samplesPerTile*bitDepth/8)
		for iy := 0; iy < ny; iy++ {
			for ix := 0; ix < nx; ix++ {
				// Read all pixels of the tile
				tileReader, err := im.TileReader(iy*nx + ix)
				if err != nil {
					return fmt.Errorf("cannot get tile %d: %v", iy*nx+ix, err)
				}
				_, err = io.ReadFull(tileReader, tileBuf)
				if err != nil {
					return fmt.Errorf("cannot read tile %d: %v", iy*nx+ix, err)
				}

				// Set up tile coordinate and bounds
				tileX0 := ix * int(tileWidth)
				tileY0 := iy * int(tileHeight)
				destOffset := (tileY0*width + tileX0) * samplePerPixel
				destStride := width * samplePerPixel

				tileDX := int(tileWidth) * samplePerPixel
				tileDY := int(tileHeight)
				if ix == nx-1 {
					tileDX = (width - ix*int(tileWidth)) * samplePerPixel
				}
				if iy == ny-1 {
					tileDY = height - iy*int(tileHeight)
				}
				srcOffset := 0
				srcStride := int(tileWidth) * samplePerPixel

				// Copy tile data line-by-line
				for iRow := 0; iRow < tileDY; iRow++ {
					if buf, ok := buffer.([]uint8); ok {
						copy(buf[destOffset:destOffset+tileDX], tileBuf[srcOffset:srcOffset+tileDX])
						destOffset += destStride
						srcOffset += srcStride
					}
					if buf, ok := buffer.([]uint16); ok {
						// Interpret data as uint16
						var tileBuf16 []uint16
						if im.Header.ByteOrder == binary.LittleEndian {
							tileBuf16 = (*(*[1 << 48]uint16)(unsafe.Pointer(&tileBuf[2*srcOffset])))[:tileDX]
							copy(buf[destOffset:destOffset+tileDX], tileBuf16[:tileDX])
							destOffset += destStride
							srcOffset += srcStride
						} else {
							for i := 0; i < tileDX; i++ {
								buf[destOffset+i] = ((uint16)(tileBuf[2*(srcOffset+i)]) << 8) | (uint16)(tileBuf[2*(srcOffset+i)+1])
							}
							destOffset += destStride
							srcOffset += srcStride
						}
					}
				}
			}
		}
		return nil
	}

	return fmt.Errorf("no strip or tile found")
}
