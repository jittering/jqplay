package jq

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"

	rice "github.com/GeertJohan/go.rice"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var Path, Version string
var riceConf = rice.Config{
	LocateOrder: []rice.LocateMethod{rice.LocateEmbedded, rice.LocateAppended, rice.LocateFS, rice.LocateWorkingDirectory},
}

// Locate jq path
func Init() error {
	// first look for existing bin
	Path = "jq"
	err := setVersion()
	if err == nil {
		log.Infof("Using system jq version %s", Version)
		return nil
	}

	// look for local copy
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = SetPath(pwd)
	if err == nil {
		log.Infof("Using pwd jq version %s", Version)
		return nil
	}

	return loadBundledJq()
}

// Fall back to using bundled JQ version
func loadBundledJq() error {
	binBox := riceConf.MustFindBox("../bin/jq")
	bin, err := binBox.Open(osBin())
	if err != nil {
		return errors.Wrapf(err, "binary not found for %s %s", runtime.GOOS, runtime.GOARCH)
	}
	defer bin.Close()

	// found the proper bin, lets copy it to a temp location
	temp, err := ioutil.TempFile("", "jq")
	if err != nil {
		return errors.Wrap(err, "failed to create temp file")
	}
	_, err = io.Copy(temp, bin)
	if err != nil {
		return errors.Wrap(err, "failed to copy bundled bin")
	}
	err = os.Chmod(temp.Name(), 0755)
	if err != nil {
		return errors.Wrap(err, "failed to to make bundled bin executable")
	}
	Path = temp.Name()
	err = setVersion()
	if err != nil {
		return err
	}
	log.Infof("Using bundled jq version %s", Version)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	go func() {
		<-stop
		log.Info("Caught stop signal; cleaning up bundled jq")
		err = os.Remove(Path)
		if err != nil {
			log.Errorf("failed to remove bundled jq at %s: %s", Path, err.Error())
		}
	}()

	return nil
}

func osDir() string {
	return fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
}

func osBin() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(osDir(), "jq.exe")
	}
	return filepath.Join(osDir(), "jq")
}

func SetPath(binDir string) error {
	jqPath := filepath.Join(binDir, "bin", "jq", osDir())
	Path = filepath.Join(jqPath, "jq")

	_, err := os.Stat(Path)
	if err != nil {
		return err
	}

	os.Setenv("PATH", fmt.Sprintf("%s%c%s", jqPath, os.PathListSeparator, os.Getenv("PATH")))

	err = setVersion()

	return err
}

func setVersion() error {
	// get version from `jq --help`
	// since `jq --version` diffs between versions
	// e.g., 1.3 & 1.4
	var b bytes.Buffer
	cmd := exec.Command(Path, "--help")
	cmd.Stdout = &b
	cmd.Stderr = &b
	cmd.Run()

	out := bytes.TrimSpace(b.Bytes())
	r := regexp.MustCompile(`\[version (.+)\]`)
	if r.Match(out) {
		m := r.FindSubmatch(out)[1]
		Version = string(m)

		return nil
	}

	return fmt.Errorf("can't find jq version: %s", out)
}
