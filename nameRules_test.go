package iotmaker_db_mssql_util

import "fmt"

func ExampleNameRules1() {
	var name string
	var err error
	name = "666 olá Mundo_estranho_666_FK"
	err, name = NameRules(name)
	if err != nil {
		panic(err)
	}

	fmt.Printf("new name: %v\n", name)

	// Output:
	// new name: OlMundoEstranho666Fk
}

func ExampleNameRules2() {
	var name string
	var err error
	name = "666 áéíóú"
	err, name = NameRules(name)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("new name: %v\n", name)

	// Output:
	// error: name is empty
}
