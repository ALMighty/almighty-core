package models

import (
	"io/ioutil"
	"os"
)

// TestDataProvider defines the simple funcion for returning data from a remote provider
type TestDataProvider func() ([]byte, error)

// LoadTestData attempt to load test data from local disk unless;
// * It does not exist or,
// * Variable REFRESH_DATA is present in ENV
//
// Data is stored under examples/test
// This is done to avoid always depending on remote systems, but also with an option
// to refresh/retest against the 'current' remote system data without manual copy/paste
func LoadTestData(filename string, provider TestDataProvider) ([]byte, error) {
	refreshLocalData := func(path string, refresh TestDataProvider) ([]byte, error) {
		content, err := refresh()
		if err != nil {
			return nil, err
		}
		err = ioutil.WriteFile(path, content, 0644)
		if err != nil {
			return nil, err
		}
		return content, nil
	}

	targetDir := "examples/test/"
	err := os.MkdirAll(targetDir, 0777)
	if err != nil {
		return nil, err
	}

	targetPath := targetDir + filename
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		// Call refresher if data does not exist locally
		return refreshLocalData(targetPath, provider)
	}
	if _, found := os.LookupEnv("REFRESH_DATA"); found {
		// Call refresher if force update of test data set in env
		return refreshLocalData(targetPath, provider)
	}

	return ioutil.ReadFile(targetPath)
}
