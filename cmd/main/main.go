package main

import (
	"fmt"
	lvm "github.com/masahiro331/go-lvm"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	info, _ := f.Stat()

	r := io.NewSectionReader(f, 0, info.Size())
	d, err := lvm.NewDriver(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(d)
}
