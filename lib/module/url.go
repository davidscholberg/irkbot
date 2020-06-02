package module

import (
	"context"
	"fmt"
	"github.com/davidscholberg/irkbot/lib/configure"
	"github.com/davidscholberg/irkbot/lib/message"
	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/net/html"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io"
	"mvdan.cc/xurls"
	"net"
	"net/http"
	urllib "net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// parseUrls attempts to fetch the title of the HTML document returned by a URL
func parseUrls(cfg *configure.Config, in *message.InboundMsg, actions *actions) bool {
	// disallow url parsing in PMs
	if !strings.HasPrefix(in.Src, "#") {
		return false
	}

	urls := xurls.Strict().FindAllString(in.Msg, -1)

	for _, urlStr := range urls {
		url, err := urllib.Parse(urlStr)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		if v, err := validateUrl(url); !v {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		title := ""
		host := ""
		if twitterConfigured(cfg) && isTweet(url) {
			host = url.Host
			var err error
			title, err = getTwitterTitle(cfg, url)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		} else {
			response, err := actions.httpGet(urlStr)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			host = response.Request.URL.Host
			title, err = getHtmlTitle(response, cfg.Http.ResponseSizeLimit)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		}
		cleanedTitle := cleanTitleWhiteSpace(title)
		actions.say(fmt.Sprintf("^ %s - [%s]", cleanedTitle, host))
	}

	// don't consume the message, in case there are commands in it
	return false
}

// cleanTitleWhiteSpace takes the contents of an html title and makes it fit
// to be printed on a single line.
func cleanTitleWhiteSpace(title string) string {
	// split string by newline characters
	titleFields := strings.FieldsFunc(title, func(c rune) bool {
		return c == '\n' || c == '\r'
	})
	// trim each field
	trimmedFields := []string{}
	for _, field := range titleFields {
		field = strings.TrimSpace(field)
		if field != "" {
			trimmedFields = append(trimmedFields, field)
		}
	}
	// join by newline separator
	return strings.Join(trimmedFields, " / ")
}

// validateUrl ensures that the given URL is safe to GET.
func validateUrl(url *urllib.URL) (bool, error) {
	privateSubnets := []string{
		"127.0.0.0/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"::1/128",
		"fc00::/7"}

	ips, err := net.LookupIP(url.Host)
	if err != nil {
		return false, err
	}
	if len(ips) == 0 {
		return false, fmt.Errorf("no IPs found for %s", url.Host)
	}

	hostNotAllowed, err := isCidrMatch(&ips[0], privateSubnets)
	if err != nil {
		return false, err
	}
	if hostNotAllowed {
		return false, fmt.Errorf("host not allowed")
	}

	return true, nil
}

// isCidrMatch tests if a given IP is in any of a given set of CIDR subnets.
func isCidrMatch(ip *net.IP, subnets []string) (bool, error) {
	var cidr *net.IPNet
	var err error

	for _, subnet := range subnets {
		_, cidr, err = net.ParseCIDR(subnet)
		if err != nil {
			return false, err
		}
		if cidr.Contains(*ip) {
			return true, nil
		}
	}

	return false, nil
}

// isTweet determines if the URL is a tweet (i.e. a twitter status)
func isTweet(url *urllib.URL) bool {
	if url.Hostname() != "twitter.com" {
		return false
	}
	// tokenize path (tweets should be in the format /:user/status/:id)
	pathElements := strings.Split(url.EscapedPath(), "/")
	if len(pathElements) != 4 {
		return false
	}
	if pathElements[2] == "status" {
		return true
	}
	return false
}

// twitterConfigured makes sure that the necessary twitter config is in place.
func twitterConfigured(cfg *configure.Config) bool {
	return cfg.Modules["url"]["twitter_client_id"] != "" && cfg.Modules["url"]["twitter_client_secret"] != ""
}

// getTwitterTitle takes the URL object and returns the title of the twitter
// status.
func getTwitterTitle(cfg *configure.Config, url *urllib.URL) (string, error) {
	if !isTweet(url) {
		return "", fmt.Errorf("URL is not a tweet")
	}
	if !twitterConfigured(cfg) {
		return "", fmt.Errorf("twitter parameters are not in config")
	}
	// tokenize path (tweets should be in the format /:user/status/:id)
	pathElements := strings.Split(url.EscapedPath(), "/")
	statusIDStr := pathElements[3]
	statusID, err := strconv.Atoi(statusIDStr)
	if err != nil {
		return "", err
	}
	oauth2Config := &clientcredentials.Config{
		ClientID:     cfg.Modules["url"]["twitter_client_id"],
		ClientSecret: cfg.Modules["url"]["twitter_client_secret"],
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	// configure http timeout for auth call
	httpClientWithTimeout := &http.Client{
		Timeout: time.Duration(cfg.Http.Timeout) * time.Second,
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClientWithTimeout)
	httpClient := oauth2Config.Client(ctx)
	// configure http timeout for api call
	httpClient.Timeout = time.Duration(cfg.Http.Timeout) * time.Second
	twitterClient := twitter.NewClient(httpClient)
	statusShowParams := &twitter.StatusShowParams{
		TweetMode: "extended",
	}
	tweet, _, err := twitterClient.Statuses.Show(int64(statusID), statusShowParams)
	if err != nil {
		return "", err
	}
	if tweet != nil {
		if tweet.User == nil {
			return tweet.FullText, nil
		}
		return fmt.Sprintf("%s: \"%s\"", tweet.User.Name, tweet.FullText), nil
	}
	return "", err
}

// getHtmlTitle returns the HTML title found in the given response body.
// This function closes the response body.
func getHtmlTitle(response *http.Response, responseSizeLimit int64) (string, error) {
	defer response.Body.Close()

	// ignore response codes 400 and above
	if response.StatusCode >= 400 {
		return "", fmt.Errorf("received status %d", response.StatusCode)
	}

	doctree, err := html.Parse(io.LimitReader(response.Body, responseSizeLimit))
	if err != nil {
		return "", err
	}

	title, err := searchForHtmlTitle(doctree)
	if err != nil {
		return "", err
	}
	if len(title) == 0 {
		return "", fmt.Errorf("title not found")
	}

	return title, nil
}

// searchForHtmlTitle searches the parsed html document for the title.
func searchForHtmlTitle(n *html.Node) (string, error) {
	if n.Type == html.ElementNode && n.Data == "title" {
		if n.FirstChild == nil {
			err := fmt.Errorf("title node has no child")
			return "", err
		}
		if n.FirstChild.Type != html.TextNode {
			err := fmt.Errorf("child of title not TextNode type")
			return "", err
		}
		return n.FirstChild.Data, nil
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title, err := searchForHtmlTitle(c)
		if len(title) > 0 || err != nil {
			return title, err
		}
	}
	return "", nil
}
