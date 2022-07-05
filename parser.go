package mango

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
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
	if r.readString(headerLength) != header {
		return errors.New("failed to read header")
	}
	// Offset handling
	offset := r.readUint32()
	fmt.Printf("Offset: %d\n", offset)
	if _, err := r.file.Seek(int64(offset), 0); err != nil {
		return err
	}
	return nil
}

func (r *ReplayParser) GetPacket() error {
	kind, err := r.readVarint32()
	if err != nil {
		return err
	}

	isCompressed := (kind & compressedMask) > 0
	if isCompressed {
		kind &= ^compressedMask
	}
	fmt.Printf("Kind: %v\n", kind)
	fmt.Printf("comp: %v\n", isCompressed)

	tick, err := r.readVarint32()
	if err != nil {
		return err
	}
	fmt.Printf("Tick: %v\n", tick)

	size, err := r.readVarint32()
	if err != nil {
		return err
	}
	fmt.Printf("Size: %v\n", size)

	message := r.readBytes(size)
	fmt.Printf("Message: %v\n", message)
	return nil
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
