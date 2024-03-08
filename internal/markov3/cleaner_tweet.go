package markov3

import (
	"regexp"
	"strings"
)

var tweetUserMatch = regexp.MustCompile("(@[^:. ]+)")
var tweetHashMatch = regexp.MustCompile("(#[^:. ]+)")
var tweetURLMatch = regexp.MustCompile("(https?://[^ ]+)")

// CleanedTweet represents a simplified tweet
type CleanedTweet struct {
	Original string
	Cleaned  string
	Mentions []string
	Hashtags []string
	URLs     []string
}

// CleanTweet takes the original tweet and cleans it
// if captureSubjects is true it will provide the subjects (Mentions,
// Hashtags, URLs) as well
func CleanTweet(orig string, captureSubjects bool) CleanedTweet {
	clean := orig
	// If you want to reply but have it show in
	// your timeline easily, a prefixed '.' does this
	if clean[0] == '.' {
		clean = clean[1:]
	}

	clean = strings.TrimSpace(clean)
	// Early RT models on Twitter used RT as a prefix
	// so we're dropping that too.
	if clean[0:2] == "RT" {
		clean = clean[3:]
	}

	clean = strings.ToLower(clean)

	var (
		mentions = make([]string, 0)
		hashes   = make([]string, 0)
		urls     = make([]string, 0)
	)

	if captureSubjects {
		mentionMatches := tweetUserMatch.FindAllStringSubmatch(orig, -1)
		for _, match := range mentionMatches {
			mentions = append(mentions, match[0])
		}

		hashMatches := tweetHashMatch.FindAllStringSubmatch(orig, -1)
		for _, hash := range hashMatches {
			hashes = append(hashes, hash[0])
		}

		urlMatches := tweetURLMatch.FindAllStringSubmatch(orig, -1)
		for _, urlm := range urlMatches {
			urls = append(urls, urlm[0])
		}
	}

	clean = tweetUserMatch.ReplaceAllString(clean, "@user")
	clean = tweetHashMatch.ReplaceAllString(clean, "#hash")
	clean = tweetURLMatch.ReplaceAllString(clean, "http://")

	for strings.Contains(clean, "  ") {
		clean = strings.ReplaceAll(clean, "  ", " ")
	}

	return CleanedTweet{
		Original: orig,
		Cleaned:  clean,
		Mentions: mentions,
		Hashtags: hashes,
		URLs:     urls,
	}
}
