package embedded

import (
	"errors"
	"io"
)

var UBitVarMap = []int{0, 4, 8, 28}

type EmbeddedDecoder struct {
	Buffer  []byte
	Length  int
	BytePos int
	TruePos int
}

/*
	Decode the packet header and extract the raw data
*/
func (p *EmbeddedDecoder) Decode(data []byte) (*EmbeddedPacket, error) {
	p.initialiseData(data)
	if kind, err := p.readUBitVar(); err != nil {
		return nil, err
	} else if size, err := p.readVarU(32); err != nil {
		return nil, err
	} else if p.Length-p.TruePos < size*8 {
		return nil, errors.New("invalid embedded size given")
	} else if payload, err := p.readByteArray(size); err != nil {
		return nil, err
	} else {
		return &EmbeddedPacket{Kind: kind, RawData: payload}, nil
	}
}

/*
	Initialise an embedded message decoder with raw byte data
*/
func (p *EmbeddedDecoder) initialiseData(data []byte) {
	p.Buffer = data
	p.Length = 8 * len(data)
	p.BytePos, p.TruePos = 0, 0
}

func (p *EmbeddedDecoder) readVarU(max int) (int, error) {
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

func (p *EmbeddedDecoder) readUBitVar() (int, error) {
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

func (p *EmbeddedDecoder) readByteArray(size int) ([]byte, error) {
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

func (p *EmbeddedDecoder) readIntNBit(n int) (int, error) {
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

func (p *EmbeddedDecoder) readBit() (b byte, err error) {
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
