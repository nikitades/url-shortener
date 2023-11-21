package main

import (
	"flag"

	"github.com/nikitades/url-shortener/internal/api"
)

func main() {
	port := flag.String("port", "8080", "http port")
	sqlconnstr := flag.String("sqlconn", "", "sql connection string")
	flag.Parse()
	api.Start(*port, *sqlconnstr)
}
