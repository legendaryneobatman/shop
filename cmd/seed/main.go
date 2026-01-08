package main

import (
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go-shop/internal/auth/repository"
	"go-shop/internal/list/entity"
	repo2 "go-shop/internal/list/repository"
	entity2 "go-shop/internal/user/entity"
	"log"
	"math/rand"
	"os"
)

func main() {
	SIZE := 100

	db, err := sqlx.Connect("pgx", os.Getenv("SHARED_DB_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer func(db *sqlx.DB) {
		_ = db.Close()
	}(db)

	log.Println("Connected to DB for seeding...")

	userIds, err := seedUsers(db, SIZE)
	if err != nil {
		return
	}
	seedLists(db, SIZE, *userIds)

	log.Println("Seeding completed successfully!")
}

func createExampleList(userId int) *entity.List {
	return &entity.List{
		UserId:      userId,
		Title:       gofakeit.BookTitle(),
		Description: gofakeit.ProductDescription(),
	}
}

func seedLists(db *sqlx.DB, size int, userIds []int) {
	repo := repo2.NewListRepository(db)

	userId := userIds[rand.Intn(size)]

	for range make([]int, size) {
		_, err := repo.Create(userId, *createExampleList(userId))
		if err != nil {
			return
		}
	}

}

func seedUsers(db *sqlx.DB, size int) (*[]int, error) {
	repo := repository.NewAuthRepository(db)
	userIds := make([]int, size)

	for index := range make([]int, size) {
		newUser := entity2.User{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Password: gofakeit.Password(false, false, false, false, false, 8),
		}
		b, err := json.MarshalIndent(newUser, "", " ")
		fmt.Print(string(b))
		createdUserId, err := repo.CreateUser(newUser)

		if err != nil {
			return nil, err
		}

		userIds[index] = createdUserId
	}

	return &userIds, nil
}
