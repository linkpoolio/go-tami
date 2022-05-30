package tami_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/linkpoolio/go-tami"
)

var (
	now                    = time.Now()
	yesterday              = now.AddDate(0, 0, -1)
	twoDaysAgo             = now.AddDate(0, 0, -2)
	threeDaysAgo           = now.AddDate(0, 0, -3)
	oneMonthAgo            = now.AddDate(0, -1, 0)
	sixWeeksAgo            = now.AddDate(0, 0, -6*7)
	twoYearsAgo            = now.AddDate(-2, 0, 0)
	eightMonthsAgo         = now.AddDate(0, -8, 0)
	nineMonthsAgo          = now.AddDate(0, -9, 0)
	mockTransactionHistory = []tami.Transaction{
		{ItemID: "Lavender", Price: 500, Timestamp: threeDaysAgo},
		{ItemID: "Hyacinth", Price: 700, Timestamp: oneMonthAgo},
		{ItemID: "Mars", Price: 1200, Timestamp: twoDaysAgo},
		{ItemID: "Nyx", Price: 612, Timestamp: twoYearsAgo},
		{ItemID: "Hyacinth", Price: 400, Timestamp: threeDaysAgo},
		{ItemID: "Nyx", Price: 1200, Timestamp: yesterday},
		{ItemID: "Mars", Price: 612, Timestamp: sixWeeksAgo},
	}
	olderTransactions = []tami.Transaction{
		{ItemID: "Hyacinth", Price: 700, Timestamp: eightMonthsAgo},
		{ItemID: "Hyacinth", Price: 400, Timestamp: nineMonthsAgo},
		{ItemID: "Mars", Price: 612, Timestamp: sixWeeksAgo},
		{ItemID: "Mars", Price: 500, Timestamp: twoDaysAgo},
		{ItemID: "Mars", Price: 999, Timestamp: yesterday},
	}
	sortedTransactions = tami.SortTransactions(mockTransactionHistory)
	validTransactions  = tami.FilterValidTransactions(sortedTransactions)
)

func TestCreateIndexValueHistory(t *testing.T) {
	indexValueHistory := tami.CreateIndexValueHistory(validTransactions)
	want := []tami.IndexValueHistoryItem{
		{
			ItemID:     "Mars",
			Price:      612,
			IndexValue: 612,
			Transaction: tami.Transaction{
				ItemID:    "Mars",
				Price:     612,
				Timestamp: sixWeeksAgo,
			},
		},
		{
			ItemID:     "Hyacinth",
			Price:      700,
			IndexValue: 612,
			Transaction: tami.Transaction{
				ItemID:    "Hyacinth",
				Price:     700,
				Timestamp: oneMonthAgo,
			},
		},
		{
			ItemID:     "Hyacinth",
			Price:      400,
			IndexValue: 472.0609756097561,
			Transaction: tami.Transaction{
				ItemID:    "Hyacinth",
				Price:     400,
				Timestamp: threeDaysAgo,
			},
		},
		{
			ItemID:     "Mars",
			Price:      1200,
			IndexValue: 746.3414634146342,
			Transaction: tami.Transaction{
				ItemID:    "Mars",
				Price:     1200,
				Timestamp: twoDaysAgo,
			},
		},
	}
	if diff := cmp.Diff(indexValueHistory, want); diff != "" {
		t.Fatalf("unexpected value history: %s", diff)
	}
}

func TestGetIndexRatios(t *testing.T) {
	indexValueHistory := tami.CreateIndexValueHistory(validTransactions)
	indexRatios := tami.GetIndexRatios(indexValueHistory)
	want := []tami.IndexValueHistoryItemWithRatio{
		{
			IndexValueHistoryItem: tami.IndexValueHistoryItem{
				ItemID:     "Mars",
				Price:      1200,
				IndexValue: 746.3414634146342,
				Transaction: tami.Transaction{
					ItemID:    "Mars",
					Price:     1200,
					Timestamp: twoDaysAgo,
				},
			},
			IndexRatio: 1.6078431372549018,
		},
		{
			IndexValueHistoryItem: tami.IndexValueHistoryItem{
				ItemID:     "Hyacinth",
				Price:      400,
				IndexValue: 472.0609756097561,
				Transaction: tami.Transaction{
					ItemID:    "Hyacinth",
					Price:     400,
					Timestamp: threeDaysAgo,
				},
			},
			IndexRatio: 0.847348161926167,
		},
	}
	if diff := cmp.Diff(indexRatios, want); diff != "" {
		t.Fatalf("unexpected index ratios: %s", diff)
	}
}

func TestGetIndexValue(t *testing.T) {
	t.Run("fetches the index value for a non-empty list of transactions", func(t *testing.T) {
		indexValueHistory := tami.CreateIndexValueHistory(validTransactions)
		indexValue := tami.GetIndexValue(indexValueHistory)
		if indexValue != 746.3414634146342 {
			t.Fatalf("unexpected index value: %f", indexValue)
		}
	})
	t.Run("fetches a zero index value for an empty list of transactions", func(t *testing.T) {
		indexValue := tami.GetIndexValue(nil)
		if indexValue != 0 {
			t.Fatalf("unexpected index value: %f", indexValue)
		}
	})
}

func TestFilterValidTransactions(t *testing.T) {
	filteredTransactions := tami.FilterValidTransactions(olderTransactions)
	want := []tami.Transaction{
		{
			Price:     612,
			ItemID:    "Mars",
			Timestamp: sixWeeksAgo,
		},
		{
			Price:     500,
			ItemID:    "Mars",
			Timestamp: twoDaysAgo,
		},
		{
			Price:     999,
			ItemID:    "Mars",
			Timestamp: yesterday,
		},
	}
	if diff := cmp.Diff(filteredTransactions, want); diff != "" {
		t.Fatalf("unexpected index ratios: %s", diff)
	}
}

func TestTAMI(t *testing.T) {
	value := tami.TAMI(validTransactions)
	if value != 1832.411067193676 {
		t.Fatalf("unexpected tami: %f", value)
	}
}
