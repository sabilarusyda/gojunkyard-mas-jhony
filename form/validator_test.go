package form

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateFlagFailedFilter(t *testing.T) {
	var got user

	err := ValidateFlag(Vfilter, got)
	assert.NotNilf(t, err, "error must be exist. err: %s", err)
}

func TestValidateFlagSuccessFilterAndValidation(t *testing.T) {
	var (
		dt   = time.Date(2018, time.December, 1, 1, 1, 1, 0, time.Local)
		got  = user{ID: 1, Name: " Risal Falah ", Email: "risal.falah@gmail.com", BirthDate: dt, Gender: "m"}
		want = user{ID: 1, Name: "Risal Falah", Email: "risal.falah@gmail.com", BirthDate: dt, Gender: "m"}
	)

	err := ValidateFlag(Vfilter, &got)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, want, got)
}

func TestValidateFailedFilter(t *testing.T) {
	var got user

	err := Validate(got)
	assert.NotNilf(t, err, "error must be exist. err: %s", err)
}

func TestValidateSuccessFilterAndValidation(t *testing.T) {
	var (
		dt   = time.Date(2018, time.December, 1, 1, 1, 1, 0, time.Local)
		got  = user{ID: 1, Name: " Risal Falah ", Email: "risal.falah@gmail.com", BirthDate: dt, Gender: "m"}
		want = user{ID: 1, Name: "Risal Falah", Email: "risal.falah@gmail.com", BirthDate: dt, Gender: "m"}
	)

	err := Validate(&got)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, want, got)
}
