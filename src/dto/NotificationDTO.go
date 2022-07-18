package dto

type NotificationType int

const (
	Message NotificationType = iota
	Follow
	Like
	Comment
)

type NotificationDTO struct {
	Message          string
	UserAuth0ID      string
	NotificationType *NotificationType
}
