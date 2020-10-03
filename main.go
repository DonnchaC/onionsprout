package main

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"strconv"
	"net"
	"net/http"
	"os"
	"golang.org/x/crypto/acme"

	"github.com/yawning/bulb"
)

func tor_control() (*bulb.Conn, error) {
	// Create a connection to the Tor control port
	conn, err := bulb.Dial("unix", "/var/run/tor/control")
	if err != nil {
		log.Fatalf("failed to connect to control port: %v", err)
		return nil, err
	}

	// See what's really going on under the hood.
	// Do not enable in production.
	conn.Debug(true)

	// Authenticate with the control port.  The password argument
	// here can be "" if no password is set (CookieAuth, no auth).
	if err := conn.Authenticate(""); err != nil {
		log.Fatalf("Authentication failed: %v", err)
		return nil, err
	}
	return conn, nil
}

func create_onion(c *bulb.Conn, port uint16, target_port int) (*bulb.OnionInfo, error) {
	targetPortStr := strconv.FormatUint((uint64)(target_port), 10)
	cfg := &bulb.NewOnionConfig{
		DiscardPK: true,
		PortSpecs: []bulb.OnionPortSpec{
				bulb.OnionPortSpec{
				VirtPort: port,
				Target:   targetPortStr,
			},
		},
		// Comment to generate a new onion
		// PrivateKey: &bulb.OnionPrivateKey{
		// 	KeyType: "ED25519-V3",
		// 	Key: "QJCGBFk/KY5+4wjmGWKRvA3G/1QdrftXMf/JbZar7EyP1JGE7RNNie/FblDNe5q3GiB+C4xwg61MeUMUbmPB4w==",
		// },
	}

	// Create the onion.
	return c.NewOnion(cfg)
}


func main() {
	// Connect to Tor control port
	tor_conn, err := tor_control()
	if err != nil {
		log.Fatalf("Failed to connect to control port: %v", err)
	}
	defer tor_conn.Close()

	if os.Getenv("DOMAIN") == "" {
		log.Fatalf("The DOMAIN environment variable must be set.")
	}

	// Set up the web server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=15768000 ; includeSubDomains")
		fmt.Fprintf(w, "Hello, HTTPS world!")
	})


	// Bind a random port to host the web server
	tls_listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("Error binding an open port: %s", err)
	}


    certManager := autocert.Manager{
        Prompt:     autocert.AcceptTOS,
        HostPolicy: autocert.HostWhitelist(os.Getenv("DOMAIN")),
        Cache:      autocert.DirCache("cache-path"),
    }

	// // create the server itself
	server := &http.Server{
		Addr: tls_listener.Addr().String(),
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
						NextProtos:     []string{acme.ALPNProto},
		},
	}

	// Start onion
	var onion_port uint16 = 443
	var listener_port = tls_listener.Addr().(*net.TCPAddr).Port
	onion, err := create_onion(tor_conn, onion_port, listener_port)
	if err != nil {
		log.Fatalf("Error starting onion: %s", err)
	}
	defer tor_conn.DeleteOnion(onion.OnionID)

	// serve HTTPS!
	log.Printf("Serving https for domains %+v at onion %s.onion", os.Getenv("DOMAIN"), onion.OnionID)
	log.Fatal(server.ServeTLS(tls_listener, "", ""))
}
