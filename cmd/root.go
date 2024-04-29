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
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/vavuthu/itr/cmd/engine"
	"github.com/vavuthu/itr/cmd/validate"
	"github.com/vavuthu/itr/logger"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "itr",
	Short: "Intelligent Test Runner (ITR) is tool that runs the test cases in parallel with user controlled queues",
	Long:  `Intelligent Test Runner (ITR) is tool that runs the test cases in parallel with user controlled queues.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: runCmd,
}

var (
	configDir			string
	executionFile 		string
	image 				string
	junitXML 			bool
	negativeTestCases 	string
	positiveTestCases 	string
	queueLength 		int
	retry 				int
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&positiveTestCases, "positive_testcases", "p", "", "Path to positive test cases to run")
	rootCmd.Flags().StringVarP(&negativeTestCases, "negative_testcases", "n", "", "Path to negative test cases to run")
	rootCmd.Flags().StringVarP(&executionFile, "execution", "e", "", "how to execute the test cases")
	rootCmd.Flags().StringVarP(&image, "image", "i", "", "image name of test framework that should exist in system")
	rootCmd.Flags().StringVarP(&configDir, "config_dir", "c", "", "path to external configuration files that are passed to test framework")
	rootCmd.Flags().IntVarP(&queueLength, "queue-length", "q", 5, "Queue length, number of test cases to run parallelly")
	rootCmd.Flags().IntVarP(&retry, "retry", "r", 0, "number of times to retry the failed test cases")
	rootCmd.PersistentFlags().BoolVarP(&junitXML, "junit-xml", "j", false, "Generate JUnit XML report")
	rootCmd.MarkFlagRequired("image")
	cobra.OnInitialize(validateFlags)

}

func runCmd(cmd *cobra.Command, args []string) {
	logger.Infof("Queue length: %d", queueLength)
	engine.RunEngine(executionFile, configDir, positiveTestCases, negativeTestCases, image, queueLength, retry, junitXML)
}

func validateFlags() {
	err := validate.Flags(rootCmd, nil)
	if err != nil {
		logger.Errorf("An error occurred: %v", err)
		os.Exit(1)
	}
}
