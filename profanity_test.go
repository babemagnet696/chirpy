package main

import (
	"testing"
)

func TestProfanityCensor(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "What in the heckity heck!!!",
			expected: "What in the heckity heck!!!",
		},{
			input:    "there is a kerfuffle over there!",
        	expected: "there is a **** over there!",
		},{
			input:    "I just SHARBERT my pants...",
			expected: "I just **** my pants...",
		},{
			input:    "What was that one game called? FoRnAx right?",
			expected: "What was that one game called? **** right?",
		},{
			input:    "I know how to cheat the system! Kerfuffle!",
			expected: "I know how to cheat the system! Kerfuffle!",
		},
	}
	for _, c := range cases {
		actual := profanityCheck(c.input)
		if actual != c.expected {
			t.Errorf("Profanity not properly censored\nActual: '%s'\nExpected: '%s'", actual, c.expected)
		}
	}
}