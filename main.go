package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Serein-sz/mission-ctrl/model"
	"github.com/Serein-sz/mission-ctrl/repository"
	"github.com/Serein-sz/mission-ctrl/scraper"
	tea "github.com/charmbracelet/bubbletea"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	fmt.Println("Connected to the sqlite database successfully.")
	repository.AutoMigrate(db)
	m := model.Model{}
	scraper.StartFetchData(&m, db)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
