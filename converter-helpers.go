package go_rnd_common

import "time"

func GetInt(key interface{}) (i int) {
	if key != nil {
		i, _ = key.(int)
	}

	return
}

// GetInt64 returns the value associated with the key as an integer.
func GetInt64(key interface{}) (i64 int64) {
	if key != nil {
		i64, _ = key.(int64)
	}

	return
}

// GetFloat64 returns the value associated with the key as a float64.
func GetFloat64(key interface{}) (f64 float64) {
	if key != nil {
		f64, _ = key.(float64)
	}

	return
}

// GetTime returns the value associated with the key as time.
func GetTime(key interface{}) (t time.Time) {
	if key != nil {
		t, _ = key.(time.Time)
	}

	return
}

// GetDuration returns the value associated with the key as a duration.
func GetDuration(key interface{}) (d time.Duration) {
	if key != nil {
		d, _ = key.(time.Duration)
	}
	return
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func GetStringSlice(key interface{}) (ss []string) {
	if key != nil {
		ss, _ = key.([]string)
	}
	return
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func GetStringMap(key interface{}) (sm map[string]interface{}) {
	if key != nil {
		sm, _ = key.(map[string]interface{})
	}
	return
}

// GetStringMapString returns the value associated with the key as a map of strings.
func GetStringMapString(key interface{}) (sms map[string]string) {
	if key != nil {
		sms, _ = key.(map[string]string)
	}
	return
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func GetStringMapStringSlice(key interface{}) (smss map[string][]string) {
	if key != nil {
		smss, _ = key.(map[string][]string)
	}
	return
}
