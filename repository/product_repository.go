package repository

import (
	"database/sql"
	"fmt"
	"krayshoping/model/entity"
)

type ProductRepository interface {
	AddProduct(product *entity.MerchantProduct) error
	UpdateProduct(product *entity.MerchantProduct) error
	DeleteProduct(productID int) error
	GetAllProducts() ([]*entity.MerchantProduct, error)
	GetByProductID(productID int) (*entity.MerchantProduct, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) AddProduct(product *entity.MerchantProduct) error {
	query := "INSERT INTO products (merchant_id, product_name, product_price, product_quantity) VALUES ($1, $2, $3, $4) RETURNING product_id"
	_, err := r.db.Exec(query, product.MerchantID, product.ProductName, product.ProductPrice, product.ProductQuantity)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepository) UpdateProduct(product *entity.MerchantProduct) error {
	query := "UPDATE products SET product_name = $1, product_price = $2, product_quantity = $3 WHERE product_id = $4"
	_, err := r.db.Exec(query, product.ProductName, product.ProductPrice, product.ProductQuantity, product.Product_ID)
	if err != nil {
		return err
	}
	return nil
}
func (r *productRepository) DeleteProduct(productID int) error {
	query := "DELETE FROM products WHERE product_id = $1"
	_, err := r.db.Exec(query, productID)
	if err != nil {
		return err
	}
	return nil
}
func (r *productRepository) GetByProductID(productID int) (*entity.MerchantProduct, error) {
	var product entity.MerchantProduct
	err := r.db.QueryRow("SELECT product_id,merchant_id,product_name,product_price, product_quantity,merchant_name FROM products WHERE product_id = $1", productID).Scan(&product.Product_ID, &product.MerchantID, &product.ProductName, &product.ProductPrice, &product.ProductQuantity, &product.MerchantName)
	if err != nil {

		return nil, fmt.Errorf("failed to get product: %v", err)
	}
	return &product, nil
}

func (r *productRepository) GetAllProducts() ([]*entity.MerchantProduct, error) {
	query := "SELECT product_id, merchant_id, product_name, product_price, product_quantity,merchant_name FROM products"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products := []*entity.MerchantProduct{}
	for rows.Next() {
		var product entity.MerchantProduct
		err := rows.Scan(&product.Product_ID, &product.MerchantID, &product.ProductName, &product.ProductPrice, &product.ProductQuantity, &product.MerchantName)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
