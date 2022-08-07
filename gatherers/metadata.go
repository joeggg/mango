package gatherers

import (
	"github.com/joeggg/mango/embedded"
	"github.com/joeggg/mango/mappings"
	"github.com/joeggg/mango/pb"
	"google.golang.org/protobuf/proto"
)

type MetadataGatherer struct {
	handlers map[int]embedded.EmbeddedHandler
	data     *pb.CDOTAMatchMetadataFile
}

func NewMetadataGatherer() embedded.Gatherer {
	mg := &MetadataGatherer{}
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
	mg.data = data.(*pb.CDOTAMatchMetadataFile)
	return nil
}
