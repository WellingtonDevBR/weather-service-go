package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/go-chi/chi/v5"
)

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func GetWeather(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	log.Printf("Received request for CEP: %s", cep)

	if !isValidCEP(cep) {
		log.Printf("Invalid CEP format: %s", cep)
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	location, err := getLocation(cep)
	if err != nil {
		log.Printf("Error getting location for CEP %s: %v", cep, err)
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	log.Printf("Found location for CEP %s: %s", cep, location)

	tempC, err := getTemperature(location)
	if err != nil {
		log.Printf("Error getting temperature for location %s: %v", location, err)
		http.Error(w, "error fetching temperature", http.StatusInternalServerError)
		return
	}

	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15

	response := WeatherResponse{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	log.Printf("Returning temperature for location %s: %+v", location, response)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func isValidCEP(cep string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(cep)
}

func getLocation(cep string) (string, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	log.Printf("Fetching location from URL: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("HTTP request failed: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected response status: %d", resp.StatusCode)
		return "", fmt.Errorf("error fetching location")
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode JSON response: %v", err)
		return "", err
	}

	if _, ok := result["erro"]; ok {
		log.Printf("CEP not found: %s", cep)
		return "", fmt.Errorf("CEP not found")
	}

	location, ok := result["localidade"].(string)
	if !ok {
		log.Printf("Failed to parse location from response: %+v", result)
		return "", fmt.Errorf("error parsing location")
	}

	return location, nil
}

func getTemperature(location string) (float64, error) {
	apiKey := "98a42a15c266432a98b25526240106"
	encodedLocation := url.QueryEscape(location)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, encodedLocation)
	log.Printf("Fetching temperature from URL: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("HTTP request failed: %v", err)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected response status: %d", resp.StatusCode)
		return 0, fmt.Errorf("error fetching temperature")
	}

	var result WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode JSON response: %v", err)
		return 0, err
	}

	return result.Current.TempC, nil
}
