package main

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"go-shop/internal/auth/repository"
	"go-shop/internal/list/entity"
	repo2 "go-shop/internal/list/repository"
	entity2 "go-shop/internal/user/entity"
	"log"
	"os"
)

type Seeder struct {
	listRepo *repo2.ListRepository
	authRepo *repository.AuthRepository
}

func NewSeeder(_listRepo *repo2.ListRepository, _authRepo *repository.AuthRepository) *Seeder {
	return &Seeder{
		listRepo: _listRepo,
		authRepo: _authRepo,
	}
}

func main() {
	const (
		UserCount        = 10
		ListPerUserCount = 100
	)

	db, err := sqlx.Connect("pgx", os.Getenv("SHARED_DB_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer func(db *sqlx.DB) {
		_ = db.Close()
	}(db)
	log.Println("Connected to DB for seeding...")

	listRepo := repo2.NewListRepository(db)
	authRepo := repository.NewAuthRepository(db)

	seed := NewSeeder(listRepo, authRepo)

	//err = seed.SeedUsers(UserCount)
	//if err != nil {
	//	logrus.Errorf("Error while seeding Users %s", err.Error())
	//}

	err = seed.SeedLists(ListPerUserCount)
	if err != nil {
		return
	}

	log.Println("Seeding completed successfully!")
}

func (s *Seeder) SeedLists(countPerUser int) error {
	const DefaultCount = 100
	users, err := s.authRepo.GetAll()
	if err != nil {
		logrus.Errorf("Error when try to get users for seedLists: %s", err.Error())
	}

	count := (func() int {
		if countPerUser > 0 {
			return countPerUser
		}

		return DefaultCount
	})()

	for _, userEl := range users {
		newList := entity.List{
			UserId:      userEl.Id,
			Title:       gofakeit.BookTitle(),
			Description: gofakeit.ProductDescription(),
		}
		for range count {
			_, err := s.listRepo.Create(userEl.Id, newList)

			logrus.WithFields(logrus.Fields{
				"UserId":      newList.UserId,
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

func (s *Seeder) SeedUsers(size int) error {

	for range make([]int, size) {
		newUser := entity2.User{
			Name:     gofakeit.Name(),
			Username: gofakeit.Username(),
			Password: gofakeit.Password(false, false, false, false, false, 8),
		}
		logrus.WithFields(logrus.Fields{
			"Id":        newUser.Id,
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

		_, err := s.authRepo.CreateUser(newUser)
		if err != nil {
			fmt.Printf("Error in seedUsers %s", err.Error())
		}
	}
	return nil
}
