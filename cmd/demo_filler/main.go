package main

import (
	"fmt"
	"os"

	"em-test-go/src/models"

	"github.com/google/uuid"
)

const demoSubscriptionCount = 40

var demoUserIDs = []uuid.UUID{
	uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba"),
	uuid.MustParse("a1b2c3d4-e5f6-4789-a012-3456789abcde"),
	uuid.MustParse("b2c3d4e5-f6a7-4890-b123-456789abcdef"),
	uuid.MustParse("c3d4e5f6-a7b8-4901-c234-56789abcdef0"),
	uuid.MustParse("d4e5f6a7-b8c9-4012-d345-6789abcdef01"),
	uuid.MustParse("e5f6a7b8-c9d0-4123-e456-789abcdef012"),
	uuid.MustParse("f6a7b8c9-d0e1-4234-f567-89abcdef0123"),
	uuid.MustParse("a7b8c9d0-e1f2-4345-a678-9abcdef01234"),
}

var demoServices = []struct {
	name  string
	price int
}{
	{"Yandex Plus", 400},
	{"Spotify", 299},
	{"Netflix", 599},
	{"Apple Music", 199},
	{"Kinopoisk", 399},
	{"ivi", 299},
	{"OKKO", 499},
	{"Wink", 349},
	{"YouTube Premium", 249},
	{"Microsoft 365", 799},
	{"Amazon Prime", 449},
	{"VK Music", 149},
	{"Telegram Premium", 299},
	{"ChatGPT Plus", 1990},
	{"Cloud Storage 100GB", 99},
}

var demoStartDates = []string{
	"01-2024", "02-2024", "03-2024", "04-2024", "05-2024", "06-2024",
	"07-2024", "08-2024", "09-2024", "10-2024", "11-2024", "12-2024",
	"01-2025", "02-2025", "03-2025", "04-2025", "05-2025", "06-2025",
	"07-2025", "08-2025", "09-2025", "10-2025", "11-2025", "12-2025",
}

func buildDemoSubscriptions() []models.Subscription {
	subs := make([]models.Subscription, 0, demoSubscriptionCount)
	for i := 0; i < demoSubscriptionCount; i++ {
		service := demoServices[i%len(demoServices)]
		subs = append(subs, models.Subscription{
			ServiceName: service.name,
			Price:       service.price + (i%5)*50,
			UserID:      demoUserIDs[i%len(demoUserIDs)],
			StartDate:   demoStartDates[i%len(demoStartDates)],
		})
	}
	return subs
}

func main() {
	models.OpenDatabaseConnection()
	models.AutoMigrateModels()

	subs := buildDemoSubscriptions()
	if err := models.Database.Create(&subs).Error; err != nil {
		fmt.Fprintf(os.Stderr, "failed to insert demo subscriptions: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("inserted %d demo subscriptions\n", len(subs))
}
