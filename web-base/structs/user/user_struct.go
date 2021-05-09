package user

// User web user struct
type User struct {
	Name             string `bson:"name" json:"name"`
	Username         string `bson:"username" json:"username"`
	Email            string `bson:"email" json:"email"`
	Password         string `bson:"password" json:"password"`
	PrivateID        string `bson:"private_id" json:"private_id"`
	PublicID         string `bson:"public_id" json:"public_id"`
	IsBanned         bool   `bson:"is_banned" json:"is_banned"`
	TOTPEnabled      bool   `bson:"totp_enabled" json:"totp_enabled"`
	SingleSession    bool   `bson:"single_session" json:"single_session"`
	CreatedAt        int64  `bson:"created_at" json:"created_at"`
	LastLoginRequest int64  `bson:"last_login_request" json:"last_login_request"`
	LastLoginSucceed int64  `bson:"last_login_succeed" json:"last_login_succeed"`
	LastLoggedIPv4   string `bson:"last_logged_ipv4" json:"last_logged_ipv4"`
}
