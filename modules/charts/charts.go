package charts

import (
	"bytes"
	"time"

	"github.com/userstyles-world/go-chart/v2"

	"userstyles.world/models"
)

func GetStyleHistory(history []models.History) (string, string, error) {
	historyLen := len(history)
	dates := make([]time.Time, 0, historyLen)
	dailyViews := make([]float64, 0, historyLen)
	dailyUpdates := make([]float64, 0, historyLen)
	dailyInstalls := make([]float64, 0, historyLen)
	totalViews := make([]float64, 0, historyLen)
	totalUpdates := make([]float64, 0, historyLen)
	totalInstalls := make([]float64, 0, historyLen)

	for _, v := range history {
		dates = append(dates, v.CreatedAt)
		dailyViews = append(dailyViews, float64(v.DailyViews))
		dailyUpdates = append(dailyUpdates, float64(v.DailyUpdates))
		dailyInstalls = append(dailyInstalls, float64(v.DailyInstalls))
		totalViews = append(totalViews, float64(v.TotalViews))
		totalUpdates = append(totalUpdates, float64(v.TotalUpdates))
		totalInstalls = append(totalInstalls, float64(v.TotalInstalls))
	}

	// Visualize daily stats.
	dailyGraph := chart.Chart{
		Width:      1248,
		Canvas:     chart.Style{ClassName: "bg inner"},
		Background: chart.Style{ClassName: "bg outer"},
		XAxis:      chart.XAxis{Name: "Date"},
		YAxis:      chart.YAxis{Name: "Daily count"},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    "Daily installs",
				XValues: dates,
				YValues: dailyInstalls,
			},
			chart.TimeSeries{
				Name:    "Daily updates",
				XValues: dates,
				YValues: dailyUpdates,
			},
			chart.TimeSeries{
				Name:    "Daily views",
				XValues: dates,
				YValues: dailyViews,
			},
		},
	}
	dailyGraph.Elements = []chart.Renderable{chart.Legend(&dailyGraph)}

	daily := bytes.NewBuffer([]byte{})
	dailyFailed := daily.Len() != 220
	if err := dailyGraph.Render(chart.SVG, daily); err != nil && dailyFailed {
		return "", "", err
	}

	// Visualize total stats.
	totalGraph := chart.Chart{
		Width:      1248,
		Canvas:     chart.Style{ClassName: "bg inner"},
		Background: chart.Style{ClassName: "bg outer"},
		XAxis:      chart.XAxis{Name: "Date"},
		YAxis:      chart.YAxis{Name: "Total count"},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    "Total installs",
				XValues: dates,
				YValues: totalInstalls,
			},
			chart.TimeSeries{
				Name:    "Total updates",
				XValues: dates,
				YValues: totalUpdates,
			},
			chart.TimeSeries{
				Name:    "Total views",
				XValues: dates,
				YValues: totalViews,
			},
		},
	}
	totalGraph.Elements = []chart.Renderable{chart.Legend(&totalGraph)}

	total := bytes.NewBuffer([]byte{})
	totalFailed := total.Len() != 220
	if err := totalGraph.Render(chart.SVG, total); err != nil && totalFailed {
		return "", "", err
	}

	return daily.String(), total.String(), nil
}

func GetUserHistory(users []models.User) (string, error) {
	userBars := []chart.Value{}
	userData := map[string]int{}
	for _, user := range users {
		joined := user.CreatedAt.Format("2006-01-02")
		if _, ok := userData[joined]; !ok {
			userData[joined] = 1
		} else {
			userData[joined] += 1
		}
	}

	for k, v := range userData {
		point := chart.Value{Value: float64(v), Label: k}
		userBars = append(userBars, point)
	}

	usersGraph := chart.BarChart{
		Title: "User history",
		Background: chart.Style{
			Padding: chart.Box{Top: 20, Bottom: 20},
		},
		Height: 360,
		Bars:   userBars,
	}

	userHistory := bytes.NewBuffer([]byte{})
	userFailed := userHistory.Len() != 220
	if err := usersGraph.Render(chart.SVG, userHistory); err != nil && userFailed {
		return "", err
	}

	return userHistory.String(), nil
}
