package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pin/tftp/v3"
)

// Hander for read (aka GET) requests.
func readHandler(filename string, rf io.ReaderFrom) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "opening %s: %v\n", filename, err)
		return err
	}
	defer file.Close() // 在函数结束时关闭文件句柄
	n, err := rf.ReadFrom(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reading %s: %v\n", filename, err)
		return err
	}
	fmt.Printf("%d bytes sent\n", n)
	return nil
}

// Handler for write (aka PUT) requests.
func writeHandler(filename string, wt io.WriterTo) error {
	if _, err := os.Stat(filename); err == nil {
		// 文件存在，尝试删除文件
		err := os.Remove(filename)
		if err != nil {
			fmt.Println("无法删除文件:", err)
			return err
		}
		fmt.Println("文件删除成功")
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "creating %s: %v\n", filename, err)
		return err
	}
	defer file.Close() // 在函数结束时关闭文件句柄
	n, err := wt.WriteTo(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "writing %s: %v\n", filename, err)
		return err
	}
	fmt.Printf("%d bytes received\n", n)
	return nil
}

// Hook for logging on every transfer completion or failure.
type logHook struct{}

func (h *logHook) OnSuccess(stats tftp.TransferStats) {
	fmt.Printf("Transfer of %s to %s complete\n", stats.Filename, stats.RemoteAddr)
}
func (h *logHook) OnFailure(stats tftp.TransferStats, err error) {
	fmt.Printf("Transfer of %s to %s failed: %v\n", stats.Filename, stats.RemoteAddr, err)
}

func main() {
	port := flag.Int("p", 69, "Local port to listen")
	flag.Parse()

	// Start the server.
	s := tftp.NewServer(readHandler, writeHandler)
	s.SetHook(&logHook{})
	go func() {
		err := s.ListenAndServe(fmt.Sprintf(":%d", *port))
		if err != nil {
			fmt.Fprintf(os.Stdout, "Can't start the server: %v\n", err)
			os.Exit(1)
		}
	}()

	// Do some other stuff.
	time.Sleep(5000 * time.Minute)

	// Eventually shutdown the server.
	s.Shutdown()
}
