package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

// Options for the CLI
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

type ReviewInput struct {
	Body struct {
		Author  string `json:"author" maxLength:"10" doc:"Author of the review"`
		Rating  int    `json:"rating" minimum:"1" maximum:"5" doc:"Rating from 1 to 5"`
		Message string `json:"message,omitempty" maxLength:"100" doc:"Review message"`
	}
}

func addRoutes(api huma.API) {
	// Register GET /greeting/{name}
	huma.Register(api, huma.Operation{
		OperationID: "get-greeting",
		Method:      http.MethodGet,
		Path:        "/greeting/{name}",
		Summary:     "Gret a greeting",
		Description: "Get a greeting for a person by name.",
		Tags:        []string{"Greetings"},
	}, func(ctx context.Context, input *struct {
		Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
	}) (*GreetingOutput, error) {
		resp := &GreetingOutput{}
		resp.Body.Message = fmt.Sprintf("Hello %s!", input.Name)
		return resp, nil
	})

	// Register POST /reviews handler.
	huma.Register(api, huma.Operation{
		OperationID:   "post-review",
		Method:        http.MethodPost,
		Path:          "/reviews",
		Summary:       "Gret a greeting",
		Description:   "Post a review.",
		Tags:          []string{"Reviews"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, i *ReviewInput) (*struct{}, error) {
		// TODO: save review in data store
		return nil, nil
	})
}

func loadExecs() {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table foo (id integer not null primary key, name text);
	delete from foo;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("こんにちは世界%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
}

func main() {
	// Create a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Create a new router & API
		router := chi.NewMux()
		api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

		addRoutes(api)

		hooks.OnStart(func() {
			fmt.Printf("Starting server on port %d...\n", options.Port)
			http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}
