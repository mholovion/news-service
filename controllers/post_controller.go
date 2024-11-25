package controllers

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mholovion/news-service/config"
	"github.com/mholovion/news-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	funcMap = template.FuncMap{
		"sub": func(a, b int) int { return a - b },
		"add": func(a, b int) int { return a + b },
	}

	validate = validator.New()
)

const pageSize = 5

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplPath := fmt.Sprintf("views/%s.html", tmpl)
	templates, err := template.New("").Funcs(funcMap).ParseFiles("views/layout.html", tmplPath)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = templates.ExecuteTemplate(&buf, "layout", data)
	if err != nil {
		log.Printf("Error executing template %s: %v", tmpl, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		data := struct {
			Error string
			Post  models.Post
		}{}
		renderTemplate(w, "new", data)
	case "POST":
		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))

		input := models.Post{
			Title:   title,
			Content: content,
		}

		err := validate.Struct(input)
		if err != nil {
			data := struct {
				Error string
				Post  models.Post
			}{
				Error: err.Error(),
				Post:  input,
			}
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate(w, "new", data)
			return
		}

		collection := config.GetCollection("posts")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		input.CreatedAt = time.Now()
		input.UpdatedAt = time.Now()
		_, err = collection.InsertOne(ctx, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func ReadPosts(w http.ResponseWriter, r *http.Request) {
	collection := config.GetCollection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	filter := bson.M{}
	if search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"title": bson.M{"$regex": search, "$options": "i"}},
				{"content": bson.M{"$regex": search, "$options": "i"}},
			},
		}
	}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	skip := (int64(page) - 1) * pageSize

	cursor, err := collection.Find(ctx, filter, options.Find().SetSkip(skip).SetLimit(pageSize).SetSort(bson.D{{"created_at", -1}}))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var posts []models.Post
	if err = cursor.All(ctx, &posts); err != nil {
		log.Printf("Error decoding posts: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Posts retrieved: %+v", posts)

	data := struct {
		Posts       []models.Post
		CurrentPage int
		TotalPages  int
		Search      string
	}{
		Posts:       posts,
		CurrentPage: page,
		TotalPages:  int((total + pageSize - 1) / pageSize),
		Search:      search,
	}

	renderTemplate(w, "index", data)
}

func ReadPost(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/post/")
	if id == "" {
		http.NotFound(w, r)
		log.Println("ID is empty")
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.NotFound(w, r)
		log.Println("Can't convert ID from Hex, ID:", id)
		return
	}

	collection := config.GetCollection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var post models.Post
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
	if err != nil {
		http.NotFound(w, r)
		log.Println("Can't find post in collection, objectID:", objectID)
		return
	}

	data := struct {
		Post models.Post
	}{
		Post: post,
	}

	renderTemplate(w, "show", data)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/post/"), "/edit")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	collection := config.GetCollection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch r.Method {
	case "GET":
		var post models.Post
		err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		data := struct {
			Error string
			Post  models.Post
		}{
			Post: post,
		}

		renderTemplate(w, "edit", data)
	case "POST":
		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))

		input := models.Post{
			ID:      objectID,
			Title:   title,
			Content: content,
		}

		err := validate.Struct(input)
		if err != nil {
			data := struct {
				Error string
				Post  models.Post
			}{
				Error: err.Error(),
				Post:  input,
			}
			renderTemplate(w, "edit", data)
			return
		}

		update := bson.M{
			"$set": bson.M{
				"title":      title,
				"content":    content,
				"updated_at": time.Now(),
			},
		}
		_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/post/"+id, http.StatusSeeOther)
	}
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/post/"), "/delete")
	if id == "" {
		http.NotFound(w, r)
		return
	}
	log.Println("ID:", id)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	collection := config.GetCollection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
