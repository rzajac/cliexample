package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NakedRoot is an example application without subcommands.
func NakedRoot(version string) *cobra.Command {
	root := &cobra.Command{}

	// Setup flags.
	root.Flags().IntP("left", "l", 0, "left number")
	root.Flags().IntP("right", "r", 0, "right number")
	root.Flags().BoolP("hex", "x", false, "display as hex")

	// Root command configuration.
	root.Use = "naked"
	root.Short = "Example naked command."
	root.Version = version + "\n"
	root.SilenceUsage = true
	root.PersistentPreRunE = rootCmdPersistentPreRunE
	root.PreRunE = func(cmd *cobra.Command, _ []string) error {
		vip := cmd.Context().Value(viperCtxKey).(*viper.Viper)
		return vip.BindPFlags(cmd.Flags())
	}
	root.RunE = func(cmd *cobra.Command, args []string) error {
		vip := cmd.Context().Value(viperCtxKey).(*viper.Viper)

		r := vip.GetInt("right")
		l := vip.GetInt("left")

		format := "%d + %d = %d\n"
		if vip.GetBool("hex") {
			format = "%X + %X = %X\n"
		}
		cmd.Printf(format, l, r, l+r)

		return nil
	}
	root.SetVersionTemplate(`{{.Version}}`)

	// Configure environment variable prefix.
	annotate(root, envVarPrefix, "CLI")

	return root
}
