package eventbus

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_MatchSubject(t *testing.T) {
	tests := []struct {
		name          string
		subject       string
		target        string
		expectedMatch bool
	}{
		{
			name:          "match full subject",
			subject:       "test.one.two.three.four.five",
			target:        "test.one.two.three.four.five",
			expectedMatch: true,
		},
		{
			name:          "match >",
			subject:       "test.one.two.>",
			target:        "test.one.two.three.four.five",
			expectedMatch: true,
		},
		{
			name:          "match > at end",
			subject:       "test.one.two.>",
			target:        "test.one.two.three",
			expectedMatch: true,
		},
		{
			name:          "match single *",
			subject:       "*",
			target:        "test",
			expectedMatch: true,
		},
		{
			name:          "match * in middle",
			subject:       "test.*.two",
			target:        "test.one.two",
			expectedMatch: true,
		},
		{
			name:          "don't match full subject",
			subject:       "test.one.two.three.four.five",
			target:        "test.one.two.four.five",
			expectedMatch: false,
		},
		{
			name:          "don't match *",
			subject:       "test.*.two",
			target:        "test.two.one",
			expectedMatch: false,
		},
		{
			name:          "don't match *, shorter than target",
			subject:       "test.*.two",
			target:        "test.one.two.three",
			expectedMatch: false,
		},
		{
			name:          "don't match >",
			subject:       "test.one.>",
			target:        "test.two.three",
			expectedMatch: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			subjectTokens := strings.Split(test.subject, ".")
			targetTokens := strings.Split(test.target, ".")
			require.Equal(t, test.expectedMatch, matchSubject(subjectTokens, targetTokens))
		})
	}
}

func Test_IsValidSubject(t *testing.T) {
	tests := []struct {
		name     string
		subject  string
		expected bool
	}{
		{
			name:     "space in token invalid",
			subject:  "test .one",
			expected: false,
		},
		{
			name:     "* in token invalid",
			subject:  "test*.one",
			expected: false,
		},
		{
			name:     "> in token invalid",
			subject:  "test>.one",
			expected: false,
		},
		{
			name:     "> not at end",
			subject:  "test.>.one",
			expected: false,
		},
		{
			name:     "blank token invalid",
			subject:  ".one",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, test.expected, isValidSubject(test.subject))
		})
	}
}
