package ttlmap

import (
	"math"
	"time"
)

// Item represents an item with an associated value and optional expiration.
type Item struct {
	value      interface{}
	expiration time.Time
	expires    bool
}

// NewItem creates an item with the specified value and optional expiration.
func NewItem(value interface{}, expiration *time.Time) Item {
	var expiration2 time.Time
	if expiration != nil {
		expiration2 = *expiration
	}
	return Item{
		value:      value,
		expiration: expiration2,
		expires:    (expiration != nil),
	}
}

// Value returns the value stored in the item.
func (item *Item) Value() interface{} {
	return item.value
}

// Expiration returns the item's expiration time.
func (item *Item) Expiration() time.Time {
	return item.expiration
}

// TTL returns the remaining duration until expiration (negative if expired).
func (item *Item) TTL() time.Duration {
	if item.expires {
		return item.expiration.Sub(time.Now())
	}
	return time.Duration(math.MaxInt64)
}

// Expired checks whether the item is already expired.
func (item *Item) Expired() bool {
	if item.expires {
		return item.expiration.Before(time.Now())
	}
	return false
}

// Expires checks whether the item has an expiration time set.
func (item *Item) Expires() bool {
	return item.expires
}

// WithExpiration creates an expiration time.
func WithExpiration(expiration time.Time) *time.Time {
	return &expiration
}

// WithTTL creates an expiration time from a specified TTL.
func WithTTL(duration time.Duration) *time.Time {
	expiration := time.Now().Add(duration)
	return &expiration
}
