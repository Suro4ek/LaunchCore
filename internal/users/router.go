package users

import (
	"LaunchCore/eu.suro/launch/protos/user"
	"LaunchCore/internal/plugins"
	"LaunchCore/pkg/mysql"
	"context"
	"os"

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
	var plugins1 []*plugins.Plugin = make([]*plugins.Plugin, 0)
	// for _, plugin := range req.Plugins {
	// 	var p plugins.Plugin
	// 	err = r.client.DB.Model(&plugins.
	// 		Plugin{}).Where("id = ?", plugin).Find(&p).Error
	// 	if err != nil {
	// 		return nil, status.Errorf(codes.NotFound, "plugin not found")
	// 	}
	// 	plugins1 = append(plugins1, p)
	// }
	err = r.client.DB.Create(&User{
		Name:     req.Name,
		Plugins:  plugins1,
		RealName: req.RealName,
	}).Error
	if err != nil {
		return nil, err
	}
	return &user.Response{
		Status: "ok",
	}, nil
}

func (r *routerUser) GetUser(ctx context.Context, req *user.GetUserRequest) (res *user.GetUserResponse, err error) {
	var use User
	err = r.client.DB.Model(&User{}).Where("name = ?", req.Name).Find(&use).Error
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	//var friendsU []User
	//err = r.client.DB.Model(&use).Association("Friends").Find(&friendsU)
	//if err != nil {
	//	return nil, status.Errorf(codes.NotFound, "user not found")
	//}

	//var friends []*user.UserM = make([]*user.UserM, 0)
	//for _, friend := range friendsU {
	//	friends = append(friends, &user.UserM{
	//		Name: friend.Name,
	//	})
	//}
	var plugin []*user.Plugin = make([]*user.Plugin, 0)
	return &user.GetUserResponse{
		User: &user.UserM{
			Name:    use.Name,
			Plugins: plugin,
		},
	}, nil
}

//func (r *routerUser) AddFriend(ctx context.Context, req *user.AddFriendRequest) (res *user.Response, err error) {
//	var use User
//	err = r.client.DB.Model(&User{}).Where("name = ?", req.Name).Find(&use).Error
//	if err != nil {
//		return nil, status.Errorf(codes.NotFound, "user not found")
//	}
//	var friends []User
//	err = r.client.DB.Model(&use).Association("Friends").Find(&friends)
//	if err != nil {
//		return nil, status.Errorf(codes.NotFound, "user not found")
//	}
//	var friend User
//	err = r.client.DB.Model(&User{}).Where("name = ?", req.FriendName).Find(&friend).Error
//	if err != nil {
//		return nil, status.Errorf(codes.NotFound, "friend not found")
//	}
//	err = r.client.DB.Model(&use).Association("Friends").Append(&friend)
//	if err != nil {
//		return nil, err
//	}
//	return &user.Response{
//		Status: "ok",
//	}, nil
//}
//
//func (r *routerUser) RemoveFriend(ctx context.Context, req *user.RemoveFriendRequest) (res *user.Response, err error) {
//	var use User
//	err = r.client.DB.Model(&User{}).Where("name = ?", req.Name).Find(&use).Error
//	if err != nil {
//		return nil, status.Errorf(codes.NotFound, "user not found")
//	}
//	var friends []User
//	err = r.client.DB.Model(&use).Association("Friends").Find(&friends)
//	if err != nil {
//		return nil, status.Errorf(codes.NotFound, "user not found")
//	}
//	var friend User
//	err = r.client.DB.Model(&User{}).Where("name = ?", req.FriendName).Find(&friend).Error
//	if err != nil {
//		return nil, status.Errorf(codes.NotFound, "friend not found")
//	}
//	err = r.client.DB.Model(&use).Association("Friends").Delete(&friend)
//	if err != nil {
//		return nil, err
//	}
//	return &user.Response{
//		Status: "ok",
//	}, nil
//}

func (r *routerUser) DeleteWorld(ctx context.Context, req *user.RemoveWorldRequest) (res *user.Response, err error) {
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is empty")
	}
	err = os.Remove("/data/minecraft/" + req.Name)
	if err != nil {
		return nil, err
	}
	return &user.Response{
		Status: "ok",
	}, nil
}
