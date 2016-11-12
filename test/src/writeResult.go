// Copyright 2016 Etix Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

////////////////////////////////////////////////
// Data declarations

// JUnitTestSuites is a collection of JUnit test suites.
type JUnitTestSuites struct {
	XMLName    xml.Name         `xml:"testsuites"`
	TestSuites []JUnitTestSuite `xml:"testsuite"`
}

// JUnitTestSuite is a single JUnit test suite which may contain many
// testcases.
type JUnitTestSuite struct {
	XMLName   xml.Name        `xml:"testsuite"`
	Tests     int             `xml:"tests,attr"`
	Failures  int             `xml:"failures,attr"`
	Time      string          `xml:"time,attr"`
	TestCases []JUnitTestCase `xml:"testcase"`
}

// JUnitTestCase is a single test case with its result.
type JUnitTestCase struct {
	XMLName xml.Name      `xml:"testcase"`
	Message string        `xml:"message,attr"`
	Time    string        `xml:"time,attr"`
	Failure *JUnitFailure `xml:"failure,omitempty"`
}

// JUnitFailure contains data related to a failed test.
type JUnitFailure struct {
	XMLName  xml.Name `xml:"failure"`
	Message  string   `xml:"message,attr"`
	Type     string   `xml:"type,attr"`
	Contents string   `xml:",chardata"`
}

// WriteResults will output the results in the standard output as well as concatenate them in an XML JUnit report
func (t *Tester) WriteResults(result Test, output string) bool {
	fmt.Printf("Displaying results...\n")
	t.writeConsoleReport(result)

	file, err := os.OpenFile(output, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error opening XML: %s\n", err)
		return false
	}
	defer file.Close()

	err = t.writeJUnitReportXML(result, file, output)
	if err != nil {
		fmt.Printf("Error writing XML: %s\n", err)
		return false
	}
	fmt.Printf("-> JUnit XML report written: %s\n", output)
	return true
}

// Write tests results under JUnit format on w
func (t *Tester) writeJUnitReportXML(result Test, rw io.ReadWriter, output string) error {
	suites := JUnitTestSuites{}

	buf, err := ioutil.ReadFile(output)

	dec := xml.NewDecoder(bytes.NewBufferString(string(buf)))
	err = dec.Decode(&suites)
	if err != nil {
		fmt.Printf("\nUnable to deserialize %s file: %s\n", output, err)
	}

	ts := JUnitTestSuite{
		Tests:     len(result.result) + len(result.expected),
		Failures:  0,
		Time:      fmt.Sprintf("%.6f", result.time.Seconds()),
		TestCases: []JUnitTestCase{},
	}

	for _, r := range result.result {
		testCase := JUnitTestCase{
			Time:    fmt.Sprintf("%.6f", result.time.Seconds()),
			Failure: nil,
		}
		testCase.Message = "The stream " + r.Address + " could be accessed and its thumbnail was properly generated"
		ts.TestCases = append(ts.TestCases, testCase)
	}

	for _, e := range result.expected {
		testCase := JUnitTestCase{
			Time:    fmt.Sprintf("%.6f", result.time.Seconds()),
			Failure: nil,
		}
		if e.err != nil {
			testCase.Failure = &JUnitFailure{
				Message: e.err.Error(),
				Type:    "",
			}
		}
	}

	suites.TestSuites = append(suites.TestSuites, ts)
	// Fix indent
	bytes, err := xml.MarshalIndent(suites, "", "\t")
	if err != nil {
		return err
	}
	// Write in param stream

	w, err := os.OpenFile(output, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	writer := io.Writer(w)
	writer.Write(bytes)
	return nil
}

func (t *Tester) writeConsoleReport(result Test) bool {
	successCount := len(result.result)
	failureCount := len(result.expected)
	fmt.Println("--- Test summary ---")
	if successCount > 0 {
		fmt.Printf("Results: %d/%d (%d%%)\n", successCount, successCount+failureCount, successCount*100/(successCount+failureCount))
		fmt.Printf("Time: %.6fs\n", result.time.Seconds())
	} else {
		fmt.Printf("No test in success\n")
	}

	return true
}
