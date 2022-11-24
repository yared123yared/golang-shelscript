package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	Command string
}

type ReturnResponse struct {
	response []string
	err      error
}

type ErrorResponse struct {
	err error
}

type OutputResponse struct {
	output string
}

func allArticles(w http.ResponseWriter, r *http.Request) {

	l := r.ContentLength
	body := make([]byte, l)
	r.Body.Read(body)
	command := &Command{}

	err := json.Unmarshal(body, command)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(command.Command)

	response := handleCommand(*&command.Command)
	fmt.Println("After response,", response)
	if response.err != nil {
		fmt.Print(response.err)
		// w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		output, err := json.MarshalIndent(response.err.Error(), "", "\t\t")

		if err != nil {
			fmt.Print(err)
		}
		_, _ = w.Write(output)
	} else {
		w.Header().Set("Content-Type", "application/json")
		output, err := json.MarshalIndent(response.response, "", "\t\t")

		if err != nil {
			// w.Header().Set("Content-Type", "application/json")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		_, _ = w.Write(output)
	}
}

func handelRequest() {
	http.HandleFunc("/api/cmd", allArticles)
	log.Fatal(http.ListenAndServe(":8001", nil))
}

func main() {
	handelRequest()
}

func handleCommand(cmd string) ReturnResponse {
	line := cmd
	if line == "exit" || line == "" {
		os.Exit(0)
	}
	response := ReturnResponse{}

	parts := strings.Split(line, " ")

	if len(parts) == 0 {
		os.Exit(0)
	}

	command := strings.TrimSpace(parts[0])
	args := parts[1:]

	output, err := exec.Command(command, args...).Output()

	response.err = err

	println(string(output))
	resp := string(output)
	formated := strings.Split(resp, "\n")
	response.response = formated

	return response
}
