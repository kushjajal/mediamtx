package logger

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const testKey = "sk_ap-south-1_AbCdEfGhIjKlMnOpQrStUvWxYz0123456789"

func TestRedact(t *testing.T) {
	for _, ca := range []struct {
		name string
		in   string
		out  string
	}{
		{"path prefix", "[path app/" + testKey + "] stream is ready",
			"[path app/sk_ap-south-1_[redacted]] stream is ready"},
		{"publishing", "is publishing to path 'app/" + testKey + "', 2 tracks",
			"is publishing to path 'app/sk_ap-south-1_[redacted]', 2 tracks"},
		{"two keys", testKey + " " + testKey,
			"sk_ap-south-1_[redacted] sk_ap-south-1_[redacted]"},
		{"no key", "[RTMP] [conn 1.2.3.4:5678] opened", "[RTMP] [conn 1.2.3.4:5678] opened"},
		{"too short to be a key", "sk_ap-south-1_short", "sk_ap-south-1_short"},
	} {
		t.Run(ca.name, func(t *testing.T) {
			require.Equal(t, ca.out, redact(ca.in))
		})
	}
}

// The key is hidden on every destination and in both formats, including when it arrives as
// an argument rather than in the format string.
func TestLoggerRedactsStreamKey(t *testing.T) {
	for _, structured := range []bool{false, true} {
		var buf bytes.Buffer
		l := &Logger{
			Destinations: []Destination{DestinationStdout},
			Structured:   structured,
			timeNow:      func() time.Time { return time.Date(2003, 11, 4, 23, 15, 8, 0, time.UTC) },
			stdout:       &buf,
		}
		require.NoError(t, l.Initialize())

		l.Log(Info, "[path app/%s] is publishing to path '%s'", testKey, "app/"+testKey)
		l.Log(Info, "a literal %% stays")
		l.Close()

		out := buf.String()
		require.NotContains(t, out, "AbCdEfGhIjKlMnOpQrStUvWxYz0123456789")
		require.Equal(t, 2, strings.Count(out, "sk_ap-south-1_[redacted]"))
		require.Contains(t, out, "a literal % stays")
	}
}
