package main

import "github.com/thatisuday/commando"

func main() {
	commando.
		SetExecutableName("reactor").
		SetVersion("v1.0.0").
		SetDescription("This CLI tool helps you create and manage React projects.")
}
