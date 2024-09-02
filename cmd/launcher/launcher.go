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

package launcher

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vavuthu/itr/cmd/mail"
	"github.com/vavuthu/itr/cmd/payload"
	"github.com/vavuthu/itr/cmd/report"
	"github.com/vavuthu/itr/cmd/statusquo"
	"github.com/vavuthu/itr/config"
	"github.com/vavuthu/itr/logger"
)

var failedTC = make(map[string]int)
var logsDir string
var exitCode int

type execute interface {
	Execute() error
}

type executeRetry interface {
	execute
	decreaseRetry()
	retriesLeft() int
}

type Command struct {
	cmd     string
	retries int
}

func (c *Command) Execute() error  {
	s := payload.LastString(strings.Split(c.cmd, "::"))
	testCaseName := strings.Split(s, " ")[0]
	logFile := filepath.Join(logsDir, testCaseName)

	// Open a file for writing (create it if not exists, truncate if exists)
	outputFile, err := os.Create(logFile)
	
	if err != nil {
		logger.Infof("Error creating output file: %v", err)
	}
	defer outputFile.Close()

	cmdSplit := strings.Fields(c.cmd)
	var testCase string
	for _, each := range cmdSplit {
		if isContains := strings.Contains(each, testCaseName); isContains {
			testCase = each
			break
		}
	}
	logger.Infof("Running test case: %s and live log streamed at %s", testCase, outputFile.Name())
	
	parts := strings.Fields(c.cmd)
	podmanCmd := exec.Command(parts[0], parts[1:]...)
	podmanCmd.Stderr = podmanCmd.Stdout
	stdoutPipe, err := podmanCmd.StdoutPipe()

	if err != nil {
		logger.Infof("Error in creating stdout pipe: %v", err)
		return fmt.Errorf("test case: %s failed", testCase)
	}

	if err := podmanCmd.Start(); err != nil {
		logger.Infof("Error in starting command: %v", err)
		return fmt.Errorf("test case: %s failed", testCase)
	}
	
	_, err = io.Copy(outputFile, stdoutPipe)
	if err != nil {
		logger.Infof("Error in copying output to file: %v", err)
		return fmt.Errorf("test case: %s failed", testCase)
	}

	err = podmanCmd.Wait()
	if err != nil {
		logger.Errorf("Error in waiting for command %v and test case is %s", err, testCase)
		if c.retriesLeft() > 0 {
			logFileName := outputFile.Name()
			backupFile := generateEpochFileName(logFileName)
			logger.Infof("moving log file %s to %s", logFileName, backupFile)
			os.Rename(logFileName, backupFile)
		}
		return fmt.Errorf("test case: %s failed", testCase)
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		// The command did not exit successfully
		logger.Infof("Command exited with non-zero status code: %v", exitErr)
		return fmt.Errorf("test case: %s failed", testCase)
	}

	logger.Infof("test case: %s executed successfully.", testCase)
	statusquo.TestCasesPassed++

	// if it passes, check test case present in failedTC and remove it if it present
	delete(failedTC, c.cmd)

	return nil
}

func (c *Command) retriesLeft() int {
	return c.retries
}

func (c *Command) decreaseRetry() {
	c.retries--
}

func (c *Command) String() string {
	return c.cmd
}

type Launcher struct {
	payload                []execute
	failedTCAfterRetry     *os.File
	payloadLock            sync.Mutex
	wg                     sync.WaitGroup
}

func (l *Launcher) LaunchCommands(queueLength, retry int) {
	workerPool := make(chan struct{}, queueLength)

	// Infinite loop to run commands till payload is completed
	out:
	for {
		time.Sleep(time.Second)
		if len(l.payload) == 0 {
			if len(failedTC) > 0  {
				// check if all test cases are retried with given Retries
				for k, v := range failedTC {
					if v != retry + 1 {
						logger.Infof("tc %s doesn't run with max retries %d", k, retry)
						time.Sleep(time.Second)
						continue out

					}
				}
				// shouldnot break for loop here because, lets say we have 3 TC's, and tc1 failed all times and
				// still tc2 and tc3 running. In this case, we shouldnot break for loop, because if we break
				// there might be chance tc2/tc3 will fail and it won't retry for retry times.
			}
			if len(failedTC) == 0 && runtime.NumGoroutine() == 2 {
				logger.Info("All the test cases are executed")
				break
			}

			// no payload and no failedTC's or no payload and all faileed TC retries time equal to max retry times
			if runtime.NumGoroutine() > 2 {
				logger.Info("still some test cases are running .....")
				time.Sleep(time.Second)
				continue out
	
			}
			
			// if we reach here which means, len(l.payload) == 0 and len(failedTC) > 0 and all TC's retries times equal to retry
			break
		}

		l.payloadLock.Lock()
		cmd := l.payload[0]
		l.payload = l.payload[1:]
		l.payloadLock.Unlock()

		workerPool <- struct{}{}
		l.wg.Add(1)
		go l.LaunchExecute(workerPool, cmd)
	}

	l.wg.Wait()
}

func (l *Launcher) LaunchExecute(workerPool chan struct{}, e execute) {
	defer func()  {
		<-workerPool
		l.wg.Done()
	}()

	if err := e.Execute(); err != nil {
		logger.Warnf("%v", err)
		if cmd, ok := e.(executeRetry); ok {
			logger.Infof("Retries left for %s is %d", strings.Fields(err.Error())[2], cmd.retriesLeft())
			if cmd.retriesLeft() > 0 {
				cmd.decreaseRetry()
				l.payloadLock.Lock()
				l.payload = append(l.payload, cmd)
				l.payloadLock.Unlock()					
			} else {
				tc := strings.Fields(err.Error())[2]
				tcWithNewLine := tc + "\n"
				logger.Warnf("test case: %v exceeded maximum retries", tc)
				statusquo.TestCasesFailed++
				logger.Error("test case:", strings.Fields(err.Error())[2], "failed")
				l.failedTCAfterRetry.Write([]byte(tcWithNewLine))
				exitCode = 3
			}
		}
	}
}

func LaunchInitiate(commands []string, configDir string, queueLength, retry int) {
	logger.Info("Intiating Launch with queueLength: ", queueLength)
	
	logsDir = configDir
	executor := []execute{}
	for _, cmd := range commands {
		executor = append(executor, &Command{cmd: cmd, retries: retry})
	}

	// Initialize the Launcher
	launch := &Launcher{
		payload: make([]execute, 0),
	}

	// Add commands to paylod
	launch.payloadLock.Lock()
	launch.payload = append(launch.payload, executor...)
	launch.payloadLock.Unlock()
	
	// statusquo
	var wg1 sync.WaitGroup
	wg1.Add(1)
	go statusquo.Statusquo(&wg1)
	
	failedTCAfterRetryFile := filepath.Join(logsDir, "failed_final_testcases.txt")
	launch.failedTCAfterRetry, _ = os.Create(failedTCAfterRetryFile)
	defer launch.failedTCAfterRetry.Close()

	startTime := time.Now()
	// Execute commands
	launch.LaunchCommands(queueLength, retry)
	wg1.Wait()

	endTime := time.Now()
	totalTime := endTime.Sub(startTime)

	// report generation
	report.GenerateSummary(configDir)
	report.GenerateHTMLReport(configDir, totalTime)

	if config.AppConfig.EmailID != "" {
		mail.SendMail()
		logger.Info("Email sent successfully to ", config.AppConfig.EmailID)
	}

	os.Exit(exitCode)
}

// suffix epoch time for the file name
func generateEpochFileName(filename string) string {
	now := time.Now()
	epochTimeStr := strconv.FormatInt(now.Unix(), 10)
	backupFile := filename + "-" + epochTimeStr
	return backupFile
}

func createFile(filename string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return file, err
}
