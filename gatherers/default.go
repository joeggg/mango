package gatherers

import "github.com/joeggg/mango/embedded"

var Default = []embedded.GathererFactory{
	NewChatGatherer,
	NewMetadataGatherer,
}
