package iotmaker_db_mssql_util

import (
	"errors"
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
