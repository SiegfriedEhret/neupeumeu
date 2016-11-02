package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
)

type pkg struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Keywords     []string          `json:"keywords"`
	Dependencies map[string]string `json:"dependencies"`
}

const (
	APP     = "neupeumeu v%s\n"
	VERSION = "1.0.0"
)

var (
	debug bool
)

func init() {
	flag.BoolVar(&debug, "d", false, "Run in debug mode")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(APP, VERSION))
		flag.PrintDefaults()
	}

	flag.Parse()

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func main() {
	logrus.Debug("neupeumeu")
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
