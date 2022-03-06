package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	config, err := parseFlags()
	if err != nil {
		panic(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addr := r.RemoteAddr
		uri := r.RequestURI
		if config.Verbose {
			log.Printf("Request URI: '%s' from '%s'\n", uri, addr)
		}
		f := http.FileServer(http.Dir(config.MountDir))
		f.ServeHTTP(w, r)
	})
	log.Printf("Start to serve. Port: '%v', MountDir: '%s'\n", config.Port, config.MountDir)
	err = http.ListenAndServe(":"+strconv.Itoa(config.Port), handler)
	if err != nil {
		log.Fatal("Could not serve: %s", err)
	}

}

func parseFlags() (Config, error) {
	var dir string
	var port int
	var isVerbose bool
	const defaultPort = 3000
	flag.IntVar(&port, "p", defaultPort, "port to listen. if not present, use 3000.")
	flag.IntVar(&port, "port", defaultPort, "port to listen. if not present, use 3000.")
	flag.BoolVar(&isVerbose, "v", false, "print request log")
	flag.BoolVar(&isVerbose, "verbose", false, "print request log")
	flag.Parse()
	args := flag.Args()

	switch l := len(args); l {
	case 0:
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalln("Could not get current dir.")
			return Config{}, err
		}
		dir = wd
	default:
		wdir, err := filepath.Abs(args[0])
		if err != nil {
			log.Fatalln("Could not get current dir. Use user home dir instead.")
		}

		dir = wdir
	}
	// if (err = filepath.Abs(dir); err != nil) {
	// }

	config := Config{
		Port:     port,
		MountDir: dir,
		Verbose:  isVerbose,
	}
	return config, nil
}

type Config struct {
	Port     int
	MountDir string
	Verbose  bool
}
