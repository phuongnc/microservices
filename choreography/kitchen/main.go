package main

import "kitchen-service/cmd"

func main() {
	runtime := cmd.NewRuntime()
	runtime.Serve()
}
