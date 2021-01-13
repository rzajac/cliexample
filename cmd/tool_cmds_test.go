package cmd

import (
	"bytes"
	"context"
	"testing"
	"time"

	kit "github.com/rzajac/testkit"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_ToolAdd_display_as_base_ten(t *testing.T) {
	// --- Given ---
	ctx, cxl := context.WithCancel(context.Background())
	defer cxl()

	vip := viper.New()
	root := ToolRoot("0.0.0")

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"add", "-l", "12", "-r", "30"})

	// --- When ---
	err := ExecuteE(ctx, vip, root)

	// --- Then ---
	assert.NoError(t, err)
	assert.Exactly(t, "12 + 30 = 42\n", buf.String())
}

func Test_ToolAdd_display_as_hex(t *testing.T) {
	// --- Given ---
	ctx, cxl := context.WithCancel(context.Background())
	defer cxl()

	vip := viper.New()
	root := ToolRoot("0.0.0")

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"add", "-l", "12", "-r", "30", "-x"})

	// --- When ---
	err := ExecuteE(ctx, vip, root)

	// --- Then ---
	assert.NoError(t, err)
	assert.Exactly(t, "C + 1E = 2A\n", buf.String())
}

func Test_ToolAdd_display_as_hex_via_environment(t *testing.T) {
	// --- Given ---
	kit.SetEnv(t, "CLI_HEX", "true")

	ctx, cxl := context.WithCancel(context.Background())
	defer cxl()

	vip := viper.New()
	root := ToolRoot("0.0.0")

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"add", "-l", "12", "-r", "30"})

	// --- When ---
	err := ExecuteE(ctx, vip, root)

	// --- Then ---
	assert.NoError(t, err)
	assert.Exactly(t, "C + 1E = 2A\n", buf.String())
}

func Test_ToolAdd_setting_env_file_has_no_effect(t *testing.T) {
	// --- Given ---
	ctx, cxl := context.WithCancel(context.Background())
	defer cxl()

	vip := viper.New()
	root := ToolRoot("0.0.0")

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"add", "-l", "10", "-r", "30", "-e", "testdata/.env"})

	// --- When ---
	err := ExecuteE(ctx, vip, root)

	// --- Then ---
	assert.NoError(t, err)
	assert.Exactly(t, "10 + 30 = 40\n", buf.String())
}

func Test_ToolAddEnv(t *testing.T) {
	// --- Given ---
	ctx, cxl := context.WithCancel(context.Background())
	defer cxl()

	vip := viper.New()
	root := ToolRoot("0.0.0")

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"add-env", "-e", "testdata/.env"})

	// --- When ---
	err := ExecuteE(ctx, vip, root)

	// --- Then ---
	assert.NoError(t, err)
	assert.Exactly(t, "32 + 10 = 42\n", buf.String())
}

func Test_ToolAddEnv_override_env_file_arg(t *testing.T) {
	// --- Given ---
	ctx, cxl := context.WithCancel(context.Background())
	defer cxl()

	vip := viper.New()
	root := ToolRoot("0.0.0")

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"add-env", "-e", "testdata/.env", "-x"})

	// --- When ---
	err := ExecuteE(ctx, vip, root)

	// --- Then ---
	assert.NoError(t, err)
	assert.Exactly(t, "20 + A = 2A\n", buf.String())
}

func Test_ToolAddEnv_override_env_file_environment(t *testing.T) {
	// --- Given ---
	kit.SetEnv(t, "CLI_HEX", "true")

	ctx, cxl := context.WithCancel(context.Background())
	defer cxl()

	vip := viper.New()
	root := ToolRoot("0.0.0")

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"add-env", "-e", "testdata/.env"})

	// --- When ---
	err := ExecuteE(ctx, vip, root)

	// --- Then ---
	assert.NoError(t, err)
	assert.Exactly(t, "20 + A = 2A\n", buf.String())
}

func Test_ToolAddEnv_path_from_environment(t *testing.T) {
	// --- Given ---
	kit.SetEnv(t, "CLI_HEX", "true")

	ctx, cxl := context.WithCancel(context.Background())
	defer cxl()

	vip := viper.New()
	root := ToolRoot("0.0.0")

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"add-env", "-e", "testdata/.env"})

	// --- When ---
	err := ExecuteE(ctx, vip, root)

	// --- Then ---
	assert.NoError(t, err)
	assert.Exactly(t, "20 + A = 2A\n", buf.String())
}

func Test_ToolAddEnv_path_from_env(t *testing.T) {
	// --- Given ---
	kit.SetEnv(t, "CLI_ENV_FILE", "testdata/.env")

	ctx, cxl := context.WithCancel(context.Background())
	defer cxl()

	vip := viper.New()
	root := ToolRoot("0.0.0")

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"add-env"})

	// --- When ---
	err := ExecuteE(ctx, vip, root)

	// --- Then ---
	assert.NoError(t, err)
	assert.Exactly(t, "32 + 10 = 42\n", buf.String())
}

func Test_ToolSpin_cancel(t *testing.T) {
	// --- Given ---
	ctx, cxl := context.WithCancel(context.Background())
	time.AfterFunc(500*time.Millisecond, func() { cxl() })

	vip := viper.New()
	root := ToolRoot("0.0.0")

	bufOut := &bytes.Buffer{}
	bufErr := &bytes.Buffer{}
	root.SetOut(bufOut)
	root.SetErr(bufErr)
	root.SetArgs([]string{"spin"})

	// --- When ---
	err := ExecuteE(ctx, vip, root)

	// --- Then ---
	assert.ErrorIs(t, err, context.Canceled)
	assert.Exactly(t, "spin\nctx err: context canceled\n", bufOut.String())
	assert.Exactly(t, "Error: context canceled\n", bufErr.String())
}

func Test_version(t *testing.T) {
	// --- Given ---
	ctx, cxl := context.WithCancel(context.Background())
	defer cxl()

	vip := viper.New()
	root := ToolRoot("0.0.0")

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"--version"})

	// --- When ---
	err := ExecuteE(ctx, vip, root)

	// --- Then ---
	assert.NoError(t, err)
	assert.Exactly(t, "0.0.0\n", buf.String())
}
