package main

import (
	"errors"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Cluster structure to setup kubeconfig
type Cluster struct {
	Name        string `yaml:"name"`
	APIServer   string `yaml:"apiserver"`
	IssuerURL   string `yaml:"issuer"`
	ClientID    string `yaml:"clientid"`
	KubeConfig  string `yaml:"kubeconfig"`
	KeepContext bool   `yaml:"keepcontext"`
	Port        int    `yaml:"port"`
	NameSpace   string `yaml:"namespace"`
	ManualInput bool   `yaml:"manualinput"`
}

func readConfig(name string) (*Cluster, error) {
	path := filepath.Join(home, kubedConf)
	confBytes, err := os.ReadFile(path)
	if err != nil {
		log.Warn("Failed in reading kubed config file ", err)
		return nil, err
	}

	var clusters []Cluster
	err = yaml.Unmarshal(confBytes, &clusters)
	if err != nil {
		log.Error("Failed in parsing config file ", err)
	}

	for _, c := range clusters {
		if c.Name == name {
			return &c, nil
		}
	}

	return nil, errors.New("provided cluster not found, run with full config parameters to configure it")
}

func setConfig(
	name string,
	apiserver string,
	issuerURL string,
	clientID string,
	kubeconfig string,
	keepContext bool,
	port int,
	namespace string,
	manualInput bool) *Cluster {

	return &Cluster{
		Name:        name,
		APIServer:   apiserver,
		IssuerURL:   issuerURL,
		ClientID:    clientID,
		KubeConfig:  kubeconfig,
		KeepContext: keepContext,
		Port:        port,
		NameSpace:   namespace,
		ManualInput: manualInput,
	}
}

func saveConfig(cluster *Cluster) error {
	path := filepath.Join(home, kubedConf)

	var clusters []Cluster

	oldConfBytes, err := os.ReadFile(path)
	if err == nil {
		err = yaml.Unmarshal(oldConfBytes, &clusters)
		if err != nil {
			log.Error("Failed in parsing config file ", err)
			clusters = nil
		}
	}

	found := false
	if clusters != nil {
		for i, c := range clusters {
			// Insert the recent config
			if c.Name == cluster.Name {
				clusters[i] = *cluster
				found = true
			}
		}
		if !found {
			clusters = append(clusters, *cluster)
		}
	} else {
		clusters = append(clusters, *cluster)
	}

	newConfBytes, err := yaml.Marshal(clusters)
	if err != nil {
		log.Warn("Failed in marshaling kubedconfig ", err)
		return err
	}

	err = os.WriteFile(path, newConfBytes, 0644)
	if err != nil {
		log.Warn("Failed in saving kubedconfig ", err)
		return err
	}

	return nil
}
