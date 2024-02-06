//go:build windows
// +build windows

package utils

import (
	"bytes"
	"context"
	"os/exec"
	"syscall"
	"time"
)

const (
	exe = "C:\\WINDOWS\\system32\\cmd.exe"
)

const (
	DefaultTimeout = 1 * time.Minute
)

// Command Compatible with complex commands by executing shell
func Command(c string, timeout time.Duration, args ...string) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	// c := exec.CommandContext(ctx, exe, "/c", cmd)
	cmd := exec.CommandContext(ctx, c, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return string(out.Bytes()), err
	}
	return string(out.Bytes()), nil
}
