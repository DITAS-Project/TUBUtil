//    Copyright 2018 TUB/*  */
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()
var log = logrus.NewEntry(logger)

func SetLogger(nLogger *logrus.Logger) {
	logger = nLogger
}

func SetLog(entty *logrus.Entry) {
	log = entty
}

func WaitForAvailible(url string, maxTimeout *time.Duration) error {
	start := time.Now()
	for {
		resp, err := http.Head(url)

		if err == nil && resp.StatusCode == 200 {
			break
		}

		if logger.Level == logrus.DebugLevel {
			if err != nil {
				log.Debugf("elastic unavailible %v", err)
			}

			if resp != nil {
				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Debugf("elastic error [%d] - %s", resp.StatusCode, string(data))
				} else {
					log.Debugf("elastic error [%d]", resp.StatusCode)
				}

			}
		}

		log.Info("ElasticSearch not availible - wating")

		if maxTimeout != nil && time.Now().Sub(start) > *maxTimeout {
			return errors.New("could not connect toElasticSearch in time")
		}

		time.Sleep(time.Duration(1e+10)) //10 seconds
	}

	return nil
}

func WaitForGreen(es *elastic.Client, maxTimeout *time.Duration) error {
	start := time.Now()
	for {
		log.Info("check if elastic is ready")
		err := es.WaitForGreenStatus("10s")
		if err == nil {
			break
		}
		if maxTimeout != nil && time.Now().Sub(start) > *maxTimeout {
			return errors.New("could not connect toElasticSearch in time")
		}
		log.Info("ElasticSearch not ready - wait")
	}

	return nil
}

func GetElasticIndex(vdcName string) string {
	t := time.Now()
	return fmt.Sprintf("%s-%d-%02d-%02d", vdcName, t.Year(), t.Month(), t.Day())
}
