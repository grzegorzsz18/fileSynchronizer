package files

import (
	"encoding/binary"
	"errors"
	"fileSender/client/internal/data"
	data2 "fileSender/pkg/data"
	"fmt"
	"net"
	"os"
)

func SendFileToServer(config data.ClientConfig, fileName string) error {

	stat, err := os.Stat(fileName)

	if os.IsNotExist(err) {
		return errors.New("file not exists")
	}

	connection, err := net.Dial("tcp", config.ServerHost+":"+config.ServerPortTCP)

	defer connection.Close()

	fileDetail := convertFileDetailToBinary(fileName, uint32(stat.Size()))

	_, err = connection.Write([]byte(fileDetail))

	if err != nil {
		fmt.Printf("Error while sending filename to server")
	}

	canTransferBinary := make([]byte, 2)
	_, err = connection.Read(canTransferBinary)
	canTransfer := binary.LittleEndian.Uint16(canTransferBinary)

	if err != nil {
		fmt.Printf("error while getting permission from server %v", err)
	}

	if canTransfer == 0 {
		return errors.New("cannot send file to server")
	}

	sendingFileLoop(fileName, connection)

	return nil
}

func sendingFileLoop(fileName string, connection net.Conn) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("error while opening a file %v", err)
	}
	fileBuffer := make([]byte, data2.FILE_TRANSFERRED_SIZE)
	for rSize := 1; ; {

		rSize, err = file.Read(fileBuffer)

		if rSize == 0 {
			break
		}

		if err != nil {
			fmt.Printf("error while reading the file %v", err)
		}

		_, err := connection.Write(fileBuffer[:rSize])

		if err != nil {
			fmt.Printf("error while sending the file %v", err)
		}
	}
}

func convertFileDetailToBinary(fileName string, fileS uint32) []byte {
	fileDetail := make([]byte, data2.FILE_NAME_SIZE+4)
	fileDetail = []byte(fileName)
	fileSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(fileSize, fileS)
	fileDetail = append(fileDetail, fileSize...)
	return fileDetail
}
