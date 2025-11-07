package entity

type PersonalAccessToken struct {
	ID        int
	UserID    int
	Token     string
	CreatedAt int64
	ExpiredAt int64
}
