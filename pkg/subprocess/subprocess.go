package subprocess

import (
	"io"
	"log"
	"os"
	"os/exec"
)

func Run(workingDirectory string, command string, arg ...string) (process *os.Process, stdin *io.WriteCloser, stdout *io.ReadCloser, runError error) {
	log.Println("Running", command, arg)
	cmd := exec.Command(command, arg...)
	cmd.Dir = workingDirectory
	cmd.Stderr = cmd.Stdout

	_stdin, _ := cmd.StdinPipe()
	stdin = &_stdin
	_stdout, _ := cmd.StdoutPipe()
	stdout = &_stdout

	if err := cmd.Start(); err != nil {
		runError = err
	}
	process = cmd.Process
	return
}
