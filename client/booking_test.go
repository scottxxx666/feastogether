package client_test

import (
	"feastogether/client"
	"feastogether/config"
	"log"
	"testing"
)

var cfg *config.Config

func init() {
	var err error
	cfg, err = config.GetConfig("..")
	if err != nil {
		log.Panicf("Failed to load config file: %v", err)
	}
}

func TestGetToken(t *testing.T) {
	if token := client.GetToken(cfg.UserConfig); token == "" {
		t.Errorf("Failed to get token")
	} else {
		t.Log(token)
	}
}

func TestGetSaveSaets(t *testing.T) {
	if saveSaets := client.GetSaveSaets(
		cfg.UserConfig,
		client.GetToken(cfg.UserConfig)); saveSaets == "" {
		t.Errorf("Failed to get saveSaets")
	} else {
		t.Log(saveSaets)
	}
}

func TestGetSaveSeats(t *testing.T) {
	if saveSeats := client.GetSaveSeats(
		cfg.UserConfig,
		client.GetToken(cfg.UserConfig),
		cfg.RestaurantConfig); saveSeats == "" {
		t.Errorf("Failed to get saveSeats")
	} else {
		t.Log(saveSeats)
	}
}

func TestB00king(t *testing.T) {
	if b00king := client.GetB00king(
		cfg.UserConfig,
		client.GetToken(cfg.UserConfig)); b00king == "" {
		t.Errorf("Failed to get b00king")
	} else {
		t.Log(b00king)
	}
}

func TestSaveBooking(t *testing.T) {
	client.GetSaveSeats(
		cfg.UserConfig,
		client.GetToken(cfg.UserConfig),
		cfg.RestaurantConfig)

	if booking := client.SaveBooking(
		cfg.UserConfig,
		client.GetToken(cfg.UserConfig),
		cfg.RestaurantConfig); booking == "" {
		t.Errorf("Booking failed")
	} else {
		t.Log(booking)
	}
}
