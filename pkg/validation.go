// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

// ValidateOptions ...
func ValidateOptions(o *Options) error {

	// validate TTL
	if o.TTL < 0 {
		return errors.New("error: cannot use negative TTL value")
	}

	// validate server / client flags
	if len(o.Address) == 0 {
		ip, err := GetOutboundIP()
		if err != nil {
			fmt.Println("WARNING: cannot obtain outbound IP. using 127.0.0.1")
		}
		o.Address = ip
	}

	// validate port
	if o.Port < 0 || o.Port > 65535 {
		return errors.New("error: port value should be in the [0-65535] range")
	}

	// validate IP
	if ip := net.ParseIP(o.Address); ip == nil {
		return fmt.Errorf("error: invalid server IP: %s", o.Address)
	}

	// make sure output path exists
	if o.OutputPath != "" && o.OutputPath != "./" {
		fmt.Printf("* creating output path: %s\n", o.OutputPath)
		if err := os.MkdirAll(o.OutputPath, os.ModePerm); err != nil {
			return err
		}
	}

	// client only validation
	if !o.Listen {
		// validate the file list
		if len(o.Files) > MaxFiles {
			return errors.New("error: maximum number of files is 128")
		}

		for i := range o.Files {
			o.Files[i] = strings.TrimSpace(o.Files[i])
			f := o.Files[i]
			if len(f) == 0 {
				if len(o.Files) == 1 {
					return errors.New("error: empty list of files")
				}
				return fmt.Errorf("error: empty filename at position: %d", i)
			}
			if len(f) > MaxFileName {
				return fmt.Errorf("error: maximum filename length is %d", MaxFileName)
			}
			// fetching files parent to the input path of the server is a securty risk
			if strings.Contains(f, "..") {
				return fmt.Errorf("error: filename cannot contain '..': %s", f)
			}
		}
	} else {
		if len(o.Token) == 0 {
			token, err := CreateToken()
			if err != nil {
				return fmt.Errorf("error: cannot auto-create new token")
			}
			o.Token = []byte(token)
		}
	}

	// validate token
	if len(o.Token) == 0 {
		return errors.New("error: token cannot be empty")
	}
	tokenRegExp := regexp.MustCompile(TokenPattern)
	if !tokenRegExp.Match(o.Token) {
		return fmt.Errorf("error: token does not match the pattern %q", TokenPattern)
	}

	return nil
}
