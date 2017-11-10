/*  This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>. 1 */


package main

import (
	"log"
	"os"
//	"fmt"
	"net/http"

	"github.com/dustin/go-coap"
	"github.com/gorilla/mux"
)

func getTemperature() []byte {
	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: 12345,
		Payload:   []byte("hello, world!"),
	}

	path := "/temperature"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

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

func temperatureHandler(w http.ResponseWriter, r *http.Request) {
	pl := getTemperature()
	log.Printf("Temperature: %s\n", string(pl))

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(string(pl)));
}

func main() {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", temperatureHandler)

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
