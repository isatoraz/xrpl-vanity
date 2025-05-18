package main

import (
	"context"
	"flag"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Peersyst/xrpl-go/pkg/crypto"
	"github.com/Peersyst/xrpl-go/xrpl/wallet"
)

func main() {
	wallets := flag.Int("wallets", 1_000_000, "Number of wallets to generate for speed test")
	threads := flag.Int("threads", runtime.NumCPU(), "Number of goroutines to use")
	flag.Parse()

	fmt.Printf("🚀 Benchmarking wallet generation with %d attempts using %d threads...\n", *wallets, *threads)

	var wg sync.WaitGroup
	var counter int64
	start := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < *threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if atomic.LoadInt64(&counter) >= int64(*wallets) || ctx.Err() != nil {
					return
				}
				_, err := wallet.New(crypto.ED25519())
				if err != nil {
					continue
				}
				atomic.AddInt64(&counter, 1)
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start).Seconds()
	fmt.Printf("✅ Generated %d wallets in %.2f seconds\n", *wallets, elapsed)
	fmt.Printf("⚡ Speed: %.2f wallets/sec\n", float64(*wallets)/elapsed)
	return
}
