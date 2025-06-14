package repository

import (
	"database/sql"
	"fmt"
	model "go-api/models"
)

type ProductRepository struct {
	connection *sql.DB
}

func NewProductRepository(connection *sql.DB) ProductRepository {
	return ProductRepository{
		connection: connection,
	}
}

func (pr *ProductRepository) GetProducts() ([]model.Product, error) {
	query := "SELECT id, product_name, price FROM products"
	rows, err := pr.connection.Query(query)

	if err != nil {
		fmt.Println(err)
		return []model.Product{}, err
	}

	var productsList []model.Product
	var productsObj model.Product

	for rows.Next() {
		err = rows.Scan(
			&productsObj.ID,
			&productsObj.Name,
			&productsObj.Price,
		)

		if err != nil {
			fmt.Println(err)
			return []model.Product{}, err
		}

		productsList = append(productsList, productsObj)
	}

	rows.Close()

	return productsList, nil
}

func (pr *ProductRepository) CreateProduct(product model.Product) (int, error) {
	var id int

	query, err := pr.connection.Prepare(
		"INSERT INTO products" +
			"(product_name, price)" +
			" VALUES ($1, $2) RETURNING id")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	err = query.QueryRow(product.Name, product.Price).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	query.Close()

	return id, nil
}

func (pr *ProductRepository) GetProductById(id_product int) (*model.Product, error) {
	query, err := pr.connection.Prepare("SELECT * FROM products WHERE id = $1")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var product model.Product

	err = query.QueryRow(id_product).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	query.Close()

	return &product, nil
}

func (pr *ProductRepository) UpdateProduct(product model.Product) (*model.Product, error) {
	query, err := pr.connection.Prepare("UPDATE products SET product_name = $1, price = $2 WHERE id = $3 RETURNING id, product_name, price")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var updatedProduct model.Product

	err = query.QueryRow(product.Name, product.Price, product.ID).Scan(
		&updatedProduct.ID,
		&updatedProduct.Name,
		&updatedProduct.Price,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		fmt.Println(err)
		return nil, err
	}

	query.Close()

	return &updatedProduct, nil
}

func (pr *ProductRepository) DeleteProduct(id_product int) (bool, error) {
	query, err := pr.connection.Prepare("DELETE FROM products WHERE id = $1")

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	result, err := query.Exec(id_product)

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	query.Close()

	return rowsAffected > 0, nil
}
