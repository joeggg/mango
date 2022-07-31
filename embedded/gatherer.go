package embedded

// Interface required for any gatherers added to the replay parser
type Gatherer interface {
	GetName() string
	GetHandlers() map[int]EmbeddedHandler
	GetResults() interface{}
}

type GathererFactory func() Gatherer
