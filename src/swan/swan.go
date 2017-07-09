/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2017-07-08 19:00:40
# File Name: main.go
# Description:
####################################################################### */

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	cfg "swan/config"
	"swan/utils"
	"time"
)

var (
	root     *string
	build    string
	stop     string
	start      string
	suffixes []string
	codeHash map[string]string
)

func main() {
	pwd, _ := os.Getwd()
	root = flag.String("p", pwd, "project directory")
	/*
		build = flag.String("b", "", "build command")
		stop = flag.String("s", "", "stop command")
		start  = flag.String("c", "", "start command")
	*/
	flag.Parse()

	cfg.Reload(*root)
	if cfg.Handle != nil {
		build = cfg.Handle.GetValue("command", "build")
	}
	if cfg.Handle != nil {
		stop = cfg.Handle.GetValue("command", "stop")
	}
	if cfg.Handle != nil {
		start = cfg.Handle.GetValue("command", "start")
	}

	for _, suffix := range strings.Split(cfg.Handle.GetValue("basic", "suffixes"), ",") {
		suffix = strings.TrimSpace(suffix)
		if suffix != "" {
			suffixes = append(suffixes, suffix)
		}
	}

	fmt.Println("Current project directory:", *root)
	fmt.Println("Config file:", cfg.Path)
	fmt.Println("Build command:", build)
	fmt.Println("Stop command:", stop)
	fmt.Println("Start command:", start)
	fmt.Println("Scan suffixes:", strings.Join(suffixes, ", "))

	codeHash = make(map[string]string)
	files, _ := utils.WalkDir(*root, suffixes)
	for _, file := range files {
		hash := utils.Md5Sum(file)
		codeHash[file] = hash
	}
	scan()
}

func scan() {
	for {
		time.Sleep(2 * time.Second)
		files, err := utils.WalkDir(*root, suffixes)
		if err != nil {
			fmt.Println("Error: %s", err)
			continue
		}

		for _, file := range files {
			newHash := utils.Md5Sum(file)
			hash, ok := codeHash[file]
			if !ok {
				fmt.Println(fmt.Sprintf("\nNew file: %s", file))
				codeHash[file] = newHash
				runCommand()
				break
			}

			if newHash != hash {
				fmt.Println(fmt.Sprintf("\nChange file: %s", file))
				codeHash[file] = newHash
				runCommand()
				break
			}
		}
	}
}

func runCommand() {
	if build != "" {
		out, err := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && %s", *root, build)).Output()
		fmt.Println("\nStart build project, out: ")
		fmt.Println("===========================")
		fmt.Println(string(out))
		if err != nil {
			fmt.Println("Err: ", err)
			return
		}
	}

	if stop != "" {
		out, err := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && %s", *root, stop)).Output()
		fmt.Println("\nStop project, out: ")
		fmt.Println("===========================")
		fmt.Println(string(out))
		if err != nil {
			fmt.Println("Err: ", err)
		}
	}

	if start != "" {
		out, err := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && %s", *root, start)).Output()
		fmt.Println("\nStart project, out: ")
		fmt.Println("===========================")
		fmt.Println(string(out))
		if err != nil {
			fmt.Println("Err: ", err)
		}
	}
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
