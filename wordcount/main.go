package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	fmt.Println(args)

	if len(args) < 3 {
		fmt.Println("ccwc requires 2 args, path to infile and -c (byte flag)")
		os.Exit(1)
	}

	f, err := os.Open(args[1])
	checkErr(err)
	defer f.Close()

	fi, err := f.Stat()
	checkErr(err)

	fmt.Println(fi.Size(), args[1])
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
