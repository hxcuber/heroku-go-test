package request

type PostFriendConnectionRequest struct {
	friends []string `json:"friends"`
}
