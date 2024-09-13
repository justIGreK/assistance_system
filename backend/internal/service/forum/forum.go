package forum

import (
	"context"
	"fmt"
	"gohelp/internal/models"
	"gohelp/internal/storage/mongo"
)

type ForumService struct {
	repo *mongo.ForumStorage
}

func NewForumService(repo *mongo.ForumStorage) *ForumService {
	return &ForumService{repo: repo}
}

func (s *ForumService) CreateDiscussion(ctx context.Context, title, content string, authorID int) (string, error) {
	discussion := &models.Discussion{
		Title:    title,
		Content:  content,
		AuthorID: authorID,
	}
	return s.repo.CreateDiscussion(ctx, discussion)
}

func (s *ForumService) CreateComment(ctx context.Context, discussionID, content string, authorID int) (string, error) {
	_, err := s.repo.GetDiscussion(ctx, discussionID)
	if err != nil {
		return "", fmt.Errorf("error during creating comment: %v", err)
	}
	
	comment := &models.Comment{
		DiscussionID: discussionID,
		Content:      content,
		AuthorID:     authorID,
	}
	return s.repo.CreateComment(ctx, comment)
}

func (s *ForumService) GetDiscussionWithComments(ctx context.Context, discussionID string) (*models.Discussion, []models.Comment, error) {
	discussion, err := s.repo.GetDiscussion(ctx, discussionID)
	if err != nil {
		return nil, nil, fmt.Errorf("error during getting discussion: %v", err)
	}

	comments, err := s.repo.GetCommentsByDiscussion(ctx, discussionID)
	if err != nil {
		return nil, nil, fmt.Errorf("error during getting comments of discussions: %v", err)
	}

	return discussion, comments, nil
}

// func (s *ForumService) GetAllDiscussions() ([]models.DiscussionWithCount, error) {

// 	var result []DiscussionWithCount
// 	for _, discussion := range discussions {
// 		count, err := commentsCollection.CountDocuments(context.TODO(), bson.M{
// 			"discussionId": discussion.ID,
// 		})
// 		if err != nil {
// 			return nil, err
// 		}

// 		result = append(result, DiscussionWithCount{
// 			Discussion:    discussion,
// 			CommentsCount: count,
// 		})
// 	}

//		return result, nil
//	}
func (s *ForumService) GetAllDiscussionsWithCountOfComments(ctx context.Context) ([]models.DiscussionWithCount, error) {
	discussions, err := s.repo.GetAllDiscussions(ctx)
	if err != nil {
		return nil, fmt.Errorf("error during getting list of discussions: %v", err)
	}
	summary, err := s.repo.GetSummaryOfDiscussions(ctx, discussions)
	if err != nil {
		return nil, fmt.Errorf("error during getting list of comments for discussions: %v", err)
	}
	return summary, nil
}

func (s *ForumService) SearchDiscussionsByName(ctx context.Context, searchTerm string) ([]models.Discussion, error) {
	
	discussions, err := s.repo.SearchDiscussionsByName(ctx, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("error during getting list of discussions ny name: %v", err)
	}
	return discussions, nil
}