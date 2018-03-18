/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2017-07-09 01:45:34
# File Name: config.go
# Description:
####################################################################### */

package config

import (
	"errors"
	"fmt"
	"os"
	"swan/utils"
)

var (
	Path   string
	Handle *C
)

const (
	name string = ".swanconfig"
)

func Reload(root string) (err error) {
	Path = root
	Path, err = findRealIniFile()
	if err != nil {
		return
	}
	Handle = SetConfig(Path)
	return
}

func findRealIniFile() (r string, err error) {
	f := fmt.Sprintf("%s/%s", Path, name)
	if utils.FileExists(f) {
		return f, nil
	}
	f = fmt.Sprintf("%s/%s", os.Getenv("HOME"), name)
	if utils.FileExists(f) {
		return f, nil
	}
	err = errors.New("config file not found.")
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
