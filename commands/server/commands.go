package server

import (
	"fmt"
	"bytes"
	"os"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/util/display"
)

func init() {
	display.ServerResponseFunc = AddResponse
}

type CmdFunc func(ccmd *cobra.Command, args []string)

type AdminCmds struct {
	cmds map[string]CmdFunc
	response *Response
}

type Request struct {
	Name string
	Args []string
}

var commands = &AdminCmds{
	cmds: map[string]CmdFunc{},
	response: nil,
}

func (comm *AdminCmds) Run(req Request, resp *Response) error {
	// make display use the buffer
	buff := &bytes.Buffer{}
	display.Out = buff
	defer func() {
		display.Out = os.Stdout
	}()

	cmd, ok := comm.cmds[req.Name]
	if !ok {
		return fmt.Errorf("invalid command %s", req.Name)
	}

	// run the indicated command, dont give it a cobra cmd because we dont have one
	cmd(nil, req.Args)

	resp.Output = fmt.Sprintf("%s%s", buff.String(), comm.response.Output)
	resp.ExitCode = comm.response.ExitCode
	return nil
}

// add a adminisration command that can be run
func AddCmd(key string, cmdFunc CmdFunc) {
	commands.cmds[key] = cmdFunc
}

// func when the command is done the response is added here
func AddResponse(out string, exitCode int) {
	commands.response = &Response{out, exitCode}
}

func RunCommand(cmd string, args []string) (*Response, error) {
	req := Request{
		Name: cmd,
		Args: args,
	}

	resp := &Response{}

	return resp, run("AdminCmds.Run", req, resp)
}