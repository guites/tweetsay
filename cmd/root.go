/*
Copyright © 2022 guites <gui.garcia67@gmail.com>

*/
package cmd

import (
	"os"
	"strings"
	"tweetsay/database"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tweetsay",
	Short: "Tweetsay shows you tweets from selected accounts on terminal opening",
	Long: `Tweetsay shows you tweets from selected accounts on terminal opening.
User timelines that are both added and toggled on will be drawn randomly each time a terminal window opens.

Start by tracking an account:
	tweetsay addTimeline @Antho_Repartie
and then toggle it:
	tweetsay toggleUser @Antho_Repartie
get a random tweet from one of your toggled users by running tweetsay without any flags:
	tweetsay
list all the words on the last tweet with a link to their wiktionary page:
	tweetsay listWords
if you disliked the last shown tweet, you can delete it from the pool by running:
	tweetsay deleteLast
you can also check your tracked users status by running:
	tweetsay listUsers

You can add shell completion running
	tweetsay completion
And following instructions.`,
	Run: func(cmd *cobra.Command, args []string) {
		database.CreateTables()
		tweet := database.GetRandomTweet()
		database.SetLastShownTweet(tweet.ID)
		words := strings.Fields(tweet.FullText)
		var wordID int
		for _, word := range words {
			// do not index user handles or urls
			if strings.HasPrefix(word, "@") ||
			strings.HasPrefix(word, "http") {
				continue
			}
			wordID = database.AddWord(word)
			if wordID != -1 {
				database.RelateWordToTweet(wordID, tweet.ID)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tweetsay.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


