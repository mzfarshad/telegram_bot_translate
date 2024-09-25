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

// Function translate text using MyMemory API
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

	// Save the initial translation
	initialTranslation := result.ResponseData.TranslatedText
	initialQuality := 80.0 // You can change this value

	var bestTranslation string
	var bestScore float64

	for _, match := range result.Matches {
		quality, err := processMatchQuality(match.Quality)
		if err != nil {
			fmt.Println("Error processing quality:", err)
			continue
		}
		log.Println("quality: ", quality)

		// Scoring based on quality, keywords and structure
		score := scoreTranslation(sourceText, match.Translation, quality)
		log.Println("Translation score: ", score)

		if score > bestScore && len(match.Translation) > 0 && !containsWeirdCharacters(match.Translation) {
			bestScore = score
			bestTranslation = match.Translation
			log.Println("Best translation so far: ", bestTranslation)
		}
	}

	// Use the original translation in the absence of a better translation
	if bestTranslation == "" || bestScore < initialQuality {
		bestTranslation = initialTranslation
		log.Println("Using initial translation: ", bestTranslation)
	}

	bestTranslation = html.UnescapeString(bestTranslation)

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

// Scoring based on keywords
func scoreTranslationByKeywords(sourceText, translatedText string) float64 {
	keyWords := extractKeywords(sourceText)

	matchedCount := 0
	for _, word := range keyWords {
		if strings.Contains(translatedText, word) {
			matchedCount++
		}
	}

	// Score based on keyword matches
	return (float64(matchedCount) / float64(len(keyWords))) * 100
}

// Extracting keywords from the text
func extractKeywords(text string) []string {
	words := strings.Fields(text)
	var keyWords []string
	for _, word := range words {
		if len(word) > 3 { // Choose words longer than three characters
			keyWords = append(keyWords, word)
		}
	}
	return keyWords
}

// Scoring based on sentence structure
func scoreTranslationBySentenceStructure(sourceText, translatedText string) float64 {
	sourceSentences := strings.Count(sourceText, ".") + strings.Count(sourceText, "!")
	translatedSentences := strings.Count(translatedText, ".") + strings.Count(translatedText, "!")

	// Scoring based on matching the number of sentences
	if sourceSentences == 0 {
		return 0
	}
	return (float64(translatedSentences) / float64(sourceSentences)) * 100
}

// Combination of scoring criteria
func scoreTranslation(sourceText, translatedText string, quality float64) float64 {
	keywordScore := scoreTranslationByKeywords(sourceText, translatedText)
	structureScore := scoreTranslationBySentenceStructure(sourceText, translatedText)

	// Combining points with weighting
	finalScore := (0.2 * keywordScore) + (0.3 * structureScore) + (0.5 * quality)

	return finalScore
}
