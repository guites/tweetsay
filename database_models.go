package main

type DatabaseTweets struct {
	ID          int64
	CreatedAt   string
	FullText    string
	UserName    string
}

type DatabaseUsers struct {
	ID         int64
	CreatedAt  string
	ScreenName string
}