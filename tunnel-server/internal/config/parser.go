package config

import (
	"net/url"
	"strings"
)

func ExtractSubdomainURL(publicURLBase *url.URL, subdomainKey string) string {
	return publicURLBase.Scheme + "://" + subdomainKey + "." + publicURLBase.Host

}

func ExtractAssignedSubdomain(host string) string {
	parts := strings.Split(host, ".")
	return parts[0]
}
