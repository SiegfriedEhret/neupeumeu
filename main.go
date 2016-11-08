package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"gitlab.com/SiegfriedEhret/neupeumeu/pkgdotjson"
	"gitlab.com/SiegfriedEhret/neupeumeu/registry"
	"gitlab.com/SiegfriedEhret/neupeumeu/utils"
)

const (
	APP          = "neupeumeu v%s\n"
	VERSION      = "1.0.0"
	NODE_MODULES = "node_modules"
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
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current working directory", err.Error())
	}

	log.Debug("Running neupeumeu in " + cwd)

	localPkg := pkgdotjson.ReadPackageDotJson("./package.json")
	log.Debug(localPkg)

	utils.InitDirs()

	for depName, depVersion := range localPkg.Dependencies {
		version, prefix := pkgdotjson.GetVersionAndPrefix(depVersion)
		fields := log.Fields{
			"name":    depName,
			"version": version,
			"prefix":  prefix,
		}

		log.WithFields(fields).Debug("Installing module...")

		err, depPkgPath, installedVersion := registry.GetPkgFromRegistry(depName, version, prefix)
		if err != nil {
			log.WithFields(fields).Errorf("Failed to get depencendy package.json", err.Error())
		} else {
			depPkg := pkgdotjson.ReadPackageDotJson(depPkgPath)

			filepath := registry.GetDepFromRegistry(depPkg.Dist.Tarball, depName, installedVersion)
			err, ok := utils.IsShasumValid(filepath, depPkg.Dist.Shasum)
			if err != nil || !ok {
				log.Error("Failed to check package")
			} else {
				utils.Extract(depName, installedVersion, filepath, cwd+"/"+NODE_MODULES+"/"+depName)
			}
		}
	}
}
