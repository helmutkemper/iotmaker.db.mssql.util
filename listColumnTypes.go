package iotmaker_db_mssql_util

import (
	"context"
	"database/sql"
	"fmt"
)

func ListColumnTypes(db *sql.DB, ctx context.Context, tableName string) (error, map[string]ColumnType) {
	var err error
	var queryReturn *sql.Rows
	var ret = make(map[string]ColumnType)
	var line []*sql.ColumnType
	var nameWithRule string
	var listPriparyKey = make(map[string]PrimaryKeyRelation)
	var isPrimaryKey bool

	err, listPriparyKey = ListPrimaryKeyColumns(db, ctx, tableName)
	if err != nil {
		return err, nil
	}

	queryReturn, err = db.QueryContext(ctx, fmt.Sprintf("SELECT * FROM %v;", tableName))
	if err != nil {
		return err, nil
	}
	defer queryReturn.Close()

	line, err = queryReturn.ColumnTypes()
	if err != nil {
		return err, nil
	}

	for _, value := range line {

		name := value.Name()
		err, nameWithRule = NameRules(name)
		if err != nil {
			return err, nil
		}

		_, isPrimaryKey = listPriparyKey[name]

		decimalSizePrecision, decimalSizeScale, decimalSizeOkToUse := value.DecimalSize()
		length, lengthOkToUse := value.Length()
		nullable, nullableOkToUse := value.Nullable()
		var toAppend = ColumnType{
			Name:                 name,
			NameWithRule:         nameWithRule,
			DatabaseTypeName:     value.DatabaseTypeName(),
			DecimalSizePrecision: decimalSizePrecision,
			DecimalSizeScale:     decimalSizeScale,
			DecimalSizeOkToUse:   decimalSizeOkToUse,
			Length:               length,
			LengthOkToUse:        lengthOkToUse,
			Nullable:             nullable,
			NullableOkToUse:      nullableOkToUse,
			IsPrimaryKey:         isPrimaryKey,
			ScanType:             value.ScanType(),
		}

		ret[name] = toAppend
	}

	return err, ret
}
