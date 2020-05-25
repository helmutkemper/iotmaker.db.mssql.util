package iotmaker_db_mssql_util

type ForeignKeyRelation struct {
	ForeignKeyName                     string
	TableName                          string
	ConstraintColumnName               string
	ReferencedObject                   string
	ReferencedColumnName               string
	IsDisabled                         string
	DeleteReferentialActionDescription string
	UpdateReferentialActionDescription string
}
