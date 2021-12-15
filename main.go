package main

import (
	"fmt"
	"isp-engine/config"
	"isp-engine/server"
)

func main() {
	apiListener := fmt.Sprintf("%s:%s", config.APIHost, config.APIPort)
	serve := server.NewServer(apiListener)
	serve.Run()
}
