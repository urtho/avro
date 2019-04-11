package sqlavro

import (
	"fmt"
	"strconv"
	"time"

	"github.com/khezen/avro"
)

func sqlColumn2AVRO(columnName string, dataType SQLType, isNullable bool, defaultValue []byte, numPrecision, numScale, charBytesLen int) (*avro.RecordFieldSchema, error) {
	var (
		fieldType avro.Schema
	)
	switch dataType {
	case Char, NChar:
		fieldType = &avro.FixedSchema{
			Name: columnName,
			Type: avro.TypeFixed,
			Size: charBytesLen,
		}
		break
	case VarChar, NVarChar,
		Text, TinyText, MediumText, LongText,
		Enum, Set:
		fieldType = avro.TypeString
		break
	case Blob, MediumBlob, LongBlob:
		fieldType = avro.TypeBytes
		break
	case TinyInt, SmallInt, MediumInt, Int, Year:
		fieldType = avro.TypeInt32
		break
	case BigInt:
		fieldType = avro.TypeInt64
		break
	case Float:
		fieldType = avro.TypeFloat32
		break
	case Double:
		fieldType = avro.TypeFloat64
		break
	case Decimal:
		fieldType = &avro.DerivedPrimitiveSchema{
			Type:        avro.TypeBytes,
			LogicalType: avro.LogicalTypeDecimal,
			Precision:   &numPrecision,
			Scale:       &numScale,
		}
		break
	case Date:
		fieldType = &avro.DerivedPrimitiveSchema{
			Type:        avro.TypeInt32,
			LogicalType: avro.LogicalTypeDate,
		}
		break
	case Time:
		fieldType = &avro.DerivedPrimitiveSchema{
			Type:        avro.TypeInt32,
			LogicalType: avro.LogicalTypeTime,
		}
		break
	case DateTime, Timestamp:
		fieldType = &avro.DerivedPrimitiveSchema{
			Type:        avro.TypeInt32,
			LogicalType: avro.LogicalTypeTimestamp,
		}
		break
	default:
		return nil, avro.ErrUnsupportedType
	}
	switch dataType {
	case Char, NChar, VarChar, NVarChar,
		Text, TinyText, MediumText, LongText,
		Enum, Set:
		if defaultValue != nil {
			defaultValue = []byte(fmt.Sprintf(`"%s"`, string(defaultValue)))
		}
		break
	case Date, Time, DateTime:
		if defaultValue != nil {
			var format string
			switch dataType {
			case Date:
				format = "2006-01-02"
				break
			case Time:
				format = "15:04:05"
				break
			case DateTime:
				format = "2006-01-02 15:04:05"
			}
			t, err := time.Parse(format, string(defaultValue))
			if err != nil {
				defaultValue = []byte(fmt.Sprintf(`"%s"`, string(defaultValue)))
			} else {
				defaultValue = []byte(strconv.Itoa(int(t.Unix())))
			}
		}
		break
	}
	if isNullable {
		fieldType = avro.UnionSchema([]avro.Schema{avro.TypeNull, fieldType})
	}
	return &avro.RecordFieldSchema{
		Name:    columnName,
		Type:    fieldType,
		Default: defaultValue,
	}, nil
}
