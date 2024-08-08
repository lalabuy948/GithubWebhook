package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var secretToken = os.Getenv("SECRET_TOKEN")

func verifySignature(payloadBody []byte, signature string) bool {
	gotHash := strings.SplitN(signature, "=", 2)
	if gotHash[0] != "sha1" {
		return false
	}

	hash := hmac.New(sha1.New, []byte(secretToken))
	if _, err := hash.Write(payloadBody); err != nil {
		log.Printf("Cannot compute the HMAC for request: %s\n", err)
		return false
	}

	expectedHash := hex.EncodeToString(hash.Sum(nil))
	return gotHash[1] == expectedHash
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	payloadBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	signature := r.Header.Get("X-Hub-Signature")
	if !verifySignature(payloadBody, signature) {
		log.Printf("Signatures didn't match!")
		http.Error(w, "Signatures didn't match!", http.StatusInternalServerError)
		return
	}

	var push map[string]interface{}
	if err := json.Unmarshal(payloadBody, &push); err != nil {
		log.Printf("Failed to parse JSON: %v", http.StatusInternalServerError)
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	headCommit, ok := push["head_commit"].(map[string]interface{})
	if !ok {
		log.Printf("Invalid webhook payload (not Head commit): %v", http.StatusInternalServerError)
		http.Error(w, "Invalid webhook payload (not Head commit)", http.StatusInternalServerError)
		return
	}

	commitID, ok := headCommit["id"].(string)
	if !ok {
		log.Printf("Invalid commit ID: %v", http.StatusInternalServerError)
		http.Error(w, "Invalid commit ID", http.StatusInternalServerError)
		return
	}

	ref, ok := push["ref"].(string)
	if !ok {
		log.Printf("Invalid ref: %v", http.StatusInternalServerError)
		http.Error(w, "Invalid ref", http.StatusInternalServerError)
		return
	}

	log.Printf("Webhook received... commit [%s]", commitID)
	if ref == "refs/heads/master" {
		log.Println("Starting pipeline...")
		cmd := exec.Command("./pipeline.sh")
		if err := cmd.Start(); err != nil {
			log.Printf("Failed to start pipeline: %v", err)
			http.Error(w, "Failed to start pipeline", http.StatusInternalServerError)
			return
		}
		log.Printf("Pipeline started with PID %d", cmd.Process.Pid)
	}
}

func main() {
	logFile, err := os.OpenFile("pipeline.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	http.HandleFunc("/webhook", webhookHandler)
	log.Println("Listening on port 4567...")
	if err := http.ListenAndServe(":4567", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
