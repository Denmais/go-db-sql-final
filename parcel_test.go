package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// Проверка на отсутствие ошибки при добавлении, наличие идентификатора
	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.Greater(t, id, 0)
	// get
	// ПРоверка на отсутствие ошибок, соответствие полей добавленной посылки
	parc, err := store.Get(id)
	require.NoError(t, err)
	parcel.Number = id
	assert.Equal(t, parcel, parc)

	// delete
	// ПРоверка на то, что посылку нельзя получить из БД
	err = store.Delete(id)
	require.NoError(t, err)
	_, err = store.Get(id)
	assert.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()
	// add
	// Проверка на отсутствие ошибки при добавлении, наличие идентификатора
	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.Greater(t, id, 0)
	// set address
	// Проверка на отсутствие ошибки при обновлении адреса
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)
	// check
	// Проверка обновленного адреса
	parc, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, newAddress, parc.Address)

}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// Проверка на отсутствие ошибки при добавлении, наличие идентификатора
	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.Greater(t, id, 0)
	// set status
	// Проверка на отсутствие ошибки при обновлении статуса
	newStatus := ParcelStatusSent
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)
	// delete new status
	// Проверка невозможности удаления посылки с новым статусом
	err = store.Delete(id)
	require.NoError(t, err)
	_, err = store.Get(id)
	require.NoError(t, err)
	// update new status
	// Проверка невозможности изменения адреса посылки с новым статусом
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)
	parc, _ := store.Get(id)
	assert.NotEqual(t, newAddress, parc.Address)
	// check
	// Проверка обновленного статуса
	parc, err = store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, newStatus, parc.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	//	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i]) // Проверка на отсутствие ошибки при добавлении, наличие идентификатора
		require.NoError(t, err)
		assert.Greater(t, id, 0)
		parcels[i].Number = id

	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	// Проверка на остутствие ошибок и количество добавленных записей
	require.NoError(t, err)
	assert.Len(t, storedParcels, len(parcels))

	// check
	// Проверка добавленных элементов
	assert.ElementsMatch(t, storedParcels, parcels)
}
