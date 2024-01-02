package main

import "payment-service/cmd"

func main() {
	runtime := cmd.NewRuntime()
	runtime.Serve()
}
