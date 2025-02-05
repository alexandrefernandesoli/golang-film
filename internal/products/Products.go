package products

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ProductsResponse represents the top-level JSON structure
type ProductsResponse struct {
	Products Products `json:"products"`
}

// Products represents the pagination and product data
type Products struct {
	CurrentPage  int       `json:"current_page"`
	Data         []Product `json:"data"`
	FirstPageURL string    `json:"first_page_url"`
	From         int       `json:"from"`
	LastPage     int       `json:"last_page"`
	LastPageURL  string    `json:"last_page_url"`
	NextPageURL  string    `json:"next_page_url"`
	Path         string    `json:"path"`
	PerPage      string    `json:"-"`
	PrevPageURL  *string   `json:"prev_page_url"`
	To           int       `json:"to"`
	Total        int       `json:"total"`
}

// Product represents a single product with all its details
type Product struct {
	ID                       int       `json:"id"`
	ShopID                   int       `json:"shop_id"`
	Title                    string    `json:"title"`
	Handle                   string    `json:"handle"`
	TemplateSuffix           string    `json:"template_suffix"`
	Description              string    `json:"description"`
	SEOTitle                 string    `json:"seo_title"`
	SEODescription           *string   `json:"seo_description"`
	Price                    float64   `json:"price"`
	CompareAtPrice           float64   `json:"compare_at_price"`
	CostPerItem              float64   `json:"cost_per_item"`
	SKU                      *string   `json:"sku"`
	Taxable                  int       `json:"taxable"`
	Barcode                  string    `json:"barcode"`
	Weight                   float64   `json:"weight"`
	WeightUnit               *string   `json:"weight_unit"`
	RequiresShipping         int       `json:"requires_shipping"`
	InventoryPolicy          int       `json:"inventory_policy"`
	Quantity                 int       `json:"quantity"`
	ProductTypeID            int       `json:"product_type_id"`
	VendorID                 int       `json:"vendor_id"`
	AddedVia                 string    `json:"added_via"`
	PublishedAt              string    `json:"published_at"`
	DeletedAt                *string   `json:"deleted_at"`
	CreatedAt                string    `json:"created_at"`
	UpdatedAt                string    `json:"updated_at"`
	Active                   int       `json:"active"`
	IsAvailableCheckout      int       `json:"is_available_checkout"`
	ActiveForBots            int       `json:"active_for_bots"`
	SimplyBookID             int       `json:"simplybook_id"`
	MinimumQuantity          int       `json:"minimum_quantity"`
	MaxQuantity              int       `json:"max_quantity"`
	MinimumQuantityCount     int       `json:"minimum_quantity_count"`
	MaxQuantityCount         int       `json:"max_quantity_count"`
	TotalUnitsSold           int       `json:"total_units_sold"`
	Kit                      int       `json:"kit"`
	PageLink                 *string   `json:"page_link"`
	EnableMultipleCurrencies int       `json:"enable_multiple_currencies"`
	SalesPage                *string   `json:"sales_page"`
	VSL                      *string   `json:"vsl"`
	ProductAccessDetails     *string   `json:"product_access_details"`
	Supplier                 *string   `json:"supplier"`
	Purpose                  *string   `json:"purpose"`
	AliProductID             *string   `json:"ali_product_id"`
	GoogleMerchantStatus     string    `json:"google_merchant_status"`
	ProductDefaultVariant    Variant   `json:"product_default_variant"`
	ProductVariants          []Variant `json:"product_variants"`
	Images                   []Image   `json:"images"`
}

// Variant represents a product variant
type Variant struct {
	ID                       int            `json:"id"`
	ProductID                int            `json:"product_id"`
	Default                  int            `json:"default"`
	Title                    string         `json:"title"`
	Price                    string         `json:"price"`
	CompareAtPrice           string         `json:"compare_at_price"`
	CostPerItem              string         `json:"cost_per_item"`
	SKU                      string         `json:"sku"`
	AliExpressSKU            *string        `json:"aliexpress_sku"`
	Position                 int            `json:"position"`
	InventoryPolicy          int            `json:"inventory_policy"`
	Quantity                 int            `json:"quantity"`
	PreventOutOfStockSelling int            `json:"prevent_out_of_stock_selling"`
	Taxable                  int            `json:"taxable"`
	Barcode                  string         `json:"barcode"`
	Swatches                 *string        `json:"swatches"`
	Length                   string         `json:"length"`
	Width                    string         `json:"width"`
	Height                   string         `json:"height"`
	DimensionUnit            *string        `json:"dimension_unit"`
	Weight                   float64        `json:"weight"`
	WeightUnit               *string        `json:"weight_unit"`
	RequiresShipping         int            `json:"requires_shipping"`
	HasDigitalAttachment     int            `json:"has_digital_attachment"`
	EnabledShopifyRedirect   int            `json:"enabled_shopify_redirect"`
	CreatedAt                string         `json:"created_at"`
	UpdatedAt                string         `json:"updated_at"`
	DeletedAt                *string        `json:"deleted_at"`
	NimbleSKU                *string        `json:"nimble_sku"`
	ImageID                  int64          `json:"-"`
	MobilePrice              float64        `json:"mobile_price"`
	VariantImage             []VariantImage `json:"variant_image"`
}

// VariantImage represents an image associated with a variant
type VariantImage struct {
	ProductImageID   int   `json:"product_image_id"`
	ProductVariantID int   `json:"product_variant_id"`
	Image            Image `json:"image"`
}

// Image represents a product image
type Image struct {
	ID          int    `json:"id"`
	ProductID   int    `json:"product_id"`
	Position    int    `json:"position"`
	Alt         string `json:"alt"`
	Src         string `json:"src"`
	MediaType   string `json:"media_type"`
	AspectRatio string `json:"aspect_ratio"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	URL         string `json:"url"`
}

func GetProducts(token string) ([]Product, error) {
	// Create a new request
	req, err := http.NewRequest("GET", "https://accounts.cartpanda.com/api/eightcomercio/products", nil)
	if err != nil {
		return nil, fmt.Errorf("error reading request. %v", err)
	}

	// Add authorization header
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + token},
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request. %v", err)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response. %v", err)
	}

	// Unmarshal the response body
	var data ProductsResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response. %v", err)
	}

	return data.Products.Data, nil
}
