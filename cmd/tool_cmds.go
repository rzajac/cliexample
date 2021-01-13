package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ToolRoot returns root command for the tool application.
func ToolRoot(version string) *cobra.Command {
	// Root command configuration.
	root := &cobra.Command{
		Use:               "tool",
		Short:             "Example tool command.",
		Version:           version + "\n",
		SilenceUsage:      true,
		PersistentPreRunE: rootCmdPersistentPreRunE,
	}
	root.SetVersionTemplate(`{{.Version}}`)

	// Add sub-commands.
	root.AddCommand(
		ToolAdd(),
		ToolAddEnv(),
		ToolSpin(),
	)

	// Configure environment variable prefix.
	annotate(root, envVarPrefix, "CLI")

	return root
}

// ToolAdd is example command adding two numbers.
func ToolAdd() *cobra.Command {
	cmd := &cobra.Command{}

	// Setup flags.
	cmd.Flags().IntP("left", "l", 0, "left number")
	cmd.Flags().IntP("right", "r", 0, "right number")
	cmd.Flags().BoolP("hex", "x", false, "display as hex")

	// Mark required flags.
	_ = cmd.MarkFlagRequired("left")
	_ = cmd.MarkFlagRequired("right")

	// Setup command.
	cmd.Use = "add"
	cmd.Short = "Add two numbers"
	cmd.PreRunE = bindCommandFlags
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
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

	return cmd
}

// ToolAddEnv is example command adding two numbers, using env file
// for missing flags.
func ToolAddEnv() *cobra.Command {
	cmd := &cobra.Command{}
	annotate(cmd, needsEnvFile, "")

	// Setup flags.
	cmd.Flags().IntP("left", "l", 0, "left number")
	cmd.Flags().IntP("right", "r", 0, "right number")
	cmd.Flags().BoolP("hex", "x", false, "display as hex")

	// Setup command.
	cmd.Use = "add-env"
	cmd.Short = "Adds two numbers using arguments, env file or both"
	cmd.PreRunE = bindCommandFlags
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		vip := cmd.Context().Value(viperCtxKey).(*viper.Viper)

		l := vip.GetInt("left")
		r := vip.GetInt("right")

		format := "%d + %d = %d\n"
		if vip.GetBool("hex") {
			format = "%X + %X = %X\n"
		}
		cmd.Printf(format, l, r, l+r)
		return nil
	}

	return cmd
}

// ToolSpin simulates long running command.
func ToolSpin() *cobra.Command {
	cmd := &cobra.Command{}

	// Setup command.
	cmd.Use = "spin"
	cmd.Short = "Spins forever"
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		for {
			cmd.Println("spin")
			select {
			case <-ctx.Done():
				err := ctx.Err()
				cmd.Printf("ctx err: %s\n", err)
				return err
			}
		}
	}

	return cmd
}
