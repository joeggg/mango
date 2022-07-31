package gatherers

import "mango/embedded"

var Default = []embedded.GathererFactory{
	NewChatGatherer,
	NewMetadataGatherer,
}
