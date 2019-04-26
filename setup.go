package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/mattn/go-mastodon"
)

func setupWizard(filename string) {
	ctx := context.Background()
	cfg := &Config{}
	in := bufio.NewScanner(os.Stdin)
	readLine := func() string {
		if !in.Scan() {
			err := in.Err()
			if err == nil {
				err = io.ErrUnexpectedEOF
			}
			log.Fatalln("Error reading input: ", err)
		}

		return in.Text()
	}

	fmt.Println("Welcome to the mstdn-aggregator setup wizard.")
	fmt.Println()
	fmt.Println("There are only a few questions you need to answer to set up this bot.")
	fmt.Println()
	fmt.Println("What is your server name?")
	fmt.Println("Something like https://mastodon.example (no slash after the instance name)")
	fmt.Print("Server name: ")
	cfg.Secrets.Server = readLine()

	fmt.Println()
	fmt.Println("Attempting to register app...")
	const scopes = "read write:statuses write:lists"
	app, err := mastodon.RegisterApp(ctx, &mastodon.AppConfig{
		Server:       cfg.Secrets.Server,
		ClientName:   "Mastodon aggregator bot",
		RedirectURIs: "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       scopes,
		Website:      "https://github.com/BenLubar/mstdn-aggregator",
	})
	if err != nil {
		log.Fatalln("Failed to register! Is your instance down? The server address should be the one you log in on rather than the one your address shows if they are different. Error:", err)
	}

	cfg.Secrets.ClientID = app.ClientID
	cfg.Secrets.ClientSecret = app.ClientSecret

	fmt.Println("Registered an app with ID ", app.ID)
	fmt.Println()
	fmt.Println("The next step is visiting this address to connect your bot's account with the bot. Follow the instructions and paste the code it gives you here when you're done.")
	fmt.Println()
	fmt.Println()
	fmt.Println(app.AuthURI)
	fmt.Println()
	fmt.Println()
	fmt.Print("Code: ")
	code := readLine()
	client := createClient(cfg)

	if err = client.AuthenticateToken(ctx, code, app.RedirectURI); err != nil {
		log.Fatalln("Failed to validate authorization code! Error: ", err)
	}

	account, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		log.Fatalln("Could not retrieve account information! Error:", err)
	}

	fmt.Println()
	fmt.Print("Successfully authenticated as @", account.Acct, "!\n")
	if !account.Bot {
		fmt.Println()
		fmt.Println("Warning: This account is not labelled as a bot. If you are not monitoring this account's notifications, it is recommended to mark it as a bot, even if it represents a real person.")
	}

	var list *mastodon.List
	lists, err := client.GetLists(ctx)
	for _, l := range lists {
		if l.Title == "mstdn-aggregator" {
			list = l
			break
		}
	}
	if err == nil && list == nil {
		list, err = client.CreateList(ctx, "mstdn-aggregator")
	}
	if err != nil {
		log.Println("Failed to create a list. Create one manually and add its ID to the config to fix this. Error:", err)
	} else {
		cfg.List = list.ID
		timeline := cfg.Secrets.Server + "/web/timelines/list/" + url.PathEscape(string(list.ID))

		fmt.Println()
		fmt.Println("Add accounts to the mstdn-aggregator list at", timeline, "and the bot will boost them.")
	}

	fail := func(msg string, err error) {
		fmt.Println()
		fmt.Println("The configuration file (you must save it manually):")
		fmt.Println()
		writeConfig(os.Stdout, cfg)
		fmt.Println()
		log.Fatalln(msg, err)
	}

	f, err := os.Create(filename)
	if err != nil {
		fail("Failed to create configuration file! Error:", err)
	}
	defer f.Close()

	if err = writeConfig(f, cfg); err != nil {
		fail("Failed to write configuration file! Error:", err)
	}

	fmt.Println("Configuration file successfully written to", filename)
	fmt.Println("Run the bot again to start it.")
}
