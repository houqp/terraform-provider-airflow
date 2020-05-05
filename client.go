package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Variables map[string]string

type AirflowClient struct {
	VariablesOutputPath string
	lock                sync.Mutex
}

func (self *AirflowClient) getOutputFilePath() string {
	if self.VariablesOutputPath != "" {
		return self.VariablesOutputPath
	}
	path, err := os.Getwd()
	if err != nil {
		panic("Failed to get current working directory")
	}
	return filepath.Join(path, "airflow_variables.json")
}

func (self *AirflowClient) loadVariables(vars *Variables) {
	output_path := self.getOutputFilePath()
	_, err := os.Stat(output_path)
	if os.IsNotExist(err) {
		return
	}
	f, err := os.Open(output_path)
	if err != nil {
		panic("Failed to open variable file")
	}
	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		panic("Failed to read variable file")
	}
	json.Unmarshal(byteValue, &vars)
}

func (self *AirflowClient) ReadVariable(key string) string {
	vars := Variables{}
	self.lock.Lock()
	defer self.lock.Unlock()

	self.loadVariables(&vars)
	return vars[key]
}

func (self *AirflowClient) DeleteVariable(key string) error {
	vars := Variables{}
	self.lock.Lock()
	defer self.lock.Unlock()

	self.loadVariables(&vars)
	delete(vars, key)

	b, err := json.Marshal(vars)
	if err != nil {
		return err
	}

	output_path := self.getOutputFilePath()
	err = ioutil.WriteFile(output_path, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (self *AirflowClient) SetVariable(key string, value string) error {
	vars := Variables{}
	self.lock.Lock()
	defer self.lock.Unlock()

	self.loadVariables(&vars)
	vars[key] = value

	b, err := json.Marshal(vars)
	if err != nil {
		return err
	}

	output_path := self.getOutputFilePath()
	err = ioutil.WriteFile(output_path, b, 0644)
	if err != nil {
		return err
	}
	return nil
}
