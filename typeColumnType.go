package iotmaker_db_mssql_util

import "reflect"

type ColumnType struct {
	FieldName                   string
	FieldNameWithRules          string
	TableName                   string
	TableNameWithRules          string
	ReferenceTableName          string
	ReferenceTableNameWithRules string
	ReferenceFieldName          string
	ReferenceFieldNameWithRules string
	SqlType                     string
	DecimalSizePrecision        int64
	DecimalSizeScale            int64
	DecimalSizeOkToUse          bool
	Length                      int64
	LengthOkToUse               bool
	Nullable                    bool
	NullableOkToUse             bool
	IsPrimaryKey                bool
	IsForeignKey                bool
	ScanType                    reflect.Type

	TypeAsString string
}

func (el ColumnType) GetFieldName() string {
	return el.FieldName
}

func (el ColumnType) GetFieldNameWithRules() string {
	return el.FieldNameWithRules
}

func (el ColumnType) GetAsDefVar() string {
	return "var " + el.TableNameWithRules + "Column" + el.FieldNameWithRules + "  " + el.GetTypeAsString()
}

func (el ColumnType) GetAsRefVar() string {
	return "&" + el.TableNameWithRules + "Column" + el.FieldNameWithRules
}

func (el ColumnType) GetAsVarName() string {
	return el.TableNameWithRules + "Column" + el.FieldNameWithRules
}

func (el ColumnType) GetAsDataVarName() string {
	return "Data" + el.GetFieldNameWithRules()
}

func (el ColumnType) GetTableName() string {
	return el.TableName
}

func (el ColumnType) GetTableNameWithRules() string {
	return el.TableNameWithRules
}

func (el ColumnType) GetReferenceTableName() string {
	return el.ReferenceTableName
}

func (el ColumnType) GetReferenceTableNameWithRules() string {
	return el.ReferenceTableNameWithRules
}

func (el ColumnType) GetReferenceFieldName() string {
	return el.ReferenceFieldName
}

func (el ColumnType) GetReferenceFieldNameWithRules() string {
	return el.ReferenceFieldNameWithRules
}

func (el ColumnType) GetTypeAsString() string {
	return el.TypeAsString
}
