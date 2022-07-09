package mango

import "fmt"

type EmbeddedParser struct {
	data []int8
}

func (p *EmbeddedParser) Parse() error {
	fmt.Println(p.data)
	kind := p.readUBitVar()
	fmt.Println(kind)
	return nil
}

func (p *EmbeddedParser) readUBitVar() int {
	return 0
}

func (p *EmbeddedParser) readUBitInt() int {
	return 0
}
