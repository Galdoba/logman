package logman

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Galdoba/logman/colorizer"
)

type formatterExpanded struct {
	fieldFormaFuncMap map[string]func(Message, Colorizer) (string, error)
	requestedFields   []string
	writerKey         string
	colorizer         Colorizer
	customColorizer   bool
}

func NewFormatter(options ...FormatterOption) *formatterExpanded {
	fe := formatterExpanded{}
	fe.fieldFormaFuncMap = make(map[string]func(Message, Colorizer) (string, error))
	opts := defaultFormatterOptions()
	for _, set := range options {
		set(&opts)
	}
	fe.requestedFields = opts.requestFields
	fe.fieldFormaFuncMap = opts.formatFuncs
	fe.colorizer = opts.colorizer
	fe.customColorizer = opts.customColorizer
	return &fe
}

func basicFormatter(field string, val interface{}) (string, error) {
	format := "%v=%v"
	return fmt.Sprintf(format, field, val), nil
}

func (fe *formatterExpanded) Format(msg Message, color bool) string {
	output := ""
	for _, field := range fe.requestedFields {

		val := msg.Value(field)

		formatted := ""
		err := errors.New("not formatted")
		fn := fe.fieldFormaFuncMap[field]
		switch fn {
		case nil:
			formatted, err = basicFormatter(field, val)
			if err != nil {
				return output + formatted
			}
		default:
			if field == "json" {
				text, err := stdJSON(msg, nil)
				if err != nil {
					return err.Error()
				}
				return text
			}
			// if field == keyCallerLong {
			// 	field = keyCaller
			// }
			// if field == keyCallerShort {
			// 	field = keyCaller
			// }

			//if msg.Value(field) == nil {
			//fmt.Println("no field", field)
			//	continue
			//}
			switch color {
			case true:
				formatted, err = fn(msg, fe.colorizer)
			case false:
				formatted, err = fn(msg, nil)
			}
			if err != nil {
				return output + formatted + "!<> " + err.Error()
			}
		}
		output += formatted + " "

	}
	return output
}

func mustIgnore(fields []string, field string) bool {
	for _, f := range fields {
		if f == field {
			return true
		}
	}
	return false
}

func (fe *formatterExpanded) AddFormatterFunc(field string, fn func(Message, Colorizer) (string, error)) {
	fe.fieldFormaFuncMap[field] = fn
}

type formatterOptions struct {
	requestFields   []string
	formatFuncs     map[string]func(Message, Colorizer) (string, error)
	colorizer       Colorizer
	customColorizer bool
}

type FormatterOption func(*formatterOptions)

func defaultFormatterOptions() formatterOptions {
	return formatterOptions{
		requestFields: Request_ShortSince,
		formatFuncs: map[string]func(Message, Colorizer) (string, error){
			keyTime:        stdFormatFunc_time,
			keySince:       stdFormatFunc_since,
			keyLevel:       stdFormatLevel,
			keyMessage:     stdFormatMessage,
			keyCallerShort: stdFormatCallerShort,
			keyCallerLong:  stdFormatCallerLong,
		},
		colorizer: nil,
	}
}

// WithRequestedFields - Add fields for writers of this level to request from message.
// Valid Requests are: keyTime, keySince, keyLevel, keyMessage, keyCallerLong, keyCallerLong
func WithRequestedFields(fields []string) FormatterOption {
	return func(fo *formatterOptions) {
		fo.requestFields = fields
	}
}

// WithColor - sets color scheme to use for this formatter.
// nil = no color
// scheme set by this func will not be overwriten by logman's global option.
func WithColor(color Colorizer) FormatterOption {
	return func(fo *formatterOptions) {
		fo.colorizer = color
		fo.customColorizer = true
	}
}

var Request_MessageOnly = []string{keyMessage}
var Request_ShortTime = []string{keyTime, keyLevel, keyMessage}
var Request_ShortSince = []string{keySince, keyLevel, keyMessage}
var Request_ShortReport = []string{keySince, keyMessage}
var Request_Medium = []string{keyTime, keyLevel, keyMessage, keyCallerShort}
var Request_Full = []string{keyTime, keySince, keyLevel, keyMessage, keyCallerLong}

func WithCustomFunc(requestKey string, fn func(Message, Colorizer) (string, error)) FormatterOption {
	return func(fo *formatterOptions) {
		fo.formatFuncs[requestKey] = fn
	}
}

func stdFormatFunc_time(msg Message, colors Colorizer) (string, error) {
	tm, err := validateTimeArg(msg.Value("time"))
	if err != nil {
		return "", err
	}
	text := formatTime(tm)
	// if len(text) < 21 {
	// 	text += "0"
	// }
	switch colors {
	case nil:
	default:
		level := fmt.Sprintf("%v", msg.Value(keyLevel))
		text = colors.ColorizeByKeys(text, colorizer.NewKey(colorizer.FG_KEY, level))
	}
	return fmt.Sprintf("[%v]", text), nil
}

func formatTime(tm time.Time) string {
	s := tm.Format("2006-01-02 15:04:05.999")
	slice := strings.Split(s, "")
	switch len(slice) {
	case 19:
		s += ".000"
	case 20:
		s += "000"
	case 21:
		s += "00"
	case 22:
		s += "0"
	case 23:
		return s
	}
	return s
}

func stdFormatFunc_since(msg Message, colors Colorizer) (string, error) {

	duration := time.Since(logMan.startTime)
	//fmt.Println("to log", duration)
	switch colors {
	case nil:
		return fmt.Sprintf("[%.3f]", float64(duration.Milliseconds())/1000), nil
	default:
		text := fmt.Sprintf("%.3f", float64(duration.Milliseconds())/1000)
		if text == "00000" {
			text = "0.000"
		}
		level := fmt.Sprintf("%v", msg.Value(keyLevel))
		text = colors.ColorizeByKeys(text, colorizer.NewKey(colorizer.FG_KEY, level))
		return fmt.Sprintf("[%v]", text), nil
	}

}

func validateTimeArg(args ...any) (time.Time, error) {
	if len(args) != 1 {
		return time.Time{}, fmt.Errorf("stdTimeFormat function expect 1 argument (have %v)", len(args))
	}
	val := args[0]
	str := fmt.Sprintf("%v", val)
	str = strings.TrimPrefix(str, "[")
	str = strings.TrimSuffix(str, "]")
	tm, err := time.Parse(time.RFC3339Nano, fmt.Sprintf("%v", str))
	if err != nil {
		return time.Time{}, err
	}
	return tm, nil
}

func stdFormatMessage(msg Message, colors Colorizer) (string, error) {
	inputs := msg.InputArgs()
	var format string
	var args []interface{}
	for i := -1; i < len(inputs)-1; i++ {
		switch i {
		case -1:
			format = fmt.Sprintf("%v", inputs[-1])
		default:
			args = append(args, inputs[i])
		}
	}
	switch colors {
	case nil:
		return fmt.Sprintf(format, args...), nil
	default:
		var coloredArgs []string
		for _, arg := range args {
			colored := colors.ColorizeByType(arg)
			coloredArgs = append(coloredArgs, colored)
		}
		text := combineColored(format, coloredArgs...)
		return fmt.Sprintf(text), nil
	}
}

func stdFormatLevel(msg Message, colors Colorizer) (string, error) {
	level := msg.Value(keyLevel)
	if level == nil {
		return "", errNoField(keyLevel)
	}
	tag := fmt.Sprintf("[%v]", level)
	switch colors {
	case nil:
		return tag, nil
	default:
		keyFg := colorizer.NewKey(colorizer.FG_KEY, fmt.Sprintf("%v", level))
		keyBg := colorizer.NewKey(colorizer.BG_KEY, fmt.Sprintf("%v", level))
		return fmt.Sprintf("[%v]", colors.ColorizeByKeys(level, keyFg, keyBg)), nil
	}
}

func errNoField(field string) error {
	return fmt.Errorf("no field with key '%v'", field)
}

func stdFormatCallerLong(msg Message, colors Colorizer) (string, error) {
	file := msg.Value(keyFile)
	line := msg.Value(keyLine)
	funk := msg.Value(keyFunc)
	if file == nil {
		return "", errNoField(keyFile)
	}
	if line == nil {
		return "", errNoField(keyLine)
	}
	text := fmt.Sprintf("\n  [caller=%v:%v]", file, line)
	if funk != nil {
		text += fmt.Sprintf(" [func=%v]", funk)
	}
	switch colors {
	case nil:
		return fmt.Sprintf("%v", text), nil
	default:
		keyFg := colorizer.NewKey(colorizer.FG_KEY, "caller")
		keyBg := colorizer.NewKey(colorizer.BG_KEY, "caller")
		return fmt.Sprintf("%v", colors.ColorizeByKeys(text, keyFg, keyBg)), nil
	}
}

func stdFormatCallerShort(msg Message, colors Colorizer) (string, error) {
	file := msg.Value(keyFile)
	line := msg.Value(keyLine)
	if file == nil {
		return "", errNoField(keyFile)
	}
	if line == nil {
		return "", errNoField(keyLine)
	}
	fileStr := filepath.Base(fmt.Sprintf("%v", file))
	text := fmt.Sprintf("\n  [caller=%v:%v]", fileStr, line)
	switch colors {
	case nil:
		return fmt.Sprintf("%v", text), nil
	default:
		keyFg := colorizer.NewKey(colorizer.FG_KEY, "caller")
		keyBg := colorizer.NewKey(colorizer.BG_KEY, "caller")
		return fmt.Sprintf("%v", colors.ColorizeByKeys(text, keyFg, keyBg)), nil
	}
}

type JSONlog struct {
	APP   string            `json:"app"`
	LVL   string            `json:"level"`
	MSG   string            `json:"message"`
	TIME  string            `json:"time"`
	AGRS1 map[string]string `json:"logman keys,omitempty"`
	AGRS2 map[string]string `json:"input arguments,omitempty"`
}

func stdJSON(msg Message, color Colorizer) (string, error) {
	jsMsg := JSONlog{}
	jsMsg.APP = logMan.appName
	jsMsg.TIME = fmt.Sprintf("%v", msg.Value(keyTime))
	jsMsg.LVL = fmt.Sprintf("%v", msg.Value(keyLevel))
	msgText, err := stdFormatMessage(msg, nil)
	if err != nil {
		return "", err
	}
	jsMsg.MSG = msgText
	keys := msg.Fields()
	for _, key := range keys {
		switch key {
		case keyTime, keyLevel, keyMessage:
		default:
			switch jsMsg.LVL {
			case ERROR, FATAL, DEBUG, TRACE:
				if jsMsg.AGRS1 == nil {
					jsMsg.AGRS1 = make(map[string]string)
				}
				jsMsg.AGRS1[key] = fmt.Sprintf("%v", msg.Value(key))
			}

		}
	}

	inputArgs := msg.InputArgs()
	switch jsMsg.LVL {
	case FATAL, TRACE:
		if len(inputArgs) > 0 {
			jsMsg.AGRS2 = make(map[string]string)
			for i := 0; i < len(inputArgs)-1; i++ {
				jsMsg.AGRS2[fmt.Sprintf("arg[%v]", i)] = fmt.Sprintf("%v", inputArgs[i])

			}
		}
	}
	bt, err := json.Marshal(&jsMsg)
	if err != nil {
		return "", err
	}
	return string(bt), nil
}
