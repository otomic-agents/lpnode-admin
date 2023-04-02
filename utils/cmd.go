package utils

import (
	"errors"
	"fmt"

	execute "github.com/alexellis/go-execute/pkg/v1"
)

type ExecuteResult struct {
	Stdout string
	Stderr string
}

func ExecuteCMD(cmdStr string, args []string) (ret *ExecuteResult, err error) {
	ret = &ExecuteResult{
		Stdout: "",
		Stderr: "",
	}
	cmd := execute.ExecTask{
		Command:     cmdStr,
		Args:        args,
		StreamStdio: false,
	}
	cmdRes, cmdErr := cmd.Execute()
	if cmdErr != nil {
		err = errors.New(cmdErr.Error())
		return
	}
	if cmdRes.ExitCode != 0 {
		err = fmt.Errorf("non-zero exit code:%s ", cmdRes.Stderr)
		return
	}
	ret = &ExecuteResult{
		Stdout: cmdRes.Stdout,
		Stderr: cmdRes.Stderr,
	}
	return
}
