package main

import (
	"log"

	"github.com/ericdebeij/akamai-review/v3/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.RootCmd, "./")
	if err != nil {
		log.Fatal(err)
	}
}
