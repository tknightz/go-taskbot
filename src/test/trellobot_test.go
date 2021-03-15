package test

import (
	simpleJson "github.com/bitly/go-simplejson"
	trellobot "github.com/tknightz/taskbotgo/src/trellobot"
	"testing"
	"fmt"
)

func TestCreateCard(t *testing.T) {
	card := simpleJson.New()
	card.Set("name", "This is a test")
	shortLink, success := trellobot.CreateCard(card, "")
	if !success {
		t.Error(shortLink)
	}
}

func TestGetCardName(t *testing.T) {
	name := trellobot.GetCardName("jL07wW5X")
	fmt.Println(name)
}

func TestMoveCard(t *testing.T) {
	success := trellobot.MoveCard("jL07wW5X", "todo")
	if !success {
		t.Error("Failed move card")
	}
}

func TestPostComment(t *testing.T) {
	success := trellobot.PostComment("jL07wW5X", "This is comment")
	if !success {
		t.Error("Failed post comment!")
	}
}

func TestGetCustomField(t *testing.T) {
	value := trellobot.GetCustomField("jL07wW5X","mr")
	if value == "" {
		t.Fail()
	}
}

func TestSetCustomField(t *testing.T) {
	result := trellobot.SetCustomField("jL07wW5X", "mr", "test from go")
	if !result {
		t.Fail()
	}
}