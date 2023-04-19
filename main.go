package main

import "Im/router"

func main() {
	r := router.Router()
	r.Run(":8080")
}
