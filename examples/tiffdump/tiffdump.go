package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Andeling/tiff"
)

func main() {
	args := os.Args
	if len(args) == 0 {
		return
	}

	f, err := os.Open(args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d, err := tiff.NewDecoder(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("0x0000  ByteOrder: <%s> Version: %d <%s>\n", d.Header.ByteOrder, d.Header.Version, d.Header.Version)

	it := d.Iter()
	for it.Next() {
		index := it.Index()
		im := it.Image()
		fmt.Printf("0x%04x  Directory %d: offset %d(0x%x) next %d(0x%x)\n", im.Offset, index, im.Offset, im.Offset, im.OffsetNext, im.OffsetNext)
		for _, id := range im.TagID() {
			tag := im.Tag[id]
			fmt.Printf("0x%04x  %5d | %-26s | %-9s | %5d | %s\n", tag.OffsetEntry, id, id, tag.Type, tag.Count, formatStringTruncated(tag))
		}
		if im.Exif != nil {
			fmt.Printf("0x%04x      ExifIFD: offset %d(0x%x) next %d\n", im.Exif.Offset, im.Exif.Offset, im.Exif.Offset, im.Exif.OffsetNext)
			for _, id := range im.Exif.TagID() {
				tag := im.Exif.Tag[id]
				fmt.Printf("0x%04x      %5d | %-26s | %-9s | %5d | %s\n", tag.OffsetEntry, id, id, tag.Type, tag.Count, formatStringTruncated(tag))
			}
		}
		if im.SubImage != nil {
			for subindex, subim := range im.SubImage {
				fmt.Printf("0x%04x      Sub-directory %d: offset %d(0x%x) next %d\n", subim.Offset, subindex, subim.Offset, subim.Offset, subim.OffsetNext)
				for _, id := range subim.TagID() {
					tag := subim.Tag[id]
					fmt.Printf("0x%04x      %5d | %-26s | %-9s | %5d | %s\n", tag.OffsetEntry, id, id, tag.Type, tag.Count, formatStringTruncated(tag))
				}
			}
		}
	}
	if err := it.Err(); err != nil {
		log.Fatal(err)
	}
}

func formatStringTruncated(v interface{}) string {
	s := fmt.Sprintf("%#v", v)
	if len(s) > 54 {
		return s[:25] + "...." + s[len(s)-25:len(s)]
	}
	return s
}
