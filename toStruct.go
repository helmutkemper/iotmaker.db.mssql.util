package iotmaker_db_mssql_util

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

type ToMakeStruct map[string][]ToMakeStructKey

type ToMakeStructKey struct {
	Name         string
	TypeString   string
	TypeReflect  reflect.Type
	SqlTag       string
	IsPrimaryKey bool
	IsForeignKey bool
}

func (el ToMakeStruct) MakeStructText(tagSql string) string {
	var ret string

	for structName, structData := range el {
		ret = fmt.Sprintf("type Table%v struct {\n", structName)
		for _, lineData := range structData {
			if lineData.IsPrimaryKey == true {
				ret += fmt.Sprintf("  %v  %v `%v:\"%v\" primaryKey:\"true\"`\n", lineData.Name, lineData.TypeString, tagSql, lineData.SqlTag)
			} else if lineData.IsForeignKey == true {
				ret += fmt.Sprintf("  %v  []Table%v `%v:\"%v\"`\n", lineData.Name, lineData.TypeString, tagSql, lineData.SqlTag)
			} else {
				ret += fmt.Sprintf("  %v  %v `%v:\"%v\"`\n", lineData.Name, lineData.TypeString, tagSql, lineData.SqlTag)
			}
		}

		ret += fmt.Sprintf("}\n")
	}

	return ret
}

func ToStruct(db *sql.DB, ctx context.Context, tableName string, data map[string]ColumnType) (error, ToMakeStruct) {
	var err error
	var tableNameWithRule string
	var tableNameList []string
	var tableDescrptionList = make(map[string]map[string]ColumnType)
	var primaryTableList = make(map[string]map[string]PrimaryKeyRelation)
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
		var listOfColType = make(map[string]ColumnType)
		var listForeignKey map[string]ForeignKeyRelation
		var listPrimaryKey map[string]PrimaryKeyRelation

		err, listOfColType = ListColumnTypes(db, ctx, name)
		if err != nil {
			return err, nil
		}

		tableDescrptionList[name] = listOfColType

		err, listPrimaryKey = ListPrimaryKeyColumns(db, ctx, name)
		if err != nil {
			return err, nil
		}

		if len(listPrimaryKey) != 0 {
			primaryTableList[name] = listPrimaryKey
		}

		err, listForeignKey = ListForeignKeyColumns(db, ctx, name)
		if err != nil {
			return err, nil
		}

		if len(listForeignKey) != 0 {
			foreignTableList[name] = listForeignKey
		}
	}

	for _, dataCol := range data {
		_, isPrimaryKey := primaryTableList[tableName][dataCol.Name]
		foreignKeyData, isForeignKey := foreignTableList[tableName][dataCol.Name]
		if isPrimaryKey == true {
			structKey, structType, structRealType := isPrimaryKeyColumn(dataCol)
			lineToRet = append(lineToRet, ToMakeStructKey{
				Name:         structKey,
				TypeString:   structType,
				TypeReflect:  structRealType,
				SqlTag:       dataCol.Name,
				IsPrimaryKey: true,
			})
		} else if isForeignKey == true {
			structKey, structType, structRealType := isForeignKeyColumn(dataCol, foreignKeyData.ReferencedObject, tableDescrptionList)
			lineToRet = append(lineToRet, ToMakeStructKey{
				Name:         structKey,
				TypeString:   structType,
				TypeReflect:  structRealType,
				SqlTag:       dataCol.Name,
				IsForeignKey: true,
			})
		} else {
			structKey, structType, structRealType := normalKeyColumn(dataCol)
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

func isPrimaryKeyColumn(dataColumn ColumnType) (string, string, reflect.Type) {
	return dataColumn.NameWithRule, dataColumn.ScanType.String(), dataColumn.ScanType
}

func isForeignKeyColumn(dataColumn ColumnType, tableName string, dbConfig map[string]map[string]ColumnType) (string, string, reflect.Type) {
	_, tableName = NameRules(tableName)
	return tableName, tableName, dataColumn.ScanType
}

func normalKeyColumn(dataColumn ColumnType) (string, string, reflect.Type) {
	return dataColumn.NameWithRule, dataColumn.ScanType.String(), dataColumn.ScanType
}
