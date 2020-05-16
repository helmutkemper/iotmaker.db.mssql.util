package iotmaker_db_mssql_util

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type ColsTypes struct {
	Name                 string
	DatabaseTypeName     string
	DecimalSizePrecision int64
	DecimalSizeScale     int64
	DecimalSizeOkToUse   bool
	Length               int64
	LengthOkToUse        bool
	Nullable             bool
	NullableOkToUse      bool
	ScanType             reflect.Type
}

func ToStruct(name string, data []ColsTypes) {
	firstChar := strings.Title(string(name[0]))
	name = firstChar + name[1:]
	fmt.Printf("type %v struct {\n", name)
	for _, dataCol := range data {
		firstChar := strings.Title(string(dataCol.Name[0]))
		name = firstChar + dataCol.Name[1:]

		fmt.Printf("  %v  %v\n", name, dataCol.ScanType.String())
	}
	fmt.Print("}\n\n")
}

func ListColsTypes(db *sql.DB, ctx context.Context, tableName string) (error, []ColsTypes) {
	var err error
	var queryReturn *sql.Rows
	var ret = make([]ColsTypes, 0)
	var line []*sql.ColumnType

	queryReturn, err = db.QueryContext(ctx, fmt.Sprintf("SELECT * FROM %v;", tableName))
	if err != nil {
		return err, nil
	}
	defer queryReturn.Close()

	line, err = queryReturn.ColumnTypes()

	for _, value := range line {

		DecimalSizePrecision, DecimalSizeScale, DecimalSizeOkToUse := value.DecimalSize()
		Length, LengthOkToUse := value.Length()
		Nullable, NullableOkToUse := value.Nullable()
		var toAppend = ColsTypes{
			Name:                 value.Name(),
			DatabaseTypeName:     value.DatabaseTypeName(),
			DecimalSizePrecision: DecimalSizePrecision,
			DecimalSizeScale:     DecimalSizeScale,
			DecimalSizeOkToUse:   DecimalSizeOkToUse,
			Length:               Length,
			LengthOkToUse:        LengthOkToUse,
			Nullable:             Nullable,
			NullableOkToUse:      NullableOkToUse,
			ScanType:             value.ScanType(),
		}

		ret = append(ret, toAppend)
	}

	return err, ret
}
