package main

import (
	"bytes"
	"fmt"
	"go/format"
	"image"
	"log"
	"os"

	"github.com/shogo82148/qrcode/internal/bitmap"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

const timingPatternOffset = 6

func main() {

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "// Code generated by genbase/main.go; DO NOT EDIT.\n\n")
	fmt.Fprintf(&buf, "package qrcode\n\n")
	fmt.Fprintf(&buf, "import (\n")
	fmt.Fprintf(&buf, "\"image\"\n")
	fmt.Fprintf(&buf, "\"github.com/shogo82148/qrcode/internal/bitmap\"\n")
	fmt.Fprintf(&buf, ")\n\n")

	genMaskList(&buf)
	genBaseList(&buf)

	out, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("base_gen.go", out, 0o644); err != nil {
		log.Fatal(err)
	}
}

func genMaskList(buf *bytes.Buffer) {
	fmt.Fprintf(buf, "var maskList = []*bitmap.Image{\n")
	genMask(buf, func(i, j int) int { return (i + j) % 2 })
	genMask(buf, func(i, j int) int { return i % 2 })
	genMask(buf, func(i, j int) int { return j % 3 })
	genMask(buf, func(i, j int) int { return (i + j) % 3 })
	genMask(buf, func(i, j int) int { return (i/2 + j/3) % 2 })
	genMask(buf, func(i, j int) int { return i*j%2 + i*j%3 })
	genMask(buf, func(i, j int) int { return (i*j%2 + i*j%3) % 2 })
	genMask(buf, func(i, j int) int { return ((i+j)%2 + i*j%3) % 2 })
	fmt.Fprintf(buf, "}\n\n")
}

func genMask(buf *bytes.Buffer, f func(i, j int) int) {
	w, h := 184, 177
	img := bitmap.New(image.Rect(0, 0, w, h))
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			img.SetBinary(j, i, f(i, j) == 0)
		}
	}
	writeImage(buf, img)
}

func genBaseList(buf *bytes.Buffer) {
	imgList := []*bitmap.Image{nil}
	usedList := []*bitmap.Image{nil}

	for version := 1; version <= 40; version++ {
		img, used := newBase(version)
		imgList = append(imgList, img)
		usedList = append(usedList, used)
	}

	fmt.Fprintf(buf, "var baseList = []*bitmap.Image{\n")
	fmt.Fprintf(buf, "nil, // dummy\n")
	for version := 1; version <= 40; version++ {
		fmt.Fprintf(buf, "\n// version %d\n", version)
		writeImage(buf, imgList[version])
	}
	fmt.Fprintf(buf, "}\n\n")

	fmt.Fprintf(buf, "var usedList = []*bitmap.Image{\n")
	fmt.Fprintf(buf, "nil, // dummy\n")
	for version := 1; version <= 40; version++ {
		fmt.Fprintf(buf, "\n// version %d\n", version)
		writeImage(buf, usedList[version])
	}
	fmt.Fprintf(buf, "}\n")
}

var positions = [][]int{
	nil, // dummy

	{},      // Version 1
	{6, 18}, // Version 2
	{6, 22}, // Version 3
	{6, 26}, // Version 4
	{6, 30}, // Version 5
	{6, 34}, // Version 6

	{6, 22, 38}, // Version 7
	{6, 24, 42}, // Version 8
	{6, 26, 46}, // Version 9
	{6, 28, 50}, // Version 10
	{6, 30, 54}, // Version 11
	{6, 32, 58}, // Version 12
	{6, 34, 62}, // Version 13

	{6, 26, 46, 66}, // Version 14
	{6, 26, 48, 70}, // Version 15
	{6, 26, 50, 74}, // Version 16
	{6, 30, 54, 78}, // Version 17
	{6, 30, 56, 82}, // Version 18
	{6, 30, 58, 86}, // Version 19
	{6, 34, 62, 90}, // Version 20

	{6, 28, 50, 72, 94},  // Version 21
	{6, 26, 50, 74, 98},  // Version 22
	{6, 30, 54, 78, 102}, // Version 23
	{6, 28, 54, 80, 106}, // Version 24
	{6, 32, 58, 84, 110}, // Version 25
	{6, 30, 58, 86, 114}, // Version 26
	{6, 34, 62, 90, 118}, // Version 27

	{6, 26, 50, 74, 98, 122},  // Version 28
	{6, 30, 54, 78, 102, 126}, // Version 29
	{6, 26, 52, 77, 104, 130}, // Version 30
	{6, 30, 56, 82, 108, 134}, // Version 31
	{6, 34, 60, 86, 112, 138}, // Version 32
	{6, 30, 58, 86, 114, 142}, // Version 33
	{6, 34, 62, 90, 118, 146}, // Version 34

	{6, 30, 46, 78, 102, 126, 150}, // Version 35
	{6, 24, 48, 70, 102, 128, 154}, // Version 36
	{6, 28, 54, 80, 106, 132, 158}, // Version 37
	{6, 32, 58, 84, 110, 136, 162}, // Version 38
	{6, 26, 54, 82, 110, 138, 166}, // Version 39
	{6, 30, 58, 86, 114, 142, 170}, // Version 40
}

func newBase(version int) (*bitmap.Image, *bitmap.Image) {
	w := 16 + 4*version
	img := bitmap.New(image.Rect(0, 0, w+1, w+1))
	used := bitmap.New(image.Rect(0, 0, w+1, w+1))

	// timing pattern
	for i := 0; i <= w; i++ {
		img.SetBinary(i, timingPatternOffset, i%2 == 0)
		img.SetBinary(timingPatternOffset, i, i%2 == 0)
		used.SetBinary(i, timingPatternOffset, true)
		used.SetBinary(timingPatternOffset, i, true)
	}

	// finder pattern
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			d := max(abs(x-3), abs(y-3))
			c := bitmap.Color(d != 2 && d != 4)
			img.SetBinary(x, y, c)
			img.SetBinary(w-x, y, c)
			img.SetBinary(x, w-y, c)
			used.SetBinary(x, y, true)
			used.SetBinary(w-x, y, true)
			used.SetBinary(x, w-y, true)
		}
	}

	// position pattern
	for j, y := range positions[version] {
		for i, x := range positions[version] {
			if (i == 0 && j == 0) || (i == len(positions[version])-1 && j == 0) || (i == 0 && j == len(positions[version])-1) {
				// finder pattern
				continue
			}
			drawPositionPattern(img, used, x, y)
		}
	}

	// reserved space for format info
	for i := 0; i < 8; i++ {
		used.SetBinary(i, 8, true)
		used.SetBinary(8, i, true)
		used.SetBinary(8, w-i, true)
		used.SetBinary(w-i, 8, true)
	}
	used.SetBinary(8, 8, true)

	// reserved space for version info
	if version >= 7 {
		for i := 0; i < 6; i++ {
			used.SetBinary(i, w-10, true)
			used.SetBinary(i, w-9, true)
			used.SetBinary(i, w-8, true)

			used.SetBinary(w-10, i, true)
			used.SetBinary(w-9, i, true)
			used.SetBinary(w-8, i, true)
		}
	}

	return img, used
}

func drawPositionPattern(img, used *bitmap.Image, cx, cy int) {
	for y := -2; y <= 2; y++ {
		for x := -2; x <= 2; x++ {
			d := max(abs(x), abs(y))
			c := bitmap.Color(d != 1)
			img.SetBinary(x+cx, y+cy, c)
			used.SetBinary(x+cx, y+cy, true)
		}
	}
}

func writeImage(buf *bytes.Buffer, img *bitmap.Image) {
	fmt.Fprintf(buf, "{\n")
	fmt.Fprintf(buf, "Stride: %d,\n", img.Stride)
	fmt.Fprintf(buf, "Rect: image.Rect(%d, %d, %d, %d),\n",
		img.Rect.Min.X, img.Rect.Min.Y, img.Rect.Max.X, img.Rect.Max.Y)
	fmt.Fprintf(buf, "Pix: []byte{\n")
	for i := 0; i < len(img.Pix); i += img.Stride {
		for _, b := range img.Pix[i : i+img.Stride] {
			fmt.Fprintf(buf, "0b%08b, ", b)
		}
		fmt.Fprintln(buf)
	}
	fmt.Fprintf(buf, "},\n")
	fmt.Fprintf(buf, "},\n")
}
