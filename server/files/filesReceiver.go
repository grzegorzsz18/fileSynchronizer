package files

import (
	"encoding/binary"
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
			fmt.Print("Error while establishing connection")
		}

		connection, err := server.Accept()

		go receiveFileFromClient(connection)

		if err != nil {
			fmt.Print("Error while accepting connection")
		}
	}
}

func receiveFileFromClient(connection net.Conn) {

	defer connection.Close()
	fileDetails := make([]byte, data.FILE_DETAILS_SIZE+4)

	size, err := connection.Read(fileDetails)

	if err != nil {
		fmt.Print("Error while reading filename")
	}

	fileName, fileSize, fileOwner := getFileDetailsFromBinary(fileDetails, uint32(size))

	canTransfer := make([]byte, 2)
	binary.LittleEndian.PutUint16(canTransfer, 1)

	_, err = connection.Write([]byte(canTransfer))

	if err != nil {
		fmt.Printf("error while sending SEND FILE permission to client")
	}

	file, err := os.Create(fileOwner + "/" + fileName)

	if err != nil {
		fmt.Printf("error while creating file %v", err)
	}

	receivingFileLoop(fileSize, connection, file)
}

func receivingFileLoop(fileSize uint32, connection net.Conn, file *os.File) {
	var receivedData uint32 = 0
	for receivedData < fileSize {
		//todo blokować infinite loop np 3x odebrano zero to wylacz
		tmp := make([]byte, data.FILE_TRANSFERRED_SIZE)
		rSize, err := connection.Read(tmp)
		receivedData += uint32(rSize)
		if err != nil {
			fmt.Printf("error while receiving the file %v", err)
		}
		_, err = file.Write(tmp[:rSize])

		if err != nil {
			fmt.Printf("error while writing the file %v", err)
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
