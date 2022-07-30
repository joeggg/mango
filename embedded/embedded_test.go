package embedded_test

import (
	"encoding/base64"
	"fmt"
	"mango"
	"mango/embedded"
	"mango/pb"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestEmbeddedParse(t *testing.T) {
	packetMap := mango.LoadExamplePacketData("embedded_packets")
	for name, info := range packetMap {
		fmt.Printf("Testing packet type %s: \n\n", name)
		p, err := processEmbeddedFromJson(info, nil)
		if err != nil {
			t.Error(err)
		}
		mango.PrintStruct(p.Data)
		fmt.Println()
	}
}

func TestGatherers(t *testing.T) {
	packetMap := mango.LoadExamplePacketData("embedded_packets")
	info := packetMap["ServerInfo"]
	total := 3
	// Register multiple gatherers
	var gatherers []embedded.Gatherer
	for i := 0; i < total; i++ {
		gatherers = append(gatherers, NewTestGatherer())
	}
	// Parse several times to simulate data gatherering over time
	for i := 0; i < total; i++ {
		_, err := processEmbeddedFromJson(info, gatherers)
		if err != nil {
			t.Error(err)
		}
	}

	count := 0
	for i, g := range gatherers {
		fmt.Printf("Handler %d ran %d times\n", i, g.GetResults())
		count += g.GetResults().(int)
	}
	// Check all ran each time
	if count != total*total {
		t.Error("Not all handlers ran!")
	}
}

func processEmbeddedFromJson(
	info map[string]interface{}, gatherers []embedded.Gatherer,
) (*embedded.EmbeddedPacket, error) {
	data, _ := base64.StdEncoding.DecodeString(info["data"].(string))
	d := embedded.EmbeddedDecoder{}
	p, err := d.Decode(data)
	if err != nil {
		return nil, err
	}
	err = p.Parse(gatherers)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// Test gatherer to count up each time a packet is handled
type TestGatherer struct {
	handlers map[int]embedded.EmbeddedHandler
	result   int
}

func NewTestGatherer() *TestGatherer {
	t := &TestGatherer{}
	t.handlers = map[int]embedded.EmbeddedHandler{
		int(pb.SVC_Messages_svc_ServerInfo): t.testHandler,
	}
	t.result = 0
	return t
}

func (t *TestGatherer) GetHandlers() map[int]embedded.EmbeddedHandler {
	return t.handlers
}

func (t *TestGatherer) GetResults() interface{} {
	return t.result
}

func (t *TestGatherer) testHandler(data proto.Message) error {
	t.result++
	return nil
}
