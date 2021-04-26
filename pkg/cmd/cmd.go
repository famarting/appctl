package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	app "github.com/famartinrh/appctl/pkg/types/app/v2"
	appctl "github.com/famartinrh/appctl/pkg/types/cmd"
)

type Command struct {
	Cmd    []string
	Env    []app.InputVar
	Path   string
	logCmd bool
}

func Execute(command *Command, stdOutFile *os.File, stdErrFile *os.File) error {
	varsMap := make(map[string]string)
	for _, v := range command.Env {
		varsMap[v.Name] = v.Value
	}

	needsExpand := func(input string) bool {
		needs := false
		os.Expand(input, func(key string) string {
			needs = true
			if appctl.Verbosity >= 10 {
				fmt.Println("input " + input + " needs expand")
			}
			return ""
		})
		return needs
	}

	expandFunc := func(input string) string {

		if needsExpand(input) {
			return os.Expand(input, func(key string) string {
				return customExpand(varsMap, key)
			})
		} else {
			return input
		}
	}

	firstArg := expandFunc(command.Cmd[0])

	args := []string{}
	for _, a := range command.Cmd[1:] {
		args = append(args, expandFunc(a))
	}

	if appctl.Verbosity > 5 {
		fmt.Println("Executing command " + strings.Join(command.Cmd, " "))
	} else if command.logCmd {
		allexpanded := []string{firstArg}
		allexpanded = append(allexpanded, args...)
		fmt.Println(strings.Join(allexpanded, " "))
	}

	cmd := exec.Command(firstArg, args...)

	if command.Env != nil {
		env := []string{}
		for _, v := range command.Env {
			env = append(env, v.Name+"="+v.Value)
		}

		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, env...)

		if appctl.Verbosity >= 11 {
			fmt.Println("All env " + strings.Join(cmd.Env, " "))
		} else if appctl.Verbosity >= 5 {
			fmt.Println("With env " + strings.Join(env, " "))
		}
	}
	if stdOutFile != nil {
		cmd.Stdout = stdOutFile
	}
	if stdErrFile != nil {
		cmd.Stderr = stdErrFile
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	if currentDir != command.Path {
		defer os.Chdir(currentDir)
		os.Chdir(command.Path)
	}

	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func RunCustomCommand(run string, projectDir string, vars []app.InputVar) error {
	splitted := strings.Split(run, " ")
	command := &Command{Cmd: splitted, Env: vars, Path: projectDir, logCmd: true}
	return Execute(command, os.Stdout, os.Stderr)
}

func customExpand(varsMap map[string]string, key string) string {
	value, ok := varsMap[key]
	if ok {
		if appctl.Verbosity >= 10 {
			fmt.Println("expand for key " + key + " is " + value)
		}
		return value
	} else {
		expand := os.ExpandEnv("$" + key)
		if appctl.Verbosity >= 10 {
			fmt.Println("expand for key " + key + " is " + expand)
		}
		return expand
	}
}
