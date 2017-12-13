package logger

import (
	//	"fmt"

	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"time"
	// "github.com/AnimusPEXUS/gosignal"
)

type (
	EntryType uint
)

const (
	TextEntryType EntryType = iota
	InfoEntryType
	WarningEntryType
	ErrorEntryType
)

type OutputOptions struct {
	TextIcon       string
	InfoIcon       string
	WarningIcon    string
	ErrorIcon      string
	InsertTime     bool
	TimeLayout     string
	ClosedByLogger bool
}

var std_output_opt = &OutputOptions{
	TextIcon:       "",
	InfoIcon:       "[i]",
	WarningIcon:    "[w]",
	ErrorIcon:      "[e]",
	InsertTime:     true,
	TimeLayout:     time.RFC3339,
	ClosedByLogger: false,
}

type WriterWrapper struct {
	out interface{}
	opt *OutputOptions
}

type LogEntry struct {
	Type EntryType
	Time time.Time
	Text string
}

func (self *LogEntry) TypeString() string {
	switch self.Type {
	case TextEntryType:
		{
			return "text"
		}
	case InfoEntryType:
		{
			return "info"
		}
	case WarningEntryType:
		{
			return "warning"
		}
	case ErrorEntryType:
		{
			return "error"
		}
	}
	return "programming error"
}

func (self *LogEntry) TypeStringT() string {
	return strings.Title(self.TypeString())
}

type LoggerCallback func(*LogEntry, *Logger)

type LoggerI interface {
	Text(string)
	Info(string)
	Warning(string)
	Error(interface{})
}

type Logger struct {
	callbacks      []LoggerCallback
	outputs        map[uint64]*WriterWrapper
	output_counter uint64
	mutex          *sync.Mutex

	stdout_lbl *LineByLineWriter
	stderr_lbl *LineByLineWriter
}

func New() *Logger {
	ret := new(Logger)
	ret.mutex = new(sync.Mutex)
	ret.ResetOutput()
	return ret
}

func (self *Logger) Close() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	for _, val := range self.outputs {
		switch val.out.(type) {
		case io.WriteCloser:
			{
				if val.opt.ClosedByLogger {
					val.out.(io.WriteCloser).Close()
				}
			}
		}
	}
	self.resetOutput()
}

func (self *Logger) ConnectCallback(callback LoggerCallback) {
	self.callbacks = append(self.callbacks, callback)
}

func (self *Logger) addOutputOpt(out interface{}, opts *OutputOptions) uint64 {
	switch out.(type) {
	case io.Writer:
	case io.WriteCloser:
	default:
		panic("only io.Writer or io.WriteCloser may be passed")
	}
	ret := self.output_counter
	self.outputs[ret] = &WriterWrapper{out, opts}
	self.output_counter++
	return ret
}

func (self *Logger) AddOutputOpt(out interface{}, opts *OutputOptions) uint64 {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.addOutputOpt(out, opts)
}

func (self *Logger) AddOutput(out interface{}) uint64 {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.addOutputOpt(out, std_output_opt)
}

func (self *Logger) DelOutput(id uint64) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	delete(self.outputs, id)
}

func (self *Logger) resetOutput() {
	self.outputs = make(map[uint64]*WriterWrapper)
}

func (self *Logger) ResetOutput() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.resetOutput()
}

func (self *Logger) PutEntry(type_ EntryType, value interface{}) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	value_str := "error"

	switch value.(type) {
	case string:
		value_str = value.(string)
	case error:
		value_str = value.(error).Error()
	default:
		value_str = reflect.ValueOf(value).String()
	}

	log_entry := &LogEntry{type_, time.Now().UTC(), value_str}

	for _, cb := range self.callbacks {
		cb(log_entry, self)
	}

	wg := &sync.WaitGroup{}

	wg.Add(len(self.outputs))
	for _, i := range self.outputs {
		go func(
			i *WriterWrapper,
			wg *sync.WaitGroup,
		) {
			self._WriteOutput(i, log_entry)
			wg.Done()
		}(i, wg)
	}

	wg.Wait()
}

func (self *Logger) _WriteOutput(
	ww *WriterWrapper,
	entry *LogEntry,
) {

	txt := ""
	switch entry.Type {
	case TextEntryType:
		txt += ww.opt.TextIcon
	case InfoEntryType:
		txt += ww.opt.InfoIcon
	case WarningEntryType:
		txt += ww.opt.WarningIcon
	case ErrorEntryType:
		txt += ww.opt.ErrorIcon
	}

	if ww.opt.InsertTime {
		if len(txt) != 0 {
			txt += " "
		}
		txt += fmt.Sprintf(
			"%-30s",
			entry.Time.Format(ww.opt.TimeLayout),
		)
	}

	if len(entry.Text) != 0 {
		if len(txt) != 0 {
			txt += " "
		}
		txt += entry.Text
	}

	txt += "\n"

	b := []byte(txt)

	ww.out.(io.Writer).Write(b)
	switch ww.out.(type) {
	case interface {
		Sync() error
	}:
		ww.out.(interface {
			Sync() error
		}).Sync()
	}

}

func (self *Logger) Text(txt string) {
	self.PutEntry(TextEntryType, txt)
}

func (self *Logger) Info(txt string) {
	self.PutEntry(InfoEntryType, txt)
}

func (self *Logger) Warning(txt string) {
	self.PutEntry(WarningEntryType, txt)
}

func (self *Logger) Error(value interface{}) {
	self.PutEntry(ErrorEntryType, value)
}

func (self *Logger) StdoutLbl() *LineByLineWriter {
	if self.stdout_lbl == nil {
		self.stdout_lbl = NewLineByLineWriter(self, false)
	}
	return self.stdout_lbl
}

func (self *Logger) StderrLbl() *LineByLineWriter {
	if self.stderr_lbl == nil {
		self.stderr_lbl = NewLineByLineWriter(self, true)
	}
	return self.stderr_lbl
}
