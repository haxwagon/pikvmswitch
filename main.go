package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/stianeikeland/go-rpio"
)

var (
	pinFlag  = flag.Int64("pin", 25, "the GPIO pin to toggle")
	portFlag = flag.Int64("port", 0, "server port.  If set, will run a service")
)

func togglePin() error {
	err := rpio.Open()
	if err != nil {
		return err
	}
	defer rpio.Close()

	switchPin := rpio.Pin(*pinFlag)
	switchPin.Toggle()

	return nil
}

func main() {
	flag.Parse()

	if *portFlag > 0 {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(405)
				fmt.Fprintf(w, "Only POST supported")
				return
			}

			if err := togglePin(); err != nil {
				w.WriteHeader(500)
				fmt.Fprintf(w, "Failed to toggle, err=%v", err)
			}
		})
		fmt.Printf("Starting service at port %d...\n", *portFlag)
		http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), nil)
	} else {
		if err := togglePin(); err != nil {
			panic(fmt.Sprint("unable to open gpio: ", err.Error()))
		}
	}
}
