package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"time"
)

////////////////////////////////////////////////
// Data declarations

// JUnitTestSuites is a collection of JUnit test suites.
type JUnitTestSuites struct {
	XMLName xml.Name `xml:"testsuites"`
	Suites  []JUnitTestSuite
}

// JUnitTestSuite is a single JUnit test suite which may contain many
// testcases.
type JUnitTestSuite struct {
	Tests     int    `xml:"tests,attr"`
	Failures  int    `xml:"failures,attr"`
	Time      string `xml:"time,attr"`
	TestCases []JUnitTestCase
}

// JUnitTestCase is a single test case with its result.
type JUnitTestCase struct {
	Message string        `xml:"message,attr"`
	Time    string        `xml:"time,attr"`
	Failure *JUnitFailure `xml:"failure,omitempty"`
}

// JUnitFailure contains data related to a failed test.
type JUnitFailure struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",chardata"`
}

func (m *manager) WriteResults(result TestCase, output string) bool {
	fmt.Printf("Displaying results...\n")
	// Write Console report
	m.writeConsoleReport(result)

	// Write XML report
	// Open xml
	file, err := os.OpenFile(output, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error opening XML: %s\n", err)
		return false
	}
	defer file.Close()
	err = m.writeJUnitReportXML(result, file, output)
	if err != nil {
		fmt.Printf("Error writing XML: %s\n", err)
		return false
	}
	fmt.Printf("-> JUnit XML report written: %s\n", output)
	return true
}

// Write tests results under JUnit format on w
func (m *manager) writeJUnitReportXML(result TestCase, r io.ReadWriter, output string) error {
	suites := JUnitTestSuites{}
	dec := xml.NewDecoder(r)
	if err := dec.Decode(&suites); err != nil {
		fmt.Printf("\nUnable to deserialize XML log file: %s\n", err)
	}
	ts := JUnitTestSuite{
		Tests:     len(result.result) + len(result.expected),
		Failures:  0,
		Time:      fmt.Sprintf("%.6f", result.time.Seconds()),
		TestCases: []JUnitTestCase{},
	}
	// Run throught all iterations
	testCase := JUnitTestCase{
		Time:    fmt.Sprintf("%.6f", result.time.Seconds()),
		Failure: nil,
	}
	if len(result.result) > 0 {
		testCase.Message = "These streams matched what we expected:"
	}
	for _, success := range result.result {
		testCase.Message += " " + success.Address
	}
	if !result.ok {
		testCase.Failure = &JUnitFailure{
			Message: "These streams did not match what we expected:",
			Type:    "",
		}
	}
	for _, fail := range result.expected {
		ts.Failures++
		testCase.Failure.Message += " " + fail.Address
	}
	ts.TestCases = append(ts.TestCases, testCase)

	suites.Suites = append(suites.Suites, ts)
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

func (m *manager) writeConsoleReport(result TestCase) bool {
	min := 50 * time.Hour
	max := 0 * time.Second
	total := 0 * time.Second
	successCount := 0
	failureCount := 0
	if result.ok {
		successCount++
		total += result.time
		if result.time < min {
			min = result.time
		}
		if result.time > max {
			max = result.time
		}
	} else {
		failureCount++
	}
	fmt.Println("--- Test summary ---")
	if successCount > 0 {
		fmt.Printf("Results: %d/%d (%d%%)\n", successCount, successCount+failureCount, successCount*100/(successCount+failureCount))
		fmt.Printf("Total time: %.6fs\n", total.Seconds())
		fmt.Printf("Average time: %.6fs\n", total.Seconds()/float64(successCount))
		fmt.Printf("Min time: %.6fs\n", min.Seconds())
		fmt.Printf("Max time: %.6fs\n", max.Seconds())
	} else {
		fmt.Printf("No test in success\n")
	}

	return true
}
