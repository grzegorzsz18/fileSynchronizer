package internal

import (
	"encoding/json"
	"fileSender/client/internal/data"
	"fmt"
	"io/ioutil"
	"os"
)

func ReadConfig() (data.ClientConfig, error) {

	file, err := os.Open("config.json")

	defer file.Close()

	if err != nil {
		fmt.Printf("Cannot find config.json file")
	}

	conf, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Printf("Error while parsing json file")
	}

	var config data.ClientConfig

	err = json.Unmarshal(conf, &config)

	if err != nil {
		fmt.Printf("Error while unmarshalling data")
	}

	return config, nil
}
