package service

import (
	"context"
	"database/sql"
	"fmt"
	pbc "like-service/genproto/comment_service"
	pb "like-service/genproto/like_service"
	pbp "like-service/genproto/post_service"
	pbu "like-service/genproto/user_service"
	"like-service/pkg/logger"
	grpcclient "like-service/service/grpc_client"
	"like-service/storage"
	"log"
)

type LikeService struct {
	storage storage.IStorage
	logger  logger.Logger
	client  grpcclient.IServiceManager
}

func NewLikeService(db *sql.DB, log logger.Logger, client grpcclient.IServiceManager) *LikeService {
	return &LikeService{
		storage:                        storage.NewStoragePg(db),
		logger:                         log,
		client:                         client,
	}
}

//rpc LikePost(PostLike) returns (Status);
//rpc LikeComment(CommentLike) returns (Status);
//rpc GetLikeOwnersByPostId(GetPostId) returns (Post);
//rpc GetLikeOwnersByCommentId(GetCommentId) returns (Comment);

func (l *LikeService) LikePost(ctx context.Context, req *pb.PostLike) (*pb.Status, error) {
	return l.storage.Like().LikePost(req)
}

func (l *LikeService) LikeComment(ctx context.Context, req *pb.CommentLike) (*pb.Status, error) {
	return l.storage.Like().LikeComment(req)
}

func (l *LikeService) GetLikeOwnersByPostId(ctx context.Context, req *pb.GetPostId) (*pb.Post, error) {
	post, err := l.storage.Like().GetLikeOwnersByPostId(req)
	if err != nil {
		fmt.Println("error service1")
		return nil, err
	}
	respPost, err := l.client.PostService().GetPostById(ctx, &pbp.GetPostId{
		PostId: req.PostId,
	})
	if err != nil {
		log.Fatal("cannot get post by id post service", err.Error())
	}

	post.Id = respPost.Id
	post.Title = respPost.Title
	post.ImageUrl = respPost.ImageUrl
	post.OwnerId = respPost.OwnerId

	for _, owner := range post.Likes {
		user, err := l.client.UserService().GetUserById(ctx, &pbu.GetUserId{
			UserId: owner.Id,
		})
		if err != nil {
			log.Fatal("error while getting user by id", err.Error())
			return nil, err
		}
		owner.Id = user.Id
		owner.FirstName = user.FirstName
		owner.LastName = user.LastName
		owner.Age = user.Age
		owner.Gender = pb.Gender(user.Gender)
	}

	return post, nil
}

func (l *LikeService) GetLikeOwnersByCommentId(ctx context.Context, req *pb.GetCommentId) (*pb.Comment, error) {
	comment, err := l.storage.Like().GetLikeOwnersByCommentId(req)
	if err != nil {
		log.Fatal("error while getting like owner by comment id", err.Error())
	}

	respComment, err := l.client.CommentService().GetCommentById(ctx, &pbc.GetCommentId{
		Id: req.CommentId,
	})
	if err != nil {
		log.Fatal("error while getting comment by id", err.Error())
	}
	comment.Id = respComment.Id
	comment.Content = respComment.Content
	comment.OwnerId = respComment.OwnerId
	comment.PostId = respComment.PostId
	comment.CreatedAt = respComment.CreatedAt

	for _, owner := range comment.Likes {
		user, err := l.client.UserService().GetUserById(ctx, &pbu.GetUserId{
			UserId: owner.Id,
		})
		if err != nil {
			log.Fatal("error while getting owner by id", err.Error())
			return nil, err
		}
		owner.Id = user.Id
		owner.FirstName = user.FirstName
		owner.LastName = user.LastName
		owner.Age = user.Age
		owner.Gender = pb.Gender(user.Gender)
	}
	return comment, nil
}
