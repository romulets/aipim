package domain

import (
	_ "embed" // embed templates
	"testing"
)

//go:embed testdata/tostring_output/generated_basic.painless
var testSimplePainless string

func TestSimpleParser(t *testing.T) {
	clm := &cloudtrailLogMapping{}
	clm.scan(testSimplePainless)
}
