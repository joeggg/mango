# Mango
A Dota 2 replay parser in Go

Inspired by Smoke and Clarity by Skadistats

A work in progress - can parse the entire replay, returning the unmarshalled protocol buffer packets. "Gatherers" can be attached to the parser to collect data from different kinds of packets. You can write your own or use the default ones provided. I currently don't have many implemented however. 

To be able to write these, check out the file `embedded/types.go` to see the different packet type names, which you can find the full definitions for in the files `pb/*.go`.

### Example Usage
```go
rp := mango.WithDefaultGatherers(mango.NewReplayParser(replayFilename))
err := rp.Initialise()
if err != nil {
    panic(err)
}
defer rp.Close()

packets, err := rp.ParseReplay()
if err != nil {
    panic(err)
}
results := rp.GetResults() // Results indexed by gatherer name

...
```
To add a gatherer
```go
rp.RegisterGatherer(myGatherer)
```
Gatherers must have the following methods (defined in `embedded/gatherer.go`):
```go
type Gatherer interface {
	GetName() string
	GetHandlers() map[int]EmbeddedHandler
	GetResults() interface{}
}
```
With `EmbeddedHandler` being a function that takes in a protobuf message (defined in `embedded/types.go`):
```go
type EmbeddedHandler func(proto.Message) error
```