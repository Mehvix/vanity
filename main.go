package main

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const PREFIX = "000000"
const VERBOSE = true

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

// todo
// * estimation calculation
// * generate only letters abcdef
// * generate only numbers 0123456789
// * suffix

func main() {
	_, err := strconv.ParseInt(PREFIX, 16, 64)
	if err != nil {
		fmt.Println("Prefix '" + PREFIX + "' is not valid hexadecimal!")
		os.Exit(128)
	}

	p := message.NewPrinter(language.English)
	INIT_T := time.Now()

	for i := 0; ; i++ {
		var privateKey = generatePrivateKey()
		var address = generateAddress(privateKey)

		var slice = address[2 : len(PREFIX)+2] // account for 0x prefix
		if slice == PREFIX {
			fmt.Println("=== Found it! ===")
			commas := p.Sprintf("%d", i)
			fmt.Println(unquote("Iteration:\t"), commas)
			fmt.Println(unquote("Time:\t\t"), time.Since(INIT_T))
			fmt.Println(unquote("Private key:\t"), hexutil.Encode(crypto.FromECDSA(privateKey)))
			fmt.Println(unquote("Address:\t"), address)
			break
		} else if VERBOSE && i%(30*477) == 0 {
			fmt.Println("Searching...")
			fmt.Println(unquote("> Slice:\t"), slice)
			commas := p.Sprintf("%d", i)
			fmt.Println(unquote("> Iteration:\t"), commas)
			fmt.Println(unquote("> Time:\t\t"), time.Since(INIT_T))
		}
	}
}
