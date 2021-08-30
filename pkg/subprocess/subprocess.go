package subprocess

import (
	"io"
	"os"
	"os/exec"
)

func Run(workingDirectory string, command string, arg ...string) (process os.Process, stdin io.WriteCloser, stdout io.ReadCloser, runError error) {
	cmd := exec.Command(command, arg...)
	cmd.Dir = workingDirectory
	cmd.Stderr = cmd.Stdout

	stdin, _ = cmd.StdinPipe()
	stdout, _ = cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		runError = err
	}
	process = *cmd.Process
	return
}
