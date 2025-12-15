package entities

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"password"`
	Email     string    `json:"email" gorm:"unique"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Role      string    `json:"role"`
	RegDate   time.Time `json:"reg_date" gorm:"default:CURRENT_TIMESTAMP"`
}

type Orders struct {
	ID          int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Items       string  `json:"items"`
	UserID      int     `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
	Address     string  `json:"address"`
	User        User    `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Review struct {
	ID      int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Rating  string `json:"rating"`
	Comment string `json:"comment"`
	BookID  int    `json:"book_id"`
	UserID  int    `json:"user_id"`
	Book    Book   `json:"book" gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE"`
	User    User   `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Discount struct {
	ID          int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Percentage  float64 `json:"percentage" gorm:"check:percentage >= 0 AND percentage <= 100"`
}

type Publisher struct {
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
	Phone string `json:"phone"`
}

type Category struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
}

type Book struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ISBN          string    `json:"isbn" gorm:"unique"`
	PublisherID   int       `json:"publisher_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	Language      string    `json:"language"`
	CategoryID    int       `json:"category_id"`
	Author        string    `json:"author"`
	Link          string    `json:"link"`
	Preorder      bool      `json:"preorder" gorm:"default:false"`
	AvailableDate time.Time `json:"available_date"`
	Rating        float64   `json:"rating" gorm:"check:rating >= 0 AND rating <= 5"`
	DiscountID    int       `json:"discount_id"`
	Publisher     Publisher `json:"publisher" gorm:"foreignKey:PublisherID;constraint:OnDelete:CASCADE"`
	Category      Category  `json:"category" gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
	Discount      Discount  `json:"discount" gorm:"foreignKey:DiscountID;constraint:OnDelete:CASCADE"`
}
