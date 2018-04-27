package gdb

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/upperwal/xgdb/locator"
	"github.com/upperwal/xgdb/parser"
	"github.com/upperwal/xgdb/ssh"
)

type GDB struct {
	ProcessInfo []locator.ProcessInfo
	ShellGroup *ssh.ShellGroup
}

func NewGDB(progname string) (*GDB, error) {
	fmt.Println("Starting xgdb...\n")

	pi, err := locator.GetProcessInfo(progname)
	if err != nil {
		if err.Error() == "exit status 1" {
			return nil, errors.New("Could not find a running process")
		}
		return nil, err
	}

	sg, err := ssh.NewShellGroup(pi)
	if err != nil {
		return nil, err
	}

	return &GDB {
		ProcessInfo: pi,
		ShellGroup: sg,
	}, nil
}

func (gdb *GDB) Start() {
	go gdb.stdinRoutine()
	go gdb.remoteStdoutRoutine()
	go gdb.remoteErrRoutine()

	parser.PROMPT = parser.SHELL

	gdb.InitGDB()
}

func (gdb *GDB) InitGDB() {
	// Attach gdb to pid's
	parser.PROMPT = parser.XGDB

	err := gdb.ShellGroup.InitGDB()
	if err != nil {
		panic(err)
	}
}

func (gdb *GDB) stdinRoutine() {
	var buf [1000]byte
	for {
		n, _ := os.Stdin.Read(buf[:])

		if buf[0] == 10 {
			fmt.Printf("%s", parser.PROMPT)
			continue
		}

		parser.PrintPrompt = 1

		cont := parser.PreProcess(buf[:n])

		if cont == true {
			continue
		}

		err := gdb.ShellGroup.Writer(buf[:n])
		if err != nil {
			panic(err)
		}

		time.Sleep(10 * time.Millisecond)

		if parser.PrintPrompt == 1 {
			fmt.Printf("%s", parser.PROMPT)
		}
		
	}
}

func (gdb *GDB) remoteStdoutRoutine() {

	for {
		s, err := gdb.ShellGroup.Reader()
		if err != nil {
			panic(err)
		}
		
		parser.ParseStdoutStream(s)
	}

}

func (gdb *GDB) remoteErrRoutine() {

	for {
		s, err := gdb.ShellGroup.ErrReader()
		if err != nil {
			panic(err)
		}
		
		parser.ParseErrStream(s)
	}

}
