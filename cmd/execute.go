package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Execute is an example execution function to be called from main package.
// It should call ExecuteE and only deal with error in application specific way.
func Execute(ctx context.Context, vip *viper.Viper, root *cobra.Command) {
	if err := ExecuteE(ctx, vip, root); err != nil {
		log.Fatal(err)
	}
}

// ExecuteE is an example execution function to be called from tests.
func ExecuteE(ctx context.Context, vip *viper.Viper, root *cobra.Command) error {
	setupRootCmd(root, vip)
	cmdCtx := context.WithValue(ctx, viperCtxKey, vip)
	if err := root.ExecuteContext(cmdCtx); err != nil {
		return err
	}
	return nil
}
