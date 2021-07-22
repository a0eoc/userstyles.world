package update

import (
	"strings"

	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/models"
	"userstyles.world/modules/log"
	"userstyles.world/search"
	"userstyles.world/utils"
)

const archiveURL = "https://raw.githubusercontent.com/33kk/uso-archive/"

func isUSo(url string) bool {
	return strings.HasPrefix(url, archiveURL)
}

func updateMeta(a, b *models.Style) bool {
	if a.Name == b.Name &&
		a.Notes == b.Notes &&
		a.Preview == b.Preview &&
		a.Description == b.Description {
		return false
	}

	return true
}

func getSourceCode(style models.Style) string {
	if style.MirrorURL != "" {
		return style.MirrorURL
	}

	return style.Original
}

func Batch(batch models.Style) {
	if batch.Archived {
		return
	}

	s := new(models.Style)
	s.ID = batch.ID

	// Update style metadata if style comes from USo-archive.
	if isUSo(batch.Original) && batch.MirrorMeta {
		importedStyle, err := utils.ImportFromArchive(batch.Original, models.APIUser{})
		if err != nil {
			log.Warn.Println("Failed to import a style from archive:", err.Error())
		}

		// Run update if fields differ.
		if updateMeta(&batch, importedStyle) {
			s.Name = importedStyle.Name
			s.Notes = importedStyle.Notes
			s.Preview = importedStyle.Preview
			s.Description = importedStyle.Description

			if err = models.UpdateStyle(s); err != nil {
				log.Warn.Printf("Failed to mirror meta for %d: %s\n", batch.ID, err.Error())
			}
			if err = search.IndexStyle(s.ID); err != nil {
				log.Warn.Printf("Failedto re-index style %d: %s\n", s.ID, err.Error())
			}
		}
	}

	// Get new style source code.
	style, err := usercss.ParseFromURL(getSourceCode(batch))
	if err != nil {
		log.Warn.Printf("Failed to parse a style %d from URL: %s\n", batch.ID, err.Error())
		return
	}

	// Exit if source code doesn't pass validation.
	if errs := usercss.BasicMetadataValidation(style); errs != nil {
		log.Warn.Printf("Failed to validate style %d.\n", batch.ID)
		return
	}

	// Mirror source code if versions don't match.
	if style.Version != usercss.ParseFromString(batch.Code).Version {
		log.Warn.Printf("Version for style %d was changed.\n", batch.ID)
		s.Code = style.SourceCode
		if err := models.UpdateStyle(s); err != nil {
			log.Warn.Printf("Failed to mirror code for %d: %s\n", batch.ID, err)
		}
	}
}
