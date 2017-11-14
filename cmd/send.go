// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	b64 "encoding/base64"
	"fmt"
	"github.com/google/go-github/github"
	helper "github.com/ianmcmahon/encoding_ssh"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"log"
	"os"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("send called")
		send()
	},
}

func init() {
	RootCmd.AddCommand(sendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func send() {
	/* baseURL, err := url.Parse(os.Getenv("GITHUB_API"))
	if err != nil {
		panic(err)
	} */

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		panic("Missing GITHUB_TOKEN environment variable")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	keyOpts := &github.ListOptions{}

	keys, _, err := client.Users.ListKeys(ctx, "mlbright", keyOpts)
	if err != nil {
		log.Fatalf("Error listing keys: %v", err)
	}

	var wormholeKey github.Key
	for _, k := range keys {
		key, _, _ := client.Users.GetKey(ctx, k.GetID())
		if key.GetTitle() == "wormhole" {
			wormholeKey = *key
		}
	}

	log.Println(wormholeKey.GetTitle())

	secret := "abcdef"

	something, err := helper.DecodePublicKey(wormholeKey.GetKey())
	if err != nil {
		log.Fatalf("Could decode public key: %v", err)
	}
	publicKey := something.(*rsa.PublicKey)

	var out []byte
	out, err = rsa.EncryptOAEP(sha1.New(), rand.Reader, publicKey, []byte(secret), []byte("test secret"))
	if err != nil {
		log.Fatalf("encrypting error: %s", err)
	}
	sEnc := b64.StdEncoding.EncodeToString(out)

	file := github.GistFile{}
	file.Content = &sEnc

	files := make(map[github.GistFilename]github.GistFile)
	files["something"] = file

	g := &github.Gist{}
	g.Files = files
	gist, _, err := client.Gists.Create(ctx, g)
	if err != nil {
		log.Fatalf("Error creating gist: %v", err)
	}
	fmt.Println(gist.GetHTMLURL())

}
