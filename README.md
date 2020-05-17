# iotmaker.db.mssql.util

Este é um script preliminar para ajudar a exportar dados do mssql.

A ideia básica é simples, você aponta um banco de dados e ele monta um código golang

```sql
USE [kemper2]
GO

/****** Object:  Table [dbo].[test]    Script Date: 17/05/2020 20:48:31 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[test](
	[id] [bigint] NOT NULL,
	[nome] [nchar](255) NULL,
 CONSTRAINT [PK_test] PRIMARY KEY CLUSTERED 
(
	[id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO
```

```sql
USE [kemper2]
GO

/****** Object:  Table [dbo].[relacao]    Script Date: 17/05/2020 20:45:38 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[relacao](
	[id] [bigint] NOT NULL,
	[id_test] [bigint] NULL,
	[name] [nchar](255) NULL,
 CONSTRAINT [PK_relacao] PRIMARY KEY CLUSTERED 
(
	[id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO

ALTER TABLE [dbo].[relacao]  WITH CHECK ADD  CONSTRAINT [FK_relacao_test] FOREIGN KEY([id_test])
REFERENCES [dbo].[test] ([id])
GO

ALTER TABLE [dbo].[relacao] CHECK CONSTRAINT [FK_relacao_test]
GO
```

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