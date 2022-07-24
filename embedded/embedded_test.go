package embedded_test

import (
	"encoding/base64"
	"fmt"
	"mango"
	"mango/embedded"
	"testing"
)

func TestEmbeddedParse(t *testing.T) {
	packetMap := mango.LoadExamplePacketData("embedded_packets")
	for name, info := range packetMap {
		fmt.Printf("Testing packet type %s: \n\n", name)
		data, _ := base64.StdEncoding.DecodeString(info["data"].(string))
		d := embedded.NewEmbeddedDecoder(data)
		p, err := d.Decode()
		if err != nil {
			t.Error(err)
		}
		err = p.Parse()
		if err != nil {
			t.Error(err)
		}
		mango.PrintStruct(p.Data)
		fmt.Println()
	}
}
