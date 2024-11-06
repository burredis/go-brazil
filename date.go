package brazil

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	dateDelimiter        = "/"
	dateDelimiterPattern = `[ ]{0,1}%s[ ]{0,1}`
	digitsPattern        = "\\d+"
	DateFormatShort      = "02/01/06"
	DateFormatLong       = "02/01/2006"
)

var (
	delimiters = []string{`\/`, `\.`, `\-`, `[ ]`}
	months     = MonthsOfYear{
		"JAN": 1, "FEV": 2, "MAR": 3, "ABR": 4, "MAI": 5, "JUN": 6, "JUL": 7, "AGO": 8, "SET": 9, "OUT": 10, "NOV": 11, "DEZ": 12,
	}
	monthPattern = `(?:` + digitsPattern + `|(` + strings.Join(months.GetMonths(), "|") + `)[\D!ç]{0,7})`
	datePattern  = `(?i)` + digitsPattern + dateDelimiterPattern + monthPattern + dateDelimiterPattern + digitsPattern
)

type MonthsOfYear map[string]int8

// GetMonths returns the months of the year
func (m MonthsOfYear) GetMonths() []string {
	var keys []string
	for name := range m {
		keys = append(keys, name)
	}
	return keys
}

// Date struct
type date struct {
	value string
}

// Time returns the date as a time.Time
func (d date) Time() time.Time {
	t, _ := time.Parse(DateFormatLong, d.value)
	return t
}

// String returns the date as a string
func (d date) String() string {
	return d.value
}

// ParseDate parses a date from a string
func ParseDate(value string) (date, error) {
	var (
		matches []string
		err     error
		time    *time.Time
	)

	// Create a pattern for each delimiter and find all matches
	for _, delimiter := range delimiters {
		pattern := fmt.Sprintf(datePattern, delimiter, delimiter)
		matches = append(matches, regexp.MustCompile(pattern).FindAllString(value, -1)...)
	}

	// Try to get the time from the matches
	for _, match := range matches {
		time, err = getTime(match)
		// If there has been an error, return it
		if err != nil {
			return date{}, err
		}
	}

	// If the time is nil, return an error
	if time == nil {
		return date{}, fmt.Errorf("invalid date")
	}

	// Return the date
	return date{value: time.Format(DateFormatLong)}, nil
}

// Get the time from the string
func getTime(value string) (*time.Time, error) {
	// Replace all non-word characters with the date delimiter
	r := regexp.MustCompile(`\W+`).ReplaceAllString(value, dateDelimiter)
	// Split the string by the date delimiter
	arr := strings.Split(r, dateDelimiter)

	// Get the day
	day, err := strconv.Atoi(arr[0])
	if err != nil {
		return nil, err
	}

	// Get the month
	month, err := strconv.Atoi(normalizeMonth(arr[1]))
	if err != nil {
		return nil, err
	}

	// Get the year
	year, err := strconv.Atoi(arr[len(arr)-1:][0])
	if err != nil {
		return nil, err
	}

	// Set the default date format and year digits
	dateFormat := DateFormatLong
	yearDigits := "%04d"

	// If the year has 2 digits, change the date format and the year digits
	if len(arr[2]) == 2 {
		dateFormat = DateFormatShort
		yearDigits = "%02d"
	}

	// Parse the date with the format obtained
	date, err := time.Parse(dateFormat, fmt.Sprintf("%02d/%02d/"+yearDigits, day, month, year))
	if err != nil {
		return nil, err
	}

	// Return the date
	return &date, nil
}

// Normalize the month name to number
func normalizeMonth(text string) string {
	// If the text is empty or has less than 3 characters, return it
	if len(text) < 3 {
		return text
	}

	// Get the first 3 characters of the text and convert them to uppercase
	month := strings.ToUpper(text)[0:3]

	// Check if the month is in the map, if not return an empty string
	value, ok := months[month]
	if !ok {
		return ""
	}

	// Return the month number with 2 digits
	return fmt.Sprintf("%02d", value)
}
