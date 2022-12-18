package rmqr

import (
	"bytes"
	"testing"
)

func TestNew1(t *testing.T) {
	qr, err := New(LevelM, []byte("123456789012"))
	if err != nil {
		t.Fatal(err)
	}
	if qr.Version != R7x43 {
		t.Errorf("unexpected version: got %v, want %v", qr.Version, R7x43)
	}
	if len(qr.Segments) != 1 {
		t.Fatalf("unexpected the length of segment: got %d, want %d", len(qr.Segments), 1)
	}
	if qr.Segments[0].Mode != ModeNumeric {
		t.Errorf("got %v, want %v", qr.Segments[0].Mode, ModeNumeric)
	}
	if !bytes.Equal(qr.Segments[0].Data, []byte("123456789012")) {
		t.Errorf("got %q, want %q", qr.Segments[0].Data, "123456789012")
	}
}

func TestEncodeToBitmap1(t *testing.T) {
	qr := &QRCode{
		Version: R15x59,
		Level:   LevelM,
		Segments: []Segment{
			{
				Mode: ModeBytes,
				Data: []byte("Rectangular Micro QR Code (rMQR)"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	got := img.Pix

	want := []byte{
		0b11111110, 0b10101010, 0b10111010, 0b10101010, 0b10101011, 0b10101010, 0b10101010, 0b11100000,
		0b10000010, 0b00110100, 0b01101110, 0b00000001, 0b10101010, 0b11101001, 0b11110011, 0b00100000,
		0b10111010, 0b11101000, 0b11111011, 0b10100100, 0b01111111, 0b11001011, 0b11101111, 0b01100000,
		0b10111010, 0b10001001, 0b11001111, 0b00001110, 0b10100010, 0b00001010, 0b10110000, 0b11000000,
		0b10111010, 0b11010101, 0b00011100, 0b01010000, 0b11101111, 0b11011100, 0b01010010, 0b00100000,
		0b10000010, 0b10010011, 0b10101011, 0b11101110, 0b01010110, 0b01010000, 0b11100110, 0b10000000,
		0b11111110, 0b00110101, 0b10110110, 0b00011011, 0b01101011, 0b01000001, 0b01101011, 0b01100000,
		0b00000000, 0b10101110, 0b00000100, 0b11110000, 0b01100000, 0b00100111, 0b10010000, 0b00000000,
		0b11101010, 0b00010101, 0b00110001, 0b10101011, 0b10011111, 0b00010111, 0b11010101, 0b01100000,
		0b01101110, 0b10100001, 0b11001011, 0b00000001, 0b11001110, 0b00010010, 0b00010101, 0b10000000,
		0b10010010, 0b10000110, 0b00111011, 0b10100100, 0b01101011, 0b01001001, 0b01010111, 0b11100000,
		0b00100001, 0b10101110, 0b01001101, 0b00001110, 0b11100110, 0b11010110, 0b10111110, 0b00100000,
		0b11100000, 0b11000000, 0b00111110, 0b01010000, 0b00011011, 0b10111101, 0b01001110, 0b10100000,
		0b10111011, 0b11111011, 0b11101100, 0b11101110, 0b10101110, 0b10100010, 0b01010110, 0b00100000,
		0b11101010, 0b10101010, 0b10111010, 0b10101010, 0b10101011, 0b10101010, 0b10101011, 0b11100000,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func TestEncodeToBitmap2(t *testing.T) {
	qr := &QRCode{
		Version: R7x43,
		Level:   LevelM,
		Segments: []Segment{
			{
				Mode: ModeNumeric,
				Data: []byte("123456789012"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	got := img.Pix

	want := []byte{
		0b11111110, 0b10101010, 0b10101110, 0b10101010, 0b10101010, 0b11100000,
		0b10000010, 0b01010111, 0b11101010, 0b00001000, 0b11011000, 0b10100000,
		0b10111010, 0b10111000, 0b10011110, 0b11010110, 0b10111111, 0b11100000,
		0b10111010, 0b01100110, 0b11100000, 0b00011111, 0b11000010, 0b00100000,
		0b10111010, 0b00101100, 0b01111110, 0b01101111, 0b10110010, 0b10100000,
		0b10000010, 0b11110111, 0b11111010, 0b11010010, 0b10111010, 0b00100000,
		0b11111110, 0b10101010, 0b10101110, 0b10101010, 0b10101011, 0b11100000,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}

func TestEncodeToBitmap3(t *testing.T) {
	qr := &QRCode{
		Version: R7x139,
		Level:   LevelM,
		Segments: []Segment{
			{
				Mode: ModeNumeric,
				Data: []byte("123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012"),
			},
		},
	}
	img, err := qr.EncodeToBitmap()
	if err != nil {
		t.Fatal(err)
	}
	got := img.Pix

	want := []byte{
		0b11111110, 0b10101010, 0b10101010, 0b10111010, 0b10101010, 0b10101010, 0b10101011, 0b10101010, 0b10101010, 0b10101010, 0b10111010, 0b10101010, 0b10101010, 0b10101011, 0b10101010, 0b10101010, 0b10101010, 0b11100000,
		0b10000010, 0b01011011, 0b00010110, 0b01101100, 0b01101011, 0b11000110, 0b01111010, 0b10101110, 0b10010010, 0b00100101, 0b10101100, 0b11001001, 0b10100000, 0b00100010, 0b11111011, 0b10110000, 0b01011000, 0b10100000,
		0b10111010, 0b01111000, 0b10011100, 0b10111001, 0b00100101, 0b10110001, 0b00101011, 0b11101001, 0b00001100, 0b00000101, 0b11111110, 0b00011110, 0b10010100, 0b10001111, 0b11011011, 0b01001010, 0b01100111, 0b11100000,
		0b10111010, 0b10100001, 0b11010111, 0b10100001, 0b01001001, 0b11000001, 0b01001100, 0b01110110, 0b00101011, 0b10001100, 0b00001100, 0b01011101, 0b01100101, 0b00101000, 0b01011100, 0b00000101, 0b01011010, 0b00100000,
		0b10111010, 0b10111000, 0b11000110, 0b01111100, 0b10011001, 0b11001001, 0b11000011, 0b10010001, 0b01100000, 0b01010010, 0b01111011, 0b11101001, 0b01101011, 0b01010011, 0b11010110, 0b10100000, 0b01100010, 0b10100000,
		0b10000010, 0b00001001, 0b10100110, 0b10101100, 0b11001001, 0b11011110, 0b00101110, 0b11011000, 0b11001101, 0b00100110, 0b11101110, 0b00101111, 0b01000000, 0b00001110, 0b10001010, 0b01101010, 0b10100110, 0b00100000,
		0b11111110, 0b10101010, 0b10101010, 0b10111010, 0b10101010, 0b10101010, 0b10101011, 0b10101010, 0b10101010, 0b10101010, 0b10111010, 0b10101010, 0b10101010, 0b10101011, 0b10101010, 0b10101010, 0b10101011, 0b11100000,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %08b, want %08b", got, want)
	}
}
