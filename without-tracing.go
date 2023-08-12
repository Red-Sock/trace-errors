//go:build rscliErrorTracingDisabled
// +build rscliErrorTracingDisabled

package errors

func init() {
	enableTracing = false
}
