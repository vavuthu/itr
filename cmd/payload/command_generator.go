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
	"strings"
)

const (
	prefix = "podman run -e BUILD_NUMBER -e BUILD_TAG -e BUILD_URL -e JOB_NAME -e NODE_NAME -e WORKSPACE --rm -v "
	PodmanPath = "/opt/cluster"
	podmanSharedMountOption = "z "
	redirectionOperator = " >"
	redirectSTDERROUT = " 2>&1"
	background = " &"
)

// Forms the podman command for test case
func CommandGenerator(testCase, testCaseName, image, configDir string) string {
	command := prefix + configDir + ":" + PodmanPath + ":" + podmanSharedMountOption + image + " " + strings.TrimSuffix(testCase, "\n")
	return command
}
