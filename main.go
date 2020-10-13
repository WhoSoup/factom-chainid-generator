package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var counter uint64
var result chan []byte
var base [][]byte
var target string

func init() {
	result = make(chan []byte)
}

func calculateID(try []byte) []byte {
	hash := sha256.New()
	for i := range base {
		eidhash := sha256.Sum256(base[i])
		hash.Write(eidhash[:])
	}

	tryhash := sha256.Sum256(try)
	hash.Write(tryhash[:])
	return hash.Sum(nil)
}

func worker() {
	for {
		count := atomic.AddUint64(&counter, 1)
		extid := []byte(strconv.FormatUint(count, 10))
		id := fmt.Sprintf("%x", calculateID(extid))

		if target == id[:len(target)] {
			result <- extid
		}
	}
}

func main() {
	workers := flag.Int("workers", 0, "Number of concurrent worker routines. Default = Maximum = number of cores")
	flag.StringVar(&target, "target", "888888", "The target byte prefix in hex (must be <= 32 bytes)")
	baseS := flag.String("base", "", "The chain id base extids (as strings), comma separated")
	flag.Parse()

	if len(target) > 64 || len(target) < 1 {
		fmt.Println("invalid target (must be at least one character and at most 64)")
		os.Exit(0)
	}
	target = strings.ToLower(target)
	targetRE := regexp.MustCompile("^[a-f0-9]+$")
	if !targetRE.MatchString(target) {
		fmt.Println("invalid target (must be hexadecimal)")
		os.Exit(0)
	}

	split := strings.Split(*baseS, ",")
	for _, s := range split {
		s = strings.TrimSpace(s)
		base = append(base, []byte(s))
	}
	fmt.Println("Using base:", split)
	fmt.Println("Using target:", target)

	start := time.Now()

	if *workers > 0 && *workers < runtime.NumCPU() {
		runtime.GOMAXPROCS(*workers)
	} else {
		*workers = runtime.NumCPU()
	}

	for i := 0; i < *workers; i++ {
		go worker()
	}

	res := <-result
	fmt.Printf("Chain found with extid %s: %x (in %s)\n", res, calculateID(res), time.Since(start))
	fmt.Println("ExtIDs:", append(split, string(res)))
}
