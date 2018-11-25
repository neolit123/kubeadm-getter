// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// RunClient ...
func RunClient(o *Options) error {

	// connection to the server
	connection, err := net.Dial("tcp", JoinHostPort(o.Address, o.Port))
	if err != nil {
		return err
	}
	defer func() {
		fmt.Println("* closing connection")
		connection.Close()
	}()

	// locals
	token := o.Token
	files := o.Files
	outputPath := o.OutputPath

	// process the handshake
	if err := processClientHandshake(connection, token); err != nil {
		return err
	}

	// send the list of files
	if err := processClientFileList(connection, files, token); err != nil {
		return err
	}

	// process files
	for _, file := range files {
		if err := processClientSingleFile(connection, file, outputPath, token); err != nil {
			return err
		}
	}
	fmt.Printf("* done transfering files\n")

	return nil
}

func processClientHandshake(connection net.Conn, token []byte) error {

	// read encrypted handshake
	fmt.Println("* receiving handshake from server")
	szBuf := make([]byte, Uint32Size)
	if err := ConnRead(connection, szBuf, Uint32Size); err != nil {
		return err
	}
	szHandshake := binary.LittleEndian.Uint32(szBuf)
	handShake := make([]byte, szHandshake)
	if err := ConnRead(connection, handShake, int(szHandshake)); err != nil {
		return err
	}

	// decrypt the handshake
	handShake, err := DecryptBytes(handShake, token)
	if err != nil {
		return err
	}

	// send the descrypted result back
	if err := ConnWrite(connection, handShake, HandShakeSize); err != nil {
		return err
	}
	fmt.Println("* handshake was successful")
	return nil
}

func processClientFileList(connection net.Conn, files []string, token []byte) error {

	// encrypt the list of files
	szBuf := make([]byte, Uint32Size)
	filesJoined := strings.Join(files, FileSeparator)
	fmt.Printf("* sending the following list of files: %s\n", files)
	encryptedFiles, err := EncryptBytes([]byte(filesJoined), token)
	if err != nil {
		return err
	}

	// send the list of files
	lenEncryptedFiles := len(encryptedFiles)
	binary.LittleEndian.PutUint32(szBuf, uint32(lenEncryptedFiles))
	if err = ConnWrite(connection, szBuf, Uint32Size); err != nil {
		return err
	}
	if err = ConnWrite(connection, encryptedFiles, lenEncryptedFiles); err != nil {
		return err
	}
	return nil
}

func processClientSingleFile(connection net.Conn, file, outputPath string, token []byte) error {

	// read the encrypted file
	szBuf := make([]byte, Uint32Size)
	if err := ConnRead(connection, szBuf, Uint32Size); err != nil {
		return err
	}
	szBlock := binary.LittleEndian.Uint32(szBuf)
	fmt.Printf("* receiving encrypted file %q of size %d bytes\n", file, szBlock)
	buf := make([]byte, int(szBlock))
	if err := ConnRead(connection, buf, int(szBlock)); err != nil {
		return err
	}

	// decrypt the file
	buf, err := DecryptBytes(buf, token)
	if err != nil {
		return err
	}

	// join output path
	filePath := filepath.Join(outputPath, file)
	fmt.Printf("* writing file %s\n", filePath)
	path, _ := filepath.Split(filePath)
	// create folder if missing
	os.MkdirAll(path, os.ModePerm)

	// write the file
	if err = WriteFile(filePath, buf); err != nil {
		return err
	}
	return nil
}
