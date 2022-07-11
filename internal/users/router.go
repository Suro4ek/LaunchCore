package users

import (
	"LaunchCore/eu.suro/launch/protos/user"
	"LaunchCore/pkg/mysql"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type routerUser struct {
	client *mysql.Client
	user.UnimplementedUserServer
}

func NewRouterUser(client *mysql.Client) user.UserServer {
	return &routerUser{
		client: client,
	}
}

func (r *routerUser) CreateUser(ctx context.Context, req *user.CreateUserRequest) (res *user.Response, err error) {
	err = r.client.DB.Create(&User{
		Name:    req.Name,
		Plugins: req.Plugins,
	}).Error
	if err != nil {
		return nil, err
	}
	return &user.Response{
		Status: "ok",
	}, nil
}

func (r *routerUser) GetUser(ctx context.Context, req *user.GetUserRequest) (res *user.GetUserResponse, err error) {
	err = r.client.DB.Model(&User{}).Where("name = ?", req.Name).Association("Friends").Find(&res.User)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return res, nil
}

func (r *routerUser) AddFriend(ctx context.Context, req *user.AddFriendRequest) (res *user.Response, err error) {
	var use User
	err = r.client.DB.Model(&User{}).Where("name = ?", req.Name).Association("Friends").Find(&use)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	var friend User
	err = r.client.DB.Model(&User{}).Where("name = ?", req.FriendName).Find(&friend).Error
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "friend not found")
	}
	err = r.client.DB.Model(&use).Association("Friends").Append(friend)
	if err != nil {
		return nil, err
	}
	return &user.Response{
		Status: "ok",
	}, nil
}

func (r *routerUser) RemoveFriend(ctx context.Context, req *user.RemoveFriendRequest) (res *user.Response, err error) {
	var use User
	err = r.client.DB.Model(&User{}).Where("name = ?", req.Name).Association("Friends").Find(&use)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	var friend User
	err = r.client.DB.Model(&User{}).Where("name = ?", req.FriendName).Find(&friend).Error
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "friend not found")
	}
	err = r.client.DB.Model(&use).Association("Friends").Delete(friend)
	if err != nil {
		return nil, err
	}
	return &user.Response{
		Status: "ok",
	}, nil
}
