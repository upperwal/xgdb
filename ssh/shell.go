package ssh

import (
	"io"
	"strconv"

	"github.com/upperwal/xgdb/locator"
	"github.com/upperwal/xgdb/parser"
)

const (
	PROMT = "(xgdb)"
)

type Shell struct {
	SSHConn *SSHConnection
	StdErr io.Reader
	Stdout io.Reader
	Stdin io.WriteCloser
}

func NewShell(hostname string) (*Shell, error) {
	sc, err := NewSSHConnection(hostname)
	if err != nil {
		return nil, err
	}

	

	in, err := sc.Session.StdinPipe()
	if err != nil {
		return nil, err
	}
	out, err := sc.Session.StdoutPipe()
	if err != nil {
		return nil, err
	}
	er, err := sc.Session.StderrPipe()
	if err != nil {
		return nil, err
	}

	sc.Shell()

	s := &Shell {
		SSHConn: sc,
		StdErr: er,
		Stdout: out,
		Stdin: in,
	}

	return s, nil
}

type ShellGroup struct {
	PID2Shell map[int]*Shell
}

func NewShellGroup(processInfo []locator.ProcessInfo) (*ShellGroup, error) {
	sg := &ShellGroup {
		PID2Shell: make(map[int]*Shell),
	}

	for _, pi := range processInfo {
		sh, err := NewShell(pi.Hostname)
		if err != nil {
			return nil, err
		}

		sg.PID2Shell[pi.PID] = sh
	}

	return sg, nil
}

func (sg *ShellGroup) InitGDB() error {
	data := "gdb --pid "
	for pid, shell := range sg.PID2Shell {
		
		_, err := shell.Stdin.Write([]byte(data + strconv.Itoa(pid) + "\n"))
		if err != nil {
			return err
		}
	}

	return nil
}

func (sg *ShellGroup) Writer(buf []byte) error {
	for _, shell := range sg.PID2Shell {
		
		_, err := shell.Stdin.Write(buf)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sg *ShellGroup) Reader() ([]*parser.RawData, error) {
	var buf [1000]byte
	combined := make([]*parser.RawData, 0)

	for pid, shell := range sg.PID2Shell {
		n, err := shell.Stdout.Read(buf[:])

		parser.PrintPrompt = 0

		if err != nil {
			return nil, err
		}

		data := &parser.RawData {
			PID: pid,
			Data: string(buf[:n]),
		}
		combined = append(combined, data)
	}

	return combined, nil
}

func (sg *ShellGroup) ErrReader() ([]*parser.RawData, error) {
	var buf [1000]byte
	combined := make([]*parser.RawData, 0)

	for pid, shell := range sg.PID2Shell {
		n, err := shell.StdErr.Read(buf[:])

		parser.PrintPrompt = 0

		if err != nil {
			return nil, err
		}

		data := &parser.RawData {
			PID: pid,
			Data: string(buf[:n]),
		}
		combined = append(combined, data)
	}

	return combined, nil
}
