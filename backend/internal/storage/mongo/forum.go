package mongo

import (
	"context"
	"errors"
	"fmt"
	"gohelp/internal/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ForumStorage struct {
	discussions *mongo.Collection
	comments    *mongo.Collection
	client      *mongo.Client
}

func NewForumStorage(db *mongo.Database, client *mongo.Client) *ForumStorage {
	return &ForumStorage{
		discussions: db.Collection("discussions"),
		comments:    db.Collection("comments"),
		client: client,
	}
}

func (s *ForumStorage) CreateDiscussion(ctx context.Context, discussion *models.Discussion) (string, error) {
	discussion.CreatedAt = time.Now()
	if discussion.Likes == nil {
		discussion.Likes = []int{}
	}
	if discussion.Dislikes == nil {
		discussion.Dislikes = []int{}
	}
	discussion.Deleted = false
	res, err := s.discussions.InsertOne(ctx, discussion)
	if err != nil {
		return "", err
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func (s *ForumStorage) GetDiscussion(ctx context.Context, id string) (*models.Discussion, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}
	var discussion models.Discussion
	err = s.discussions.FindOne(ctx, bson.M{
		"$and": []bson.M{
			{"_id": oid},
			{"deleted": false},
		},
	}).Decode(&discussion)
	if err != nil {
		return nil, err
	}

	discussion.LikesCount = len(discussion.Likes)
	discussion.DisikesCount = len(discussion.Dislikes)

	return &discussion, nil
}

func (s *ForumStorage) GetComment(ctx context.Context, id string) (*models.Comment, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	var comments models.Comment
	err = s.comments.FindOne(ctx, bson.M{
		"$and": []bson.M{
			{"_id": oid},
			{"deleted": false},
		},
	}).Decode(&comments)
	if err != nil {
		return nil, err
	}

	comments.LikesCount = len(comments.Likes)
	comments.DisikesCount = len(comments.Dislikes)

	return &comments, nil
}

func (s *ForumStorage) GetAllDiscussions(ctx context.Context) ([]models.DiscussionTopic, error) {
	var discussions []models.DiscussionTopic
	filter := bson.M{"deleted": false}
	cursor, err := s.discussions.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &discussions); err != nil {
		return nil, err
	}
	for i := range discussions {
		discussions[i].LikesCount = len(discussions[i].Likes)
		discussions[i].DisikesCount = len(discussions[i].Dislikes)
	}
	return discussions, nil
}

func (s *ForumStorage) CreateComment(ctx context.Context, comment *models.Comment) (string, error) {
	comment.CreatedAt = time.Now()
	if comment.Likes == nil {
		comment.Likes = []int{}
	}
	if comment.Dislikes == nil {
		comment.Dislikes = []int{}
	}
	comment.Deleted = false
	fmt.Println(comment)
	res, err := s.comments.InsertOne(ctx, comment)
	if err != nil {
		return "", err
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func (s *ForumStorage) GetSummaryOfDiscussions(ctx context.Context, discussions []models.DiscussionTopic) ([]models.DiscussionWithCount, error) {
	var result []models.DiscussionWithCount

	for _, discussion := range discussions {
		log.Println(discussion.ID)
		count, err := s.comments.CountDocuments(context.TODO(), bson.M{
			"$and": []bson.M{
				{"discussion_id": discussion.ID},
				{"deleted": false},
			},
		})
		if err != nil {
			return nil, err
		}

		result = append(result, models.DiscussionWithCount{
			Discussion:    discussion,
			CommentsCount: count,
		})
	}
	return result, nil
}

func (s *ForumStorage) GetCommentsByDiscussion(ctx context.Context, discussionID string) ([]models.Comment, error) {
	var comments []models.Comment
	cursor, err := s.comments.Find(context.TODO(), bson.M{"discussion_id": discussionID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &comments); err != nil {
		return nil, err
	}
	for i := range comments {
		comments[i].LikesCount = len(comments[i].Likes)
		comments[i].DisikesCount = len(comments[i].Dislikes)
	}
	return comments, nil
}

func (s *ForumStorage) SearchDiscussionsByName(ctx context.Context, searchTerm string) ([]models.Discussion, error) {

	filter := bson.M{
		"$text": bson.M{
			"$search": searchTerm,
		},
	}
	opts := options.Find().SetSort(bson.D{{"score", bson.M{"$meta": "textScore"}}})

	cursor, err := s.discussions.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, fmt.Errorf("error finding discussions: %v", err)
	}
	defer cursor.Close(context.TODO())

	var discussions []models.Discussion
	if err := cursor.All(context.TODO(), &discussions); err != nil {
		return nil, fmt.Errorf("error decoding discussions: %v", err)
	}

	return discussions, nil
}

func (s *ForumStorage) RemoveVote(userID int, discussionID string, voteType string) error {
	oid, err := primitive.ObjectIDFromHex(discussionID)
	if err != nil {
		return errors.New("invalid discussionID")
	}
	filter := bson.M{"_id": oid}
	removeVote := bson.M{
		"$pull": bson.M{
			"likes":    userID,
			"dislikes": userID,
		},
	}
	_, err = s.discussions.UpdateOne(context.TODO(), filter, removeVote)
	if err != nil {
		return err
	}
	return nil
}
func (s *ForumStorage) DiscAddVote(userID int, discussionID string, voteType string) error {
	oid, err := primitive.ObjectIDFromHex(discussionID)
	if err != nil {
		return errors.New("invalid discussionID")
	}
	filter := bson.M{"_id": oid}
	var update bson.M
	if voteType == "like" {
		fmt.Println("\nlike")
		update = bson.M{
			"$push": bson.M{
				"likes": userID,
			},
		}
	} else if voteType == "dislike" {
		fmt.Println("dislike")
		update = bson.M{
			"$push": bson.M{
				"dislikes": userID,
			},
		}
	}

	result, err := s.discussions.UpdateOne(context.TODO(), filter, update)
	fmt.Println(result)
	if err != nil {
		return err
	}
	return nil
}

func (s *ForumStorage) ComAddVote(userID int, commentID string, voteType string) error {
	oid, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return errors.New("invalid discussionID")
	}
	filter := bson.M{"_id": oid}
	var update bson.M
	if voteType == "like" {
		fmt.Println("\nlike")
		update = bson.M{
			"$push": bson.M{
				"likes": userID,
			},
		}
	} else if voteType == "dislike" {
		fmt.Println("dislike")
		update = bson.M{
			"$push": bson.M{
				"dislikes": userID,
			},
		}
	}

	result, err := s.comments.UpdateOne(context.TODO(), filter, update)
	fmt.Println(result)
	if err != nil {
		return err
	}
	return nil
}

func (s *ForumStorage) UpdateDiscussion(ctx context.Context, discussionID, content string) error {
	oid, err := primitive.ObjectIDFromHex(discussionID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}

	update := bson.M{
		"$set": bson.M{
			"content": content,
			"edited":  true,
		},
	}

	result, err := s.discussions.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	log.Println(result)

	return nil
}

func (s *ForumStorage) UpdateComment(ctx context.Context, commentID, content string) error {
	oid, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}

	update := bson.M{
		"$set": bson.M{
			"content": content,
			"edited":  true,
		},
	}

	result, err := s.comments.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	log.Println(result)

	return nil
}

func (s *ForumStorage) DeleteFullDiscussion(ctx context.Context, discussionID string) error {
	oid, err := primitive.ObjectIDFromHex(discussionID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}

	update := bson.M{
		"$set": bson.M{
			"deleted": true,
		},
	}
	

	result, err := s.discussions.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	log.Println(result)

	filter = bson.M{"discussion_id": discussionID}

	result, err = s.comments.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	log.Println(result)


	// session, err := s.client.StartSession()
	// if err != nil {
	// 	return err
	// }
	// defer session.EndSession(ctx)
	// err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
	// 	if err := session.StartTransaction(); err != nil {
	// 		return err
	// 	}
	// 	update := bson.M{
	// 		"$set": bson.M{
	// 			"deleted": true,
	// 		},
	// 	}
	// 	_, err := s.discussions.UpdateOne(ctx, bson.M{"_id": oid}, update)
	// 	if err != nil {
	// 		session.AbortTransaction(sc)
	// 		return err
	// 	}
	// 	_, err = s.comments.UpdateMany(sc, bson.M{"discussion_id": oid}, update)
	// 	if err != nil {
	// 		session.AbortTransaction(sc)
	// 		return err
	// 	}

	// 	if err := session.CommitTransaction(sc); err != nil {
	// 		return err
	// 	}

	// 	return nil
	// })

	// if err != nil {
	// 	log.Println("Transaction failed:", err)
	// 	return err
	// }


	return nil
}

func (s *ForumStorage) DeleteComment(ctx context.Context, commentID string) error {
	oid, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}

	update := bson.M{
		"$set": bson.M{
			"deleted": true,
		},
	}

	result, err := s.comments.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	log.Println(result)

	return nil
}

func (s *ForumStorage) DeleteAllComments(ctx context.Context, userID int) error {
	filter := bson.M{"author_id": userID}

	update := bson.M{
		"$set": bson.M{
			"deleted": true,
		},
	}

	result, err := s.comments.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	log.Println(result)
	return nil
}

func (s *ForumStorage) DeleteAllDiscussions(ctx context.Context, userID int) error {
	var discussions []models.DiscussionTopic
	filter := bson.M{"author_id": userID}
	cursor, err := s.discussions.Find(context.TODO(), filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &discussions); err != nil {
		return err
	}
	for _, discussion := range discussions {
		s.DeleteFullDiscussion(ctx, discussion.ID)
	}
	
	return nil
}


