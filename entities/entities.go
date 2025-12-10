package entities

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"name"`
	LastName  string `json:"lastname"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Role      string `json:"role"`
	RegDate   string `json:"reg_date"`
}

type Orders struct {
	ID          string `json:"id" gorm:"column:id;primaryKey"`
	Items       string `json:"items"`
	UserID      string `json:"user_id"`
	TotalAmount string `json:"total_amount"`
	Status      string `json:"status"`
	Address     string `json:"address"`
}

type Review struct {
	ID      string `json:"id"`
	Rating  string `json:"rating"`
	Comment string `json:"comment"`
	BookID  string `json:"book_id"`
	UserID  string `json:"user_id"`
}

type Discount struct {
	ID          string  `json:"id" gorm:"column:id;primaryKey"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Percentage  float64 `json:"percentage"`
}

type Publisher struct {
	ID    string `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Book struct {
	ID                string    `json:"id" gorm:"column:id;primaryKey"`
	ISBN              string    `json:"isbn"`
	PublisherID       string    `json:"publisher_id"`
	Title             string    `json:"title"`
	AuthorID          string    `json:"author_id"`
	Description       string    `json:"description"`
	Price             string    `json:"price"`
	Language          string    `json:"language"`
	CategoryID        string    `json:"category_id"`
	Category          string    `json:"category"`
	Author            string    `json:"author"`
	Link              string    `json:"link"`
	Preorder          bool      `json:"preorder"`
	AvailablePreorder string    `json:"available"`
	Rating            float64   `json:"rating"`
	Publisher         Publisher `json:"publisher" gorm:"foreignKey:PublisherID"`
}
