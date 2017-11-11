/*
Go CoAP/HTTP Proxy
Copyright (C) 2017 Simon WIETRICH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>. 1 
*/


package main

import (
	"log"
	"net/http"

	"github.com/dustin/go-coap"
	"github.com/gorilla/mux"
)

func getResource(resource string) []byte {
	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: 12345,
		Payload:   []byte ("Hello, world"),
	}

	path := "/" + resource

	req.SetOption(coap.ETag, "weetag")
	req.SetOption(coap.MaxAge, 3)
	req.SetPathString(path)

	c, err := coap.Dial("udp", "10.0.0.100:5683")
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	rv, err := c.Send(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	return rv.Payload
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("%s %s\n", r.Method, r.URL)
	pl := getResource(vars["resource"])


	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(pl)));
	log.Printf("\t-> Respond %s\n", string(pl))
}

func main() {
	// declare new multiplexer
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/proxy/{resource}", proxyHandler)

	// Pass multiplexer to go stdlib http server
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
