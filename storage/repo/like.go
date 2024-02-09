package repo

import (
	pb "like-service/genproto/like_service"
)

//rpc LikePost(PostLike) returns (Status);
//rpc LikeComment(CommentLike) returns (Status);
//rpc GetLikeOwnersByPostId(GetPostId) returns (Post);
//rpc GetLikeOwnersByCommentId(GetCommentId) returns (Comment);

//LikeServiceI

type LikeServiceI interface {
	LikePost(*pb.PostLike) (*pb.Status, error)
	LikeComment(*pb.CommentLike) (*pb.Status, error)
	GetLikeOwnersByPostId(*pb.GetPostId) (*pb.Post, error)
	GetLikeOwnersByCommentId(*pb.GetCommentId) (*pb.Comment, error)
}
