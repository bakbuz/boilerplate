package validation_test

import (
	"testing"

	"codegen/utils/validation"
)

type (
	ok struct {
		A string `validate:"required,eq=1"`
	}
	ok1 struct {
		A string `validate:"required,eq=1"`
		B string `validate:"required,eq=1"`
	}
	fail struct {
		A string `validate:"get"`
	}
)

func TestValidate(t *testing.T) {
	// should pass
	if err := validation.Validate(&ok{A: "1"}); err != nil {
		t.FailNow()
	}
	// should fail
	if err := validation.Validate(&ok1{A: "1", B: "y"}); err == nil {
		t.FailNow()
	}
	// should fail
	if err := validation.Validate(&ok{A: "x"}); err == nil {
		t.FailNow()
	}
	// should panic
	defer func() {
		if r := recover(); r == nil {
			t.FailNow()
		}
	}()
	validation.Validate(&fail{A: "x"})
}
