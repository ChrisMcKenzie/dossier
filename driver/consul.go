package driver

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/consul/api"
)

type Response struct {
	*api.KVPair

	Error error
}

type ConsulDriver struct {
	*api.KV
	waitIndex uint64
}

type Watcher struct {
	*api.KV
	index uint64
}

func NewConsulDriver() (ConsulDriver, error) {
	client, err := api.NewClient(api.DefaultConfig())
	return ConsulDriver{client.KV(), 0}, err
}

func (d ConsulDriver) WatchFS(prefix, basePath string) error {
	for {
		log.Println("getting fs")
		keys, meta, err := d.List(prefix, &api.QueryOptions{
			WaitIndex: d.waitIndex,
		})
		if err != nil {
			return err
		}
		d.waitIndex = meta.LastIndex

		for _, file := range keys {
			fmt.Println("file changed:", file.Key)
			makeFile(filepath.Join(basePath, file.Key), file.Value)
		}

	}
	return nil
}

func makeFile(path string, value []byte) error {
	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, value, 0777)
}
