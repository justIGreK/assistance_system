package util

import (
	"errors"
	"regexp"
	"strings"
)

var validTitle = regexp.MustCompile(`(?i)^[\p{L}0-9]{3,}.*$`)

func ValidateTitle(title string) error {
	title = strings.TrimSpace(title)
	if !validTitle.MatchString(title) {
		return errors.New("title must contain at least one word with more than 2 letters and cannot consist of only spaces or punctuation")
	}
	return nil
}
