package files

import (
	"encoding/binary"
	"fileSender/client/internal"
	"fileSender/client/internal/data"
	"fileSender/pkg"
	data2 "fileSender/pkg/data"
	"fmt"
	"net"
	"os"
	"sync"
)

func LocalFilesSendingManager(conf data.ClientConfig) {

	var wg sync.WaitGroup
	getInfoAndSendStructure(conf, "", &wg)
	wg.Wait()
}

func getInfoAndSendStructure(conf data.ClientConfig, path string, wg *sync.WaitGroup) {

	fileDetails, err := internal.RetrieveFilesInfoFromServer(conf, path)

	if err != nil {
		fmt.Printf("error while connecting to server %v \n", err)
		panic(1)
	}

	var localFiles map[uint32]data2.FileDetails
	localFiles = pkg.GetFilesList(conf.LocalDirectoryPath + path)

	for k, v := range localFiles {
		if _, ok := fileDetails[k]; !ok {
			if v.IsDirectory {
				getInfoAndSendStructure(conf, path+"/"+v.Name, wg)
			} else {
				wg.Add(1)
				go sendFileToServer(conf, path+"/"+v.Name, wg)
			}
		} else {
			if v.IsDirectory {
				getInfoAndSendStructure(conf, path+"/"+v.Name, wg)
			}
		}
	}
}

func sendFileToServer(config data.ClientConfig, filePath string, wg *sync.WaitGroup) {

	defer wg.Done()

	localFilePath := config.LocalDirectoryPath + filePath

	stat, err := os.Stat(localFilePath)

	if os.IsNotExist(err) {
		fmt.Printf("file not exists %v \n", err)
	}

	connection, err := net.Dial("tcp", config.ServerHost+":"+config.ServerPortTCP)

	defer connection.Close()

	fileDetail := convertFileDetailToBinary(filePath, uint32(stat.Size()), config.UserName, config.GetUserPasswordHash())

	_, err = connection.Write([]byte(fileDetail))

	if err != nil {
		fmt.Println("Error while sending filename to server")
	}

	canTransferBinary := make([]byte, 2)
	_, err = connection.Read(canTransferBinary)
	canTransfer := binary.LittleEndian.Uint16(canTransferBinary)

	if err != nil {
		fmt.Printf("error while getting permission from server %v \n", err)
	}

	if canTransfer == 0 {
		fmt.Println("cannot send file to server")
	}

	file, err := os.Open(localFilePath)
	if err != nil {
		fmt.Printf("error while opening a file %v \n", err)
	}

	fmt.Println("Synchronizing file... " + filePath)
	sendingFileLoop(file, connection)
	fmt.Println("File " + filePath + " successfully synchronized")

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
			fmt.Printf("error while reading the file %v \n", err)
		}
		_, err = connection.Write(fileBuffer[:rSize])

		if err != nil {
			fmt.Printf("error while sending the file %v \n", err)
		}
	}
}

func convertFileDetailToBinary(fileName string, fileS uint32, userName string, userPassword string) []byte {
	fileDetail := make([]byte, data2.FILE_DETAILS_SIZE)
	fileDetail = []byte(fileName + ":" + userName + ":" + userPassword)
	fileSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(fileSize, fileS)
	fileDetail = append(fileDetail, fileSize...)
	return fileDetail
}
