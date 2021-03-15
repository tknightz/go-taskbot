package gitlabbot

import (
	"regexp"
)

func GetCardId(description string) (string, bool) {
	cardIdRegex := regexp.MustCompile("https://trello.com/c/(\\w*)")
	cardId      := cardIdRegex.FindStringSubmatch(description)
	if len(cardId) > 0 {
		return cardId[1], true
	}
	return "", false
}
