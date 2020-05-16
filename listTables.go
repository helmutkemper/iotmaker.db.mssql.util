package iotmaker_db_mssql_util

import (
	"context"
	"database/sql"
)

func ListTables(db *sql.DB, ctx context.Context) (error, []string) {
	var returnList = make([]string, 0)
	var err error
	var queryReturn *sql.Rows

	queryReturn, err = db.QueryContext(ctx, "SELECT name FROM SYSOBJECTS WHERE xtype = 'U';")
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
