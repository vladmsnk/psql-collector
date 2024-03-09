package loader

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"
)

const testColumn = "data"

type Implementation struct {
	db *sql.DB
}

func New(db *sql.DB) *Implementation {
	return &Implementation{db: db}
}

type Loader interface {
	FillWithTestData(ctx context.Context)
	ProvokeTableBloat(ctx context.Context, tableName string) error
}

func (i *Implementation) FillWithTestData(ctx context.Context) error {
	const (
		usersNum      = 10000
		categoriesNum = 100
		productsNum   = 5000
		ordersNum     = 20000
		orderItemsNum = 50000
		reviewsNum    = 50000
	)

	isFilled, err := i.checkDataIsFilled(ctx)
	if err != nil {
		return fmt.Errorf("i.checkDataFilled: %w", err)
	}
	if isFilled {
		return nil
	}

	i.insertUsers(usersNum)
	i.insertCategories(categoriesNum)
	i.insertProducts(productsNum)
	time.Sleep(time.Minute)
	i.insertOrders(ordersNum)
	i.insertOrderItems(orderItemsNum)
	i.insertReviews(reviewsNum)

	log.Println("database filled")
	return nil
}

func (i *Implementation) ProvokeTableBloat(ctx context.Context, tableName string) error {

	//var (
	//	createColumnQuery = fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s VARCHAR(255)", tableName, testColumn)
	//	updateQuery       = fmt.Sprintf("UPDATE %s SET %s = %s || $1 WHERE id %% 2 = 0", tableName, testColumn, testColumn)
	//	deleteQuery       = fmt.Sprintf("DELETE FROM %s WHERE id %% 10 = $1", tableName)
	//)
	//
	//i.db.ExecContext(ctx, createColumnQuery)

	return nil
}

func (impl *Implementation) checkDataIsFilled(ctx context.Context) (bool, error) {
	query := `
select count(*) from orders
`
	var count int
	err := impl.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("impl.db.QueryRowContext: %w", err)
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (impl *Implementation) insertUsers(count int) {
	for i := 0; i < count; i++ {
		_, err := impl.db.Exec("INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)",
			fmt.Sprintf("user%d", i), fmt.Sprintf("user%d@example.com", i), "hashedpassword")
		if err != nil {
			log.Fatal("Error inserting into users: ", err)
		}
	}
}

func (impl *Implementation) insertCategories(count int) {
	for i := 0; i < count; i++ {
		_, err := impl.db.Exec("INSERT INTO categories (name, description) VALUES ($1, $2)",
			fmt.Sprintf("category%d", i), "A generic category description")
		if err != nil {
			log.Fatal("Error inserting into categories: ", err)
		}
	}
}

func (impl *Implementation) insertProducts(count int) {
	for i := 0; i < count; i++ {
		_, err := impl.db.Exec("INSERT INTO products (name, description, price, stock_quantity, category_id) VALUES ($1, $2, $3, $4, $5)",
			fmt.Sprintf("product%d", i), "A generic product description", rand.Float64()*100, rand.Intn(100), rand.Intn(100)+1)
		if err != nil {
			log.Fatal("Error inserting into products: ", err)
		}
	}
}

func (impl *Implementation) insertOrders(count int) {
	for i := 0; i < count; i++ {
		_, err := impl.db.Exec("INSERT INTO orders (user_id, order_date, status, total) VALUES ($1, $2, $3, $4)",
			rand.Intn(10000)+1, time.Now(), "Completed", rand.Float64()*100)
		if err != nil {
			log.Fatal("Error inserting into orders: ", err)
		}
	}
}

func (impl *Implementation) insertOrderItems(count int) {
	for i := 0; i < count; i++ {
		_, err := impl.db.Exec("INSERT INTO order_items (order_id, product_id, quantity, price_at_purchase) VALUES ($1, $2, $3, $4)",
			rand.Intn(20000)+1, rand.Intn(5000)+1, rand.Intn(5)+1, rand.Float64()*100)
		if err != nil {
			log.Fatal("Error inserting into order_items: ", err)
		}
	}
}

func (impl *Implementation) insertReviews(count int) {
	for i := 0; i < count; i++ {
		_, err := impl.db.Exec("INSERT INTO reviews (product_id, user_id, rating, comment) VALUES ($1, $2, $3, $4)",
			rand.Intn(5000)+1, rand.Intn(10000)+1, rand.Intn(5)+1, "A generic review comment")
		if err != nil {
			log.Fatal("Error inserting into reviews: ", err)
		}
	}
}
