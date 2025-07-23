package logman

import (
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
)

// This is a convinience function for ProcessMessage.
// Printf formats message according to a format specifier and writes to output writers of Level INFO.
// It returns message processing error encountered.
// Same as func: Info(format, args...) (error) .
func Printf(format string, args ...interface{}) error {
	msg := NewMessage(format, args...)
	if err := process(msg, logMan.logLevels[INFO]); err != nil {
		return err
	}
	return nil
}

// This is a convinience function for ProcessMessage.
// Println writes to output writers of Level INFO. Args are separated
// with spaces.
// It returns message processing error encountered.
func Println(args ...interface{}) error {
	format := ""
	for range args {
		format += "%v "
	}
	format = strings.TrimSuffix(format, " ")
	msg := NewMessage(format, args...)
	if err := process(msg, logMan.logLevels[INFO]); err != nil {
		return err
	}
	return nil
}

// This is a convinience function for ProcessMessage.
// Fatalf formats message according to a format specifier and writes to output writers of Level FATAL.
// It returns message processing error encountered or error created if processing is success.
// By default calling level Fatal cause os.Exit(1) after completion (subject to change during logger setup process).
func Fatalf(format string, args ...interface{}) error {
	msg := NewMessage(format, args...)
	if err := process(msg, logMan.logLevels[FATAL]); err != nil {
		return err
	}
	return nil
}

// This is a convinience function for ProcessMessage.
// Errorf formats message according to a format specifier and writes to output writers of Level ERROR.
// It returns message processing error encountered or error created if processing is success.
func Errorf(format string, args ...interface{}) error {
	errCreated := fmt.Errorf(format, args...)
	msg := NewMessage(format, args...)
	if errProcessing := process(msg, logMan.logLevels[ERROR]); errProcessing != nil {
		return errProcessing
	}
	return errCreated
}

// This is a convinience function for ProcessMessage.
// Error creates message input argument and writes to output writers of Level ERROR.
// It returns message processing error encountered or input error if processing is success.
func Error(errInput error) error {
	msg := NewMessage(errInput.Error())
	if errProcessing := process(msg, logMan.logLevels[ERROR]); errProcessing != nil {
		return errProcessing
	}
	return errInput
}

// This is a convinience function for ProcessMessage.
// Warn formats message according to a format specifier and writes to output writers of Level WARN.
// It returns message processing error encountered.
func Warn(format string, args ...interface{}) error {
	msg := NewMessage(format, args...)
	if err := process(msg, logMan.logLevels[WARN]); err != nil {
		return err
	}
	return nil
}

// This is a convinience function for ProcessMessage.
// Info formats message according to a format specifier and writes to output writers of Level INFO.
// It returns message processing error encountered.
func Info(format string, args ...interface{}) error {
	msg := NewMessage(format, args...)
	if err := process(msg, logMan.logLevels[INFO]); err != nil {
		return err
	}
	return nil
}

// This is a convinience function for ProcessMessage.
// Debug receives message with additional comments. Message will be written to writers of Level DEBUG.
// Comments will be printed to os.Stderr EVEN if message will not be processed.
// It returns message processing error encountered.
func Debug(msg Message, comments ...string) error {
	for _, comment := range comments {
		comment = color.S256(253).Sprintf("%v", comment)
		fmt.Fprintf(os.Stderr, "#%v\n", comment)
	}
	if err := process(msg, logMan.logLevels[DEBUG]); err != nil {
		return err
	}
	return nil
}

// This is a convinience function for ProcessMessage.
// Trace receives message with  additional comments. Message will be written to writers of Level TRACE.
// Comments will be printed to os.Stderr EVEN if message will not be processed.
// It returns message processing error encountered.
func Trace(msg Message, comments ...string) error {
	for _, comment := range comments {
		comment = color.S256(253).Sprintf("%v", comment)
		fmt.Fprintf(os.Stderr, "#%v\n", comment)
	}
	if err := process(msg, logMan.logLevels[TRACE]); err != nil {
		return err
	}
	return nil
}

// This is a convinience function for ProcessMessage.
// Ping receives comments. Message with code location will be created and written to writers of Level PING.
// Comments will be printed to os.Stderr EVEN if message will not be processed.
// Message processing error encountered will be printed as comment.
//
// Never return error.
func Ping(comments ...string) error {
	msg := NewMessage("")
	if err := process(msg, logMan.logLevels[PING]); err != nil {
		fmt.Fprintf(os.Stderr, "ping error: %v\n", err)
	}
	for _, comment := range comments {
		comment = color.S256(239).Sprintf("%v", comment)
		fmt.Fprintf(os.Stderr, "%v\n", comment)
	}
	return nil
}
