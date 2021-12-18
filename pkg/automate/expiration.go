package automate

import (
	"log"
	"os"
	"time"
)

const (
	Layout                      = "Mon, 02 Jan 2006 15:04:05 -0700"
	ExpirationNameEnvVar        = "EXPIRATION_TAG_NAME"
	SlidingExpirationNameEnvVar = "SLIDING_EXPIRATION"
)

func parseExpiration(v string) (time.Time, error) {
	return time.Parse(Layout, v)
}

func getSlidingExpiration() string {
	se, ok := os.LookupEnv(SlidingExpirationNameEnvVar)
	if !ok {
		se = "168h"
	}
	return se
}

func isExpired(tags map[string]*string) bool {
	tv, ok := tags[getExpirationTagName()]
	if !ok {
		return false
	}
	expiresAt, err := parseExpiration(*tv)
	if err != nil {
		log.Printf("Error while parsing '%s' expected time.Time string", *tv)
		return false
	}
	return time.Now().After(expiresAt)
}

func getExpiration() string {
	t := calculateExpiration()
	return t.Format(Layout)
}

func getExpirationTagName() string {
	e, ok := os.LookupEnv(ExpirationNameEnvVar)
	if !ok {
		e = "com.thorsten-hans.expiration"
	}
	return e
}

func calculateExpiration() time.Time {
	d, _ := time.ParseDuration(getSlidingExpiration())
	expiresAt := time.Now().Add(d)
	return expiresAt
}
