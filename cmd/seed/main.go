package main

import (
	"fmt"
	repo2 "shop/internal/list"
	"shop/internal/models"
	entity2 "shop/internal/user"
	"log"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const passwordLength = 8

type Seeder struct {
	listRepo *repo2.Repository
	authRepo *entity2.Repository
}

func NewSeeder(_listRepo *repo2.Repository, _authRepo *entity2.Repository) *Seeder {
	return &Seeder{
		listRepo: _listRepo,
		authRepo: _authRepo,
	}
}

func main() {
	const (
		userCount        = 10
		listPerUserCount = 100
	)

	db, err := sqlx.Connect("pgx", os.Getenv("SHARED_DB_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer func(db *sqlx.DB) {
		_ = db.Close()
	}(db)
	log.Println("Connected to DB for seeding...")

	listRepo := repo2.NewRepository(db)
	authRepo := entity2.NewRepository(db)

	seed := NewSeeder(listRepo, authRepo)

	// err = seed.SeedUsers(userCount)
	// if err != nil {
	//	logrus.Errorf("Error while seeding Users %s", err.Error())
	//}

	err = seed.SeedLists(listPerUserCount)
	if err != nil {
		return
	}

	log.Println("Seeding completed successfully!")
}

func (s *Seeder) SeedLists(countPerUser int) error {
	const defaultCount = 100
	users, err := s.authRepo.GetAll()
	if err != nil {
		logrus.Errorf("Error when try to get users for seedLists: %s", err.Error())
	}

	count := (func() int {
		if countPerUser > 0 {
			return countPerUser
		}

		return defaultCount
	})()

	for _, userEl := range users {
		newList := models.List{
			UserID:      userEl.ID,
			Title:       gofakeit.BookTitle(),
			Description: gofakeit.ProductDescription(),
		}
		for range count {
			_, err := s.listRepo.Create(userEl.ID, newList)

			logrus.WithFields(logrus.Fields{
				"UserID":      newList.UserID,
				"Title":       newList.Title,
				"Description": newList.Description,
			}).Info("Created list content")

			if err != nil {
				fmt.Printf("Error in seedingLists %s", err.Error())
				return err
			}
		}
	}

	return nil
}

func (s *Seeder) SeedUsers(size int) {
	for range make([]int, size) {
		newUser := models.User{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Password: gofakeit.Password(false, false, false, false, false, passwordLength),
		}
		logrus.WithFields(logrus.Fields{
			"ID":        newUser.ID,
			"Name":      newUser.Name,
			"Username":  newUser.Username,
			"Password":  newUser.Password,
			"AvatarURL": newUser.AvatarURL,
			"Phone":     newUser.Phone,
			"Role":      newUser.Role,
			"IsActive":  newUser.IsActive,
			"CreatedAt": newUser.CreatedAt,
			"UpdatedAt": newUser.UpdatedAt,
		}).Info("Created user content")

		_, err := s.authRepo.CreateUser(&newUser)
		if err != nil {
			fmt.Printf("Error in seedUsers %s", err.Error())
		}
	}
}
