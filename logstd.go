package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/aletheia7/ul"
)

var (
	subsystem      = flag.String("subsystem", "com.github.jum.logstd", "os_log subsystem")
	stdoutCategory = flag.String("stdout", "stdout", "os_log category for stdout")
	stderrCategory = flag.String("stderr", "stderr", "os_log category for stderr")
)

func main() {
	cmdLog := ul.New_object("com.github.jum.logstd", "cmdLog")
	defer cmdLog.Release()
	flag.CommandLine.SetOutput(cmdLog)
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintf(cmdLog, "Usage: logstd [flags] -- shell command\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	stdoutLog := ul.New_object(*subsystem, *stdoutCategory)
	defer stdoutLog.Release()
	stderrLog := ul.New_object(*subsystem, *stderrCategory)
	defer stderrLog.Release()

	cmd := exec.Command("sh", "-c", strings.Join(flag.Args(), " "))
	cmdLog.Debugf("about to exec %#v", cmd)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		cmdLog.Error(err.Error())
		os.Exit(1)
	}
	go func(l io.Writer, p io.Reader) {
		_, err := io.Copy(l, p)
		if err != nil {
			cmdLog.Error(err.Error())
		}
	}(stderrLog, stderr)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cmdLog.Error(err.Error())
		os.Exit(1)
	}
	go func(l io.Writer, p io.Reader) {
		_, err := io.Copy(l, p)
		if err != nil {
			cmdLog.Error(err.Error())
		}
	}(stdoutLog, stdout)

	if err := cmd.Start(); err != nil {
		cmdLog.Error(err.Error())
		os.Exit(1)
	}

	if err := cmd.Wait(); err != nil {
		cmdLog.Error(err.Error())
	}
	os.Exit(cmd.ProcessState.ExitCode())
}
