package ulog
import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	Mb             int = 1024000
	DefaultMaxSize     = 100
	_depth             = 9
)

type Formatter struct {
	TimestampFormat string
}

var strScanID = "NONE"

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC1123Z
	}

	_, file, line, ok := runtime.Caller(_depth)
	var strFile string
	if !ok {
		strFile = "???:?"
	} else {
		strFile = fmt.Sprint(filepath.Base(file), ":", line)
	}
	tTime := time.Now()
	t1 := tTime.Unix()
	t2 := tTime.UnixNano()
	t3 := t2 - t1*1000000000
	strTime := tTime.Format("2006-01-02 15:04:05")
	var serialized string

	if strScanID == "NONE" {
		serialized = fmt.Sprintf("{\"time\":\"%s%d\",\"level\":\"%s\",\"msg\":\"%s\",\"filename\":\"%s\"}",
			strTime, t3, entry.Level.String(), entry.Message, strFile)
	} else {
		serialized = fmt.Sprintf("{\"time\":\"%s%d\",\"level\":\"%s\",\"scanID\":\"%s\",\"msg\":\"%s\",\"filename\":\"%s\"}",
			strTime, t3, entry.Level.String(), strScanID, entry.Message, strFile)
	}

	return append([]byte(serialized), '\n'), nil
}

type FileWriter struct {
	lock        sync.Mutex
	fileName    string
	currentData string
	fp          *os.File
}

func NewFileWriter(fileName string) (*FileWriter, error) {
	fileName = fileName + "." + time.Now().Format("20060102") + ".log"
	w := &FileWriter{fileName: fileName}
	var err error
	w.fp, err = os.OpenFile(w.fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	w.currentData = time.Now().Format("20060102")
	return w, err
}

func (w *FileWriter) Rotate() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	fmt.Println("Rotate")
	if w.fp != nil {
		fmt.Println("fp close")
		err = w.fp.Close()
		w.fp = nil
		if err != nil {
			fmt.Println("err!=nil")
			return
		}
	}

	array := strings.SplitN(w.fileName, ".20", -1)

	strFileName := array[0]

	strFileName = strFileName + "." + w.currentData + ".log"
	w.fp, err = os.OpenFile(strFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	return
}

type LogHook struct {
	fileWriter *FileWriter
	maxSize    int
}

func SetScanID(scanID string) {
	strScanID = scanID
}

func DisAbleScanID() {
	SetScanID("NONE")
}

func NewLogHook(fileName string, size ...int) (*LogHook, error) {
	fw, err := NewFileWriter(fileName)
	if err != nil {
		return nil, err
	}
	maxSize := DefaultMaxSize
	if len(size) > 0 {
		maxSize = size[0]
	}
	logHook := &LogHook{
		fileWriter: fw,
		maxSize:    maxSize,
	}
	return logHook, nil
}

func InitLog(logfile string, loglevel string) error {
	if logfile == "" {
		return errors.New("logfile is empty")
	}

	hook, err := NewLogHook(logfile)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	logrus.AddHook(hook)

	level, err := logrus.ParseLevel(loglevel)
	if err != nil {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(level)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	return nil
}

func (hook *LogHook) Fire(entry *logrus.Entry) error {
	formatter := Formatter{}
	line, err := formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	if hook.fileWriter != nil {
		fileInfo, err := hook.fileWriter.fp.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "1Unable to get file info, %v", err)
			return err
		}
		nowDay := time.Now().Format("20060102")
		if fileInfo.Size()+int64(len(line)) > int64(hook.maxSize*Mb) || hook.fileWriter.currentData != nowDay {
			fmt.Println("nowData:", nowDay)
			fmt.Println("oldData:", hook.fileWriter.currentData)
			hook.fileWriter.currentData = nowDay
			if err := hook.fileWriter.Rotate(); err != nil {
				fmt.Println("rotate err")
				fmt.Fprintf(os.Stderr, "Unable to rotate, %v", err)
				return err
			}
		}
		if _, err = hook.fileWriter.fp.Write(line); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to write data, %v", err)
			return err
		}
	}

	return nil

}

func (hook *LogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

