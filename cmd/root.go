package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

// needsEnvFile is a command annotation requiring env file.
const needsEnvFile = "needs_env_file"

// envVarPrefix is a command annotation key holding environment variable prefix.
const envVarPrefix = "env_var_prefix"

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

// viperCtxKey is a context key. It can be used to access viper instance.
var viperCtxKey = &contextKey{"viper"}

// rootCmdPersistentPreRunE is persistent pre run function for root commands.
// It loads env file if command has needsEnvFile annotation
// and configures viper.
func rootCmdPersistentPreRunE(cmd *cobra.Command, _ []string) error {
	// Check if command needs env file.
	if _, ok := cmd.Annotations[needsEnvFile]; !ok {
		return nil
	}

	vip := cmd.Context().Value(viperCtxKey).(*viper.Viper)
	vip.SetConfigFile(vip.GetString("env-file"))
	vip.SetConfigType("env")
	if err := vip.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading env file, %w", err)
	}

	// We want env file variables to have the same names as
	// environmental variables (with prefix). I have not found a better
	// way to achieve this.
	if prefix, _ := cmd.Root().Annotations[envVarPrefix]; prefix != "" {
		fil, err := os.Open(vip.ConfigFileUsed())
		if err != nil {
			return fmt.Errorf("error reading env file, %w", err)
		}

		env, err := gotenv.StrictParse(fil)
		if err != nil {
			return fmt.Errorf("error parsing env file, %w", err)
		}

		for k := range env {
			if strings.HasPrefix(k, prefix) {
				nk := strings.Replace(k, prefix+"_", "", 1)
				nk = strings.ToLower(nk)
				vip.RegisterAlias(strings.ToLower(k), nk)
			}
		}
	}
	return nil
}

// setupRootCmd configures viper and root command.
func setupRootCmd(root *cobra.Command, vip *viper.Viper) {
	pfs := root.PersistentFlags()
	pfs.Bool("version", false, "version")
	pfs.StringP("env-file", "e", "./.env", "env file")

	_ = vip.BindPFlag("version", pfs.Lookup("version"))
	_ = vip.BindPFlag("env-file", pfs.Lookup("env-file"))

	prefix, _ := root.Annotations[envVarPrefix]
	vip.SetEnvPrefix(prefix)
	vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	vip.AutomaticEnv()
}

// annotate is convenience function to set annotation on the command.
func annotate(cmd *cobra.Command, key, val string) {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string, 1)
	}
	cmd.Annotations[key] = val
}
