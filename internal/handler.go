package internal

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Controller struct {
	q *MessageQueueMap
}

func NewController(q *MessageQueueMap) *Controller {
	return &Controller{q: q}
}

func (c *Controller) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// Лень писать нормальный роутер
	fmt.Printf("get new %s request\n", req.Method)

	urlParts := strings.Split(req.URL.Path, "/")
	queue := c.q.GetQueue(urlParts[len(urlParts)-1])
	switch req.Method {
	case "PUT":
		c.PutToQueue(resp, req, queue)
	case "GET":
		c.GetFromQueue(resp, req, queue)
	default:
		resp.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *Controller) PutToQueue(resp http.ResponseWriter, req *http.Request, queue *MessageQueue) {
	item := req.URL.Query().Get("v")
	if item == "" {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	queue.Put(item)
	resp.WriteHeader(http.StatusOK)
}

func (c *Controller) GetFromQueue(resp http.ResponseWriter, req *http.Request, queue *MessageQueue) {
	timeout, _ := strconv.Atoi(req.URL.Query().Get("timeout"))
	var message string
	if timeout != 0 {
		ctx, cancel := context.WithTimeout(req.Context(), time.Second*time.Duration(timeout))
		defer cancel()

		message = queue.GetCtx(ctx)
	} else {
		message = queue.GetOrReturn()
	}

	if message == "" {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(message))
}
