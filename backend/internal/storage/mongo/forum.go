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
}

func NewForumStorage(db *mongo.Database) *ForumStorage {
	return &ForumStorage{
		discussions: db.Collection("discussions"),
		comments:    db.Collection("comments"),
	}
}

func (s *ForumStorage) CreateDiscussion(ctx context.Context, discussion *models.Discussion) (string, error) {
	discussion.CreatedAt = time.Now()

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
	err = s.discussions.FindOne(ctx, bson.M{"_id": oid}).Decode(&discussion)
	if err != nil {
		return nil, err
	}
	return &discussion, nil
}

func (s *ForumStorage) GetAllDiscussions(ctx context.Context) ([]models.DiscussionTopic, error) {
	var discussions []models.DiscussionTopic
	cursor, err := s.discussions.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &discussions); err != nil {
		return nil, err
	}
	return discussions, nil
}

func (s *ForumStorage) CreateComment(ctx context.Context, comment *models.Comment) (string, error) {
	comment.CreatedAt = time.Now()

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
			"discussion_id": discussion.ID,
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
	cursor, err := s.comments.Find(context.TODO(), bson.M{
		"discussion_id": discussionID,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &comments); err != nil {
		return nil, err
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
