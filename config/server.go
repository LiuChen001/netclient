// Package config provides functions for reading the config.
package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gravitl/netmaker/models"
	"github.com/kr/pretty"
	"gopkg.in/yaml.v3"
)

var Servers map[string]Server
var ServerNodes map[string]struct{}

type Server struct {
	Name        string
	Version     string
	API         string
	CoreDNSAddr string
	Broker      string
	MQPort      string
	MQID        string
	Password    string
	DNSMode     bool
	Is_EE       bool
	Nodes       []string
}

type Networks map[string]struct{}

// ReadServerConfig reads a server configuration file and returns it as a
// Server instance. If no configuration file is found, nil and no error will be
// returned. The configuration must live in one of the directories specified in
// with AddConfigPath()
//
// In case multiple configuration files are found, the one in the most specific
// or "closest" directory will be preferred.
func GetServers() error {
	file := GetNetclientPath() + "servers.yml"
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(&Servers); err != nil {
		return err
	}
	return nil
}

func WriteServerConfig() error {
	file := GetNetclientPath() + "servers.yml"
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(GetNetclientPath(), os.ModePerm)
		} else if err != nil {
			return err
		}
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	log.Println("servers to be saved")
	pretty.Println(Servers)
	err = yaml.NewEncoder(f).Encode(Servers)
	if err != nil {
		return err
	}
	return f.Sync()
}

func WriteInitialServerConfig(cfg *models.ServerConfig) error {
	var s Server
	s.Name = strings.Replace(cfg.Server, "broker.", "", 1)
	s.Broker = cfg.Server
	s.MQPort = cfg.MQPort
	s.API = cfg.API
	s.DNSMode, _ = strconv.ParseBool(cfg.DNSMode)
	s.CoreDNSAddr = cfg.CoreDNSAddr
	s.Version = cfg.Version
	s.Is_EE = cfg.Is_EE
	s.MQID = Netclient.HostID
	s.Password = Netclient.HostPass
	Servers[s.Name] = s
	return WriteServerConfig()
}

func ConvertServerCfg(cfg *models.ServerConfig) Server {
	var server Server
	server.Name = cfg.Server
	server.Version = cfg.Version
	server.Broker = cfg.Broker
	server.MQPort = cfg.MQPort
	server.MQID = Netclient.HostID
	server.Password = Netclient.HostPass
	server.API = cfg.API
	server.CoreDNSAddr = cfg.CoreDNSAddr
	server.Is_EE = cfg.Is_EE
	server.DNSMode, _ = strconv.ParseBool(cfg.DNSMode)
	return server
}
