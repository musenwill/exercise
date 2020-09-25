package priv

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// DateFormat represents the format for date literals.
	DateFormat = "2006-01-02"

	// DateTimeFormat represents the format for date time literals.
	DateTimeFormat = "2006-01-02 15:04:05.999999"
)

var (
	// ErrInvalidIP represents the error when parse IP
	ErrInvalidIP = errors.New("IP address invalid")
)

// Value represents a value that can be bound
// to a parameter when parsing the query.
type Value interface {
	TokenType() Token
	Value() string
}

type (
	// Identifier is an identifier value.
	Identifier string

	// StringValue is a string literal.
	StringValue string

	// RegexValue is a regexp literal.
	RegexValue string

	// NumberValue is a number literal.
	NumberValue float64

	// IntegerValue is an integer literal.
	IntegerValue int64

	// BooleanValue is a boolean literal.
	BooleanValue bool

	// DurationValue is a duration literal.
	DurationValue string

	// ErrorValue is a special value that returns an error during parsing
	// when it is used.
	ErrorValue string
)

// Parser represents an fql parser.
type Parser struct {
	s      *bufScanner
	params map[string]Value
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: newBufScanner(r)}
}

// peekRune returns the next rune that would be read by the scanner.
func (p *Parser) peekRune() rune {
	r, _, _ := p.s.s.r.ReadRune()
	if r != eof {
		_ = p.s.s.r.UnreadRune()
	}

	return r
}

// ParseIdent parses an identifier.
func (p *Parser) ParseIdent() (string, error) {
	tok, pos, lit := p.ScanIgnoreWhitespace()
	if tok != IDENT {
		return "", newParseError(tokstr(tok, lit), []string{"identifier"}, pos)
	}
	return lit, nil
}

// parseSegmentedIdents parses a segmented identifiers.
// e.g.,  "db"."rp".measurement  or  "db"..measurement
func (p *Parser) parseSegmentedIdents() ([]string, error) {
	ident, err := p.ParseIdent()
	if err != nil {
		return nil, err
	}
	idents := []string{ident}

	// Parse remaining (optional) identifiers.
	for {
		if tok, _, _ := p.Scan(); tok != DOT {
			// No more segments so we're done.
			p.Unscan()
			break
		}

		if ch := p.peekRune(); ch == '/' {
			// Next segment is a regex so we're done.
			break
		} else if ch == ':' {
			// Next segment is context-specific so let caller handle it.
			break
		} else if ch == '.' {
			// Add an empty identifier.
			idents = append(idents, "")
			continue
		}

		// Parse the next identifier.
		if ident, err = p.ParseIdent(); err != nil {
			return nil, err
		}

		idents = append(idents, ident)
	}

	if len(idents) > 3 {
		msg := fmt.Sprintf("too many segments in %s", QuoteIdent(idents...))
		return nil, &ParseError{Message: msg}
	}

	return idents, nil
}

// ParseOptionalTokenAndInt parses the specified token followed
// by an int, if it exists.
func (p *Parser) ParseOptionalTokenAndInt(t Token) (int, error) {
	// Check if the token exists.
	if tok, _, _ := p.ScanIgnoreWhitespace(); tok != t {
		p.Unscan()
		return 0, nil
	}

	// Scan the number.
	tok, pos, lit := p.ScanIgnoreWhitespace()
	if tok != INTEGER {
		return 0, newParseError(tokstr(tok, lit), []string{"integer"}, pos)
	}

	// Parse number.
	n, _ := strconv.ParseInt(lit, 10, 64)
	if n < 0 {
		msg := fmt.Sprintf("%s must be >= 0", t.String())
		return 0, &ParseError{Message: msg, Pos: pos}
	}

	return int(n), nil
}

// parseResample parses a RESAMPLE [EVERY <duration>] [FOR <duration>].
// This function assumes RESAMPLE has already been consumed.
// EVERY and FOR are optional, but at least one of the two has to be used.
func (p *Parser) parseResample() (time.Duration, time.Duration, error) {
	var interval time.Duration
	if tok, _, _ := p.ScanIgnoreWhitespace(); tok == EVERY {
		tok, pos, lit := p.ScanIgnoreWhitespace()
		if tok != DURATIONVAL {
			return 0, 0, newParseError(tokstr(tok, lit), []string{"duration"}, pos)
		}

		d, err := ParseDuration(lit)
		if err != nil {
			return 0, 0, &ParseError{Message: err.Error(), Pos: pos}
		}
		interval = d
	} else {
		p.Unscan()
	}

	var maxDuration time.Duration
	if tok, _, _ := p.ScanIgnoreWhitespace(); tok == FOR {
		tok, pos, lit := p.ScanIgnoreWhitespace()
		if tok != DURATIONVAL {
			return 0, 0, newParseError(tokstr(tok, lit), []string{"duration"}, pos)
		}

		d, err := ParseDuration(lit)
		if err != nil {
			return 0, 0, &ParseError{Message: err.Error(), Pos: pos}
		}
		maxDuration = d
	} else {
		p.Unscan()
	}

	// Neither EVERY or FOR were read, so read the next token again
	// so we can return a suitable error message.
	if interval == 0 && maxDuration == 0 {
		tok, pos, lit := p.ScanIgnoreWhitespace()
		return 0, 0, newParseError(tokstr(tok, lit), []string{"EVERY", "FOR"}, pos)
	}
	return interval, maxDuration, nil
}

// Scan returns the next token from the underlying scanner.
func (p *Parser) Scan() (tok Token, pos Pos, lit string) {
	return p.scan(p.s.Scan)
}

// ScanRegex returns the next token from the underlying scanner
// using the regex scanner.
func (p *Parser) ScanRegex() (tok Token, pos Pos, lit string) {
	return p.scan(p.s.ScanRegex)
}

type scanFunc func() (tok Token, pos Pos, lit string)

func (p *Parser) scan(fn scanFunc) (tok Token, pos Pos, lit string) {
	tok, pos, lit = fn()
	if tok == BOUNDPARAM {
		// If we have a bound parameter, attempt to
		// replace it in the scanner. If the bound parameter
		// isn't valid, do not perform the replacement.
		k := strings.TrimPrefix(lit, "$")
		if len(k) != 0 {
			if v, ok := p.params[k]; ok {
				tok, lit = v.TokenType(), v.Value()
			}
		}
	}
	return tok, pos, lit
}

// ScanIgnoreWhitespace scans the next non-whitespace and non-comment token.
func (p *Parser) ScanIgnoreWhitespace() (tok Token, pos Pos, lit string) {
	for {
		tok, pos, lit = p.Scan()
		if tok == WS || tok == COMMENT {
			continue
		}
		return
	}
}

// consumeWhitespace scans the next token if it's whitespace.
func (p *Parser) consumeWhitespace() {
	if tok, _, _ := p.Scan(); tok != WS {
		p.Unscan()
	}
}

// Unscan pushes the previously read token back onto the buffer.
func (p *Parser) Unscan() { p.s.Unscan() }

// ParseDuration parses a time duration from a string.
// This is needed instead of time.ParseDuration because this will support
// the full syntax that fql supports for specifying durations
// including weeks and days.
func ParseDuration(s string) (time.Duration, error) {
	// Return an error if the string is blank or one character
	if len(s) < 2 {
		return 0, ErrInvalidDuration
	}

	// Split string into individual runes.
	a := []rune(s)

	// Start with a zero duration.
	var d time.Duration
	i := 0

	// Check for a negative.
	isNegative := false
	if a[i] == '-' {
		isNegative = true
		i++
	}

	var measure int64
	var unit string

	// Parsing loop.
	for i < len(a) {
		// Find the number portion.
		start := i
		for ; i < len(a) && isDigit(a[i]); i++ {
			// Scan for the digits.
		}

		// Check if we reached the end of the string prematurely.
		if i >= len(a) || i == start {
			return 0, ErrInvalidDuration
		}

		// Parse the numeric part.
		n, err := strconv.ParseInt(string(a[start:i]), 10, 64)
		if err != nil {
			return 0, ErrInvalidDuration
		}
		measure = n

		// Extract the unit of measure.
		// If the last two characters are "ms" then parse as milliseconds.
		// Otherwise just use the last character as the unit of measure.
		unit = string(a[i])
		switch a[i] {
		case 'n':
			if i+1 < len(a) && a[i+1] == 's' {
				unit = string(a[i : i+2])
				d += time.Duration(n)
				i += 2
				continue
			}
			return 0, ErrInvalidDuration
		case 'u', 'µ':
			d += time.Duration(n) * time.Microsecond
		case 'm':
			if i+1 < len(a) && a[i+1] == 's' {
				unit = string(a[i : i+2])
				d += time.Duration(n) * time.Millisecond
				i += 2
				continue
			}
			d += time.Duration(n) * time.Minute
		case 's':
			d += time.Duration(n) * time.Second
		case 'h':
			d += time.Duration(n) * time.Hour
		case 'd':
			d += time.Duration(n) * 24 * time.Hour
		case 'w':
			d += time.Duration(n) * 7 * 24 * time.Hour
		default:
			return 0, ErrInvalidDuration
		}
		i++
	}

	// Check to see if we overflowed a duration
	if d < 0 && !isNegative {
		return 0, fmt.Errorf("overflowed duration %d%s: choose a smaller duration or INF", measure, unit)
	}

	if isNegative {
		d = -d
	}
	return d, nil
}

// FormatDuration formats a duration to a string.
func FormatDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	} else if d%(7*24*time.Hour) == 0 {
		return fmt.Sprintf("%dw", d/(7*24*time.Hour))
	} else if d%(24*time.Hour) == 0 {
		return fmt.Sprintf("%dd", d/(24*time.Hour))
	} else if d%time.Hour == 0 {
		return fmt.Sprintf("%dh", d/time.Hour)
	} else if d%time.Minute == 0 {
		return fmt.Sprintf("%dm", d/time.Minute)
	} else if d%time.Second == 0 {
		return fmt.Sprintf("%ds", d/time.Second)
	} else if d%time.Millisecond == 0 {
		return fmt.Sprintf("%dms", d/time.Millisecond)
	} else if d%time.Microsecond == 0 {
		// Although we accept both "u" and "µ" when reading microsecond durations,
		// we output with "u", which can be represented in 1 byte,
		// instead of "µ", which requires 2 bytes.
		return fmt.Sprintf("%du", d/time.Microsecond)
	}
	return fmt.Sprintf("%dns", d)
}

// parseTokens consumes an expected sequence of tokens.
func (p *Parser) parseTokens(toks []Token) error {
	for _, expected := range toks {
		if tok, pos, lit := p.ScanIgnoreWhitespace(); tok != expected {
			return newParseError(tokstr(tok, lit), []string{tokens[expected]}, pos)
		}
	}
	return nil
}

var (
	// Quote String replacer.
	qsReplacer = strings.NewReplacer("\n", `\n`, `\`, `\\`, `'`, `\'`)

	// Quote Ident replacer.
	qiReplacer = strings.NewReplacer("\n", `\n`, `\`, `\\`, `"`, `\"`)
)

// QuoteString returns a quoted string.
func QuoteString(s string) string {
	return `'` + qsReplacer.Replace(s) + `'`
}

// QuoteIdent returns a quoted identifier from multiple bare identifiers.
func QuoteIdent(segments ...string) string {
	var buf bytes.Buffer
	for i, segment := range segments {
		needQuote := IdentNeedsQuotes(segment) ||
			((i < len(segments)-1) && segment != "") || // not last segment && not ""
			((i == 0 || i == len(segments)-1) && segment == "") // the first or last segment and an empty string

		if needQuote {
			_ = buf.WriteByte('"')
		}

		_, _ = buf.WriteString(qiReplacer.Replace(segment))

		if needQuote {
			_ = buf.WriteByte('"')
		}

		if i < len(segments)-1 {
			_ = buf.WriteByte('.')
		}
	}
	return buf.String()
}

// IdentNeedsQuotes returns true if the ident string given would require quotes.
func IdentNeedsQuotes(ident string) bool {
	// check if this identifier is a keyword
	tok := Lookup(ident)
	if tok != IDENT {
		return true
	}
	for i, r := range ident {
		if i == 0 && !isIdentFirstChar(r) {
			return true
		} else if i > 0 && !isIdentChar(r) {
			return true
		}
	}
	return false
}

// isDateString returns true if the string looks like a date-only time literal.
func isDateString(s string) bool { return dateStringRegexp.MatchString(s) }

// isDateTimeString returns true if the string looks like a date+time time literal.
func isDateTimeString(s string) bool { return dateTimeStringRegexp.MatchString(s) }

var dateStringRegexp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
var dateTimeStringRegexp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}.+`)

// ErrInvalidDuration is returned when parsing a malformed duration.
var ErrInvalidDuration = errors.New("invalid duration")

// ParseError represents an error that occurred during parsing.
type ParseError struct {
	Message  string
	Found    string
	Expected []string
	Pos      Pos
}

// newParseError returns a new instance of ParseError.
func newParseError(found string, expected []string, pos Pos) *ParseError {
	return &ParseError{Found: found, Expected: expected, Pos: pos}
}

// Error returns the string representation of the error.
func (e *ParseError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s at line %d, char %d", e.Message, e.Pos.Line+1, e.Pos.Char+1)
	}
	return fmt.Sprintf("found %s, expected %s at line %d, char %d", e.Found, strings.Join(e.Expected, ", "), e.Pos.Line+1, e.Pos.Char+1)
}
