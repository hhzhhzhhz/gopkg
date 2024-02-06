//go:build linux
// +build linux

package utils

import (
	"bytes"
	"os/exec"
	"syscall"
	"time"
)

const (
	sh = "/bin/sh"
)

const (
	DefaultTimeout = 1 * time.Minute
)

// Command Compatible with complex commands by executing shell
func Command(c string, timeout time.Duration, args ...string) (string, error) {
	cmd := exec.Command(c, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	time.AfterFunc(timeout,
		func() {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		})

	if err := cmd.Start(); err != nil {
		return string(out.Bytes()), err
	}

	if err := cmd.Wait(); err != nil {
		return string(out.Bytes()), err
	}
	return string(out.Bytes()), nil
}
