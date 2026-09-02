package main

import (
	"fmt"
	"sync"
)

func main() {
	m := NewMetrics()
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Inc("req_total")
			m.Inc("req_ok")
			m.Inc("conn_total")
		}()
	}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Inc("conn_total")
			m.Inc("conn_ok")
		}()
	}
	wg.Wait()
	m.Add("total", 19)

	fmt.Println("snapshot (sorted):")
	for _, kv := range m.Snapshot() {
		fmt.Printf("%s=%s\n", kv.Key, kv.Val)
	}
	a, _ := m.Get("req_total")
	b, _ := m.Get("conn_total")
	fmt.Printf("total a+b=%d\n", a+b)
}