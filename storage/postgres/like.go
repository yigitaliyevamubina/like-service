package postgres

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	pb "like-service/genproto/like_service"
)

type likeRepo struct {
	db *sql.DB
}

// NewLikeRepo

func NewLikeRepo(db *sql.DB) *likeRepo {
	return &likeRepo{db: db}
}

//rpc LikePost(PostLike) returns (Status);
//rpc LikeComment(CommentLike) returns (Status);
//rpc GetLikeOwnersByPostId(GetPostId) returns (Post);
//rpc GetLikeOwnersByCommentId(GetCommentId) returns (Comment);

func (l *likeRepo) LikePost(reqLike *pb.PostLike) (*pb.Status, error) {
	status := pb.Status{
		Liked: false,
	}
	query := `INSERT INTO postlikes(post_id, user_id) VALUES($1, $2)`
	_, err := l.db.Exec(query, reqLike.PostId, reqLike.OwnerId)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code.Name() == "unique_violation" {
			query := `DELETE FROM postlikes WHERE post_id = $1 AND user_id = $2`
			_, err := l.db.Exec(query, reqLike.PostId, reqLike.OwnerId)
			if err != nil {
				fmt.Println("error")
				return nil, err
			}
			status.Liked = false
		}
	} else {
		status.Liked = true
	}
	return &status, nil
}

func (l *likeRepo) LikeComment(reqLike *pb.CommentLike) (*pb.Status, error) {
	status := pb.Status{
		Liked: false,
	}
	query := `INSERT INTO commentlikes(comment_id, user_id) VALUES($1, $2)`
	_, err := l.db.Exec(query, reqLike.CommentId, reqLike.OwnerId)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code.Name() == "unique_violation" {
			query := `DELETE FROM commentlikes WHERE comment_id = $1 AND user_id = $2`
			_, err := l.db.Exec(query, reqLike.CommentId, reqLike.OwnerId)
			if err != nil {
				fmt.Println("error")
				return nil, err
			}
			status.Liked = false
		}
	} else {
		status.Liked = true
	}

	return &status, nil
}

func (l *likeRepo) GetLikeOwnersByPostId(postId *pb.GetPostId) (*pb.Post, error) {
	post := pb.Post{
		Id:       "",
		Title:    "",
		ImageUrl: "",
		OwnerId:  "",
		Likes:    []*pb.Owner{},
	}

	query := `SELECT user_id FROM postlikes WHERE post_id = $1`

	rows, err := l.db.Query(query, postId.PostId)
	if err != nil {
		fmt.Println("error2")
		return nil, err
	}

	for rows.Next() {
		owner := pb.Owner{
			Id:        "",
			FirstName: "",
			LastName:  "",
			Age:       0,
			Gender:    0,
		}

		if err := rows.Scan(&owner.Id); err != nil {
			fmt.Println("error3")
			return nil, err
		}

		post.Likes = append(post.Likes, &owner)
	}

	return &post, err
}

func (l *likeRepo) GetLikeOwnersByCommentId(commentId *pb.GetCommentId) (*pb.Comment, error) {
	comment := pb.Comment{
		Id:        "",
		Content:   "",
		OwnerId:   "",
		PostId:    "",
		CreatedAt: "",
		UpdatedAt: "",
		DeletedAt: "",
		Likes:     []*pb.Owner{},
	}

	query := `SELECT user_id FROM commentlikes WHERE comment_id = $1`

	rows, err := l.db.Query(query, commentId.CommentId)
	if err != nil {
		fmt.Println("error4")
		return nil, err
	}

	for rows.Next() {
		owner := pb.Owner{
			Id:        "",
			FirstName: "",
			LastName:  "",
			Age:       0,
			Gender:    0,
		}

		if err := rows.Scan(&owner.Id); err != nil {
			fmt.Println("error5")
			return nil, err
		}

		comment.Likes = append(comment.Likes, &owner)
	}

	return &comment, err
}
