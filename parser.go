package mango

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"os"

	"github.com/golang/snappy"
	"github.com/joeggg/mango/embedded"
	"github.com/joeggg/mango/gatherers"
	"github.com/joeggg/mango/packet"
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
	file      *os.File
	decoder   *embedded.EmbeddedDecoder
	Gatherers map[string]embedded.Gatherer
}

func NewReplayParser() *ReplayParser {
	return &ReplayParser{
		decoder:   &embedded.EmbeddedDecoder{},
		Gatherers: map[string]embedded.Gatherer{},
	}
}

func WithDefaultGatherers(rp *ReplayParser) *ReplayParser {
	for _, gFactory := range gatherers.Default {
		rp.RegisterGatherer(gFactory())
	}
	return rp
}

func (rp *ReplayParser) RegisterGatherer(g embedded.Gatherer) {
	rp.Gatherers[g.GetName()] = g
}

func (rp *ReplayParser) Initialise(filename string) error {
	// Read file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	rp.file = file
	// Header handling
	if data, _ := rp.readString(headerLength); data != header {
		return errors.New("failed to read header")
	}
	return nil
}

func (rp *ReplayParser) Close() error {
	return rp.file.Close()
}

/*

 */
func (rp *ReplayParser) GetSummary() (proto.Message, error) {
	// Offset handling
	if offset, err := rp.readUint32(); err != nil {
		return nil, err
	} else if _, err = rp.file.Seek(int64(offset), 0); err != nil {
		return nil, err
	} else if packet, err := rp.getPacket(); err != nil {
		return nil, err
	} else if err := packet.Parse(); err != nil {
		return nil, err
	} else {
		return packet.Message, nil
	}
}

/*
	Parse through the entire replay
*/
func (rp *ReplayParser) ParseReplay() ([]*packet.Packet, error) {
	var packets []*packet.Packet
	rp.readBytes(headerLength) // Read past summary offset
	for {
		// Get next packet and parse
		p, err := rp.getPacket()
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
		// Handle embedded message
		if p.RawEmbed != nil {
			embed, err := rp.decoder.Decode(p.RawEmbed)
			if err != nil {
				return packets, err
			}
			embed.Parse(rp.Gatherers)
			p.Embed = embed
		}
		packets = append(packets, p)
	}
}

/*
	Return the results of all gatherers in one object, indexed by their names
*/
func (rp *ReplayParser) GetResults() map[string]interface{} {
	results := map[string]interface{}{}
	for name, g := range rp.Gatherers {
		results[name] = g.GetResults()
	}
	return results
}

func (rp *ReplayParser) getPacket() (*packet.Packet, error) {
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
		packet := &packet.Packet{
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
