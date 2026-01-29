package cmd

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
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}

	fmt.Println(fi.Size(), args[1])
}
