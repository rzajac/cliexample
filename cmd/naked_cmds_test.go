package cmd

import (
	"bytes"
	"context"
	"testing"

	kit "github.com/rzajac/testkit"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_Naked_display_as_base_ten(t *testing.T) {
	// --- Given ---
	ctx, cxl := context.WithCancel(context.Background())
	defer cxl()

	vip := viper.New()
	root := NakedRoot("0.0.0")

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"-l", "12", "-r", "30"})

	// --- When ---
	err := ExecuteE(ctx, vip, root)

	// --- Then ---
	assert.NoError(t, err)
	assert.Exactly(t, "12 + 30 = 42\n", buf.String())
}

func Test_Naked_display_right_from_env(t *testing.T) {
	// --- Given ---
	kit.SetEnv(t, "CLI_RIGHT", "30")

	ctx, cxl := context.WithCancel(context.Background())
	defer cxl()

	vip := viper.New()
	root := NakedRoot("0.0.0")

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"-l", "12"})

	// --- When ---
	err := ExecuteE(ctx, vip, root)

	// --- Then ---
	assert.NoError(t, err)
	assert.Exactly(t, "12 + 30 = 42\n", buf.String())
}
