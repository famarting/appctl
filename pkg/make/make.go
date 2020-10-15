package make

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/famartinrh/appctl/pkg/types/app"
	appctl "github.com/famartinrh/appctl/pkg/types/cmd"
)

//Command encapsulates parameters of ExecuteCmd
type Command struct {
	Cmd []string
	Env []string
}

func execute(command *Command, stdOutFile *os.File, stdErrFile *os.File) error {
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

func BuildProject(makefilePath string, projectDir string, vars []app.InputVar) error {
	// make -f $(pwd)/examples/QuarkusJVMMakefile.mk -C examples/simple-app
	env := []string{}
	for _, v := range vars {
		//TODO verify does this need quotes?
		env = append(env, v.Name+"="+v.Value)
	}
	cmd := &Command{Cmd: []string{"make", "-f", makefilePath, "-C", projectDir}, Env: env}
	return execute(cmd, os.Stdout, os.Stderr)
}
