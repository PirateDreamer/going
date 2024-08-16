package api

import "github.com/PirateDreamer/going/example/internal/rest"

func InitApi() {
	InitIoc()

	rest.InitUserRest()
}

func InitIoc() {
	rest.InitUserIoc()
}
