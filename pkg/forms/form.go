package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Form creates a custom Form struct, which anonymously embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initializes a custom Form struct.
func New(data url.Values, db ...string) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checks if field is in post
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}

// Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// IsEmail checks for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}

// IsIP checks for valid ip address
func (f *Form) IsIP(field string) {
	if !govalidator.IsIP(f.Get(field)) {
		f.Errors.Add(field, "Invalid ip address")
	}
}

// MaxLength validates max length of field
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}

// EqualTo validates the start of a string
func (f *Form) EqualTo(field1, field2 string) {
	if f.Get(field1) != f.Get(field2) {
		f.Errors.Add(field2, "Please enter the same value again")
	}
}

// StartsWith validates the start of a string
func (f *Form) StartsWith(field string, matcher string) {
	if !strings.HasPrefix(f.Get(field), matcher) {
		f.Errors.Add(field, fmt.Sprintf("This field must start with %s)", matcher))
	}
}

// EndsWith validates the end of a string
func (f *Form) EndsWith(field string, matcher string) {
	if !strings.HasSuffix(f.Get(field), matcher) {
		f.Errors.Add(field, fmt.Sprintf("This field must end with %s)", matcher))
	}
}

// Contains validates existence of substring
func (f *Form) Contains(field string, matcher string) {
	if !strings.Contains(f.Get(field), matcher) {
		f.Errors.Add(field, fmt.Sprintf("This field must contain %s)", matcher))
	}
}

// MatchRegEx validates regex pattern match
func (f *Form) MatchRegEx(field string, pattern string) {
	rxPat := regexp.MustCompile(pattern)

	if !rxPat.MatchString(f.Get(field)) {
		f.Errors.Add(field, "This field is in an invalid format")
	}
}

// IsURL validates that field is a url
func (f *Form) IsURL(field string) {
	u, err := url.Parse(f.Get(field))

	if err != nil {
		f.Errors.Add(field, "This field is not a valid URL")
	} else if u.Scheme == "" || u.Host == "" {
		f.Errors.Add(field, "This field must be an absolute URL")
	} else if u.Scheme != "http" && u.Scheme != "https" {
		f.Errors.Add(field, "The URL must start with http or https")
	}

}

// IsInt validates int
func (f *Form) IsInt(field string) {
	n, err := strconv.Atoi(f.Get(field))
	if err != nil {
		f.Errors.Add(field, "This field must be an integer")
	} else if n < 0 {
		// just here so the compiler does not throw an error
	}
}

// IsFloat validates float
func (f *Form) IsFloat(field string) {
	n, err := strconv.ParseFloat(f.Get(field), 64)
	if err != nil {
		f.Errors.Add(field, fmt.Sprintf("%f is not a float", n))
	} else if n < 0 {
		// just here so the compiler does not throw an error
	}
}

// IsDate validates dateISO
func (f *Form) IsDate(field string) {
	d, err := time.Parse("2006-01-02", f.Get(field))
	if err != nil {
		f.Errors.Add(field, "This field must be a date in the form YYYY-MM-DD")
	} else if d.Year() != 2017 {
		// just here so the compiler does not throw an error
	}
}

// NoSpaces checks for white space
func (f *Form) NoSpaces(field string) {
	if govalidator.HasWhitespace(f.Get(field)) {
		f.Errors.Add(field, "Spaces are not permitted")
	}
}

// PermittedValues validation
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

// MinLength method to check that a specific field in the form
func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", d))
	}
}

// MatchesPattern checks that a specific field in the form matches regex pattern
func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}

// HasFile returns true if post has a an attachment, otherwise false
func (f *Form) HasFile(field string, r *http.Request) bool {
	_, _, err := r.FormFile(field)
	if err != nil {
		if err == http.ErrMissingFile {
			return false
		}
	}
	return true
}

// Valid method which returns true if there are no errors.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
