package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	flag.Usage = func() {
		os.Stderr.WriteString("mstdn-aggregator requires the path to a configuration file.\n\nIf the file does not exist, this program will walk you through creating it.\n\n")
		flag.PrintDefaults()
	}

	flagConfigPath := flag.String("c", "", "Configuration file path (suggested name: accountname.yml)")

	flag.Parse()
	if *flagConfigPath == "" {
		flag.Usage()
		os.Exit(2)
	}

	f, err := os.Open(*flagConfigPath)
	if os.IsNotExist(err) {
		setupWizard(*flagConfigPath)
		return
	}
	if err != nil {
		log.Fatalln("Failed to open config file:", err)
	}

	cfg, err := readConfig(f)
	if err != nil {
		_ = f.Close()
		log.Fatalln("Failed to read config file:", err)
	}
	if err = f.Close(); err != nil {
		log.Println("Warning: error closing config file:", err)
	}

	runBot(cfg)
}
