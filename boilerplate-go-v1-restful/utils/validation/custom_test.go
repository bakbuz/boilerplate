package validation_test

import (
	"codegen/utils/validation"
	"testing"

	"github.com/go-playground/validator/v10"
)

type person struct {
	FirstName      string `validate:"required,max=10"`
	LastName       string `validate:"required,max=10"`
	Email          string `validate:"required,email,max=64"`
	Age            uint8  `validate:"gte=0,lte=130"`
	FavouriteColor string `validate:"iscolor"` // alias for 'hexcolor|rgb|rgba|hsl|hsla'
}

func TestCustomValidate_Fail(t *testing.T) {
	e := &person{
		FirstName: "",
		LastName:  "1234567890_000",
	}

	cv := validation.NewCustomValidator()

	// shouldn't pass
	if err := cv.Validate(e); err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
}

func TestCustomValidate_Fail_ErrorCode(t *testing.T) {
	e := &person{
		FirstName: "",
		LastName:  "1234567890_000",
	}

	cv := validation.NewCustomValidator()

	// shouldn't pass
	if err := cv.Validate2(e); err != nil {
		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			t.Log(e.Namespace())
			t.Log(e.Field())
			t.Log(e.StructNamespace())
			t.Log(e.StructField())
			t.Log(e.Tag())
			t.Log(e.ActualTag())
			t.Log(e.Kind())
			t.Log(e.Type())
			t.Log(e.Value())
			t.Log(e.Param())
		}

		//t.Log(err.Error())
		t.FailNow()
	}
}

func TestCustomValidate_Pass(t *testing.T) {
	e := &person{
		FirstName:      "Deli",
		LastName:       "Dumrul",
		Email:          "a@a.com",
		Age:            33,
		FavouriteColor: "#000000",
	}

	cv := validation.NewCustomValidator()

	// should pass
	if err := cv.Validate(e); err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
}
