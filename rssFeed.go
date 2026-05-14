package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
)

// Struct defining xml unmarshalling
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// Nested xml struct
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// Function to fetch a feed from a given URL, returns RSSFeed struct
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("http request failed with error: %w", err)
	}

	//Set user agent
	req.Header.Set("User-Agent", "gogetterRSS")
	//Create http client
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http client failed with error: %w", err)
	}
	defer res.Body.Close()
	//Create RSSFeed struct
	var xmlRSS RSSFeed
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http failed: %w", err)
	}
	xmlBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body with error: %v\n", err)
	}
	//Unmarshal XML to created struct above
	err = xml.Unmarshal(xmlBody, &xmlRSS)
	if err != nil {
		return nil, errors.New("can't unmarshal XML to xmlRSS struct\n")
	}
	//Run certain struct fields through html's unescape string method
	xmlRSS.Channel.Title = html.UnescapeString(xmlRSS.Channel.Title)
	xmlRSS.Channel.Description = html.UnescapeString(xmlRSS.Channel.Description)
	//Loop through nested RSSItem struct
	for i := range xmlRSS.Channel.Item {
		xmlRSS.Channel.Item[i].Title = html.UnescapeString(xmlRSS.Channel.Item[i].Title)
		xmlRSS.Channel.Item[i].Description = html.UnescapeString(xmlRSS.Channel.Item[i].Description)
	}
	return &xmlRSS, nil
}
