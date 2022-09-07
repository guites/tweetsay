/*
Copyright © 2022 guites <gui.garcia67@gmail.com>

*/
package cmd

import (
	"fmt"
	"tweetsay/database"

	"github.com/spf13/cobra"
)

// listUsersCmd represents the listUsers command
var listUsersCmd = &cobra.Command{
	Use:   "listUsers",
	Short: "Shows currently tracked twitter accounts",
	Long: `Shows currently tracked twitter accounts alongside the number of indexed tweets and whether they are actively being drawn on new terminal windows.`,
	Run: func(cmd *cobra.Command, args []string) {
		users := database.GetAllUsersWithTweetCount()
		var ActiveStr string
		for _, user := range users {
			if user.Active {
				ActiveStr = "[x]"
			} else {
				ActiveStr = "[ ]"
			}
			fmt.Printf("@%s (%s) %s - %d tweets\n", user.User.ScreenName, user.User.Name, ActiveStr, user.DbTweetCount)
		}
	},
}

func init() {
	rootCmd.AddCommand(listUsersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listUsersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listUsersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
