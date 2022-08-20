package gatherers

import (
	"fmt"

	"github.com/joeggg/mango/embedded"
	"github.com/joeggg/mango/mappings"
	"github.com/joeggg/mango/pb"
	"google.golang.org/protobuf/proto"
)

type Player struct {
	AccountID       uint32
	Username        string
	Hero            string
	NetWorthGraph   []float32
	AbilityUpgrades []uint32
}

type Team struct {
	Players       map[string]*Player
	NetWorthGraph []float32
	XPGraph       []float32
	GoldGraph     []float32
}

type Metadata struct {
	MatchID uint64
	Radiant *Team
	Dire    *Team
}

type MetadataGatherer struct {
	handlers map[int]embedded.EmbeddedHandler
	data     *Metadata
}

func NewMetadataGatherer() embedded.Gatherer {
	mg := &MetadataGatherer{data: &Metadata{}}
	mg.handlers = map[int]embedded.EmbeddedHandler{
		int(pb.EDotaUserMessages_DOTA_UM_MatchMetadata): mg.handleData,
	}
	return mg
}

func (mg *MetadataGatherer) GetName() string { return "Metadata" }

func (mg *MetadataGatherer) GetHandlers() map[int]embedded.EmbeddedHandler {
	return mg.handlers
}

func (mg *MetadataGatherer) GetResults() interface{} {
	return mg.data
}

func (mg *MetadataGatherer) handleData(data proto.Message, lk *mappings.LookupObjects) error {
	match := data.(*pb.CDOTAMatchMetadataFile)
	mg.data.MatchID = match.GetMatchId()
	metadata := match.GetMetadata()

	for _, team := range metadata.GetTeams() {
		t := &Team{
			NetWorthGraph: team.GetGraphNetWorth(),
			XPGraph:       team.GetGraphExperience(),
			GoldGraph:     team.GetGraphGoldEarned(),
			Players:       make(map[string]*Player),
		}
		var isRadiant bool
		if num := team.GetDotaTeam(); num == 2 {
			isRadiant = true
			mg.data.Radiant = t
		} else if num == 3 {
			isRadiant = false
			mg.data.Dire = t
		} else {
			return fmt.Errorf("unknown team number: %d", num)
		}
		for i, player := range team.GetPlayers() {
			mg.addPlayerData(player, i, lk, isRadiant)
		}
	}

	return nil
}

func (mg *MetadataGatherer) addPlayerData(
	player *pb.CDOTAMatchMetadata_Team_Player,
	index int,
	lk *mappings.LookupObjects,
	isRadiant bool,
) error {
	info := &Player{}
	var playerSummary *pb.CGameInfo_CDotaGameInfo_CPlayerInfo
	if isRadiant {
		playerSummary = lk.Players[index]
	} else {
		playerSummary = lk.Players[5+index]
	}
	// Basic player data
	info.AccountID = player.GetAccountId()
	info.Username = playerSummary.GetPlayerName()
	info.Hero = playerSummary.GetHeroName()
	info.NetWorthGraph = player.GetGraphNetWorth()

	if isRadiant {
		mg.data.Radiant.Players[info.Hero] = info
	} else {
		mg.data.Dire.Players[info.Hero] = info
	}
	return nil
}
