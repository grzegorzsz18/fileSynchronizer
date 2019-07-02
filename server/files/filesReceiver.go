package files

import (
	"encoding/binary"
	"fileSender/pkg/data"
	"fmt"
	"net"
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
	fileDetails := make([]byte, data.FILE_NAME_SIZE+4)

	size, err := connection.Read(fileDetails)

	if err != nil {
		fmt.Print("Error while reading filename")
	}

	bytesFileName := fileDetails[:size-4]
	bytesFileSize := fileDetails[size-4:]

	fileName := strings.Trim(string(bytesFileName), ":")
	fileSize := binary.LittleEndian.Uint32(bytesFileSize)

	fmt.Println(fileName, fileSize)
}
