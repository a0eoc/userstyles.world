package init

import (
	"os"
	"path"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"

	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
	"userstyles.world/utils"
)

var tables = []struct {
	name  string
	model any
}{
	{"users", &models.User{}},
	{"styles", &models.Style{}},
	{"stats", &models.Stats{}},
	{"oauths", &models.OAuth{}},
	{"histories", &models.History{}},
	{"logs", &models.Log{}},
	{"reviews", &models.Review{}},
	{"notifications", &models.Notification{}},
}

func connect() (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger: logger.New(
			log.Database,
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logLevel(),
				Colorful:      config.DBColor,
			},
		),
	}

	file := path.Join(config.DataDir, config.DB)
	conn, err := gorm.Open(sqlite.Open(file), gormConfig)
	if err != nil {
		return nil, err
	}

	conn.Use(prometheus.New(prometheus.Config{
		DBName:          "userstyles",
		RefreshInterval: 15,
		StartServer:     false,
	}))

	return conn, nil
}

// Initialize the database connection.
func Initialize() {
	conn, err := connect()
	if err != nil {
		log.Warn.Fatal("Failed to connect database:", err.Error())
	}

	database.Conn = conn
	log.Info.Println("Database successfully connected.")

	db, err := conn.DB()
	if err != nil {
		log.Warn.Fatal(err)
	}

	// GORM doesn't set a maximum of open connections by default.
	db.SetMaxOpenConns(config.DBMaxOpenConns)

	shouldSeed := false
	// Generate data for development.
	if config.DBDrop && !config.Production {
		for _, table := range tables {
			if err := drop(table.model); err != nil {
				log.Warn.Fatalf("Failed to drop %s, err: %s\n", table.name, err.Error())
			}
			log.Info.Println("Dropped database table", table.name)
		}
		shouldSeed = true
	}

	// Run one-time migrations.
	if _, ok := os.LookupEnv("MAGIC"); ok {
		migration()
	}

	// Migrate tables.
	if config.DBMigrate {
		for _, table := range tables {
			if err := migrate(table.model); err != nil {
				log.Warn.Fatalf("Failed to migrate %s, err: %s\n", table.name, err.Error())
			}

			log.Info.Println("Migrated table", table.name)
		}
	}

	if shouldSeed {
		seed()
	}

	// TODO: Simplify the entire process, including dropping and seeding data.
	if config.DBMigrate {
		log.Info.Println("Database migration complete.")
		os.Exit(0)
	}
}

// Close closes the database connection.
func Close() error {
	db, err := database.Conn.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

func generateData(amount int) ([]models.Style, []models.User) {
	randomData := utils.RandomString(amount * 7 * 4)
	var styleStructs []models.Style
	for i := 0; i < amount; i++ {
		startData := randomData[(i * 7 * 4):]
		styleStructs = append(styleStructs, models.Style{
			UserID:      uint(amount),
			Category:    startData[:4],
			Name:        startData[4:8],
			Description: startData[8:12],
			Notes:       startData[12:16],
			Preview:     startData[16:20],
			Code:        startData[20:24],
			Homepage:    startData[24:28],
			Featured:    true,
		})
	}

	var userStructs []models.User
	randomData = utils.RandomString(amount * 4 * 4)
	for i := 0; i < amount; i++ {
		startData := randomData[(i * 4 * 4):]
		userStructs = append(userStructs, models.User{
			Username:  startData[:4],
			Email:     startData[4:8],
			Biography: startData[8:12],
			Password:  startData[12:16],
		})
	}

	return styleStructs, userStructs
}

func seed() {
	log.Info.Println("Seeding database mock data.")
	defer log.Info.Println("Finished seeding mock data.")

	users := []models.User{
		{
			Username:  "admin",
			Email:     "admin@usw.local",
			Biography: "Admin of USw.",
			Password:  utils.GenerateHashedPassword("admin123"),
			Role:      models.Admin,
		},
		{
			Username:  "moderator",
			Email:     "moderator@usw.local",
			Biography: "I'm a moderator.",
			Password:  utils.GenerateHashedPassword("moderator"),
		},
		{
			Username: "regular",
			Email:    "regular@usw.local",
			Password: utils.GenerateHashedPassword("regular"),
		},
	}

	styles := []models.Style{
		{
			UserID:      1,
			Name:        "Dark-GitHub",
			Description: "Customizable dark theme for GitHub.",
			Notes:       "Some notes go here.",
			Preview:     "https://userstyles.world/preview/2/0t.webp",
			Original:    "https://raw.githubusercontent.com/vednoc/dark-github/main/github.user.styl",
			Homepage:    "https://github.com/vednoc/dark-github",
			Category:    "github.com",
			MirrorCode:  true,
			Featured:    true,
		},
		{
			UserID:      1,
			Name:        "Dark-GitLab",
			Description: "Customizable dark theme for GitLab.",
			Notes:       "Some notes go here.",
			Preview:     "https://userstyles.world/preview/3/0t.webp",
			Original:    "https://gitlab.com/vednoc/dark-gitlab/raw/master/gitlab.user.styl",
			Homepage:    "https://gitlab.com/vednoc/dark-gitlab",
			Category:    "gitlab.com",
			MirrorCode:  true,
			Featured:    true,
		},
		{
			UserID:      1,
			Name:        "Dark-WhatsApp",
			Description: "Customizable dark theme for WhatsApp.",
			Notes:       "Some notes go here.",
			Preview:     "https://userstyles.world/preview/4/0t.webp",
			Original:    "https://raw.githubusercontent.com/vednoc/dark-whatsapp/master/wa.user.styl",
			Homepage:    "https://github.com/vednoc/dark-whatsapp",
			Category:    "web.whatsapp.com",
			MirrorCode:  true,
			Featured:    true,
		},
		{
			UserID:   2,
			Name:     "Archived userstyle",
			Archived: true,
		},
		{
			UserID:   3,
			Name:     "Featured userstyle",
			Featured: true,
		},
		{
			UserID: 3,
			Name:   "Temporary userstyle",
		},
	}

	oauths := []models.OAuth{
		{
			UserID:       1,
			Name:         "USw integration",
			Description:  "Just some integration",
			Scopes:       []string{"user", "style"},
			ClientID:     "publicccc_client",
			ClientSecret: "secreettUwU",
			RedirectURI:  "https://gusted.xyz/callback_helper",
		},
	}

	logs := []models.Log{
		{
			UserID:         1,
			Reason:         "I like to abuse powers.",
			Kind:           models.LogBanUser,
			TargetUserName: "gusted",
		},
		{
			UserID:         1,
			Reason:         "My style is superior",
			Kind:           models.LogRemoveStyle,
			TargetUserName: "gusted",
			TargetData:     "Black-Discord",
		},
	}

	if config.DBRandomData {
		s, u := generateData(config.DBRandomDataAmount)
		styles = append(styles, s...)
		users = append(users, u...)
	}

	for i := range users {
		database.Conn.Create(&users[i])
	}
	for i := range styles {
		database.Conn.Create(&styles[i])
	}
	for i := range oauths {
		database.Conn.Create(&oauths[i])
	}
	for i := range logs {
		database.Conn.Create(&logs[i])
	}
}
