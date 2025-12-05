package entities

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"name"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Role      string `json:"role"`
	RegDate   string `json:"reg_date"`
}

type Publisher struct {
	PublisherID string `json:"publisher_id" gorm:"column:publisher_id;primaryKey"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
}

type Orders struct {
	ID          string `json:"id" gorm:"column:id;primaryKey"`
	Description string `json:"description"`
	DateStart   string `json:"date_start"`
	DateEnd     string `json:"date_end"`
	Percent     string `json:"percent"`
}

type Review struct {
	ID      string `json:"id"`
	Rating  string `json:"rating"`
	Comment string `json:"comment"`
}

type Discount struct {
	DiscountID  string  `json:"discount_id" gorm:"column:discount_id;primaryKey"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Percentage  float64 `json:"percentage"`
}
type Book struct {
	ID          string    `json:"id" gorm:"column:id;primaryKey"`
	ISBN        string    `json:"isbn"`
	PublisherID string    `json:"publisher_id"`
	Title       string    `json:"title"`
	AuthorID    string    `json:"author_id"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	Language    string    `json:"language"`
	CategoryID  string    `json:"category_id"`
	Category    string    `gorm:"category"`
	Author      string    `gorm:"author"`
	Link        string    `gorm:"link"`
	Publisher   Publisher `gorm:"foreignKey:PublisherID"`
}
