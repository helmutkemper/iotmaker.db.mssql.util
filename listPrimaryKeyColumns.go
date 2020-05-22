package iotmaker_db_mssql_util

import (
	"database/sql"
	"fmt"
)

func (el GoToMSSqlCode) ListPrimaryKeyColumns(tableName string) (error, map[string]PrimaryKeyRelation) {
	var returnList = make(map[string]PrimaryKeyRelation)
	var err error
	var queryReturn *sql.Rows

	queryReturn, err = el.Db.QueryContext(el.Ctx, fmt.Sprintf(`SELECT 
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
    AND Col.Table_Name = '%v' ORDER BY columnName ASC`, tableName))
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
