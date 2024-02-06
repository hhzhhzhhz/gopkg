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
	c := exec.CommandContext(ctx, c, args)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out

	if err := c.Run(); err != nil {
		return string(out.Bytes()), err
	}
	return string(out.Bytes()), nil
}
