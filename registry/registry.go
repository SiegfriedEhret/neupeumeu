package registry

import (
	"io"
	"net/http"
	"os"

	"encoding/json"
	"errors"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"gitlab.com/SiegfriedEhret/neupeumeu/pkgdotjson"
	"gitlab.com/SiegfriedEhret/neupeumeu/utils"
)

const (
	REGISTRY = "https://registry.npmjs.org"
)

func GetDepFromRegistry(from string, name string, version string) string {
	log.WithFields(log.Fields{
		"from":    from,
		"name":    name,
		"version": version,
	}).Debug("GetDepFromRegistry")

	to := utils.CacheDir + "/" + name

	err := utils.CreateDir(to)
	if err != nil {
		log.Errorf("Failed to create %s", to, err.Error())
		return ""
	}

	filepath := to + "/" + version + ".tgz"

	download(from, filepath)

	return filepath
}

func GetPkgFromRegistry(name string, version string, prefix string) (err error, toReal string, installedVersion string) {
	from := REGISTRY + "/" + name + "/" + prefix + version
	to := utils.CacheDir + "/" + name

	log.WithFields(log.Fields{
		"name":    name,
		"version": version,
		"prefix":  prefix,
		"from":    from,
		"to":      to,
	}).Debug("GetPkgFromRegistry")

	err = utils.CreateDir(to)
	if err != nil {
		log.Errorf("Failed to create %s", to, err.Error())
		return err, "", ""
	}

	tempJSON := &pkgdotjson.Pkg{}
	err = downloadJSON(from, tempJSON)
	if err != nil {
		log.Errorf("Failed to get %s", from, err.Error())
		return err, "", ""
	} else if tempJSON.Error != "" {
		log.Errorf("Failed to get %s", from, errors.New(tempJSON.Error).Error())
	} else {
		log.Debug("Got JSON", tempJSON)

		installedVersion = tempJSON.Version

		toReal = to + "/" + installedVersion + ".json"

		jsonData, err := json.Marshal(tempJSON)
		if err != nil {
			log.Errorf("Failed to convert JSON data", tempJSON, err.Error())
			return err, "", ""
		} else {
			ioutil.WriteFile(toReal, jsonData, utils.FILEMODE)
		}
	}

	return
}

func downloadJSON(from string, to interface{}) error {
	log.WithFields(log.Fields{
		"from": from,
	}).Debug("downloadJSON")

	resp, err := http.Get(from)
	if err != nil {
		log.Errorf("Failed to get", from, err.Error())
		return err
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(to)
}

func download(from string, to string) {
	log.WithFields(log.Fields{
		"from": from,
		"to":   to,
	}).Debug("download")

	resp, err := http.Get(from)
	if err != nil {
		log.Errorf("Failed to get", from, err.Error())
	}

	defer resp.Body.Close()

	filepath := to
	out, err := os.Create(filepath)
	if err != nil {
		log.Errorf("Failed to create %s", filepath, err.Error())
	}
	defer out.Close()

	io.Copy(out, resp.Body)
}
