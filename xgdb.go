package main

import (
	"github.com/upperwal/xgdb/gdb"
)
func main() {

	gdb, err := gdb.NewGDB("main")
	if err != nil {
		panic(err)
	}

	gdb.Start()

	select {}
}
