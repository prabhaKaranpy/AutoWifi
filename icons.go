package main

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

func getIcon(c color.Color) []byte {
	width := 256
	height := 256
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)

	var buf bytes.Buffer
	png.Encode(&buf, img)

	return pngToIco(buf.Bytes())
}

func pngToIco(pngBytes []byte) []byte {
	// Create a simple ICO format wrapper around the PNG data
	// See: https://en.wikipedia.org/wiki/ICO_(file_format)

	header := []byte{
		0, 0, // Reserved
		1, 0, // Type (1=Icon)
		1, 0, // Count (1 Image)
	}

	entry := make([]byte, 16)
	// Width (0 = 256px) - Note: 0 byte means 256
	entry[0] = 0
	// Height (0 = 256px)
	entry[1] = 0
	// ColorCount (0 = No palette)
	entry[2] = 0
	// Reserved
	entry[3] = 0
	// Planes (1)
	entry[4] = 1
	entry[5] = 0
	// BitCount (32 = RGBA)
	entry[6] = 32
	entry[7] = 0
	// BytesInRes (Size of PNG)
	binary.LittleEndian.PutUint32(entry[8:], uint32(len(pngBytes)))
	// ImageOffset (Header 6 + Entry 16 = 22)
	binary.LittleEndian.PutUint32(entry[12:], 22)

	var ico bytes.Buffer
	ico.Write(header)
	ico.Write(entry)
	ico.Write(pngBytes)

	return ico.Bytes()
}

var (
	iconGreen  = getIcon(color.RGBA{0, 255, 0, 255})
	iconYellow = getIcon(color.RGBA{255, 255, 0, 255})
	iconRed    = getIcon(color.RGBA{255, 0, 0, 255})
)
