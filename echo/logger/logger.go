package logger

import (
	"encoding/json"
	"fmt"
	"github.com/brickman-source/golang-utilities/log"
	gecholog "github.com/labstack/gommon/log"
	"github.com/mattn/go-colorable"
	"io"
	"os"
)

type (
	Logger struct {
	}
)

func New() (l *Logger) {
	l = &Logger{}
	return
}

func (l *Logger) Prefix() string {
	return ""
}

func (l *Logger) SetHeader(h string) {
	//l.template = l.newTemplate(h)
}

func (l *Logger) SetPrefix(p string) {
	//l.prefix = p
}

func (l *Logger) Level() gecholog.Lvl {
	return gecholog.DEBUG //l.level
}

func (l *Logger) SetLevel(v gecholog.Lvl) {
	//l.level = v
}

func (l *Logger) Output() io.Writer {
	//return l.output
	return colorable.NewColorableStdout()
}

func (l *Logger) SetOutput(w io.Writer) {
	//l.output = w
}

func (l *Logger) Print(i ...interface{}) {
	log.Infof(nil, "%v", i)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	log.Infof(nil, format, args)
}

func (l *Logger) Printj(j gecholog.JSON) {
	b, err := json.Marshal(j)
	if err == nil {
		log.Infof(nil, "%v", string(b))
	}
}

func (l *Logger) Debug(i ...interface{}) {
	log.Debugf(nil, "%v", i)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	log.Debugf(nil, format, args)
}

func (l *Logger) Debugj(j gecholog.JSON) {
	b, err := json.Marshal(j)
	if err == nil {
		log.Debugf(nil, "%v", string(b))
	}
}

func (l *Logger) Info(i ...interface{}) {
	log.Infof(nil, "%v", i)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	log.Infof(nil, format, args)
}

func (l *Logger) Infoj(j gecholog.JSON) {
	b, err := json.Marshal(j)
	if err == nil {
		log.Infof(nil, "%v", string(b))
	}
}

func (l *Logger) Warn(i ...interface{}) {
	log.Warnf(nil, "%v", i)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	log.Warnf(nil, format, args)
}

func (l *Logger) Warnj(j gecholog.JSON) {
	b, err := json.Marshal(j)
	if err == nil {
		log.Warnf(nil, "%v", string(b))
	}
}

func (l *Logger) Error(i ...interface{}) {
	log.Errorf(nil, "%v", i)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	log.Errorf(nil, format, args)
}

func (l *Logger) Errorj(j gecholog.JSON) {
	b, err := json.Marshal(j)
	if err == nil {
		log.Errorf(nil, "%v", string(b))
	}
}

func (l *Logger) Fatal(i ...interface{}) {
	log.Fatalf(nil, "%v", i)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(nil, format, args)
}

func (l *Logger) Fatalj(j gecholog.JSON) {
	b, err := json.Marshal(j)
	if err == nil {
		log.Fatalf(nil, "%v", string(b))
	}
	os.Exit(1)
}

func (l *Logger) Panic(i ...interface{}) {
	log.Errorf(nil, "%v", i)
	panic(fmt.Sprint(i...))
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	log.Errorf(nil, format, args)
	panic(fmt.Sprintf(format, args))
}

func (l *Logger) Panicj(j gecholog.JSON) {
	b, err := json.Marshal(j)
	if err == nil {
		log.Errorf(nil, "%v", string(b))
	}
	panic(j)
}

//
//func Prefix() string {
//	return ""
//}
//
//func SetPrefix(p string) {
//	""
//}
//
//func Level() Lvl {
//	return global.Level()
//}
//
//func SetLevel(v Lvl) {
//	global.SetLevel(v)
//}
//
//func Output() io.Writer {
//	return global.Output()
//}
//
//func SetOutput(w io.Writer) {
//	global.SetOutput(w)
//}
//
//func SetHeader(h string) {
//	global.SetHeader(h)
//}
//
//func Print(i ...interface{}) {
//	global.Print(i...)
//}
//
//func Printf(format string, args ...interface{}) {
//	global.Printf(format, args...)
//}
//
//func Printj(j gecholog.JSON) {
//	global.Printj(j)
//}
//
//func Debug(i ...interface{}) {
//	global.Debug(i...)
//}
//
//func Debugf(format string, args ...interface{}) {
//	global.Debugf(format, args...)
//}
//
//func Debugj(j gecholog.JSON) {
//	global.Debugj(j)
//}
//
//func Info(i ...interface{}) {
//	global.Info(i...)
//}
//
//func Infof(format string, args ...interface{}) {
//	global.Infof(format, args...)
//}
//
//func Infoj(j gecholog.JSON) {
//	global.Infoj(j)
//}
//
//func Warn(i ...interface{}) {
//	global.Warn(i...)
//}
//
//func Warnf(format string, args ...interface{}) {
//	global.Warnf(format, args...)
//}
//
//func Warnj(j gecholog.JSON) {
//	global.Warnj(j)
//}
//
//func Error(i ...interface{}) {
//	global.Error(i...)
//}
//
//func Errorf(format string, args ...interface{}) {
//	global.Errorf(format, args...)
//}
//
//func Errorj(j gecholog.JSON) {
//	global.Errorj(j)
//}
//
//func Fatal(i ...interface{}) {
//	global.Fatal(i...)
//}
//
//func Fatalf(format string, args ...interface{}) {
//	global.Fatalf(format, args...)
//}
//
//func Fatalj(j gecholog.JSON) {
//	global.Fatalj(j)
//}
//
//func Panic(i ...interface{}) {
//	global.Panic(i...)
//}
//
//func Panicf(format string, args ...interface{}) {
//	global.Panicf(format, args...)
//}
//
//func Panicj(j gecholog.JSON) {
//	global.Panicj(j)
//}
//
//func (l *Logger) log(v Lvl, format string, args ...interface{}) {
//	l.mutex.Lock()
//	defer l.mutex.Unlock()
//	buf := l.bufferPool.Get().(*bytes.Buffer)
//	buf.Reset()
//	defer l.bufferPool.Put(buf)
//	_, file, line, _ := runtime.Caller(3)
//
//	if v >= l.level || v == 0 {
//		message := ""
//		if format == "" {
//			message = fmt.Sprint(args...)
//		} else if format == "json" {
//			b, err := json.Marshal(args[0])
//			if err != nil {
//				panic(err)
//			}
//			message = string(b)
//		} else {
//			message = fmt.Sprintf(format, args...)
//		}
//
//		_, err := l.template.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
//			switch tag {
//			case "time_rfc3339":
//				return w.Write([]byte(time.Now().Format(time.RFC3339)))
//			case "time_rfc3339_nano":
//				return w.Write([]byte(time.Now().Format(time.RFC3339Nano)))
//			case "level":
//				return w.Write([]byte(l.levels[v]))
//			case "prefix":
//				return w.Write([]byte(l.prefix))
//			case "long_file":
//				return w.Write([]byte(file))
//			case "short_file":
//				return w.Write([]byte(path.Base(file)))
//			case "line":
//				return w.Write([]byte(strconv.Itoa(line)))
//			}
//			return 0, nil
//		})
//
//		if err == nil {
//			s := buf.String()
//			i := buf.Len() - 1
//			if s[i] == '}' {
//				// gecholog.JSON header
//				buf.Truncate(i)
//				buf.WriteByte(',')
//				if format == "json" {
//					buf.WriteString(message[1:])
//				} else {
//					buf.WriteString(`"message":`)
//					buf.WriteString(strconv.Quote(message))
//					buf.WriteString(`}`)
//				}
//			} else {
//				// Text header
//				buf.WriteByte(' ')
//				buf.WriteString(message)
//			}
//			buf.WriteByte('\n')
//			l.output.Write(buf.Bytes())
//		}
//	}
//}
