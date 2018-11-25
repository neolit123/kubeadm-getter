// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"path/filepath"
	"strings"
)

// RunServer ...
func RunServer(o *Options) error {

	// create a TCP server
	server, err := net.Listen("tcp", JoinHostPort(o.Address, o.Port))
	if err != nil {
		return err
	}
	defer server.Close()

	// start a TTL handler
	go TTLHandler(o.TTL)

	// listen for connections
	for {
		connection, err := server.Accept()
		if err != nil {
			return err
		}
		go processServerConnectionWrapper(connection, o.Token, o.InputPath)
	}
}

func processServerConnectionWrapper(connection net.Conn, token []byte, inputPath string) {

	remote := "[" + connection.RemoteAddr().String() + "]:"
	fmt.Printf("* accepted connection from %s\n", remote)
	if err := processServerConnection(connection, remote, token, inputPath); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s closing connection\n", remote)
	connection.Close()
}

func processServerConnection(connection net.Conn, remote string, token []byte, inputPath string) error {

	// handshake
	if err := processServerHandshake(connection, remote, token); err != nil {
		return err
	}

	// file list
	fileList, err := processServerFileList(connection, remote, token)
	if err != nil {
		return err
	}

	// processServer the files
	for _, file := range fileList {
		if err := processServerSingleFIle(connection, remote, inputPath, file, token); err != nil {
			return err
		}
	}

	fmt.Printf("%s done sending files\n", remote)
	return nil
}

func processServerHandshake(connection net.Conn, remote string, token []byte) error {

	// create encrypted handshake
	fmt.Printf("%s sending handshake to client\n", remote)
	handShake, err := CreateHandshakeBytes()
	if err != nil {
		return err
	}
	encHandShake, err := EncryptBytes(handShake, token)
	if err != nil {
		return err
	}

	// send the encrypted handshake
	szBuf := make([]byte, Uint32Size)
	lenHandshake := len(encHandShake)
	binary.LittleEndian.PutUint32(szBuf, uint32(lenHandshake))
	if err = ConnWrite(connection, szBuf, Uint32Size); err != nil {
		return err
	}
	if err = ConnWrite(connection, encHandShake, lenHandshake); err != nil {
		return err
	}

	// wait for the decerypted handshake to be returned
	handShakeFromClient := make([]byte, HandShakeSize)
	if err = ConnRead(connection, handShakeFromClient, HandShakeSize); err != nil {
		return err
	}
	if !bytes.Equal(handShakeFromClient, handShake) {
		return fmt.Errorf("%s handshake did not succeed", remote)
	}
	fmt.Printf("%s handshake was successful\n", remote)
	return nil
}

func processServerFileList(connection net.Conn, remote string, token []byte) ([]string, error) {

	// read the list of requested files
	szFilesBuf := make([]byte, Uint32Size)
	if err := ConnRead(connection, szFilesBuf, Uint32Size); err != nil {
		return nil, err
	}
	szFiles := binary.LittleEndian.Uint32(szFilesBuf)
	fileListBytes := make([]byte, szFiles)
	if err := ConnRead(connection, fileListBytes, int(szFiles)); err != nil {
		return nil, err
	}

	// decrypt the list
	fileListBytes, err := DecryptBytes(fileListBytes, token)
	if err != nil {
		return nil, err
	}
	fileList := strings.Split(string(fileListBytes), FileSeparator)
	fmt.Printf("%s requested the following list of files: %s\n", remote, string(fileListBytes))
	return fileList, nil
}

func processServerSingleFIle(connection net.Conn, remote string, inputPath, file string, token []byte) error {

	// read the requested file
	requestedFile := filepath.Join(inputPath, file)
	raw, err := ReadFile(requestedFile)
	if err != nil {
		fmt.Printf("%s WARNING: requested missing file: %s\n", remote, requestedFile)
		return nil
	}

	// encrypt the file
	encrypted, err := EncryptBytes(raw, token)
	if err != nil {
		return nil
	}
	lenEnc := len(encrypted)

	// send the file
	fmt.Printf("%s sending encrypted file %q of size %d bytes\n", remote, requestedFile, lenEnc)
	szBuf := make([]byte, Uint32Size)
	binary.LittleEndian.PutUint32(szBuf, uint32(lenEnc))
	if err := ConnWrite(connection, szBuf, Uint32Size); err != nil {
		return err
	}
	if err := ConnWrite(connection, encrypted, lenEnc); err != nil {
		return err
	}
	fmt.Printf("%s transfer ended for file: %s\n", remote, requestedFile)
	return nil
}
