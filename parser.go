package mango

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"mango/pb"
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

func (r *ReplayParser) Initialise() error {
	// Header handling
	if data, _ := r.readString(headerLength); data != header {
		return errors.New("failed to read header")
	}
	return nil
}

func (r *ReplayParser) GetSummary() (proto.Message, error) {
	// Offset handling
	if offset, err := r.readUint32(); err != nil {
		return nil, err
	} else if _, err = r.file.Seek(int64(offset), 0); err != nil {
		return nil, err
	} else if packet, err := r.GetPacket(); err != nil {
		return nil, err
	} else if summary, err := packet.Parse(); err != nil {
		return nil, err
	} else {
		return summary, nil
	}
}

func (r *ReplayParser) ParseReplay() error {
	r.readBytes(headerLength)
	for i := 0; i < 10; i++ {
		p, err := r.GetPacket()
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
		res, err := p.Parse()
		if err != nil {
			return err
		}
		if p.Size < 1000 {
			fmt.Println(pb.EDemoCommands(p.Kind))
			PrintStruct(res)

		} else {
			fmt.Println("- Too large to show :(")
		}
		fmt.Println()
	}
	return nil
}

func (r *ReplayParser) GetPacket() (*Packet, error) {
	if kind, err := r.readVarint32(); err != nil {
		return nil, err
	} else if tick, err := r.readVarint32(); err != nil {
		return nil, err
	} else if size, err := r.readVarint32(); err != nil {
		return nil, err
	} else if message, err := r.readBytes(size); err != nil {
		return nil, err
	} else {
		isCompressed := (kind & compressedMask) > 0
		packet := &Packet{
			Kind:         kind,
			Tick:         tick,
			Size:         size,
			IsCompressed: isCompressed,
			Message:      message,
		}
		if isCompressed {
			packet.Kind &= ^compressedMask
			packet.Message, err = snappy.Decode(nil, message)
			if err != nil {
				return nil, err
			}
		}
		return packet, nil
	}
}

func (r *ReplayParser) readVarint32() (int, error) {
	total, shift := 0, 0
	for {
		current, err := r.readByte()
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

func (r *ReplayParser) readUint32() (uint32, error) {
	data, err := r.readBytes(4)
	return binary.LittleEndian.Uint32(data), err
}

func (r *ReplayParser) readString(size int) (string, error) {
	data, err := r.readBytes(size)
	return string(data), err
}

func (r *ReplayParser) readByte() (byte, error) {
	data, err := r.readBytes(1)
	return data[0], err
}

func (r *ReplayParser) readBytes(size int) (data []byte, err error) {
	data = make([]byte, size)
	_, err = r.file.Read(data)
	return
}
