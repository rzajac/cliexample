package main

import (
	"context"

	"github.com/spf13/viper"

	"github.com/rzajac/cliexample/cmd"
)

func main() {
	cmd.Execute(context.Background(), viper.New(), cmd.ToolRoot("0.0.0"))
}
