package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/hashicorp/go-version"
)

const staticVersion string = "1.0.0"
const fortiUrl string = "https://docs.fortinet.com"
const fortidocument string = fortiUrl + "/document/fortigate/"
const fileCSV string = "records.csv"
const fileUnDuplicated string = "final.csv"
const headerRow string = "BugID,Description,Status,Version\n"

type fortiTable struct {
	BugID       string `header:"Bug ID"`
	Description string `header:"Description"`
	Status      string
	Version     string
}
type arrayVersions []string

var versions arrayVersions

func (i *arrayVersions) String() string {
	return "my string representation"
}

func (i *arrayVersions) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func getUrlIssues(version string) (string, string) {
	log.Printf("Starting gathering links for version %s\n", version)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// https://docs.fortinet.com/document/fortigate/6.4.0/fortios-release-notes
	url := fortidocument + version + "/fortios-release-notes"

	// Make request
	response, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	bodyString := ""

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString = string(bodyBytes)
	}

	r, _ := regexp.Compile("href=.*known-issues")
	knownIssues := r.FindString(bodyString)
	knownIssues = strings.ReplaceAll(knownIssues, "href=\"", "")

	r, _ = regexp.Compile("href=.*resolved-issues")
	resolvedIssues := r.FindString(bodyString)
	resolvedIssues = strings.ReplaceAll(resolvedIssues, "href=\"", "")

	knownIssuesUrl := fortiUrl + knownIssues
	resolvedIssuesUrl := fortiUrl + resolvedIssues

	log.Printf("The knownIssuesUrl is %s\n", knownIssuesUrl)
	log.Printf("The resolvedIssuesUrl is %s\n", resolvedIssuesUrl)

	return knownIssuesUrl, resolvedIssuesUrl
}

func returnTable(url string, version string, status string) []fortiTable {
	log.Printf("Starting parsing %s data to table", url)
	var table []fortiTable

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("tbody tr").Each(func(_ int, tr *goquery.Selection) {

		e := fortiTable{}

		tr.Find("td").Each(func(ix int, td *goquery.Selection) {
			switch ix {
			case 0:
				e.BugID = strings.TrimSpace(td.Text())
			case 1:
				// Source
				// https://stackoverflow.com/questions/37290693/how-to-remove-redundant-spaces-whitespace-from-a-string-in-golang
				desc := strings.TrimSpace(td.Text())
				desc = strings.Join(strings.Fields(desc), " ")
				e.Description = strings.TrimSpace(desc)
			}
		})

		e.Version = version
		e.Status = status

		table = append(table, e)
	})
	return table
}

func getResolvedIssues(resolvedIssuesUrl string, version string) []fortiTable {
	log.Printf("Getting the resolved issue for version %s", version)
	resolvedIssues := returnTable(resolvedIssuesUrl, version, "resolved")
	return resolvedIssues
}

func getKnownIssues(knownIssuesUrl string, version string) []fortiTable {
	log.Printf("Getting the known issue for version %s", version)
	knownIssues := returnTable(knownIssuesUrl, version, "unresolved")
	return knownIssues
}

func writeToCSV(table []fortiTable, fileName string) {
	log.Printf("Starting writing the to file %s", fileName)

	var isCreated bool

	if _, err := os.Stat(fileName); err == nil {
		isCreated = false
	} else {
		isCreated = true
	}

	// Source
	// https://articles.wesionary.team/read-and-write-csv-file-in-go-b445e34968e9

	file, err := os.OpenFile(fileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	if isCreated {
		file.WriteString(headerRow)
	}

	defer file.Close()

	if err != nil {
		log.Fatalln("Failed to open file", err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()

	for _, record := range table {
		row := []string{record.BugID, record.Description, record.Status, record.Version}
		if err := w.Write(row); err != nil {
			log.Fatalln("Error writing record to file", err)
		}
	}
}

func createFortiList(data [][]string) []fortiTable {
	var table []fortiTable
	for i, line := range data {
		if i > 0 { // omit header line
			var rec fortiTable
			for j, field := range line {
				switch j {
				case 0:
					rec.BugID = field
				case 1:
					rec.Description = field
				case 2:
					rec.Status = field
				case 3:
					rec.Version = field
				}
			}
			table = append(table, rec)
		}
	}
	return table
}

func removeDuplicates(input string, output string) {
	f, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	listOrg := createFortiList(data)

	var listFinal []fortiTable
	var listPreFinal []fortiTable
	var listWithoutDuplicates []fortiTable

	visited := make(map[fortiTable]bool, 0)
	visitedFinal := make(map[fortiTable]bool, 0)

	// Removing a duplicates
	for _, item := range listOrg {
		if _, value := visited[item]; !value {
			visited[item] = true
			listWithoutDuplicates = append(listWithoutDuplicates, item)
		}
	}

	// Logic for issues
	versionList := []string{}
	for _, itemOrg := range listOrg {
		finalItem := fortiTable{}
		for _, item := range listWithoutDuplicates {
			if item.BugID == itemOrg.BugID {
				versionList = append(versionList, item.Version)
			}
		}
		sort.Slice(versionList, func(i, j int) bool {
			v1, _ := version.NewVersion(versionList[i])
			v2, _ := version.NewVersion(versionList[j])
			return v1.GreaterThanOrEqual(v2)
		})

		for _, item := range listWithoutDuplicates {
			if item.BugID == itemOrg.BugID && versionList[0] == item.Version {
				finalItem.BugID = item.BugID
				finalItem.Description = item.Description
				finalItem.Version = item.Version
				finalItem.Status = item.Status
			}
		}
		listPreFinal = append(listPreFinal, finalItem)
		versionList = nil

	}

	// Removing a duplicates
	for _, item := range listPreFinal {
		if _, value := visitedFinal[item]; !value {
			visitedFinal[item] = true
			listFinal = append(listFinal, item)
		}
	}

	writeToCSV(listFinal, output)
}

func main() {
	flag.Var(&versions, "version", "Version(s) of the FortiOS")
	recordsFile := flag.String("recordsFile", fileCSV, "Name of the unsorted records from versions")
	sorted := flag.Bool("sorted", false, "Get a sorted release notes")
	sortedFile := flag.String("sortedFile", fileUnDuplicated, "Name of the sorted output file")
	flag.Parse()

	for _, version := range versions {
		knownIssuesUrl, resolvedIssuesUrl := getUrlIssues(version)
		resolvedIssues := getResolvedIssues(resolvedIssuesUrl, version)
		writeToCSV(resolvedIssues, *recordsFile)
		knownIssuses := getKnownIssues(knownIssuesUrl, version)
		writeToCSV(knownIssuses, *recordsFile)
	}

	if *sorted {
		removeDuplicates(*recordsFile, *sortedFile)
	}
}
