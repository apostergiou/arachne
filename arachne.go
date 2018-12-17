package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/urfave/cli"
)

func main() {
	var (
		listen   string
		upstream string
		network  string
	)

	cli.AppHelpTemplate = fmt.Sprintf(`%s
WEBSITE: https://apostergiou.com

EXAMPLES:
	1. Listen on tcp/53 and forward to Cloudflare
		$ {{.Name}}

	2. Listen on tcp/5300 and forward to Cloudflare server 1.1.0.1:853
		$ {{.Name}} -u 1.1.0.1:853 -l 0.0.0.0:5300
`, cli.AppHelpTemplate)

	app := cli.NewApp()
	app.Usage = "Receive DNS queries and forward them upstream using TLS"
	app.Name = "arachne"
	app.Version = "1.0.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Apostolis Stergiou",
			Email: "apostergiou.com",
		},
	}

	app.Copyright = "(c) 2018 Apostolis Stergiou"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "upstream, u",
			Value:       "1.1.1.1:853",
			Usage:       "upstream DNS",
			EnvVar:      "UPSTREAM",
			Destination: &upstream,
		},
		cli.StringFlag{
			Name:        "listen, l",
			Value:       "127.0.0.1:53",
			Usage:       "address to listen for requests",
			EnvVar:      "LISTEN",
			Destination: &listen,
		},
		cli.StringFlag{
			Name:        "network, n",
			Value:       "tcp",
			Usage:       "network to listen to",
			EnvVar:      "NETWORK",
			Destination: &network,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	app.Action = func(c *cli.Context) error {
		log.Printf("Starting arachne")
		cfg, err := SetupConfig(listen, upstream, network)
		if err != nil {
			return err
		}
		return StartServer(cfg)
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// StartServer starts a new server that listens for Network (e.g. TCP) DNS queries.
func StartServer(cfg *Config) error {
	var wg sync.WaitGroup
	s, err := NewServer(cfg, log.New(os.Stderr, fmt.Sprintf("[server] "), log.LstdFlags))
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()
	wg.Wait()

	return nil
}
