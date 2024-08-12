package types

import (
	"flag"
	"log"
)

func Comment() {
	vfs := flag.NewFlagSet("n", flag.PanicOnError)
	err := vfs.Parse([]string{})
	if err != nil {
		log.Fatal(err)
	}
}
