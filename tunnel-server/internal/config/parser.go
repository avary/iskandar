package config

import (
	"errors"
	"net/url"
	"strings"
)

var ErrInvalidSubdomain = errors.New("invalid subdomain: host must contain at least one subdomain part")

func ExtractSubdomainURL(publicURLBase *url.URL, subdomainKey string) string {
	return publicURLBase.Scheme + "://" + subdomainKey + "." + publicURLBase.Host

}

func ExtractAssignedSubdomain(host string) (string, error) {
	// Remove port if present
	if colonIdx := strings.Index(host, ":"); colonIdx != -1 {
		host = host[:colonIdx]
	}

	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return "", ErrInvalidSubdomain
	}
	return parts[0], nil
}
