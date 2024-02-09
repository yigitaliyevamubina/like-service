package grpc_client

import (
	"fmt"
	"google.golang.org/grpc"
	"like-service/config"
	pbc "like-service/genproto/comment_service"
	pbp "like-service/genproto/post_service"
	pbu "like-service/genproto/user_service"
	"like-service/pkg/logger"
	"log"
)

type IServiceManager interface {
	UserService() pbu.UserServiceClient
	PostService() pbp.PostServiceClient
	CommentService() pbc.CommentServiceClient
}

type serviceManager struct {
	cfg            config.Config
	userService    pbu.UserServiceClient
	postService    pbp.PostServiceClient
	commentService pbc.CommentServiceClient
}

func New(cfg config.Config) (IServiceManager, error) {
	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithInsecure())
	if err != nil {
		log.Fatal("error while dialing to the user service", logger.Error(err))
	}

	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.PostServiceHost, cfg.PostServicePort),
		grpc.WithInsecure())
	if err != nil {
		log.Fatal("error while dialing to the post service", logger.Error(err))
	}

	connComment, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.CommentServiceHost, cfg.CommentServicePort),
		grpc.WithInsecure())
	if err != nil {
		log.Fatal("error while dialing to the comment service", logger.Error(err))
	}

	return &serviceManager{
		cfg:            cfg,
		userService:    pbu.NewUserServiceClient(connUser),
		postService:    pbp.NewPostServiceClient(connPost),
		commentService: pbc.NewCommentServiceClient(connComment),
	}, nil
}

func (s *serviceManager) UserService() pbu.UserServiceClient {
	return s.userService
}

func (s *serviceManager) PostService() pbp.PostServiceClient {
	return s.postService
}

func (s *serviceManager) CommentService() pbc.CommentServiceClient {
	return s.commentService
}
