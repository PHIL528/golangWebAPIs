package main

import (
	"os/exec"
	"time"
)

var terminal string = "Terminal"
var script1 string = `tell application "` + terminal + `" to do script "cd go/src/proto-playground ; make pubsub"`
var script2 string = `tell application "` + terminal + `" to do script "cd go/src/proto-playground ; make s"`
var script3 string = `tell application "` + terminal + `" to do script "cd go/src/proto-playground ; make l"`
var script4 string = `tell application "` + terminal + `" to do script "cd go/src/proto-playground ; make c"`

func main() {
	cmd := exec.Command("osascript", "-s", "h", "-e", script1)
	cmd.Run()
	<-time.After(1)
	cmd = exec.Command("osascript", "-s", "h", "-e", script2)
	cmd.Run()
	<-time.After(1)
	cmd = exec.Command("osascript", "-s", "h", "-e", script3)
	cmd.Run()
	cmd = exec.Command("osascript", "-s", "h", "-e", script4)
	cmd.Run()
}
