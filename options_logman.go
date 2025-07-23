package logman

import "fmt"

// LogmanOptions - settings for logMan object.
type LogmanOptions func(*options)

type options struct {
	appMinimumLoglevel int
	appName            string
	longCallerNames    bool
	logLevels          map[string]*loggingLevel
	colorizer          Colorizer
	globalWriterKeys   []string
	globalFormatters   []*formatterExpanded
}

func defaultOpts() options {
	return options{
		appMinimumLoglevel: ImportanceALL,
		logLevels:          defaultLoggingLevels(),
	}

}

// WithLogLevels sets loglevels to logman with slice of NewLogLevel functions.
// Used to create custom logLevels.
// Caution: It overrides default levels if new loglevel has standard key ("fatal", "error", "warn", "info", "debug", "trace").
func WithLogLevels(lvls ...*loggingLevel) LogmanOptions {
	return func(o *options) {
		//o.logLevels = make(map[string]*loggingLevel)
		for _, lvl := range lvls {
			o.logLevels[lvl.name] = lvl
		}
	}
}

// WithAppLogLevelImportance sets minimum message importance level logMan will process.
// If input is below ImportanceNone importance will be set to ImportanceNone.
// If input is above ImportanceALL importance will be set to ImportanceALL.
func WithAppLogLevelImportance(importance int) LogmanOptions {
	return func(o *options) {
		if importance > ImportanceNONE {
			importance = ImportanceNONE
		}
		if importance < ImportanceALL {
			importance = ImportanceALL
		}
		o.appMinimumLoglevel = importance
	}
}

// WithGlobalColorizer - sets global color scheme for logman
func WithGlobalColorizer(colorizer Colorizer) LogmanOptions {
	return func(o *options) {
		o.colorizer = colorizer
	}
}

// WithGlobalWriterFormatter - Add writer to all level.
// Useful to setup logfile.
func WithGlobalWriterFormatter(writer string, formatter *formatterExpanded) LogmanOptions {
	return func(o *options) {
		o.globalWriterKeys = append(o.globalWriterKeys, writer)
		o.globalFormatters = append(o.globalFormatters, formatter)
	}
}

// WithJSON - Add json writer to all levels.
// Useful to setup logfile.
func WithJSON(directory string) LogmanOptions {
	formatter := NewFormatter(WithCustomFunc("json", stdJSON), WithRequestedFields([]string{"json"}))
	return func(o *options) {
		o.globalWriterKeys = append(o.globalWriterKeys, directory)
		o.globalFormatters = append(o.globalFormatters, formatter)
	}
}

func WithAppName(name string) LogmanOptions {
	return func(o *options) {
		o.appName = name
	}
}

//AFTER SETUP CONTROL

func SetLevelWriterFormatter(level, writer string, formatter *formatterExpanded) error {
	if _, ok := logMan.logLevels[level]; !ok {
		return fmt.Errorf("logman has no level '%v'", level)
	}
	logMan.logLevels[level].writerFormatterMap[writer] = formatter
	return nil
}

func ResetWriters(levels ...string) error {
	for _, level := range levels {
		if _, ok := logMan.logLevels[level]; !ok {
			return fmt.Errorf("logman has no level '%v'", level)
		}
		logMan.logLevels[level].writerFormatterMap = make(map[string]*formatterExpanded)
	}
	return nil
}

func RemovetWriter(level, writer string) error {
	if _, ok := logMan.logLevels[level]; !ok {
		return fmt.Errorf("logman has no level '%v'", level)
	}
	if _, ok := logMan.logLevels[level].writerFormatterMap[writer]; !ok {
		return fmt.Errorf("logman level '%v' has no writer '%v'", level, writer)
	}
	delete(logMan.logLevels[level].writerFormatterMap, writer)
	return nil
}
