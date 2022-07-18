package dto

type NotificationsUpdateDTO struct {
	MessageNotifications bool
	FollowNotifications  bool
	LikeNotifications    bool
	CommentNotifications bool
}
