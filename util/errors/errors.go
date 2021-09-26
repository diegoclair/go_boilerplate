package errors

import (
	"regexp"
)

//noSQLRowsRE - to check if sql error is because that there are no rows
var noSQLRowsRE = regexp.MustCompile(noSQLRows)

//noRecordsFindRE - to check if sql error is because that there are no records find with the parameters
var noRecordsFindRE = regexp.MustCompile(noRecordsFind)

const noSQLRows string = "no rows in result set"
const noRecordsFind string = "No records find"

//SQLNotFound - Check if the error is because there are no sql rows or
//no records find with given parameters
func SQLNotFound(err string) bool {
	noRowsIdx := noSQLRowsRE.FindStringIndex(err)
	if len(noRowsIdx) > 0 {
		return true
	}

	noRecorsIdx := noRecordsFindRE.FindStringIndex(err)

	return len(noRecorsIdx) > 0
}
