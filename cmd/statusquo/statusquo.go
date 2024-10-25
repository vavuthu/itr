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

package statusquo

import (
	"runtime"
	"sync"
	"time"
	
	"github.com/vavuthu/itr/logger"
)

var (
	TestCasesPassed 	int
	TestCasesFailed 	int
	TestCasesNotSelected int
	TotalTestCases 		int
)

// status quo for the test case execution
func Statusquo(wg *sync.WaitGroup) {
	defer func()  {
		wg.Done()
	}()

	for {
		time.Sleep(60 * time.Second)
		// main function and Statusquo are the goroutines
		// exit if there is no test case execution
		if runtime.NumGoroutine() == 2 {
			printStatus()
			break
		}
		printStatus()
	}
}

func printStatus() {
	tcRunning := runtime.NumGoroutine() - 2
	toExecute := TotalTestCases - TestCasesPassed - TestCasesFailed - TestCasesNotSelected - tcRunning
	logger.Info("Total Test cases:", TotalTestCases)
	logger.Info("Passed:", TestCasesPassed)
	logger.Info("Failed:", TestCasesFailed)
	logger.Info("Not selected:", TestCasesNotSelected)
	logger.Info("Test cases running:", tcRunning)
	logger.Info("To Execute:", toExecute)
}
