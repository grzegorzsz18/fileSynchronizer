package internal

import (
	"bytes"
	"encoding/json"
	"fileSender/client/internal/data"
	data2 "fileSender/pkg/data"
	"fmt"
	"net/http"
)

func RetrieveFilesInfoFromServer(config data.ClientConfig) (map[uint32]data2.FileDetails, error) {

	jsonData, _ := json.Marshal(data2.UserDetails{
		Name:     config.UserName,
		Password: config.GetUserPasswordHash(),
	})

	resp, err := http.Post("http://"+config.ServerHost+":"+config.ServerPortRest+"/files", "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Printf("request err %v", err)
		return nil, err
	}

	if resp.StatusCode == 401 {
		fmt.Printf("WRONG PASSWORD")
		return nil, nil
	}

	var clientDetails map[uint32]data2.FileDetails

	err = json.NewDecoder(resp.Body).Decode(&clientDetails)

	if err != nil {
		fmt.Printf("decode err %v", err)
		return nil, err
	}

	return clientDetails, nil
}
