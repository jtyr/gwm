package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"

	log "github.com/sirupsen/logrus"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func httpError(w http.ResponseWriter, msg string, code int) {
	finalMsg := fmt.Sprintf("%d %s - %s", code, http.StatusText(code), msg)

	http.Error(w, finalMsg, code)
}

func replaceString(input, search, replace string) (result string) {
	log.Trace("Replacing string")

	m := regexp.MustCompile(search)

	log.Tracef("String before replacement: %s", input)

	result = m.ReplaceAllString(input, replace)

	log.Tracef("String after replacement: %s", result)

	return
}

func replaceStringInArray(data []any, search, replace string) {
	log.Trace("Replacing strings in Array")

	for i, v := range data {
		if v == nil {
			continue
		}

		rt := reflect.TypeOf(v)

		switch rt.Kind() {
		case reflect.Map:
			replaceStringInMap(v.(map[string]any), search, replace)
		case reflect.Slice:
			replaceStringInArray(v.([]any), search, replace)
		case reflect.String:
			v = replaceString(v.(string), search, replace)
		}

		data[i] = v
	}
}

func replaceStringInMap(data map[string]any, search, replace string) {
	log.Trace("Replacing strings in Map")

	for k, v := range data {
		if v == nil {
			continue
		}

		rt := reflect.TypeOf(v)

		switch rt.Kind() {
		case reflect.Map:
			replaceStringInMap(v.(map[string]any), search, replace)
		case reflect.Slice:
			replaceStringInArray(v.([]any), search, replace)
		case reflect.String:
			v = replaceString(v.(string), search, replace)
		}

		data[k] = v
	}
}

func processWebhookRequest(r http.Request) (string, http.Header, error) {
	// Read query body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read request body: %s", err)
	}

	searchFor := getEnv("GWM_SEARCH", "^http://git.localhost/")
	replaceWith := getEnv("GWM_REPLACE", "http://gitea-http.gitea:3000/")

	// Convert JSON string to data structure
	data := make(map[string]any)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", nil, fmt.Errorf("failed to unmarshal JSON body: %s", err)
	}

	log.Tracef("Data before replacement: %s", data)

	// Replace strings
	replaceStringInMap(data, searchFor, replaceWith)

	log.Tracef("Data after replacement: %s", data)

	// Convert data structure to JSON string
	newBody, err := json.Marshal(data)
	if err != nil {
		return "", nil, fmt.Errorf("failed to marshal JSON data: %s", err)
	}

	forwardTo := getEnv("GWM_FORWARD", "http://argocd-server.argocd/api/webhook")

	// Create a new client
	req, err := http.NewRequest("POST", forwardTo, bytes.NewBuffer(newBody))
	if err != nil {
		return "", nil, fmt.Errorf("failed to create a new client POST request: %s", err)
	}

	// Copy Headers from the request to the new client
	req.Header = r.Header

	log.Trace("Forwarding request")

	// Run the client query
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("failed to perform client POST request: %s", err)
	}

	// Read body from the client response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read response body: %s", err)
	}

	return string(respBody), resp.Header, nil
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpError(w, "Only POST method is allowed", 405)

		return
	}

	log.Debug("Got POST request")

	// Process the request
	respBody, respHeader, err := processWebhookRequest(*r)
	if err != nil {
		msg := fmt.Sprintf("Failed to process webhook request: %s", err)
		log.Error(msg)
		httpError(w, msg, 500)

		return
	}

	// Copy Headers from the client response
	for key, vals := range respHeader {
		for _, v := range vals {
			w.Header().Set(key, v)
		}
	}

	if _, err := io.WriteString(w, respBody); err != nil {
		log.Errorf("Failed to write string: %s", err)
	}
}

func healthyHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Got request")

	if _, err := io.WriteString(w, "healthy\n"); err != nil {
		log.Errorf("Failed to write string: %s", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpError(w, "Only POST method is allowed", 405)

		return
	}

	log.Debug("Got POST request")

	msg, _ := io.ReadAll(r.Body)

	log.Debugf("Body: %s", msg)

	if _, err := io.WriteString(w, "Webhook received!\n"); err != nil {
		log.Errorf("Failed to write string: %s", err)
	}
}

func main() {
	// Configure logger
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fnName := fmt.Sprintf("%s()", frame.Function)
			fileName := fmt.Sprintf("%s:%d", path.Base(frame.File), frame.Line)

			return fnName, fileName
		},
	})

	logLevel := getEnv("GWM_LOG_LEVEL", "info")

	switch logLevel {
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "trace":
		log.SetLevel(log.TraceLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	// Define web seerver endpoints
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthy", healthyHandler)
	http.HandleFunc("/webhook", webhookHandler)

	host := getEnv("GWM_HOST", "0.0.0.0")
	port := getEnv("GWM_PORT", "8080")

	hostPort := fmt.Sprintf("%s:%s", host, port)

	log.Infof("Starting server on %s", hostPort)

	// Start HTTP server
	err := http.ListenAndServe(hostPort, nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Error("Server closed")
	} else if err != nil {
		log.Errorf("Error starting server: %s", err)

		os.Exit(1)
	}
}
