package mango

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"os"

	"github.com/golang/snappy"
	"google.golang.org/protobuf/proto"
)

var (
	headerLength    = 8
	header          = "PBDEMS2\x00"
	compressedMask  = 64
	varintMask      = int(math.Pow(2, 32) - 1) // 32 bit mask
	varintBlockSize = 7
	varintMaxSize   = 32
)

type ReplayParser struct {
	file *os.File
}

func NewReplayParser(filename string) (*ReplayParser, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	r := &ReplayParser{
		file: file,
	}
	return r, nil
}

func (rp *ReplayParser) Initialise() error {
	// Header handling
	if data, _ := rp.readString(headerLength); data != header {
		return errors.New("failed to read header")
	}
	return nil
}

func (rp *ReplayParser) GetSummary() (proto.Message, error) {
	// Offset handling
	if offset, err := rp.readUint32(); err != nil {
		return nil, err
	} else if _, err = rp.file.Seek(int64(offset), 0); err != nil {
		return nil, err
	} else if packet, err := rp.GetPacket(); err != nil {
		return nil, err
	} else if err := packet.Parse(); err != nil {
		return nil, err
	} else {
		return packet.Message, nil
	}
}

func (rp *ReplayParser) ParseReplay() ([]*Packet, error) {
	var packets []*Packet
	rp.readBytes(headerLength)
	for i := 0; i < 10; i++ {
		p, err := rp.GetPacket()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return packets, nil
		}
		err = p.Parse()
		if err != nil {
			return packets, err
		}
		packets = append(packets, p)
	}
	return packets, nil
}

func (rp *ReplayParser) GetPacket() (*Packet, error) {
	if kind, err := rp.readVarint32(); err != nil {
		return nil, err
	} else if tick, err := rp.readVarint32(); err != nil {
		return nil, err
	} else if size, err := rp.readVarint32(); err != nil {
		return nil, err
	} else if message, err := rp.readBytes(size); err != nil {
		return nil, err
	} else {
		isCompressed := (kind & compressedMask) > 0
		packet := &Packet{
			Kind:         kind,
			Tick:         tick,
			Size:         size,
			IsCompressed: isCompressed,
			RawMessage:   message,
		}
		if isCompressed {
			packet.Kind &= ^compressedMask
			packet.RawMessage, err = snappy.Decode(nil, message)
			if err != nil {
				return nil, err
			}
		}
		return packet, nil
	}
}

func (rp *ReplayParser) readVarint32() (int, error) {
	total, shift := 0, 0
	for {
		current, err := rp.readByte()
		if err != nil {
			return total, err
		}
		total |= (int(current) & 0x7F) << shift

		if current&0x80 == 0 {
			return total & varintMask, nil
		}

		shift += varintBlockSize
		if shift >= varintMaxSize {
			return total, errors.New("invalid varint")
		}
	}
}

func (rp *ReplayParser) readUint32() (uint32, error) {
	data, err := rp.readBytes(4)
	return binary.LittleEndian.Uint32(data), err
}

func (rp *ReplayParser) readString(size int) (string, error) {
	data, err := rp.readBytes(size)
	return string(data), err
}

func (rp *ReplayParser) readByte() (byte, error) {
	data, err := rp.readBytes(1)
	return data[0], err
}

func (rp *ReplayParser) readBytes(size int) (data []byte, err error) {
	data = make([]byte, size)
	_, err = rp.file.Read(data)
	return
}
