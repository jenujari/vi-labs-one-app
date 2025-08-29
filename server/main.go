package main

import (
	"fmt"
	"os"
	"server/config"
	"server/helpers"
	"server/repository"

	"server/router"

	"github.com/goforj/godump"
)

var masterProcess *helpers.ProcessContext

func init() {
	helpers.InitProcessContext()
}

func main() {
	masterProcess = helpers.GetMainProcess()

	config.SetuDbConnection(masterProcess.CTX)

	defer func() {
		config.CloseDbConnection()
	}()

	repository.InitRepository()

	masterProcess.AddWorker(1)
	srv := router.GetServer()

	go router.RunServer(masterProcess)
	config.GetLogger().Println("Server is running at ", srv.Addr)

	masterProcess.WaitForFinish()
}

type Profile struct {
	Age   int
	Email string
}

type User struct {
	Name    string
	Profile Profile
}

func dump() {
	user := User{
		Name: "Alice",
		Profile: Profile{
			Age:   30,
			Email: "alice@example.com",
		},
	}

	// Pretty-print to stdout
	godump.Dump(user)

	// Get dump as string
	output := godump.DumpStr(user)
	fmt.Println("str", output)

	// HTML for web UI output
	html := godump.DumpHTML(user)
	fmt.Println("html", html)

	// Print JSON directly to stdout
	godump.DumpJSON(user)

	// Write to any io.Writer (e.g. file, buffer, logger)
	godump.Fdump(os.Stderr, user)

	// Dump and exit
	godump.Dd(user) // this will print the dump and exit the program
}
