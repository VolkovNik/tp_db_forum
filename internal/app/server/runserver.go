package server

import (
	"TPForum/internal/pkg/config"
	_forumDelivery 		"TPForum/internal/pkg/forum/delivery"
	_forumRepository 	"TPForum/internal/pkg/forum/repository"
	_forumUsecase 		"TPForum/internal/pkg/forum/usecase"
	_postDelivery 		"TPForum/internal/pkg/post/delivery"
	_postRepository 	"TPForum/internal/pkg/post/repository"
	_postUsecase 		"TPForum/internal/pkg/post/usecase"
	_serviceDelivery 	"TPForum/internal/pkg/service/delivery"
	_serviceRepository 	"TPForum/internal/pkg/service/repository"
	_serviceUsecase 	"TPForum/internal/pkg/service/usecase"
	_threadDelivery 	"TPForum/internal/pkg/thread/delivery"
	_threadRepository 	"TPForum/internal/pkg/thread/repository"
	_threadUsecase 		"TPForum/internal/pkg/thread/usecase"
	_userDelivery 		"TPForum/internal/pkg/user/delivery"
	_userRepository 	"TPForum/internal/pkg/user/repository"
	_userUsecase 		"TPForum/internal/pkg/user/usecase"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"net/http"
)

func RunServer(addr string) {
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()

	configPostgres := pgx.ConnConfig{
		User:                 config.Get().Postgres.User,
		Database:             config.Get().Postgres.DBName,
		Password:             config.Get().Postgres.Password,
		PreferSimpleProtocol: false,
	}

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     configPostgres,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}

	connPool, err := pgx.NewConnPool(poolConfig)
	if err != nil {
		fmt.Println(err)
	}

	forumRepo 	:= _forumRepository.NewForumRepository(connPool)
	userRepo 	:= _userRepository.NewUserRepository(connPool)
	threadRepo 	:= _threadRepository.NewThreadRepository(connPool)
	postRepo 	:= _postRepository.NewPostRepository(connPool)
	serviceRepo := _serviceRepository.NewServiceRepository(connPool)

	forumUC	 	:= _forumUsecase.NewForumUsecase(forumRepo)
	userUC 		:= _userUsecase.NewUserUsecase(userRepo)
	threadUC 	:= _threadUsecase.NewThreadUsecase(threadRepo)
	postUC 		:= _postUsecase.NewPostUsecase(postRepo)
	serviceUC	:= _serviceUsecase.NewServiceUsecase(serviceRepo)

	_forumDelivery.NewForumDelivery(s, forumUC)
	_userDelivery.NewUserDelivery(s, &userUC)
	_threadDelivery.NewThreadDelivery(s, threadUC)
	_postDelivery.NewPostDelivery(s, postUC)
	_serviceDelivery.NewServiceDelivery(s, serviceUC)

	defer connPool.Close()
	if err = http.ListenAndServe(":5000", s); err != nil {
		fmt.Println(err)
	}
}
