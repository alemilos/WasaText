package validation

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

// Breakdown of the pattern:
//  1. Core Emojis: [\x{2000}-\x{2BFF}\x{E000}-\x{FFFF}\x{1F000}-\x{1FFFF}] - Covers all major emojis
//  2. Modifiers and Joiners (zero or more):
//     - [\x{1F3FB}-\x{1F3FF}]: Skin tone modifiers.
//     - \x{200D}: Zero Width Joiner (ZWJ) to form compound emojis.
//     - \x{FE0F}: Variation Selector (ensures a character renders as an emoji, not text).
//     - Regional indicators (flags) are included in the wider 1F000 range.
var emojiRegex = regexp.MustCompile(
	`^` +
		// Match one or more core emoji characters/symbols
		`([\x{2000}-\x{2BFF}\x{E000}-\x{FFFF}\x{1F000}-\x{1FFFF}])+` +
		// Allow zero or more modifier/joiner characters
		`([\x{1F3FB}-\x{1F3FF}\x{200D}\x{FE0F}])*` +
		`$`,
)

func ValidateEmoji(s string) error {
	// Basic UTF-8 validation
	if !utf8.ValidString(s) {
		return fmt.Errorf("invalid_utf8_sequence: must be valid UTF-8")
	}

	// Comprehensive regex match
	if !emojiRegex.MatchString(s) {
		// Note: This check implicitly handles empty strings because of the '+' quantifier.
		return fmt.Errorf("invalid_emoji_unicode_or_sequence: does not match known emoji patterns")
	}

	// Check max length (using 10 runes, as defined in your OpenAPI pattern)
	// This prevents excessively long/malicious input, even if it's composed of valid parts.
	if utf8.RuneCountInString(s) > 10 {
		return fmt.Errorf("emoji_sequence_too_long: maximum 10 runes allowed")
	}

	return nil
}
