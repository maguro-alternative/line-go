package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type LineResponses struct {
	Events []struct {
		ReplyToken string `json:"replyToken"`
		Type       string `json:"type"`
		Source     struct {
			GroupID string `json:"groupId"`
			UserID  string `json:"userId"`
			Type    string `json:"type"`
		} `json:"source"`
		Timestamp float64 `json:"timestamp"`
		Message   struct {
			ID                  string   `json:"id"`
			Text                string   `json:"text"`
			Duration            int64    `json:"duration"`
			FileName            string   `json:"fileName"`
			FileSize            int64    `json:"fileSize"`
			Title               string   `json:"title"`
			Address             string   `json:"address"`
			Latitude            float64  `json:"latitude"`
			Longitude           float64  `json:"longitude"`
			PackageID           string   `json:"packageId"`
			StickerID           string   `json:"stickerId"`
			StickerResourceType string   `json:"stickerResourceType"`
			Keywords            []string `json:"keywords"`
			ImageSet            struct {
				ID    string  `json:"id"`
				Index float64 `json:"index"`
				Total float64 `json:"total"`
			} `json:"imageSet"`
			ContentProvider struct {
				Type               string `json:"type"`
				OriginalContentURL string `json:"originalContentUrl"`
				PreviewImageURL    string `json:"previewImageUrl"`
			} `json:"contentProvider"`
		} `json:"message"`
		Mode            string `json:"mode"`
		WebhookEventID  string `json:"webhookEventId"`
		DeliveryContext struct {
			IsRedelivery bool `json:"isRedelivery"`
		} `json:"isRedelivery"`
	} `json:"events"`
}

func main() {
	http.HandleFunc("/", Handler)
	port := "8080"
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}


func Handler(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Verify request
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	xLineSignature := r.Header.Get("X-Line-Signature")
	// macの生成
	mac := hmac.New(sha256.New, []byte(channelSecret))
	mac.Write(body)
	validSignByte := mac.Sum(nil)

	signature := base64.StdEncoding.EncodeToString(validSignByte)

	if xLineSignature != signature {
		log.Printf("Invalid signature: %s = %s", xLineSignature, signature)
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}
	log.Printf("signature:%s = %s", xLineSignature, signature)

	// Parse request
	var lineResponses LineResponses
	if err := json.Unmarshal(body, &lineResponses); err != nil {
		http.Error(w, "can't parse body", http.StatusBadRequest)
		return
	}

	// Do something with the request
	for _, event := range lineResponses.Events {
		log.Printf("Got message: %s", event.Message.Text)
	}

	// Respond to LINE
	w.WriteHeader(http.StatusOK)
}