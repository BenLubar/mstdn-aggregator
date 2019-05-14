module github.com/BenLubar/mstdn-aggregator

go 1.12

require (
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/mattn/go-mastodon v0.0.4
	github.com/tomnomnom/linkheader v0.0.0-20180905144013-02ca5825eb80 // indirect
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/mattn/go-mastodon => github.com/BenLubar-PR/go-mastodon v0.0.3-0.20190427015647-5ff4ecfd52d2
