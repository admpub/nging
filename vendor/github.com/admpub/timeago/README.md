[![Code Coverage](https://scrutinizer-ci.com/g/SerhiiCho/timeago/badges/coverage.png?b=master)](https://scrutinizer-ci.com/g/SerhiiCho/timeago/?branch=master)
[![Build Status](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2FSerhiiCho%2Ftimeago%2Fbadge%3Fref%3Dmaster&style=flat)](https://actions-badge.atrox.dev/SerhiiCho/timeago/goto?ref=master)
[![Build Status](https://scrutinizer-ci.com/g/SerhiiCho/timeago/badges/build.png?b=master)](https://scrutinizer-ci.com/g/SerhiiCho/timeago/build-status/master)
[![Scrutinizer Code Quality](https://scrutinizer-ci.com/g/SerhiiCho/timeago/badges/quality-score.png?b=master)](https://scrutinizer-ci.com/g/SerhiiCho/timeago/?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/SerhiiCho/timeago)](https://goreportcard.com/report/github.com/SerhiiCho/timeago)
![GitHub](https://img.shields.io/github/license/SerhiiCho/timeago)

100% test covered, fast and lightweight datetime converter that converts given datetime into "n time ago" format from [YEAR-MONTH-DAY HOURS:MINUTES:SECONDS]. Supports Russian and English languages.

[Example usage on repl.it](https://repl.it/@SerhiiCho/Usage-of-timeago-package)

## Configurations
##### Language
Default language is English. Optionally you can set the language in your application by calling `timeago.Set()` function and pass the flag "ru" or "en" in your init function.

```go
func init() {
    timeago.Set("language", "ru")
}
```

##### Location
Default location is Europe/Kiev. Optionally you can set the location in your application by calling `timeago.Set()` function and pass the location you need in your init function.

```go
func init() {
    timeago.Set("location", "America/New_York")
}
```

> Please make sure that timezone configuration is correct for your location. It is very important for displaying the correct datetime.

## Usage

For outputting post publishing date or something else you can just pass the date to method `timeago.Take()`. It will count the interval between now and given date and returns converted format.

```go
timeago.Take("2019-10-23 10:46:00") // after 10 seconds outputs: 10 seconds ago
```

If you want to show last user login like if user is online or not, you can optionally add `|online` to the datetime string. All it does is just displaying **Online** if date interval withing 60 seconds.

```go
timeago.Take("2019-10-23 10:46:00|online")
```

## Quick Start

```bash
go get -u github.com/SerhiiCho/timeago
```
