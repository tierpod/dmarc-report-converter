package dmarc

import "reflect"

// isRecordsEqual compares two records (exclude Row.Count from comparation)
func isRecordsEqual(r1, r2 Record) bool {
	// make copy of source structs. Be careful about reference types!
	rr1 := r1
	rr2 := r2

	// set Row.Count to zero value (exclude this values from comparation)
	rr1.Row.Count = 0
	rr2.Row.Count = 0

	// compare copied structs
	if reflect.DeepEqual(rr1, rr2) {
		return true
	}

	return false
}

// MergeRecord merges new record rec to the report r.
func (r *Report) MergeRecord(rec Record) {
	for i, record := range r.Records {
		if isRecordsEqual(record, rec) {
			curCount := record.Row.Count
			r.Records[i].Row.Count = curCount + rec.Row.Count
			return
		}
	}

	r.Records = append(r.Records, rec)
}

// MergeReport merges another report rep to the report r. Keeps the earliest Begin date and the
// oldest End date. Merges all records from report rep to the r. Doesn't touch another r fields.
func (r *Report) MergeReport(rep Report) {
	if rep.ReportMetadata.DateRange.Begin.Before(r.ReportMetadata.DateRange.Begin.Time) {
		r.ReportMetadata.DateRange.Begin = rep.ReportMetadata.DateRange.Begin
	}

	if rep.ReportMetadata.DateRange.End.After(r.ReportMetadata.DateRange.End.Time) {
		r.ReportMetadata.DateRange.End = rep.ReportMetadata.DateRange.End
	}

	for _, record := range rep.Records {
		r.MergeRecord(record)
	}
}
