package main

import (
	"fmt"
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

func groupReportsKey(r dmarc.Report) string {
	return fmt.Sprintf("%v!%v!%v", r.ReportMetadata.OrgName, r.ReportMetadata.Email, r.PolicyPublished.Domain)
}

func groupMergeReports(reports []dmarc.Report) ([]dmarc.Report, error) {
	if len(reports) == 0 {
		return nil, fmt.Errorf("reports list is empty")
	}

	if len(reports) == 1 {
		return reports, nil
	}

	grouped := make(map[string][]dmarc.Report)

	// group reports by key
	for _, r := range reports {
		key := groupReportsKey(r)
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
