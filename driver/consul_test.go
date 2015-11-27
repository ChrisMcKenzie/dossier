package driver

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/consul/api"
)

const File = "file.ext"

var Value = []byte("hello, world")

func createTestKeys(basePath string, driver *ConsulDriver) error {
	_, err := driver.KV().Put(&api.KVPair{
		Key:   filepath.Join(basePath, File),
		Value: Value,
	}, nil)

	return err
}

func cleanTestKeys(basePath string, driver *ConsulDriver) error {
	_, err := driver.KV().Delete(filepath.Join(basePath, File), nil)
	return err
}

func TestConsulDriver(t *testing.T) {
	driver, err := NewConsulDriver()
	if err != nil {
		t.Fatal(err)
	}

	// Setup base directory
	basePath, err := ioutil.TempDir(".", "fs")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(basePath)

	if err := createTestKeys(basePath, driver); err != nil {
		t.Fatal(err)
	}
	defer cleanTestKeys(basePath, driver)

	files, err := driver.BuildFS(basePath, basePath)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file)
		data, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}

		if string(data) != string(Value) {
			t.Error("Remote file and local file do not match")
		}
	}
}
