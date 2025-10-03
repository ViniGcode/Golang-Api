package repository

import (
	"context"
	"errors"
	"gin-rest-api/api"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectionURL = "mongodb://localhost:xxxx" //porta do banco
	poolSize      = 100
)

//BookRepository repositorio para interagir com o mongo
type BookRepository struct{}

var db *mongo.Database

// Estabelece a conexão com o banco de dados e cria o handle do banco ao carregar
// para todas as interações subsequentes com o banco de dados
func init() {
	clientOptions := options.Client().ApplyURI(connectionURL).SetMaxPoolSize(poolSize)
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	db = client.Database("library")
}

//List lista todos livros
func (br *BookRepository) List() []api.Book {
	log.Println("[BookRepository] List() Obtendo livros do Mongo.")
	collection := db.Collection("books")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		log.Fatal("Erro ao recuperar livros..")
	}
	defer cursor.Close(ctx)

	var books []api.Book
	for cursor.Next(ctx) {
		var book api.Book
		cursor.Decode(&book)
		books = append(books, book)
	}
	if err := cursor.Err(); err != nil {
		log.Println("Erro ao percorrer o cursor")
		//return error
	}
	return books
}

//Get Book by ID
func (br *BookRepository) Get(id string) api.Book {
	col := db.Collection("books")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)

	var b api.Book
	col.FindOne(ctx, bson.M{"_id": objID}).Decode(&b)
	return b
}

//Create cria um novo documento de book no Mongo
func (br *BookRepository) Create(book api.Book) interface{} {
	log.Println("[BookRepository] Create() Criando novo documento de livro no Mongo.")
	collection := db.Collection("books")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, book)

	if err != nil {
		log.Fatal("Erro ao inserir novo livro..")
	}
	return result.InsertedID
}

//Delete Book by ID
func (br *BookRepository) Delete(id string) (int64, error) {
	col := db.Collection("books")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)

	result, err := col.DeleteOne(ctx, bson.M{"_id": objID})

	if err != nil || result.DeletedCount != 1 {
		return -1, errors.New("Nenhum registro encontrado")
	}
	return result.DeletedCount, err

}

//Update Book by ID
func (br *BookRepository) Update(book api.Book) (int64, error) {
	col := db.Collection("books")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := col.UpdateOne(
		ctx,
		bson.M{"_id": book.ID},
		bson.D{
			{"$set", bson.D{
				{"author", book.Author},
				{"title", book.Title},
				{"pages", book.Pages},
			}},
		},
	)

	return result.ModifiedCount, err

}
