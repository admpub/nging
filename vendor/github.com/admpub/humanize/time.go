package humanize

// Time values humanization functions.

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Time constants.
const (
	Second   = 1
	Minute   = 60
	Hour     = 60 * Minute
	Day      = 24 * Hour
	Week     = 7 * Day
	Month    = 30 * Day
	Year     = 12 * Month
	LongTime = 35 * Year
)

// buildTimeInputRe will build a regular expression to match all possible time inputs.
func (humanizer *Humanizer) buildTimeInputRe() {
	// Get all possible time units.
	units := make([]string, 0, len(humanizer.provider.Times.Units))
	for unit := range humanizer.provider.Times.Units {
		units = append(units, unit)
	}
	// Regexp will match: number, optional coma or dot, optional second number, unit name
	humanizer.timeInputRe = regexp.MustCompile("([0-9]+)[.,]?([0-9]*?) (" + strings.Join(units, "|") + ")")
}

// humanizeDuration will return a humanized form of time duration.
func (humanizer *Humanizer) humanizeDuration(seconds int64, precise int) string {
	if seconds < 1 {
		panic("Cannot humanize durations < 1 sec.")
	}
	secondsLeft := seconds
	humanized := []string{}

	for i := -1; secondsLeft > 0 && (precise < 0 || i < precise); i++ {
		// Find the ranges closest but bigger then diff.
		n := sort.Search(len(humanizer.provider.Times.Ranges), func(i int) bool {
			return humanizer.provider.Times.Ranges[i].UpperLimit > secondsLeft
		})

		// Within the ranges find the one matching our time best.
		timeRanges := humanizer.provider.Times.Ranges[n]
		k := sort.Search(len(timeRanges.Ranges), func(i int) bool {
			return timeRanges.Ranges[i].UpperLimit > secondsLeft
		})
		timeRange := timeRanges.Ranges[k]
		actualTime := secondsLeft / timeRanges.DivideBy // Integer division!

		// If range has a placeholder for a number, insert it.
		if strings.Contains(timeRange.Format, "%d") {
			humanized = append(humanized, fmt.Sprintf(timeRange.Format, actualTime))
		} else {
			humanized = append(humanized, timeRange.Format)
		}

		// Subtract the time span covered by this part.
		secondsLeft -= actualTime * timeRanges.DivideBy
		if precise == 0 { // We don't care about the reminder.
			secondsLeft = 0
		}
	}

	if len(humanized) == 1 {
		return humanized[0]
	}
	return fmt.Sprintf(
		"%s%s%s",
		strings.Join(humanized[:len(humanized)-1], ", "),
		humanizer.provider.Times.RemainderSep,
		humanized[len(humanized)-1],
	)
}

// TimeDiffNow is a convenience method returning humanized time from now till date.
func (humanizer *Humanizer) TimeDiffNow(date time.Time, precise int) string {
	return humanizer.TimeDiff(time.Now(), date, precise)
}

// TimeDiff will return the humanized time difference between the given dates.
// Precise setting determines whether a rough approximation or exact description should be returned, e.g.:
//   precise=0 -> "3 months"
//   precise=2  -> "2 months, 1 week and 3 days"
//
// TODO: in precise mode some ranges should be skipped, like weeks in the example above.
func (humanizer *Humanizer) TimeDiff(startDate, endDate time.Time, precise int) string {
	diff := endDate.Unix() - startDate.Unix()

	if diff == 0 {
		return humanizer.provider.Times.Now
	}

	// Don't bother with Math.Abs
	absDiff := diff
	if absDiff < 0 {
		absDiff = -absDiff
	}

	humanized := humanizer.humanizeDuration(absDiff, precise)

	// Past or future?
	if diff > 0 {
		return fmt.Sprintf(humanizer.provider.Times.Future, humanized)
	}
	return fmt.Sprintf(humanizer.provider.Times.Past, humanized)
}

// ParseDuration will return time duration as parsed from input string.
func (humanizer *Humanizer) ParseDuration(input string) (time.Duration, error) {
	allMatched := humanizer.timeInputRe.FindAllStringSubmatch(input, -1)
	if len(allMatched) == 0 {
		return time.Duration(0), fmt.Errorf("Cannot parse '%s'", input)
	}

	totalDuration := time.Duration(0)
	for _, matched := range allMatched {
		// 0 - full match, 1 - number, 2 - decimal, 3 - unit
		if matched[2] == "" { // Decimal component is empty.
			matched[2] = "0"
		}
		// Parse first two groups into a float.
		number, err := strconv.ParseFloat(matched[1]+"."+matched[2], 64)
		if err != nil {
			return time.Duration(0), err
		}
		// Get the value of the unit in seconds.
		seconds, _ := humanizer.provider.Times.Units[matched[3]]
		// Parser will simply sum up all the found durations.
		totalDuration += time.Duration(int64(number * float64(seconds) * float64(time.Second)))
	}

	return totalDuration, nil
}
