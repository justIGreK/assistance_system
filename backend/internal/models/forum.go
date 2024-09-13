package models

import "time"

type Discussion struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Title     string    `json:"title" bson:"title"`
	Content   string    `json:"content" bson:"content"`
	AuthorID  int    `json:"author_id" bson:"author_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type Comment struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	DiscussionID string    `json:"discussion_id" bson:"discussion_id"`
	Content      string    `json:"content" bson:"content"`
	AuthorID     int    `json:"author_id" bson:"author_id"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
}

type DiscussionTopic struct {
    ID      string `bson:"_id"`
    Title   string `bson:"title"`
    Content string `bson:"content"`
}

type DiscussionWithCount struct {
    Discussion    DiscussionTopic
    CommentsCount int64
}