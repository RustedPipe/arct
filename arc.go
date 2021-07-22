// Package arc provides functionality for parsing arctool
// log files and generating reports in csv format.
package arc

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// GenerateReport takes a path to the logfile and generate
// csv report in the given output file. If the operation
// is not successful it returns an error.
func GenerateReport(filein, fileout string) error {
	fin, err := os.Open(filein)
	if err != nil {
		return fmt.Errorf("opening log file: %s, err: %v", filein, err)
	}
	defer fin.Close()

	fout, err := os.Create(fileout)
	if err != nil {
		return fmt.Errorf("creating output file: %s, err: %v", fileout, err)
	}
	defer fout.Close()

	if err := processData(fin, fout); err != nil {
		return fmt.Errorf("creating data file: %s from log file: %s, err: %v", fileout, filein, err)
	}

	return nil
}

func processData(r io.Reader, w io.Writer) error {
	csvwriter := csv.NewWriter(w)
	csvheader := []string{"Sr.No", "WPRN", "PremiseID"}
	csvwriter.Write(csvheader)

	// Variables to hold values extracted from a log line.
	var srno, wprn, premiseid string
	linesNumber := 1

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		l := scanner.Text()
		if !strings.HasPrefix(l, "Sr.No") {
			continue
		}

		linesNumber++

		// We need to clean ';' to prepare string for Sscanf function.
		l = strings.ReplaceAll(l, ";", "")
		_, err := fmt.Sscanf(l, "Sr.No = %s WPRN = %s PremiseID = %s", &srno, &wprn, &premiseid)
		if err != nil {
			return fmt.Errorf("processing log line: %s, %v", l, err)
		}

		if err := csvwriter.Write([]string{srno, wprn, premiseid}); err != nil {
			return fmt.Errorf("writing line to csv file: %v", []string{srno, wprn, premiseid})
		}
	}

	// We are not interested in an empty csv file with a header only.
	// So, we return error and abandon creating an empty csv file.
	if linesNumber == 1 {
		return errors.New("no data in the input file")
	}
	csvwriter.Flush()
	if err := csvwriter.Error(); err != nil {
		return fmt.Errorf("writing csv file: %v", err)
	}

	return nil
}
