package main

import (
	"github.com/catbugdemo/gee"
)

func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(200, "<h1>Hello Gee</h1>")
	})

	r.Run(":9999")
}
