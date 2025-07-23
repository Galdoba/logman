package logman

const (
	stdTagFATAL = "fatal"
	stdTagERROR = "error"
	stdTagWARN  = "warn"
	stdTagINFO  = "info"
	stdTagDEBUG = "debug"
	stdTagTRACE = "trace"
	stdTagPing  = "ping"
)

var LogLevelFATAL = &loggingLevel{
	name:       FATAL,
	tag:        stdTagFATAL,
	importance: ImportanceFATAL,
	callerInfo: true,
	osExit:     true,
	writerFormatterMap: map[string]*formatterExpanded{
		Stdout: NewFormatter(WithRequestedFields(Request_ShortTime)),
	},
}

var LogLevelERROR = &loggingLevel{
	name:       ERROR,
	tag:        stdTagERROR,
	importance: ImportanceERROR,
	callerInfo: true,
	osExit:     false,
	writerFormatterMap: map[string]*formatterExpanded{
		Stdout: NewFormatter(WithRequestedFields(Request_ShortTime)),
	},
}

var LogLevelWARN = &loggingLevel{
	name:       WARN,
	tag:        stdTagWARN,
	importance: ImportanceWARN,
	callerInfo: false,
	osExit:     false,
	//formatFunc: formatTextSimple,
	writerFormatterMap: map[string]*formatterExpanded{
		Stderr: NewFormatter(WithRequestedFields(Request_ShortTime)),
	},
}

var LogLevelINFO = &loggingLevel{
	name:       INFO,
	tag:        stdTagINFO,
	importance: ImportanceINFO,
	callerInfo: false,
	osExit:     false,
	//writers:    map[string]io.Writer{Stderr: os.Stderr},
	//formatFunc: formatTextSimple,
	writerFormatterMap: map[string]*formatterExpanded{
		Stderr: NewFormatter(WithRequestedFields(Request_ShortTime)),
	},
}

var LogLevelDEBUG = &loggingLevel{
	name:       DEBUG,
	tag:        stdTagDEBUG,
	importance: ImportanceDEBUG,
	callerInfo: false,
	osExit:     false,
	//writers:    map[string]io.Writer{Stderr: os.Stderr},
	writerFormatterMap: map[string]*formatterExpanded{
		Stderr: NewFormatter(WithRequestedFields(Request_ShortTime)),
	},
}

var LogLevelTRACE = &loggingLevel{
	name:       TRACE,
	tag:        stdTagTRACE,
	importance: ImportanceTRACE,
	callerInfo: true,
	osExit:     false,
	//writers:    map[string]io.Writer{Stderr: os.Stderr},
	writerFormatterMap: map[string]*formatterExpanded{
		Stderr: NewFormatter(WithRequestedFields(Request_ShortTime)),
	},
}

var LogLevelPING = &loggingLevel{
	name:       PING,
	tag:        stdTagPing,
	importance: ImportancePING,
	callerInfo: true,
	osExit:     false,
	//writers:    map[string]io.Writer{Stderr: os.Stderr},
	//formatFunc: formatPing,
}

func defaultLoggingLevels() map[string]*loggingLevel {
	levels := make(map[string]*loggingLevel)
	levels[FATAL] = LogLevelFATAL
	levels[ERROR] = LogLevelERROR
	levels[WARN] = LogLevelWARN
	levels[INFO] = LogLevelINFO
	levels[DEBUG] = LogLevelDEBUG
	levels[TRACE] = LogLevelTRACE
	return levels
}
