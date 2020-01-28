package tiff_test

import (
	"bytes"
	"encoding/binary"
	"image"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	tiff "github.com/Andeling/tiff"
	_ "golang.org/x/image/tiff"
)

func BenchmarkParseIFD(b *testing.B) {
	buf, err := ioutil.ReadFile("../tiff_testdata/rgb16_le.tif")
	if err != nil {
		b.Fatal(err)
	}
	rd := bytes.NewReader(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d, err := tiff.NewDecoder(rd)
		if err != nil {
			b.Fatal(err)
		}

		it := d.Iter()
		for it.Next() {
		}
		if err := it.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode_RGB16_LE_GoImage(b *testing.B) {
	buf, err := ioutil.ReadFile("../tiff_testdata/rgb16_le.tif")
	if err != nil {
		b.Fatal(err)
	}
	rd := bytes.NewReader(buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rd.Seek(0, io.SeekStart)
		_, _, err = image.Decode(rd)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode_RGB16_BE_GoImage(b *testing.B) {
	buf, err := ioutil.ReadFile("../tiff_testdata/rgb16_be.tif")
	if err != nil {
		b.Fatal(err)
	}
	rd := bytes.NewReader(buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rd.Seek(0, io.SeekStart)
		_, _, err = image.Decode(rd)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode_RGB16_LE(b *testing.B) {
	buf, err := ioutil.ReadFile("../tiff_testdata/rgb16_le.tif")
	if err != nil {
		b.Fatal(err)
	}
	rd := bytes.NewReader(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d, err := tiff.NewDecoder(rd)
		if err != nil {
			b.Fatal(err)
		}

		it := d.Iter()
		for it.Next() {
			im := it.Image()
			width, height := im.WidthHeight()
			samplesPerPixel := im.SamplesPerPixel()
			buf := make([]uint16, width*height*samplesPerPixel)
			err = im.DecodeImage(buf)
			if err != nil {
				b.Fatal(err)
			}
		}
		if err := it.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode_RGB16_LE_Strip(b *testing.B) {
	buf, err := ioutil.ReadFile("../tiff_testdata/rgb16_le_strip.tif")
	if err != nil {
		b.Fatal(err)
	}
	rd := bytes.NewReader(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d, err := tiff.NewDecoder(rd)
		if err != nil {
			b.Fatal(err)
		}

		it := d.Iter()
		for it.Next() {
			im := it.Image()
			width, height := im.WidthHeight()
			samplesPerPixel := im.SamplesPerPixel()
			buf := make([]uint16, width*height*samplesPerPixel)
			err = im.DecodeImage(buf)
			if err != nil {
				b.Fatal(err)
			}
		}
		if err := it.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode_RGB16_LE_Tile(b *testing.B) {
	buf, err := ioutil.ReadFile("../tiff_testdata/rgb16_le_tile.tif")
	if err != nil {
		b.Fatal(err)
	}
	rd := bytes.NewReader(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d, err := tiff.NewDecoder(rd)
		if err != nil {
			b.Fatal(err)
		}

		it := d.Iter()
		for it.Next() {
			im := it.Image()
			width, height := im.WidthHeight()
			samplesPerPixel := im.SamplesPerPixel()
			buf := make([]uint16, width*height*samplesPerPixel)
			err = im.DecodeImage(buf)
			if err != nil {
				b.Fatal(err)
			}
		}
		if err := it.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode_RGB16_BE(b *testing.B) {
	buf, err := ioutil.ReadFile("../tiff_testdata/rgb16_be.tif")
	if err != nil {
		b.Fatal(err)
	}
	rd := bytes.NewReader(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d, err := tiff.NewDecoder(rd)
		if err != nil {
			b.Fatal(err)
		}

		it := d.Iter()
		for it.Next() {
			im := it.Image()
			width, height := im.WidthHeight()
			samplesPerPixel := im.SamplesPerPixel()
			buf := make([]uint16, width*height*samplesPerPixel)
			err = im.DecodeImage(buf)
			if err != nil {
				b.Fatal(err)
			}
		}
		if err := it.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode_RGB16_Zip_LE(b *testing.B) {
	buf, err := ioutil.ReadFile("../tiff_testdata/rgb16_le_zip.tif")
	if err != nil {
		b.Fatal(err)
	}
	rd := bytes.NewReader(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d, err := tiff.NewDecoder(rd)
		if err != nil {
			b.Fatal(err)
		}

		it := d.Iter()
		for it.Next() {
			im := it.Image()
			width, height := im.WidthHeight()
			samplesPerPixel := im.SamplesPerPixel()
			buf := make([]uint16, width*height*samplesPerPixel)
			err = im.DecodeImage(buf)
			if err != nil {
				b.Fatal(err)
			}
		}
		if err := it.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode_RGB16_Zstd_LE(b *testing.B) {
	buf, err := ioutil.ReadFile("../tiff_testdata/rgb16_le_zstd.tif")
	if err != nil {
		b.Fatal(err)
	}
	rd := bytes.NewReader(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d, err := tiff.NewDecoder(rd)
		if err != nil {
			b.Fatal(err)
		}

		it := d.Iter()
		for it.Next() {
			im := it.Image()
			width, height := im.WidthHeight()
			samplesPerPixel := im.SamplesPerPixel()
			buf := make([]uint16, width*height*samplesPerPixel)
			err = im.DecodeImage(buf)
			if err != nil {
				b.Fatal(err)
			}
		}
		if err := it.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode_RGB16_LZW_LE(b *testing.B) {
	buf, err := ioutil.ReadFile("../tiff_testdata/rgb16_le_lzw.tif")
	if err != nil {
		b.Fatal(err)
	}
	rd := bytes.NewReader(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d, err := tiff.NewDecoder(rd)
		if err != nil {
			b.Fatal(err)
		}

		it := d.Iter()
		for it.Next() {
			im := it.Image()
			width, height := im.WidthHeight()
			samplesPerPixel := im.SamplesPerPixel()
			buf := make([]uint16, width*height*samplesPerPixel)
			err = im.DecodeImage(buf)
			if err != nil {
				b.Fatal(err)
			}
		}
		if err := it.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode_RGB8_JPEG_LE(b *testing.B) {
	buf, err := ioutil.ReadFile("../tiff_testdata/rgb8_le_jpeg_90.tif")
	if err != nil {
		b.Fatal(err)
	}
	rd := bytes.NewReader(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d, err := tiff.NewDecoder(rd)
		if err != nil {
			b.Fatal(err)
		}

		it := d.Iter()
		for it.Next() {
			im := it.Image()
			width, height := im.WidthHeight()
			samplesPerPixel := im.SamplesPerPixel()
			buf := make([]uint8, width*height*samplesPerPixel)
			err = im.DecodeImage(buf)
			if err != nil {
				b.Fatal(err)
			}
		}
		if err := it.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode_RGB8_LE_GoImage(b *testing.B) {
	buf, err := ioutil.ReadFile("../tiff_testdata/rgb8_le.tif")
	if err != nil {
		b.Fatal(err)
	}
	rd := bytes.NewReader(buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rd.Seek(0, io.SeekStart)
		_, _, err = image.Decode(rd)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode_RGB8_LE(b *testing.B) {
	buf, err := ioutil.ReadFile("../tiff_testdata/rgb8_le.tif")
	if err != nil {
		b.Fatal(err)
	}
	rd := bytes.NewReader(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d, err := tiff.NewDecoder(rd)
		if err != nil {
			b.Fatal(err)
		}

		it := d.Iter()
		for it.Next() {
			im := it.Image()
			width, height := im.WidthHeight()
			samplesPerPixel := im.SamplesPerPixel()
			buf := make([]uint8, width*height*samplesPerPixel)
			err = im.DecodeImage(buf)
			if err != nil {
				b.Fatal(err)
			}
		}
		if err := it.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode_RGB8_ZIP(b *testing.B) {
	buf, err := ioutil.ReadFile("../tiff_testdata/rgb8_le_zip.tif")
	if err != nil {
		b.Fatal(err)
	}
	rd := bytes.NewReader(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d, err := tiff.NewDecoder(rd)
		if err != nil {
			b.Fatal(err)
		}

		it := d.Iter()
		for it.Next() {
			im := it.Image()
			width, height := im.WidthHeight()
			samplesPerPixel := im.SamplesPerPixel()
			buf := make([]uint8, width*height*samplesPerPixel)
			err = im.DecodeImage(buf)
			if err != nil {
				b.Fatal(err)
			}
		}
		if err := it.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode_RGB8_Zstd(b *testing.B) {
	buf, err := ioutil.ReadFile("../tiff_testdata/rgb8_le_zstd.tif")
	if err != nil {
		b.Fatal(err)
	}
	rd := bytes.NewReader(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d, err := tiff.NewDecoder(rd)
		if err != nil {
			b.Fatal(err)
		}

		it := d.Iter()
		for it.Next() {
			im := it.Image()
			width, height := im.WidthHeight()
			samplesPerPixel := im.SamplesPerPixel()
			buf := make([]uint8, width*height*samplesPerPixel)
			err = im.DecodeImage(buf)
			if err != nil {
				b.Fatal(err)
			}
		}
		if err := it.Err(); err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkDecode_RGB8_LZW(b *testing.B) {
	buf, err := ioutil.ReadFile("../tiff_testdata/rgb8_le_lzw.tif")
	if err != nil {
		b.Fatal(err)
	}
	rd := bytes.NewReader(buf)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d, err := tiff.NewDecoder(rd)
		if err != nil {
			b.Fatal(err)
		}

		it := d.Iter()
		for it.Next() {
			im := it.Image()
			width, height := im.WidthHeight()
			samplesPerPixel := im.SamplesPerPixel()
			buf := make([]uint8, width*height*samplesPerPixel)
			err = im.DecodeImage(buf)
			if err != nil {
				b.Fatal(err)
			}
		}
		if err := it.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func TestDecode_RGB16(t *testing.T) {
	f_groundtruth, err := os.Open("../tiff_testdata/rgb16_le.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer f_groundtruth.Close()

	buf_groundtruth := make([]uint16, 1365*2048*3)
	err = binary.Read(f_groundtruth, binary.LittleEndian, buf_groundtruth)
	if err != nil {
		t.Fatal(err)
	}

	filenames := []string{
		// Byte order
		"../tiff_testdata/rgb16_le.tif",
		"../tiff_testdata/rgb16_be.tif",

		// Strip
		"../tiff_testdata/rgb16_le_strip.tif",
		"../tiff_testdata/rgb16_le_tile.tif",
		"../tiff_testdata/rgb16_be_tile.tif",

		// Lossless compression
		"../tiff_testdata/rgb16_le_zip.tif",
		"../tiff_testdata/rgb16_le_lzw.tif",
		"../tiff_testdata/rgb16_le_zstd.tif",
	}

	for _, filename := range filenames {
		func() {
			f, err := os.Open(filename)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			d, err := tiff.NewDecoder(f)
			if err != nil {
				t.Fatal(err)
			}

			it := d.Iter()
			it.Next()
			im := it.Image()
			width, height := im.WidthHeight()
			samplesPerPixel := im.SamplesPerPixel()
			buf := make([]uint16, width*height*samplesPerPixel)
			err = im.DecodeImage(buf)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(buf, buf_groundtruth) {
				t.Fatalf("%s failed: data inconsistent with ground truth", filename)
			}
		}()
	}
}

func TestDecode_RGB8(t *testing.T) {
	f_groundtruth, err := os.Open("../tiff_testdata/rgb8_le.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer f_groundtruth.Close()

	buf_groundtruth := make([]uint8, 1365*2048*3)
	err = binary.Read(f_groundtruth, binary.LittleEndian, buf_groundtruth)
	if err != nil {
		t.Fatal(err)
	}

	filenames := []string{
		// Byte order
		"../tiff_testdata/rgb8_le.tif",

		"../tiff_testdata/rgb8_le_strip.tif",
		"../tiff_testdata/rgb8_le_tile.tif",

		// Lossless compression
		"../tiff_testdata/rgb8_le_lzw.tif",
		"../tiff_testdata/rgb8_le_zip.tif",
		"../tiff_testdata/rgb8_le_zstd.tif",
	}

	for _, filename := range filenames {
		func() {
			f, err := os.Open(filename)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			d, err := tiff.NewDecoder(f)
			if err != nil {
				t.Fatal(err)
			}

			it := d.Iter()
			it.Next()
			im := it.Image()
			width, height := im.WidthHeight()
			samplesPerPixel := im.SamplesPerPixel()
			buf := make([]uint8, width*height*samplesPerPixel)
			err = im.DecodeImage(buf)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(buf, buf_groundtruth) {
				t.Fatalf("%s failed: data inconsistent with ground truth", filename)
			}
		}()
	}
}
