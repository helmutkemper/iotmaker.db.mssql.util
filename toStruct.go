package iotmaker_db_mssql_util

import (
	"context"
	"database/sql"
	"errors"
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
	Name              string
	NameOfPrimaryKey  string
	TypeString        string
	TypeReflect       reflect.Type
	SqlTag            string
	IsPrimaryKey      bool
	IsForeignKey      bool
	ToPrimaryKeyQuery string
	ToForeignKeyQuery string
	ToPrimaryKeyVars  string
	ToForeignKeyVars  string
	ToPrimaryKeyScan  string
	ToForeignKeyScan  string
}

func (el ToMakeStruct) MakeStructText(tagSql string) string {
	var ret string

	for structName, structData := range el {
		ret += fmt.Sprintf("type Table%v struct {\n", structName)
		for _, lineData := range structData {
			if lineData.IsPrimaryKey == true {
				ret += fmt.Sprintf("  %v  %v `\n"+
					"%v:\"%v\" \n"+
					"primaryKey:\"true\" \n"+
					"primaryKeyQuery:\"%v\" \n"+
					"primaryKeyScan:\"%v\" \n"+
					"primaryKeyVars:\"%v\" \n"+
					"`\n",
					lineData.Name,
					lineData.TypeString,
					tagSql,
					lineData.SqlTag,
					lineData.ToPrimaryKeyQuery,
					lineData.ToPrimaryKeyScan,
					lineData.ToPrimaryKeyVars,
				)
			} else if lineData.IsForeignKey == true {
				ret += fmt.Sprintf("  %v  []Table%v `"+
					"%v:\"Table%v\" "+
					"primaryKeyQuery:\"%v\" \n"+
					"foreignKeyQuery:\"%v\" \n"+
					"foreignKeyScan:\"%v\" \n"+
					"foreignKeyVars:\"%v\" \n"+
					"primaryKeyScan:\"%v\" \n"+
					"primaryKeyVars:\"%v\" \n"+
					"primaKeyFieldName:\"%v\"\n"+
					"`\n",
					lineData.Name,
					lineData.TypeString,
					tagSql,
					lineData.SqlTag,
					lineData.ToPrimaryKeyQuery,
					lineData.ToForeignKeyQuery,
					lineData.ToForeignKeyScan,
					lineData.ToForeignKeyVars,
					lineData.ToPrimaryKeyScan,
					lineData.ToPrimaryKeyVars,
					lineData.NameOfPrimaryKey,
				)
			} else {
				ret += fmt.Sprintf("  %v  %v `"+
					"%v:\"%v\""+
					"`\n",
					lineData.Name,
					lineData.TypeString,
					tagSql,
					lineData.SqlTag,
				)
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
	var primaryKeyTableName string

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
		err, primaryKeyTableName = el.getPrimaryKeyName(tableName)
		if err != nil {
			return err, nil
		}

		err, tableNameWithRules := NameRules(tableName)
		if err != nil {
			return err, nil
		}

		var lineToRet = make([]ToMakeStructKey, 0)

		for _, dataCol := range tableData {
			_, isPrimaryKey := primaryTableList[tableName][dataCol.Name]
			_, isForeignKey := foreignTableList[tableName][dataCol.Name]
			if isPrimaryKey == true {
				lineToRet = append(lineToRet, ToMakeStructKey{
					Name:              dataCol.NameWithRule,
					NameOfPrimaryKey:  "",
					TypeString:        dataCol.GetScanTypeAsString(),
					TypeReflect:       dataCol.GetScanType(),
					SqlTag:            dataCol.Name,
					IsPrimaryKey:      true,
					IsForeignKey:      false,
					ToForeignKeyQuery: "",
					ToPrimaryKeyQuery: el.mountQuery(tableName, false),
					ToPrimaryKeyScan:  el.mountScanVars(tableName),
					ToForeignKeyScan:  "",
					ToPrimaryKeyVars:  el.mountVars(tableName),
					ToForeignKeyVars:  "",
				})
			} else if isForeignKey == true {
				referenceTableName := foreignTableList[tableName][dataCol.Name].ReferencedObject
				_, referenceTableNameWithRule := NameRules(referenceTableName)

				lineToRet = append(lineToRet, ToMakeStructKey{
					Name:              referenceTableNameWithRule,
					NameOfPrimaryKey:  primaryKeyTableName,
					TypeString:        referenceTableNameWithRule,
					TypeReflect:       dataCol.GetScanType(),
					SqlTag:            referenceTableNameWithRule,
					IsPrimaryKey:      false,
					IsForeignKey:      true,
					ToPrimaryKeyQuery: el.mountQuery(tableName, false),
					ToForeignKeyQuery: el.mountQuery(referenceTableName, true),
					ToPrimaryKeyScan:  el.mountScanVars(tableName),
					ToForeignKeyScan:  el.mountScanVars(referenceTableName),
					ToPrimaryKeyVars:  el.mountVars(tableName),
					ToForeignKeyVars:  el.mountVars(referenceTableName),
				})
			} else {
				lineToRet = append(lineToRet, ToMakeStructKey{
					Name:              dataCol.NameWithRule,
					NameOfPrimaryKey:  "",
					TypeString:        dataCol.GetScanTypeAsString(),
					TypeReflect:       dataCol.GetScanType(),
					SqlTag:            dataCol.Name,
					IsPrimaryKey:      false,
					IsForeignKey:      false,
					ToPrimaryKeyQuery: "",
					ToForeignKeyQuery: "",
					ToPrimaryKeyScan:  "",
					ToForeignKeyScan:  "",
					ToPrimaryKeyVars:  "",
					ToForeignKeyVars:  "",
				})
			}
		}

		ret[tableNameWithRules] = lineToRet
	}

	return nil, ret
}

func (el GoToMSSqlCode) getColumnName(columnName, tableName string) string {
	tableNameWithRules := el.dbConfig[tableName][columnName].TableNameWithRule
	columnNameWithRules := el.dbConfig[tableName][columnName].NameWithRule

	return "table" + tableNameWithRules + "Column" + columnNameWithRules
}

func (el GoToMSSqlCode) getColumnNameWithSqlAsStatement(columnName, tableName string) string {
	tableNameWithRules := el.dbConfig[tableName][columnName].TableNameWithRule
	columnNameWithRules := el.dbConfig[tableName][columnName].NameWithRule

	return columnName + " AS table" + tableNameWithRules + "Column" + columnNameWithRules
}

func (el GoToMSSqlCode) getPrimaryKeyName(tableName string) (error, string) {

	for columnName, columnData := range el.dbConfig[tableName] {

		if columnData.IsPrimaryKey == true {
			return nil, el.getColumnName(columnName, tableName)
		}
	}

	return errors.New("primary key not found"), ""
}

func (el GoToMSSqlCode) mountQuery(tableName string, foreignKey bool) string {
	var primaryKey string
	var query string

	query += fmt.Sprintf("SELECT ")

	l := len(el.dbConfig[tableName]) - 1
	c := 0
	for columnName, columnData := range el.dbConfig[tableName] {
		query += fmt.Sprintf("%v", columnName)

		if columnData.IsPrimaryKey == true {
			primaryKey = columnName
		}

		if c != l {
			query += fmt.Sprint(", ")
		} else {
			query += fmt.Sprint(" ")
		}

		c += 1
	}

	if foreignKey == true {
		query += fmt.Sprintf("FROM %v WHERE %v = %%v", tableName, primaryKey)
	} else {
		query += fmt.Sprintf("FROM %v", tableName)
	}

	return query
}

func (el GoToMSSqlCode) mountVars(tableName string) string {
	var varFromQuery string

	for columnName, columnData := range el.dbConfig[tableName] {
		varFromQuery += fmt.Sprintf("var %v %v\\n", el.getColumnName(columnName, tableName), columnData.ScanType.String())
	}

	return varFromQuery
}

func (el GoToMSSqlCode) mountScanVars(tableName string) string {
	var scanFromQuery string

	l := len(el.dbConfig[tableName]) - 1
	c := 0
	for columnName := range el.dbConfig[tableName] {
		scanFromQuery += fmt.Sprintf("&%v", el.getColumnName(columnName, tableName))

		if c != l {
			scanFromQuery += fmt.Sprint(", ")
		}

		c += 1
	}

	return scanFromQuery
}
