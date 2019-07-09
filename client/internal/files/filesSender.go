package files

import (
	"encoding/binary"
	"fileSender/client/internal/data"
	data2 "fileSender/pkg/data"
	"fmt"
	"net"
	"os"
	"sync"
)

func SendFileToServer(config data.ClientConfig, fileName string, wg *sync.WaitGroup) {

	defer wg.Done()

	path := config.DirectoryPath + "/" + fileName

	stat, err := os.Stat(path)

	if os.IsNotExist(err) {
		fmt.Printf("file not exists")
	}

	connection, err := net.Dial("tcp", config.ServerHost+":"+config.ServerPortTCP)

	defer connection.Close()

	fileDetail := convertFileDetailToBinary(fileName, uint32(stat.Size()), config.UserName)

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
		fmt.Printf("cannot send file to server")
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("error while opening a file %v", err)
	}

	sendingFileLoop(file, connection)

}

func sendingFileLoop(sendingFile *os.File, connection net.Conn) {

	var err error

	fileBuffer := make([]byte, data2.FILE_TRANSFERRED_SIZE)
	for rSize := 1; ; {

		rSize, err = sendingFile.Read(fileBuffer)
		if rSize == 0 {
			break
		}

		if err != nil {
			fmt.Printf("error while reading the file %v", err)
		}
		_, err = connection.Write(fileBuffer[:rSize])

		if err != nil {
			fmt.Printf("error while sending the file %v", err)
		}
	}
}

func convertFileDetailToBinary(fileName string, fileS uint32, userName string) []byte {
	fileDetail := make([]byte, data2.FILE_DETAILS_SIZE)
	fileDetail = []byte(fileName + ":" + userName)
	fileSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(fileSize, fileS)
	fileDetail = append(fileDetail, fileSize...)
	return fileDetail
}
