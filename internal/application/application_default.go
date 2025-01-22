package application

import (
	"database/sql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/buyer"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/carry"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/country"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/employee"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/inbound_order"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/locality"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/product"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/product_batch"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/product_record"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/product_type"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/province"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/purchase_order"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/section"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/seller"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/warehouse"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-sql-driver/mysql"
	_ "github.com/meli-fresh-products-api-backend-go-t2/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// ConfigApplicationDefault is the configuration for NewApplicationDefault.
type ConfigApplicationDefault struct {
	// DB is the database configuration.
	DB *mysql.Config
	// Addr is the server address.
	Addr string
}

// NewApplicationDefault creates a new ApplicationDefault.
func NewApplicationDefault(config *ConfigApplicationDefault) *ApplicationDefault {
	// default values
	defaultCfg := &ConfigApplicationDefault{
		DB:   nil,
		Addr: ":8080",
	}

	if config != nil {
		if config.DB != nil {
			defaultCfg.DB = config.DB
		}

		if config.Addr != "" {
			defaultCfg.Addr = config.Addr
		}
	}

	return &ApplicationDefault{
		cfgDB:   defaultCfg.DB,
		cfgAddr: defaultCfg.Addr,
	}
}

// ApplicationDefault is an implementation of the Application interface.
type ApplicationDefault struct {
	// cfgDB is the database configuration.
	cfgDB *mysql.Config
	// cfgAddr is the server address.
	cfgAddr string
	// db is the database connection.
	db *sql.DB
	// router is the chi router.
	router *chi.Mux
}

// TearDown tears down the application.
func (a *ApplicationDefault) TearDown() {
	// close db
	if a.db != nil {
		a.db.Close()
	}
}

// SetUp initializes the application by setting up the database connection,
// configuring the router, and registering various routes and services.
func (a *ApplicationDefault) SetUp() (err error) {
	// connect to db
	a.db, err = sql.Open("mysql", a.cfgDB.FormatDSN())

	if err != nil {
		log.Fatalf("error opening db: %s", err.Error())
	}

	if err = a.db.Ping(); err != nil {
		log.Fatalf("error pinging db: %s", err.Error())
	}

	router := chi.NewRouter()
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition"
	))

	localityRepo := locality.NewMysqlLocalityRepository(a.db)
	provinceRepo := province.NewMysqlProvinceRepository(a.db)
	countryRepo := country.NewMysqlCountryRepository(a.db)
	localityService := locality.NewBasicLocalityService(localityRepo, provinceRepo, countryRepo)
	err = locality.NewLocalityRoutes(router, localityService)

	if err != nil {
		panic(err)
	}

	// Requisito 1 - Seller
	// ldSellers := internal.NewSellerJSONFile("./internal/sellers.json")
	// dbSellers, err := ldSellers.Load()
	// if err != nil {
	// return
	// }
	// sellerRepo := repository.NewSellerDBRepository(dbSellers)
	sellerRepo := seller.NewSellerRepository(a.db)
	sellerService := seller.NewSellerService(sellerRepo, localityRepo)

	if err := seller.RegisterSellerRoutes(router, sellerService); err != nil {
		panic(err)
	}

	// Requisito 4 - ProductType
	productTypeRepo := product_type.NewProductTypeDB(a.db)

	productTypeService := product_type.NewProductTypeService(productTypeRepo)
	if err := product_type.NewProductTypeRoutes(router, productTypeService); err != nil {
		panic(err)
	}

	// Requisito 4 - Product
	productRepo := product.NewProductDB(a.db)
	productService := product.NewProductService(productRepo, productTypeService)
	err = product.NewProductRoutes(router, productService)

	if err != nil {
		panic(err)
	}

	//Requisito 4 - Product Records
	productRecordsRepo := product_record.NewProductRecordDB(a.db)
	productRecordsService := product_record.NewProductRecordService(productRecordsRepo, productService)

	err = product_record.NewProductRecordsRoutes(router, productRecordsService)
	if err != nil {
		panic(err)
	}

	// Requisito 2 - Warehouses
	warehouseRepo := warehouse.NewWarehouseRepository(a.db)
	warehouseService := warehouse.NewWarehouseService(warehouseRepo, localityRepo)

	err = warehouse.NewWarehouseRoutes(router, warehouseService)
	if err != nil {
		panic(err)
	}

	// Requisito 3 - Section

	sectionRepo := section.NewSectionMysql(a.db)
	sectionService := section.NewBasicSectionService(sectionRepo, warehouseService, productTypeService)

	err = section.RegisterSectionRoutes(router, sectionService)
	if err != nil {
		panic(err)
	}

	// Requisito 5 - Employees
	employeesRepo := employee.NewEmployeeRepository(a.db)

	employeesService := employee.NewEmployeeService(employeesRepo, warehouseService)
	if err := employee.RegisterEmployeesRoutes(router, employeesService); err != nil {
		panic(err)
	}

	// Requisito 6 - Buyers
	buyersRepo := buyer.NewBuyerDB(a.db)
	buyersService := buyer.NewBuyer(buyersRepo)
	// Create the routes and deps
	if err = buyer.BuyerRoutes(router, buyersService); err != nil {
		panic(err)
	}

	// Requisito 6 - Purchase Orders
	purchaseOrdersRepo := purchase_order.NewPurchaseOrderDB(a.db)
	purchaseOrdersService := purchase_order.NewPurchaseOrderService(purchaseOrdersRepo, buyersService, productRecordsRepo)

	err = purchase_order.RegisterPurchaseOrdersRoutes(router, purchaseOrdersService)
	if err != nil {
		panic(err)
	}

	// Sprint2 Requisito 2 - Carry
	carriesRepo := carry.NewMySQLCarryRepository(a.db)

	//carryService := carry.NewMySQLCarryService(carriesRepo, localityRepo)
	carryService := carry.NewMySQLCarryService(carriesRepo, localityRepo)
	if err = carry.CarryRoutes(router, carryService); err != nil {
		panic(err)
	}

	// Sprint2 Requisito 3 - Product Batch
	productBatchRepo := product_batch.NewProductBatchRepository(a.db)
	productBatchService := product_batch.NewProductBatchesService(productBatchRepo, productRepo, sectionRepo)

	if err = product_batch.ProductBatchRoutes(router, productBatchService); err != nil {
		panic(err)
	}

	inboundOrderRepo := inbound_order.NewInboundOrderRepository(a.db)
	inboundOrderService := inbound_order.NewInboundOrderService(inboundOrderRepo)

	if err := inbound_order.RegisterInboundOrderRoutes(router, inboundOrderService); err != nil {
		panic(err)
	}

	a.router = router

	return nil
}

// Run runs the application.
func (a *ApplicationDefault) Run() (err error) {
	defer a.db.Close()
	log.Printf("starting server at %s\n", a.cfgAddr)

	err = http.ListenAndServe(a.cfgAddr, a.router)

	return
}
