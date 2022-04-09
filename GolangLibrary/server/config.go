package server

import (
	"BDSWebsocket/server/logger"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type ServerConfig struct {
	ListenAddr string `json:"ListenAddr"`
	Endpoint   string `json:"Endpoint"`
	Token      string `json:"Token"`
	UsingTLS   bool   `json:"UsingTLS"`
	CertFile   string `json:"CertFile"`
	KeyFile    string `json:"KeyFile"`
}

// ResetDefault reset default config and generate a random token
func (s *ServerConfig) ResetDefault() {
	logger.Warn.Printf("Config file failed to load, using default value")
	s.ListenAddr = ":8080"
	s.Endpoint = "/ws"
	s.UsingTLS = false
	s.CertFile = "cert.pem"
	s.KeyFile = "key.pem"
	//generate random token and print to logger
	s.Token = fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String())))
	logger.Printf("RandomToken: %s", s.Token)

}

func (s *ServerConfig) WriteToFile(file string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, data, os.ModePerm)
}

// LoadConfig load config from file and check if it is valid
func (s *ServerConfig) LoadConfig(file string) {
	var (
		err  error
		data []byte
	)

	//load config from file
	//using json format
	if _, err = os.Stat(file); err != nil {
		logger.Error.Printf("Config file not found: %s", file)
		s.ResetDefault()
		logger.Warn.Printf("Create new config file: %s", file)
		err = s.WriteToFile(file)
		if err != nil {
			logger.Error.Printf("Failed to create config file: %s", err)
		}
		return
	}
	if data, err = ioutil.ReadFile(file); err != nil {
		logger.Error.Printf("Read config file error: %s", err.Error())
		s.ResetDefault()
		return
	}
	err = json.Unmarshal(data, s)
	if err != nil {
		logger.Error.Printf("Parse config file error: %s", err.Error())
		s.ResetDefault()
		return
	}

	//reset default value if config file is not valid
	if !ValidateStructPrint(s) {
		s.ResetDefault()
	}
	if s.UsingTLS {
		if _, err = os.Stat(s.CertFile); err != nil {
			//parse and load cert file

			logger.Error.Printf("Cert file not found: %s", s.CertFile)
			logger.Warn.Printf("TLS will be disabled")
			s.UsingTLS = false
			return
		}
		if _, err = os.Stat(s.KeyFile); err != nil {
			logger.Error.Printf("Key file not found: %s", s.KeyFile)
			logger.Warn.Printf("TLS will be disabled")
			s.UsingTLS = false
			return
		}
	}
	return
}

// Config store the Server Config
var Config ServerConfig
