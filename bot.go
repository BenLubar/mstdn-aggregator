package main

import (
	"context"
	"log"

	"github.com/mattn/go-mastodon"
)

func runBot(cfg *Config) {
	ctx := context.Background()
	client := createClient(cfg)

	ch, err := client.NewWSClient().StreamingWSList(ctx, cfg.List)
	if err != nil {
		log.Fatalln("Failed to open list timeline stream! Error:", err)
	}

	go backfill(ctx, client, cfg.List)

	for event := range ch {
		switch e := event.(type) {
		case *mastodon.UpdateEvent:
			boost(ctx, client, e.Status)
		case *mastodon.ErrorEvent:
			log.Println("streaming error:", e)
		}
	}
}

func backfill(ctx context.Context, client *mastodon.Client, id mastodon.ID) {
	var pg mastodon.Pagination
	for {
		statuses, err := client.GetTimelineList(ctx, id, &pg)
		if err != nil {
			log.Println("stopping backfill: encountered error:", err)
			return
		}

		for _, status := range statuses {
			if status.Reblogged == true {
				log.Println("stopping backfill: reached an already-boosted post")
				return
			}

			if status.Reblog != nil && status.Reblog.Reblogged == true {
				// already boosted; don't re-boost on backfill
				continue
			}

			boost(ctx, client, status)
		}

		if pg.MaxID == "" {
			log.Println("stopping backfill: reached end of list feed")
			return
		}
	}
}

func boost(ctx context.Context, client *mastodon.Client, status *mastodon.Status) {
	if status.Reblog != nil {
		status = status.Reblog
		if status.Reblogged == true {
			if _, err := client.Unreblog(ctx, status.ID); err != nil {
				log.Println("failed to re-boost:", status.ID, "Error:", err)
				return
			}
		}
	}

	_, err := client.Reblog(ctx, status.ID)
	if err != nil {
		log.Println("failed to boost:", status.ID, "Error:", err)
	}
}
