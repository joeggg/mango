package mappings

import (
	"github.com/joeggg/mango/pb"
)

type Field struct {
	Type         *string
	Name         *string
	SerializerId *SerializerId
	Properties   *Properties
}

type SerializerId struct {
	Name    *string
	Version *int32
}

type Properties struct {
	EncodeFlags *int32
	BitCount    *int32
	LowValue    *float32
	HighValue   *float32
	EncoderType *string
}

type TableDeserializer struct {
	serializer *pb.CSVCMsg_FlattenedSerializer
	Fields     []*Field
}

func NewTableDeserializer(s *pb.CSVCMsg_FlattenedSerializer) *TableDeserializer {
	return &TableDeserializer{
		serializer: s,
		Fields:     []*Field{},
	}
}

func (d *TableDeserializer) CreateFields() {
	for _, field := range d.serializer.GetFields() {
		f := &Field{
			Type: d.getSymbol(field.VarTypeSym),
			Name: d.getSymbol(field.VarNameSym),
			SerializerId: &SerializerId{
				Name:    d.getSymbol(field.FieldSerializerNameSym),
				Version: field.FieldSerializerVersion,
			},
			Properties: &Properties{
				EncodeFlags: field.EncodeFlags,
				BitCount:    field.BitCount,
				LowValue:    field.LowValue,
				HighValue:   field.HighValue,
				EncoderType: d.getSymbol(field.VarEncoderSym),
			},
		}
		d.Fields = append(d.Fields, f)
	}
}

/*
	Lookup symbol from a reference
*/
func (d *TableDeserializer) getSymbol(num *int32) *string {
	if num == nil {
		return nil
	}
	return &d.serializer.GetSymbols()[*num]
}
