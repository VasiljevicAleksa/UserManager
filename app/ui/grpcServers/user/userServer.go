package server

import (
	"context"

	"usermanager/app/domain"
	"usermanager/app/services"
	proto "usermanager/app/ui/protos/user"
	v "usermanager/app/ui/validations"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userServer struct {
	proto.UserServiceServer
	userService services.UserService
}

func NewUserGrpcServer(g *grpc.Server, u services.UserService) *userServer {
	userGrpcServer := userServer{
		userService: u,
	}
	proto.RegisterUserServiceServer(g, &userGrpcServer)
	return &userGrpcServer
}

func (s *userServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	// validate request
	if err := v.ValidateCreateUserReq(req); err != nil {
		log.Error().Err(err).Msg("validation failed for create user request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// add user
	id, err := s.userService.Add(req)
	if err != nil {
		log.Error().Err(err).Msg("create user failed")
		return nil, err
	}

	log.Info().Msgf("user %v sucessfully added", req.Nickname)
	return &proto.CreateUserResponse{Id: id.String()}, nil
}

func (s *userServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	// validate request
	if err := v.ValidateUpdateUserReq(req); err != nil {
		log.Error().Err(err).Msg("validation failed for update user request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// update user
	if err := s.userService.Update(req); err != nil {
		log.Error().Err(err).Msgf("update user with id %v failed", req.Id)
		return nil, err
	}

	log.Info().Msgf("user with id %v successfully updated", req.Id)
	return &proto.UpdateUserResponse{Id: req.Id}, nil
}

func (s *userServer) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	// validate request
	if err := v.ValidateDeleteUserReq(req); err != nil {
		log.Error().Err(err).Msg("validation failed for delete user request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// delete user
	if err := s.userService.Delete(req.Id); err != nil {
		log.Error().Err(err).Msgf("delete user with id %v failed", req.Id)
		return nil, err
	}

	log.Info().Msgf("user with id %v successfully deleted", req.Id)
	return &proto.DeleteUserResponse{Id: req.Id}, nil
}

func (s *userServer) GetUserPage(ctx context.Context, req *proto.UserPageRequest) (*proto.UserPageResponse, error) {
	// validate request
	if err := v.ValidateUserPageReq(req); err != nil {
		log.Error().Err(err).Msg("validation failed for user page request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// get user page
	users, err := s.userService.GetPage(req)
	if err != nil {
		log.Error().Err(err).Msg("get user page failed")
		return nil, err
	}

	return userPageResponse(users), nil
}

func userPageResponse(users []domain.User) *proto.UserPageResponse {
	response := proto.UserPageResponse{
		Users: make([]*proto.UserPageResponse_User, 0, len(users)),
	}

	for _, u := range users {
		response.Users = append(response.Users, &proto.UserPageResponse_User{
			Id:        u.Id.String(),
			Firstname: u.Firstname,
			Lastname:  u.Lastname,
			Nickname:  u.Nickname,
			Email:     u.Email,
			Country:   u.Country,
			Created:   timestamppb.New(u.CreatedAt),
		})
	}

	return &response
}
