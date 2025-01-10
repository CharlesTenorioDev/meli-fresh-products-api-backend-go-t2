package application

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/loader"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/repository"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/routes"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/service"
)

// ConfigApplicationDefault is the configuration for NewApplicationDefault.
type ConfigApplicationDefault struct {
	// Db is the database configuration.
	Db *mysql.Config
	// Addr is the server address.
	Addr string
}

// NewApplicationDefault creates a new ApplicationDefault.
func NewApplicationDefault(config *ConfigApplicationDefault) *ApplicationDefault {
	// default values
	defaultCfg := &ConfigApplicationDefault{
		Db:   nil,
		Addr: ":8080",
	}
	if config != nil {
		if config.Db != nil {
			defaultCfg.Db = config.Db
		}
		if config.Addr != "" {
			defaultCfg.Addr = config.Addr
		}
	}

	return &ApplicationDefault{
		cfgDb:   defaultCfg.Db,
		cfgAddr: defaultCfg.Addr,
	}
}

// ApplicationDefault is an implementation of the Application interface.
type ApplicationDefault struct {
	// cfgDb is the database configuration.
	cfgDb *mysql.Config
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

func (a *ApplicationDefault) SetUp() (err error) {
	// connect to db
	a.db, err = sql.Open("mysql", a.cfgDb.FormatDSN())
	if err != nil {
		log.Fatalf("error opening db: %s", err.Error())
	}
	if err = a.db.Ping(); err != nil {
		log.Fatalf("error pinging db: %s", err.Error())
	}

	router := chi.NewRouter()
	a.router = router

	// Requisito 1 - Seller
	ldSellers := internal.NewSellerJSONFile("./internal/sellers.json")
	dbSellers, err := ldSellers.Load()
	if err != nil {
		return
	}
	sellerRepo := repository.NewSellerDbRepository(dbSellers)
	sellerService := service.NewSellerService(sellerRepo)
	if err := routes.RegisterSellerRoutes(router, sellerService); err != nil {
		panic(err)
	}

	// Requisito 4 - ProductType
	productTypeRepo := repository.NewProductTypeDB(nil)
	productTypeService := service.NewProductTypeService(productTypeRepo)
	if err := routes.NewProductTypeRoutes(router, productTypeService); err != nil {
		panic(err)
	}

	// Requisito 4 - Product
	productRepo := repository.NewProductDB(nil)
	productService := service.NewProductService(productRepo, productTypeService)
	err = routes.NewProductRoutes(router, productService)
	if err != nil {
		panic(err)
	}

	// Requisito 2 - Warehouses
	warehouseRepo := repository.NewWarehouseDB(a.db)
	warehouseService := service.NewWarehouseService(warehouseRepo)
	err = routes.NewWarehouseRoutes(router, warehouseService)
	if err != nil {
		panic(err)
	}

	// Requisito 3 - Section
	sectionRepo := repository.NewMemorySectionRepository(nil)
	sectionService := service.NewBasicSectionService(sectionRepo, warehouseService, productTypeService)
	err = routes.RegisterSectionRoutes(router, sectionService)
	if err != nil {
		panic(err)
	}

	// Requisito 5 - Employees
	filePath := "docs/db/employees.json"
	ld := loader.NewEmployeeJsonFile(filePath)
	db, err := ld.Load()
	if err != nil {
		fmt.Println(err)
		return
	}
	employeesRepo := repository.NewEmployeeRepository(db)
	employeesService := service.NewEmployeeService(employeesRepo, warehouseService)
	if err := routes.RegisterEmployeesRoutes(router, employeesService); err != nil {
		panic(err)
	}

	// Requisito 6 - Buyers
	buyersRepo := repository.NewBuyerDb(nil)
	buyersService := service.NewBuyer(buyersRepo)
	// Create the routes and deps
	if err = routes.BuyerRoutes(router, buyersService); err != nil {
		panic(err)
	}

	// Sprint2 Requisito 2 - Locality
	localitiesRepo := repository.NewMysqlLocalityRepository(a.db)
	localityService := service.NewMysqlLocalityService(localitiesRepo)
	if err = routes.LocalityRoutes(router, localityService); err != nil {
		panic(err)
	}

	// Sprint2 Requisito 2 - Carry
	carriesRepo := repository.NewMySQLCarryRepository(a.db)
	carryService := service.NewMySQLCarryService(carriesRepo, localitiesRepo)
	if err = routes.CarryRoutes(router, carryService); err != nil {
		panic(err)
	}

	return nil
}

// Run runs the application.
func (a *ApplicationDefault) Run() (err error) {
	defer a.db.Close()
	log.Printf("starting server at %s\n", a.cfgAddr)

	err = http.ListenAndServe(a.cfgAddr, a.router)
	return
}
