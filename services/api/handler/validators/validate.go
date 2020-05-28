package validators

import (
	"errors"
	"fmt"
	"github.com/lexbond13/api_core/util"
	"github.com/araddon/dateparse"
	"regexp"
	"strings"
)

var RegexpIsEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validate struct {
	errors []error
}

// Add
func (ve *Validate) Add(err error) {
	ve.errors = append(ve.errors, err)
}

// Count
func (ve *Validate) Count() int {
	if ve.errors == nil {
		return 0
	}
	return len(ve.errors)
}

// Errors
func (ve *Validate) Errors() []string {
	errs := make([]string, 0, len(ve.errors))
	for _, err := range ve.errors {
		errs = append(errs, err.Error())
	}
	return errs
}

func (ve *Validate) String() string {
	return strings.Join(ve.Errors(), ",")
}

// Required error for when a value is missing
func (ve *Validate) Required(field, value interface{}) {
	if value == "" {
		ve.Add(errors.New(fmt.Sprintf("%s is required", field)))
	}
}

// IsEmail
func (ve *Validate) IsEmail(email string) {
	if !RegexpIsEmail.MatchString(email) {
		ve.Add(errors.New(fmt.Sprintf("%s is not correct email address", email)))
	}
}

// IsDate
func (ve *Validate) IsDate(date string) {
	_, err := dateparse.ParseAny(date)
	if err != nil {
		ve.Add(errors.New(fmt.Sprintf("%s is not correct date", date)))
	}
}

// Size
func (ve *Validate) Size(size, allow int64) {
	if size > allow {
		ve.Add(errors.New(fmt.Sprintf("max image size allow: %dMB.", util.ConvertCountBytesToCountMegabytes(allow))))
	}
}

// Extension
func (ve *Validate) Extension(extension string, allow []string) {
	if len(allow) > 0 {
		for _, ext := range allow {
			if strings.TrimSpace(ext) == extension {
				return
			}
		}
		ve.Add(errors.New(fmt.Sprintf("only extensions %s allowed.", strings.Join(allow, ","))))
	}
}
