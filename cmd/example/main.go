package main

import (
	"maragu.dev/goo/service"
)

func main() {
	service.Start(service.Options{
		Migrate: true,
	})
}
