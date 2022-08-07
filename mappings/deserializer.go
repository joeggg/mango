package mappings

import (
	"github.com/joeggg/mango/pb"
)

type Field struct {
	Type *string
	Name *string
	*SerializerId
	*Properties
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

type Deserializer struct {
	serializer *pb.CSVCMsg_FlattenedSerializer
	fields     []*Field
}

func (d *Deserializer) CreateFields() {
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
		d.fields = append(d.fields, f)
	}
}

/*
	Lookup symbol from a reference
*/
func (d *Deserializer) getSymbol(num *int32) *string {
	if num == nil {
		return nil
	}
	return &d.serializer.GetSymbols()[*num]
}
