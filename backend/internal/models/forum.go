package models

import "time"

type Discussion struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	Title        string    `json:"title" bson:"title"`
	Content      string    `json:"content" bson:"content"`
	AuthorID     int       `json:"author_id" bson:"author_id"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	Likes        []int     `json:"-" bson:"likes"`
	LikesCount   int       `json:"likes" bson:"-"`
	Dislikes     []int     `json:"-" bson:"dislikes"`
	DisikesCount int       `json:"dislikes" bson:"-"`
	Edited       bool      `json:"edited" bson:"edited"`
}

type Comment struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	DiscussionID string    `json:"discussion_id" bson:"discussion_id"`
	Content      string    `json:"content" bson:"content"`
	AuthorID     int       `json:"author_id" bson:"author_id"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	Likes        []int     `json:"-" bson:"likes"`
	LikesCount   int       `json:"likes" bson:"-"`
	Dislikes     []int     `json:"-" bson:"dislikes"`
	DisikesCount int       `json:"dislikes" bson:"-"`
	Edited       bool      `json:"edited" bson:"edited"`
}
type DiscussionTopic struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Title        string `json:"title" bson:"title"`
	Content      string `json:"content" bson:"content"`
	Likes        []int  `json:"-" bson:"likes"`
	LikesCount   int    `json:"likes" bson:"-"`
	Dislikes     []int  `json:"-" bson:"dislikes"`
	DisikesCount int    `json:"dislikes" bson:"-"`
}

type DiscussionWithCount struct {
	Discussion    DiscussionTopic
	CommentsCount int64
}
