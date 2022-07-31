package gatherers

import (
	"mango/embedded"
	"mango/pb"

	"google.golang.org/protobuf/proto"
)

type Message struct {
	Text string
	Time float64
}

type ChatGatherer struct {
	handlers     map[int]embedded.EmbeddedHandler
	total        int
	timeOffset   float64
	seconds      float64
	GameMessages map[int][]*Message
}

func NewChatGatherer() *ChatGatherer {
	cg := &ChatGatherer{}
	cg.handlers = map[int]embedded.EmbeddedHandler{
		int(pb.EDotaUserMessages_DOTA_UM_ChatMessage):           cg.handleChatMessage,
		int(pb.EDotaUserMessages_DOTA_UM_ChatEvent):             cg.handleChatEvent,
		int(pb.EDotaUserMessages_DOTA_UM_ChatWheel):             cg.handleChatWheel,
		int(pb.NET_Messages_net_Tick):                           cg.handleTick,
		int(pb.EDotaUserMessages_DOTA_UM_GamerulesStateChanged): cg.handleGameRules,
	}
	cg.total = 0
	cg.timeOffset = 0
	cg.GameMessages = map[int][]*Message{}
	for i := -1; i < 10; i++ {
		cg.GameMessages[i] = []*Message{}
	}
	return cg
}

func (cg *ChatGatherer) GetHandlers() map[int]embedded.EmbeddedHandler {
	return cg.handlers
}

func (cg *ChatGatherer) GetResults() interface{} {
	return cg.GameMessages
}

func (cg *ChatGatherer) handleTick(data proto.Message) error {
	message := data.(*pb.CNETMsg_Tick)
	cg.seconds = float64(message.GetTick())*1.0/30.0 - cg.timeOffset
	return nil
}

func (cg *ChatGatherer) handleGameRules(data proto.Message) error {
	message := data.(*pb.CDOTAUserMsg_GamerulesStateChanged)
	if message.GetState() == 10 {
		cg.timeOffset = cg.seconds + 90
		cg.seconds = -90
	}
	return nil
}

func (cg *ChatGatherer) handleChatMessage(data proto.Message) error {
	message := data.(*pb.CDOTAUserMsg_ChatMessage)
	id := int(message.GetSourcePlayerId())
	cg.GameMessages[id] = append(
		cg.GameMessages[id], &Message{message.GetMessageText(), cg.seconds / 60},
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
			cg.GameMessages[int(id)] = append(
				cg.GameMessages[int(id)], &Message{message.GetType().String(), cg.seconds / 60},
			)
		}
	} else {
		cg.GameMessages[-1] = append(
			cg.GameMessages[-1], &Message{message.GetType().String(), cg.seconds / 60},
		)
	}
	cg.total++
	return nil
}

func (cg *ChatGatherer) handleChatWheel(data proto.Message) error {
	message := data.(*pb.CDOTAUserMsg_ChatWheel)
	chatWheel := pb.EDOTAChatWheelMessage(message.GetChatMessageId()).String()
	id := int(message.GetPlayerId())
	cg.GameMessages[id] = append(cg.GameMessages[id], &Message{chatWheel, cg.seconds / 60})
	cg.total++
	return nil
}
