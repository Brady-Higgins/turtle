package main

import (
	"github.com/Brady-Higgins/turtle/cmd"
	"github.com/Brady-Higgins/turtle/internal/config"
)

func main() {
	cmd.Execute()
	config.ReadConfig()
}
