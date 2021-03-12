package validator

import (
	"errors"
	"net/http"
	"strings"

	"github.com/thedevsaddam/govalidator"
	"github.com/txfs19260817/scopelens/server/models"
)

func TeamValidator(team *models.Team, r *http.Request) (error, bool) {

	rules := govalidator.MapData{
		"id":          []string{},
		"title":       []string{"required", "between:1,100"},
		"author":      []string{"between:1,80"},
		"format":      []string{"required", "in_formats"},
		"pokemon":     []string{"required", "between:1,6"},
		"showdown":    []string{"between:200,1800", "max_word:350"}, // 7 * 50 words
		"image":       []string{},
		"description": []string{"max:3000"},
		"uploader":    []string{"between:1,50"},
		"created_at":  []string{},
		"likes":       []string{"numeric"},
		"state":       []string{"numeric", "bool"},
	}

	// Showdown
	messages := govalidator.MapData{
		"showdown": []string{"between: Not a valid Showdown paste. "},
	}

	opts := govalidator.Options{
		Request:  r,
		Data:     team,
		Rules:    rules,
		Messages: messages, // custom message map (Optional)
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()

	// valid
	if len(e) == 0 {
		return nil, true
	}

	var errMsgs []string
	for _, i := range e {
		for _, j := range i {
			errMsgs = append(errMsgs, j)
		}
	}

	return errors.New(strings.Join(errMsgs, "\\n")), false
}
