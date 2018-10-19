package machine

import "regexp"

var InstructionHandlers map[string]InstructionHandler = make(map[string]InstructionHandler)

var noSkip InstructionFilter = make(InstructionFilter)

// pre-compiled regular expressions.
var stringEscaper *regexp.Regexp = regexp.MustCompile(`\[\[[^\]]+\]\]`)
var stringUnescaper *regexp.Regexp = regexp.MustCompile(`___REP___\d{9}___PER___`)
var lineExploder *regexp.Regexp = regexp.MustCompile(`[\n\r]`)
