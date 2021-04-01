package handlers

import (
	"testing"
)

func TestFilterStringsPhotometa(t *testing.T) {
	inputs := []string{
		"This is a test, and replaced should be this: [[\"Jl. SMA Aek Kota Batu\",\"id\"],[\"Sumatera Utara\",\"de\"]], yes that is what should be replaced.",
		"[[\"а/д Вятка\",\"ru\"]]",
	}
	outputs := []string{
		"This is a test, and replaced should be this: [[\"\",\"\"],[\"\",\"\"]], yes that is what should be replaced.",
		"[[\"\",\"\"]]",
	}

	for i := range inputs {
		out := string(filterPhotometa([]byte(inputs[i])))
		if out != outputs[i] {
			t.Fatal("Expected\n", outputs[i], "\nbut got\n", out)
		}
	}
}

func TestFilterStringsGoogle(t *testing.T) {
	inputs := []string{
		"<meta content=\"https://maps.google.com/maps/api/staticmap?center=50.774016%2C6.1014016&amp;zoom=12&amp;size=256x256&amp;language=en&amp;sensor=false&amp;client=google-maps-frontend&amp;signature=Itxkc4DzDYPENGYKu1558fxDqwk\" itemprop=\"image\">",
	}
	outputs := []string{
		"<meta content=\"/maps/api/staticmap?center=50.774016%2C6.1014016&amp;zoom=12&amp;size=256x256&amp;language=en&amp;sensor=false&amp;client=google-maps-frontend&amp;signature=Itxkc4DzDYPENGYKu1558fxDqwk\" itemprop=\"image\">",
	}

	for i := range inputs {
		out := string(filterUrls([]byte(inputs[i])))
		if out != outputs[i] {
			t.Fatal("Expected\n", outputs[i], "\nbut got\n", out)
		}
	}
}
