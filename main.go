package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

var sessions = make(map[string]string)

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	// 32bytesの乱数を生成
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func setSession(w http.ResponseWriter, userID string) error {
	sessionID, err := generateSessionID()
	if err != nil {
		return err
	}
	sessions[sessionID] = userID
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
	return nil
}

func getSession(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}
	userID, ok := sessions[cookie.Value]
	if !ok {
		return "", fmt.Errorf("invalid session")
	}
	return userID, nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	err := setSession(w, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "logged in as %s\n", userID)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	fmt.Fprintf(w, "Welcome, %s!", userID)
}

func main() {
	http.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "login.html")
	})
	http.HandleFunc("POST /login", loginHandler)
	http.HandleFunc("/", mainHandler)
	http.ListenAndServe(":8080", nil)
}
