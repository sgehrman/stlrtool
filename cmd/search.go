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
	"strings"
	"sync"
	"time"

	"github.com/stellar/go/keypair"
)

var searches = 0

// searchCmd represents the vanity command
var searchCmd *cobra.Command

func init() {
	searchCmd = &cobra.Command{
		Use:   "search",
		Short: "Search for lost private key",
		Long:  `Search for lost private key.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Searching...")
			searchForKey(args)
		},
	}

	rootCmd.AddCommand(searchCmd)
}

func searchForKey(args []string) {
	if len(args) != 1 {
		usage2()
		os.Exit(1)
	}

	matchArg := strings.ToUpper(args[0])
	fmt.Printf("Looking for: %s\n", matchArg)

	var wg sync.WaitGroup

	wg.Add(1)

	for i := 0; i < 40; i++ {
		go search2(&wg, matchArg)
	}

	wg.Wait()
}

func search2(wg *sync.WaitGroup, matchArg string) {
	defer wg.Done()

	t0 := time.Now()

	for {
		searches++

		if searches%500000 == 0 {
			fmt.Printf("searches: %v\n", searches)
		}

		kp, err := keypair.Random()

		if err != nil {
			log.Fatal(err)
		}

		if kp.Address() == matchArg {
			t1 := time.Now()
			fmt.Printf("Search took %v seconds. searches: %v\n", t1.Sub(t0).Seconds(), searches)

			fmt.Printf("Secret seed: %s\n", kp.Seed())
			fmt.Printf("Public: %s\n", kp.Address())
			os.Exit(0)
		}
	}
}

func usage2() {
	fmt.Printf("Usage:\n\tstlrtool PREFIX -s (for suffix)\n")
}
