package repository

import (
	"encoding/json"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"log"
	"os"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type BuyerRepo struct {
	buyerTable map[int]internal.Buyer
}

func NewBuyerDb(buyerTab map[int]internal.Buyer) *BuyerRepo {

	BuyerDb := make(map[int]internal.Buyer)
	if buyerTab != nil {
		BuyerDb = buyerTab
	}
	return &BuyerRepo{buyerTable: BuyerDb}
}

var buyersFile = "./docs/db/buyers.json"

func (repo *BuyerRepo) LoadBuyers() (map[int]internal.Buyer, error) {
	file, err := os.ReadFile(buyersFile)
	if err != nil {
		log.Println("Error to read file", err)
		return nil, err
	}

	var buyers []internal.Buyer

	if err := json.Unmarshal(file, &buyers); err != nil {
		log.Println("Error to unmarshal - ")
		return nil, err
	}

	repo.buyerTable = make(map[int]internal.Buyer)
	for _, buyer := range buyers {
		repo.buyerTable[int(buyer.ID)] = buyer
	}

	return repo.buyerTable, nil
}

func (repo *BuyerRepo) GetAll() ([]internal.Buyer, error) {
	buyersMap, err := repo.LoadBuyers()
	if err != nil {
		return nil, err
	}

	var buyers []internal.Buyer
	for _, buyer := range buyersMap {
		buyers = append(buyers, buyer)
	}

	return buyers, nil
}

func (repo *BuyerRepo) GetOne(id int) (*internal.Buyer, error) {
	buyersMap, err := repo.LoadBuyers()
	if err != nil {
		log.Println("Error to Load Buyers - ", err)
		return nil, err
	}

	if buyer, exists := buyersMap[id]; exists {
		return &buyer, err
	}
	return nil, err
}

func (repo *BuyerRepo) CreateBuyer(newBuyer internal.Buyer) (*internal.Buyer, error) {
	buyers, err := repo.GetAll()

	if err != nil {
		log.Println("Error to load - ", err)
		return nil, err
	}

	for _, buyer := range buyers {
		if buyer.ID == newBuyer.ID {
			log.Println("There is an user with this ID", err)
			return nil, utils.ErrConflict
		}
	}
	buyers = append(buyers, newBuyer)

	file, err := os.Create(buyersFile)
	if err != nil {
		log.Println("Erro ao reabrir o arquivo:", err)
		return nil, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(buyers); err != nil {
		log.Println("Erro ao codificar JSON:", err)
		return nil, err
	}
	log.Println("Comprador salvo!")

	return &newBuyer, nil
}

func (repo *BuyerRepo) UpdateBuyer(updatedBuyer *internal.Buyer) (*internal.Buyer, error) {
	buyers, err := repo.GetAll()

	if err != nil {
		log.Println("Error to load - ", err)
		return nil, err
	}

	for i, buyer := range buyers {
		if buyer.ID == updatedBuyer.ID {
			buyers[i] = *updatedBuyer
		}
	}

	file, err := os.Create(buyersFile)

	if err != nil {
		log.Println("Error to open file:", err)
		return nil, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(buyers); err != nil {
		log.Println("Error to encoder JSON:", err)
		return nil, err
	}
	log.Println("Comprador atualizado!")

	return updatedBuyer, nil
}

func (repo *BuyerRepo) DeleteBuyer(id int) error {
	buyers, err := repo.GetAll()
	if err != nil {
		log.Println("Erro ao carregar os compradores:", err)
		return err
	}

	indexToDelete := -1
	for i, buyer := range buyers {
		if int(buyer.ID) == id {
			indexToDelete = i
			break
		}
	}

	if indexToDelete == -1 {
		log.Println("Comprador com o ID fornecido não encontrado")
		return utils.ErrInvalidArguments
	}

	buyers = append(buyers[:indexToDelete], buyers[indexToDelete+1:]...)

	file, err := os.Create(buyersFile)
	if err != nil {
		log.Println("Erro ao abrir o arquivo para escrita:", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(buyers); err != nil {
		log.Println("Erro ao codificar JSON:", err)
		return err
	}

	log.Println("Comprador excluído com sucesso!")
	return nil
}
