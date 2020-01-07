// +build windows

package signal

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// make pid file and log pid
func LogPid() {
	pid := os.Getpid()
	log.Printf("windows pid [%d]", pid)
	file, err := os.OpenFile("pid", os.O_CREATE|os.O_TRUNC,0)
	if err != nil {
		log.Fatalln("pid file open error : ", err)
	}
	_, err = file.WriteString(strconv.Itoa(pid))
	if err != nil {
		log.Fatalln("pid write error : ", err)
	}
	err = file.Close()
	if err != nil {
		log.Fatalln("pid file close error : ", err)
	}
}

// remove pid file
func RemovePid() error {
	return os.Remove("pid")
}

// get pid from pid file
func Pid() (int, error) {
	r, err := ioutil.ReadFile("pid")
	if err != nil {
		return -1, err
	}
	i, err := strconv.Atoi(string(r))
	if err != nil {
		return -1, err
	}
	return i, nil
}

// stop process of pid
// send kill signal to referred process
func StopProcess() error {
	pid, err := Pid()
	if err != nil {
		return err
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	//Sending Interrupt on Windows is not implemented
	err = process.Signal(syscall.SIGKILL)
	if err != nil {
		return err
	}
	return nil
}

func ProcessSignal(stop chan bool, usr1 chan bool, usr2 chan bool) {
	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGINT, syscall.SIGKILL)
	for {
		msg := <-s
		switch msg {
		case syscall.SIGINT, syscall.SIGKILL:
			log.Printf("windows signal [%s]", msg)
			if stop != nil {
				stop <- true
			}
			break
		case nil:
			break
		default:
			log.Printf("windows signal [%s]", msg)
			break
		}
	}
}