package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
)

func main() {
	threads := runtime.NumCPU()

	t1 := time.Now()
	var count uint64 = 0
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				res := generate()
				count++
				if strings.HasPrefix(res.address, "0xp0817") {
					fmt.Printf("Mnemonic:%s\r\nAddress:%s", res.mnemonic, res.address)
					os.Exit(0)
				}
			}
		}()
	}
	go func() {
		for {
			elapsed := time.Since(t1)
			c := atomic.LoadUint64(&count)
			rate := 0
			if elapsed.Seconds() != 0 {
				rate = int(c / uint64(elapsed.Seconds()))
			}
			log.Printf("rate:%v/s,time elapsed:%v", rate, elapsed)
			time.Sleep(10 * time.Second)
		}
	}()
	wg.Wait()
}

func generate() walletdata {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		log.Fatal(err)
	}

	mnemonic, _ := bip39.NewMnemonic(entropy)
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}

	return walletdata{
		mnemonic: mnemonic,
		address:  account.Address.Hex(),
	}
}

type walletdata struct {
	mnemonic string
	address  string
}
