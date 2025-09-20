package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"time"

	escpos "github.com/justinmichaelvieira/escpos"
	"golang.org/x/image/font/opentype"
)

type Config struct {
	FetchUrl       string    `json:"fetch_url"`
	FontFile       string    `json:"font_file"`
	PrinterAddress string    `json:"printer_address"`
	FromTime       time.Time `json:"from_time"`
	ToTime         time.Time `json:"to_time"`
}

func loadConfig(filename string) (Config, error) {
	configBytes, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(configBytes, &config)

	return config, err
}

func loadFont(filename string) (*opentype.Font, error) {
	fontBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	loadedFont, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	return loadedFont, nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <config_file.json>\n", os.Args[0])
	}

	config, err := loadConfig(os.Args[1])
	if err != nil {
		log.Fatalf("could not read config file: %s", err.Error())
	}

	sights, err := getSights(config.FetchUrl)
	if err != nil {
		log.Fatalf("could not fetch sights: %s", err.Error())
	}

	loadedFont, err := loadFont(config.FontFile)
	if err != nil {
		log.Fatalf("could not load font: %s", err.Error())
	}

	conn, err := net.Dial("tcp", config.PrinterAddress)
	if err != nil {
		log.Fatalf("could not connect to printer: %s", err.Error())
	}
	defer conn.Close()

	printer := escpos.New(conn)
	printer.SetConfig(escpos.ConfigEpsonTMT88II)

	for _, sight := range sights {
		if sight.Timestamp.Before(time.Now()) {
			continue
		}

		img, err := renderSight(sight, loadedFont)
		if err != nil {
			log.Fatalf("could not render sight: %s", err.Error())
		}

		_, err = printer.PrintImage(img)
		if err != nil {
			log.Fatalf("could not print sight image: %s", err.Error())
		}
	}

	err = printer.PrintAndCut()
	if err != nil {
		log.Fatalf("print and cut failed: %s", err.Error())
	}
}
