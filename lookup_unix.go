package group // import "sigint.ca/group"

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	groupModTime time.Time
	groupsByName map[string]*Group
	groupsByID   map[int]*Group
)

func lookup(groupname string) (*Group, error) {
	return lookupUnix(-1, groupname, true)
}

func lookupId(gid string) (*Group, error) {
	i, e := strconv.Atoi(gid)
	if e != nil {
		return nil, e
	}
	return lookupUnix(i, "", false)
}

func lookupUnix(gid int, groupname string, lookupByName bool) (*Group, error) {
	stat, err := os.Stat("/etc/group")
	if err != nil {
		return nil, fmt.Errorf("group: %s", err)
	}
	if stat.ModTime() != groupModTime || groupsByID == nil {
		if err := populateMap(); err != nil {
			return nil, fmt.Errorf("group: %s", err)
		}
		groupModTime = stat.ModTime()
	}
	if lookupByName {
		if g, ok := groupsByName[groupname]; ok {
			return g, nil
		} else {
			return nil, UnknownGroupError(groupname)
		}
	} else {
		if g, ok := groupsByID[gid]; ok {
			return g, nil
		} else {
			return nil, UnknownGroupIdError(gid)
		}
	}
}

func populateMap() error {
	groupsByName = make(map[string]*Group)
	groupsByID = make(map[int]*Group)
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
		id, err := strconv.Atoi(fields[2])
		if err != nil {
			return err
		}
		g := &Group{
			Id:   id,
			Name: fields[0],
		}
		groupsByID[g.Id] = g
		groupsByName[g.Name] = g
	}
	return scanner.Err()
}
