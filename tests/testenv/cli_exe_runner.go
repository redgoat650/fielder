package testenv

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	FielderExeEnv = "FIELDER_EXE"
)

type ExeRunner struct {
	Exe     string
	TestDir string
}

func NewExeRunner(t *testing.T) *ExeRunner {
	exe := os.Getenv(FielderExeEnv)
	if exe == "" {
		exe = "fielder"
	}

	testDir := t.TempDir()

	return &ExeRunner{
		Exe:     filepath.FromSlash(exe),
		TestDir: testDir,
	}
}

func (r *ExeRunner) Run(t *testing.T, expectedErr bool, args ...string) (string, error) {
	t.Logf("running 'fielder %v'", strings.Join(args, " "))

	cmd := exec.Command(r.Exe, append([]string{
		"--config", filepath.Join(r.TestDir, ".fielder.yaml"),
	}, args...)...)

	out, err := cmd.CombinedOutput()

	outStr := string(out)
	if expectedErr {
		require.Error(t, err, "unexpected success\ncombinedOut: %s\nerr: %s\n", outStr, err)
	} else {
		require.NoError(t, err, "unexpected error\ncombinedOut: %s\nerr: %s\n", outStr, err)
	}

	return string(out), err
}
