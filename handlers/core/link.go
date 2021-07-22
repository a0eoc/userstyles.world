package core

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/modules/config"
)

func GetLinkedSite(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	switch c.Params("site") {
	case "chat":
		return c.Redirect(config.AppLinkChat, fiber.StatusSeeOther)
	case "source":
		return c.Redirect(config.AppLinkSource, fiber.StatusSeeOther)
	default:
		return c.Render("err", fiber.Map{
			"Title": "Invalid link",
			"User":  u,
		})
	}
}