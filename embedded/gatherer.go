package embedded

// Interface required for any gatherers added to the replay parser
type Gatherer interface {
	GetHandlers() map[int]EmbeddedHandler
	GetResults() interface{}
}
