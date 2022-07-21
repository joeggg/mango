package embedded

import (
	"errors"
	"fmt"
	"io"
	"mango"
	"mango/pb"
	"reflect"

	"google.golang.org/protobuf/proto"
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
	data := reflect.New(t).Elem().Interface().(pb.CSVCMsg_ServerInfo)

	size, err := p.readVarU(32)
	if err != nil {
		return err
	}
	if p.Length-p.TruePos < size*8 {
		return errors.New("invalid embedded size given")
	}

	payload, err := p.readByteArray(size)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(payload, &data)
	if err != nil {
		return err
	}
	mango.PrintStruct(data)
	return nil
}

func (p *EmbeddedParser) readVarU(max int) (int, error) {
	max = ((max + 6) / 7) * 7
	result := 0
	shift := 0
	for {
		num, err := p.readIntNBit(8)
		if err != nil {
			return result, err
		}
		result |= (num & 0x7F) << shift
		shift += 7
		if num&0x80 == 0 || shift == max {
			return result, nil
		}
	}
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

func (p *EmbeddedParser) readByteArray(size int) ([]byte, error) {
	data := make([]byte, 0, size)
	for i := 0; i < size; i++ {
		num, err := p.readIntNBit(8)
		if err != nil {
			return data, err
		}
		data = append(data, byte(num))
	}
	return data, nil
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
