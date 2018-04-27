package main

import (
	"fmt"
	"flag"
	"os"

	"github.com/upperwal/xgdb/gdb"
)
func main() {

	progName := flag.String("attach", "", "Program name to attach xGDB")
	flag.Parse()

	if *progName == "" {
		fmt.Println("Give a program name to attach gdb")
		fmt.Println("Ex: xgdb -attach main")
		os.Exit(1)
	}

	gdb, err := gdb.NewGDB(*progName)
	if err != nil {
		panic(err)
	}

	gdb.Start()

	select {}
}
