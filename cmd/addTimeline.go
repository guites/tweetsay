/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"tweetsay/database"
	"tweetsay/helper"
	"tweetsay/twiapi"

	"github.com/spf13/cobra"
)

// addTimelineCmd represents the addTimeline command
var addTimelineCmd = &cobra.Command{
	Use:   "addTimeline @user_handle",
	Short: "Adds a user timeline to the pool",
	Long: `Adds and starts tracking a user timeline. This means that their tweets can be toggled and drawn on new terminal windows.
	
New tweets will also be fetched from the API and added to the pool.`,
	Args: func (cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("usage: add_timeline @user_handle")
		}
		return nil
	},
Run: func(cmd *cobra.Command, args []string) {
		client := helper.GetTwitterClient()

		username := helper.RemoveAt(args[0])
		
		_, db_err := database.GetUser(username)
		
		if db_err == nil {
			fmt.Printf("User @%s already being tracked!\n", username)
			os.Exit(0)
		}
		
		twitter_user := twiapi.GetByUsername(username, client)

		user, user_err := database.AddUser(twitter_user)

		if user_err != nil {
			log.Fatal("Error while adding user to database", user_err)
		}

		fmt.Printf("Account: @%s (%s) (%d)\n", user.User.ScreenName, user.User.Name, user.User.ID)

		tweets := twiapi.GetUserTimeline(&user, client)
		
		for _, tweet := range tweets {
			database.AddTweet(&tweet)
		}

		fmt.Printf("Indexed %d tweets from user @%s\n", len(tweets), username)
	},
}

func init() {
	rootCmd.AddCommand(addTimelineCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addTimelineCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addTimelineCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
