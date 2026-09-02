package main

import (
	"fmt"
	"os"
)

func main() {
	path := os.TempDir() + "/wal_demo.log"
	// The log is opened fresh to simulate the state after a restart.
	os.Remove(path)

	first := []LogEntry{{Key: "a", Val: "1"}, {Key: "b", Val: "2"}}
	if err := AppendLog(path, first); err != nil {
		fmt.Println("append error:", err)
		return
	}
	fmt.Println("append ok")

	// Simulate the process being restarted: a fresh AppendLog call on the
	// same file must append without corrupting prior records.
	second := []LogEntry{{Key: "c", Val: "3"}, {Key: "d", Val: "4"}}
	if err := AppendLog(path, second); err != nil {
		fmt.Println("append error:", err)
		return
	}

	entries, err := ReadLog(path)
	if err != nil {
		fmt.Println("read error:", err)
		return
	}
	fmt.Printf("log entries: %d\n", len(entries))

	fmt.Println("replay:")
	for _, kv := range Replay(entries) {
		fmt.Printf("%s=%s\n", kv.Key, kv.Val)
	}

	os.Remove(path)
}