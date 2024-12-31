/*
Copyright © 2024 vavuthu@redhat.com

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

package config

type Config struct {
	ConfigDir string
	EmailID string
	RunID string
	Subject string
	Retry int
	Env map[string]interface{} // For dynamic parameters
}

var AppConfig = Config{
    Env: make(map[string]interface{}), // Initialize Extras to avoid nil map
}

func InitializeConfig(retry int, email string, runid string, configdir string, subject string, env map[string]interface{}) {
	AppConfig.Retry = retry
	AppConfig.EmailID = email
	AppConfig.RunID = runid
	AppConfig.ConfigDir = configdir
	AppConfig.Subject = subject
	AppConfig.Env = env
}

// UpdateConfigEnv allows updating the Env field dynamically
func UpdateConfigEnv(key string, value interface{}) {
    if AppConfig.Env == nil {
        AppConfig.Env = make(map[string]interface{}) // Ensure the map is initialized
    }
    AppConfig.Env[key] = value
}

