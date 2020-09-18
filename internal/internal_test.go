package internal_test

import (
	"testing"

	"github.com/elgohr/go-localstack/internal"
	"github.com/stretchr/testify/assert"
)

func TestMustParseConstraint(t *testing.T) {
	c := internal.MustParseConstraint(">= 1.0.0")

	assert.NotNil(t, c)
}

func TestMustParseConstraint__Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected the invalid constraint to cause a panic")
		}
	}()

	_ = internal.MustParseConstraint(">=<><>< foo.bba.0")
}
