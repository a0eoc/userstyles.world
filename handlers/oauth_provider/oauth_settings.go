package oauth_provider

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func OAuthSettingsGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")
	isEdit := id != ""

	var method string
	var oauth *models.APIOAuth = nil
	var err error
	if isEdit {
		method = "edit"
		oauth, err = models.GetOAuthByID(database.DB, id)
	} else {
		method = "add"
	}

	if err != nil {
		fmt.Printf("Failed to oauth, err: %#+v\n", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	if isEdit {
		if u.ID != oauth.UserID {
			return c.Render("err", fiber.Map{
				"Title": "Users don't match",
				"User":  u,
			})
		}

	}
	oauths, err := models.ListOAuthsOfUser(database.DB, u.Username)
	if err != nil {
		if isEdit {
			return c.Render("oauth_settings", fiber.Map{
				"Title":  "OAuth Settings",
				"User":   u,
				"OAuth":  oauth,
				"OAuths": nil,
				"Method": method,
			})
		}
		return c.Render("oauth_settings", fiber.Map{
			"Title":  "OAuth Settings",
			"User":   u,
			"OAuths": nil,
			"Method": method,
		})
	}

	if isEdit {
		return c.Render("oauth_settings", fiber.Map{
			"Title":  "OAuth Settings",
			"User":   u,
			"OAuth":  oauth,
			"OAuths": oauths,
			"Method": method,
		})
	}
	return c.Render("oauth_settings", fiber.Map{
		"Title":  "OAuth Settings",
		"User":   u,
		"OAuths": oauths,
		"Method": method,
	})
}

func OAuthSettingsPost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

	q := models.OAuth{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		RedirectURI: strings.TrimSuffix(c.FormValue("redirect_uri"), "/"),
		Scopes: utils.Filter([]string{"style", "user"}, func(name interface{}) bool {
			return c.FormValue(name.(string)) == "on"
		}).([]string),
		UserID: u.ID,
	}

	// TODO: Validate this.

	var err error
	if id != "" {
		err = models.UpdateOAuth(database.DB, &q)
	} else {
		q.ClientID = utils.B2s((utils.RandStringBytesMaskImprSrcUnsafe(32)))
		q.ClientSecret = utils.B2s((utils.RandStringBytesMaskImprSrcUnsafe(128)))
		_, err = models.CreateOAuth(database.DB, &q)
	}

	if err != nil {
		log.Println("Updating style failed, err:", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	return c.Redirect("/oauth_settings/"+id, fiber.StatusSeeOther)

}