package translation

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode"
)

// ResponseData represents the response data from MyMemory API
type ResponseData struct {
	TranslatedText string `json:"translatedText"`
}

// Match represents individual match data from MyMemory API
type Match struct {
	Segment     string          `json:"segment"`
	Translation string          `json:"translation"`
	Source      string          `json:"source"`
	Target      string          `json:"target"`
	Quality     json.RawMessage `json:"quality"`
	Match       float64         `json:"match"`
}

// MyMemoryResponse represents the complete response from MyMemory API
type MyMemoryResponse struct {
	ResponseData ResponseData `json:"responseData"`
	Matches      []Match      `json:"matches"`
}

// processMatchQuality processes the quality field and converts it to float64
func processMatchQuality(rawQuality json.RawMessage) (float64, error) {
	var qualityFloat float64
	var qualityString string

	// Attempt to unmarshal the field as float64
	if err := json.Unmarshal(rawQuality, &qualityFloat); err == nil {
		return qualityFloat, nil
	}

	// If that fails, try to unmarshal the field as string
	if err := json.Unmarshal(rawQuality, &qualityString); err == nil {
		return strconv.ParseFloat(qualityString, 64)
	}

	return 0, fmt.Errorf("failed to parse quality field")
}

// Function translate text
func TranslateText(sourceText, sourceLang, targetLang string) (string, error) {
	baseURL := "https://api.mymemory.translated.net/get"

	params := url.Values{}
	params.Add("q", sourceText)
	params.Add("langpair", sourceLang+"|"+targetLang)

	finalURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(finalURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get valid response, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("Response Body:", string(body))

	var result MyMemoryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	log.Println("result: ", result)
	var bestTranslation string
	var bestQuality float64

	for _, match := range result.Matches {
		quality, err := processMatchQuality(match.Quality)
		if err != nil {
			fmt.Println("Error processing quality:", err)
			continue
		}
		log.Println("qulity is : ", quality)

		if quality > bestQuality && len(match.Translation) > 0 && !containsWeirdCharacters(match.Translation) {
			bestQuality = quality
			bestTranslation = match.Translation
			log.Println("translate text after qulity: ", bestTranslation)
		}
	}

	if bestTranslation == "" {
		bestTranslation = result.ResponseData.TranslatedText
		log.Println("best translate in if 2: ", bestTranslation)
		if bestTranslation == "" {
			return "", fmt.Errorf("no valid translation found")
		}
	}
	log.Println("best translate text before html: ", bestTranslation)
	bestTranslation = html.UnescapeString(bestTranslation)
	log.Println("best translate text after html: ", bestTranslation)

	return bestTranslation, nil
}

// A function to detect unusual strings
func containsWeirdCharacters(s string) bool {

	// List of HTML characters to check
	htmlChars := []string{"<", ">", "&", "\"", "'"}

	for _, char := range htmlChars {
		if strings.Contains(s, char) {
			return true
		}
	}

	for _, char := range s {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) && !unicode.IsPunct(char) {
			return true
		}
	}

	return false
}
