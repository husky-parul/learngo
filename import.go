package main

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/containers/image/v5/types"
	"github.com/containers/storage"
	"github.com/containers/storage/pkg/reexec"
	"github.com/sirupsen/logrus"

	lb "github.com/containers/common/libimage"
)

func getRuntime() (runtime *lb.Runtime, cleanup func()) {

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	storeOptions := &storage.StoreOptions{
		RunRoot:         "/run/containers/storage",
		GraphRoot:       "/var/lib/containers/storage",
		GraphDriverName: "overlay",
	}

	systemContext := &types.SystemContext{
		SystemRegistriesConfPath:    "/etc/containers/registries.conf",
		SystemRegistriesConfDirPath: "/dev/null",
	}
	runtime, err := lb.RuntimeFromStoreOptions(&lb.RuntimeOptions{SystemContext: systemContext}, storeOptions)
	if err != nil {
		runtime.LookupImage("8522d622299c", nil)

	}

	cleanup = func() {
		runtime.Shutdown(true)
	}

	return runtime, cleanup
}

func importFromTar() {

	tars, err := ioutil.ReadDir("./tarball")
	if err != nil {
		logrus.Warningf("can't read the tarball dir")
	}
	for _, tar := range tars {
		logrus.Infof("Tar name", tar.Name())
	}
	runtime, cleanup := getRuntime()
	defer cleanup()
	options := lb.ImportOptions{}
	options.Writer = os.Stdout
	logrus.Infof("-----------------------------------------------------")
	name, err := runtime.Import(context.Background(), "./tarball/busybox.tar", &options)
	logrus.Infof("--------------------------STOPP-----------------------------------------")
	if err != nil {
		return
	}
	logrus.Infof("Name:", name)

	logrus.Infof("______+++++++++++++++++++")
	return
}

func main() {
	if reexec.Init() {
		return
	}
	importFromTar()

}
