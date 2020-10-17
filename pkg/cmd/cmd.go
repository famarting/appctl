package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	appctl "github.com/famartinrh/appctl/pkg/types/cmd"
)

type Command struct {
	Cmd []string
	Env []string
}

func Execute(command *Command, stdOutFile *os.File, stdErrFile *os.File) error {
	if appctl.Verbosity > 5 {
		fmt.Println("Executing command " + strings.Join(command.Cmd, " "))
	}
	cmd := exec.Command(command.Cmd[0], command.Cmd[1:]...)
	if command.Env != nil {
		if appctl.Verbosity > 5 {
			fmt.Println("With env " + strings.Join(command.Env, " "))
		}
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, command.Env...)
	}
	if stdOutFile != nil {
		cmd.Stdout = stdOutFile
	}
	if stdErrFile != nil {
		cmd.Stderr = stdErrFile
	}
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
