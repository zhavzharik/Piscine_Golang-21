package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func displayInfo(dir string, ext string, useF *bool, useExt *bool, useSl *bool, useD *bool) {
	if ext == "." && *useExt {
		fmt.Println("Please specify extension!")
	}
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		link := info.Mode() & fs.ModeSymlink
		currExt := filepath.Ext(path)
		if link == fs.ModeSymlink && *useSl {
			f, err := filepath.EvalSymlinks(path)
			if err != nil {
				fmt.Printf(path + " -> [broken]\n")
			} else {
				fmt.Printf(path + " -> " + f + "\n")
			}
		}
		if strings.Contains(currExt, ".") && *useF && !*useExt {
			fmt.Println(path)
		}
		if currExt == ext && *useF && *useExt {
			fmt.Println(path)
		}
		if info.IsDir() && *useD {
			fmt.Println(path)
		}
		return nil
	})
}

func setFlags() (*bool, *bool, *bool, *bool) {
	useF := flag.Bool("f", false, "print only files")
	useExt := flag.Bool("ext", false, "print only files with a certain extension")
	useSl := flag.Bool("sl", false, "print only symlinks")
	useD := flag.Bool("d", false, "print only directories")
	flag.Parse()
	return useF, useExt, useSl, useD
}

func checkArgs() (int, string, string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	count := flag.NArg()
	path := dir
	ext := "."
	if count == 1 {
		path += flag.Args()[0]
	} else if count == 2 {
		ext += flag.Args()[0]
		path += flag.Args()[1]
	} else {
		fmt.Println("Wrong number of arguments! Required name of folder at least!")
	}
	_, err = ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	return count, ext, path
}

func main() {
	useF, useExt, useSl, useD := setFlags()
	_, ext, path := checkArgs()
	if !*useF && !*useExt && !*useSl && !*useD {
		*useF = true
		*useSl = true
		*useD = true
	}
	displayInfo(path, ext, useF, useExt, useSl, useD)
}
