package main

import "order-service/cmd"

func main() {
	runtime := cmd.NewRuntime()
	runtime.Serve()
}
