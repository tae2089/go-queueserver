package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var (
	MaxWorker       = 3  //os.Getenv("MAX_WORKERS")
	MaxQueue        = 20 //os.Getenv("MAX_QUEUE")
	MaxLength int64 = 2048
)
var JobQueue chan Job
var data []int

func main() {
	JobQueue = make(chan Job, MaxQueue)
	log.Println("main start")
	dispatcher := NewDispatcher(MaxWorker)
	dispatcher.Run()
	data = make([]int, 0)
	r := gin.Default()
	r.POST("/", func(ctx *gin.Context) {
		fmt.Println("A valid payload request received")
		form, err := ctx.MultipartForm()
		if err != nil {
			panic(err)
		}
		files := form.File["files"]
		for _, file := range files {
			work := Job{
				Payload: Payload{
					data: file.Filename,
				},
			}
			fmt.Println("sending payload  to workque")
			JobQueue <- work
			fmt.Println("sent payload  to workque")
		}
		ctx.String(200, "finish")
	})
	r.GET("/", func(ctx *gin.Context) {
		log.Println(data)
		ctx.String(http.StatusOK, "good")
	})
	r.Run(":3000")
}
