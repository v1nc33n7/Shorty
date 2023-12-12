package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/gin-gonic/gin"
)

type KeyValue struct {
	Key   string `uri:"url"`
	Value string `form:"url"`
}

func checkPrefix(url string) string {
	if len(url) < 7 {
		return "https://" + url
	}

	if url[:7] != "http://" || url[:8] != "https://" {
		return "https://" + url
	}
	return url
}

func ShortUrl(msg KeyValue) (string, error) {
	msg.Value = checkPrefix(msg.Value)

	h := sha256.New()
	h.Write([]byte(msg.Value))

	bs := h.Sum(nil)
	newUrl := base64.StdEncoding.EncodeToString(bs)[:8]

	msg.Key = newUrl

	go func() {
		pq.pipe <- msg
		rdc.pipe <- msg
	}()

	return newUrl, nil
}

func AddShort(c *gin.Context) {
	var keyvalue KeyValue

	err := c.ShouldBind(&keyvalue)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Use /add/?url=",
		})
		return
	}

	newUrl, err := ShortUrl(keyvalue)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Internal error",
		})
		return
	}

	c.JSON(200, gin.H{
		"url": fmt.Sprintf("http://%s/r/%s", c.Request.Host, newUrl),
	})
}

func RedirectUrl(c *gin.Context) {
	var keyvalue KeyValue

	err := c.ShouldBindUri(&keyvalue)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Use /r/<code>",
		})
		return
	}

	err = rdc.FindCacheUrl(&keyvalue)
	if err != nil {
		err = pq.FindUrl(&keyvalue)
		if err != nil {
			c.JSON(404, gin.H{
				"error": "Url don't exists",
			})
			return
		}
	}

	c.Redirect(301, keyvalue.Value)
}

func RunServer() {
	r := gin.Default()

	r.GET("/add", AddShort)
	r.GET("/r/:url", RedirectUrl)

	r.Run()
}
