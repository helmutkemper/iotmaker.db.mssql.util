package iotmaker_db_mssql_util

import (
	"context"
	"database/sql"
	"fmt"
)

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

func ListForeignKeyColumns(db *sql.DB, ctx context.Context, tableName string) (error, []ForeignKeyRelation) {
	var returnList = make([]ForeignKeyRelation, 0)
	var err error
	var queryReturn *sql.Rows

	queryReturn, err = db.QueryContext(ctx, fmt.Sprintf(`SELECT
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
  WHERE f.parent_object_id = OBJECT_ID('%v');`, tableName))
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

		returnList = append(returnList, line)
	}

	return nil, returnList
}
