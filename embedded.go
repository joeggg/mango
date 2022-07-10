package mango

import (
	"fmt"
	"io"
)

type EmbeddedParser struct {
	Buffer  []byte
	Length  int
	BytePos int
	TruePos int
}

func NewEmbeddedParser(data []int8) *EmbeddedParser {
	buffer := make([]byte, 0, len(data))
	for _, item := range data {
		buffer = append(buffer, byte(item))
	}
	return &EmbeddedParser{Buffer: buffer, Length: 8 * len(buffer), BytePos: 0, TruePos: 0}
}

func (p *EmbeddedParser) Parse() error {
	fmt.Println(p.Buffer)
	kind, err := p.readIntNBit(33)
	if err != nil && err != io.EOF {
		return err
	}
	fmt.Println(kind)
	return nil
}

func (p *EmbeddedParser) readIntNBit(n int) (int, error) {
	total := 0
	for i := 0; i < n; i++ {
		val, err := p.readBit()
		if err != nil {
			return total, err
		}
		total |= int(val) << i
	}
	return total, nil
}

func (p *EmbeddedParser) readBit() (b byte, err error) {
	if p.TruePos >= p.Length {
		return b, io.EOF
	}
	shift := p.TruePos - 8*p.BytePos
	b = (p.Buffer[p.BytePos] & (1 << shift)) >> shift
	p.TruePos++

	if p.TruePos%8 == 0 {
		p.BytePos++
	}
	return
}
