package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Peersyst/xrpl-go/pkg/crypto"
	"github.com/Peersyst/xrpl-go/xrpl/wallet"
)

const LogStep = 100000

type walletResult struct {
	Address string
	Seed    string
}

func generateWallets(ctx context.Context, prefix, suffix string, isCaseSensitive bool, resultChan chan<- walletResult, attempts, foundCount *int64, total int, wg *sync.WaitGroup, startTime time.Time) {
	defer wg.Done()

	if !isCaseSensitive {
		prefix = strings.ToLower(prefix)
		suffix = strings.ToLower(suffix)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			newWallet, err := wallet.New(crypto.ED25519())
			if err != nil {
				continue
			}

			addr := string(newWallet.ClassicAddress)
			currentAttempts := atomic.AddInt64(attempts, 1)

			if currentAttempts%LogStep == 0 {
				currentFound := atomic.LoadInt64(foundCount)
				elapsed := time.Since(startTime).Truncate(time.Second)
				fmt.Printf("⏳ Checked %d keypairs in %s. Found %d / %d\n", currentAttempts, elapsed, currentFound, total)
			}

			compareAddr := addr
			if !isCaseSensitive {
				compareAddr = strings.ToLower(addr)
			}

			if (prefix == "" || strings.HasPrefix(compareAddr, prefix)) &&
				(suffix == "" || strings.HasSuffix(compareAddr, suffix)) {
				atomic.AddInt64(foundCount, 1)
				resultChan <- walletResult{Address: addr, Seed: newWallet.Seed}
				return
			}
		}
	}
}

func main() {
	prefix := flag.String("starts-with", "", "Prefix to search for (e.g., 'rLOL')")
	suffix := flag.String("ends-with", "", "Suffix to search for (e.g., 'xyz')")
	count := flag.Int("count", 1, "Number of matching wallets to find")
	threads := flag.Int("threads", runtime.NumCPU(), "Number of goroutines to use")
	isCaseSensitive := flag.Bool("is-case-sensitive", false, "Match prefix/suffix with exact case")
	flag.Parse()

	if *prefix == "" && *suffix == "" {
		fmt.Println("❌ You must provide at least one of --starts-with or --ends-with")
		os.Exit(1)
	}

	if *count < 1 {
		fmt.Println("❌ Count must be at least 1")
		os.Exit(1)
	}

	if *threads > 1000 {
		*threads = 1000
	} else if *threads < 1 {
		*threads = runtime.NumCPU()
	}

	fmt.Printf("🚀 Searching with %d goroutines for %d wallet(s) matching:\n", *threads, *count)
	if *prefix != "" {
		fmt.Printf("  🔍 Prefix: %q\n", *prefix)
	}
	if *suffix != "" {
		fmt.Printf("  🔍 Suffix: %q\n", *suffix)
	}
	fmt.Printf("  🔠 Case-sensitive: %v\n", *isCaseSensitive)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultChan := make(chan walletResult, *threads)
	var attempts int64
	var foundCount int64
	var wg sync.WaitGroup
	startTime := time.Now()

	for i := 0; i < *threads; i++ {
		wg.Add(1)
		go generateWallets(ctx, *prefix, *suffix, *isCaseSensitive, resultChan, &attempts, &foundCount, *count, &wg, startTime)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	found := 0
	for res := range resultChan {
		found++
		fmt.Printf("\n🎯 Wallet %d found!\n", found)
		fmt.Printf("  📍 Address: %s\n", res.Address)
		fmt.Printf("  🔑 Seed: %s\n", res.Seed)

		if found >= *count {
			cancel()
			break
		}
	}

	fmt.Printf("\n🔁 Total attempts: %d\n", atomic.LoadInt64(&attempts))
	fmt.Printf("⏰  Elapsed time: %s\n", time.Since(startTime).Truncate(time.Second))
}
