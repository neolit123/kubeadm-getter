// SPDX-License-Identifier: Apache-2.0

package pkg

import "io/ioutil"

// ReadFile ...
func ReadFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// WriteFile ...
func WriteFile(filePath string, data []byte) error {
	if err := ioutil.WriteFile(filePath, data, 0660); err != nil {
		return err
	}
	return nil
}
