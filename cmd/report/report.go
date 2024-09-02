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

package report

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/vavuthu/itr/config"
	"github.com/vavuthu/itr/cmd/statusquo"
	"github.com/vavuthu/itr/cmd/utils"
	"github.com/vavuthu/itr/logger"
)

const (
	passed = "passed_testcases.txt"
	failed = "failed_final_testcases.txt"
	skipped = "skipped_testcases.txt"
)

var (
	totalTCPassed int
	totalTCFailed int
	totalTCSkipped int
)

func GenerateSummary(configDir string) {
	logger.Info("########################### SUMMARY ###########################")
	
	passedFilePath := filepath.Join(configDir, passed)
	failedFilePath := filepath.Join(configDir, failed)
	skippedFilePath := filepath.Join(configDir, skipped)

	// Create a map to hold the different statuses
	statuses := map[string]string{
		"Passed":  text.Colors{text.FgGreen}.Sprint("Passed"),
		"Failed":  text.Colors{text.FgRed}.Sprint("Failed"),
		"Skipped": text.Colors{text.FgYellow}.Sprint("Skipped"),
	}

	// Create a new table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Set column names and widths
	t.AppendHeader(table.Row{"Test Case", "Status"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Test Case", WidthMax: 160},
		{Name: "Status", WidthMax: 20},
	})

	isPassedFilePathExist := utils.CheckFileExists(passedFilePath)
	if isPassedFilePathExist {
		totalTCPassed, _ = utils.CountLines([]string{passedFilePath})
		processFile(t, passedFilePath, statuses["Passed"])

	}
	
	isFailedFilePathExist := utils.CheckFileExists(failedFilePath)
	if isFailedFilePathExist {
		totalTCFailed, _ = utils.CountLines([]string{failedFilePath})
		processFile(t, failedFilePath, statuses["Failed"])
	}

	isSkippedFilePathExist := utils.CheckFileExists(skippedFilePath)
	if isSkippedFilePathExist {
		totalTCSkipped, _ = utils.CountLines([]string{skippedFilePath})
		processFile(t, skippedFilePath, statuses["Skipped"])
	}

	logger.Info("Total Test cases: ", statusquo.TotalTestCases)
	logger.Info("Passed: ", totalTCPassed)
	logger.Info("Failed: ", totalTCFailed)
	logger.Info("Skipped: ", totalTCSkipped)
	logger.Info("###############################################################")

	t.Render()
}

// GenerateHTMLReport generates the HTML report
func GenerateHTMLReport(configDir string, totalTime time.Duration) {

	passedFilePath := filepath.Join(configDir, passed)
	failedFilePath := filepath.Join(configDir, failed)
	skippedFilePath := filepath.Join(configDir, skipped)

	// Generate HTML content
	passedTests, _ := readLines(passedFilePath)
	skippedTests, _ := readLines(skippedFilePath)
	failedTests, _ := readLines(failedFilePath)

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Test Summary</title>
</head>
<body>
    <h1>Summary</h1>
    <p>%d tests ran in %.2f minutes</p>
    <p>%d passed, %d skipped, %d failed</p>
    <table border="1">
        <tr>
            <th>Test</th>
            <th>Result</th>
        </tr>
`, statusquo.TotalTestCases, totalTime.Minutes(), totalTCPassed, totalTCSkipped, totalTCFailed)

	for _, test := range passedTests {
		htmlContent += fmt.Sprintf(`
        <tr>
            <td>%s</td>
            <td>Passed</td>
        </tr>
`, test)
	}

	for _, test := range skippedTests {
		htmlContent += fmt.Sprintf(`
        <tr>
            <td>%s</td>
            <td>Skipped</td>
        </tr>
`, test)
	}

	for _, test := range failedTests {
		htmlContent += fmt.Sprintf(`
        <tr>
            <td>%s</td>
            <td>Failed</td>
        </tr>
`, test)
	}

	htmlContent += `
    </table>
</body>
</html>
`
	// Write the HTML content to a file
	htmlReport := "report_" + config.AppConfig.RunID + ".html"
	os.WriteFile(htmlReport, []byte(htmlContent), 0644)
	logger.Infof("HTML file '%s' generated successfully.", htmlReport)

}

// processFile reads a file and adds rows to the table with the corresponding status color
func processFile(t table.Writer, filePath string, status string) {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Errorf("Error opening file %s: %s\n", filePath, err)
		return
	}
	defer file.Close()

	// Read file content line by line and add rows to the table
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		testCase := scanner.Text()
		// Use Colors to set the text color based on status
		statusString := status
		t.AppendRow(table.Row{testCase, statusString})
	}

	// Check for any errors encountered while reading the file
	if err := scanner.Err(); err != nil {
		logger.Errorf("Error reading file %s: %s\n", filePath, err)
	}
}

// unicodeTitle capitalizes the first letter of a string using the Unicode-aware cases package
func unicodeTitle(s string) string {
	title := cases.Title(language.Und)
	return title.String(s)
}


func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}