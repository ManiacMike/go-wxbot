package main

import (
	"flag"
	"fmt"
	"github.com/ManiacMike/go-wxbot/gowxbot"
)

var debug = flag.String("d", "off", "if on debug mode")

func main() {

	flag.Parse()

	fmt.Printf("debug mode %s\n", *debug)

	wx := gowxbot.Wxweb{}
	wx.Start()
}
