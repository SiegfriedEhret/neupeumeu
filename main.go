package main

import (
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
)

type pkg struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Keywords    []string `json:"keywords"`
}

func init() {
	logrus.Debug("Init bargl")
}

func main() {
	logrus.Debug("bargl")
	logrus.Info(readPackageDotJson())
}

func readPackageDotJson() pkg {
	raw, err := ioutil.ReadFile("./package.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c pkg
	json.Unmarshal(raw, &c)
	return c
}