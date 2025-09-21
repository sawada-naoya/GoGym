package user

type IDProvider interface{ NewUserID() string }
