package api

import "demo/internal/rest"

func InitApi() {
	InitIoc()

	rest.InitUserRest()
}

func InitIoc() {
	rest.InitUserIoc()
}
