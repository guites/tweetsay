/*
Copyright © 2022 guites <gui.garcia67@gmail.com>

*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"tweetsay/database"
	"unicode/utf8"

	"github.com/spf13/cobra"
)

// listWordsCmd represents the listWords command
var listWordsCmd = &cobra.Command{
	Use:   "listWords",
	Short: "Lists all words on the last shown tweet with a link to their wiktionary page",
	Long: `Lists all words on the last shown tweet with a link to their wiktionary page.

The output is formatted as a list, and the wiktionary language prefix (ie. fr, en, pt, de) will be inferred from the tweet language.`,
	Run: func(cmd *cobra.Command, args []string) {
		tweet := database.GetLastShownTweet()
		// formats wiktionary url based on tweet Lang
		var wiktionary_url string
		if tweet.Lang == "und" ||
		tweet.Lang == "qme" || 
		tweet.Lang == "" {
			wiktionary_url = "https://wiktionary.org/wiki/"
		} else {
			wiktionary_url =  "https://" + tweet.Lang + ".wiktionary.org/wiki/"
		}

		words := strings.Fields(tweet.FullText)

		var max_length int
		for _, word := range words {
			curr_length := utf8.RuneCountInString(word)
			if curr_length > max_length {
				max_length = curr_length
			}
		}

		for index, word := range words {
			if word == "!" ||
			word == "?" ||
			word == "-" ||
			word == ":" {
				continue
			}
			wiktionary_link := wiktionary_url + word
			if strings.HasPrefix(word, "http") {
				wiktionary_link = ""
			}
			wrapped_word := word
			if utf8.RuneCountInString(word) < max_length {
				wrapped_word = word + strings.Repeat(" ", max_length - utf8.RuneCountInString(word))
			}
			padded_index := "0"+strconv.Itoa(index + 1)
			fmt.Println(padded_index[len(padded_index)-2:], wrapped_word, wiktionary_link)
		}
	},
}

func init() {
	rootCmd.AddCommand(listWordsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listWordsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listWordsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
