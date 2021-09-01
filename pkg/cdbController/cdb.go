package cdbcontroller

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"

	"github.com/nonetype/gocdb/subprocess"
)

type Cdb struct {
	cdbPath   string
	process   *os.Process
	targetPID int
	stdin     *io.WriteCloser
	stdout    *io.ReadCloser
}

func NewCdb(cdbPath string) *Cdb {
	cdb := &Cdb{
		cdbPath:   cdbPath,
		process:   nil,
		targetPID: 0,
		stdin:     nil,
		stdout:    nil,
	}
	return cdb
}

func (c *Cdb) Run(programArgs ...string) (runError error) {
	c.process, c.stdin, c.stdout, runError = subprocess.Run(".", c.cdbPath, programArgs...)
	if runError != nil {
		log.Fatal(runError)
	}

	return
}

func (c *Cdb) Read() (output string, runError error) {
	if c.stdout != nil {
		reader := bufio.NewReader(*c.stdout)
		buf := make([]byte, 1024)
		_, err := reader.Read(buf)
		if err != nil {
			runError = err
		}
		output = string(buf)
	}
	return
}

func (c *Cdb) Write(input string) (runError error) {
	if c.stdout != nil {
		writer := bufio.NewWriter(*c.stdin)
		writer.WriteString(input)
		writer.WriteString("\n")
		writer.Flush()
	}
	return
}
