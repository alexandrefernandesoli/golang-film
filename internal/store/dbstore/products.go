package dbstore

import (
	"fmt"
	"goth/internal/products"
	"time"
)

type ProductsStore struct {
	Products       []products.Product
	LastTimeLoaded time.Time
}

type NewProductStoreParams struct {
	Products []products.Product
}

func NewProductStore(params NewProductStoreParams) *ProductsStore {
	return &ProductsStore{
		Products: params.Products,
	}
}

func (s *ProductsStore) LoadProducts() error {
	products, err := products.GetProducts("YOMlENvfDbt7RUTvXe0hvgKMY7YSb58YwbDvMJjCx5bi7C4RyQbl4gjNqMrk")

	if err != nil {
		fmt.Println("Error loading products")
		return err
	}

	s.Products = products
	s.LastTimeLoaded = time.Now()

	return nil
}

func (s *ProductsStore) GetProducts() []products.Product {
	if time.Since(s.LastTimeLoaded) > 5*time.Second {
		go func() {
			err := s.LoadProducts()
			if err != nil {
				fmt.Println("Error loading products")
			}
		}()
	}

	return s.Products
}

func (s *ProductsStore) GetLastTimeLoaded() string {
	// as DD/MM/YYYY HH:MM:SS
	return s.LastTimeLoaded.Format("02/01/2006 15:04:05")
}
