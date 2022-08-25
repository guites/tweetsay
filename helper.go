package main

import (
	"envholder"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// remove @ if user included in command line arg
func removeAt(username string) (string) {
	if username[0:1] == "@" {
		username = username[1:]
	}
	return username
}

func prepare_twitter_api() (*twitter.Client){
	var envVars envholder.EnvHolder = envholder.LoadEnv()
	consumerKey := envVars.GetVar("TWITTER_API_KEY")
	consumerSecret := envVars.GetVar("TWITTER_API_KEY_SECRET")
	accessToken := envVars.GetVar("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := envVars.GetVar("TWITTER_ACCESS_TOKEN_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	return client
}