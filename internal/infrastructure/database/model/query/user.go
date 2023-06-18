package query_model

type Client struct {
	ID        int    `json:"id" db:"ـid"`
	UserID    int    `json:"user_id" db:"user_id"`
	PersonID  int    `json:"person_id" db:"person_id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Age       int    `json:"age"  db:"age"`
	Username  string `json:"username" db:"username"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"-" db:"password"`
}

func (c *Client) GetFullName() string {
	return c.FirstName + " " + c.LastName
}
