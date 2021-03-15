package slackbot

import (
	"fmt"
	"regexp"
	"strings"

	simpleJson "github.com/bitly/go-simplejson"
	util "github.com/tknightz/taskbotgo/src/util"
)

func SendRequest(method string, endpoint string, params *simpleJson.Json) *simpleJson.Json {
	var bearer = "Bearer " + TOKEN
	var header = map[string]string{
		"Authorization": bearer,
	}
	params.Set("token", TOKEN)
	response := util.SendRequest(method, endpoint, params, header)
	return response
}

func GetUserName(userId string) string {
	endpoint := "https://slack.com/api/users.info?token=" + TOKEN + "&user=" + userId
	params := simpleJson.New()
	response := util.SendRequest("GET", endpoint, params, nil)
	userName, err := response.GetPath("user", "real_name").String()

	if err != nil {
		return "UserId_NOT_FOUND"
	}

	return userName
}

func GetLinkToThread(channelId string, threadTs string) string {
	endpoint := "https://slack.com/api/chat.getPermalink?token=" + TOKEN + "&channel=" + channelId + "&message_ts=" + threadTs
	params := simpleJson.New()
	response := util.SendRequest("GET", endpoint, params, nil)
	linkToThread, err := response.Get("permalink").String()

	if err != nil {
		return "ERROR_GET_THREAD_LINK"
	}

	return linkToThread
}

func GetMessageCreatedCard(userId string, cardId string) string {
	message := fmt.Sprintf(
		`"<@%v>\n:ok_hand: Tạo card thành công!\n"
   ":card_index: Id Card: %v\n"
   ":link: Link : https://trello.com/c/%v\n"`,
		userId, cardId, cardId)
	return message
}

func ReplyInThread(channelId string, threadTs string, content string) bool {
	endpoint := "https://slack.com/api/chat.postMessage"
	params := simpleJson.New()
	params.Set("channel", channelId)
	params.Set("thread_ts", threadTs)
	params.Set("text", content)

	response := SendRequest("POST", endpoint, params)

	return response != nil
}

func PostReaction(channelId string, timestamp string, reactionName string) bool {
	endpoint := "https://slack.com/api/reactions.add"
	params := simpleJson.New()
	params.Set("channel", channelId)
	params.Set("timestamp", timestamp)
	params.Set("name", reactionName)

	response := SendRequest("POST", endpoint, params)
	success, _ := response.Get("ok").Bool()
	return success
}

func GetContentThread(channelId string, threadTs string) *simpleJson.Json {
	endpoint := "https://slack.com/api/conversations.replies?token=" + TOKEN + "&channel=" + channelId + "&ts=" + threadTs
	params := simpleJson.New()
	response := util.SendRequest("GET", endpoint, params, nil)
	bytes, err := response.Get("messages").Encode()

	if err != nil {
		fmt.Println("ERROR_GET_MESSAGES")
		return simpleJson.New()
	}
	messages, _ := simpleJson.NewJson(bytes)
	return messages
}

func GetTaggedUsers(message string) map[string]string {
	userRegex := regexp.MustCompile("<@(?P<userId>\\w*)>")
	users := userRegex.FindAllStringSubmatch(message, -1)

	if len(users) == 0 {
		return map[string]string{}
	}

	var userIdsMap = map[string]string{}
	for _, userId := range users {
		userName := GetUserName(userId[1])
		userIdsMap[userId[1]] = userName
	}

	return userIdsMap
}

func GetTypeCard(message string) (string, bool) {
	typeRegex := regexp.MustCompile("type=(?P<typeCard>\\d)")
	typeCard := typeRegex.FindStringSubmatch(message)
	if len(typeCard) > 1 {
		return typeCard[1], true
	}
	return "", false
}

func GetCardName(message string) (string, bool) {
	cardNameRegex := regexp.MustCompile("name\\s\"(?P<cardName>\\w.*)\"")
	cardName := cardNameRegex.FindStringSubmatch(message)
	if len(cardName) > 0 {
		return cardName[1], true
	}
	return "", false
}

func GetRequestor(messages *simpleJson.Json) (string, string) {
	for index := range messages.MustArray() {
		_, isBot := messages.GetIndex(index).CheckGet("bot_id")
		if !isBot {
			requestorId, _ := messages.GetIndex(index).Get("user").String()
			requestorName := GetUserName(requestorId)
			return requestorId, requestorName
		}
	}
	return "", ""
}

func ProcessMesssages(messages *simpleJson.Json) string {
	var output string
	for index := range messages.MustArray() {
		sender, _ := messages.GetIndex(index).Get("user").String()
		text, _ := messages.GetIndex(index).Get("text").String()
		message := fmt.Sprintf("<@%v> : %v", sender, text)
		var taggedUsers map[string]string = GetTaggedUsers(message)

		processedMessage := message

		if taggedUsers != nil {
			for key, val := range taggedUsers {
				processedMessage = strings.ReplaceAll(processedMessage, "<@"+key+">", "*@"+val+"*")
			}
		}
		output += processedMessage + "\n"
	}
	return output
}

func RenderComment(channelId string, threadTs string) string {
	messages := GetContentThread(channelId, threadTs)

	channelLinkToThread := make(chan string)
	channelRequestorName := make(chan string)
	channelCommentBody := make(chan string)

	go func() {
		link := GetLinkToThread(channelId, threadTs)
		channelLinkToThread <- link
	}()
	go func() {
		_, name := GetRequestor(messages)
		channelRequestorName <- name
	}()
	go func() {
		channelCommentBody <- ProcessMesssages(messages)
	}()

	linkToThread := <-channelLinkToThread
	requestorName := <-channelRequestorName
	commentBody := <-channelCommentBody

	output := fmt.Sprintf("- Người yêu cầu : %v\n- Link to thread : %v\n----\n%v", requestorName, linkToThread, commentBody)
	return output
}
