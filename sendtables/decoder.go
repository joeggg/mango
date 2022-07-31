package sendtables

import (
	"bytes"
	"errors"
	"io"
	"mango/pb"
	"math"

	"google.golang.org/protobuf/proto"
)

var (
	varintMask      = int(math.Pow(2, 32) - 1) // 32 bit mask
	varintBlockSize = 7
	varintMaxSize   = 32
)

type TableDecoder struct {
	buffer io.Reader
}

func (d *TableDecoder) Decode(data []byte) (*pb.CSVCMsg_FlattenedSerializer, error) {
	d.buffer = bytes.NewReader(data)
	result := &pb.CSVCMsg_FlattenedSerializer{}
	if size, err := d.readVarint32(); err != nil {
		return nil, err
	} else if message, err := d.readBytes(size); err != nil {
		return nil, err
	} else if err = proto.Unmarshal(message, result); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (d *TableDecoder) readVarint32() (int, error) {
	total, shift := 0, 0
	for {
		current, err := d.readByte()
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

func (d *TableDecoder) readByte() (byte, error) {
	data, err := d.readBytes(1)
	return data[0], err
}

func (d *TableDecoder) readBytes(size int) (data []byte, err error) {
	data = make([]byte, size)
	_, err = d.buffer.Read(data)
	return
}
