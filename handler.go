package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Orginal struct {
	Url string `form:"url" binding:"required"`
}

type Shorter struct {
	Url string `uri:"url" binding:"required"`
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

func ShortUrl(orginal string) (string, error) {
	orginal = checkPrefix(orginal)

	h := sha256.New()
	h.Write([]byte(orginal))

	bs := h.Sum(nil)
	newUrl := base64.StdEncoding.EncodeToString(bs)[:8]

	err := rdc.AddUrl(newUrl, orginal)
	if err != nil {
		return "", err
	}

	return newUrl, nil
}

func AddShort(c *gin.Context) {
	var orginal Orginal

	err := c.ShouldBind(&orginal)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Use /add/?url=",
		})
		return
	}

	newUrl, err := ShortUrl(orginal.Url)
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
	var shorter Shorter

	err := c.ShouldBindUri(&shorter)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Use /r/<code>",
		})
		return
	}

	orginalUrl, err := rdc.FindUrl(shorter.Url)
	if err != nil {
		c.JSON(404, gin.H{
			"error": "Url don't exists",
		})
		return
	}

	c.Redirect(301, orginalUrl)
}

func RunServer() {
	r := gin.Default()

	r.GET("/add", AddShort)
	r.GET("/r/:url", RedirectUrl)

	r.Run()
}
