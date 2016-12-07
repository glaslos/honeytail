package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hpcloud/tail"
)

var (
	fileName   string
	followFile bool
	reopenFile bool
	fileWhence int
)

func process() {
	t, err := tail.TailFile(fileName,
		tail.Config{
			Follow: followFile,
			ReOpen: reopenFile,
			Location: &tail.SeekInfo{
				Offset: 0,
				Whence: fileWhence,
			},
		})
	if err != nil {
		log.Println(err)
		return
	}
	for line := range t.Lines {
		// Process log line
		fmt.Println(line.Text)
	}
}

func main() {
	flag.StringVar(&fileName, "file-name", "/var/log/syslog", "File to tail")
	flag.BoolVar(&followFile, "follow-file", true, "Follow the file, tailing new lines")
	flag.BoolVar(&reopenFile, "reopen-file", true, "Re-open the file if renamed (logrotate)")
	flag.IntVar(&fileWhence, "file-whence", 2, "0: start, 1: current, 2:end")
	flag.Parse()
	process()
}
