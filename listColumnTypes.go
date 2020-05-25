package iotmaker_db_mssql_util

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
)

type Table struct {
	Data map[string]ColumnType
	Ref  *Database
}

func (el Table) getListAsc() []string {
	listOrder := make([]string, 0, len(el.Data))
	for fieldName := range el.Data {
		listOrder = append(listOrder, fieldName)
	}

	sort.Strings(listOrder)
	return listOrder
}

func (el Table) GetTableOrderByFieldAsc() []ColumnType {
	ret := make([]ColumnType, 0, len(el.Data))
	for _, fieldName := range el.getListAsc() {
		ret = append(ret, el.Data[fieldName])
	}

	return ret
}

func (el Table) GetSelectQuery() string {
	return "\"SELECT " + el.GetFieldList() + " FROM [" + el.GetTableName() + "]\""
}

func (el Table) GetForeignKeyFunction(field ColumnType) string {
	return "Get" + field.GetReferenceTableNameWithRules() + "AsArray(db, ctx, " + field.GetAsVarName() + ")"
}

func (el Table) GetSelectWhereQuery() string {
	err, pk := el.GetPrimaryKeyRef()
	if err != nil {
		panic(errors.New("table " + el.GetTableName() + "has't primary key"))
	}
	return "\"SELECT " + el.GetFieldList() + " FROM [" + el.GetTableName() + "] WHERE " + pk.GetFieldNameWithRules() + " = %v\""
	//return "\"SELECT " + el.GetFieldList() + " FROM " + el.GetTableName() + " WHERE " + pk.GetFieldNameWithRules() + " = %v\"" + pk.GetAsVarName()
}

func (el Table) GetFieldList() string {
	ret := ""

	for k, v := range el.getListAsc() {
		if k != 0 {
			ret += ", "
		}
		ret += el.Data[v].GetFieldName()
	}

	return ret
}

func (el Table) GetPrimaryKeyRef() (error, ColumnType) {
	for _, v := range el.getListAsc() {
		if el.Data[v].IsPrimaryKey == true {
			return nil, el.Data[v]
		}
	}

	return errors.New("primary key not found"), ColumnType{}
}

func (el Table) GetRefVarList() string {
	ret := ""

	for k, v := range el.getListAsc() {
		if k != 0 {
			ret += ", "
		}
		ret += el.Data[v].GetAsRefVar()
	}

	return ret
}

func (el Table) GetStructFromComplexTable() string {
	ret := ""

	for _, v := range el.getListAsc() {
		ret += "  "

		if el.Data[v].IsForeignKey == true {
			ret += el.Data[v].GetReferenceTableNameWithRules()
			ret += "  "
			ret += "[]Data" + el.Data[v].GetReferenceTableNameWithRules()
		} else {
			ret += el.Data[v].GetFieldNameWithRules()
			ret += "  "
			ret += el.Data[v].GetTypeAsString()
		}

		ret += "\n"
	}

	return "{\n" + ret + "}\n"
}

func (el Table) GetStruct() string {
	ret := ""

	for _, v := range el.getListAsc() {
		ret += "  "
		ret += el.Data[v].GetFieldNameWithRules()
		ret += "  "
		ret += el.Data[v].GetTypeAsString()
		ret += "\n"
	}

	return "{\n" + ret + "}\n"
}

func (el Table) GetStructWidtVars() string {
	ret := ""

	for _, v := range el.getListAsc() {
		ret += "  "
		ret += el.Data[v].GetFieldNameWithRules()
		ret += ": "
		ret += el.Data[v].GetAsVarName()
		ret += ",\n"
	}

	return "Table" + el.GetTableNameWithRules() + "{\n" + ret + "}\n"
}

func (el Table) GetStructDataWidtVars() string {
	ret := ""

	for _, v := range el.getListAsc() {
		ret += "  "
		ret += el.Data[v].GetFieldNameWithRules()
		ret += ": "
		ret += el.Data[v].GetAsVarName()
		ret += ",\n"
	}

	return "Data" + el.GetTableNameWithRules() + "{\n" + ret + "}\n"
}

func (el Table) GetStructWidtForeignKey() string {
	ret := ""

	for _, v := range el.getListAsc() {
		ret += "  "

		if el.Data[v].IsForeignKey == true {
			ret += (*el.Ref).Data[el.Data[v].ReferenceTableName].GetTableNameWithRules()
			ret += ": "
			ret += el.GetForeignKeyFunction(el.Data[v])
		} else {
			ret += el.Data[v].GetFieldNameWithRules()
			ret += ": "
			ret += el.Data[v].GetAsVarName()
		}

		ret += ",\n"
	}

	return "Data" + el.GetTableNameWithRules() + "{\n" + ret + "}\n"
}

func (el Table) GetDefStructFromComplexTable() string {
	return "type Data" + el.GetTableNameWithRules() + " struct " + el.GetStructFromComplexTable()
}

func (el Table) GetDefStruct() string {
	return "type Table" + el.GetTableNameWithRules() + " struct " + el.GetStruct()
}

func (el Table) GetTableName() string {
	for _, fieldData := range el.Data {
		return fieldData.TableName
	}

	return ""
}

func (el Table) GetTableNameWithRules() string {
	for _, fieldData := range el.Data {
		return fieldData.TableNameWithRules
	}

	return ""
}

func (el Table) GetRowDefVar() string {
	return "var Row" + el.GetTableNameWithRules() + " *sql.Rows"
}

func (el Table) GetRowVarName() string {
	return "Row" + el.GetTableNameWithRules()
}

func (el Table) GetTableAsVarName() string {
	for _, fieldData := range el.Data {
		return "Table" + fieldData.TableNameWithRules
	}

	return ""
}

func (el Table) GetDataAsVarName() string {
	for _, fieldData := range el.Data {
		return "Data" + fieldData.TableNameWithRules
	}

	return ""
}

func (el Table) GetTableAsDefArrayVar() string {

	for range el.Data {
		return "var ArrayOf" + el.GetTableAsVarName() + " = make([]" + el.GetTableAsVarName() + ", 0)"
	}

	return ""
}

func (el Table) GetDataAsDefArrayVar() string {

	for range el.Data {
		return "var ArrayOf" + el.GetTableAsVarName() + " = make([]" + el.GetDataAsVarName() + ", 0)"
	}

	return ""
}

func (el Table) GetNameFromArrayVar() string {

	for range el.Data {
		return "ArrayOf" + el.GetTableAsVarName()
	}

	return ""
}

func (el Table) GetAppendVar() string {

	for range el.Data {
		return "ArrayOf" + el.GetTableAsVarName() + " = append(ArrayOf" + el.GetTableAsVarName() + ", data)"
	}

	return ""
}

func (el Table) GetTableAsDefVarToPopulate() string {

	for range el.Data {
		return "var " + el.GetTableAsVarName() + " = " + el.GetStructWidtVars()
	}

	return ""
}

func (el Table) GetDataAsDefVarToPopulate() string {

	for range el.Data {
		return "var data = " + el.GetStructWidtForeignKey()
	}

	return ""
}

type Database struct {
	db   *sql.DB
	ctx  context.Context
	Data map[string]Table
}

func (el Database) GetTable(tableName string) Table {
	return el.Data[tableName]
}

func (el Database) GetCodeFromPrimaryKey(tableName string) string {
	ret := ""
	ret += "func Get" + el.Data[tableName].GetTableNameWithRules() + "(db *sql.DB, ctx context.Context) error {\n"
	ret += "  var err error\n"
	ret += "  " + el.Data[tableName].GetRowDefVar() + "\n"
	ret += "  " + el.Data[tableName].GetRowVarName() + ", err = db.QueryContext(ctx, " + el.Data[tableName].GetSelectQuery() + ")\n"
	ret += "  if err != nil {\n"
	ret += "    return err\n"
	ret += "  }\n"
	ret += "  if " + el.Data[tableName].GetRowVarName() + ".Next() {\n"

	for _, fieldData := range el.Data[tableName].GetTableOrderByFieldAsc() {
		ret += "    " + fieldData.GetAsDefVar() + "\n"
	}

	ret += "    err = " + el.Data[tableName].GetRowVarName() + ".Scan(" + el.Data[tableName].GetRefVarList() + ")\n"
	ret += "    if err != nil {\n"
	ret += "      return err\n"
	ret += "    }\n"

	ret += "    " + el.Data[tableName].GetDataAsDefVarToPopulate() + "\n"

	ret += "    YourFunction(data)\n"
	ret += "  }\n"
	ret += "  return nil\n"

	return ret
}

func (el Database) GetCodeFromForeignKey(tableName string, field ColumnType) string {
	ret := "\n"
	ret += "func Get" + el.Data[tableName].GetTableNameWithRules() + "AsArray(db *sql.DB, ctx context.Context, " + field.GetAsVarName() + " " + field.TypeAsString + ") []Data" + el.Data[tableName].GetTableNameWithRules() + " {\n"
	ret += "  var err error\n"
	ret += "  " + el.Data[tableName].GetDataAsDefArrayVar() + "\n"
	ret += "  " + el.Data[tableName].GetRowDefVar() + "\n"
	ret += "  " + el.Data[tableName].GetRowVarName() + ", err = db.QueryContext(ctx, fmt.Sprintf(" + el.Data[tableName].GetSelectWhereQuery() + "," + " " + field.GetAsVarName() + "))\n"
	ret += "  if err != nil {\n"
	ret += "    panic(err)\n"
	ret += "  }\n"
	ret += "  defer " + el.Data[tableName].GetRowVarName() + ".Close()\n"
	ret += "  if " + el.Data[tableName].GetRowVarName() + ".Next() {\n"

	for _, fieldData := range el.Data[tableName].GetTableOrderByFieldAsc() {
		ret += "    " + fieldData.GetAsDefVar() + "\n"
	}

	ret += "    err = " + el.Data[tableName].GetRowVarName() + ".Scan(" + el.Data[tableName].GetRefVarList() + ")\n"
	ret += "    if err != nil {\n"
	ret += "      panic(err)\n"
	ret += "    }\n"

	ret += "    " + "var data = " + el.Data[tableName].GetStructDataWidtVars() + "\n"
	ret += "    " + el.Data[tableName].GetAppendVar() + "\n"

	ret += "  }\n"
	ret += "  return " + el.Data[tableName].GetNameFromArrayVar() + "\n"
	ret += "}\n"

	return ret
}

func (el Database) ToCode(tableName, packageName string) string {
	var err error

	ret := "package " + packageName + "\n"
	ret += "\n"
	ret += "import (\n"
	ret += "  \"context\"\n"
	ret += "  \"database/sql\"\n"
	ret += "  \"fmt\"\n"
	ret += "  _ \"github.com/denisenkom/go-mssqldb\"\n"
	//ret += "  mssqlUtil \"github.com/helmutkemper/iotmaker.db.mssql.util\"\n"
	ret += ")\n"
	ret += "func main() {\nvar err error\nvar db *sql.DB\nvar ctx context.Context\n\nctx = context.TODO()\n\nconnString := fmt.Sprintf(\"server=%s;port=%d;database=%s;user id=%s;password=%s\",\n\"localhost\",\n1434,\n\"toExport\",\n\"CS\\\\helmut.kemper\",\n\"temp@123\",\n)\n\ndb, err = sql.Open(\"sqlserver\", connString)\nif err != nil {\npanic(err)\n}\nerr = db.PingContext(ctx)\nif err != nil {\npanic(err)\n}\nerr = GetUser(db, ctx)\nif err != nil {\npanic(err)\n}\n}\n"
	ret += "\n"
	ret += "func YourFunction(data interface{}) {\n  fmt.Printf(\"%+v\", data)\n}\n"

	for _, v := range el.Data {
		ret += v.GetDefStructFromComplexTable() + "\n"
	}

	for _, v := range el.Data {
		ret += v.GetDefStruct() + "\n"
	}

	ret += el.GetCodeFromPrimaryKey(tableName)

	ret += "\n}\n"

	//err, fk = el.getForeignKeyColumns(tableName)
	if err != nil {
		panic(err)
	}

	for _, v := range el.Data {
		for _, d := range v.Data {
			if d.IsForeignKey == true {
				ret += el.GetCodeFromForeignKey(d.ReferenceTableName, d)
			}
		}
	}

	return ret
}

func (el *Database) New(db *sql.DB, ctx context.Context) error {
	el.db = db
	el.ctx = ctx
	el.Data = make(map[string]Table)

	var err error
	var listOfTables []string

	err, listOfTables = el.listTables()
	if err != nil {
		return err
	}

	for _, tableName := range listOfTables {
		err, el.Data[tableName] = el.mounteTableData(tableName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (el Database) listTables() (error, []string) {
	var returnList = make([]string, 0)
	var err error
	var queryReturn *sql.Rows

	queryReturn, err = el.db.QueryContext(el.ctx, "SELECT name FROM SYSOBJECTS WHERE xtype = 'U';")
	if err != nil {
		return err, []string{}
	}
	for queryReturn.Next() {
		var colNameFromQuery string
		err = queryReturn.Scan(&colNameFromQuery)
		if err != nil {
			return err, []string{}
		}

		returnList = append(returnList, colNameFromQuery)
	}

	return nil, returnList
}

func (el Database) mounteTableData(tableName string) (error, Table) {
	var err error
	var typesList Table
	var primaryKeyList map[string]PrimaryKeyRelation
	var foreignKeyList map[string]ForeignKeyRelation

	err, primaryKeyList = el.getPrimaryKeysFromTable(tableName)
	if err != nil {
		return err, Table{}
	}

	err, foreignKeyList = el.getForeignKeyColumns(tableName)
	if err != nil {
		return err, Table{}
	}

	err, typesList = el.getFieldsTypeFromTable(tableName, primaryKeyList, foreignKeyList)
	if err != nil {
		return err, Table{}
	}

	return nil, typesList
}

func (el Database) getForeignKeyColumns(tableName string) (error, map[string]ForeignKeyRelation) {
	var returnList = make(map[string]ForeignKeyRelation)
	var err error
	var queryReturn *sql.Rows

	queryReturn, err = el.db.QueryContext(el.ctx, fmt.Sprintf(`SELECT
      f.name AS foreign_key_name
     ,OBJECT_NAME(f.parent_object_id) AS table_name
     ,COL_NAME(fc.parent_object_id, fc.parent_column_id) AS constraint_column_name
     ,OBJECT_NAME (f.referenced_object_id) AS referenced_object
     ,COL_NAME(fc.referenced_object_id, fc.referenced_column_id) AS referenced_column_name
     ,is_disabled
     ,delete_referential_action_desc
     ,update_referential_action_desc
  FROM sys.foreign_keys AS f
  INNER JOIN sys.foreign_key_columns AS fc
     ON f.object_id = fc.constraint_object_id
  WHERE f.parent_object_id = OBJECT_ID('%v') ORDER BY foreign_key_name ASC;`, tableName))
	if err != nil {
		return err, nil
	}

	for queryReturn.Next() {

		var line ForeignKeyRelation
		err = queryReturn.Scan(
			&line.ForeignKeyName,
			&line.TableName,
			&line.ConstraintColumnName,
			&line.ReferencedObject,
			&line.ReferencedColumnName,
			&line.IsDisabled,
			&line.DeleteReferentialActionDescription,
			&line.UpdateReferentialActionDescription,
		)
		if err != nil {
			return err, nil
		}

		returnList[line.ConstraintColumnName] = line
	}

	return nil, returnList
}

// Return a list of primary keys
// Format:
// map[tableName]
//   TableName: nome da tabela. Ex.: user
//   ColumnName: nome do campo. Ex.: id
//   ConstraintName: nome do grupo. Ex.: dbo
func (el Database) getPrimaryKeysFromTable(tableName string) (error, map[string]PrimaryKeyRelation) {
	var returnList = make(map[string]PrimaryKeyRelation)
	var err error
	var queryReturn *sql.Rows

	queryReturn, err = el.db.QueryContext(
		el.ctx,
		fmt.Sprintf(
			`SELECT 
	tab.TABLE_NAME AS tableName, 
	col.COLUMN_NAME AS columnName, 
	col.CONSTRAINT_SCHEMA AS constraintName
FROM 
    INFORMATION_SCHEMA.TABLE_CONSTRAINTS Tab, 
    INFORMATION_SCHEMA.CONSTRAINT_COLUMN_USAGE Col 
WHERE
    Col.Constraint_Name = Tab.Constraint_Name
    AND Col.Table_Name = Tab.Table_Name
    AND Constraint_Type = 'PRIMARY KEY'
    AND Col.Table_Name = '%v' ORDER BY columnName ASC`,
			tableName,
		),
	)
	if err != nil {
		return err, nil
	}

	for queryReturn.Next() {

		var line PrimaryKeyRelation
		err = queryReturn.Scan(
			&line.TableName,
			&line.ColumnName,
			&line.ConstraintName,
		)
		if err != nil {
			return err, nil
		}

		returnList[line.ColumnName] = line
	}

	return nil, returnList
}

func (el Database) getFieldsTypeFromTable(tableName string, primaryKeyList map[string]PrimaryKeyRelation, foreingKeyList map[string]ForeignKeyRelation) (error, Table) {
	var err error
	var tableNameWithRules string
	var fieldNameWithRules string
	var referenceTableNameWithRules string
	var referenceFieldNameWithRules string
	var queryReturn *sql.Rows
	var ret Table
	ret.Data = make(map[string]ColumnType)
	var line []*sql.ColumnType
	var isPrimaryKey bool
	var isForeignKey bool
	//var primaryKeyRelation PrimaryKeyRelation
	var foreignKeyRelation ForeignKeyRelation

	err, tableNameWithRules = NameRules(tableName)
	if err != nil {
		return err, Table{}
	}

	queryReturn, err = el.db.QueryContext(el.ctx, fmt.Sprintf("SELECT * FROM [%v];", tableName))
	if err != nil {
		return err, Table{}
	}
	defer queryReturn.Close()

	line, err = queryReturn.ColumnTypes()
	if err != nil {
		return err, Table{}
	}

	for _, value := range line {

		fieldName := value.Name()
		err, fieldNameWithRules = NameRules(fieldName)
		if err != nil {
			return err, Table{}
		}

		_, isPrimaryKey = primaryKeyList[fieldName]
		foreignKeyRelation, isForeignKey = foreingKeyList[fieldName]

		if isForeignKey == true {
			err, referenceTableNameWithRules = NameRules(foreignKeyRelation.ReferencedObject)
			if err != nil {
				return err, Table{}
			}

			err, referenceFieldNameWithRules = NameRules(foreignKeyRelation.ConstraintColumnName)
			if err != nil {
				return err, Table{}
			}
		} else {
			referenceTableNameWithRules = ""
			referenceFieldNameWithRules = ""
		}

		decimalSizePrecision, decimalSizeScale, decimalSizeOkToUse := value.DecimalSize()
		length, lengthOkToUse := value.Length()
		nullable, nullableOkToUse := value.Nullable()
		var toAppend = ColumnType{
			FieldName:                   fieldName,
			FieldNameWithRules:          fieldNameWithRules,
			TableName:                   tableName,
			TableNameWithRules:          tableNameWithRules,
			ReferenceTableName:          foreignKeyRelation.ReferencedObject,
			ReferenceTableNameWithRules: referenceTableNameWithRules,
			ReferenceFieldName:          foreignKeyRelation.ConstraintColumnName,
			ReferenceFieldNameWithRules: referenceFieldNameWithRules,
			SqlType:                     value.DatabaseTypeName(),
			TypeAsString:                value.ScanType().String(),
			DecimalSizePrecision:        decimalSizePrecision,
			DecimalSizeScale:            decimalSizeScale,
			DecimalSizeOkToUse:          decimalSizeOkToUse,
			Length:                      length,
			LengthOkToUse:               lengthOkToUse,
			Nullable:                    nullable,
			NullableOkToUse:             nullableOkToUse,
			ScanType:                    value.ScanType(),
			IsPrimaryKey:                isPrimaryKey,
			IsForeignKey:                isForeignKey,
		}

		ret.Data[fieldName] = toAppend
	}

	ret.Ref = &el

	return err, ret
}
