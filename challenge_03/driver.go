package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
)

type RDriver struct {
	volumes    map[string]string
	m          *sync.Mutex
	mountPoint string
}

func NewRDriver() RDriver {
	return RDriver{
		volumes:    make(map[string]string),
		m:          &sync.Mutex{},
		mountPoint: "/tmp/RDriver",
	}
}

func (d RDriver) Create(r volume.Request) volume.Response {
	logrus.Infof("Create volume: %s", r.Name)
	d.m.Lock()
	defer d.m.Unlock()

	if _, ok := d.volumes[r.Name]; ok {
		return volume.Response{}
	}

	volumePath := filepath.Join(d.mountPoint, r.Name)

	_, err := os.Lstat(volumePath)
	if err != nil {
		logrus.Errorf("Error %s %v", volumePath, err.Error())
		return volume.Response{Err: fmt.Sprintf("Error: %s: %s", volumePath, err.Error())}
	}

	d.volumes[r.Name] = volumePath

	return volume.Response{}
}

func (d RDriver) List(r volume.Request) volume.Response {
	logrus.Info("Volumes list ", r)

	volumes := []*volume.Volume{}

	for name, path := range d.volumes {
		volumes = append(volumes, &volume.Volume{
			Name:       name,
			Mountpoint: path,
		})
	}

	return volume.Response{Volumes: volumes}

}

func (d RDriver) Get(r volume.Request) volume.Response {
	logrus.Info("Get volume ", r)
	if path, ok := d.volumes[r.Name]; ok {
		return volume.Response{
			Volume: &volume.Volume{
				Name:       r.Name,
				Mountpoint: path,
			},
		}
	}
	return volume.Response{
		Err: fmt.Sprintf("volume named %s not found", r.Name),
	}
}

func (d RDriver) Remove(r volume.Request) volume.Response {
	logrus.Info("Remove volume ", r)

	d.m.Lock()
	defer d.m.Unlock()

	if _, ok := d.volumes[r.Name]; ok {
		delete(d.volumes, r.Name)
	}

	return volume.Response{}
}

func (d RDriver) Path(r volume.Request) volume.Response {
	logrus.Info("Get volume path", r)

	if path, ok := d.volumes[r.Name]; ok {
		return volume.Response{
			Mountpoint: path,
		}
	}
	return volume.Response{}
}

func (d RDriver) Mount(r volume.Request) volume.Response {
	logrus.Info("Mount volume ", r)

	if path, ok := d.volumes[r.Name]; ok {
		return volume.Response{
			Mountpoint: path,
		}
	}

	return volume.Response{}

}

func (d RDriver) Unmount(r volume.Request) volume.Response {
	logrus.Info("Unmount ", r)
	return volume.Response{}
}

func (d RDriver) Capabilities(r volume.Request) volume.Response {
	logrus.Info("Capabilities ", r)
	return volume.Response{}
}
