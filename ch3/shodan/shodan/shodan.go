package shodan

const BaseURL = "https://api.shodan.io"

type Client struct {
	apiKey string
	url    string
}

func New(apiKey string) *Client {
	return &Client{apiKey: apiKey, url: BaseURL}
}
