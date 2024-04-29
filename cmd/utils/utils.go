/*
Copyright © 2023 vavuthu@redhat.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package utils

import (
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/vavuthu/itr/logger"
)

func CountLines(filenames []string) (int, error) {
	totalCount := 0
	for _, filePath := range filenames {
		if len(filePath) == 0 {
			continue
		}
		
		file, err := os.OpenFile(filePath, os.O_RDONLY, 0444)
		if err != nil {
			logger.Errorf("Failed to open file %s with error %s", filePath, err)
		}
		defer file.Close()

		count, _ := CountLinesInFile(file)
		totalCount = totalCount + count
	}
	return totalCount, nil
}

func CountLinesInFile(r io.Reader) (int, error) {
	var count int
	var read int
	var err error
	var target []byte = []byte("\n")

	buffer := make([]byte, 32*1024)
	
	for {
		read, err = r.Read(buffer)
		if err != nil {
			break
		}

		count += bytes.Count(buffer[:read], target)
	}

	if err == io.EOF {
		return count, nil
	}
	
	return count, err
}

func CheckFileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !errors.Is(err, os.ErrNotExist)
}
