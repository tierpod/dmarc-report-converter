package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/tierpod/dmarc-report-converter/pkg/dmarc"
)

func mergeReports(reports []dmarc.Report) dmarc.Report {
	first := reports[0]
	for _, rep := range reports[1:] {
		first.MergeReport(rep)
	}

	// resort records and recalculate message statistic after merge
	first.SortRecords()
	first.CalculateStats()
	return first
}

func groupReportsKey(r dmarc.Report, t *template.Template) (string, error) {
	var buf bytes.Buffer

	err := t.Execute(&buf, r)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func groupMergeReports(reports []dmarc.Report, t *template.Template) ([]dmarc.Report, error) {
	if len(reports) == 0 {
		return nil, fmt.Errorf("reports list is empty")
	}

	if len(reports) == 1 {
		return reports, nil
	}

	grouped := make(map[string][]dmarc.Report)

	// group reports by key
	for _, r := range reports {
		key, err := groupReportsKey(r, t)
		if err != nil {
			return reports, fmt.Errorf("error generating merge key: %s", err)
		}

		//lint:ignore S1036 we dont want to add nil value
		if _, found := grouped[key]; found {
			grouped[key] = append(grouped[key], r)
		} else {
			grouped[key] = []dmarc.Report{r}
		}
	}

	// flatten grouped reports to slice of reports
	var result []dmarc.Report
	for k, r := range grouped {
		log.Printf("[INFO] merge: %v report(s), grouped by key '%v'", len(r), k)
		result = append(result, mergeReports(r))
	}

	return result, nil
}
