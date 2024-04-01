package main

import (
	"perpus_api/db"
	"perpus_api/routes"
)

func main() {

	if err := db.Init(); err != nil {
        panic(err)
    }

	e := routes.Init()

	e.Logger.Fatal(e.Start(":1323"))
}
