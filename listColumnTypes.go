package iotmaker_db_mssql_util

import (
	"database/sql"
	"fmt"
)

func (el GoToMSSqlCode) ListColumnTypes(tableName string) (error, map[string]ColumnType) {
	var err error
	var queryReturn *sql.Rows
	var ret = make(map[string]ColumnType)
	var line []*sql.ColumnType
	var nameWithRule string
	var listPrimaryKey = make(map[string]PrimaryKeyRelation)
	var listForeignKey = make(map[string]ForeignKeyRelation)
	var isPrimaryKey bool
	var isForeignKey bool
	var tableNameWithRule string

	err, tableNameWithRule = NameRules(tableName)
	if err != nil {
		return err, nil
	}

	err, listForeignKey = el.ListForeignKeyColumns(tableName)

	err, listPrimaryKey = el.ListPrimaryKeyColumns(tableName)
	if err != nil {
		return err, nil
	}

	queryReturn, err = el.Db.QueryContext(el.Ctx, fmt.Sprintf("SELECT * FROM [%v];", tableName))
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

		_, isPrimaryKey = listPrimaryKey[name]
		_, isForeignKey = listForeignKey[name]

		decimalSizePrecision, decimalSizeScale, decimalSizeOkToUse := value.DecimalSize()
		length, lengthOkToUse := value.Length()
		nullable, nullableOkToUse := value.Nullable()
		var toAppend = ColumnType{
			Name:                 name,
			NameWithRule:         nameWithRule,
			TableNameWithRule:    tableNameWithRule,
			DatabaseTypeName:     value.DatabaseTypeName(),
			DecimalSizePrecision: decimalSizePrecision,
			DecimalSizeScale:     decimalSizeScale,
			DecimalSizeOkToUse:   decimalSizeOkToUse,
			Length:               length,
			LengthOkToUse:        lengthOkToUse,
			Nullable:             nullable,
			NullableOkToUse:      nullableOkToUse,
			IsPrimaryKey:         isPrimaryKey,
			IsForeignKey:         isForeignKey,
			ScanType:             value.ScanType(),
		}

		ret[name] = toAppend
	}

	return err, ret
}
