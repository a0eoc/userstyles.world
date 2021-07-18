package style

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/database"
	"userstyles.world/utils"
)

func sendPromotionEmail(userID uint, style *models.APIStyle, modName, baseURL string) {
	user, err := models.FindUserByID(strconv.Itoa(int(userID)))
	if err != nil {
		log.Println("Couldn't find user:", err)
		return
	}

	modProfile := baseURL + "/user/" + modName

	partPlain := utils.NewPart().
		SetBody("Hi " + user.Username + ",\n" +
			"We'd like to notice you about a recent action from our moderation team:\n\n" +
			"Your style \"" + style.Name + "\" has been promoted to be featured on the homepage!\n" +
			"It has been promoted by the moderator" + modName + ", checkout their profile: " + modProfile + ".\n\n" +
			"Regards,\n" + "The Moderation Team")
	partHTML := utils.NewPart().
		SetBody("<p>Hi " + user.Username + ",</p>\n" +
			"<br>\n" +
			"<p>We'd like to notice you about a recent action from our moderation team:</p>\n" +
			"<br><br>\n" +
			"<p>Your style \"" + style.Name + "\" has been promoted to be featured on the homepage!\n</p>\n" +
			"<p>It has been promoted by the moderator " +
			"<a target=\"_blank\" clicktracking=\"off\" href=\"" + modProfile + "\">" + modName + "</a>.</p>\n" +
			"<br><br>\n" +
			"<p>Regards,</p>\n" + "<p>The Moderation Team</p>").
		SetContentType("text/html")

	err = utils.NewEmail().
		SetTo(user.Email).
		SetSubject("Style is being featured").
		AddPart(*partPlain).
		AddPart(*partHTML).
		SendEmail()
	if err != nil {
		log.Println("Couldn't send email:", err)
		return
	}
}

func Promote(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	p := c.Params("id")

	// Only moderator and above have permissions to promote styles.
	if u.Role < models.Moderator {
		return c.Render("err", fiber.Map{
			"Title": "You don't have enough permission for this.",
			"User":  u,
		})
	}

	id, err := strconv.Atoi(p)
	if err != nil {
		log.Printf("Couldn't convert %s to int, err: %s\n", p, err)
		return c.Render("err", fiber.Map{
			"Title": "Couldn't convert style ID",
			"User":  u,
		})
	}

	style, err := models.GetStyleByID(p)
	if err != nil {
		log.Println("Couldn't get the style:", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	err = database.Conn.
		Model(models.Style{}).
		Where("id = ?", p).
		Update("featured", !style.Featured).
		Error

	if err != nil {
		log.Println("Couldn't feature style:", err)
		return c.Render("err", fiber.Map{
			"Title": "Failed to promote a style",
			"User":  u,
		})
	}

	// Ahem!!! We don't save the new value of Featured to the current style.
	// So we have to reverse check it ;)
	if !style.Featured {
		go sendPromotionEmail(style.UserID, style, u.Username, c.BaseURL())

		// Create a notification.
		notification := models.Notification{
			Seen:     false,
			Kind:     models.KindStylePromotion,
			TargetID: int(style.UserID),
			UserID:   int(u.ID),
			StyleID:  id,
		}

		go func(notification models.Notification) {
			if err := notification.Create(); err != nil {
				log.Printf("Failed to create a notification for %d, err: %v", id, err)
			}
		}(notification)
	}

	return c.Redirect("/style/"+p, fiber.StatusSeeOther)
}
