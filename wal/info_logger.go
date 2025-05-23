/*	Copyright [2024] [Elia Renzoni]
*
*   Licensed under the Apache License, Version 2.0 (the "License");
*   you may not use this file except in compliance with the License.
*   You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*   Unless required by applicable law or agreed to in writing, software
*   distributed under the License is distributed on an "AS IS" BASIS,
*   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*   See the License for the specific language governing permissions and
*   limitations under the License.
*/




package wal

import (
	"log"
	"os"
)

type InfoLogger struct {
	logger *log.Logger
}

func NewInfoLogger() *InfoLogger {
	return &InfoLogger{
		logger: log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
	}
}

func (i *InfoLogger) ReportInfo(info string) {
	i.logger.Println(info)
}
