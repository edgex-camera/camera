package process

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

func (p *process) handleOutput(stdout io.ReadCloser, stderr io.ReadCloser) error {

	stdoutFilename := p.stdoutFilename
	if stdoutFilename == "" {
		stdoutFilename = fmt.Sprintf("/tmp/%s-%d.stdout", filepath.Base(p.cmd.Path), p.cmd.Process.Pid)
	}

	stderrFilename := p.stderrFilename
	if stderrFilename == "" {
		stderrFilename = fmt.Sprintf("/tmp/%s-%d.stderr", filepath.Base(p.cmd.Path), p.cmd.Process.Pid)
	}

	p.lc.Info(fmt.Sprintf("stdout writing to %s", stdoutFilename))
	p.lc.Info(fmt.Sprintf("stderr writing to %s", stderrFilename))
	go pipeToFile(p.lc, stdout, stdoutFilename)
	go pipeToFile(p.lc, stderr, stderrFilename)
	return nil
}

func pipeToFile(lc logger.LoggingClient, pipe io.ReadCloser, filename string) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		lc.Error(err.Error())
		return
	}
	defer file.Close()

	if _, err := io.Copy(file, pipe); err != nil {
		lc.Error(err.Error())
	}
}
