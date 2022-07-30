package gatherers

import (
	"fmt"
	"mango/embedded"
	"mango/pb"

	"google.golang.org/protobuf/proto"
)

type ChatGatherer struct {
	handlers map[int]embedded.EmbeddedHandler
	total    int
}

func NewChatGatherer() *ChatGatherer {
	cg := &ChatGatherer{}
	cg.handlers = map[int]embedded.EmbeddedHandler{
		int(pb.EDotaUserMessages_DOTA_UM_ChatMessage): cg.handleChatMessage,
		int(pb.EDotaUserMessages_DOTA_UM_ChatEvent):   cg.handleChatEvent,
		int(pb.EDotaUserMessages_DOTA_UM_ChatWheel):   cg.handleChatWheel,
	}
	cg.total = 0
	return cg
}

func (cg *ChatGatherer) GetHandlers() map[int]embedded.EmbeddedHandler {
	return cg.handlers
}

func (cg *ChatGatherer) GetResults() interface{} {
	return cg.total
}

func (cg *ChatGatherer) handleChatMessage(data proto.Message) error {
	message := data.(*pb.CDOTAUserMsg_ChatMessage)
	fmt.Printf("ChatMessage: %v\n", message)
	cg.total++
	return nil
}

func (cg *ChatGatherer) handleChatEvent(data proto.Message) error {
	message := data.(*pb.CDOTAUserMsg_ChatEvent)
	fmt.Printf("ChatEvent: %v\n", message)
	cg.total++
	return nil
}

func (cg *ChatGatherer) handleChatWheel(data proto.Message) error {
	message := data.(*pb.CDOTAUserMsg_ChatWheel)
	fmt.Printf("ChatWheel: %v\n", message)
	cg.total++
	return nil
}
