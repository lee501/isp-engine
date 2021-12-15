package config

import "isp-engine/utils"

var (
	APIHost string
	APIPort string
)

func init() {
	APIHost = utils.Getenv("API_HOST", "127.0.0.1")
	APIPort = utils.Getenv("API_PORT", "8006")
}
