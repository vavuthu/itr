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

package mail

import (
	"os"
	"path/filepath"

	"github.com/vavuthu/itr/config"
	"github.com/wneessen/go-mail"
	"github.com/vavuthu/itr/cmd/utils"
	"github.com/vavuthu/itr/logger"
	
)

const (
	subjectFile = "description.txt"
)

func SendMail() {
	emailID := config.AppConfig.EmailID
	htmlReport := "report_" + config.AppConfig.RunID + ".html"
	htmlContent, err := os.ReadFile(htmlReport)
	if err != nil {
		logger.Errorf("Failed to read HTML file: %s", htmlReport)
	}

	m := mail.NewMsg()
	if err := m.From("ocs-ci@redhat.com"); err != nil {
		logger.Errorf("failed to set From address: %s", err)
	}
	if err := m.To(emailID); err != nil {
		logger.Errorf("failed to set To address: %s", err)
	}

	subject := config.AppConfig.Subject
	defaultSubject := "[ITR RUN ID: " + config.AppConfig.RunID + "]"

	if subject != "" {
		defaultSubject = "[ITR RUN ID: " + config.AppConfig.RunID + "] " + subject
		
	} else {
		subjectFilePath := filepath.Join(config.AppConfig.ConfigDir, subjectFile)
		isSubjectFilePathExist := utils.CheckFileExists(subjectFilePath)
		if isSubjectFilePathExist {
			content, err := os.ReadFile(subjectFilePath)
            if err != nil {
                logger.Errorf("Error reading file: %s", err)
            } else {
				defaultSubject = "[ITR RUN ID: " + config.AppConfig.RunID + "] " + string(content)
            }

		}

	}

	m.Subject(defaultSubject)
    
	m.SetBodyString(mail.TypeTextHTML, string(htmlContent))
	if err := m.WriteToSendmail(); err != nil {
        logger.Errorf("failed to write mail to local sendmail: %s", err)
    }
	
    
}
