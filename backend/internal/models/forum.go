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
	Deleted      bool      `json:"-" bson:"deleted"`
}

type Comment struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	DiscussionID string    `json:"-" bson:"discussion_id"`
	RelatedTo    string    `json:"-" bson:"related_to"`
	Content      string    `json:"content" bson:"content"`
	AuthorID     int       `json:"author_id" bson:"author_id"`
	Likes        []int     `json:"-" bson:"likes"`
	LikesCount   int       `json:"likes" bson:"-"`
	Dislikes     []int     `json:"-" bson:"dislikes"`
	DisikesCount int       `json:"dislikes" bson:"-"`
	Edited       bool      `json:"edited" bson:"edited"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	Deleted      bool      `json:"-" bson:"deleted"`
	Children     []Comment `json:"children,omitempty" bson:"-"`
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
