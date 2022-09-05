package twiapi

import (
	"fmt"
	"log"
	"tweetsay/model"

	"github.com/dghubble/go-twitter/twitter"
)

// fetch user tweets from API and save to database
func GetUserTimeline(user *model.User, client *twitter.Client) ([]twitter.Tweet) {

	fmt.Printf("Fetching user @%s tweets from Twitter API\n", user.User.ScreenName)
	tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		UserID: user.User.ID,
    	Count: 3200,
		TweetMode: "extended",
	})
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	return tweets
}

func GetByUsername(username string, client *twitter.Client) (*twitter.User){

	fmt.Printf("Fetching user @%s from Twitter API\n", username)
	twitter_user, _, err := client.Users.Show(&twitter.UserShowParams{
		ScreenName: username,
	})

	if err != nil {
		log.Fatalf("err: %s\n", err)
	}

	return twitter_user
}