# iotmaker.db.mssql.util

Este é um script preliminar para ajudar a exportar dados do mssql.

A ideia básica é simples, você aponta um banco de dados e ele monta um código golang

```golang
func main() {

  var err error
  var db *sql.DB
  var listOfTables []string
  var list2 map[string]mssqlUtil.ColumnType
  var toPrint mssqlUtil.ToMakeStruct

  ctx = mssqlUtil.GetContextBackground()

  connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s",
    "localhost",
    1434,
    "kemper2",
    "CS\\helmut.kemper",
    "temp@123")
  db, err = sql.Open("sqlserver", connString)
  if err != nil {
    panic(err)
  }
  fmt.Printf("%v\n", connString)
  err = db.PingContext(ctx)
  if err != nil {
    panic(err)
  }

  defer db.Close()

  var listOfColType map[string]mssqlUtil.ColumnType
  err, listOfTables = mssqlUtil.ListTables(db, ctx)
  if err != nil {
    panic(err)
  }
  for _, tableName := range listOfTables {

    err, listOfColType = mssqlUtil.ListColumnTypes(db, ctx, tableName)
    if err != nil {
      panic(err)
    }
    err, toPrint = mssqlUtil.ToStruct(db, ctx, tableName, listOfColType)
    if err != nil {
      panic(err)
    }

    fmt.Printf("%v\n", toPrint.MakeStructText("sql"))
  }
}
```
OutPut atual:
```golang
type TableTest struct {
  Nome  string `sql:"nome"`
  Id  int64 `sql:"id" primaryKey:"true"`
}

type TableRelacao struct {
  Id  int64 `sql:"id" primaryKey:"true"`
  Test  []TableTest `sql:"id_test" query:"SELECT id, nome FROM test WHERE id = %v" scan:"&columnId, &columnNome" vars:"  var columnId int64\n  var columnNome string\n"`
  Name  string `sql:"name"`
}

```