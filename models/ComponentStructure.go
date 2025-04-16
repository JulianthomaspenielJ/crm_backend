package models

import "database/sql"

type Component struct {
	ID            int    `json:"id"`
	ComponentName string `json:"componentName"`
	Description   string `json:"description"`
	ProductType   string `json:"productType"`
	CreatedBy     string `json:"createdBy"`
	CreatedAt     string `json:"createdAt"`
}

type VariantData struct {
	VariantName        string  `json:"variantName"`
	VariantType        string  `json:"variantType"`
	VariantPrice       float64 `json:"variantPrice"`
	VariantDescription string  `json:"variantDescription"`
	WarrantyInfo       string  `json:"warrantyInfo"`
	Brand              string  `json:"brand"`
	VariantImg         string  `json:"variantImg"`
	VariantID          string  `json:"variantId"`
	CreatedBy          string  `json:"createdBy"`
}

type VariantItems struct {
	ID                 int     `json:"id"`
	VariantName        string  `json:"variantName"`
	VariantType        string  `json:"variantType"`
	VariantPrice       float64 `json:"variantPrice"`
	VariantDescription string  `json:"variantDescription"`
	VarientBrand       string  `json:"brand"`
	CreatedBy          string  `json:"createdBy"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedBy          sql.NullString
	UpdatedAt          sql.NullString
}

type ProductData struct {
	ProcessorId []struct {
		Id       int `json:"id"`
		Quantity int `json:"quantity"`
	} `json:"processorId"`
	MotherboardId []struct {
		Id       int `json:"id"`
		Quantity int `json:"quantity"`
	} `json:"motherboardId"`
	RamId []struct {
		Id       int `json:"id"`
		Quantity int `json:"quantity"`
	} `json:"ramId"`
	StorageId []struct {
		Id       int `json:"id"`
		Quantity int `json:"quantity"`
	} `json:"storageId"`
	SmpsId []struct {
		Id       int `json:"id"`
		Quantity int `json:"quantity"`
	} `json:"smpsId"`
	MonitorId []struct {
		Id       int `json:"id"`
		Quantity int `json:"quantity"`
	} `json:"monitorId"`
	KeyboardMouseId []struct {
		Id       int `json:"id"`
		Quantity int `json:"quantity"`
	} `json:"keyboardMouseId"`
	OsId []struct {
		Id       int `json:"id"`
		Quantity int `json:"quantity"`
	} `json:"osId"`
	ChassisId []struct {
		Id       int `json:"id"`
		Quantity int `json:"quantity"`
	} `json:"chassisId"`
	LanId []struct {
		Id       int `json:"id"`
		Quantity int `json:"quantity"`
	} `json:"lanId"`
	WifiId []struct {
		Id       int `json:"id"`
		Quantity int `json:"quantity"`
	} `json:"wifiId"`
	DvdDriverId []struct {
		Id       int `json:"id"`
		Quantity int `json:"quantity"`
	} `json:"dvdDriverId"`
	GraphicsId []struct {
		Id       int `json:"id"`
		Quantity int `json:"quantity"`
	} `json:"graphicsId"`
	CoolerId []struct {
		Id       int `json:"id"`
		Quantity int `json:"quantity"`
	} `json:"coolerId"`
	WarrantyInfo string `json:"warrantyInfo"`
	Packages     string `json:"packages"`
	Logistics    string `json:"logistics"`
}

type Products struct {
	Id       int `json:"id"`
	Quantity int `json:"quantity"`
}

type ProductDataInsert struct {
	ProcessorId      int    `json:"processorId"`
	ProcessorQty     int    `json:"processorQty"`
	MotherboardId    int    `json:"motherboardId"`
	MotherboardQty   int    `json:"motherboardQty"`
	RamId            int    `json:"ramId"`
	RamQty           int    `json:"ramQty"`
	StorageId1       int    `json:"storageId1"`
	StorageQty1      int    `json:"storageQty1"`
	StorageId2       int    `json:"storageId2"`
	StorageQty2      int    `json:"storageQty2"`
	StorageId3       int    `json:"storageId3"`
	StorageQty3      int    `json:"storageQty3"`
	SmpsId           int    `json:"smpsId"`
	SmpsQty          int    `json:"smpsQty"`
	MonitorId        int    `json:"monitorId"`
	MonitorQty       int    `json:"monitorQty"`
	KeyboardMouseId  int    `json:"keyboardMouseId"`
	KeyboardMouseQty int    `json:"keyboardMouseQty"`
	OsId             int    `json:"osId"`
	OsQty            int    `json:"osQty"`
	ChassisId        int    `json:"chassisId"`
	ChassisQty       int    `json:"chassisQty"`
	LanId            int    `json:"lanId"`
	LanQty           int    `json:"lanQty"`
	WifiId           int    `json:"wifiId"`
	WifiQty          int    `json:"wifiQty"`
	DvdDriverId      int    `json:"dvdDriverId"`
	DvdDriverQty     int    `json:"dvdDriverQty"`
	GraphicsId       int    `json:"graphicsId"`
	GraphicsQty      int    `json:"graphicsQty"`
	CoolerId         int    `json:"coolerId"`
	CoolerQty        int    `json:"coolerQty"`
	WarrantyInfo     string `json:"warrantyInfo"`
	Packages         string `json:"packages"`
	Logistics        string `json:"logistics"`
}

// Define the structure that matches the JSON format in the Products column
type DealProductData struct {
	ID                   int           `json:"id"`
	Data                 []Product     `json:"data"`
	DealID               int           `json:"dealID"`
	IsDeal               sql.NullInt64 `json:"isDeal"`
	CreatedAt            string        `json:"CreatedAt"`
	CreatedBy            string        `json:"CreatedBy"`
	AddressTo            string        `json:"addressTo"`
	MarginType           string        `json:"marginType"`
	MarginTypePercentage int           `json:"marginTypePercentage"`
	TotalAmount          float64       `json:"totalAmount"`
	TierFinalValue       string        `json:"tierFinalValue"` // Change this to string
	TotalDiscountPrice   float64       `json:"totalDiscountPrice"`
}

// Product holds individual product details in the `data` array
type DealProduct struct {
	ID           int     `json:"id"`
	Qty          int     `json:"qty"`
	Name         string  `json:"Name"`
	Price        float64 `json:"price"`
	Descriptions string  `json:"descriptions"`
}

type ProductComponent struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Count       int     `json:"count"`
}
type Product struct {
	ProductID     int               `json:"product_id"`
	Processor     ProductComponent  `json:"processor"`
	MotherBoard   ProductComponent  `json:"motherboard"`
	RAM           ProductComponent  `json:"ram"`
	SMPS          ProductComponent  `json:"smps"`
	Storage1      ProductComponent  `json:"storage1"`
	Storage2      *ProductComponent `json:"storage2,omitempty"`
	Storage3      *ProductComponent `json:"storage3,omitempty"`
	Monitor       ProductComponent  `json:"monitor"`
	OS            ProductComponent  `json:"os"`
	Chassis       ProductComponent  `json:"chassis"`
	LAN           *ProductComponent `json:"lan,omitempty"`
	WiFi          *ProductComponent `json:"wifi,omitempty"`
	Keyboard      *ProductComponent `json:"keyboard,omitempty"`
	DVD           *ProductComponent `json:"dvd,omitempty"`
	Graphics      *ProductComponent `json:"graphics,omitempty"`
	Cooler        *ProductComponent `json:"cooler,omitempty"`
	CreatedBy     string            `json:"created_by"`
	CreatedAt     string            `json:"created_at"`
	Warranty      string            `json:"warranty"`
	Logistics     string            `json:"logistics"`
	Package       string            `json:"package"`
	ReferenceName string            `json:"referenceName"`
}

type ProductsList struct {
	ProductID     int    `json:"product_id"`
	ReferenceName string `json:"package"`
	CreatedBy     string `json:"created_by"`
	CreatedAt     string `json:"created_at"`
}

type UpdateComponent struct {
	ID                 int    `json:"id"`
	Brand              string `json:"brand"`
	VarientName        string `json:"variantName"`
	VarientDescription string `json:"variantDescription"`
	VariantPrice       int64  `json:"variantPrice"`
	UpdatedBy          string `json:"createdBy"`
}

type GetProductsDetails struct {
	ID        int    `json:"id"`
	DeletedBy string `json:"deletedBy"`
	Reason    string `json:"reason"`
}

type GetProductList struct {
	ID          int    `json:"id"`
	GemNumber   string `json:"gemNumber"`
	DueDate     string `json:"dueDate"`
	Quantity    string `json:"quantity"`
	DealerName  string `json:"dealerName"`
	Customer    string `json:"customer"`
	GemRaisedAt string `json:"gemRaisedAt"`
	Level       string `json:"level"`
	Reason      string `json:"reason"`
	GemEndDate  string `json:"GemEndDate"`
	ProductName string `json:"productName"`
	Description string `json:"description"`
	ProductID   int    `json:"productId"`
	CreatedBy   string `json:"createdBy"`
	CreatedAt   string `json:"createdAt"`
	Status      string `json:"status"`
}

type CustomError struct {
	Code    int
	Message string
}
type PricingProduct struct {
	Data                   []PricingComponent `json:"data"`
	ID                     int                `json:"id"`
	WarrantyInfo           int                `json:"warrantyInfo"`
	WarrantyAmount         int                `json:"warrantyAmount"`
	ProductsQuantity       int                `json:"productsQuantity"`
	TotalPrice             float64            `json:"totalAmount"`
	TotalAmountWithGST     float64            `json:"totalAmountWithGST"`
	MarginAmount           float64            `json:"marginAmount"`
	ADP                    int                `json:"adpAmount"`
	Bag                    int                `json:"bagAmount"`
	LogisticPackage        int                `json:"logisticsPackages"`
	LogisticValue          string             `json:"logisticsValue"`
	TotalAmountWithMargin  float64            `json:"totalAmountWithmargin"`
	ReferenceName          string             `json:"referenceName"`
	PackagesValue          string             `json:"packagesValue"`
	MarginName             string             `json:"marginName"`
	CreatedBy              string             `json:"CreatedBy"`
	CreatedAt              string             `json:"CreatedAt"`
	AddressTo              string             `json:"addressTo"`
	GstAmount              string             `json:"gstAmount"`
	GstPercentage          string             `json:"gstPercentage"`
	ReferenceID            int                `json:"referenceID"`
	UpdatedBy              string             `json:"updatedBy"`
	MarginType             string             `json:"marginType"`
	MarginPercentage       int                `json:"marginTypePercentage"`
	MarginWithGST          string             `json:"marginwithGST"`
	TierFinalValue         string             `json:"tierFinalValue"`
	HddReturns             int                `json:"hddReturns"`
	LaptopWarranty         int                `json:"laptopWarranty"`
	LaptopWarrantyPrice    int                `json:"laptopWarrantyAmount"`
	CustomerName           string             `json:"customerName"`
	DiscountAmount         string             `json:"discountAmount"`
	TotalPriceWithDiscount int                `json:"totalDiscountPrice"`
	IsDeal                 sql.NullInt64      `json:"isDeal"`
	DealID                 int                `json:"dealID"`
	TicketNo               string             `json:"ticketNo"`
}

type PricingComponent struct {
	ID           int    `json:"id"`
	Name         string `json:"Name"`
	Descriptions string `json:"descriptions"`
	Price        int    `json:"price"`
	Quantity     int    `json:"qty"`
}

type ProductMainList struct {
	ID        int    `json:"id"`
	Product   string `json:"product"`
	CreatedBy string `json:"createdBy"`
	CreatedAt string `json:"createdAt"`
}

type ProductID struct {
	ID int `json:"productID"`
}

type PricingMainComponent struct {
	ProductName string `json:"productName"`
	Description string `json:"description"`
}

type PricingMainListProduct struct {
	CreatedBy          string                 `json:"createdBy"`
	ProductName        string                 `json:"productName"`
	SelectedComponents []PricingMainComponent `json:"selectedComponents"`
}

type ProductReferData struct {
	ID           int    `json:"id"`
	Qty          int    `json:"qty"`
	Name         string `json:"Name"`
	Price        int    `json:"price"`
	Descriptions string `json:"descriptions"`
}

type ProductReferenceStruct struct {
	ReferenceName         string             `json:"referenceName"`
	Products              []ProductReferData `json:"data"`
	AddressTo             string             `json:"addressTo"`
	AdpAmount             int                `json:"adpAmount"`
	BagAmount             int                `json:"bagAmount"`
	GstAmount             string             `json:"gstAmount"`
	UpdatedBy             string             `json:"updatedBy"`
	HddReturns            int                `json:"hddReturns"`
	MarginName            string             `json:"marginName"`
	MarginType            string             `json:"marginType"`
	ReferenceID           int                `json:"referenceID"`
	TotalAmount           int                `json:"totalAmount"`
	CustomerName          string             `json:"customerName"`
	MarginAmount          int                `json:"marginAmount"`
	WarrantyInfo          int                `json:"warrantyInfo"`
	GstPercentage         string             `json:"gstPercentage"`
	MarginWithGST         string             `json:"marginWithGST"`
	PackagesValue         string             `json:"packagesValue"`
	DiscountAmount        string             `json:"discountAmount"`
	LaptopWarranty        int                `json:"laptopWarranty"`
	LogisticsValue        string             `json:"logisticsValue"`
	TierFinalValue        string             `json:"tierFinalValue"`
	WarrantyAmount        int                `json:"warrantyAmount"`
	ProductsQuantity      int                `json:"productsQuantity"`
	LogisticsPackages     int                `json:"logisticsPackages"`
	TotalAmountWithGST    int                `json:"totalAmountWithGST"`
	TotalDiscountPrice    int                `json:"totalDiscountPrice"`
	LaptopWarrantyAmount  int                `json:"laptopWarrantyAmount"`
	MarginTypePercentage  int                `json:"marginTypePercentage"`
	TotalAmountWithMargin int                `json:"totalAmountWithMargin"`
}

type ProductDetailedData struct {
	ID            int    `json:"id"`
	Products      string `json:"products"`
	ReferenceName string `json:"referenceName"`
}

type ComponentDetail struct {
	ID                 int    `json:"id"`
	VarientName        string `json:"varientName"`
	VarientDescription string `json:"varientDescription"`
	VarientPrice       int    `json:"varientPrice"`
	VarientType        string `json:"varientType"`
}
