package trellobot

import (
	simpleJson "github.com/bitly/go-simplejson"
	util "github.com/tknightz/taskbotgo/src/util"
)

// SendRequest interface to send a request with json body.
func SendRequest(method string, endpoint string, params *simpleJson.Json) *simpleJson.Json {
	params.Set("key", KEY)
	params.Set("token", TOKEN)
	response := util.SendRequest(method, endpoint, params, nil)
	return response
}

func CreateCard(card *simpleJson.Json, typeCard string) (string, bool) {
	endpoint := "https://api.trello.com/1/cards"
	params := card
	params.Set("idList", TYPELIST_MAP[typeCard])
	params.Set("desc", DESC)

	if typeCard != "" {
		params.Set("idLabels", LABELS_MAP[typeCard])
	}

	response := SendRequest("POST", endpoint, params)
	shortLink, err := response.Get("shortLink").String()
	if err != nil {
		return "", false
	}
	return shortLink, true
}

func GetCardName(cardId string) string {
	endpoint := "https://api.trello.com/1/cards/" + cardId + "/name"
	response := SendRequest("GET", endpoint, simpleJson.New())
	name, _ := response.Get("_value").String()
	return name
}

func GetCustomField(cardId string, customFieldName string) string {
	endpoint := "https://api.trello.com/1/cards/" + cardId + "/customFieldItems"
	response := SendRequest("GET", endpoint, simpleJson.New())

	if response == nil {
		return ""
	}

	customFieldItems := response.MustArray()

	for index := range customFieldItems {
		idField, _ := response.GetIndex(index).Get("idCustomField").String()
		if idField == CUSTOMFIELDS_MAP[customFieldName] {
			textValue, _ := response.GetIndex(index).GetPath("value", "text").String()
			return textValue
		}
	}
	return ""
}

func SetCustomField(cardId string, customFieldName string, value string) bool {
	oldValue := GetCustomField(cardId, customFieldName)
	if oldValue == value {
		return true
	}

	endpoint := "https://api.trello.com/1/cards/" + cardId + "/customField/" + CUSTOMFIELDS_MAP[customFieldName] + "/item"
	params := simpleJson.New()
	params.SetPath([]string{"value", "text"}, value)

	response := SendRequest("PUT", endpoint, params)
	return response != nil
}

func MoveCard(cardId string, listName string) bool {
	endpoint := "https://api.trello.com/1/cards/" + cardId
	params := simpleJson.New()
	params.Set("idList", LIST_MAP[listName])

	response := SendRequest("PUT", endpoint, params)
	_, success := response.CheckGet("shortLink")
	return success
}

func PostComment(cardId string, comment string) bool {
	endpoint := "https://api.trello.com/1/cards/" + cardId + "/actions/comments"
	params := simpleJson.New()
	params.Set("text", comment)

	response := SendRequest("POST", endpoint, params)
	_, err := response.Get("data").Get("card").Get("shortLink").String()

	if err != nil {
		return false
	}
	return true
}

func CreateWebhook(idModel string) bool {
	params := simpleJson.New()
	params.Set("idModel", idModel)
	params.Set("callbackURL", CALLBACK_URL)

	endpoint := "https://api.trello.com/1/webhooks"
	response := SendRequest("POST", endpoint, params)
	return response != nil
}
