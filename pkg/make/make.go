package make

import (
	"os"

	"github.com/famartinrh/appctl/pkg/cmd"
	app "github.com/famartinrh/appctl/pkg/types/app/v2"
)

func BuildProject(makefilePath string, projectDir string, vars []app.InputVar) error {
	command := &cmd.Command{Cmd: []string{"make", "-f", makefilePath, "-C", projectDir}, Env: vars}
	return cmd.Execute(command, os.Stdout, os.Stderr)
}
