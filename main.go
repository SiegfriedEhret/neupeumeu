package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"os/user"
	"io"
	"net/http"
	"strings"
)

type pkg struct {
	Name           string            `json:"name"`
	Version        string            `json:"version"`
	Description    string            `json:"description"`
	Keywords       []string          `json:"keywords"`
	Dependencies   map[string]string `json:"dependencies"`
	DevDpendencies map[string]string `json:"devDependencies"`
}

const (
	APP      = "neupeumeu v%s\n"
	VERSION  = "1.0.0"
	REGISTRY = "https://registry.npmjs.org"
)

var (
	debug bool
	cacheDir string
)

func init() {
	flag.BoolVar(&debug, "d", false, "Run in debug mode")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(APP, VERSION))
		flag.PrintDefaults()
	}

	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	log.Debug("neupeumeu")

	localPkg := readPackageDotJson()
	log.Debug(localPkg)

	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Failed to get current user info", err.Error())
	}

	homeDir := currentUser.HomeDir

	log.Debug(homeDir)
	err, cacheDir = createNeupeumeuCacheDir(homeDir)
	if err != nil {
		log.Error("Error creating directory", cacheDir, err.Error())
	}

	for dep, version := range localPkg.Dependencies {
		log.WithFields(log.Fields{
			"dep": dep,
			"version": version,
		}).Debug("Installing module...")

		var versionToDownload string

		if strings.HasPrefix(version, "^") || strings.HasPrefix(version, "~"){
			versionToDownload = version[1:]
		}

		download(dep, versionToDownload)
	}
}

func createNeupeumeuCacheDir(homeDir string) (error, string) {
	dir := homeDir + "/.neupeumeu"
	err := createDir(dir)
	return err, dir
}

func createDir(dir string) error {
	err := os.MkdirAll(dir, 0711)
	return err
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

func download(n string, v string) {
	// http://registry.npmjs.org/beulogue/-/beulogue-4.0.2.tgz
	from := REGISTRY + "/" + n + "/-/" + n + "-" + v + ".tgz"
	to := cacheDir + "/" + n;

	err := createDir(to)
	if err != nil {
		log.Errorf("Failed to create %s", to, err.Error())
		return
	}

	resp, err := http.Get(from)
	if err != nil {
		log.Errorf("Failed to get", from, err.Error())
	}
	defer resp.Body.Close()

	filepath := to + "/" + v + ".tgz"
	out, err := os.Create(filepath)
	if err != nil {
		log.Errorf("Failed to create %s", filepath, err.Error())
	}
	defer out.Close()

	io.Copy(out, resp.Body)
}