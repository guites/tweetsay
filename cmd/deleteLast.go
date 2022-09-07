/*
Copyright © 2022 guites <gui.garcia67@gmail.com>

*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"tweetsay/database"

	"github.com/spf13/cobra"
)

// deleteLastCmd represents the deleteLast command
var deleteLastCmd = &cobra.Command{
	Use:   "deleteLast",
	Short: "Delete the last shown tweet from the pool",
	Long: `Deletes the last shown tweet from the pool by setting its SoftDeleted flag to True.
The tweet will therefore not be drawn on new terminal windows.
You will be asked to confirm the action by typing 'delete'.`,
	Run: func(cmd *cobra.Command, args []string) {
		tweet := database.GetLastShownTweet()

		fmt.Printf("Delete tweet from @%s?\n",tweet.User.ScreenName)
		if len(tweet.FullText) > 30 {
			fmt.Println(tweet.FullText[:30],"...")
		} else {
			fmt.Println(tweet.FullText,"...")
		}
	
		var confirmDelete string
		fmt.Print("Please type 'delete' to confirm removal:")
		fmt.Scanf("%s", &confirmDelete)
		if strings.ToLower(confirmDelete) != "delete" {
			fmt.Println("Tweet was not deleted.")
			os.Exit(0)
		}

		database.DeleteTweet(tweet.ID)

		fmt.Println("Tweet was deleted successfully.")
},
}

func init() {
	rootCmd.AddCommand(deleteLastCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteLastCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteLastCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
