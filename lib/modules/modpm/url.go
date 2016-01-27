package modpm

import (
    "fmt"
    "golang.org/x/net/html"
    "net/http"
    "github.com/davidscholberg/irkbot/lib"
    "github.com/mvdan/xurls"
)

// Url attempts to fetch the title of the HTML document returned by a URL
func Url(p *lib.Privmsg) bool {
    urls := xurls.Strict.FindAllString(p.Msg, -1)

    for _, url := range urls {
        title, err := getHtmlTitle(url)
        if err != nil {
            continue
        }
        lib.Say(p, fmt.Sprintf("^ %s - [%s]", title, url))
    }

    // don't consume the message, in case there are commands in it
    return false
}

// getHtmlTitle returns the HTML title found at the given URL.
func getHtmlTitle(url string) (string, error) {
    response, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer response.Body.Close()

    doctree, err := html.Parse(response.Body)
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
