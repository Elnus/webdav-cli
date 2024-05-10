package chksum

import (
	"crypto/md5"
	"log"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

func chkMD5(data []byte) []byte {
	chunkSize := runtime.NumCPU()
	hash := md5.New()
	for i := 0; i < chunkSize; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			defer func() {
				recover()
			}()
			start := index * len(data) / chunkSize
			end := (index + 1) * len(data) / chunkSize
			_, err := hash.Write(data[start:end])
			if err != nil {
				log.Panic("hash failed", err)
			}
		}(i)
	}
	wg.Wait()
	return hash.Sum(nil)
}
