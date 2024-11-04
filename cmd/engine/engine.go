
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

package engine

import (
	"github.com/vavuthu/itr/cmd/launcher"
	"github.com/vavuthu/itr/cmd/payload"
	"github.com/vavuthu/itr/cmd/statusquo"
	"github.com/vavuthu/itr/cmd/utils"
	"github.com/vavuthu/itr/logger"
)

func RunEngine(execution, configDir, nonDisruptiveTestCases, disruptiveTestCases, image string, queueLength, retry int, junitXML bool) {
	logger.Info("Starting ITR engine")
	
	filenames := []string{nonDisruptiveTestCases, disruptiveTestCases}
	totalCount, err := utils.CountLines(filenames)
	if err != nil {
			logger.Error("error in reading file %s ", err)
			return
	}
	statusquo.TotalTestCases = totalCount

	if len(nonDisruptiveTestCases) != 0 {
		RunEngineParallely(execution, configDir, nonDisruptiveTestCases, image, queueLength, retry, junitXML)
	}

	if len(disruptiveTestCases) != 0 {
		RunEngineSerially(execution, configDir, disruptiveTestCases, image, queueLength, retry, junitXML)
	}
}

func RunEngineParallely(execution, configDir, nonDisruptiveTestCases, image string, queueLength, retry int, junitXML bool) {
	logger.Info("Running engine parallely")
	commands := payload.GenerateAllPodmanCommands(execution, configDir, nonDisruptiveTestCases, image, junitXML)
	launcher.LaunchInitiate(commands, configDir, queueLength, retry)
}

func RunEngineSerially(execution, configDir, disruptiveTestCases, image string, queueLength, retry int, junitXML bool) {
	logger.Info("Running engine serially")
	commands := payload.GenerateAllPodmanCommands(execution, configDir, disruptiveTestCases, image, junitXML)
	queueLength = 1
	launcher.LaunchInitiate(commands, configDir, queueLength, retry)
}
