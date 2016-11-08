package pkgdotjson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type pkgDist struct {
	Shasum  string `json:"shasum"`
	Tarball string `json:"tarball"`
}

type Pkg struct {
	Error          string            `json:"error"`
	Name           string            `json:"name"`
	Version        string            `json:"version"`
	Description    string            `json:"description"`
	Keywords       []string          `json:"keywords"`
	Dependencies   map[string]string `json:"dependencies"`
	DevDpendencies map[string]string `json:"devDependencies"`
	Dist           pkgDist           `json:"dist"`
}

func ReadPackageDotJson(path string) Pkg {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c Pkg
	json.Unmarshal(raw, &c)

	return c
}

func GetVersionAndPrefix(depVersion string) (version string, prefix string) {
	if strings.HasPrefix(depVersion, "^") || strings.HasPrefix(depVersion, "~") {
		version = depVersion[1:]
		prefix = depVersion[:1]
	} else {
		version = depVersion
		prefix = ""
	}

	return
}
