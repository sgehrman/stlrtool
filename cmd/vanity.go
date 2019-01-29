// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/stellar/go/keypair"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"

var tries = 0

var useSuffixP *bool

// vanityCmd represents the vanity command
var vanityCmd *cobra.Command

func init() {
	vanityCmd = &cobra.Command{
		Use:   "vanity",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Searching...")
			vanity(args)
		},
	}

	rootCmd.AddCommand(vanityCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// vanityCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	useSuffixP = vanityCmd.Flags().BoolP("suffix", "s", false, "find suffix instead of prefix")
}

func vanity(args []string) {
	if len(args) != 1 {
		usage()
		os.Exit(1)
	}

	matchArg := strings.ToUpper(args[0])
	checkPlausible(matchArg)

	var wg sync.WaitGroup

	wg.Add(1)

	for i := 0; i < 42; i++ {
		go search(&wg, *useSuffixP, matchArg)
	}

	wg.Wait()
}

func search(wg *sync.WaitGroup, useSuffix bool, matchArg string) {
	defer wg.Done()

	t0 := time.Now()

	for {
		tries++

		if tries%100000 == 0 {
			fmt.Printf("tries: %v\n", tries)
		}

		kp, err := keypair.Random()

		if err != nil {
			log.Fatal(err)
		}

		found := false

		if useSuffix {
			if strings.HasSuffix(kp.Address(), matchArg) {
				found = true
			}
		} else {
			// NOTE: the first letter of an address will always be G, and the second letter will be one of only a few
			// possibilities in the base32 alphabet, so we are actually searching for the vanity value after this 2
			// character prefix.
			if strings.HasPrefix(kp.Address()[2:], matchArg) {
				found = true
			}
		}

		if found {
			t1 := time.Now()
			fmt.Printf("Search took %v seconds. tries: %v\n", t1.Sub(t0).Seconds(), tries)

			fmt.Printf("Secret seed: %s\n", kp.Seed())
			fmt.Printf("Public: %s\n", kp.Address())
			os.Exit(0)
		}
	}
}

func usage() {
	fmt.Printf("Usage:\n\tstlrtool PREFIX -s (for suffix)\n")
}

// aborts the attempt if a desired character is not a valid base32 digit
func checkPlausible(prefix string) {
	for _, r := range prefix {
		if !strings.ContainsRune(alphabet, r) {
			fmt.Printf("Invalid prefix: %s is not in the base32 alphabet\n", strconv.QuoteRune(r))
			os.Exit(1)
		}
	}
}
