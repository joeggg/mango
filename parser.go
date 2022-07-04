package mango

import (
	"encoding/binary"
	"errors"
	"os"
)

var (
	headerLength = 8
	header       = "PBDEMS2\x00"
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
		return errors.New("Failed to read header")
	}
	// Offset handling
	if _, err := r.file.Seek(int64(r.readUint32()), 1); err != nil {
		return err
	}
	return nil
}

func (r *ReplayParser) readUint32() uint32 {
	return binary.LittleEndian.Uint32(r.readBytes(4))
}

func (r *ReplayParser) readString(size int) string {
	return string(r.readBytes(size))
}

func (r *ReplayParser) readBytes(size int) (data []byte) {
	data = make([]byte, size)
	r.file.Read(data)
	return
}
