package main

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const PREFIX = "0000"
const N = len(PREFIX)
const VERBOSE = false
const THREADS = 3

func unquote(inp string) string {
	str, _ := strconv.Unquote(`"` + inp + `"`)
	return str
}

func generatePrivateKey() *ecdsa.PrivateKey {
	privateKey, _ := crypto.GenerateKey() // Generate a new private key
	return privateKey
}

func generateAddress(privateKey *ecdsa.PrivateKey) string {
	publicKey := privateKey.Public()                         // Generate the public key from the private key
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)        // Convert it to ECDSA format
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex() // Find the address from the public key
	return address
}

func logPow(x float32, n int) float32 {
	if n == 0 {
		return 1
	}
	if n == 1 {
		return x
	}
	y := logPow(x, n/2)
	if n%2 == 0 {
		return y * y
	}
	return x * y * y
}

func speedTest(tests int, delta time.Duration) {
	fmt.Println("It took", delta, "to generate", tests, "addresses")
	fmt.Println("Prefix is ", N, "characters long")
	// one_iter = logPow((1 / 16), N)
}

// todo
// * estimation calculation
// * generate only letters abcdef
// * generate only numbers 0123456789
// * suffix

func search(thread int, wg *sync.WaitGroup) {
	fmt.Println("Thread", thread, "started")
	printer := message.NewPrinter(language.English)
	INIT_T := time.Now()

	for i := 0; ; /*loop til I say don't*/ i++ {
		var privateKey = generatePrivateKey()
		var address = generateAddress(privateKey)

		var slice = address[2 : N+2] // account for 0x prefix
		if slice == PREFIX {
			fmt.Println("=== Found it! ===")
			commas := printer.Sprintf("%d", i)
			fmt.Println(unquote("Iteration:\t"), commas)
			fmt.Println(unquote("Time:\t\t"), time.Since(INIT_T))
			fmt.Println(unquote("Private key:\t"), hexutil.Encode(crypto.FromECDSA(privateKey)))
			fmt.Println(unquote("Address:\t"), address)
			wg.Done()
			break
		} else if VERBOSE && i%(30*477) == 0 {
			fmt.Println("Searching...")
			fmt.Println(unquote("> Slice:\t"), slice)
			commas := printer.Sprintf("%d", i)
			fmt.Println(unquote("> Iteration:\t"), commas)
			fmt.Println(unquote("> Time:\t\t"), time.Since(INIT_T))
			speedTest(i, time.Since(INIT_T))
		}
	}
}

func main() {
	_, err := strconv.ParseInt(PREFIX, 16, 64)
	if err != nil {
		fmt.Println("Prefix '" + PREFIX + "' is not valid hexadecimal!")
		os.Exit(128)
	}

	var wg sync.WaitGroup
	wg.Add(THREADS)

	for thread := 0; thread < THREADS; thread++ {
		go search(thread, &wg)
		defer wg.Done()
	}
	wg.Wait()

}
