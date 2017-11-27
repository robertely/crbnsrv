package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	Host = "localhost"
	Port = "2003"
)

func main() {
	// Listen for incoming connections.
	l, err := net.Listen("tcp", Host+":"+Port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + Host + ":" + Port)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

type Datum struct {
	Key   string
	Value interface{}
	Time  time.Time
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Close connection when this function ends
	defer func() {
		fmt.Println("Closing connection...")
		conn.Close()
	}()

	fmt.Println("Handling new connection...")
	timeoutDuration := 10 * time.Second
	bufReader := bufio.NewReader(conn)

	for {
		// Set a deadline for reading. Read operation will fail if no data
		// is received after deadline.
		conn.SetReadDeadline(time.Now().Add(timeoutDuration))

		// Read tokens delimited by newline
		lineraw, err := bufReader.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		// Process line
		Cleaned := strings.Trim(string(lineraw), "\n")
		s := strings.Split(Cleaned, " ")
		d := Datum{}
		d.Key = strings.Trim(s[0], ".")
		d.Value, _ = strconv.ParseFloat(s[1], 64)
		if len(s) == 3 {
			t, err := strconv.ParseInt(s[2], 10, 64)
			if err == nil {
				d.Time = time.Unix(t, 0)
			}
		}
		fmt.Println(d)
	}
}
