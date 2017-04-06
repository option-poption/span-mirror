package main

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"regexp"
)

var fileRE = regexp.MustCompile("^cme\\.20\\d+\\.s\\.pa2\\.zip$")
var numberRE = regexp.MustCompile("^\\d+$")

func main() {
	client, err := ftp.Dial("ftp.cmegroup.com:21")
	if err != nil {
		panic(err)
	}
	defer client.Quit()

	err = client.Login("anonymous", "anonymous")
	if err != nil {
		panic(err)
	}

	err = client.ChangeDir("/span/data/cme")
	if err != nil {
		panic(err)
	}

	download(client)
}

func download(client *ftp.ServerConn) {

	entries, err := client.List("")

	folders := []string{}
	for _, entry := range entries {
		if entry.Type == ftp.EntryTypeFile {
			if fileRE.MatchString(entry.Name) {
				fmt.Println(entry.Name)
			}
		} else if entry.Type == ftp.EntryTypeFolder {
			if numberRE.MatchString(entry.Name) {
				folders = append(folders, entry.Name)
			}
		}
	}

	// traverse into folders
	for _, folder := range folders {
		err = client.ChangeDir(folder)
		if err != nil {
			panic(err)
		}

		download(client)

		err = client.ChangeDir("..")
		if err != nil {
			panic(err)
		}
	}
}
