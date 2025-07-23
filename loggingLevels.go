package logman

const (
	FATAL = "fatal"
	ERROR = "error"
	WARN  = "warn"
	INFO  = "info"
	DEBUG = "debug"
	TRACE = "trace"
	PING  = "ping"
	ALL   = "All_Levels"
)

var STDLevels = []string{TRACE, DEBUG, INFO, WARN, ERROR, FATAL}

type loggingLevel struct {
	name               string
	tag                string
	importance         int
	callerInfo         bool
	osExit             bool
	colorSchemes       map[string]uint8
	formatFunc         func(Message) (string, error)
	FMTE               *formatterExpanded
	writerFormatterMap map[string]*formatterExpanded
}

func NewLoggingLevel(name string, opts ...LevelOpts) *loggingLevel {
	lo := loggingLevel{}
	lo.name = name
	options := defaultLevel()
	lo.tag = name
	for _, enrich := range opts {
		enrich(&options)
	}
	lo.callerInfo = options.callerInfo
	lo.osExit = options.osExit
	lo.importance = options.importance
	lo.writerFormatterMap = options.writerFormatterMap
	if options.tag != "" {
		lo.tag = options.tag
	}
	return &lo
}

func LevelTag(tag string) LevelOpts {
	return func(lvl *lvlOpts) {
		lvl.tag = tag
	}
}

func LevelImportance(imp int) LevelOpts {
	return func(lvl *lvlOpts) {
		lvl.importance = imp
	}
}

func LevelCallerInfo(callerInfo bool) LevelOpts {
	return func(lvl *lvlOpts) {
		lvl.callerInfo = callerInfo
	}
}

func LevelExitWhenDone(osExit bool) LevelOpts {
	return func(lvl *lvlOpts) {
		lvl.osExit = osExit
	}
}

type lvlOpts struct {
	tag                string
	importance         int
	callerInfo         bool
	osExit             bool
	writerFormatterMap map[string]*formatterExpanded
}

func defaultLevel() lvlOpts {
	return lvlOpts{
		tag:                "",
		importance:         ImportanceINFO,
		callerInfo:         false,
		writerFormatterMap: make(map[string]*formatterExpanded),
	}
}

type LevelOpts func(*lvlOpts)

func WithWriter(writerKey string, expandedFormatter *formatterExpanded) LevelOpts {
	return func(lo *lvlOpts) {
		lo.writerFormatterMap[writerKey] = expandedFormatter
	}
}
