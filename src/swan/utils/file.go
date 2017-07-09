/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2017-07-09 02:08:30
# File Name: src/swan/utils/file.go
# Description:
####################################################################### */

package utils

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func FileExists(file string) bool {
	finfo, err := os.Stat(file)
	if err == nil && !finfo.IsDir() {
		return true
	}
	return false
}

func WalkDir(path string, suffixes []string) (files []string, err error) {
	for k, suffix := range suffixes {
		suffixes[k] = strings.ToUpper(suffix)
	}

	err = filepath.Walk(path, func(filename string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}

		if len(suffixes) == 0 {
			files = append(files, filename)
		}
		for _, suffix := range suffixes {
			if strings.HasSuffix(strings.ToUpper(fi.Name()), fmt.Sprintf(".%s", suffix)) {
				files = append(files, filename)
			}
		}
		return nil
	})

	return files, err
}

func Md5Sum(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()

	r := bufio.NewReader(f)
	md5hash := md5.New()

	_, err = io.Copy(md5hash, r)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", md5hash.Sum(nil))
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
