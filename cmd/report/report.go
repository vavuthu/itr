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
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"golang.org/x/net/html"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/vavuthu/itr/cmd/statusquo"
	"github.com/vavuthu/itr/cmd/utils"
	"github.com/vavuthu/itr/config"
	"github.com/vavuthu/itr/logger"
)

const (
	passed = "passed_testcases.txt"
	failed = "failed_final_testcases.txt"
	skipped = "skipped_testcases.txt"
	testReport = "test_report.html"
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

	passedTests, _ := readLines(passedFilePath)
	skippedTests, _ := readLines(skippedFilePath)
	failedTests, _ := readLines(failedFilePath)

	envMap := make(map[string]string)

	// check test_report.html exists or not
	testReportFilePath := filepath.Join(configDir, testReport)
	istestReportFilePathExist := utils.CheckFileExists(testReportFilePath)
	if istestReportFilePathExist {
		reportContent, err := os.ReadFile(testReportFilePath)
		if err != nil {
			logger.Errorf("Failed to read %s file", testReportFilePath)
		} else {
			doc, err := html.Parse(bytes.NewReader(reportContent))
			if err != nil {
				logger.Errorf("Failed to parse html file %s", testReportFilePath)
			} else {
				// Function to traverse the DOM and extract keys and values from the Environment section
				var extractKeyValues func(*html.Node)
				extractKeyValues = func(n *html.Node) {
					if n.Type == html.ElementNode && n.Data == "table" {
						// Check if this is the "environment" table
						for _, attr := range n.Attr {
							if attr.Key == "id" && attr.Val == "environment" {
								for c := n.FirstChild; c != nil; c = c.NextSibling {
									if c.Type == html.ElementNode && c.Data == "tbody" {
										for tr := c.FirstChild; tr != nil; tr = tr.NextSibling {
											if tr.Type == html.ElementNode && tr.Data == "tr" {
												var key, value string
												for td := tr.FirstChild; td != nil; td = td.NextSibling {
													if td.Type == html.ElementNode && td.Data == "td" {
														if key == "" {
															key = td.FirstChild.Data
														} else {
															value = td.FirstChild.Data
														}
													}
												}
												// Store the key-value pair in the map
												envMap[key] = value
											}
										}
									}
								}
								return // Exit after processing the environment table
							}
						}
					}
					// Traverse the rest of the tree
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						extractKeyValues(c)
					}
				}
				// Start extraction
				extractKeyValues(doc)
			}
		}
	}

	// Generate HTML content
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
	<h2>Environment</h2>
	<table border="1" id="environment">
	`, statusquo.TotalTestCases, totalTime.Minutes(), totalTCPassed, totalTCSkipped, totalTCFailed)

	for key, value := range envMap {
		htmlContent += fmt.Sprintf(`
        <tr>
            <td>%s</td>
            <td>%s</td>
        </tr>
		`, key, value)
	}

	htmlContent += `
    </table>
	<h2>Results</h2>
    <table border="1" id="results-table">
        <tr>
            <th>Test</th>
            <th>Result</th>
        </tr>
	`

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