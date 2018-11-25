// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/neolit123/kubeadm-getter/pkg"
)

const example = `
    # server example:
    sudo kubeadm-getter --listen --address=<server-ip> --port=11000 --ttl=240 \
    --token=abcdef.1234567890abcdef --input-path=/etc/kubernetes/pki

    # client example:
    sudo kubeadm-getter --address=<server-ip> --port=11000 --ttl=240 \
    --token=abcdef.1234567890abcdef --output-path=/etc/kubernetes/pki \
    --files=ca.crt;ca.key
`

// Run ...
func Run() error {

	// parse and validate flags
	o := &pkg.Options{}
	if err := parseFlags(o); err != nil {
		fmt.Println("see --help")
		return err
	}

	fmt.Println("* kubeadm-getter")
	fmt.Printf("* using the following token:\n%s\n", string(o.Token))

	// act like a server or client
	if o.Listen {
		fmt.Printf("* server listenting on %s\n", pkg.JoinHostPort(o.Address, o.Port))
		if err := pkg.RunServer(o); err != nil {
			return err
		}
	} else {
		fmt.Printf("* connecting to server %s\n", pkg.JoinHostPort(o.Address, o.Port))
		if err := pkg.RunClient(o); err != nil {
			return err
		}
	}

	return nil
}

func parseFlags(o *pkg.Options) error {
	// add usage
	oldUsage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, example)
		oldUsage()
	}

	var files, token string
	var createToken bool

	// parse flags
	// shared
	flag.IntVar(&o.TTL, "ttl", pkg.DefaultTTL, "maximum time (seconds) for the process to be active. use 0 for no limit")
	flag.IntVar(&o.Port, "port", pkg.DefaultPort, "port to connect to or listen on")
	flag.StringVar(&o.Address, "address", "", "address of a server to connect to as a client or listen to as a server")
	flag.StringVar(&token, "token", "", "token to be used for authorization - e.g. abcdef.1234567890abcdef")
	flag.BoolVar(&createToken, "create-token", false, "creates a secure token")
	// server
	flag.StringVar(&o.InputPath, "input-path", "./", "(server) the path where uploaded files are location")
	flag.BoolVar(&o.Listen, "listen", false, "(server) if set this process will act like a server, otherwise it will act like a client. if the value is empty, the prefered outbound IP will be used.")
	// client
	flag.StringVar(&files, "files", "", "(client) list of files to get from the server folder separated by ';'")
	flag.StringVar(&o.OutputPath, "output-path", "./", "(client) the path where the downloaded files will be written")

	flag.Parse()

	// if create-token is called skip validation
	if createToken {
		token, err := pkg.CreateToken()
		if err != nil {
			fmt.Println("could not create a token")
			os.Exit(1)
		} else {
			fmt.Println(token)
			os.Exit(0)
		}
	}

	filesList := strings.Split(files, pkg.FileSeparator)
	o.Files = filesList
	o.Token = []byte(token)

	return pkg.ValidateOptions(o)
}
