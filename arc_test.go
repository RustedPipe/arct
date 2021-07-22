package arc

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func content() io.Reader {
	return strings.NewReader(`—————————————START————————————19/02/2021 12:51:06
File loading is started for file dfpub0457_DLRZone1_residential_2021-02-18.csv
Sr.No = SKSZPUB0257-1; WPRN = 2607303; PremiseID = 2306982
Sr.No = SKSZPUB0257-2; WPRN = 2607304; PremiseID = 3104983
Sr.No = SKSZPUB0257-3; WPRN = 2607305; PremiseID = 5616984
Sr.No = SKSZPUB0257-4; WPRN = 2607306; PremiseID = 1626985
Total number of records to be processed 4
Total number of processed records successfully 4
Total number of failed records 0
—————————————END—————————————19/02/2021 12:52:21
—————————————START————————————19/02/2021 12:53:48
File loading is started for file dfpub0441_DLRZone4_residential_2021-02-18.csv
Sr.No = ALSZPUB0241-1; WPRN = 1507307; PremiseID = 2601986
Sr.No = ALSZPUB0241-2; WPRN = 1507308; PremiseID = 2601987
Sr.No = ALSZPUB0241-3; WPRN = 1507309; PremiseID = 2601988
Total number of records to be processed 3
Total number of processed records successfully 3
Total number of failed records 0
—————————————END—————————————19/02/2021 12:58:21`)
}

func contentNoRecords() io.Reader {
	return strings.NewReader(`—————————————START————————————19/02/2021 12:51:06
File loading is started for file dfpub0457_DLRZone1_residential_2021-02-18.csv
Total number of records to be processed 0
Total number of processed records successfully 0
Total number of failed records 0
—————————————END—————————————19/02/2021 12:52:21
—————————————START————————————19/02/2021 12:53:48
File loading is started for file dfpub0441_DLRZone4_residential_2021-02-18.csv
Total number of records to be processed 0
Total number of processed records successfully 0
Total number of failed records 0
—————————————END—————————————19/02/2021 12:58:21`)
}

func TestProcessData(t *testing.T) {
	wantCSV := `Sr.No,WPRN,PremiseID
SKSZPUB0257-1,2607303,2306982
SKSZPUB0257-2,2607304,3104983
SKSZPUB0257-3,2607305,5616984
SKSZPUB0257-4,2607306,1626985
ALSZPUB0241-1,1507307,2601986
ALSZPUB0241-2,1507308,2601987
ALSZPUB0241-3,1507309,2601988
`

	out, err := ioutil.TempFile(".", "fout.csv")
	if err != nil {
		t.Fatalf("can't create a temp file, got error: %s", err)
	}
	defer os.Remove(out.Name())

	err = processData(content(), out)
	if err != nil {
		t.Errorf("processData() got error: %s", err)
	}

	got, err := ioutil.ReadFile(out.Name())
	if err != nil {
		t.Errorf("reading from out file")
	}

	if !cmp.Equal(string(got), wantCSV) {
		t.Errorf("processData() \n%s", cmp.Diff(string(got), wantCSV))
	}
}

func TestProcessData_Empty(t *testing.T) {
	out, err := ioutil.TempFile(".", "fout.csv")
	if err != nil {
		t.Fatalf("can't create a temp file, got error: %s", err)
	}
	defer os.Remove(out.Name())

	wantError := "no data in the input file"
	err = processData(contentNoRecords(), out)
	if err != nil {
		if !cmp.Equal(err.Error(), wantError) {
			t.Errorf("processData() got: \n%s", cmp.Diff(err.Error(), wantError))
		}
	}
}
