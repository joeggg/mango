package embedded_test

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"testing"

	"github.com/joeggg/mango"
	"github.com/joeggg/mango/embedded"
	"github.com/joeggg/mango/mappings"
	"github.com/joeggg/mango/pb"
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

func TestGatherersSuccess(t *testing.T) {
	packetMap := mango.LoadExamplePacketData("embedded_packets")
	info := packetMap["ServerInfo"]
	total := 3
	// Register multiple gatherers
	gatherers := map[string]embedded.Gatherer{}
	for i := 0; i < total; i++ {
		gatherers[strconv.Itoa(i)] = NewTestGatherer()
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
		fmt.Printf("Handler %s ran %d times\n", i, g.GetResults())
		count += g.GetResults().(int)
	}
	// Check all ran each time and the data accumulated
	if count != total*total {
		t.Error("Not all handlers ran!")
	}
}

func TestGatherersFirstNull(t *testing.T) {
	packetMap := mango.LoadExamplePacketData("embedded_packets")
	info := packetMap["ServerInfo"]
	gatherers := map[string]embedded.Gatherer{
		"0": &TestGatherer{map[int]embedded.EmbeddedHandler{}, 0}, "1": NewTestGatherer(),
	}
	_, err := processEmbeddedFromJson(info, gatherers)
	if err != nil {
		t.Error(err)
	}
	count := 0
	for i, g := range gatherers {
		fmt.Printf("Handler %s ran %d times\n", i, g.GetResults())
		count += g.GetResults().(int)
	}
	// Only second one should run
	if count != 1 {
		t.Errorf("invalid number of handlers ran: %d", count)
	}
}

func processEmbeddedFromJson(
	info map[string]interface{}, gatherers map[string]embedded.Gatherer,
) (*embedded.EmbeddedPacket, error) {
	data, _ := base64.StdEncoding.DecodeString(info["data"].(string))
	d := embedded.EmbeddedDecoder{}
	p, err := d.Decode(data)
	if err != nil {
		return nil, err
	}
	err = p.Parse(gatherers, &mappings.LookupObjects{})
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

func (t *TestGatherer) GetName() string { return "Test" }

func (t *TestGatherer) GetHandlers() map[int]embedded.EmbeddedHandler {
	return t.handlers
}

func (t *TestGatherer) GetResults() interface{} {
	return t.result
}

func (t *TestGatherer) testHandler(data proto.Message, lk *mappings.LookupObjects) error {
	t.result++
	return nil
}
