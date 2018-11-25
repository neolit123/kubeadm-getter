// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

// JoinHostPort ...
func JoinHostPort(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}

// ConnRead ...
func ConnRead(conn net.Conn, buf []byte, sz int) error {
	var start, n int
	var err error
	defer func() {
		fmt.Println("")
	}()
	for {
		if n, err = conn.Read(buf[start:]); err != nil {
			return err
		}
		start += n
		fmt.Printf("\r                                      ")
		fmt.Printf("\rread %d bytes", start)
		if n == 0 || start == sz {
			break
		}
	}
	return nil
}

// ConnWrite ...
func ConnWrite(conn net.Conn, buf []byte, sz int) error {
	var start, n int
	var err error
	defer func() {
		fmt.Println("")
	}()
	for {
		if n, err = conn.Write(buf[start:]); err != nil {
			return err
		}
		start += n
		fmt.Printf("\r                                      ")
		fmt.Printf("\rwrote %d bytes", start)
		if n == 0 || start == sz {
			break
		}
	}
	return nil
}

// TTLHandler ...
func TTLHandler(ttl int) {

	if ttl == 0 {
		fmt.Println("* WARNING: TTL value is zero! this process will remain open")
		return
	}

	fmt.Printf("* this process will remain open for %v (TTL)\n", time.Duration(ttl)*time.Second)

	var timer uint64
	for {
		if timer > uint64(ttl) {
			fmt.Println("TTL reached!")
			os.Exit(1)
		}
		timer++
		time.Sleep(time.Duration(1) * time.Second)
	}
}

// GetOutboundIP ...
func GetOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1", err // fall back
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}
