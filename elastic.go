/*
 * Copyright 2018 Information Systems Engineering, TU Berlin, Germany
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *                     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * This is being developed for the DITAS Project: https://www.ditas-project.eu/
 */

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

//SetLogger can be used to use a custom logrus for this package
func SetLogger(nLogger *logrus.Logger) {
	logger = nLogger
}

//SetLog can be used to use a custom logrus for this package
func SetLog(entty *logrus.Entry) {
	log = entty
}

//WaitForAvailible utility function that tries to connect to the given url for a given time. Using HTTP Head. Fails for all but 200 responses.
func WaitForAvailible(url string, maxTimeout *time.Duration) error {
	return WaitForAvailibleWithAuth(url, nil, maxTimeout)
}

//WaitForAvailibleWithAuth utility function that tries to connect to the given url for a given time. Using HTTP Head. Fails for all but 200 responses.
func WaitForAvailibleWithAuth(url string, auth []string, maxTimeout *time.Duration) error {
	start := time.Now()
	for {

		req, err := http.NewRequest("HEAD", url, nil)

		if err != nil {
			continue
		}

		if auth != nil && len(auth) == 2 {
			req.SetBasicAuth(auth[0], auth[1])
		}

		resp, err := http.DefaultClient.Do(req)

		if err == nil && resp.StatusCode == 200 {
			break
		}

		if logger.Level == logrus.DebugLevel {
			if err != nil {
				log.Debugf("%s unavailible %v", url, err)
			}

			if resp != nil {
				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Debugf("error [%d] - %s", resp.StatusCode, string(data))
				} else {
					log.Debugf("error [%d]", resp.StatusCode)
				}

			}
		}

		log.Info("not availible - wating")

		if maxTimeout != nil && time.Now().Sub(start) > *maxTimeout {
			return fmt.Errorf("could not connect to %s in time", url)
		}

		time.Sleep(time.Duration(1e+10)) //10 seconds
	}

	return nil
}

//WaitForGreen blocks until an elasticserach instance reaches the green status.
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

//GetElasticIndex calculates the current VDC-Index-Name
func GetElasticIndex(vdcName string) string {
	t := time.Now()
	return fmt.Sprintf("%s-%d-%02d-%02d", vdcName, t.Year(), t.Month(), t.Day())
}

//GetElasticIndexByIDs retuns the current VDC-Index-Name using Id's
func GetElasticIndexByIDs(vdcID, blueprintID string) string {
	return GetElasticIndex(fmt.Sprintf("%s-%s", blueprintID, vdcID))
}
