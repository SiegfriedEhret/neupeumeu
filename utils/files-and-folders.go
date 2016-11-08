package utils

import (
	"crypto/sha1"
	"io/ioutil"
	"os"

	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"fmt"
	"io"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

const (
	FILEMODE = 0711
)

var (
	CacheDir string
)

func InitDirs() {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Failed to get current user info", err.Error())
	}

	err, CacheDir = CreateNeupeumeuCacheDir(currentUser.HomeDir)
	if err != nil {
		log.Error("Error creating directory", CacheDir, err.Error())
	}
}

func CreateNeupeumeuCacheDir(homeDir string) (error, string) {
	dir := homeDir + "/.neupeumeu"
	err := CreateDir(dir)
	return err, dir
}

func CreateDir(dir string) error {
	err := os.MkdirAll(dir, FILEMODE)
	return err
}

func Extract(name, version, path, pathTo string) {
	log.Debugf("Extracting to %s", path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Failed to open %s", path, err.Error())
		return
	}

	reader := bytes.NewReader(data)

	gzipReader, err := gzip.NewReader(reader)
	defer gzipReader.Close()
	if err != nil {
		panic(err)
	}

	tarReader := tar.NewReader(gzipReader)

	tempFolder := CacheDir + "/tmp/" + name + "-" + version + "-" + getTime()
	tempPackageFolder := tempFolder + "/package"
	err = CreateDir(tempPackageFolder)
	if err != nil {
		log.Errorf("Can't create temp folder %s", tempFolder, err.Error())
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		log.Debugf("Header %+v", header)

		err = CreateDir(pathTo)
		if err != nil {
			log.Errorf("Failed to create dir %s", pathTo, err.Error())
		}

		path := filepath.Join(tempFolder, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				log.Errorf("Failed to create path %s", path, err.Error())
			}
			continue
		}

		var fileDir string
		if header.Typeflag == tar.TypeReg {
			fileDir = filepath.Dir(path)
		} else {
			fileDir = path
		}

		err = CreateDir(fileDir)
		if err != nil {
			log.Debugf("Failed to create directory %s", fileDir)
			log.Error(err.Error())
		}

		f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			log.Debugf("Failed to open/create file %s", path)
			log.Error(err.Error())
		}
		defer f.Close()

		_, err = io.Copy(f, tarReader)
		if err != nil {
			log.Debugf("Failed to copy file", f)
			log.Error(err.Error())
		}
	}

	err = CreateDir(pathTo)
	if err != nil {
		log.Debugf("Can't create output folder %s", pathTo)
		log.Error(err.Error())
	}

	err = os.Rename(tempPackageFolder, pathTo)
	if err != nil {
		log.Debugf("Can't move from %s to %s", tempFolder, pathTo)
		log.Error(err.Error())
	}

	err = os.Remove(tempFolder)
	if err != nil {
		log.Debugf("Failed to cleanup temp folder: %s", tempFolder)
		log.Warn(err.Error())
	}
}

func IsShasumValid(path, pkgSha string) (error, bool) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return err, false
	}

	hash := sha1.New()
	hash.Write(raw)

	computed := strings.ToLower(hex.EncodeToString(hash.Sum(nil)))

	log.WithFields(log.Fields{
		"computed": computed,
		"original": pkgSha,
	}).Debug("Checking shasums...")

	return nil, computed == pkgSha
}

func getTime() string {
	return strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
}
