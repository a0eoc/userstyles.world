package user

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
	"userstyles.world/utils"
)

func setSocials(u *models.User, k, v string) {
	switch k {
	case "github":
		if v != "" {
			u.Socials.Github = v
		}
	case "gitlab":
		if v != "" {
			u.Socials.Gitlab = v
		}
	case "codeberg":
		if v != "" {
			u.Socials.Codeberg = v
		}
	}
}

func Account(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	log.Info.Println("User", u.ID, "visited account page.")

	styles, err := models.GetStylesByUser(u.Username)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Server error",
			"User":  u,
		})
	}

	user, err := models.FindUserByName(u.Username)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "User not found",
			"User":  u,
		})
	}

	return c.Render("user/account", fiber.Map{
		"Title":  "Account",
		"User":   u,
		"Params": user,
		"Styles": styles,
	})
}

func EditAccount(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	styles, err := models.GetStylesByUser(u.Username)
	if err != nil {
		return c.Render("err", fiber.Map{
			"User":  u,
			"Title": "Server error",
		})
	}

	user, err := models.FindUserByName(u.Username)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "User not found",
			"User":  u,
		})
	}

	name := strings.TrimSpace(c.FormValue("name"))
	if name != "" {
		prev := user.DisplayName
		user.DisplayName = name

		if err := utils.Validate().StructPartial(user, "DisplayName"); err != nil {
			var validationError validator.ValidationErrors
			if ok := errors.As(err, &validationError); ok {
				log.Info.Println("Validation errors:", validationError)
			}
			user.DisplayName = prev

			l := len(name)
			var e string
			switch {
			case l < 5 || l > 20:
				e = "Display name must be between 5 and 20 characters."
			default:
				e = "Make sure your input contains valid characters."
			}

			return c.Render("user/account", fiber.Map{
				"Title":  "Validation Error",
				"User":   u,
				"Params": user,
				"Styles": styles,
				"Error":  e,
			})
		}
	}

	bio := strings.TrimSpace(c.FormValue("bio"))
	if bio != "" {
		prev := user.Biography
		user.Biography = bio

		if err := utils.Validate().StructPartial(user, "Biography"); err != nil {
			var validationError validator.ValidationErrors
			if ok := errors.As(err, &validationError); ok {
				log.Info.Println("Validation errors:", validationError)
			}
			user.Biography = prev

			return c.Render("user/account", fiber.Map{
				"Title":  "Validation Error",
				"User":   u,
				"Params": user,
				"Styles": styles,
				"Error":  "Biography must be shorter than 512 characters.",
			})
		}
	}

	setSocials(user, "github", c.FormValue("github"))
	setSocials(user, "gitlab", c.FormValue("gitlab"))
	setSocials(user, "codeberg", c.FormValue("codeberg"))

	dbErr := database.Conn.
		Model(models.User{}).
		Where("id", user.ID).
		Updates(user).
		Error

	if dbErr != nil {
		log.Warn.Println("Updating user profile failed, err:", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	return c.Render("user/account", fiber.Map{
		"Title":  "Account",
		"User":   u,
		"Params": user,
		"Styles": styles,
	})
}
