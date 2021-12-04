// Copyright 2021 - Lars Gohr
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
