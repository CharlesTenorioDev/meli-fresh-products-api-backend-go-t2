package application

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/repository"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/routes"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/service"
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

	localityRepo := repository.NewMysqlLocalityRepository(a.db)
	provinceRepo := repository.NewMysqlProvinceRepository(a.db)
	countryRepo := repository.NewMysqlCountryRepository(a.db)
	localityService := service.NewBasicLocalityService(localityRepo, provinceRepo, countryRepo)
	err = routes.NewLocalityRoutes(router, localityService)

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
	sellerRepo := repository.NewSellerMysql(a.db)
	sellerService := service.NewSellerService(sellerRepo, localityRepo)

	if err := routes.RegisterSellerRoutes(router, sellerService); err != nil {
		panic(err)
	}

	// Requisito 4 - ProductType
	productTypeRepo := repository.NewProductTypeDB(a.db)

	productTypeService := service.NewProductTypeService(productTypeRepo)
	if err := routes.NewProductTypeRoutes(router, productTypeService); err != nil {
		panic(err)
	}

	// Requisito 4 - Product
	productRepo := repository.NewProductDB(a.db)
	productService := service.NewProductService(productRepo, productTypeService)
	err = routes.NewProductRoutes(router, productService)

	if err != nil {
		panic(err)
	}

	//Requisito 4 - Product Records
	productRecordsRepo := repository.NewProductRecordDB(a.db)
	productRecordsService := service.NewProductRecordService(productRecordsRepo, productService)

	err = routes.NewProductRecordsRoutes(router, productRecordsService)
	if err != nil {
		panic(err)
	}

	// Requisito 2 - Warehouses
	warehouseRepo := repository.NewWarehouseRepository(a.db)
	warehouseService := service.NewWarehouseService(warehouseRepo, localityRepo)

	err = routes.NewWarehouseRoutes(router, warehouseService)
	if err != nil {
		panic(err)
	}

	// Requisito 3 - Section

	sectionRepo := repository.NewSectionMysql(a.db)
	sectionService := service.NewBasicSectionService(sectionRepo, warehouseService, productTypeService)

	err = routes.RegisterSectionRoutes(router, sectionService)
	if err != nil {
		panic(err)
	}

	// Requisito 5 - Employees
	employeesRepo := repository.NewEmployeeRepository(a.db)

	employeesService := service.NewEmployeeService(employeesRepo, warehouseService)
	if err := routes.RegisterEmployeesRoutes(router, employeesService); err != nil {
		panic(err)
	}

	// Requisito 6 - Buyers
	buyersRepo := repository.NewBuyerDb(a.db)
	buyersService := service.NewBuyer(buyersRepo)
	// Create the routes and deps
	if err = routes.BuyerRoutes(router, buyersService); err != nil {
		panic(err)
	}

	// Requisito 6 - Purchase Orders
	purchaseOrdersRepo := repository.NewPurchaseOrderDb(a.db)
	purchaseOrdersService := service.NewPurchaseOrderService(purchaseOrdersRepo, buyersService, productRecordsRepo)

	err = routes.RegisterPurchaseOrdersRoutes(router, purchaseOrdersService)
	if err != nil {
		panic(err)
	}

	// Sprint2 Requisito 2 - Carry
	carriesRepo := repository.NewMySQLCarryRepository(a.db)

	carryService := service.NewMySQLCarryService(carriesRepo, localityRepo)
	if err = routes.CarryRoutes(router, carryService); err != nil {
		panic(err)
	}

	// Sprint2 Requisito 3 - Product Batch
	productBatchRepo := repository.NewProductBatchRepository(a.db)
	productBatchService := service.NewProductBatchesService(productBatchRepo, productRepo, sectionRepo)

	if err = routes.ProductBatchRoutes(router, productBatchService); err != nil {
		panic(err)
	}

	inboundOrderRepo := repository.NewInboundOrderRepository(a.db)
	inboundOrderService := service.NewInboundOrderService(inboundOrderRepo)

	if err := routes.RegisterInboundOrderRoutes(router, inboundOrderService); err != nil {
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
