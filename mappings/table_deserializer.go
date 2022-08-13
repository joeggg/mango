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

type Serializer struct {
	SerializerId *SerializerId
	Fields       []*Field
}

type TableDeserializer struct {
	serializer  *pb.CSVCMsg_FlattenedSerializer
	Fields      []*Field
	Serializers map[string]*Serializer
}

func NewTableDeserializer(s *pb.CSVCMsg_FlattenedSerializer) *TableDeserializer {
	return &TableDeserializer{
		serializer:  s,
		Fields:      []*Field{},
		Serializers: make(map[string]*Serializer),
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
	for _, serializer := range d.serializer.GetSerializers() {
		s := &Serializer{
			SerializerId: &SerializerId{
				Name:    d.getSymbol(serializer.SerializerNameSym),
				Version: serializer.SerializerVersion,
			},
			Fields: []*Field{},
		}
		for _, fieldIndex := range serializer.GetFieldsIndex() {
			if int(fieldIndex) >= len(d.Fields) {
				continue
			}
			s.Fields = append(s.Fields, d.Fields[fieldIndex])
		}
		d.Serializers[*s.SerializerId.Name] = s
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
