package logger

import (
	"fmt"
	"github.com/shirou/gopsutil/v4/process"
	"log"
	"os"
	"sync"
	"template/pgk/utils"
	"time"
)

const (
	directory  = "logs/"
	applPrefix = directory + "appl"
	logFormat  = ".log"
	currentLog = directory + "appl_current.log"
)

var instance Logger
var once sync.Once

type logger struct {
	logLevel string
	file     *os.File
	logg     *log.Logger
}

type Logger interface {
	Write(...interface{})
	RenameLog() error
	WriteStatisticToLog(args map[int]int)
	WriteCpuInfoToLog([]float64)
	WriteMemoryInfoToLog(*process.MemoryInfoStat)
	GetLogFile() *os.File
	GetLogLevel() string
	SetLogLevel(logLevel string) error
}

func LoggerInstance() Logger {
	once.Do(func() {
		instance = NewLogger()
	})

	return instance
}

func NewLogger() Logger {
	// create if directory doesnt exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err = os.Mkdir(directory, 0755); err != nil {
			log.Fatal(err)
		}
	}

	file, err := os.OpenFile(currentLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	logger := &logger{
		file: file,
		logg: log.New(file, "", log.LstdFlags),
	}

	return logger
}

func (log *logger) Write(args ...any) {
	if args == nil {
		log.logg.Print("no data")
		return
	}

	log.logg.Print(fmt.Sprint(args...))
}

func (log *logger) RenameLog() error {
	if err := log.file.Close(); err != nil {
		return err
	}

	if err := os.Rename(currentLog, applPrefix+time.Now().Format("2006.01.02_150405")+logFormat); err != nil {
		return err
	}

	return nil
}

func (log *logger) WriteStatisticToLog(args map[int]int) {
	for k, v := range args {
		log.Write("---------------------------------------------------------------------")
		log.Write("<STAT> ", "source", "\t | \t ", "count", "\t|")
		log.Write("<STAT> ", k, "\t\t | \t ", v, "\t\t|")
		log.Write("---------------------------------------------------------------------")
	}
}

func (log *logger) WriteCpuInfoToLog(cpuData []float64) {
	log.Write("---------------------------------------------------------------------")
	log.Write("<STAT>\t", "source\t", "thread num", "\t|\t ", "percent used", "\t|\t", "percent free", "\t|")
	for k, v := range cpuData {
		log.Write("<STAT>\t", "CPU\t\t\t", k, "\t\t|\t\t", utils.RoundTo(v, 2), "\t\t|\t\t", utils.RoundTo(100-v, 2), "\t\t|")
	}
	log.Write("---------------------------------------------------------------------")
}

func (log *logger) WriteMemoryInfoToLog(memInfo *process.MemoryInfoStat) {
	log.Write("<STAT>\t", "source", "\t", "rss MB used", "\t|\t\t", "vms MB used", "\t|")
	log.Write("<STAT>\t", "MEM\t\t\t", memInfo.RSS/1024/1024, "\t\t|\t\t\t", memInfo.VMS/1024/1024, "\t\t|")
}

func (log *logger) GetLogFile() *os.File {
	return log.file
}

func (log *logger) GetLogLevel() string {
	return log.logLevel
}

// from config file
func (log *logger) SetLogLevel(logLevel string) error {
	// change to ENUM
	switch logLevel {
	case "info":
		log.logLevel = "info"
	case "debug":
		log.logLevel = "debug"
	case "warn":
		log.logLevel = "warn"
	case "error":
		log.logLevel = "error"
	default:
		log.logLevel = "info"
	}
	return nil
}
