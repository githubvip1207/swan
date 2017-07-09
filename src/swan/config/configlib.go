/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2017-07-08 19:14:04
# File Name: config.go
# Description:
####################################################################### */

package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type C struct {
	path string
	info []map[string]map[string]string
}

func SetConfig(path string) *C {
	c := &C{}
	c.path = path
	return c
}

func (this *C) GetValue(section, name string) string {
	this.ReadList()
	conf := this.ReadList()
	for _, v := range conf {
		for k, v := range v {
			if k == section {
				return v[name]
			}
		}
	}
	return ""
}

func (this *C) SetValue(section, key, value string) bool {
	this.ReadList()
	var ok bool
	var index = make(map[int]bool)
	var conf = make(map[string]map[string]string)
	for i, v := range this.info {
		_, ok = v[section]
		index[i] = ok
	}

	i, ok := func(m map[int]bool) (i int, v bool) {
		for i, v := range m {
			if v == true {
				return i, true
			}
		}
		return 0, false
	}(index)

	if ok {
		this.info[i][section][key] = value
		return true
	} else {
		conf[section] = make(map[string]string)
		conf[section][key] = value
		this.info = append(this.info, conf)
		return true
	}

	return false
}

func (this *C) DeleteValue(section, name string) bool {
	this.ReadList()
	for i, v := range this.info {
		for key, _ := range v {
			if key == section {
				delete(this.info[i][key], name)
				return true
			}
		}
	}
	return false
}

func (this *C) ReadList() []map[string]map[string]string {

	file, err := os.Open(this.path)
	if err != nil {
		CheckErr(err)
	}
	defer file.Close()
	var data map[string]map[string]string
	var section string
	buf := bufio.NewReader(file)
	for {
		l, err := buf.ReadString('\n')
		line := strings.TrimSpace(l)
		if err != nil {
			if err != io.EOF {
				CheckErr(err)
			}
			if len(line) == 0 {
				break
			}
		}
		switch {
		case len(line) == 0:
		case string(line[0]) == "#": //增加配置文件备注
		case line[0] == '[' && line[len(line)-1] == ']':
			section = strings.TrimSpace(line[1 : len(line)-1])
			data = make(map[string]map[string]string)
			data[section] = make(map[string]string)
		default:
			i := strings.IndexAny(line, "=")
			value := strings.TrimSpace(line[i+1 : len(line)])
			data[section][strings.TrimSpace(line[0:i])] = value
			if this.uniquappend(section) == true {
				this.info = append(this.info, data)
			}
		}

	}

	return this.info
}

func CheckErr(err error) string {
	if err != nil {
		return fmt.Sprintf("Error is :'%s'", err.Error())
	}
	return "Notfound this error"
}

//Ban repeated appended to the slice method
func (this *C) uniquappend(conf string) bool {
	for _, v := range this.info {
		for k, _ := range v {
			if k == conf {
				return false
			}
		}
	}
	return true
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
