package iotmaker_db_mssql_util

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

type ToMakeStruct map[string][]ToMakeStructKey

type ToMakeStructKey struct {
	Name        string
	TypeString  string
	TypeReflect reflect.Type
	SqlTag      string
}

func (el ToMakeStruct) MakeStructText(tagSql string) string {
	var ret string
	for structName, structData := range el {
		ret = fmt.Sprintf("type %v struct {\n", structName)

		for _, lineData := range structData {
			ret += fmt.Sprintf("  %v  %v `%v:\"%v\"`\n", lineData.Name, lineData.TypeString, tagSql, lineData.SqlTag)
		}

		ret += fmt.Sprintf("}\n")
	}

	return ret
}

func ToStruct(db *sql.DB, ctx context.Context, tableName string, data map[string]ColumnType) (error, ToMakeStruct) {
	var err error
	var tableNameWithRule string
	var tableNameList []string
	var foreignTableList = make(map[string]map[string]ForeignKeyRelation)
	var ret = make(ToMakeStruct)
	var lineToRet = make([]ToMakeStructKey, 0)

	err, tableNameWithRule = NameRules(tableName)
	if err != nil {
		return err, nil
	}

	err, tableNameList = ListTables(db, ctx)
	if err != nil {
		return err, nil
	}

	for _, name := range tableNameList {
		var list map[string]ForeignKeyRelation
		err, list = ListForeignKeyColumns(db, ctx, name)
		if err != nil {
			return err, nil
		}

		if len(list) == 0 {
			continue
		}

		foreignTableList[name] = list
	}

	for _, dataCol := range data {
		key, ok := foreignTableList[tableName][dataCol.Name]
		if ok != true {
			structKey, structType, structRealType := notForeignKeyColumn(dataCol)
			lineToRet = append(lineToRet, ToMakeStructKey{
				Name:        structKey,
				TypeString:  structType,
				TypeReflect: structRealType,
				SqlTag:      dataCol.Name,
			})
		} else {
			structKey, structType, structRealType := isForeignKeyColumn(dataCol, key.ReferencedObject)
			lineToRet = append(lineToRet, ToMakeStructKey{
				Name:        structKey,
				TypeString:  structType,
				TypeReflect: structRealType,
				SqlTag:      dataCol.Name,
			})
		}
	}

	ret[tableNameWithRule] = lineToRet

	return nil, ret
}

func notForeignKeyColumn(dataColumn ColumnType) (string, string, reflect.Type) {
	return dataColumn.NameWithRule, dataColumn.ScanType.String(), dataColumn.ScanType
}

func isForeignKeyColumn(dataColumn ColumnType, tableName string) (string, string, reflect.Type) {
	_, tableName = NameRules(tableName)
	return tableName, "[]" + tableName, dataColumn.ScanType
}
