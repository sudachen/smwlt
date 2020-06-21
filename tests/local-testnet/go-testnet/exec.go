package go_testnet

import (
	"bufio"
	"fmt"
	"github.com/sudachen/smwlt/fu/errstr"
	"os"
	exec2 "os/exec"
	"syscall"
)

func stringify(opts map[string]interface{}) []string {
	a := []string{}
	for k,v := range opts {
		if v != nil {
			k = k + "=" + fmt.Sprint(v)
		}
		a = append(a,k)
	}
	return a
}

func exec(path string, opts map[string]interface{}, sigterm chan struct{}) (chan string, error) {
	c := make(chan string,1)
	fmt.Println(path, stringify(opts))
	cmd := exec2.Command(path, stringify(opts)...)
	stdout, err := cmd.StdoutPipe()
	scan := bufio.NewScanner(stdout)
	if err != nil { return nil, errstr.Wrapf(1,err, "failed to create stdout pipe: %v", err.Error())}
	stderr, err := cmd.StderrPipe()
	if err != nil { return nil, errstr.Wrapf(1,err, "failed to create stderr pipe: %v", err.Error())}
	escan := bufio.NewScanner(stderr)
	go func() {
		for escan.Scan() {
			fmt.Fprintln(os.Stderr,cmd.ProcessState.String(),scan.Text())
		}
	}()
	go func() {
		for scan.Scan() {
			//fmt.Fprintln(os.Stderr,cmd.Process.Pid,path,scan.Text())
			c <- scan.Text()
		}
		close(c)
		// take result from OS, will not create zombie!
		_ = cmd.Wait()
	}()
	if err := cmd.Start(); err != nil {
		return nil, errstr.Wrapf(1,err, "failed to run command: %v", err.Error())
	}
	go func() {
		<- sigterm
		fmt.Println("KILLED",path,opts)
		_ = cmd.Process.Signal(syscall.SIGTERM)
	}()
	return c, nil
}
