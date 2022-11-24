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
// define an object command to recive from json

type Command struct {
	Command string
}

// define an object for response

type ReturnResponse struct {
	response []string
	err      error
}

// controller that will be called at api/cmd

func CommandController(w http.ResponseWriter, r *http.Request) {
	// parse the bodi of request

	l := r.ContentLength
	body := make([]byte, l)
	r.Body.Read(body)

	// define an object to parse
	command := &Command{}

	// parse body to command object
	err := json.Unmarshal(body, command)
	if err != nil {
		fmt.Println(err)
	}

	// printing for debugging
	fmt.Println(command.Command)

	// handle command will be called 
	response := handleCommand(*&command.Command)

	// printing to get see response result
	fmt.Println("After response,", response)

	// if the out put of a command is not exist and have error
	if response.err != nil {
		fmt.Print(response.err)

		// set return type to json
		w.Header().Set("Content-Type", "application/json")

		// set the status bar to be 404
		w.WriteHeader(404)

		// parse the object to json
		output, err := json.MarshalIndent(response.err.Error(), "", "\t\t")

		// if result have an error
		if err != nil {
			fmt.Print(err)
		}

		// write the output
		_, _ = w.Write(output)

	// else if the command excuted well
	} else {

		// set the response to be json
		w.Header().Set("Content-Type", "application/json")

		// parse the object to json
		output, err := json.MarshalIndent(response.response, "", "\t\t")

		// if there is an error on parsing
		if err != nil {
			
			// w.Header().Set("Content-Type", "application/json")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		
		// return result 
		_, _ = w.Write(output)
	}
}

func handelRequest() {
	// when ever user redirect to /api/cmd this will control the request
	http.HandleFunc("/api/cmd", CommandController)

	// this will server the api at localhost:8001
	log.Fatal(http.ListenAndServe(":8001", nil))
}

func main() {
	handelRequest()
}

func handleCommand(cmd string) ReturnResponse {
	// assign command string to variable line
	line := cmd

	// if the command value is exit the program should exit
	if line == "exit" || line == "" {
		os.Exit(0)
	}

	// define an object for return
	response := ReturnResponse{}

	// split the command argument for command and args
	parts := strings.Split(line, " ")

	// take out the command from the command object, as the first string is alwasy the command parts[0] will be it
	command := strings.TrimSpace(parts[0])

	// take out the args
	args := parts[1:]

	// excute the command
	output, err := exec.Command(command, args...).Output()

	// if there is an error add that error to the returned object
	response.err = err

	// debug the output
	println(string(output))

	// convert the out put to string
	resp := string(output)

	// if the out put is in the form of array split it with new line
	formated := strings.Split(resp, "\n")

	// add the response out put to the returned object
	response.response = formated

	// return the response
	return response
}
