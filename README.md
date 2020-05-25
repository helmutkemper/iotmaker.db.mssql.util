# iotmaker.db.mssql.util

Esse pacote golang tem por finalidade ajudar a exportar dados contidos em banco de dados 
MS SQL Server e foi idealizado para ajudar na criação de códigos Golang.

Basicamente, basta apontar uma tabela e o código vai escrever o código básico para a 
acessar e exportar a tabela linha a linha em forma de um struct populado.

Veja um exemplo abaixo:

```sql
USE [toExport]
GO

SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[work_phone](
    [id] [bigint] NOT NULL,
    [number] [nchar](20) NULL,
    CONSTRAINT [PK_work_phone] PRIMARY KEY CLUSTERED
       (
        [id] ASC
           )WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO

CREATE TABLE [dbo].[home_phone](
	[id] [bigint] NOT NULL,
	[number] [nchar](20) NULL,
 CONSTRAINT [PK_home_phone] PRIMARY KEY CLUSTERED
(
	[id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO

CREATE TABLE [dbo].[mobile_phone](
     [id] [bigint] NOT NULL,
     [number] [nchar](20) NULL,
     CONSTRAINT [PK_mobile_phone] PRIMARY KEY CLUSTERED
         (
          [id] ASC
             )WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO

CREATE TABLE [dbo].[user](
     [id] [bigint] NOT NULL,
     [id_mobile_phone] [bigint] NULL,
     [id_home_phone] [bigint] NULL,
     [id_work_phone] [bigint] NULL,
     [name] [nchar](255) NULL,
     CONSTRAINT [PK_user] PRIMARY KEY CLUSTERED
         (
          [id] ASC
             )WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO

ALTER TABLE [dbo].[user]  WITH CHECK ADD  CONSTRAINT [FK_user_home_phone] FOREIGN KEY([id_home_phone])
    REFERENCES [dbo].[home_phone] ([id])
GO

ALTER TABLE [dbo].[user] CHECK CONSTRAINT [FK_user_home_phone]
GO

ALTER TABLE [dbo].[user]  WITH CHECK ADD  CONSTRAINT [FK_user_mobile_phone] FOREIGN KEY([id_mobile_phone])
    REFERENCES [dbo].[mobile_phone] ([id])
GO

ALTER TABLE [dbo].[user] CHECK CONSTRAINT [FK_user_mobile_phone]
GO

ALTER TABLE [dbo].[user]  WITH CHECK ADD  CONSTRAINT [FK_user_work_phone] FOREIGN KEY([id_work_phone])
    REFERENCES [dbo].[work_phone] ([id])
GO

ALTER TABLE [dbo].[user] CHECK CONSTRAINT [FK_user_work_phone]
GO

INSERT INTO [dbo].[home_phone]([id],[number])VALUES(1,'home 1')
GO
INSERT INTO [dbo].[home_phone]([id],[number])VALUES(2,'home 2')
GO
INSERT INTO [dbo].[home_phone]([id],[number])VALUES(3,'home 3')
GO


INSERT INTO [dbo].[mobile_phone]([id],[number])VALUES(1,'mobile 1')
GO
INSERT INTO [dbo].[mobile_phone]([id],[number])VALUES(2,'mobile 2')
GO
INSERT INTO [dbo].[mobile_phone]([id],[number])VALUES(3,'mobile 3')
GO


INSERT INTO [dbo].[work_phone]([id],[number])VALUES(1,'work 1')
GO
INSERT INTO [dbo].[work_phone]([id],[number])VALUES(2,'work 2')
GO
INSERT INTO [dbo].[work_phone]([id],[number])VALUES(3,'work 3')
GO

INSERT INTO [dbo].[user]([id],[id_mobile_phone],[id_home_phone],[id_work_phone],[name])VALUES(1,1,1,1,'home 1; phone 1; work 1')
INSERT INTO [dbo].[user]([id],[id_mobile_phone],[id_home_phone],[id_work_phone],[name])VALUES(2,2,2,2,'home 2; phone 2; work 2')
INSERT INTO [dbo].[user]([id],[id_mobile_phone],[id_home_phone],[id_work_phone],[name])VALUES(3,3,3,3,'home 3; phone 3; work 3')
GO
```

```golang
package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	mssqlUtil "github.com/helmutkemper/iotmaker.db.mssql.util"
)

var (
	ctx context.Context
)

func main() {
	var err error
	var db *sql.DB

	ctx = context.Background()

	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s",
		"localhost",
		1434,
		"toExport",
		"____user_name____",
		"____password____",
  )
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		panic(err)
	}
	err = db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	var tbl mssqlUtil.Database
	tbl.New(db, ctx)
	fmt.Printf("%v", tbl.ToCode("user", "main"))
}
```
Saída:
```golang
package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	var err error
	var db *sql.DB
	var ctx context.Context

	ctx = context.Background()

	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s",
		"localhost",
		1434,
		"toExport",
		"CS\\helmut.kemper",
		"temp@123",
	)

	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		panic(err)
	}
	err = db.PingContext(ctx)
	if err != nil {
		panic(err)
	}
	err = GetUser(db, ctx)
	if err != nil {
		panic(err)
	}
}

func YourFunction(data interface{}) {
	fmt.Printf("%+v", data)
}

type DataUser struct {
	Id          int64
	HomePhone   []DataHomePhone
	MobilePhone []DataMobilePhone
	WorkPhone   []DataWorkPhone
	Name        string
}

type DataMobilePhone struct {
	Id     int64
	Number string
}

type DataHomePhone struct {
	Id     int64
	Number string
}

type DataWorkPhone struct {
	Id     int64
	Number string
}

type TableMobilePhone struct {
	Id     int64
	Number string
}

type TableHomePhone struct {
	Id     int64
	Number string
}

type TableWorkPhone struct {
	Id     int64
	Number string
}

type TableUser struct {
	Id            int64
	IdHomePhone   int64
	IdMobilePhone int64
	IdWorkPhone   int64
	Name          string
}

func GetUser(db *sql.DB, ctx context.Context) error {
	var err error
	var RowUser *sql.Rows
	RowUser, err = db.QueryContext(ctx, "SELECT id, id_home_phone, id_mobile_phone, id_work_phone, name FROM [user]")
	if err != nil {
		return err
	}
	if RowUser.Next() {
		var UserColumnId int64
		var UserColumnIdHomePhone int64
		var UserColumnIdMobilePhone int64
		var UserColumnIdWorkPhone int64
		var UserColumnName string
		err = RowUser.Scan(&UserColumnId, &UserColumnIdHomePhone, &UserColumnIdMobilePhone, &UserColumnIdWorkPhone, &UserColumnName)
		if err != nil {
			return err
		}
		var data = DataUser{
			Id:          UserColumnId,
			HomePhone:   GetHomePhoneAsArray(db, ctx, UserColumnIdHomePhone),
			MobilePhone: GetMobilePhoneAsArray(db, ctx, UserColumnIdMobilePhone),
			WorkPhone:   GetWorkPhoneAsArray(db, ctx, UserColumnIdWorkPhone),
			Name:        UserColumnName,
		}

		YourFunction(data)
	}
	return nil

}

func GetMobilePhoneAsArray(db *sql.DB, ctx context.Context, UserColumnIdMobilePhone int64) []DataMobilePhone {
	var err error
	var ArrayOfTableMobilePhone = make([]DataMobilePhone, 0)
	var RowMobilePhone *sql.Rows
	RowMobilePhone, err = db.QueryContext(ctx, fmt.Sprintf("SELECT id, number FROM [mobile_phone] WHERE Id = %v", UserColumnIdMobilePhone))
	if err != nil {
		panic(err)
	}
	defer RowMobilePhone.Close()
	if RowMobilePhone.Next() {
		var MobilePhoneColumnId int64
		var MobilePhoneColumnNumber string
		err = RowMobilePhone.Scan(&MobilePhoneColumnId, &MobilePhoneColumnNumber)
		if err != nil {
			panic(err)
		}
		var data = DataMobilePhone{
			Id:     MobilePhoneColumnId,
			Number: MobilePhoneColumnNumber,
		}

		ArrayOfTableMobilePhone = append(ArrayOfTableMobilePhone, data)
	}
	return ArrayOfTableMobilePhone
}

func GetHomePhoneAsArray(db *sql.DB, ctx context.Context, UserColumnIdHomePhone int64) []DataHomePhone {
	var err error
	var ArrayOfTableHomePhone = make([]DataHomePhone, 0)
	var RowHomePhone *sql.Rows
	RowHomePhone, err = db.QueryContext(ctx, fmt.Sprintf("SELECT id, number FROM [home_phone] WHERE Id = %v", UserColumnIdHomePhone))
	if err != nil {
		panic(err)
	}
	defer RowHomePhone.Close()
	if RowHomePhone.Next() {
		var HomePhoneColumnId int64
		var HomePhoneColumnNumber string
		err = RowHomePhone.Scan(&HomePhoneColumnId, &HomePhoneColumnNumber)
		if err != nil {
			panic(err)
		}
		var data = DataHomePhone{
			Id:     HomePhoneColumnId,
			Number: HomePhoneColumnNumber,
		}

		ArrayOfTableHomePhone = append(ArrayOfTableHomePhone, data)
	}
	return ArrayOfTableHomePhone
}

func GetWorkPhoneAsArray(db *sql.DB, ctx context.Context, UserColumnIdWorkPhone int64) []DataWorkPhone {
	var err error
	var ArrayOfTableWorkPhone = make([]DataWorkPhone, 0)
	var RowWorkPhone *sql.Rows
	RowWorkPhone, err = db.QueryContext(ctx, fmt.Sprintf("SELECT id, number FROM [work_phone] WHERE Id = %v", UserColumnIdWorkPhone))
	if err != nil {
		panic(err)
	}
	defer RowWorkPhone.Close()
	if RowWorkPhone.Next() {
		var WorkPhoneColumnId int64
		var WorkPhoneColumnNumber string
		err = RowWorkPhone.Scan(&WorkPhoneColumnId, &WorkPhoneColumnNumber)
		if err != nil {
			panic(err)
		}
		var data = DataWorkPhone{
			Id:     WorkPhoneColumnId,
			Number: WorkPhoneColumnNumber,
		}

		ArrayOfTableWorkPhone = append(ArrayOfTableWorkPhone, data)
	}
	return ArrayOfTableWorkPhone
}
```