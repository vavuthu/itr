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

package payload

import (
	"os"
	"strings"

	"github.com/vavuthu/itr/logger"
)

var (
	TestCaseName 	string
	commands 		[]string
)

// FormPayload generates the payload content by replacing placeholders.
func FormPayload(basePayload, testCase, testCaseName, configDir string, junitXML bool) string {
	basePayload = strings.ReplaceAll(basePayload, "<MY_TEST_CASE>", testCase)
	basePayload  = strings.ReplaceAll(basePayload, configDir, PodmanPath)
	if junitXML {
		junitFile := PodmanPath + "/" + testCaseName + ".xml"
		basePayload += " --junit-xml " + junitFile
	}
	return basePayload
}

// GenerateAllPodmanCommands generates podman commands for all test cases.
func GenerateAllPodmanCommands(execution, configDir, nonDisruptiveTestCases, image string, junitXML bool) []string {

	// Read the content of the file
	content, err := os.ReadFile(nonDisruptiveTestCases)
	if err != nil {
		logger.Infof("Error in reading file %v", err)
									
		return nil
	}
	
	// Get the test case names
	testCaseNames := strings.Split(strings.TrimSpace(string(content)), "\n")

	// Read the content of the execution file
	executionContent, err := os.ReadFile(execution)
	if err != nil {
		logger.Infof("Error in reading fle %v", err)
	}

	// Replace <MY_TEST_CASE> with actual test case names
	for _, testCase := range testCaseNames {
		testCaseName := LastString(strings.Split(testCase, "::"))
		modifiedContent := FormPayload(string(executionContent), testCase, testCaseName, configDir, junitXML)
		podmanCommand := CommandGenerator(modifiedContent, testCaseName, image, configDir)
		commands = append(commands, podmanCommand)
	}

	return commands
}

func LastString(s []string) string {
	return s[len(s) - 1]
}
