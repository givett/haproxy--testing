package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
)

// Httpstruct comment
type Httpstruct struct {
	Httpcode string `json:"httpcode"`
}

var httpstruct Httpstruct
var codenum int

// loadConfig comment
func loadConfig() {
	rBody, err := ioutil.ReadFile("availcode.json")

	if err != nil {
		fmt.Printf("Error reading availcode.json, %s", rBody)
	}

	err = json.Unmarshal(rBody, &httpstruct)
	if err != nil {
		fmt.Printf("Error unmarshaling availcode.json, %s", err)
	}

	fmt.Printf("httpcode is %s\n", httpstruct.Httpcode)
}

// handler comment
func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "GET" {
		codenum, err := strconv.Atoi(httpstruct.Httpcode)
		if err != nil {
			panic(err)
		}

		// w.WriteHeader(http.StatusOK)
		w.WriteHeader(codenum)

		host, _ := os.Hostname()
		fmt.Fprintf(w, "os.Hostname: %s\n", host)

		addrs, err := net.LookupIP(host)
		if err != nil {
			//panic(err)
			fmt.Fprintf(w, "error: %s\n", err)
		}

		fmt.Fprintf(w, "addr: %s\n", addrs)

		for _, addr := range addrs {

			if ipv4 := addr.To4(); ipv4 != nil {
				fmt.Fprintf(w, "IPv4: %s\n", ipv4)
			}
		}
		println(r.URL.String() + " " + r.Method + httpstruct.Httpcode)
	} else if r.Method == "POST" {
		fmt.Fprintf(w, "POST, %q", html.EscapeString(r.URL.Path))

		rBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "Error reading app request, %v", rBody)
		}

		err = json.Unmarshal(rBody, &httpstruct)
		if err != nil {
			fmt.Fprintf(w, "Error unmarshaling app %v", err)
		}

		fmt.Fprintf(w, "httpcode is %s", httpstruct.Httpcode)

	} else {
		println(r.URL.String() + " " + r.Method + "405")
		http.Error(w, "Invalid request method.", 405)
	}

}

func main() {
	loadConfig()

	http.HandleFunc("/", handler)
	fmt.Print("Listening on http://localhost:8070")
	http.ListenAndServe(":8070", nil)

	// curl -X GET localhost:80 -vv
}
