package mango_test

import (
	"fmt"
	"mango"
	"testing"

	"google.golang.org/protobuf/proto"
)

const testReplayFilename = "testdata/test.dem"

func TestSummary(t *testing.T) {
	p := mango.NewReplayParser(testReplayFilename)
	err := p.Initialise()
	if err != nil {
		t.Error(err)
	}
	defer p.Close()
	if summary, err := p.GetSummary(); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("Summary packet:\n\n")
		mango.PrintStruct(summary)
	}
}

func TestParse(t *testing.T) {
	p := mango.NewReplayParser(testReplayFilename)
	err := p.Initialise()
	if err != nil {
		t.Error(err)
	}
	defer p.Close()

	packets, err := p.ParseReplay()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("All replay parsed through without errors!")
	fmt.Printf("Sample packets: \n\n")
	count := 0
	for i := 0; i < 10; i++ {
		var toShow proto.Message
		if packets[i].Embed != nil {
			fmt.Printf("%s:\n", packets[i].Command)
			fmt.Printf("Embedded packet! Type: %s\n", packets[i].Embed.Command)
			toShow = packets[i].Embed.Data
			count++
		} else {
			fmt.Printf("%s:\n", packets[i].Command)
			toShow = packets[i].Message
			count++
		}

		if toShow != nil {
			mango.PrintStruct(toShow)
			fmt.Println()
		}
	}
	fmt.Printf("Showing %d packets\n", count)
}

func TestParseWithGatherers(t *testing.T) {
	rp := mango.WithDefaultGatherers(mango.NewReplayParser(testReplayFilename))
	err := rp.Initialise()
	if err != nil {
		t.Error(err)
	}
	defer rp.Close()

	if _, err := rp.ParseReplay(); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("\nGatherer results: \n\n")
		results := rp.GetResults()
		for name := range results {
			if results == nil {
				t.Errorf("%s gatherer results were null", name)
			}
		}
		mango.PrintStruct(results)
	}
}
