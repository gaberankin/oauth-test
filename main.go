package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/skratchdot/open-golang/open"
)

var srv *http.Server
var listener net.Listener
var wg sync.WaitGroup

func startServer() {
	srv = &http.Server{Addr: ":8081"}
	srv.ListenAndServe()
}
func stopServer() {
	go func() {
		time.Sleep(1 * time.Second)
		wg.Done()
		srv.Shutdown(nil)
	}()
}
func main() {

	http.HandleFunc("/complete", func(w http.ResponseWriter, r *http.Request) {
		defer stopServer() // this will stop the server in, yet again, another separate thread.  probably a bad idea
		_, err := w.Write([]byte("OAUTH Authentication complete.  You can now close your browser and return to the terminal."))
		if err != nil {
			fmt.Printf(err.Error())
		}
	})
	wg.Add(1)
	//have to start the server in a seperate thread to be able to tell our browser when to open the oauth authentication flow
	go startServer()

	//opening this url in your browser will immediately kill the app.  in reality, you want to point this url to the
	// oauth applicaiton login page, and in your oauth service, register this localhost url as the callback.
	open.Run("http://localhost:8081/complete")
	wg.Wait()
	//place remainder of application logic here.
}
