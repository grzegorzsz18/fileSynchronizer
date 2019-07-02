package files

import (
	"encoding/binary"
	"fileSender/client/internal/data"
	data2 "fileSender/pkg/data"
	"fmt"
	"net"
)

func SendFileToServer(config data.ClientConfig, fileName string) error {

	connection, err := net.Dial("tcp", config.ServerHost + ":" + config.ServerPortTCP)

	fileDetail := make([]byte, data2.FILE_NAME_SIZE + 4)
	fileDetail = []byte(fileName)

	fileSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(fileSize, 31415926)
	fileDetail = append(fileDetail, fileSize...)

	_, err = connection.Write([]byte(fileDetail))

	if err != nil {
		fmt.Printf("Error while sending filename to server")
	}

	return nil
}
