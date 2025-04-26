package validators

import (
	"net/url"
	"regexp"
	"strings"
)

var validSlugRegexp = regexp.MustCompile(
	`^(?:` +
		`[A-Za-z0-9]{8,12}` +
		`|` +
		`[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}` +
		`)$`,
)

func Slug(s string) error {
	valid := validSlugRegexp.MatchString(s)
	if !valid {
		return ErrInvalidSlug
	}

	return nil
}

func OriginalURL(u string) error {
	parsed, err := url.Parse(u)
	if err != nil {
		return ErrInvalidOriginalURL
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ErrInvalidOriginalURL
	}

	if parsed.Host == "" || strings.IndexRune(parsed.Host, '\\') > 0 {
		return ErrInvalidOriginalURL
	}

	return nil
}
