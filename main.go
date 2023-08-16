/*
Copyright © 2022, 2023 NAME HERE hdebeij@akamai.com
*/
package main

import (
	"github.com/ericdebeij/akamai-review/v3/cmd"
)

func main() {
	cmd.Execute()
	cmd.Cleanup()
}
