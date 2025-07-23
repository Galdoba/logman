package logman

import (
	"testing"

	"github.com/Galdoba/logman/colorizer"
)

func TestLogMan(t *testing.T) {
	//	logMan.appName = "scribe"
	//jsonFormatter := NewFormatter(WithCustomFunc("json", stdJSON), WithRequestedFields([]string{"json"}))

	Setup(
		WithGlobalColorizer(colorizer.DefaultScheme()),
		// WithGlobalWriterFormatter(Stderr, stdFormatFunc_time),
		// WithRequestedFields(Request_ShortTime),
		// WithGlobalWriterFormatter(Stderr, NewFormatter(WithRequestedFields(Request_ShortTime))),
		//WithGlobalWriterFormatter(`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\pkg\logman\v2\`, jsonFormatter),
		// WithJSON(`c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\pkg\logman\v2\`),
		WithAppName("scribe_test"),
	)
	// SetOutput(Stderr, ALL)
	Info("test error: %v", "some error")
}
