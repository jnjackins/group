package group

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

func Name(gid string) (string, error) {
	groupf, err := os.Open("/etc/group")
	if err != nil {
		return "", errors.New("group file does not exsts")
	}
	b := bufio.NewReader(groupf)
	for {
		line, err := b.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return "", errors.New("error reading group file")
		}
		if line[0] == '#' {
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) != 4 {
			continue
		}
		if fields[2] == gid {
			return fields[0], nil
		}
	}
	return "", errors.New("group not found")
}
