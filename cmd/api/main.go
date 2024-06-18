package main

import (
	"ISO_Auditing_Tool/cmd/server"
	"fmt"
)

func main() {

	// server := server.NewServer()
	//
	// err := server.ListenAndServe()
	srv := server.NewServer()
	httpServer := srv.Start()
	err := httpServer.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
