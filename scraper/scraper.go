// Package scraper
package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Serein-sz/mission-ctrl/model"
	"github.com/Serein-sz/mission-ctrl/repository"
	"github.com/gocolly/colly/v2"
	"gorm.io/gorm"
)

func StartFetchData(m *model.Model, db *gorm.DB) {
	colly.Async(true)
	c := colly.NewCollector(
		colly.AllowedDomains("192.168.2.240", "192.168.2.111"),
	)
	c.Limit(&colly.LimitRule{
		Parallelism: 50,
	})

	c.OnHTML("body > div:nth-child(3) > div:nth-child(2) > table > tbody:nth-child(2) > tr", func(e *colly.HTMLElement) {
		e.ForEach("td.left > span:nth-child(2) > a", func(i int, h *colly.HTMLElement) {
			fmt.Printf("Handle repository: %s\n", h.Text)
			e.Request.Visit(extractURI(e) + "/branches/" + h.Text + ".git")
		})
	})

	c.OnHTML("body > div:nth-child(4) > div > table > tbody > tr > td:nth-child(2) > span", func(e *colly.HTMLElement) {
		e.ForEach("a", func(i int, h *colly.HTMLElement) {
			for i := range 1 {
				e.Request.Visit(extractURI(e) + h.Attr("href")[2:] + "?pg=" + strconv.Itoa(i+1))
			}
		})
	})
	tasks := []repository.Task{}
	c.OnHTML("body > div:nth-child(4) > div:nth-child(2) > table > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr.commit", func(i int, h *colly.HTMLElement) {
			commiter := h.ChildText("td.hidden-phone.author > span > a")
			date := formatDate(h.ChildText("td.date > span"))
			description := h.ChildAttr("td.message span", "title")
			if description == "" {
				description = h.ChildText("td.message a")
			}
			task := repository.NewTask(commiter, date, description, extractRepository(e.Request.URL.EscapedPath()))
			tasks = append(tasks, task)
			m.AddValue(task.String())
		})
	})

	c.Visit("http://192.168.2.111:19999/repositories")
	c.Visit("http://192.168.2.240:9999/repositories")
	fmt.Printf("len(tasks): %v\n", len(tasks))
	db.CreateInBatches(tasks, 50)
}

func extractURI(e *colly.HTMLElement) string {
	return e.Request.URL.Scheme + "://" + e.Request.Host
}

func formatDate(originDate string) string {
	now := time.Now()
	pattern := "2006-01-02"
	switch originDate {
	case "刚刚":
		return now.Format(pattern)
	case "昨天":
		return now.AddDate(0, 0, -1).Format(pattern)
	default:
		if strings.HasSuffix(originDate, "小时以前") {
			return now.Format(pattern)
		}
		if strings.HasSuffix(originDate, "天以前") {
			day, _ := strconv.Atoi(strings.Split(originDate, " ")[0])
			return now.AddDate(0, 0, -day).Format(pattern)
		}
	}
	return originDate
}

func extractRepository(url string) string {
	re := regexp.MustCompile(`/([^/]*).git/`)
	group := re.FindStringSubmatch(url)
	if len(group) > 1 {
		return group[1]
	}
	return ""
}
