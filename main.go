/*
Copyright © 2022, 2023 NAME HERE hdebeij@akamai.com
*/
package main

import (
	"time"

	"github.com/ericdebeij/akamai-review/v3/cmd"
)

func main() {
	cmd.Execute()
	cmd.Cleanup()

	// giving time to close all files correctly
	time.Sleep(time.Second * 3)
}
