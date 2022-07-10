package mango_test

import (
	"mango"
	"testing"
)

func TestEmbeddedParser(t *testing.T) {
	data := []int8{1, 2, 3, 4}
	p := mango.NewEmbeddedParser(data)
	err := p.Parse()
	if err != nil {
		t.Error(err)
	}
}
