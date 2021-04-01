package handlers

import (
	"regexp"
)

// GOAL : MATCH [["Jl. SMA Aek Kota Batu","id"],["Sumatera Utara","de"]]
var stringRegex = "(\\p{L}| |\\d|\\_|\\-|\\,|\\.|/)"
var languageRegex = "\\[\"" + stringRegex + "+\"+,\"" + stringRegex + "{1,10}\"\\]"

var compiledRegexp *regexp.Regexp = regexp.MustCompile(languageRegex)

var googleMapsRegex = "https:\\/\\/(www\\.|maps\\.)?google\\.com/"
var compiledGoogleMapsRegex *regexp.Regexp = regexp.MustCompile(googleMapsRegex)

var googleConsentRegex = "https:\\/\\/consent\\.google\\.com/"
var compiledGoogleConsentRegex *regexp.Regexp = regexp.MustCompile(googleConsentRegex)

// filterStrings filters all string contents from a given string (as byte array),
// used to strip all localization information from a specific street view packet
func filterPhotometa(body []byte) []byte {
	result := compiledRegexp.ReplaceAllString(string(body), "[\"\",\"\"]")
	return []byte(result)
}

func filterUrls(body []byte) []byte {
	result := compiledGoogleMapsRegex.ReplaceAllString(string(body), "/")
	result = compiledGoogleConsentRegex.ReplaceAllString(result, "/")
	return []byte(result)
}
