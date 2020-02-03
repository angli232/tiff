# tiff

[![GoDoc](https://godoc.org/github.com/Andeling/tiff?status.svg)](https://godoc.org/github.com/Andeling/tiff)

This is a Go package to read and write TIFF or TIFF-like files.

This package is still experimental. Some features are missing, especially in the case of the encoder.

## Features

| Category                 | Feature           | Decode      | Encode |
| ------------------------ | ----------------- | ----------- | ------ |
| **Format**               | Classic TIFF      | Yes         | Yes    |
|                          | BigTIFF           | Yes         | Yes    |
| **Metadata**             | TIFF tags         | Yes         | Yes    |
|                          | Exif tags         | Yes         | Yes    |
|                          | GPS tags          | Yes         | -      |
| **Lossless Compression** | LZW               | Yes         | Yes    |
|                          | Deflate           | Yes         | Yes    |
|                          | zstd              | Yes         | Yes    |
|                          | Lossless JPEG     | -           | -      |
| **Lossy Compression**    | Lossy JPEG        | 8-bit YCbCr | -      |
| **Bit depth**            | 1-bit             | -           | -      |
|                          | 8-bit             | Yes         | Yes    |
|                          | 16-bit            | Yes         | Yes    |
| **Planar Configuration** | Contig            | Yes         | Yes    |
|                          | Separate (planar) | -           | -      |
| **Segmented Images**     | Strip             | Yes         | Yes    |
|                          | Tile              | Yes         | -      |