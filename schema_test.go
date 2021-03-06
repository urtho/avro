package avro

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestMarshaling(t *testing.T) {
	cases := []struct {
		typeName    Type
		schemaBytes []byte
		expectedErr error
	}{
		{
			TypeRecord,
			[]byte(`{"type":"record","name":"person","fields":[{"name":"id","type":"int"},{"name":"guid","type":"string"},{"name":"isActive","type":"boolean"},{"name":"age","type":"int"},{"name":"name","type":"string"},{"name":"address","type":"string"},{"name":"latitude","type":"double"},{"name":"longitude","type":"double"},{"name":"tags","type":{"type":"array","items":"string"}},{"name":"friends","type":{"type":"array","items":{"type":"record","name":"friends_record","fields":[{"name":"id","type":"int"},{"name":"name","type":"string"}]}}},{"name":"randomArrayItem","type":"string"}]}`),
			nil,
		},
		{
			TypeRecord,
			[]byte(`{"type":"record","namespace":"test","name":"LongList","aliases":["LinkedLongs"],"doc":"list of 64 bits integers","fields":[{"name":"value","type":"long"},{"name":"next","type":["null","LongList"]}]}`),
			nil,
		},
		{
			TypeRecord,
			[]byte(`{"type":"record","name":"LongList","aliases":["LinkedLongs"],"fields":[{"name":"value","type":"long"},{"name":"next","type":["null","LongList"]}]}`),
			nil,
		},
		{
			TypeRecord,
			[]byte(`{"type":"record","name":"LongList","aliases":["LinkedLongs"],"fields":[{"name":"value","type":"long"},{"name":"next","type":["null","LongList"]}]}`),
			nil,
		},
		{
			TypeRecord,
			[]byte(`{"type":"record","name":"LongList","aliases":["LinkedLongs"],"fields":[{"name":"value","type":"long","order":"ignore"},{"name":"next","type":["null","LongList"]}]}`),
			nil,
		},
		{
			TypeRecord,
			[]byte(`{"type":"record","name":"LongList","aliases":["LinkedLongs"],"fields":[{"name":"value","type":"long","default":0},{"name":"next","type":["null","LongList"]}]}`),
			nil,
		},
		{
			TypeRecord,
			[]byte(`{"type":"record","namespace":"test","name":"LongList","aliases":["LinkedLongs"],"fields":[{"name":"value","type":"something"},{"name":"next","type":["null","LongList"]}]}`),
			ErrUnsupportedType,
		},
		{
			TypeRecord,
			[]byte(`{"type":"record","namespace":"test","name":"LongList","aliases":["LinkedLongs"],"fields":[{"name":"value"},{"name":"next","type":["null","LongList"]}]}`),
			ErrInvalidSchema,
		},
		{
			TypeRecord,
			[]byte(`{"type":"record","name":"LongList","aliases":["LinkedLongs"],"fields":[{"name":"value","type":"long","order":"something"},{"name":"next","type":["null","LongList"]}]}`),
			ErrInvalidSchema,
		},
		{
			TypeRecord,
			[]byte(`{"type":"record","name":"LongList","aliases":["LinkedLongs"],"fields":[{"name":"value","type":"long","order":0},{"name":"next","type":["null","LongList"]}]}`),
			ErrInvalidSchema,
		},
		{
			TypeRecord,
			[]byte(`{"type":"record","fields":[{"name":"value","type":"long"},{"name":"next","type":["null","LongList"]}]}`),
			ErrInvalidSchema,
		},
		{
			TypeRecord,
			[]byte(`{"type":"record","name":"LongList","fields":[{"type":"long"},{"name":"next","type":["null","LongList"]}]}`),
			ErrInvalidSchema,
		},
		{
			TypeRecord,
			[]byte(`{"type":"record","name":"LongList","aliases":"something","fields":[{"name":"value","type":"long"},{"name":"next","type":["null","LongList"]}]}`),
			ErrInvalidSchema,
		},
		{
			TypeRecord,
			[]byte(`{"type":"record","name":"LongList","fields":[{"name":"value","aliases":"something","type":"long"},{"name":"next","type":["null","LongList"]}]}`),
			ErrInvalidSchema,
		},
		{
			TypeRecord,
			[]byte(`{"type":"record","name":"LongList"}`),
			ErrInvalidSchema,
		},
		{
			TypeRecord,
			[]byte(`{"type":"record","name":"LongList","fields":"something"}`),
			ErrInvalidSchema,
		},
		{
			TypeRecord,
			[]byte(`{"type":"record","name":"LongList","aliases":[0],"fields":[{"name":"value","type":"long"},{"name":"next","type":["null","LongList"]}]}`),
			ErrInvalidSchema,
		},
		{
			TypeArray,
			[]byte(`{"type":"array","items":"string"}`),
			nil,
		},
		{
			TypeArray,
			[]byte(`{"type":"array","items":["null","string"]}`),
			nil,
		},
		{
			TypeArray,
			[]byte(`{"type":"array","values":"long"}`),
			ErrInvalidSchema,
		},
		{
			TypeArray,
			[]byte(`{"type":"array","items":"something"}`),
			ErrUnsupportedType,
		},
		{
			TypeMap,
			[]byte(`{"type":"map","values":"long"}`),
			nil,
		},
		{
			TypeMap,
			[]byte(`{"type":"map","values":["null","long"]}`),
			nil,
		},
		{
			TypeMap,
			[]byte(`{"type":"map","values":["null","something"]}`),
			ErrUnsupportedType,
		},
		{
			TypeMap,
			[]byte(`{"type":"map","items":"long"}`),
			ErrInvalidSchema,
		},
		{
			TypeEnum,
			[]byte(`{"type":"enum","name":"Suit","symbols":["SPADES","HEARTS","DIAMONDS","CLUBS"]}`),
			nil,
		},
		{
			TypeEnum,
			[]byte(`{"type":"enum","name":"Suit"}`),
			ErrInvalidSchema,
		},
		{
			TypeEnum,
			[]byte(`{"type":"enum","name":"Suit","symbols":"something"}`),
			ErrInvalidSchema,
		},
		{
			TypeEnum,
			[]byte(`{"type":"enum","name":0,"symbols":["SPADES"]}`),
			ErrInvalidSchema,
		},
		{
			TypeEnum,
			[]byte(`{"type":"enum","name":"Suit","symbols":["SPADES",11,"DIAMONDS","CLUBS"]}`),
			ErrInvalidSchema,
		},
		{
			TypeEnum,
			[]byte(`{"type":"enum","name":"Suit","namespace":0,"symbols":["SPADES"]}`),
			ErrInvalidSchema,
		},
		{
			TypeEnum,
			[]byte(`{"type":"enum","name":"Suit","doc":0,"symbols":["SPADES"]}`),
			ErrInvalidSchema,
		},
		{
			TypeFixed,
			[]byte(`{"type":"fixed","name":"md5","size":16}`),
			nil,
		},
		{
			TypeFixed,
			[]byte(`{"type":"fixed","name":"md5","size":-16}`),
			ErrInvalidSchema,
		},
		{
			TypeFixed,
			[]byte(`{"type":"fixed","name":"md5"}`),
			ErrInvalidSchema,
		},
		{
			TypeFixed,
			[]byte(`{"type":"fixed","name":"md5","size":"16"}`),
			ErrInvalidSchema,
		},
		{
			TypeFixed,
			[]byte(`{"type":"fixed","name":0,"size":16}`),
			ErrInvalidSchema,
		},
		{
			TypeFixed,
			[]byte(`{"type":"fixed","name":"md5","size":16}`),
			nil,
		},
		{
			Type(LogialTypeDuration),
			[]byte(`{"type":"fixed","logicalType":"duration","name":"md5","size":12}`),
			nil,
		},
		{
			Type(LogialTypeDuration),
			[]byte(`{"type":"fixed","logicalType":"duration","name":"md5","size":16}`),
			ErrInvalidSchema,
		},
		{
			Type(LogicalTypeTimestamp),
			[]byte(`{"type":"fixed","logicalType":"timestamp","name":"md5","size":12}`),
			ErrInvalidSchema,
		},
		{
			TypeUnion,
			[]byte(`["null","string"]`),
			nil,
		},
		{
			TypeUnion,
			[]byte(`["something","string"]`),
			ErrUnsupportedType,
		},
		{
			TypeUnion,
			[]byte(`[0,"string"]`),
			ErrInvalidSchema,
		},
		{
			Type("something"),
			[]byte(`{"type":"something","name":"something"}`),
			ErrUnsupportedType,
		},
		{
			Type(LogicalTypeTimestamp),
			[]byte(`{"type":"int","logicalType":"timestamp"}`),
			nil,
		},
		{
			Type(LogicalTypeTime),
			[]byte(`{"type":"int","logicalType":"time"}`),
			nil,
		},
		{
			Type(LogicalTypeDate),
			[]byte(`{"type":"int","logicalType":"date"}`),
			nil,
		},
		{
			TypeInt32,
			[]byte(`{"type":"int","logicalType":0}`),
			ErrInvalidSchema,
		},
		{
			Type(LogicalTypeDecimal),
			[]byte(`{"type":"int","logicalType":"decimal"}`),
			ErrInvalidSchema,
		},
		{
			Type(LogicalTypeTimestamp),
			[]byte(`{"type":"long","logicalType":"timestamp"}`),
			nil,
		},
		{
			Type(LogicalTypeTime),
			[]byte(`{"type":"long","logicalType":"time"}`),
			nil,
		},
		{
			Type(LogicalTypeDate),
			[]byte(`{"type":"long","logicalType":"date"}`),
			nil,
		},
		{
			TypeInt64,
			[]byte(`{"type":"long","logicalType":0}`),
			ErrInvalidSchema,
		},
		{
			Type(LogicalTypeDecimal),
			[]byte(`{"type":"long","logicalType":"decimal"}`),
			ErrInvalidSchema,
		},
		{
			Type(LogicalTypeTimestamp),
			[]byte(`{"type":"int","logicalType":"timestamp"}`),
			nil,
		},
		{
			Type(LogicalTypeTime),
			[]byte(`{"type":"int","logicalType":"time"}`),
			nil,
		},
		{
			Type(LogicalTypeDate),
			[]byte(`{"type":"int","logicalType":"date"}`),
			nil,
		},
		{
			TypeInt32,
			[]byte(`{"type":"int","logicalType":0}`),
			ErrInvalidSchema,
		},
		{
			Type(LogicalTypeTimestamp),
			[]byte(`{"type":"int","logicalType":"timestamp"}`),
			nil,
		},
		{
			Type(LogicalTypeTime),
			[]byte(`{"type":"int","logicalType":"time"}`),
			nil,
		},
		{
			Type(LogicalTypeDate),
			[]byte(`{"type":"int","logicalType":"date"}`),
			nil,
		},
		{
			TypeInt32,
			[]byte(`{"type":"int","logicalType":0}`),
			ErrInvalidSchema,
		},
		{
			TypeBytes,
			[]byte(`{"type":"bytes"}`),
			ErrInvalidSchema,
		},
		{
			Type(LogicalTypeTime),
			[]byte(`{"type":"bytes","logicalType":"time"}`),
			ErrInvalidSchema,
		},
		{
			Type(LogicalTypeDecimal),
			[]byte(`{"type":"bytes","logicalType":"decimal"}`),
			ErrInvalidSchema,
		},
		{
			Type(LogicalTypeDecimal),
			[]byte(`{"type":"bytes","logicalType":"decimal","precision":5}`),
			nil,
		},
		{
			Type(LogicalTypeDecimal),
			[]byte(`{"type":"bytes","logicalType":"decimal","precision":-5}`),
			ErrInvalidSchema,
		},
		{
			Type(LogicalTypeDecimal),
			[]byte(`{"type":"bytes","logicalType":"decimal","precision":"something"}`),
			ErrInvalidSchema,
		},
		{
			Type(LogicalTypeDecimal),
			[]byte(`{"type":"bytes","logicalType":"decimal","precision":5,"scale":2}`),
			nil,
		},
		{
			Type(LogicalTypeDecimal),
			[]byte(`{"type":"bytes","logicalType":"decimal","precision":5,"scale":-2}`),
			ErrInvalidSchema,
		},
		{
			Type(LogicalTypeDecimal),
			[]byte(`{"type":"bytes","logicalType":"decimal","precision":5,"scale":"something"}`),
			ErrInvalidSchema,
		},
	}
	var (
		anySchema        AnySchema
		underlyingSchema Schema
		schemaBytes      []byte
	)
	for i, c := range cases {
		err := json.Unmarshal(c.schemaBytes, &anySchema)
		if err != nil && err != c.expectedErr {
			t.Errorf("case %d - error %v, got %v", i, c.expectedErr, err)
		}
		if err != nil {
			continue
		}
		underlyingSchema = anySchema.Schema()
		if underlyingSchema.TypeName() != c.typeName {
			t.Errorf("case %d - expected:%s got:%s", i, c.typeName, underlyingSchema.TypeName())
		}
		schemaBytes, err = json.Marshal(underlyingSchema)
		if err != nil {
			t.Errorf("case %d - error %v, got %v", i, c.expectedErr, err)
			continue
		}
		if !bytes.EqualFold(schemaBytes, c.schemaBytes) {
			t.Errorf("case %d -\nexpected:\n%s\ngot:\n%s\n", i, c.schemaBytes, schemaBytes)
		}
	}
}
