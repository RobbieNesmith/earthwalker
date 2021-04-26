package handlers

import (
	"testing"
)

func TestFilterStringsPhotometa(t *testing.T) {
	inputs := []string{
		"This is a test, and replaced should be this: [[\"Jl. SMA Aek Kota Batu\",\"id\"],[\"Sumatera Utara\",\"de\"]], yes that is what should be replaced.",
		"[[\"а/д Вятка\",\"ru\"]]",
		"[[null,[\"5169892334252101127\",\"2545753247963106749\"],null,\"2545753247963106749\"],null,[\"The Body Shop\",\"de\"],[\"Cosmetics store\",\"en\"],\"https://maps.gstatic.com/mapfiles/annotations/icons/shopping_closed_2x.5.png\",null,null,4]",
		"\"https://maps.gstatic.com/mapfiles/annotations/icons/medical_2x.5.png\"",
	}
	outputs := []string{
		"This is a test, and replaced should be this: [[\"\",\"\"],[\"\",\"\"]], yes that is what should be replaced.",
		"[[\"\",\"\"]]",
		"[[null,[\"5169892334252101127\",\"2545753247963106749\"],null,\"2545753247963106749\"],null,[\"\",\"\"],[\"\",\"\"],\"\",null,null,4]",
		"\"\"",
	}

	for i := range inputs {
		out := filterPhotometa(inputs[i])
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
		out := filterUrls(inputs[i])
		if out != outputs[i] {
			t.Fatal("Expected\n", outputs[i], "\nbut got\n", out)
		}
	}
}
