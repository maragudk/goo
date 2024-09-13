package main

import (
	"maragu.dev/goo/service"
)

func main() {
	service.Start(service.Options{
		Log:     service.NewLogger(),
		Migrate: true,
	})
}
