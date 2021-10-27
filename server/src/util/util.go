package util

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var (
	steamIDExp   = regexp.MustCompile("^s[0-9]{17,}")
	timeStampExp = regexp.MustCompile("t[0-9]{10,}$")
	tokenExp     = regexp.MustCompile("^s[0-9]{17,}t[0-9]{10,}$")
)

// ParseToken Parse sSTEAMIDtMASTERCOOKIE style GOTV+ Token. For example, if input is "s845489096165654t8799308478907", this should return "845489096165654" and "8799308478907".
func ParseToken(token string) (string, time.Time, error) {
	match := tokenExp.MatchString(token)
	if !match {
		return "", time.Time{}, fmt.Errorf("invalid token type")
	}
	steamid := steamIDExp.FindString(token)[1:]
	timestamp := timeStampExp.FindString(token)
	ts, err := strconv.ParseInt(timestamp[1:], 10, 64)
	if err != nil {
		return "", time.Time{}, err
	}
	t := time.Unix(ts, 0)

	return steamid, t, nil
}
