package log

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap/zapcore"
	"runtime"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func NewNormalCaller() logrus.Hook {
	return &callerHook{}
}

type callerHook struct {
}

func (hook *callerHook) Fire(entry *logrus.Entry) error {
	entry.Data["caller"] = hook.caller(entry)
	return nil
}

func (hook *callerHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		//logrus.DebugLevel,
	}
}

func (hook *callerHook) caller(entry *logrus.Entry) string {
	return string(stack(8))
	//if _, file, line, ok := runtime.Caller(7); ok {
	//    return strings.Join([]string{file, strconv.Itoa(line)}, ":")
	//} else {
	//    return ""
	//}
}

func stack(skip int) []byte {
	buf := new(bytes.Buffer)
	//var lines [][]byte
	//var lastFile string
	for i := skip; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		//fmt.Fprintf(buf, "%s:%d ", file, line)
		//if file != lastFile {
		//    data, err := ioutil.ReadFile(file)
		//    if err != nil {
		//        continue
		//    }
		//    lines = bytes.Split(data, []byte{'\n'})
		//    lastFile = file
		//}
		//line--
		//fmt.Fprintf(buf, "%s:%s:%d\n", function(pc), file, line)
		//file = strings.Replace(file, "$HOME/work/go/", "", -1)
		fmt.Fprintf(buf, "%s : %d ", file, line)
	}
	return buf.Bytes()
}

//func source(lines [][]byte, n int) []byte {
//    if n < 0 || n >= len(lines) {
//        return dunno
//    }
//    return bytes.Trim(lines[n], " \t")
//}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())

	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func zapCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(string(stack(8)))
}
