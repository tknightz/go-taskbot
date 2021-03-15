package test

import (
	slackbot "github.com/tknightz/taskbotgo/src/slackbot"
	"testing"
	"fmt"
)

func TestReplyInThread(t *testing.T) {
	result := slackbot.ReplyInThread("C01CU5JCK4L","1609126591.000900", "Test in go")
	fmt.Println(result)
}

func TestPostReaction(t *testing.T) {
	success := slackbot.PostReaction("C01CU5JCK4L","1609126591.000900", "white_check_mark")
	if !success {
		t.Fail()
	}
}

func TestGetContentThread(t *testing.T) {
	result := slackbot.GetContentThread("C01CU5JCK4L","1609126591.000900")
	fmt.Println(result)
}

func TestRenderComment(t *testing.T) {
	result := slackbot.RenderComment("C01CU5JCK4L","1609126591.000900")
	fmt.Println(result)
}