package handlers

import (
	"regexp"
)

var stringRegex = "(\\p{L}| |\\d|\\_|\\-|\\,|\\.|/)"

// GOAL : MATCH [["Jl. SMA Aek Kota Batu","id"],["Sumatera Utara","de"]]
var listOfStringsRegex = "\\[\"" + stringRegex + "+\"+,\"" + stringRegex + "{1,10}\"\\]"
var replaceListOfStringsWith = "[\"\",\"\"]"
var compiledListOfStringsRegex *regexp.Regexp = regexp.MustCompile(listOfStringsRegex)

var iconRegex = "\"https://maps.gstatic.com/mapfiles/annotations/icons/" + stringRegex + "+\""
var replaceIconWith = "\"\""
var compiledIconRegex *regexp.Regexp = regexp.MustCompile(iconRegex)

var googleMapsRegex = "https:\\/\\/(www\\.|maps\\.)?google\\.com/"
var compiledGoogleMapsRegex *regexp.Regexp = regexp.MustCompile(googleMapsRegex)

var googleConsentRegex = "https:\\/\\/consent\\.google\\.com/"
var compiledGoogleConsentRegex *regexp.Regexp = regexp.MustCompile(googleConsentRegex)

// filterStrings filters all string contents from a given string (as byte array),
// used to strip all localization information from a specific street view packet
func filterPhotometa(body string) string {
	result := compiledListOfStringsRegex.ReplaceAllString(body, replaceListOfStringsWith)
	result = compiledIconRegex.ReplaceAllString(result, replaceIconWith)
	return result
}

func filterUrls(body string) string {
	result := compiledGoogleMapsRegex.ReplaceAllString(body, "/")
	result = compiledGoogleConsentRegex.ReplaceAllString(result, "/")
	return result
}
