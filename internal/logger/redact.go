package logger

import "regexp"

// WatchZA fork: a Quazar stream key is the path a publisher publishes to (app/<key>), and
// MediaMTX writes path names into its logs — "[path app/<key>]" on every line a path logs,
// "is publishing to path '<key>'", and errors that name the path. Anyone who can read the
// journal could then publish as that channel.
//
// streamKeySecret matches the key Quazar's control issues, sk_<region>_ and 36 letters and
// digits. The region is kept, so a line still reads as a stream key; the 36-character
// secret is replaced.
var streamKeySecret = regexp.MustCompile(`(sk_[a-z0-9-]+_)[A-Za-z0-9]{36}`)

const redactedSecret = "[redacted]"

// redact hides every stream key's secret in a formatted log line.
func redact(msg string) string {
	return streamKeySecret.ReplaceAllString(msg, "${1}"+redactedSecret)
}
