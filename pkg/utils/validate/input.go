package validate

import (
	"errors"
	"regexp"
)

func IsWikipediaUrl(  url string) error{
	matched, _ := regexp.MatchString(`^https?:\/\/[a-z]{2}\.wikipedia\.org\/wiki\/`, url)
	if !matched {
		return errors.New("The input should be a Wikipedia page")
	}
	return nil
}