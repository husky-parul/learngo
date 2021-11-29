package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/containers/image/v5/types"
	"github.com/containers/storage"
	"github.com/sirupsen/logrus"

	lb "github.com/containers/common/libimage"
)

func main() {

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
		fmt.Printf("error occured: ", err)
		os.Exit(1)
	}
	sys := runtime.SystemContext()
	fmt.Printf("runtime: ", sys.RegistriesDirPath)
	isEx, _ := runtime.Exists("registry.fedoraproject.org/fedora")
	fmt.Printf("isEx:::   ", isEx)

	copyOptions := lb.CopyOptions{
		SystemContext: sys,
	}
	options := lb.ImportOptions{}
	options.CopyOptions = copyOptions
	options.Tag = "tarred"
	var name string
	logrus.Infof("Logrus___________-")

	fmt.Printf("*****start******", time.Now())

	name, err = runtime.Import(context.Background(), "/home/parsingh/go/src/github.com/learngo/fed.tar", &options)

	fmt.Printf("*****stop******", time.Now())

	if err != nil {
		fmt.Printf("error oc::  ", err)
		runtime.Shutdown(true)
		os.Exit(1)
	}
	fmt.Printf("name::  ", name)
	runtime.Shutdown(true)
}
