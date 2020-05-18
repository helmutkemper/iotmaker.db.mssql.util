package iotmaker_db_mssql_util

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

type GoToMSSqlCode struct {
	Db  *sql.DB
	Ctx context.Context

	dbConfig map[string]map[string]ColumnType
}

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
		ret += fmt.Sprintf("type Table%v struct {\n", structName)
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

func (el *GoToMSSqlCode) ToStruct() (error, ToMakeStruct) {
	var err error
	var tableNameList []string
	var primaryTableList = make(map[string]map[string]PrimaryKeyRelation)
	var foreignTableList = make(map[string]map[string]ForeignKeyRelation)
	var ret = make(ToMakeStruct)

	el.dbConfig = make(map[string]map[string]ColumnType)

	err, tableNameList = el.ListTables()
	if err != nil {
		return err, nil
	}

	for _, tableName := range tableNameList {
		var listOfColType = make(map[string]ColumnType)
		var listForeignKey map[string]ForeignKeyRelation
		var listPrimaryKey map[string]PrimaryKeyRelation

		err, listOfColType = el.ListColumnTypes(tableName)
		if err != nil {
			return err, nil
		}

		el.dbConfig[tableName] = listOfColType

		err, listPrimaryKey = el.ListPrimaryKeyColumns(tableName)
		if err != nil {
			return err, nil
		}

		if len(listPrimaryKey) != 0 {
			primaryTableList[tableName] = listPrimaryKey
		}

		err, listForeignKey = el.ListForeignKeyColumns(tableName)
		if err != nil {
			return err, nil
		}

		if len(listForeignKey) != 0 {
			foreignTableList[tableName] = listForeignKey
		}
	}

	for tableName, tableData := range el.dbConfig {

		_, tableNameWithRules := NameRules(tableName)
		var lineToRet = make([]ToMakeStructKey, 0)

		for _, dataCol := range tableData {
			_, isPrimaryKey := primaryTableList[tableName][dataCol.Name]
			_, isForeignKey := foreignTableList[tableName][dataCol.Name]
			if isPrimaryKey == true {
				lineToRet = append(lineToRet, ToMakeStructKey{
					Name:         dataCol.NameWithRule,
					TypeString:   dataCol.GetScanTypeAsString(),
					TypeReflect:  dataCol.GetScanType(),
					SqlTag:       dataCol.Name,
					IsPrimaryKey: true,
				})
			} else if isForeignKey == true {
				referenceTableName := foreignTableList[tableName][dataCol.Name].ReferencedObject
				_, referenceTableNameWithRule := NameRules(referenceTableName)

				lineToRet = append(lineToRet, ToMakeStructKey{
					Name:         referenceTableNameWithRule,
					TypeString:   referenceTableNameWithRule,
					TypeReflect:  dataCol.GetScanType(),
					SqlTag:       dataCol.Name,
					IsForeignKey: true,
					ToQuery:      el.mountQuery(referenceTableName),
					ToScan:       el.mountScanVars(referenceTableName),
					ToVars:       el.mountVars(referenceTableName),
				})
			} else {
				lineToRet = append(lineToRet, ToMakeStructKey{
					Name:        dataCol.NameWithRule,
					TypeString:  dataCol.GetScanTypeAsString(),
					TypeReflect: dataCol.GetScanType(),
					SqlTag:      dataCol.Name,
				})
			}
		}

		ret[tableNameWithRules] = lineToRet
	}

	return nil, ret
}

func (el GoToMSSqlCode) mountQuery(tableName string) string {
	var primaryKey string
	var query string

	query += fmt.Sprintf("SELECT ")

	l := len(el.dbConfig) - 1
	c := 0
	for k, v := range el.dbConfig[tableName] {
		query += fmt.Sprintf("%v", k)

		if v.IsPrimaryKey == true {
			primaryKey = k
		}

		if c != l {
			query += fmt.Sprint(", ")
		} else {
			query += fmt.Sprint(" ")
		}

		c += 1
	}

	query += fmt.Sprintf("FROM %v WHERE %v = %%v", tableName, primaryKey)

	return query
}

func (el GoToMSSqlCode) mountVars(tableName string) string {
	var varFromQuery string

	for _, v := range el.dbConfig[tableName] {
		varFromQuery += fmt.Sprintf("  var column%v %v\\n", v.NameWithRule, v.ScanType.String())
	}

	return varFromQuery
}

func (el GoToMSSqlCode) mountScanVars(tableName string) string {
	var scanFromQuery string

	l := len(el.dbConfig) - 1
	c := 0
	for _, v := range el.dbConfig[tableName] {
		scanFromQuery += fmt.Sprintf("&column%v", v.NameWithRule)

		if c != l {
			scanFromQuery += fmt.Sprint(", ")
		}

		c += 1
	}

	return scanFromQuery
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

/*
element := reflect.ValueOf(el).Elem()
*/
