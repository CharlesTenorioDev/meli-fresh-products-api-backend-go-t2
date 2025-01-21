package section

// import (
// 	"github.com/meli-fresh-products-api-backend-go-t2/internal"
// 	"testing"

// 	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
// 	"github.com/stretchr/testify/require"
// )

// var simpleSection = internal.Section{
// 	ID:                 1,
// 	SectionNumber:      1,
// 	CurrentTemperature: 1,
// 	MinimumTemperature: 1,
// 	CurrentCapacity:    1,
// 	MinimumCapacity:    1,
// 	MaximumCapacity:    1,
// 	WarehouseID:        1,
// 	ProductTypeID:      1,
// }

// func Test_GetAll(t *testing.T) {
// 	repo := NewMemorySectionRepository(map[int]internal.Section{
// 		1: simpleSection,
// 	})
// 	sections, _ := repo.GetAll()
// 	require.Equal(t, 1, len(sections))
// 	require.Equal(t, simpleSection, sections[0])
// }

// func Test_GetById_WhenExists(t *testing.T) {
// 	repo := NewMemorySectionRepository(map[int]internal.Section{
// 		1: simpleSection,
// 	})
// 	section, _ := repo.GetByID(1)
// 	require.Equal(t, simpleSection, section)
// }

// func Test_GetById_WhenNotExists(t *testing.T) {
// 	repo := NewMemorySectionRepository(nil)
// 	section, _ := repo.GetByID(1)
// 	require.Empty(t, section)
// }

// func Test_GetBySectionNumber_WhenExists(t *testing.T) {
// 	repo := NewMemorySectionRepository(map[int]internal.Section{
// 		1: simpleSection,
// 	})
// 	section, _ := repo.GetBySectionNumber(1)
// 	require.Equal(t, simpleSection, section)
// }

// func Test_GetBySectionNumber_WhenNotExists(t *testing.T) {
// 	repo := NewMemorySectionRepository(nil)
// 	section, _ := repo.GetBySectionNumber(99)
// 	require.Empty(t, section)
// }

// func Test_Delete_WhenExist(t *testing.T) {
// 	repo := NewMemorySectionRepository(map[int]internal.Section{
// 		1: simpleSection,
// 	})
// 	err := repo.Delete(1)
// 	require.Nil(t, err)
// }

// func Test_Delete_WhenNotExist(t *testing.T) {
// 	repo := NewMemorySectionRepository(nil)
// 	err := repo.Delete(1)
// 	require.ErrorIs(t, err, utils.ErrNotFound)
// }

// func Test_Save_WhenLoadedDb(t *testing.T) {
// 	repo := NewMemorySectionRepository(map[int]internal.Section{
// 		1: simpleSection,
// 	})
// 	newSection, _ := repo.Save(internal.Section{
// 		SectionNumber:      2,
// 		CurrentTemperature: 2,
// 		MinimumTemperature: 2,
// 		CurrentCapacity:    2,
// 		MinimumCapacity:    2,
// 		MaximumCapacity:    2,
// 		WarehouseID:        2,
// 		ProductTypeID:      2,
// 	})
// 	require.Equal(t, 2, newSection.ID)
// }

// func Test_Save_WhenEmptyDb(t *testing.T) {
// 	repo := NewMemorySectionRepository(nil)
// 	newSection, _ := repo.Save(internal.Section{
// 		SectionNumber:      2,
// 		CurrentTemperature: 2,
// 		MinimumTemperature: 2,
// 		CurrentCapacity:    2,
// 		MinimumCapacity:    2,
// 		MaximumCapacity:    2,
// 		WarehouseID:        2,
// 		ProductTypeID:      2,
// 	})
// 	require.Equal(t, 1, newSection.ID)
// }

// func Test_Update_WhenExists(t *testing.T) {
// 	repo := NewMemorySectionRepository(map[int]internal.Section{
// 		1: simpleSection,
// 	})
// 	sectionToUpdate := simpleSection
// 	sectionToUpdate.MaximumCapacity = 99
// 	updatedSection, _ := repo.Update(sectionToUpdate)
// 	require.Equal(t, 99, updatedSection.MaximumCapacity)
// 	require.Equal(t, sectionToUpdate, updatedSection)
// }
