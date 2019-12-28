package repository

type Result struct {
	// Error is the summary level error.
	Error error
	// ErrorDetails contains the details of the error.
	// Sample values assigned are array of UpsertFailDetail,
	// array of DeleteFailDetail if result means partial failure
	// on Set or Remove, or any 3rd party details/error info.
	ErrorDetails interface{}
}

type UpsertFailDetail struct {
	KeyValue KeyValue
	Error    error
}
type DeleteFailDetail struct {
	Key   string
	Error error
}

// IsSuccessful returns false if result status signifies failure.
func (result Result) IsSuccessful() bool {
	return result.Error == nil
}
