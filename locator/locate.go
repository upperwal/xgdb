package locator

import (
	"bytes"
	"os/exec"
	"strings"
	"strconv"
)

type ProcessInfo struct {
	Hostname string
	PID int
}

func GetProcessInfo(prog_name string) ([]ProcessInfo, error) {
	pi := make([]ProcessInfo, 0)

	pidofCMD := exec.Command("pidof", prog_name)
	var out bytes.Buffer
	pidofCMD.Stdout = &out
	
	if err := pidofCMD.Run(); err != nil {
		return nil, err
	} else {
		o, _ := out.ReadString('\n')

		withoutNewline := strings.Replace(o, "\n", "", -1)
		sa := strings.Split(withoutNewline, " ")

		for _, s := range sa {

			pid_i, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}

			lpi := ProcessInfo {
				Hostname: "localhost",
				PID: pid_i,
			}
			pi = append(pi, lpi)
		}
	}

	return pi, nil
}