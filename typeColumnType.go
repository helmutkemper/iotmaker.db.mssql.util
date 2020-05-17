package iotmaker_db_mssql_util

import "reflect"

type ColumnType struct {
	Name                 string
	NameWithRule         string
	DatabaseTypeName     string
	DecimalSizePrecision int64
	DecimalSizeScale     int64
	DecimalSizeOkToUse   bool
	Length               int64
	LengthOkToUse        bool
	Nullable             bool
	NullableOkToUse      bool
	IsPrimaryKey         bool
	ScanType             reflect.Type
}
