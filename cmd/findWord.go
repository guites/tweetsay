/*
Copyright © 2022 guites <gui.garcia67@gmail.com>

*/
package cmd

import (
	"errors"
	"fmt"
	"tweetsay/database"

	"github.com/spf13/cobra"
)

// findWordCmd represents the findWord command
var findWordCmd = &cobra.Command{
	Use:   "findWord word",
	Short: "Find tweets containing the given word",
	Long: `Find tweets containing the given word amongst the tweets you have seen so far.`,
	Args: func (cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please include the word you want to search for")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		tweetExcerpts := database.FindTweetsContainingWord(args[0])
		for _, excerpt := range tweetExcerpts {
			fmt.Printf("@%s (%d): %s\n", excerpt.Username, excerpt.TweetID, excerpt.TweetExcerpt)
		}
	},
}

func init() {
	rootCmd.AddCommand(findWordCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// findWordCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// findWordCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
