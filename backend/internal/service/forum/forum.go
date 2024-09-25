package forum

import (
	"context"
	"errors"
	"fmt"
	"gohelp/internal/models"
	"gohelp/internal/storage/mongo"
	"log"
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

func (s *ForumService) CreateComment(ctx context.Context, related_to, discussionID, content string, authorID int) (string, error) {
	var err error
	if _, err = s.repo.GetDiscussion(ctx, discussionID); err != nil {
		return "", fmt.Errorf("error during searching related disc: %v", err)
	}
	if related_to != "" {
		comment, err := s.repo.GetComment(ctx, related_to)
		if err != nil {
			return "", fmt.Errorf("error during searching related comment: %v", err)
		}
		if comment.DiscussionID != discussionID {
			return "", fmt.Errorf("invalid related comment: %v", err)
		}
	}

	comment := &models.Comment{
		DiscussionID: discussionID,
		RelatedTo:    related_to,
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
	commentsTree := buildCommentTree(comments)

	return discussion, commentsTree, nil
}

func buildCommentTree(comments []models.Comment) []models.Comment {

	tree := make(map[string][]models.Comment)
	var rootComments []models.Comment

	for _, comment := range comments {
		if comment.Deleted {
			comment.Content = "<deleted>"
		}
		parentID := ""
		if comment.RelatedTo != "" {
			parentID = comment.RelatedTo
		}
		tree[parentID] = append(tree[parentID], comment)
	}

	var addChildren func(string) []models.Comment
	addChildren = func(parentID string) []models.Comment {
		if children, ok := tree[parentID]; ok {
			for i := range children {
				children[i].Children = addChildren(children[i].ID)
			}
			return children
		}
		return nil
	}

	rootComments = addChildren("")
	return rootComments
}
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

func (s *ForumService) Vote(ctx context.Context, userID int, element_id, voteType string) error {
	_, err1 := s.repo.GetDiscussion(ctx, element_id)
	_, err2 := s.repo.GetComment(ctx, element_id)
	if (err1 != nil && err2 != nil) || (err1 == nil && err2 == nil) {
		return errors.New("nothing was found or discussion with comment has equal ids")
	} else if err1 == nil {
		err := s.VoteDiscussion(ctx, userID, element_id, voteType)
		if err != nil {
			return err
		}
		return nil
	} else if err2 == nil {
		err := s.VoteComment(ctx, userID, element_id, voteType)
		if err != nil {
			return err
		}
		return nil
	}
	log.Println("function was ended suspicious")
	return nil
}

func (s *ForumService) VoteDiscussion(ctx context.Context, userID int, discussionID, voteType string) error {

	err := s.repo.RemoveVote(userID, discussionID, voteType)
	if err != nil {
		return fmt.Errorf("error during removing votes: %v", err)
	}
	err = s.repo.DiscAddVote(userID, discussionID, voteType)
	if err != nil {
		return fmt.Errorf("error during adding vote: %v", err)
	}
	return nil
}

func (s *ForumService) VoteComment(ctx context.Context, userID int, commentID, voteType string) error {

	err := s.repo.RemoveVote(userID, commentID, voteType)
	if err != nil {
		return fmt.Errorf("error during removing votes: %v", err)
	}
	err = s.repo.ComAddVote(userID, commentID, voteType)
	if err != nil {
		return fmt.Errorf("error during adding vote: %v", err)
	}
	return nil
}

func (s *ForumService) UpdateDiscussion(ctx context.Context, discussionID, content string, authorID int) (*models.Discussion, error) {

	disc, err := s.repo.GetDiscussion(ctx, discussionID)
	if err != nil {
		return nil, fmt.Errorf("error during getting discussion: %v", err)
	}
	if disc.AuthorID != authorID {
		return nil, errors.New("you have no permissions to do this")
	}
	err = s.repo.UpdateDiscussion(ctx, discussionID, content)
	if err != nil {
		return nil, fmt.Errorf("error during updating discussion: %v", err)
	}
	disc, err = s.repo.GetDiscussion(ctx, discussionID)
	if err != nil {
		return nil, fmt.Errorf("error during getting discussion: %v", err)
	}

	return disc, nil
}
func (s *ForumService) UpdateComment(ctx context.Context, commentID, content string, authorID int) (*models.Comment, error) {

	comm, err := s.repo.GetComment(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("error during getting discussion: %v", err)
	}
	if comm.AuthorID != authorID {
		return nil, errors.New("you have no permissions to do this")
	}
	err = s.repo.UpdateDiscussion(ctx, commentID, content)
	if err != nil {
		return nil, fmt.Errorf("error during updating discussion: %v", err)
	}
	comm, err = s.repo.GetComment(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("error during getting discussion: %v", err)
	}
	return comm, nil
}

func (s *ForumService) DeleteFullDiscussion(ctx context.Context, discussionID string) error {
	_, err := s.repo.GetDiscussion(ctx, discussionID)
	if err != nil {
		return fmt.Errorf("error during getting discussion: %v", err)
	}
	err = s.repo.DeleteFullDiscussion(ctx, discussionID)
	if err != nil {
		return fmt.Errorf("error during updating discussion: %v", err)
	}

	return nil
}
func (s *ForumService) DeleteComment(ctx context.Context, commentID, userRole string, authorID int) error {

	comm, err := s.repo.GetComment(ctx, commentID)
	if err != nil {
		return fmt.Errorf("error during getting discussion: %v", err)
	}
	if userRole == models.CustomerRole {
		if comm.AuthorID != authorID {
			return errors.New("you have no permissions to do this")
		}
	}
	err = s.repo.DeleteComment(ctx, commentID)
	if err != nil {
		return fmt.Errorf("error during updating discussion: %v", err)
	}

	return nil
}

func (s *ForumService) DeleteFullHistory(ctx context.Context, userID int) error{
	err := s.repo.DeleteAllComments(ctx, userID)
	if err != nil {
		return  fmt.Errorf("error during deleting comments: %v", err)
	}

	err = s.repo.DeleteAllDiscussions(ctx, userID)
	if err != nil {
		return  fmt.Errorf("error during deleting discussions: %v", err)
	}
	return nil

}