package handler

import (
	"fmt"
	simpleJson "github.com/bitly/go-simplejson"
	trellobot "github.com/tknightz/taskbotgo/src/trellobot"
	slackbot "github.com/tknightz/taskbotgo/src/slackbot"
	util "github.com/tknightz/taskbotgo/src/util"
	"github.com/valyala/fasthttp"
	"regexp"
)

type SlackInfo struct {
	message string
	channelId string
	threadTs string
	user string
}

func ExtractCard(message string) *simpleJson.Json {
	card        := simpleJson.New()
	cardName, _ := slackbot.GetCardName(message)
	members     := GetAssignedUser(message)
	card.Set("name", cardName)
	card.Set("idMembers", members)
	return card
}

func GetAssignedUser(message string) []string {
	assignRegex := regexp.MustCompile("assign\\s(<@\\w*> ?)+")
	assignMatch := assignRegex.FindStringSubmatch(message)
	var assignString string

	if len(assignMatch) > 0 {
		assignString = assignMatch[0]
	}

	assignedMembersMap := slackbot.GetTaggedUsers(assignString)
	var users []string

	for _, user := range(assignedMembersMap) {
		users = append(users, user)
	}
	return users
}

func CreateNewTask(message string) (string, bool) {
	typeCard, _        := slackbot.GetTypeCard(message)
	card               := ExtractCard(message)
	shortLink, success := trellobot.CreateCard(card, typeCard)
	return shortLink, success
}

func DetectKeyword(message string) string {
	taskRegex := regexp.MustCompile("^<@\\w*>\\s(\\w*)")
	keyword   := taskRegex.FindStringSubmatch(message)

	if len(keyword) < 1 {
		return ""
	}

	return keyword[1]
}

func ExtractInfo(data *simpleJson.Json) SlackInfo {
	message, _    := data.GetPath("event", "text").String()
	channelId, _  := data.GetPath("event", "channel").String()
	user, _       := data.GetPath("event", "user").String()
	threadTs, err := data.GetPath("event", "thread_ts").String()

	if err != nil {
		ts, _ := data.GetPath("event", "ts").String()
		threadTs = ts
	}

	message  = util.Standardlize(message)
	info    := SlackInfo{
		message: message,
		channelId: channelId,
		threadTs: threadTs,
		user: user,
	}
	return info
}

func SlackHandler(ctx *fasthttp.RequestCtx) {
	data, _      := simpleJson.NewJson(ctx.Request.Body())
	typeReq, err := data.Get("type").String()

	if err == nil && typeReq == "url_verification" {
		challenge, _ := data.Get("challenge").String()
		fmt.Fprintf(ctx, challenge)
		return
	}

	slackInfo := ExtractInfo(data)
	keyword   := DetectKeyword(slackInfo.message)

	switch keyword {
	case "": 
		slackbot.ReplyInThread(slackInfo.channelId, slackInfo.threadTs, "Wrong")
		fmt.Fprintf(ctx, "Success")
		return

	case "name": 
		cardId, success := CreateNewTask(slackInfo.message)
		if !success {
			fmt.Fprintf(ctx, cardId)
			return
		}
		message := slackbot.GetMessageCreatedCard(slackInfo.user, cardId)
		comment := slackbot.RenderComment(slackInfo.channelId, slackInfo.threadTs)
		go trellobot.PostComment(cardId, comment)
		go slackbot.PostReaction(slackInfo.channelId, slackInfo.threadTs, "heavy_check_mark")
		go slackbot.ReplyInThread(slackInfo.channelId, slackInfo.threadTs, message)

		fmt.Fprintf(ctx, "Success")
		return 

	case "update":
		slackbot.ReplyInThread(slackInfo.channelId, slackInfo.threadTs, "Underworking")
		fmt.Fprintf(ctx, "Underworking")
		return

	case "status":
		slackbot.ReplyInThread(slackInfo.channelId, slackInfo.threadTs, "Underworking")
		fmt.Fprintf(ctx, "Underworking")
		return

	default:
		slackbot.ReplyInThread(slackInfo.channelId, slackInfo.threadTs, "Unknown command")
		fmt.Fprintf(ctx, "Unknown command")
		return
	}
}
