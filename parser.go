package mango

import (
	"encoding/binary"
	"errors"
	"fmt"
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
	file   *os.File
	offset int64
}

type Packet struct {
	Kind         int
	Tick         int
	Size         int
	IsCompressed bool
	Message      []byte
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
	if r.readString(headerLength) != header {
		return errors.New("failed to read header")
	}
	// Offset handling
	r.offset = int64(r.readUint32())
	fmt.Printf("Offset: %d\n", r.offset)
	return nil
}

func (r *ReplayParser) GetSummary() error {
	if _, err := r.file.Seek(r.offset, 0); err != nil {
		return err
	} else if packet, err := r.GetPacket(); err != nil {
		return err
	} else if summary, err := packet.Parse(); err != nil {
		return err
	} else {
		PrintStruct(summary)
		return nil
	}
}

func (p *Packet) Parse() (proto.Message, error) {
	if pb.EDemoCommands(p.Kind) == pb.EDemoCommands_DEM_FileInfo {
		data := &pb.CDemoFileInfo{}
		err := proto.Unmarshal(p.Message, data)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		return nil, errors.New("unknown protobuf type")
	}
}

func (r *ReplayParser) GetPacket() (*Packet, error) {
	if kind, err := r.readVarint32(); err != nil {
		return nil, err
	} else if tick, err := r.readVarint32(); err != nil {
		return nil, err
	} else if size, err := r.readVarint32(); err != nil {
		return nil, err
	} else {
		message := r.readBytes(size)
		isCompressed := (kind & compressedMask) > 0
		if isCompressed {
			kind &= ^compressedMask
			decompMessage, err := snappy.Decode(nil, message)
			if err != nil {
				return nil, err
			}
			return &Packet{kind, tick, size, isCompressed, decompMessage}, nil
		}
		return &Packet{kind, tick, size, isCompressed, message}, nil
	}
}

func (r *ReplayParser) readVarint32() (int, error) {
	total, shift := 0, 0
	for {
		current := int(r.readByte())
		total |= (current & 0x7F) << shift

		if current&0x80 == 0 {
			return total & varintMask, nil
		}

		shift += varintBlockSize
		if shift >= varintMaxSize {
			return total, errors.New("invalid varint")
		}
	}
}

func (r *ReplayParser) readUint32() uint32 {
	return binary.LittleEndian.Uint32(r.readBytes(4))
}

func (r *ReplayParser) readString(size int) string {
	return string(r.readBytes(size))
}

func (r *ReplayParser) readByte() byte {
	return r.readBytes(1)[0]
}

func (r *ReplayParser) readBytes(size int) (data []byte) {
	data = make([]byte, size)
	_, err := r.file.Read(data)
	if err != nil {
		panic(err)
	}
	return
}
