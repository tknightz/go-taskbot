package main

import (
	"fmt"
	"log"
	"os"
	handler "github.com/tknightz/taskbotgo/src/handler"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)


// Index : home page
func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "App is running!\n")
}

// GetPort : get port in os environment
func GetPort() string {
	port := os.Getenv("PORT")
	if port != "" {
		return ":" + port
	}
	return ":5000"
}

func main (){
	router := fasthttprouter.New()
	router.GET("/", Index)
	router.POST("/gitlabbot", handler.GitlabHandler)
	router.POST("/slackbot", handler.SlackHandler)
	port := GetPort()
	log.Fatal(fasthttp.ListenAndServe(port, router.Handler))
}
