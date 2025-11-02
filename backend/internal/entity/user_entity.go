package entity

type UserEntity struct {
	ID        int
	Name      string
	Email     string
	Password  string
	CreatedAt int64
	UpdatedAt int64
}

type PersonalAccessToken struct {
	ID        int
	UserID    int
	Token     string
	CreatedAt int64
	ExpiredAt int64
}
