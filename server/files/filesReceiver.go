package files

import (
	"encoding/binary"
	"errors"
	"fileSender/pkg/data"
	"fmt"
	"net"
	"os"
	"strings"
)

func HandleFilesReceining() {

	server, err := net.Listen("tcp", "localhost:22222")
	defer server.Close()

	for {

		if err != nil {
			fmt.Println("Error while establishing connection")
		}

		connection, err := server.Accept()

		go receiveFileFromClient(connection)

		if err != nil {
			fmt.Println("Error while accepting connection")
		}
	}
}

func receiveFileFromClient(connection net.Conn) {

	defer connection.Close()
	fileDetails := make([]byte, data.FILE_DETAILS_SIZE+4)

	size, err := connection.Read(fileDetails)

	if err != nil {
		fmt.Println("Error while reading filename")
	}

	fileName, fileSize, fileOwner := getFileDetailsFromBinary(fileDetails, uint32(size))
	path := fileOwner + "/" + fileName
	canTransfer := make([]byte, 2)
	binary.LittleEndian.PutUint16(canTransfer, 1)

	_, err = connection.Write([]byte(canTransfer))

	if err != nil {
		fmt.Println("error while sending SEND FILE permission to client")
	}

	err = createDirStructureIfNotExists(path)

	if err != nil {
		fmt.Printf("Error creating structures %v \n", err)
	}

	file, err := os.Create(path)

	if err != nil {
		fmt.Printf("error while creating file %v \n", err)
	}

	receivingFileLoop(fileSize, connection, file)
}

func receivingFileLoop(fileSize uint32, connection net.Conn, file *os.File) {
	var receivedData uint32 = 0
	for receivedData < fileSize {
		tmp := make([]byte, data.FILE_TRANSFERRED_SIZE)
		rSize, err := connection.Read(tmp)
		receivedData += uint32(rSize)
		if err != nil {
			fmt.Printf("error while receiving the file %v \n", err)
			break
		}
		_, err = file.Write(tmp[:rSize])

		if err != nil {
			fmt.Printf("error while writing the file %v \n", err)
			break
		}
	}
}

func getFileDetailsFromBinary(fileDetails []byte, size uint32) (string, uint32, string) {
	bytesFileDetails := fileDetails[:size-4]
	bytesFileSize := fileDetails[size-4:]
	fileSize := binary.LittleEndian.Uint32(bytesFileSize)
	fileDetailsSplitted := strings.Split(string(bytesFileDetails), ":")
	fileName := fileDetailsSplitted[0]
	fileOwner := fileDetailsSplitted[1]
	return fileName, fileSize, fileOwner
}

func createDirStructureIfNotExists(filePath string) error {

	splitted := strings.Split(filePath, "/")

	switch len(splitted) {
	case 0:
		return errors.New("Wrong file path format")
	case 1:
		return nil //just file without dir
	default:
		return os.MkdirAll(strings.Join(splitted[:len(splitted)-1], "/"), 0777)

	}
}
