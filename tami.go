package tami

import (
	"sort"
	"time"
)

type Transaction struct {
	Price     float64
	ItemID    interface{}
	Timestamp time.Time
}

type IndexValueHistoryItem struct {
	Price       float64
	ItemID      interface{}
	IndexValue  float64
	Transaction Transaction
}

type IndexValueHistoryItemWithRatio struct {
	IndexValueHistoryItem IndexValueHistoryItem
	IndexRatio            float64
}

// SortTransactions returns the given list of transactions sorted in chronological order.
func SortTransactions(transactionHistory []Transaction) []Transaction {
	sort.Slice(transactionHistory, func(i, j int) bool {
		return transactionHistory[i].Timestamp.Before(transactionHistory[j].Timestamp)
	})
	return transactionHistory
}

// FilterValidTransactions returns only transactions that have at least two sales in the last year, and at least one
// sale in the last six months from the given list of transactions.
func FilterValidTransactions(transactionHistory []Transaction) []Transaction {
	now := time.Now()
	oneYearAgo := now.AddDate(-1, 0, 0)
	sixMonthsAgo := now.AddDate(0, -6, 0)

	type include struct {
		PastYearSaleCount      int
		HasSaleInLastSixMonths bool
		IsValid                bool
	}

	inclusionMap := map[interface{}]*include{}

	for _, transaction := range transactionHistory {
		if _, ok := inclusionMap[transaction.ItemID]; !ok {
			inclusionMap[transaction.ItemID] = &include{
				PastYearSaleCount:      0,
				HasSaleInLastSixMonths: false,
				IsValid:                false,
			}
		}

		currentMapItem := inclusionMap[transaction.ItemID]
		if currentMapItem.IsValid {
			continue
		}

		if transaction.Timestamp.Before(oneYearAgo) {
			continue
		}

		currentMapItem.PastYearSaleCount++

		if transaction.Timestamp.Before(sixMonthsAgo) {
			continue
		}

		currentMapItem.HasSaleInLastSixMonths = true

		if currentMapItem.PastYearSaleCount >= 2 {
			currentMapItem.IsValid = true
		}

	}

	var res []Transaction

	for _, transaction := range transactionHistory {
		if inclusionMap[transaction.ItemID].IsValid {
			res = append(res, transaction)
		}
	}

	return res
}

// CreateIndexValueHistory creates a list that contains the index value at the time of each transaction, and includes
// the transaction as well.
func CreateIndexValueHistory(transactionHistory []Transaction) []IndexValueHistoryItem {
	transactionMap := map[interface{}]Transaction{}

	lastIndexValue, lastDivisor := float64(0), float64(1)

	var result []IndexValueHistoryItem

	for i := 0; i < len(transactionHistory); i++ {
		transaction := transactionHistory[i]
		_, isNotFirstSale := transactionMap[transaction.ItemID]
		isFirstSale := !isNotFirstSale

		transactionMap[transaction.ItemID] = transaction

		itemCount := len(transactionMap)

		allLastSoldValue := float64(0)

		for _, t := range transactionMap {
			allLastSoldValue += t.Price
		}

		indexValue := float64(allLastSoldValue) / (float64(itemCount) * lastDivisor)

		if i == 0 {
			lastIndexValue = indexValue
			result = append(result, IndexValueHistoryItem{
				ItemID:      transaction.ItemID,
				Price:       transaction.Price,
				IndexValue:  indexValue,
				Transaction: transaction,
			})

			continue
		}

		nextDivisor := lastDivisor
		if isFirstSale {
			nextDivisor = lastDivisor * (indexValue / lastIndexValue)
		}

		weightedIndexValue := float64(allLastSoldValue) / (float64(itemCount) * nextDivisor)

		lastIndexValue = weightedIndexValue
		lastDivisor = nextDivisor
		result = append(result, IndexValueHistoryItem{
			ItemID:      transaction.ItemID,
			Price:       transaction.Price,
			IndexValue:  weightedIndexValue,
			Transaction: transaction,
		})
	}

	return result
}

// GetIndexValue returns the index value of the last item from a given list of transactions.
func GetIndexValue(indexValueHistory []IndexValueHistoryItem) float64 {
	if len(indexValueHistory) == 0 {
		return 0
	}
	return indexValueHistory[len(indexValueHistory)-1].IndexValue
}

// GetIndexRatios calculates the index ratio for the last transaction of each item in the list of IndexValueHistoryItem.
// GetIndexRatios returns a list of IndexValueHistoryItemWithRatio where each instance is the IndexValueHistoryItem
// with an additional `IndexRatio` property added.
func GetIndexRatios(indexValueHistory []IndexValueHistoryItem) []IndexValueHistoryItemWithRatio {
	lastSaleMap := map[interface{}]IndexValueHistoryItem{}
	for _, item := range indexValueHistory {
		lastSaleMap[item.ItemID] = item
	}

	var res []IndexValueHistoryItemWithRatio
	for _, v := range lastSaleMap {
		res = append(res, IndexValueHistoryItemWithRatio{
			IndexValueHistoryItem: IndexValueHistoryItem{
				Price:       v.Price,
				ItemID:      v.ItemID,
				IndexValue:  v.IndexValue,
				Transaction: v.Transaction,
			},
			IndexRatio: float64(v.Price) / v.IndexValue,
		})
	}
	return res
}

// TAMI calculates the Time Adjusted Market Index for a list of transactions for a given collection.
func TAMI(transactionHistory []Transaction) float64 {
	sortedTransactions := SortTransactions(transactionHistory)
	validTransactions := FilterValidTransactions(sortedTransactions)
	indexValueHistory := CreateIndexValueHistory(validTransactions)
	indexValue := GetIndexValue(indexValueHistory)
	indexRatios := GetIndexRatios(indexValueHistory)
	var timeAdjustedValues []float64
	for _, item := range indexRatios {
		timeAdjustedValues = append(timeAdjustedValues, indexValue*item.IndexRatio)
	}
	timeAdjustedMarketIndex := float64(0)
	for _, v := range timeAdjustedValues {
		timeAdjustedMarketIndex += v
	}
	return timeAdjustedMarketIndex
}
