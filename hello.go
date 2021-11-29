package main

import (
	"context"
	"os"
	"time"

	"github.com/containers/image/v5/types"
	"github.com/containers/storage"
	"github.com/sirupsen/logrus"

	lb "github.com/containers/common/libimage"
)

func main() {

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	storeOptions := &storage.StoreOptions{
		RunRoot:         "/run/containers/storage",
		GraphRoot:       "/var/lib/containers/storage",
		GraphDriverName: "overlay",
	}

	// Make sure that the tests do not use the host's registries.conf.
	systemContext := &types.SystemContext{
		SystemRegistriesConfPath:    "/etc/containers/registries.conf",
		SystemRegistriesConfDirPath: "/dev/null",
	}

	runtime, err := lb.RuntimeFromStoreOptions(&lb.RuntimeOptions{SystemContext: systemContext}, storeOptions)
	if err != nil {
		logrus.Errorf("error occured: ", err)
		os.Exit(1)
	}
	sys := runtime.SystemContext()
	logrus.Infof("runtime: ", sys.RegistriesDirPath)
	isEx, _ := runtime.Exists("registry.fedoraproject.org/fedora")
	logrus.Infof("isEx:::   ", isEx)

	copyOptions := lb.CopyOptions{
		SystemContext:    sys,
		DirForceCompress: true,
	}
	logrus.Infof("DirForceCompress: ", copyOptions.DirForceCompress)
	options := lb.ImportOptions{}
	options.CopyOptions = copyOptions
	options.PolicyAllowStorage = true
	options.InsecureSkipTLSVerify = types.OptionalBoolTrue
	// options.ManifestMIMEType = options.ManifestMIMEType

	options.Tag = "tarred"
	var name string

	logrus.Infof("*****start******", time.Now())

	name, err = runtime.Import(context.Background(), "fed.tar", &options)

	logrus.Infof("*****stop******", time.Now())

	if err != nil {
		logrus.Errorf("error oc::  ", err.Error())
		runtime.Shutdown(true)
		os.Exit(1)
	}
	logrus.Infof("name::  ", name)
	runtime.Shutdown(true)
}
