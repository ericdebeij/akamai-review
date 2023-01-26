/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"
	"os"

	"github.com/ericdebeij/akamai-review/v2/cmd"
)

func setlogfile(filename string) (file *os.File, err error) {
	if filename != "" {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
		file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		log.SetOutput(file)
	}
	return
}

func main() {
	setlogfile(os.Getenv("AK_LOG_PATH"))
	cmd.Execute()
}
