package data

import (
	"fmt"
	"os"
	"time"
)

const NAME = "alpaca"

func createDirAlpaca() {
	err := os.MkdirAll(os.Getenv("HOME")+"/"+NAME, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}
}

func CreateFileAndSave(data string) {
	homeDir := os.Getenv("HOME")

	if homeDir == "" {
		fmt.Println("Cannot get the home dir for user")
		return
	}

	file, err := os.Create(homeDir + "/" + NAME + "/response_" + time.Now().UTC().Format(time.RFC3339) + ".alp")
	if err != nil {
		fmt.Println("Error creating file:", err)
		createDirAlpaca()
		CreateFileAndSave(data)
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
