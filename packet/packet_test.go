package packet_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/joeggg/mango"
	"github.com/joeggg/mango/mappings"
	"github.com/joeggg/mango/packet"
)

func TestPacketParse(t *testing.T) {
	packetMap := mango.LoadExamplePacketData("packets")
	for name, info := range packetMap {
		fmt.Printf("Testing packet type %s: \n\n", name)
		b, err := base64.StdEncoding.DecodeString(info["data"].(string))
		if err != nil {
			t.Error(err)
		}
		p := packet.Packet{RawMessage: b, Kind: int(info["code"].(float64))}
		err = p.Parse(&mappings.LookupObjects{})
		if err != nil {
			t.Error(err)
		}
		mango.PrintStruct(p.Message)
		if p.Embed != nil {
			fmt.Println("Embdedded data:")
			mango.PrintStruct(p.Embed.Data)
		}
		fmt.Println()
	}
}
