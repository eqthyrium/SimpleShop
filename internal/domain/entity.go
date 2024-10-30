package domain

type User struct {
	UserId         int    `json:"user_id"`
	Nickname       string `json:"nickname"`
	MemberIdentity string `json:"member_identity"`
	Password       string `json:"password"`
	Role           string `json:"role"`
}

type Product struct {
	ProductId   int    `json:"product_id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Cost        int    `json:"cost"`
}

// Notification, Reaction struct must be created
