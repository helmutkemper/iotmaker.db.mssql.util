package iotmaker_db_mssql_util

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
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
	ToForeignKeyList  string
}

func (el ToMakeStruct) MakeStructText(tagSql string) string {
	var ret string

	for structName, structData := range el {
		ret += fmt.Sprintf("type Table%v struct {\n", structName)
		for _, lineData := range structData {
			if lineData.IsPrimaryKey == true {
				ret += fmt.Sprintf("  %v  %v `"+
					"%v:\"%v\" "+
					"primaryKey:\"true\" "+
					"primaKeyFieldName:\"%v\""+
					"primaryKeyQuery:\"%v\" "+
					"primaryKeyScan:\"%v\" "+
					"primaryKeyVars:\"%v\" "+
					"toForeignKeyList:\"%v\" "+
					"`\n",
					lineData.Name,
					lineData.TypeString,
					tagSql,
					lineData.SqlTag,
					structName,
					lineData.ToPrimaryKeyQuery,
					lineData.ToPrimaryKeyScan,
					lineData.ToPrimaryKeyVars,
					lineData.ToForeignKeyList,
				)
			} else if lineData.IsForeignKey == true {
				ret += fmt.Sprintf("  %v  []Table%v `"+
					"%v:\"Table%v\" "+
					"primaryKeyQuery:\"%v\" "+
					"foreignKeyQuery:\"%v\" "+
					"foreignKeyScan:\"%v\" "+
					"foreignKeyVars:\"%v\" "+
					"primaryKeyScan:\"%v\" "+
					"primaryKeyVars:\"%v\" "+
					"primaKeyFieldName:\"%v\""+
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

func (el GoToMSSqlCode) MakeFile(tableName string) error {
	/*err, tableNameWithRule := NameRules(tableName)
	  if err != nil {
	    return err
	  }

	  err, foreignKeyColumnsList := el.ListForeignKeyColumns(tableName)
	  if err != nil {
	    panic(err)
	  }

	  toFile := ""
		toFile += fmt.Sprintf("func main() {\n")
		toFile += fmt.Sprintf("  var err error\n")
	  toFile += fmt.Sprintf("  var db *sql.DB\n")
	  toFile += fmt.Sprintf("  var ctx context.Context\n")
		toFile += fmt.Sprintf("  var %vRows *sql.Rows\n", tableNameWithRule)

		for _, fkData := range foreignKeyColumnsList {
	    toFile += fmt.Sprintf("  var %vRows *sql.Rows\n", fkData.ReferencedObjectWithRule)
	  }

		toFile += fmt.Sprintf("\n")

		data := el.prepare(tableName)
		toFile += fmt.Sprintf(
		  "%v\n",
		  el.queryCodePrimaryKey(
		    2,
		    tableName,
		    data,
	    ),
	  )

		for _, structConfig := range foreignKeyColumnsList {
			//data = el.prepare(structConfig.ReferencedObject)

			/*var primaryKeyName, primaryKeyNameWithRule string
			for _, v1 := range el.dbConfig[structConfig.ReferencedObject] {
				if v1.IsPrimaryKey == true {
					primaryKeyName = v1.Name
					primaryKeyNameWithRule = v1.NameWithRule
					break
				}
			}*/

	/*toFile += fmt.Sprintf("%v\n",
			el.queryCodeForeignKey(4, structConfig),
		)
	}

	fmt.Printf("%v\n", toFile)*/
	return nil
}

type prepareData struct {
	fieldTypeAll                map[string]string
	fieldNameAll                map[string]string
	fieldNameWithRuleAll        map[string]string
	fieldNameAllForeign         map[string]string
	fieldNameWithRuleAllForeign map[string]string
	fieldNameAllPrimary         map[string]string
	fieldNameWithRuleAllPrimary map[string]string
}

func (el *GoToMSSqlCode) queryCodeForeignKey(spaces int, data ForeignKeyRelation) string {

	/*spacesString := ""
		for i := 0; i != spaces; i += 1 {
			spacesString += " "
		}

	  dataPrepared := el.prepare(data.ReferencedObject)

		toFile := ""
		toFile += fmt.Sprintf(
		  "%v%vRows, err = db.QueryContext(ctx, fmt.Sprintf(\"SELECT %v FROM %v WHERE %v = %%v\", %v))\n",
		  spacesString,
		  data.ReferencedObjectWithRule,
		  el.joinMap(dataPrepared.fieldNameAll, ", "),
	    data.ReferencedObject,
	    el.joinMap(dataPrepared.fieldNameAllPrimary, ""),
	    data.TableNameWithRule+"Column"+data.ConstraintColumnNameWithRule,
	  )
		toFile += fmt.Sprintf("%vif err != nil {\n", spacesString)
		toFile += fmt.Sprintf("%v  panic(err)\n", spacesString)
		toFile += fmt.Sprintf("%v}\n", spacesString)

		for _, line := range dataPrepared.fieldNameWithRuleAllForeign {
			toFile += fmt.Sprintf("%vvar arrayOfStruct%v = make([]%v, 0)\n", spacesString, line, line)
		}
		toFile += fmt.Sprintf("\n")

		toFile += fmt.Sprintf("%vfor %vRows.Next() {\n", spacesString, data.ReferencedObjectWithRule)
		for k, line := range dataPrepared.fieldNameWithRuleAll {
			toFile += fmt.Sprintf("%v  var %vColumn%v %v\n", spacesString, data.ReferencedObjectWithRule, line, dataPrepared.fieldTypeAll[k])
		}

		toFile += fmt.Sprintf("%v  %vRows.Scan(", spacesString, data.ReferencedObjectWithRule)
		for _, line := range dataPrepared.fieldNameWithRuleAll {
			toFile += fmt.Sprintf("&%vColumn%v, ", data.ReferencedObjectWithRule, line)
		}
		toFile = strings.TrimRight(toFile, ", ")
		toFile += fmt.Sprintf(")\n\n")

		toStruct := ""
		for _, v := range el.dbConfig[data.ReferencedObject] {
	    toStruct += fmt.Sprintf("%v: %v, ", v.NameWithRule, v.VarName)
	  }
	  toStruct = strings.TrimRight(toStruct, ", ")
	  toStruct = fmt.Sprintf("{%v}", toStruct)

	  toFile += fmt.Sprintf("%v  arrayOfStruct%v = append(arrayOfStruct%v, Table%v%v)\n", spacesString, data.ReferencedObjectWithRule, data.ReferencedObjectWithRule, data.ReferencedObjectWithRule, toStruct)

		toFile += fmt.Sprintf("%v}\n", spacesString)

		return toFile*/
	return ""
}

func (el *GoToMSSqlCode) queryCodePrimaryKey(spaces int, tableName string, data prepareData) string {
	/*  var err error
	  var tableNameWithRule string

	  err, tableNameWithRule = NameRules(tableName)
	  if err != nil {
	    panic(err)
	  }

	  err, foreignKeyColumnsList := el.ListForeignKeyColumns(tableName)
	  if err != nil {
	    panic(err)
	  }

		spacesString := ""
		for i := 0; i != spaces; i += 1 {
			spacesString += " "
		}

		toFile := ""
		toFile += fmt.Sprintf(
		  "%v%vRows, err = db.QueryContext(ctx, \"-SELECT %v FROM %v\")\n",
		  spacesString,
		  tableNameWithRule,
		  el.joinMap(data.fieldNameAll, ", "),
		  tableName,
	  )
		toFile += fmt.Sprintf("%vif err != nil {\n", spacesString)
		toFile += fmt.Sprintf("%v  panic(err)\n", spacesString)
		toFile += fmt.Sprintf("%v}\n", spacesString)

	  for _, fkData := range foreignKeyColumnsList {
	    toFile += fmt.Sprintf("%vvar arrayOfStruct%v = make([]Table%v, 0)\n", spacesString, fkData.ReferencedObjectWithRule, fkData.ReferencedObjectWithRule)
	  }

		toFile += fmt.Sprintf("\n")

		toFile += fmt.Sprintf("%vfor %vRows.Next() {\n", spacesString, tableNameWithRule)
		for k, line := range data.fieldNameWithRuleAll {
			toFile += fmt.Sprintf("%v  var %vColumn%v %v\n", spacesString, tableNameWithRule, line, data.fieldTypeAll[k])
		}

		toFile += fmt.Sprintf("%v  %vRows.Scan(", spacesString, tableNameWithRule)

		keys := make([]string, 0, len(data.fieldNameWithRuleAll))
		for k := range data.fieldNameWithRuleAll {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			toFile += fmt.Sprintf("&%vColumn%v, ", tableNameWithRule, data.fieldNameWithRuleAll[k])
		}
		toFile = strings.TrimRight(toFile, ", ")
		toFile += fmt.Sprintf(")\n")

		return toFile

	*/
	return ""
}

func (el *GoToMSSqlCode) joinMap(list map[string]string, glue string) string {
	var ret = ""

	keys := make([]string, 0, len(list))
	for k := range list {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		ret += list[k] + glue
	}
	return strings.TrimRight(ret, glue)
}

func (el *GoToMSSqlCode) Analyze() (error, ToMakeStruct) {
	/*var err error
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

		//err, listOfColType = el.GetFieldsTypeFromTable(tableName)
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

			//fixme: to function - start
			var toForeignKeyList = ""
			for _, columnData := range foreignTableList[tableName] {
				toForeignKeyList += columnData.ReferencedObjectWithRule + ","
			}
			toForeignKeyList = strings.TrimRight(toForeignKeyList, ",")
			//fixme: to function - end

			if isPrimaryKey == true {
				lineToRet = append(lineToRet, ToMakeStructKey{
					Name:              dataCol.NameWithRule,
					NameOfPrimaryKey:  dataCol.NameWithRule,
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
					ToForeignKeyList:  toForeignKeyList,
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
					ToPrimaryKeyQuery: "", //el.mountQuery(tableName, false),
					ToForeignKeyQuery: el.mountQuery(referenceTableName, true),
					ToPrimaryKeyScan:  "", //el.mountScanVars(tableName),
					ToForeignKeyScan:  el.mountScanVars(referenceTableName),
					ToPrimaryKeyVars:  el.mountVars(tableName),
					ToForeignKeyVars:  "", //el.mountVars(referenceTableName),
					ToForeignKeyList:  "",
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
					ToForeignKeyList:  "",
				})
			}
		}

		ret[tableNameWithRules] = lineToRet
	}

	return nil, ret

	*/
	return nil, nil
}

func (el GoToMSSqlCode) getColumnName(columnName, tableName string) string {
	//tableNameWithRules := el.dbConfig[tableName][columnName].TableNameWithRule
	//columnNameWithRules := el.dbConfig[tableName][columnName].NameWithRule

	return "table" //+ /*tableNameWithRules +*/ "Column" + columnNameWithRules
}

func (el GoToMSSqlCode) getColumnNameWithSqlAsStatement(columnName, tableName string) string {
	//tableNameWithRules := el.dbConfig[tableName][columnName].TableNameWithRule
	//columnNameWithRules := el.dbConfig[tableName][columnName].NameWithRule

	return "" //columnName + " AS table" /*+ tableNameWithRules*/ + "Column" + columnNameWithRules
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
