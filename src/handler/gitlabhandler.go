package handler

import (
	"fmt"
	simpleJson "github.com/bitly/go-simplejson"
	trellobot "github.com/tknightz/taskbotgo/src/trellobot"
	gitlabbot "github.com/tknightz/taskbotgo/src/gitlabbot"
	slackbot "github.com/tknightz/taskbotgo/src/slackbot"
	"github.com/valyala/fasthttp"
)

func ProcessGitlabOpenedMerge(cardId string, mergeUrl string, isWorkInProgress bool) bool {
	if !isWorkInProgress {
		trellobot.MoveCard(cardId, "qc")
		trellobot.SetCustomField(cardId, "mr", mergeUrl)
		return true
	} else {
		trellobot.MoveCard(cardId, "edit")
		return true
	}
}

func ProcessGitlabMerged(cardId string) bool {
	trellobot.MoveCard(cardId, "done")
	cardName := trellobot.GetCardName(cardId)
	content := fmt.Sprintf("Task \"**%v**\" đã được làm xong!", cardName)
	return slackbot.ReplyInThread("channel", "thread", content)
}

func GitlabHandler(ctx *fasthttp.RequestCtx) {
	data, _             := simpleJson.NewJson(ctx.Request.Body())
	objectKind, _       := data.Get("object_kind").String()
	state, _            := data.GetPath("object_attributes", "state").String()
	isWorkInProgress, _ := data.GetPath("object_attributes", "work_in_progress").Bool()
	description, _      := data.GetPath("object_attributes", "description").String()
	urlMergeRequest, _  := data.GetPath("object_attributes", "url").String()
	cardId, exists      := gitlabbot.GetCardId(description)

	if !exists {
		fmt.Fprintf(ctx, "errr")
	}

	if objectKind == "merge_request" && state == "opened" {
		ProcessGitlabOpenedMerge(cardId, urlMergeRequest, isWorkInProgress)
	} else if objectKind == "merge_request" && state == "merged" {
		ProcessGitlabMerged(cardId)
	}
}
