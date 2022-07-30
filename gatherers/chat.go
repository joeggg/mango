package gatherers

import (
	"fmt"
	"mango/embedded"
	"mango/pb"

	"google.golang.org/protobuf/proto"
)

type ChatGatherer struct {
	handlers     map[int]embedded.EmbeddedHandler
	total        int
	GameMessages map[int][]string
}

func NewChatGatherer() *ChatGatherer {
	cg := &ChatGatherer{}
	cg.handlers = map[int]embedded.EmbeddedHandler{
		int(pb.EDotaUserMessages_DOTA_UM_ChatMessage): cg.handleChatMessage,
		int(pb.EDotaUserMessages_DOTA_UM_ChatEvent):   cg.handleChatEvent,
		int(pb.EDotaUserMessages_DOTA_UM_ChatWheel):   cg.handleChatWheel,
	}
	cg.total = 0
	cg.GameMessages = map[int][]string{}
	for i := -1; i < 10; i++ {
		cg.GameMessages[i] = []string{}
	}
	return cg
}

func (cg *ChatGatherer) GetHandlers() map[int]embedded.EmbeddedHandler {
	return cg.handlers
}

func (cg *ChatGatherer) GetResults() interface{} {
	return cg.GameMessages
}

func (cg *ChatGatherer) handleChatMessage(data proto.Message) error {
	message := data.(*pb.CDOTAUserMsg_ChatMessage)
	id := int(message.GetSourcePlayerId())
	cg.GameMessages[id] = append(cg.GameMessages[id], message.GetMessageText())
	fmt.Printf(
		"ChatMessage:\n\tplayer=%d, message=%s\n\n",
		message.GetSourcePlayerId(),
		message.GetMessageText(),
	)
	cg.total++
	return nil
}

func (cg *ChatGatherer) handleChatEvent(data proto.Message) error {
	message := data.(*pb.CDOTAUserMsg_ChatEvent)
	// Get all potential player IDs and remove unnused ones (== -1)
	playerIds := []int32{
		message.GetPlayerid_1(),
		message.GetPlayerid_2(),
		message.GetPlayerid_3(),
		message.GetPlayerid_4(),
		message.GetPlayerid_5(),
		message.GetPlayerid_6(),
	}
	for i := 0; i < len(playerIds); i++ {
		if playerIds[i] == -1 {
			playerIds = append(playerIds[:i], playerIds[i+1:]...)
			i--
		}
	}
	// Add message to store
	if len(playerIds) > 0 {
		for _, id := range playerIds {
			cg.GameMessages[int(id)] = append(cg.GameMessages[int(id)], message.GetType().String())
		}
	} else {
		cg.GameMessages[-1] = append(cg.GameMessages[-1], message.GetType().String())
	}
	fmt.Printf(
		"ChatEvent:\n\tplayers=%v, type=%s\n\n",
		playerIds,
		message.GetType(),
	)
	cg.total++
	return nil
}

func (cg *ChatGatherer) handleChatWheel(data proto.Message) error {
	message := data.(*pb.CDOTAUserMsg_ChatWheel)
	chatWheel := pb.EDOTAChatWheelMessage(message.GetChatMessageId()).String()
	id := int(message.GetPlayerId())
	cg.GameMessages[id] = append(cg.GameMessages[id], chatWheel)
	fmt.Printf(
		"ChatWheel:\n\tplayer=%d, message=%s, emoticon=%d\n\n",
		message.GetPlayerId(),
		chatWheel,
		message.GetEmoticonId(),
	)
	cg.total++
	return nil
}
