package runner

import (
	"io"
	"os/exec"
	"strings"
)

func runBefore() bool {
	runnerLog("Run Before...")

	beforeCommand := buildBefore()
	runnerLog("before")
	if strings.TrimSpace(beforeCommand) == "" {
		return false
	}
	cmd := exec.Command(beforeCommand)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	io.Copy(appLogWriter{}, stderr)
	io.Copy(appLogWriter{}, stdout)
	return true
}

func runAfter() bool {
	runnerLog("Run After...")

	afterCommand := buildAfter()
	if strings.TrimSpace(afterCommand) == "" {
		return false
	}
	cmd := exec.Command(afterCommand)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	io.Copy(appLogWriter{}, stderr)
	io.Copy(appLogWriter{}, stdout)
	return true
}

func run() bool {
	runnerLog("Running...")

	cmd := exec.Command(buildPath())

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	go io.Copy(appLogWriter{}, stderr)
	go io.Copy(appLogWriter{}, stdout)

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		runnerLog("Killing PID %d", pid)
		cmd.Process.Kill()
	}()

	return true
}
