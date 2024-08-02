package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	customframework "github.com/Nishad4140/assignment_1/internal/customFramework"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	app := customframework.NewServer()

	app.Get("/users/{id}", func(req *customframework.Request, res *customframework.Response) {
		id := req.PathParam("id")
		user := User{
			ID:   id,
			Name: "Nishad",
			Age:  22,
		}
		res.Status(http.StatusOK).Json(user)
	})

	app.Post("/users", func(req *customframework.Request, res *customframework.Response) {
		var user User
		if err := req.Body(&user); err != nil {
			res.Status(http.StatusBadRequest).Json(map[string]string{"error": "Invalid request body"})
			return
		}

		res.Status(http.StatusCreated).Json(user)
	})

	app.Put("/users/{id}", func(req *customframework.Request, res *customframework.Response) {
		id := req.PathParam("id")
		var user User
		if err := req.Body(&user); err != nil {
			res.Status(http.StatusBadRequest).Json(map[string]string{"error": "Invalid request body"})
			return
		}
		user.ID = id 

		res.Status(http.StatusOK).Json(user)
	})

	app.Delete("/users/{id}", func(req *customframework.Request, res *customframework.Response) {
		id := req.PathParam("id")
		res.Status(http.StatusOK).Json(map[string]string{"message": fmt.Sprintf("User %s deleted", id)})
	})

	app.Any("/any", func(req *customframework.Request, res *customframework.Response) {
		queries := req.Query()
		headers := req.Headers()
		response := map[string]interface{}{
			"method":  req.Method,
			"queries": queries,
			"headers": headers,
		}
		res.Status(http.StatusOK).Json(response)
	})

	app.Get("/stream", func(req *customframework.Request, res *customframework.Response) {
		res.Header("Content-Type", "text/plain")
		for i := 0; i < 5; i++ {
			res.Write([]byte(fmt.Sprintf("Chunk %d\n", i)))
			time.Sleep(1 * time.Second)
		}
		res.End()
	})

	fmt.Println("Server starting on port 8080...")

	go func() {
		if err := app.Listen(8080); err != nil {
			fmt.Println("Server failed:", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	if err := app.Shutdown(15 * time.Second); err != nil {
		fmt.Println("Server shutdown failed:", err)
	} else {
		fmt.Println("Server gracefully stopped")
	}
}
