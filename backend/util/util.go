package util

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

var validTitle = regexp.MustCompile(`(?i)\b[a-zA-Z]{3,}\b`)

func ValidateTitle(title string) error {
	title = strings.TrimSpace(title)
	if !validTitle.MatchString(title) {
		return errors.New("title must contain at least one word with more than 2 letters and cannot consist of only spaces or punctuation")
	}
	return nil
}

var adjectives = []string{"Swift", "Mighty", "Brave", "Clever", "Silent", "Fierce", "Bright", "Shadow", "Lunar", "Wild"}
var nouns = []string{"Tiger", "Eagle", "Wizard", "Knight", "Panther", "Phoenix", "Dragon", "Wolf", "Hawk", "Fox"}

func GenerateNickname() string {
	adjective := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	number := rand.Intn(1000)

	return fmt.Sprintf("%s%s%d", adjective, noun, number)
}
