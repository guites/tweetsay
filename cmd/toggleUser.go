/*
Copyright © 2022 guites <gui.garcia67@gmail.com>

*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"tweetsay/database"
	"tweetsay/helper"

	"github.com/spf13/cobra"
)

// toggleUserCmd represents the toggleUser command
var toggleUserCmd = &cobra.Command{
	Use:   "toggleUser @user_handle",
	Short: "Toggles whether user tweets should be drawn on new terminal windows",
	Long: `Toggles whether user tweets should be drawn on new terminal windows.

The user should be previously tracked by using the addTimeline command.`,
	Args: func (cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("user handle is required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		username := helper.RemoveAt(args[0])
		user, err := database.GetUser(username)
		if err != nil {
			fmt.Printf("User @%s not registered in database. Please run <tweetsay addTimeline @%s>\n", username, username)
			os.Exit(1)
		}
		database.ToggleUser(user)
	},
}

func init() {
	rootCmd.AddCommand(toggleUserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// toggleUserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// toggleUserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
