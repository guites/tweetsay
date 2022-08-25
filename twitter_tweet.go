package main

import (
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
)

// fetch user tweets from API and save to database
func get_user_timeline(user *twitter.User, client *twitter.Client) () {

	if get_user_timeline_from_db(user) {
		return
	}
	log.Printf("Fetching user @%s tweets from Twitter API\n", user.ScreenName)
	tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		UserID: user.ID,
    	Count: 3200,
		TweetMode: "extended",
	})
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	for _, tweet := range tweets {
		// fmt.Printf("%s\n", tweet.FullText)
		add_tweet_to_db(&tweet)
	}
	return
}