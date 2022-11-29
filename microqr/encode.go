package microqr

import (
	"errors"
	"fmt"

	bitmap "github.com/shogo82148/qrcode/bitmap"
	"github.com/shogo82148/qrcode/internal/bitstream"
	"github.com/shogo82148/qrcode/internal/reedsolomon"
)

// EncodeToBitmap encodes QR Code into bitmap image.
func (qr *QRCode) EncodeToBitmap() (*bitmap.Image, error) {
	if qr.Version < 1 || qr.Version > 4 {
		return nil, fmt.Errorf("microqr: invalid version: %d", qr.Version)
	}
	if qr.Level < 0 || qr.Level >= 4 {
		return nil, fmt.Errorf("microqr: invalid level: %d", qr.Level)
	}
	format := formatTable[qr.Version][qr.Level]
	if format < 0 {
		return nil, fmt.Errorf("microqr: invalid version-level pair: %d-%s", qr.Version, qr.Level)
	}

	var buf bitstream.Buffer
	if err := qr.encodeSegments(&buf); err != nil {
		return nil, err
	}

	w := 8 + 2*int(qr.Version)
	img := baseList[qr.Version].Clone()
	used := usedList[qr.Version]

	dy := -1
	x, y := w, w
	for {
		if !used.BinaryAt(x, y) {
			bit, err := buf.ReadBit()
			if err != nil {
				break
			}
			img.SetBinary(x, y, bit != 0)
		}
		x--
		if x < 0 {
			break
		}

		if !used.BinaryAt(x, y) {
			bit, err := buf.ReadBit()
			if err != nil {
				break
			}
			img.SetBinary(x, y, bit != 0)
		}
		x, y = x+1, y+dy
		if y < 0 || y > w {
			dy *= -1
			x, y = x-2, y+dy
		}
		if x < 0 {
			break
		}
	}

	mask := qr.Mask
	encoded := encodedFormat[(format<<2)|int(mask)]
	for i := 0; i < 8; i++ {
		img.SetBinary(8, i+1, (encoded>>i)&1 != 0)
		img.SetBinary(i+1, 8, (encoded>>(14-i))&1 != 0)
	}

	img.Mask(img, used, maskList[mask])

	return img.Export(), nil
}

func (qr *QRCode) encodeSegments(buf *bitstream.Buffer) error {
	for _, s := range qr.Segments {
		if err := s.encode(qr.Version, buf); err != nil {
			return err
		}
	}
	l := buf.Len()
	buf.WriteBitsLSB(0x00, int(8-l%8))

	capacity := capacityTable[qr.Version][qr.Level]
	n := capacity.Total - capacity.Data
	rs := reedsolomon.New(n)
	rs.Write(buf.Bytes())
	correction := rs.Sum(make([]byte, 0, n))
	for _, b := range correction {
		buf.WriteBitsLSB(uint64(b), 8)
	}
	return nil
}

func (s *Segment) encode(version Version, buf *bitstream.Buffer) error {
	switch s.Mode {
	case ModeNumeric:
		return s.encodeNumber(version, buf)
	case ModeAlphanumeric:
		return s.encodeAlphabet(version, buf)
	case ModeBytes:
		return s.encodeBytes(version, buf)
	default:
		return errors.New("qrcode: unknown mode")
	}
}

func (s *Segment) encodeNumber(version Version, buf *bitstream.Buffer) error {
	// validation
	var n int
	data := s.Data
	switch version {
	case 1:
		n = 3
	case 2:
		n = 4
	case 3:
		n = 5
	case 4:
		n = 6
	default:
		return fmt.Errorf("qrcode: invalid version: %d", version)
	}
	if len(data) >= 1<<n {
		return fmt.Errorf("qrcode: data is too long: %d", len(data))
	}

	for _, ch := range data {
		if ch < '0' || ch > '9' {
			return fmt.Errorf("qrcode: invalid character in number mode: %02x", ch)
		}
	}

	// mode
	switch version {
	case 2:
		buf.WriteBitsLSB(uint64(ModeBytes), 1)
	case 3:
		buf.WriteBitsLSB(uint64(ModeBytes), 2)
	case 4:
		buf.WriteBitsLSB(uint64(ModeBytes), 3)
	}

	// data length
	buf.WriteBitsLSB(uint64(len(s.Data)), n)

	// data
	for i := 0; i+2 < len(data); i += 3 {
		n1 := int(data[i+0] - '0')
		n2 := int(data[i+1] - '0')
		n3 := int(data[i+2] - '0')
		bits := n1*100 + n2*10 + n3
		buf.WriteBitsLSB(uint64(bits), 10)
	}

	switch len(data) % 3 {
	case 1:
		bits := data[len(data)-1] - '0'
		buf.WriteBitsLSB(uint64(bits), 4)
	case 2:
		n1 := int(data[len(data)-2] - '0')
		n2 := int(data[len(data)-1] - '0')
		bits := n1*10 + n2
		buf.WriteBitsLSB(uint64(bits), 7)
	}
	return nil
}

var bitToAlphanumeric = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:")
var alphabets [256]int

func init() {
	for i := range alphabets {
		alphabets[i] = -1
	}
	for i, ch := range bitToAlphanumeric {
		alphabets[ch] = i
	}
}

func (s *Segment) encodeAlphabet(version Version, buf *bitstream.Buffer) error {
	// validation
	var n int
	data := s.Data
	switch version {
	case 2:
		n = 3
	case 3:
		n = 4
	case 4:
		n = 5
	default:
		return fmt.Errorf("qrcode: invalid version: %d", version)
	}
	if len(data) >= 1<<n {
		return fmt.Errorf("qrcode: data is too long: %d", len(data))
	}

	for _, ch := range data {
		if alphabets[ch] < 0 {
			return fmt.Errorf("qrcode: invalid character in alphabet mode: %02x", ch)
		}
	}

	// mode
	switch version {
	case 2:
		buf.WriteBitsLSB(uint64(ModeBytes), 1)
	case 3:
		buf.WriteBitsLSB(uint64(ModeBytes), 2)
	case 4:
		buf.WriteBitsLSB(uint64(ModeBytes), 3)
	}

	// data length
	buf.WriteBitsLSB(uint64(len(data)), n)

	// data
	for i := 0; i+1 < len(data); i += 2 {
		n1 := alphabets[data[i+0]]
		n2 := alphabets[data[i+1]]
		bits := n1*45 + n2
		buf.WriteBitsLSB(uint64(bits), 11)
	}

	if len(data)%2 != 0 {
		bits := alphabets[data[len(data)-1]]
		buf.WriteBitsLSB(uint64(bits), 6)
	}
	return nil
}

func (s *Segment) encodeBytes(version Version, buf *bitstream.Buffer) error {
	// validation
	var n int
	data := s.Data
	switch version {
	case 3:
		n = 4
	case 4:
		n = 5
	default:
		return fmt.Errorf("qrcode: invalid version: %d", version)
	}
	if len(data) >= 1<<n {
		return fmt.Errorf("qrcode: data is too long: %d", len(data))
	}

	// mode
	switch version {
	case 3:
		buf.WriteBitsLSB(uint64(ModeBytes), 2)
	case 4:
		buf.WriteBitsLSB(uint64(ModeBytes), 3)
	}

	// data length
	buf.WriteBitsLSB(uint64(len(data)), n)

	// data
	for _, bits := range data {
		buf.WriteBitsLSB(uint64(bits), 8)
	}
	return nil
}
