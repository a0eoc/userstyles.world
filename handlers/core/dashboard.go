package core

import (
	"log"
	"sort"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/charts"
)

func Dashboard(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	// Don't allow regular users to see this page.
	if u.Role < models.Moderator {
		return c.Render("err", fiber.Map{
			"Title": "Page not found",
			"User":  u,
		})
	}

	// Get styles.
	styles, err := models.GetAllAvailableStyles()
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Styles not found",
			"User":  u,
		})
	}

	sort.Slice(styles, func(i, j int) bool {
		return styles[i].ID > styles[j].ID
	})

	// Get users.
	users, err := models.FindAllUsers()
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Users not found",
			"User":  u,
		})
	}

	// Render user history.
	var userHistory string
	if len(users) > 0 {
		userHistory, err = charts.GetUserHistory(users)
		if err != nil {
			log.Printf("Failed to render user history, err: %s\n", err.Error())
		}
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID > users[j].ID
	})

	// Get history data.
	history, err := new(models.History).GetStatsForAllStyles()
	if err != nil {
		log.Printf("Couldn't find style histories, err: %s", err.Error())
	}

	// Render style history.
	var dailyHistory, totalHistory string
	if len(*history) > 0 {
		dailyHistory, totalHistory, err = charts.GetStyleHistory(*history)
		if err != nil {
			log.Printf("Failed to render history for all styles, err: %s\n", err.Error())
		}
	}

	return c.Render("core/dashboard", fiber.Map{
		"Title":        "Dashboard",
		"User":         u,
		"Styles":       styles,
		"Users":        users,
		"DailyHistory": dailyHistory,
		"TotalHistory": totalHistory,
		"UserHistory":  userHistory,
	})
}
