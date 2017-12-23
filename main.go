package main

import (
	"bufio"
	"github.com/BurntSushi/toml"
	"github.com/Juev/test-tcp/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

// Config type
type Config struct {
	Host string
	Port int
}

var config Config
var port string
var host string

var (
	cfg     = kingpin.Flag("config", "Config file in TOML format. Used only ip and port variables.").Short('c').Default("config.toml").String()
	logFile = kingpin.Flag("log", "Log file name.").Short('l').Default("logfile").String()
	mode    = kingpin.Flag("mode", "Client/Server mode. Example: -m \"server\"").Short('m').Default("server").String()
)

func main() {
	kingpin.Version(version.Release)
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.Parse()

	f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Cannot open Log File")
	}
	defer f.Close()

	log.SetOutput(f)

	if *mode != "server" && *mode != "client" {
		log.Fatal("Mode can be only \"server\" or \"client\"")
	}

	log.Println("Using config file: ", *cfg)
	if _, err := toml.DecodeFile(*cfg, &config); err != nil {
		log.Fatal(err)
	}

	if config.Host == "" || config.Port == 0 {
		log.Fatal("Host and Port must be present in config file, and not be \"\" and 0 value")
	}

	port = strconv.Itoa(config.Port)
	host = config.Host

	if *mode == "server" {
		log.Printf(
			"Starting the server, build time: %s, release: %s",
			version.BuildTime, version.Release,
		)
		server()
	} else {
		log.Printf(
			"Starting the client, build time: %s, release: %s",
			version.BuildTime, version.Release,
		)
		client()
	}
}

func client() {
	for {
		time.Sleep(1 * time.Second)
		conn, err := net.Dial("tcp", host+":"+port)
		if err != nil {
			log.Println("Something wrong: ", err)
			continue
		}
		defer func() {
			conn.Close()
			log.Println("Connection closed")
		}()
		// send to socket
		conn.Write([]byte("Ping\n"))
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		log.Print("Message from server: " + message)
	}
}

func server() {
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listener created on " + host + ":" + port)

	defer func() {
		listener.Close()
		log.Println("Listener closed")
	}()

	for {
		// Get net.TCPConn object
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	log.Println("Handling new connection...")

	// Close connection when this function ends
	defer func() {
		log.Println("Closing connection...")
		conn.Close()
	}()

	timeoutDuration := 5 * time.Second
	bufReader := bufio.NewReader(conn)

	// Set a deadline for reading. Read operation will fail if no data
	// is received after deadline.
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	// Read tokens delimited by newline
	bytes, err := bufReader.ReadBytes('\n')
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("%s", bytes)
	conn.Write([]byte("Pong"))
}
