package embedded

import (
	"errors"
	"fmt"
	"io"
	"reflect"
)

var UBitVarMap = []int{0, 4, 8, 28}

type EmbeddedParser struct {
	Buffer  []byte
	Length  int
	BytePos int
	TruePos int
}

func NewEmbeddedParser(data []byte) *EmbeddedParser {
	buffer := make([]byte, 0, len(data))
	for _, item := range data {
		buffer = append(buffer, byte(item))
	}
	return &EmbeddedParser{Buffer: buffer, Length: 8 * len(buffer), BytePos: 0, TruePos: 0}
}

func (p *EmbeddedParser) Parse() error {
	fmt.Println(p.Length)
	kind, err := p.readUBitVar()
	if err != nil {
		return err
	}
	t, ok := EmbeddedTypeMap[kind]
	if !ok {
		return errors.New("unknown embedded message type")
	}
	fmt.Println(kind, t)
	data := reflect.New(t).Elem().Interface()
	fmt.Println(data)
	size, err := p.readUBitVar()
	if err != nil {
		return err
	}
	fmt.Println(size, p.Length-p.TruePos)
	return nil
}

func (p *EmbeddedParser) readUBitVar() (int, error) {
	num, err := p.readIntNBit(6)
	if err != nil {
		return num, err
	}
	if ntype := num >> 4; ntype == 0 {
		return num, nil
	} else {
		extra, err := p.readIntNBit(UBitVarMap[ntype])
		if err != nil {
			return num, err
		}
		// Select
		return (num & 15) | extra<<4, nil
	}

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
