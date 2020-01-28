package tiff

//
// TagID
//

// TagID identifies the field
type TagID uint16

func (t TagID) String() string {
	return tagName[t]
}

type byTagID []TagID

func (a byTagID) Len() int           { return len(a) }
func (a byTagID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTagID) Less(i, j int) bool { return a[i] < a[j] }

//
// ExifTagID
//

// ExifTagID identifies the Exif tag field
type ExifTagID uint16

func (t ExifTagID) String() string {
	return exifTagName[t]
}

type byExifTagID []ExifTagID

func (a byExifTagID) Len() int           { return len(a) }
func (a byExifTagID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byExifTagID) Less(i, j int) bool { return a[i] < a[j] }

//
// GPSTagID
//

// GPSTagID identifies the Exif tag field
type GPSTagID uint16

func (t GPSTagID) String() string {
	return gpsTagName[t]
}

type byGPSTagID []GPSTagID

func (a byGPSTagID) Len() int           { return len(a) }
func (a byGPSTagID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byGPSTagID) Less(i, j int) bool { return a[i] < a[j] }

//
// GPSTaInteroperabilityTagIDgID
//

// InteroperabilityTagID identifies the Exif tag field
type InteroperabilityTagID uint16

func (t InteroperabilityTagID) String() string {
	return interoperabilityTagName[t]
}

type byInteroperabilityTagID []InteroperabilityTagID

func (a byInteroperabilityTagID) Len() int           { return len(a) }
func (a byInteroperabilityTagID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byInteroperabilityTagID) Less(i, j int) bool { return a[i] < a[j] }

const (
	NewFiletypeImage        = 0x0 // full resolution version
	NewFiletypeReducedImage = 0x1 // reduced resolution version
	NewFiletypePage         = 0x2 // one page of many
	NewFiletypeMask         = 0x4 // transparency mask

	CompressionNone       = 1     // Baseline. None
	CompressionCCITTRLE   = 2     // CCITT modified Huffman RLE
	CompressionCCITTG3    = 3     // CCITT T.4 (Group 3)
	CompressionCCITTG4    = 4     // CCITT T.6 (Group 4)
	CompressionLZW        = 5     // LZW
	CompressionJPEGOld    = 6     // Old-style JPEG
	CompressionJPEG       = 7     // JPEG
	CompressionDeflate    = 8     // Deflate compression
	CompressionDeflateOld = 32946 // Non-standard. Deflate compression
	CompressionPackBits   = 32773 // Baseline. PackBits compression
	CompressionLossyJPEG  = 34892 // Lossy JPEG
	CompressionLZMA       = 34925 // LZMA2
	CompressionZstd       = 50000 // ZSTD: WARNING not registered in Adobe-maintained registry
	CompressionWebP       = 50001 // WEBP: WARNING not registered in Adobe-maintained registry

	PhotometricWhiteIsZero = 0     // WhiteIsZero
	PhotometricBlackIsZero = 1     // BlackIsZero
	PhotometricRGB         = 2     // RGB
	PhotometricPalette     = 3     // Palette color
	PhotometricMask        = 4     // Transparency mask
	PhotometricSeparated   = 5     // Separated (usually CMYK)
	PhotometricYCbCr       = 6     // YCbCr
	PhotometricCIELab      = 8     // 1976 CIE L*a*b*
	PhotometricICCLab      = 9     // ICC L*a*b*
	PhotometricITULab      = 10    // ITU L*a*b*
	PhotometricCFA         = 32803 // Color filter array
	PhotometricLogL        = 32844 // CIE Log2(L)
	PhotometricLogLuv      = 32845 // CIE Log2(L) (u',v')
	PhotometricLinearRaw   = 34925 // LinearRaw (or de-mosaiced CFA data)

	PlanarConfigContig   = 1 // The component values for each pixel are stored contiguously, e.g, RGBRGB...RGB
	PlanarConfigSeparate = 2 // The components are stored in separate component planes, e.g., RR..RGG..GBB..B
)

const (
	// Baseline & Extended TIFF
	TagNewSubfileType         TagID = 254   // Baseline. A general indication of the kind of data contained in this subfile.
	TagSubfileType            TagID = 255   // Baseline. A general indication of the kind of data contained in this subfile.
	TagImageWidth             TagID = 256   // Baseline. The number of columns in the image, i.e., the number of pixels per row.
	TagImageLength            TagID = 257   // Baseline. The number of rows of pixels in the image.
	TagBitsPerSample          TagID = 258   // Baseline. Number of bits per component.
	TagCompression            TagID = 259   // Baseline. Compression scheme used on the image data.
	TagPhotometric            TagID = 262   // Baseline. The color space of the image data.
	TagThreshholding          TagID = 263   // Baseline. For black and white TIFF files that represent shades of gray, the technique used to convert from gray to black and white pixels.
	TagCellWidth              TagID = 264   // Baseline. The width of the dithering or halftoning matrix used to create a dithered or halftoned bilevel file.
	TagCellLength             TagID = 265   // Baseline. The length of the dithering or halftoning matrix used to create a dithered or halftoned bilevel file.
	TagFillOrder              TagID = 266   // Baseline. The logical order of bits within a byte.
	TagDocumentName           TagID = 269   // Extended. The name of the document from which this image was scanned.
	TagImageDescription       TagID = 270   // Baseline. A string that describes the subject of the image.
	TagMake                   TagID = 271   // Baseline. The scanner manufacturer.
	TagModel                  TagID = 272   // Baseline. The scanner model name or number.
	TagStripOffsets           TagID = 273   // Baseline. For each strip, the byte offset of that strip.
	TagOrientation            TagID = 274   // Baseline. The orientation of the image with respect to the rows and columns.
	TagSamplesPerPixel        TagID = 277   // Baseline. The number of components per pixel.
	TagRowsPerStrip           TagID = 278   // Baseline. The number of rows per strip.
	TagStripByteCounts        TagID = 279   // Baseline. For each strip, the number of bytes in the strip after compression.
	TagMinSampleValue         TagID = 280   // Baseline. The minimum component value used.
	TagMaxSampleValue         TagID = 281   // Baseline. The maximum component value used.
	TagXResolution            TagID = 282   // Baseline. The number of pixels per ResolutionUnit in the ImageWidth direction.
	TagYResolution            TagID = 283   // Baseline. The number of pixels per ResolutionUnit in the ImageLength direction.
	TagPlanarConfig           TagID = 284   // Baseline. How the components of each pixel are stored.
	TagPageName               TagID = 285   // Extended. The name of the page from which this image was scanned.
	TagXPosition              TagID = 286   // Extended. X position of the image.
	TagYPosition              TagID = 287   // Extended. Y position of the image.
	TagFreeOffsets            TagID = 288   // Baseline. For each string of contiguous unused bytes in a TIFF file, the byte offset of the string.
	TagFreeByteCounts         TagID = 289   // Baseline. For each string of contiguous unused bytes in a TIFF file, the number of bytes in the string.
	TagGrayResponseUnit       TagID = 290   // Baseline. The precision of the information contained in the GrayResponseCurve.
	TagGrayResponseCurve      TagID = 291   // Baseline. For grayscale data, the optical density of each possible pixel value.
	TagT4Options              TagID = 292   // Extended. Options for Group 3 Fax compression
	TagT6Options              TagID = 293   // Extended. Options for Group 4 Fax compression
	TagResolutionUnit         TagID = 296   // Baseline. The unit of measurement for XResolution and YResolution.
	TagPageNumber             TagID = 297   // Extended. The page number of the page from which this image was scanned.
	TagTransferFunction       TagID = 301   // Extended. Describes a transfer function for the image in tabular style.
	TagSoftware               TagID = 305   // Baseline. Name and version number of the software package(s) used to create the image.
	TagDateTime               TagID = 306   // Baseline. Date and time of image creation.
	TagArtist                 TagID = 315   // Baseline. Person who created the image.
	TagHostComputer           TagID = 316   // Baseline. The computer and/or operating system in use at the time of image creation.
	TagPredictor              TagID = 317   // Extended. A mathematical operator that is applied to the image data before an encoding scheme is applied.
	TagWhitePoint             TagID = 318   // Extended. The chromaticity of the white point of the image.
	TagPrimaryChromaticities  TagID = 319   // Extended. The chromaticities of the primaries of the image.
	TagColorMap               TagID = 320   // Baseline & Supplement1. A color map for palette color images.
	TagHalftoneHints          TagID = 321   // Extended. Conveys to the halftone function the range of gray levels within a colorimetrically-specified image that should retain tonal detail.
	TagTileWidth              TagID = 322   // Extended. The tile width in pixels. This is the number of columns in each tile.
	TagTileLength             TagID = 323   // Extended. The tile length (height) in pixels. This is the number of rows in each tile.
	TagTileOffsets            TagID = 324   // Extended. For each tile, the byte offset of that tile, as compressed and stored on disk.
	TagTileByteCounts         TagID = 325   // Extended. For each tile, the number of (compressed) bytes in that tile.
	TagBadFaxLines            TagID = 326   // Extended. Used in the TIFF-F standard, denotes the number of 'bad' scan lines encountered by the facsimile device.
	TagCleanFaxData           TagID = 327   // Extended. Used in the TIFF-F standard, indicates if 'bad' lines encountered during reception are stored in the data, or if 'bad' lines have been replaced by the receiver.
	TagConsecutiveBadFaxLines TagID = 328   // Extended. Used in the TIFF-F standard, denotes the maximum number of consecutive 'bad' scanlines received.
	TagSubIFDs                TagID = 330   // Supplement1. Offset to child IFDs.
	TagInkSet                 TagID = 332   // Extended. The set of inks used in a separated (Photometric = 5) image.
	TagInkNames               TagID = 333   // Extended. The name of each ink used in a separated image.
	TagNumberOfInks           TagID = 334   // Extended. The number of inks.
	TagDotRange               TagID = 336   // Extended. The component values that correspond to a 0% dot and 100% dot.
	TagTargetPrinter          TagID = 337   // Extended. A description of the printing environment for which this separation is intended.
	TagExtraSamples           TagID = 338   // Baseline. Description of extra components.
	TagSampleFormat           TagID = 339   // Extended. Specifies how to interpret each data sample in a pixel.
	TagSMinSampleValue        TagID = 340   // Extended. Specifies the minimum sample value.
	TagSMaxSampleValue        TagID = 341   // Extended. Specifies the maximum sample value.
	TagTransferRange          TagID = 342   // Extended. Expands the range of the TransferFunction.
	TagClipPath               TagID = 343   // Supplement1. Mirrors the essentials of PostScript's path creation functionality.
	TagXClipPathUnits         TagID = 344   // Supplement1. The number of units that span the width of the image, in terms of integer ClipPath coordinates.
	TagYClipPathUnits         TagID = 345   // Supplement1. The number of units that span the height of the image, in terms of integer ClipPath coordinates.
	TagIndexed                TagID = 346   // Supplement1. Aims to broaden the support for indexed images to include support for any color space.
	TagJPEGTables             TagID = 347   // Supplement1. JPEG quantization and/or Huffman tables.
	TagOPIProxy               TagID = 351   // Supplement1. OPI-related.
	TagGlobalParametersIFD    TagID = 400   // Extended. Used in the TIFF-FX standard to point to an IFD containing tags that are globally applicable to the complete TIFF file.
	TagProfileType            TagID = 401   // Extended. Used in the TIFF-FX standard, denotes the type of data stored in this file or IFD.
	TagFaxProfile             TagID = 402   // Extended. Used in the TIFF-FX standard, denotes the 'profile' that applies to this file.
	TagCodingMethods          TagID = 403   // Extended. Used in the TIFF-FX standard, indicates which coding methods are used in the file.
	TagVersionYear            TagID = 404   // Extended. Used in the TIFF-FX standard, denotes the year of the standard specified by the FaxProfile field.
	TagModeNumber             TagID = 405   // Extended. Used in the TIFF-FX standard, denotes the mode of the standard specified by the FaxProfile field.
	TagDecode                 TagID = 433   // Extended. Used in the TIFF-F and TIFF-FX standards, holds information about the ITULAB (Photometric = 10) encoding.
	TagDefaultImageColor      TagID = 434   // Extended. Defined in the Mixed Raster Content part of RFC 2301, is the default color needed in areas where no image is available.
	TagJPEGProc               TagID = 512   // Extended. Old-style JPEG compression field. TechNote2 invalidates this part of the specification. But it is still quite widely used.
	TagJPEGIFOffset           TagID = 513   // Extended. Old-style JPEG compression field. TechNote2 invalidates this part of the specification. But it is still quite widely used.
	TagJPEGIFByteCount        TagID = 514   // Extended. Old-style JPEG compression field. TechNote2 invalidates this part of the specification. But it is still quite widely used.
	TagJPEGRestartInterval    TagID = 515   // Extended. Old-style JPEG compression field. TechNote2 invalidates this part of the specification. But it is still quite widely used.
	TagJPEGLosslessPredictors TagID = 517   // Extended. Old-style JPEG compression field. TechNote2 invalidates this part of the specification. But it is still quite widely used.
	TagJPEGPointTransforms    TagID = 518   // Extended. Old-style JPEG compression field. TechNote2 invalidates this part of the specification. But it is still quite widely used.
	TagJPEGQTables            TagID = 519   // Extended. Old-style JPEG compression field. TechNote2 invalidates this part of the specification. But it is still quite widely used.
	TagJPEGDCTables           TagID = 520   // Extended. Old-style JPEG compression field. TechNote2 invalidates this part of the specification. But it is still quite widely used.
	TagJPEGACTables           TagID = 521   // Extended. Old-style JPEG compression field. TechNote2 invalidates this part of the specification. But it is still quite widely used.
	TagYCbCrCoefficients      TagID = 529   // Extended. The transformation from RGB to YCbCr image data.
	TagYCbCrSubSampling       TagID = 530   // Extended. Specifies the subsampling factors used for the chrominance components of a YCbCr image.
	TagYCbCrPositioning       TagID = 531   // Extended. Specifies the positioning of subsampled chrominance components relative to luminance samples.
	TagReferenceBlackWhite    TagID = 532   // Extended. Specifies a pair of headroom and footroom image data values (codes) for each pixel component.
	TagStripRowCounts         TagID = 559   // Extended. Defined in the Mixed Raster Content part of RFC 2301, used to replace RowsPerStrip for IFDs with variable-sized strips.
	TagXMP                    TagID = 700   // Extended. XML packet containing XMP metadata
	TagImageID                TagID = 32781 // Supplement1. OPI-related.
	TagCopyright              TagID = 33432 // Baseline. Copyright notice.
	TagImageLayer             TagID = 34732 // Extended. Defined in the Mixed Raster Content part of RFC 2301, used to denote the particular function of this Image in the mixed raster scheme.

	// From libtiff
	TagColorResponseUnit TagID = 300
	TagT82Options        TagID = 435

	// TIFF/EP
	TagCFARepeatPatternDim      TagID = 33421
	TagCFAPattern               TagID = 33422
	TagSelfTimeMode             TagID = 34859
	TagFocalPlaneXResolution    TagID = 37390
	TagFocalPlaneYResolution    TagID = 37391
	TagFocalPlaneResolutionUnit TagID = 37392
	TagImageNumber              TagID = 37393
	TagSecurityClassification   TagID = 37394
	TagImageHistory             TagID = 37395
	TagExposureIndex            TagID = 37397
	TagTIFFEPStandardID         TagID = 37398
	TagSensingMethod            TagID = 37399

	// EXIF Private IFDs
	TagExifIFD             TagID = 34665
	TagGPSIFD              TagID = 34853
	TagInteroperabilityIFD TagID = 40965

	// EXIF Tags
	ExifTagExposureTime             ExifTagID = 33434
	ExifTagFNumber                  ExifTagID = 33437
	ExifTagExposureProgram          ExifTagID = 34850
	ExifTagSpectralSensitivity      ExifTagID = 34852
	ExifTagISOSpeedRatings          ExifTagID = 34855
	ExifTagOECF                     ExifTagID = 34856
	ExifTagSensitivityType          ExifTagID = 34864
	ExifTagRecommendedExposureIndex ExifTagID = 34866
	ExifTagExifVersion              ExifTagID = 36864
	ExifTagDateTimeOriginal         ExifTagID = 36867
	ExifTagDateTimeDigitized        ExifTagID = 36868
	ExifTagComponentsConfiguration  ExifTagID = 37121
	ExifTagCompressedBitsPerPixel   ExifTagID = 37122
	ExifTagShutterSpeedValue        ExifTagID = 37377
	ExifTagApertureValue            ExifTagID = 37378
	ExifTagBrightnessValue          ExifTagID = 37379
	ExifTagExposureBiasValue        ExifTagID = 37380
	ExifTagMaxApertureValue         ExifTagID = 37381
	ExifTagSubjectDistance          ExifTagID = 37382
	ExifTagMeteringMode             ExifTagID = 37383
	ExifTagLightSource              ExifTagID = 37384
	ExifTagFlash                    ExifTagID = 37385
	ExifTagFocalLength              ExifTagID = 37386
	ExifTagSubjectArea              ExifTagID = 37396
	ExifTagMakerNote                ExifTagID = 37500
	ExifTagUserComment              ExifTagID = 37510
	ExifTagSubsecTime               ExifTagID = 37520
	ExifTagSubsecTimeOriginal       ExifTagID = 37521
	ExifTagSubsecTimeDigitized      ExifTagID = 37522
	ExifTagFlashpixVersion          ExifTagID = 40960
	ExifTagColorSpace               ExifTagID = 40961
	ExifTagPixelXDimension          ExifTagID = 40962
	ExifTagPixelYDimension          ExifTagID = 40963
	ExifTagRelatedSoundFile         ExifTagID = 40964
	ExifTagFlashEnergy              ExifTagID = 41483
	ExifTagSpatialFrequencyResponse ExifTagID = 41484
	ExifTagFocalPlaneXResolution    ExifTagID = 41486
	ExifTagFocalPlaneYResolution    ExifTagID = 41487
	ExifTagFocalPlaneResolutionUnit ExifTagID = 41488
	ExifTagSubjectLocation          ExifTagID = 41492
	ExifTagExposureIndex            ExifTagID = 41493
	ExifTagSensingMethod            ExifTagID = 41495
	ExifTagFileSource               ExifTagID = 41728
	ExifTagSceneType                ExifTagID = 41729
	ExifTagCFAPattern               ExifTagID = 41730
	ExifTagCustomRendered           ExifTagID = 41985
	ExifTagExposureMode             ExifTagID = 41986
	ExifTagWhiteBalance             ExifTagID = 41987
	ExifTagDigitalZoomRatio         ExifTagID = 41988
	ExifTagFocalLengthIn35mmFilm    ExifTagID = 41989
	ExifTagSceneCaptureType         ExifTagID = 41990
	ExifTagGainControl              ExifTagID = 41991
	ExifTagContrast                 ExifTagID = 41992
	ExifTagSaturation               ExifTagID = 41993
	ExifTagSharpness                ExifTagID = 41994
	ExifTagDeviceSettingDescription ExifTagID = 41995
	ExifTagSubjectDistanceRange     ExifTagID = 41996
	ExifTagImageUniqueID            ExifTagID = 42016
	ExifTagCameraOwnerName          ExifTagID = 42032
	ExifTagBodySerialNumber         ExifTagID = 42033
	ExifTagLensSpecification        ExifTagID = 42034
	ExifTagLensMake                 ExifTagID = 42035
	ExifTagLensModel                ExifTagID = 42036
	ExifTagLensSerialNumber         ExifTagID = 42037

	// DNG 1.0
	TagDNGVersion             TagID = 50706
	TagDNGBackwardVersion     TagID = 50707
	TagUniqueCameraModel      TagID = 50708
	TagLocalizedCameraModel   TagID = 50709
	TagCFAPlaneColor          TagID = 50710
	TagCFALayout              TagID = 50711
	TagLinearizationTable     TagID = 50712
	TagBlackLevelRepeatDim    TagID = 50713
	TagBlackLevel             TagID = 50714
	TagBlackLevelDeltaH       TagID = 50715
	TagBlackLevelDeltaV       TagID = 50716
	TagWhiteLevel             TagID = 50717
	TagDefaultScale           TagID = 50718
	TagBestQualityScale       TagID = 50780
	TagDefaultCropOrigin      TagID = 50719
	TagDefaultCropSize        TagID = 50720
	TagCalibrationIlluminant1 TagID = 50778
	TagCalibrationIlluminant2 TagID = 50779
	TagColorMatrix1           TagID = 50721
	TagColorMatrix2           TagID = 50722
	TagCameraCalibration1     TagID = 50723
	TagCameraCalibration2     TagID = 50724
	TagReductionMatrix1       TagID = 50725
	TagReductionMatrix2       TagID = 50726
	TagAnalogBalance          TagID = 50727
	TagAsShotNeutral          TagID = 50728
	TagAsShotWhiteXY          TagID = 50729
	TagBaselineExposure       TagID = 50730
	TagBaselineNoise          TagID = 50731
	TagBaselineSharpness      TagID = 50732
	TagBayerGreenSplit        TagID = 50733
	TagLinearResponseLimit    TagID = 50734
	TagCameraSerialNumber     TagID = 50735
	TagLensInfo               TagID = 50736
	TagChromaBlurRadius       TagID = 50737
	TagAntiAliasStrength      TagID = 50738
	TagDNGPrivateData         TagID = 50740
	TagMakerNoteSafety        TagID = 50741

	// DNG 1.1
	TagShadowScale             TagID = 50739
	TagRawDataUniqueID         TagID = 50781
	TagOriginalRawFileName     TagID = 50827
	TagOriginalRawFileData     TagID = 50828
	TagActiveArea              TagID = 50829
	TagMaskedAreas             TagID = 50830
	TagAsShotICCProfile        TagID = 50831
	TagAsShotPreProfileMatrix  TagID = 50832
	TagCurrentICCProfile       TagID = 50833
	TagCurrentPreProfileMatrix TagID = 50834

	// DNG 1.2
	TagColorimetricReference       TagID = 50879
	TagCameraCalibrationSignature  TagID = 50931
	TagProfileCalibrationSignature TagID = 50932
	TagExtraCameraProfiles         TagID = 50933
	TagAsShotProfileName           TagID = 50934
	TagNoiseReductionApplied       TagID = 50935
	TagProfileName                 TagID = 50936
	TagProfileHueSatMapDims        TagID = 50937
	TagProfileHueSatMapData1       TagID = 50938
	TagProfileHueSatMapData2       TagID = 50939
	TagProfileToneCurve            TagID = 50940
	TagProfileEmbedPolicy          TagID = 50941
	TagProfileCopyright            TagID = 50942
	TagForwardMatrix1              TagID = 50964
	TagForwardMatrix2              TagID = 50965
	TagPreviewApplicationName      TagID = 50966
	TagPreviewApplicationVersion   TagID = 50967
	TagPreviewSettingsName         TagID = 50968
	TagPreviewSettingsDigest       TagID = 50969
	TagPreviewColorSpace           TagID = 50970
	TagPreviewDateTime             TagID = 50971
	TagRawImageDigest              TagID = 50972
	TagOriginalRawFileDigest       TagID = 50973
	TagSubTileBlockSize            TagID = 50974
	TagRowInterleaveFactor         TagID = 50975
	TagProfileLookTableDims        TagID = 50981
	TagProfileLookTableData        TagID = 50982

	// DNG 1.3
	TagOpcodeList1  TagID = 51008
	TagOpcodeList2  TagID = 51009
	TagOpcodeList3  TagID = 51022
	TagNoiseProfile TagID = 51041

	// DNG 1.4
	TagDefaultUserCrop              TagID = 51125
	TagDefaultBlackRender           TagID = 51110
	TagBaselineExposureOffset       TagID = 51109
	TagProfileLookTableEncoding     TagID = 51108
	TagProfileHueSatMapEncoding     TagID = 51107
	TagOriginalDefaultFinalSize     TagID = 51089
	TagOriginalBestQualityFinalSize TagID = 51090
	TagOriginalDefaultCropSize      TagID = 51091
	TagNewRawImageDigest            TagID = 51111
	TagRawToPreviewGain             TagID = 51112

	// DNG (Adobe DNG SDK)
	TagCacheBlob    TagID = 51113
	TagCacheVersion TagID = 51114
)

var tagName = map[TagID]string{
	// Baseline & Extended TIFF
	TagNewSubfileType:         "NewSubfileType",
	TagSubfileType:            "SubfileType",
	TagImageWidth:             "ImageWidth",
	TagImageLength:            "ImageLength",
	TagBitsPerSample:          "BitsPerSample",
	TagCompression:            "Compression",
	TagPhotometric:            "Photometric",
	TagThreshholding:          "Threshholding",
	TagCellWidth:              "CellWidth",
	TagCellLength:             "CellLength",
	TagFillOrder:              "FillOrder",
	TagDocumentName:           "DocumentName",
	TagImageDescription:       "ImageDescription",
	TagMake:                   "Make",
	TagModel:                  "Model",
	TagStripOffsets:           "StripOffsets",
	TagOrientation:            "Orientation",
	TagSamplesPerPixel:        "SamplesPerPixel",
	TagRowsPerStrip:           "RowsPerStrip",
	TagStripByteCounts:        "StripByteCounts",
	TagMinSampleValue:         "MinSampleValue",
	TagMaxSampleValue:         "MaxSampleValue",
	TagXResolution:            "XResolution",
	TagYResolution:            "YResolution",
	TagPlanarConfig:           "PlanarConfig",
	TagPageName:               "PageName",
	TagXPosition:              "XPosition",
	TagYPosition:              "YPosition",
	TagFreeOffsets:            "FreeOffsets",
	TagFreeByteCounts:         "FreeByteCounts",
	TagGrayResponseUnit:       "GrayResponseUnit",
	TagGrayResponseCurve:      "GrayResponseCurve",
	TagT4Options:              "T4Options",
	TagT6Options:              "T6Options",
	TagResolutionUnit:         "ResolutionUnit",
	TagPageNumber:             "PageNumber",
	TagTransferFunction:       "TransferFunction",
	TagSoftware:               "Software",
	TagDateTime:               "DateTime",
	TagArtist:                 "Artist",
	TagHostComputer:           "HostComputer",
	TagPredictor:              "Predictor",
	TagWhitePoint:             "WhitePoint",
	TagPrimaryChromaticities:  "PrimaryChromaticities",
	TagColorMap:               "ColorMap",
	TagHalftoneHints:          "HalftoneHints",
	TagTileWidth:              "TileWidth",
	TagTileLength:             "TileLength",
	TagTileOffsets:            "TileOffsets",
	TagTileByteCounts:         "TileByteCounts",
	TagBadFaxLines:            "BadFaxLines",
	TagCleanFaxData:           "CleanFaxData",
	TagConsecutiveBadFaxLines: "ConsecutiveBadFaxLines",
	TagSubIFDs:                "SubIFDs",
	TagInkSet:                 "InkSet",
	TagInkNames:               "InkNames",
	TagNumberOfInks:           "NumberOfInks",
	TagDotRange:               "DotRange",
	TagTargetPrinter:          "TargetPrinter",
	TagExtraSamples:           "ExtraSamples",
	TagSampleFormat:           "SampleFormat",
	TagSMinSampleValue:        "SMinSampleValue",
	TagSMaxSampleValue:        "SMaxSampleValue",
	TagTransferRange:          "TransferRange",
	TagClipPath:               "ClipPath",
	TagXClipPathUnits:         "XClipPathUnits",
	TagYClipPathUnits:         "YClipPathUnits",
	TagIndexed:                "Indexed",
	TagJPEGTables:             "JPEGTables",
	TagOPIProxy:               "OPIProxy",
	TagGlobalParametersIFD:    "GlobalParametersIFD",
	TagProfileType:            "ProfileType",
	TagFaxProfile:             "FaxProfile",
	TagCodingMethods:          "CodingMethods",
	TagVersionYear:            "VersionYear",
	TagModeNumber:             "ModeNumber",
	TagDecode:                 "Decode",
	TagDefaultImageColor:      "DefaultImageColor",
	TagJPEGProc:               "JPEGProc",
	TagJPEGIFOffset:           "JPEGIFOffset",
	TagJPEGIFByteCount:        "JPEGIFByteCount",
	TagJPEGRestartInterval:    "JPEGRestartInterval",
	TagJPEGLosslessPredictors: "JPEGLosslessPredictors",
	TagJPEGPointTransforms:    "JPEGPointTransforms",
	TagJPEGQTables:            "JPEGQTables",
	TagJPEGDCTables:           "JPEGDCTables",
	TagJPEGACTables:           "JPEGACTables",
	TagYCbCrCoefficients:      "YCbCrCoefficients",
	TagYCbCrSubSampling:       "YCbCrSubSampling",
	TagYCbCrPositioning:       "YCbCrPositioning",
	TagReferenceBlackWhite:    "ReferenceBlackWhite",
	TagStripRowCounts:         "StripRowCounts",
	TagXMP:                    "XMP",
	TagImageID:                "ImageID",
	TagCopyright:              "Copyright",
	TagImageLayer:             "ImageLayer",

	// TIFF/EP
	TagCFARepeatPatternDim:      "CFARepeatPatternDim",
	TagCFAPattern:               "CFAPattern",
	TagSelfTimeMode:             "SelfTimeMode",
	TagFocalPlaneXResolution:    "FocalPlaneXResolution",
	TagFocalPlaneYResolution:    "FocalPlaneYResolution",
	TagFocalPlaneResolutionUnit: "FocalPlaneResolutionUnit",
	TagImageNumber:              "ImageNumber",
	TagSecurityClassification:   "SecurityClassification",
	TagImageHistory:             "ImageHistory",
	TagExposureIndex:            "ExposureIndex",
	TagTIFFEPStandardID:         "TIFF/EPStandardID",
	TagSensingMethod:            "SensingMethod",

	// EXIF Private IFDs
	TagExifIFD:             "ExifIFD",
	TagGPSIFD:              "GPSIFD",
	TagInteroperabilityIFD: "InteroperabilityIFD",

	// DNG 1.0
	TagDNGVersion:             "DNGVersion",
	TagDNGBackwardVersion:     "DNGBackwardVersion",
	TagUniqueCameraModel:      "UniqueCameraModel",
	TagLocalizedCameraModel:   "LocalizedCameraModel",
	TagCFAPlaneColor:          "CFAPlaneColor",
	TagCFALayout:              "CFALayout",
	TagLinearizationTable:     "LinearizationTable",
	TagBlackLevelRepeatDim:    "BlackLevelRepeatDim",
	TagBlackLevel:             "BlackLevel",
	TagBlackLevelDeltaH:       "BlackLevelDeltaH",
	TagBlackLevelDeltaV:       "BlackLevelDeltaV",
	TagWhiteLevel:             "WhiteLevel",
	TagDefaultScale:           "DefaultScale",
	TagBestQualityScale:       "BestQualityScale",
	TagDefaultCropOrigin:      "DefaultCropOrigin",
	TagDefaultCropSize:        "DefaultCropSize",
	TagCalibrationIlluminant1: "CalibrationIlluminant1",
	TagCalibrationIlluminant2: "CalibrationIlluminant2",
	TagColorMatrix1:           "ColorMatrix1",
	TagColorMatrix2:           "ColorMatrix2",
	TagCameraCalibration1:     "CameraCalibration1",
	TagCameraCalibration2:     "CameraCalibration2",
	TagReductionMatrix1:       "ReductionMatrix1",
	TagReductionMatrix2:       "ReductionMatrix2",
	TagAnalogBalance:          "AnalogBalance",
	TagAsShotNeutral:          "AsShotNeutral",
	TagAsShotWhiteXY:          "AsShotWhiteXY",
	TagBaselineExposure:       "BaselineExposure",
	TagBaselineNoise:          "BaselineNoise",
	TagBaselineSharpness:      "BaselineSharpness",
	TagBayerGreenSplit:        "BayerGreenSplit",
	TagLinearResponseLimit:    "LinearResponseLimit",
	TagCameraSerialNumber:     "CameraSerialNumber",
	TagLensInfo:               "LensInfo",
	TagChromaBlurRadius:       "ChromaBlurRadius",
	TagAntiAliasStrength:      "AntiAliasStrength",
	TagDNGPrivateData:         "DNGPrivateData",
	TagMakerNoteSafety:        "MakerNoteSafety",

	// DNG 1.1
	TagShadowScale:             "ShadowScale",
	TagRawDataUniqueID:         "RawDataUniqueID",
	TagOriginalRawFileName:     "OriginalRawFileName",
	TagOriginalRawFileData:     "OriginalRawFileData",
	TagActiveArea:              "ActiveArea",
	TagMaskedAreas:             "MaskedAreas",
	TagAsShotICCProfile:        "AsShotICCProfile",
	TagAsShotPreProfileMatrix:  "AsShotPreProfileMatrix",
	TagCurrentICCProfile:       "CurrentICCProfile",
	TagCurrentPreProfileMatrix: "CurrentPreProfileMatrix",

	// DNG 1.2
	TagColorimetricReference:       "ColorimetricReference",
	TagCameraCalibrationSignature:  "CameraCalibrationSignature",
	TagProfileCalibrationSignature: "ProfileCalibrationSignature",
	TagExtraCameraProfiles:         "ExtraCameraProfiles",
	TagAsShotProfileName:           "AsShotProfileName",
	TagNoiseReductionApplied:       "NoiseReductionApplied",
	TagProfileName:                 "ProfileName",
	TagProfileHueSatMapDims:        "ProfileHueSatMapDims",
	TagProfileHueSatMapData1:       "ProfileHueSatMapData1",
	TagProfileHueSatMapData2:       "ProfileHueSatMapData2",
	TagProfileToneCurve:            "ProfileToneCurve",
	TagProfileEmbedPolicy:          "ProfileEmbedPolicy",
	TagProfileCopyright:            "ProfileCopyright",
	TagForwardMatrix1:              "ForwardMatrix1",
	TagForwardMatrix2:              "ForwardMatrix2",
	TagPreviewApplicationName:      "PreviewApplicationName",
	TagPreviewApplicationVersion:   "PreviewApplicationVersion",
	TagPreviewSettingsName:         "PreviewSettingsName",
	TagPreviewSettingsDigest:       "PreviewSettingsDigest",
	TagPreviewColorSpace:           "PreviewColorSpace",
	TagPreviewDateTime:             "PreviewDateTime",
	TagRawImageDigest:              "RawImageDigest",
	TagOriginalRawFileDigest:       "OriginalRawFileDigest",
	TagSubTileBlockSize:            "SubTileBlockSize",
	TagRowInterleaveFactor:         "RowInterleaveFactor",
	TagProfileLookTableDims:        "ProfileLookTableDims",
	TagProfileLookTableData:        "ProfileLookTableData",

	// DNG 1.3
	TagOpcodeList1:  "OpcodeList1",
	TagOpcodeList2:  "OpcodeList2",
	TagOpcodeList3:  "OpcodeList3",
	TagNoiseProfile: "NoiseProfile",

	// DNG 1.4
	TagDefaultUserCrop:              "DefaultUserCrop",
	TagDefaultBlackRender:           "DefaultBlackRender",
	TagBaselineExposureOffset:       "BaselineExposureOffset",
	TagProfileLookTableEncoding:     "ProfileLookTableEncoding",
	TagProfileHueSatMapEncoding:     "ProfileHueSatMapEncoding",
	TagOriginalDefaultFinalSize:     "OriginalDefaultFinalSize",
	TagOriginalBestQualityFinalSize: "OriginalBestQualityFinalSize",
	TagOriginalDefaultCropSize:      "OriginalDefaultCropSize",
	TagNewRawImageDigest:            "NewRawImageDigest",
	TagRawToPreviewGain:             "RawToPreviewGain",

	// DNG SDK
	TagCacheBlob:    "CacheBlob",
	TagCacheVersion: "CacheVersion",
}

var exifTagName = map[ExifTagID]string{
	// EXIF
	ExifTagExposureTime:             "ExposureTime",
	ExifTagFNumber:                  "FNumber",
	ExifTagExposureProgram:          "ExposureProgram",
	ExifTagSpectralSensitivity:      "SpectralSensitivity",
	ExifTagISOSpeedRatings:          "ISOSpeedRatings",
	ExifTagOECF:                     "OECF",
	ExifTagSensitivityType:          "SensitivityType",
	ExifTagRecommendedExposureIndex: "RecommendedExposureIndex",
	ExifTagExifVersion:              "ExifVersion",
	ExifTagDateTimeOriginal:         "DateTimeOriginal",
	ExifTagDateTimeDigitized:        "DateTimeDigitized",
	ExifTagComponentsConfiguration:  "ComponentsConfiguration",
	ExifTagCompressedBitsPerPixel:   "CompressedBitsPerPixel",
	ExifTagShutterSpeedValue:        "ShutterSpeedValue",
	ExifTagApertureValue:            "ApertureValue",
	ExifTagBrightnessValue:          "BrightnessValue",
	ExifTagExposureBiasValue:        "ExposureBiasValue",
	ExifTagMaxApertureValue:         "MaxApertureValue",
	ExifTagSubjectDistance:          "SubjectDistance",
	ExifTagMeteringMode:             "MeteringMode",
	ExifTagLightSource:              "LightSource",
	ExifTagFlash:                    "Flash",
	ExifTagFocalLength:              "FocalLength",
	ExifTagSubjectArea:              "SubjectArea",
	ExifTagMakerNote:                "MakerNote",
	ExifTagUserComment:              "UserComment",
	ExifTagSubsecTime:               "SubsecTime",
	ExifTagSubsecTimeOriginal:       "SubsecTimeOriginal",
	ExifTagSubsecTimeDigitized:      "SubsecTimeDigitized",
	ExifTagFlashpixVersion:          "FlashpixVersion",
	ExifTagColorSpace:               "ColorSpace",
	ExifTagPixelXDimension:          "PixelXDimension",
	ExifTagPixelYDimension:          "PixelYDimension",
	ExifTagRelatedSoundFile:         "RelatedSoundFile",
	ExifTagFlashEnergy:              "FlashEnergy",
	ExifTagSpatialFrequencyResponse: "SpatialFrequencyResponse",
	ExifTagFocalPlaneXResolution:    "FocalPlaneXResolution",
	ExifTagFocalPlaneYResolution:    "FocalPlaneYResolution",
	ExifTagFocalPlaneResolutionUnit: "FocalPlaneResolutionUnit",
	ExifTagSubjectLocation:          "SubjectLocation",
	ExifTagExposureIndex:            "ExposureIndex",
	ExifTagSensingMethod:            "SensingMethod",
	ExifTagFileSource:               "FileSource",
	ExifTagSceneType:                "SceneType",
	ExifTagCFAPattern:               "CFAPattern",
	ExifTagCustomRendered:           "CustomRendered",
	ExifTagExposureMode:             "ExposureMode",
	ExifTagWhiteBalance:             "WhiteBalance",
	ExifTagDigitalZoomRatio:         "DigitalZoomRatio",
	ExifTagFocalLengthIn35mmFilm:    "FocalLengthIn35mmFilm",
	ExifTagSceneCaptureType:         "SceneCaptureType",
	ExifTagGainControl:              "GainControl",
	ExifTagContrast:                 "Contrast",
	ExifTagSaturation:               "Saturation",
	ExifTagSharpness:                "Sharpness",
	ExifTagDeviceSettingDescription: "DeviceSettingDescription",
	ExifTagSubjectDistanceRange:     "SubjectDistanceRange",
	ExifTagImageUniqueID:            "ImageUniqueID",
	ExifTagCameraOwnerName:          "CameraOwnerName",
	ExifTagBodySerialNumber:         "BodySerialNumber",
	ExifTagLensSpecification:        "LensSpecification",
	ExifTagLensMake:                 "LensMake",
	ExifTagLensModel:                "LensModel",
	ExifTagLensSerialNumber:         "LensSerialNumber",
}

var gpsTagName = map[GPSTagID]string{}
var interoperabilityTagName = map[InteroperabilityTagID]string{}
