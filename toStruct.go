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
	ToQuery      string
	ToVars       string
	ToScan       string
}

func (el ToMakeStruct) MakeStructText(tagSql string) string {
	var ret string

	for structName, structData := range el {
		ret = fmt.Sprintf("type Table%v struct {\n", structName)
		for _, lineData := range structData {
			if lineData.IsPrimaryKey == true {
				ret += fmt.Sprintf("  %v  %v `%v:\"%v\" primaryKey:\"true\"`\n", lineData.Name, lineData.TypeString, tagSql, lineData.SqlTag)
			} else if lineData.IsForeignKey == true {
				ret += fmt.Sprintf("  %v  []Table%v `%v:\"%v\" query:\"%v\" scan:\"%v\" vars:\"%v\"`\n", lineData.Name, lineData.TypeString, tagSql, lineData.SqlTag, lineData.ToQuery, lineData.ToScan, lineData.ToVars)
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
			structKey, structType, query, scanFromQuery, varFromQuery, structRealType := isForeignKeyColumn(dataCol, foreignKeyData.ReferencedObject, tableDescrptionList)
			lineToRet = append(lineToRet, ToMakeStructKey{
				Name:         structKey,
				TypeString:   structType,
				TypeReflect:  structRealType,
				SqlTag:       dataCol.Name,
				IsForeignKey: true,
				ToQuery:      query,
				ToScan:       scanFromQuery,
				ToVars:       varFromQuery,
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

func isForeignKeyColumn(dataColumn ColumnType, tableName string, dbConfig map[string]map[string]ColumnType) (string, string, string, string, string, reflect.Type) {
	var primaryKey string
	var query string
	var scanFromQuery string
	var varFromQuery string

	query += fmt.Sprintf("SELECT ")

	l := len(dbConfig) - 1
	c := 0
	for k, v := range dbConfig[tableName] {
		query += fmt.Sprintf("%v", k)
		scanFromQuery += fmt.Sprintf("&column%v", v.NameWithRule)
		varFromQuery += fmt.Sprintf("  var column%v %v\\n", v.NameWithRule, v.ScanType.String())

		if v.IsPrimaryKey == true {
			primaryKey = k
		}

		if c != l {
			query += fmt.Sprint(", ")
			scanFromQuery += fmt.Sprint(", ")
		} else {
			query += fmt.Sprint(" ")
		}

		c += 1
	}

	query += fmt.Sprintf("FROM %v WHERE %v = %%v", tableName, primaryKey)
	_, tableName = NameRules(tableName)
	return tableName, tableName, query, scanFromQuery, varFromQuery, dataColumn.ScanType
}

func normalKeyColumn(dataColumn ColumnType) (string, string, reflect.Type) {
	return dataColumn.NameWithRule, dataColumn.ScanType.String(), dataColumn.ScanType
}
