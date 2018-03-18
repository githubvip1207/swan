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

	if err := cfg.Reload(*root); err != nil {
		fmt.Println(err.Error())
		return
	}
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
		cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && %s", *root, build))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Println("\nStart build...")
		if err := cmd.Run(); err != nil {
			fmt.Println("============== Build failed ===================")
			return
		}
		fmt.Println("Build was successful")
	}

	if stop != "" {
		cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && %s", *root, stop))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Println("\nStop project...")
		if err := cmd.Run(); err != nil {
			fmt.Println("============== Stop failed ===================")
			return
		}
		fmt.Println("Stop was successful")
	}

	if start != "" {
		cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && %s", *root, start))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Println("\nStart project...")
		if err := cmd.Run(); err != nil {
			fmt.Println("============== Start failed ===================")
			return
		}
		fmt.Println("Start was successful")
	}
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
