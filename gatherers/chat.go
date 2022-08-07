package gatherers

import (
	"github.com/joeggg/mango/embedded"
	"github.com/joeggg/mango/mappings"
	"github.com/joeggg/mango/pb"
	"google.golang.org/protobuf/proto"
)

type MessageSet struct {
	Player   *pb.CGameInfo_CDotaGameInfo_CPlayerInfo
	Messages []*Message
}

type Message struct {
	Text string
	Time float64
}

type ChatGatherer struct {
	handlers     map[int]embedded.EmbeddedHandler
	total        int
	timeOffset   float64
	seconds      float64
	GameMessages map[int]*MessageSet
}

func NewChatGatherer() embedded.Gatherer {
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
	cg.GameMessages = make(map[int]*MessageSet)
	for i := -1; i < 10; i++ {
		cg.GameMessages[i] = &MessageSet{nil, []*Message{}}
	}
	return cg
}

func (cg *ChatGatherer) GetName() string { return "Chat" }

func (cg *ChatGatherer) GetHandlers() map[int]embedded.EmbeddedHandler {
	return cg.handlers
}

func (cg *ChatGatherer) GetResults() interface{} {
	return cg.GameMessages
}

func (cg *ChatGatherer) handleTick(data proto.Message, lk *mappings.LookupObjects) error {
	message := data.(*pb.CNETMsg_Tick)
	cg.seconds = float64(message.GetTick())/30.0 - cg.timeOffset
	return nil
}

func (cg *ChatGatherer) handleGameRules(data proto.Message, lk *mappings.LookupObjects) error {
	message := data.(*pb.CDOTAUserMsg_GamerulesStateChanged)
	if message.GetState() == 10 {
		cg.timeOffset = cg.seconds
		cg.seconds = -90
	}
	return nil
}

func (cg *ChatGatherer) handleChatMessage(data proto.Message, lk *mappings.LookupObjects) error {
	message := data.(*pb.CDOTAUserMsg_ChatMessage)
	id := int(message.GetSourcePlayerId())
	player, ok := lk.Players[id]
	if !ok {
		return nil
	}
	cg.GameMessages[id].Player = player
	cg.GameMessages[id].Messages = append(cg.GameMessages[id].Messages, &Message{message.GetMessageText(), cg.seconds / 60})
	cg.total++
	return nil
}

func (cg *ChatGatherer) handleChatEvent(data proto.Message, lk *mappings.LookupObjects) error {
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
			messageSet, ok := cg.GameMessages[int(id)]
			if !ok {
				continue
			}
			messageSet.Messages = append(
				messageSet.Messages,
				&Message{message.GetType().String(), cg.seconds / 60},
			)
		}
	} else {
		cg.GameMessages[-1].Messages = append(
			cg.GameMessages[-1].Messages, &Message{message.GetType().String(), cg.seconds / 60},
		)
	}
	cg.total++
	return nil
}

func (cg *ChatGatherer) handleChatWheel(data proto.Message, lk *mappings.LookupObjects) error {
	message := data.(*pb.CDOTAUserMsg_ChatWheel)
	chatWheel := pb.EDOTAChatWheelMessage(message.GetChatMessageId()).String()
	id := int(message.GetPlayerId())
	player, ok := lk.Players[id]
	if !ok {
		return nil
	}
	cg.GameMessages[id].Player = player
	cg.GameMessages[id].Messages = append(
		cg.GameMessages[id].Messages, &Message{chatWheel, cg.seconds / 60},
	)
	cg.total++
	return nil
}
