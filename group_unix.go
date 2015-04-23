package group // import "sigint.ca/group"

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	groupModTime time.Time
	groupNames   map[string]string
)

func Name(gid string) (string, error) {
	stat, err := os.Stat("/etc/group")
	if err != nil {
		return "", fmt.Errorf("group: %s", err)
	}
	if stat.ModTime() != groupModTime || groupNames == nil {
		if err := populateMap(); err != nil {
			return "", fmt.Errorf("group: %s", err)
		}
		groupModTime = stat.ModTime()
	}
	if name, ok := groupNames[gid]; ok {
		return name, nil
	} else {
		return "", fmt.Errorf("group: gid not found: %d", gid)
	}
}

func populateMap() error {
	groupNames = make(map[string]string)
	f, err := os.Open("/etc/group")
	if err != nil {
		return err
	}
	defer f.Close()
	b := bufio.NewReader(f)
	scanner := bufio.NewScanner(b)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == '#' {
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) != 4 {
			continue
		}
		groupNames[fields[2]] = fields[0]
	}
	return scanner.Err()
}
