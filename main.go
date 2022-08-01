package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

type RootResolver struct{}

var (
	opts     = []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	articles = []Article{
		{
			ID:          "0",
			Title:       "Test book title",
			Description: "Test book description",
			Content:     "test book content",
		},
		{
			ID:          "1",
			Title:       "Test book title 2",
			Description: "Test book description 2",
			Content:     "test book content 2",
		},
	}
)

type Article struct {
	ID          graphql.ID
	Title       string
	Description string
	Content     string
}

//Index
func (r *RootResolver) Feed() ([]Article, error) {
	return articles, nil
}

//Find by ID
func (r *RootResolver) Find(id struct{ ID graphql.ID }) (Article, error) {
	for _, article := range articles {
		if article.ID == id.ID {
			return article, nil
		}

	}

	//Error
	return articles[0], errors.New("No article found with ID " + string(id.ID))
}

//Create Article
func (r *RootResolver) Post(args struct {
	ID          graphql.ID
	Title       string
	Description string
	Content     string
}) (Article, error) {

	newArticle := Article{
		ID:          graphql.ID(fmt.Sprint(len(articles))),
		Title:       args.Title,
		Description: args.Description,
		Content:     args.Content,
	}

	articles = append(articles, newArticle)
	return newArticle, nil
}

//Update Article

func (r *RootResolver) Update(args struct {
	ID          graphql.ID
	Title       string
	Description string
	Content     string
},
) (Article, error) {

	isUpdated := false

	updatedArticle := Article{
		ID:          args.ID,
		Title:       args.Title,
		Description: args.Description,
		Content:     args.Content,
	}

	for index, article := range articles {
		if article.ID == args.ID {
			articles[index] = updatedArticle
			isUpdated = true
		}
	}

	if isUpdated {
		return updatedArticle, nil
	} else {
		return updatedArticle, errors.New("No article found with ID" + string(args.ID))
	}
}

//DELETE

func (r *RootResolver) Delete(id struct{ ID graphql.ID }) ([]Article, error) {
	deleted := false

	for index, article := range articles {
		if article.ID == id.ID {
			articles = append(articles[:index], articles[index+1:]...)
			deleted = true
		}
	}

	if deleted {
		return articles, nil
	} else {
		return articles, errors.New("No article with the ID: " + string(id.ID))
	}
}

func (r *RootResolver) Info() (string, error) {
	return "this is a thing", nil
}

//CORS - Cross Origin Resource Sharing

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func parseSchema(path string, resolver interface{}) *graphql.Schema {
	bstr, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	schemaString := string(bstr)
	parsedSchema, err := graphql.ParseSchema(
		schemaString,
		resolver,
		opts...,
	)
	if err != nil {
		panic(err)
	}
	return parsedSchema
}


func main() {
	http.Handle("/graphql", Cors(&relay.Handler{
		Schema: parseSchema("./schema.graphql", &RootResolver{}),
	}))

	fmt.Println("serving on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
