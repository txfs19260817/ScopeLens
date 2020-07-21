package validator

import (
	"errors"
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"scopelens-server/models"
	"strconv"
	"strings"
)

var (
	formats []string
)

// custom rules
func init() {
	// load formats and pokemon names
	formats = models.GetFormats()

	// custom rules to take fixed length word.
	// e.g: max_word:5 will throw error if the field contains more than 5 words
	govalidator.AddCustomRule("max_word", func(field string, rule string, message string, value interface{}) error {
		valSlice := strings.Fields(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "max_word:")) //handle other error
		if len(valSlice) > l {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("The %s field must not be greater than %d words. ", field, l)
		}
		return nil
	})

	// custom rules to check if an element in a slice
	govalidator.AddCustomRule("in_formats", func(field string, rule string, message string, value interface{}) error {
		fname, ok := value.(string)
		if !ok {
			// wrong use case
			return fmt.Errorf("Incorrect Format type. ")
		}
		for _, f := range formats {
			if f == fname {
				return nil
			}
		}
		if message != "" {
			return errors.New(message)
		}
		return fmt.Errorf("The format %s is temporarily not supported. ", fname)
	})

}