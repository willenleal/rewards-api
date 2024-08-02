package api

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"unicode"
)

func PointsCalculator(receipt Receipt) (int64, error) {
	totalPoints := 0

	rules := []func(Receipt) (int, error){
		pointsForOddPurchaseDate,
		pointsForRetailerName,
		pointsForTotalPrice,
		pointsRelativeToItemDescription,
		pointsPerTwoItemsOnReceipt,
		pointsForTimeWindow,
	}

	for _, rule := range rules {
		points, err := rule(receipt)
		if err != nil {
			return 0, err
		}

		totalPoints += points
	}

	return int64(totalPoints), nil
}

func pointsForTimeWindow(Receipt Receipt) (int, error) {

	points := 0
	parts := strings.Split(Receipt.PurchaseTime, ":")

	if len(parts) != 2 {
		return 0, errors.New("invalid time format")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	if (hours == 14 && minutes >= 0) || (hours == 15) || (hours == 16 && minutes == 0) {
		points += 10
		return points, nil
	}
	return 0, nil
}

func pointsForOddPurchaseDate(receipt Receipt) (int, error) {
	points := 0

	// If day is odd add 6 points
	tempDay := receipt.PurchaseDate[len(receipt.PurchaseDate)-2:]
	day, err := strconv.Atoi(tempDay)

	if err != nil {
		return 0, err
	}

	if day%2 != 0 {
		points += 6
	}

	return points, nil

}

func pointsRelativeToItemDescription(receipt Receipt) (int, error) {
	points := 0
	// If the trimmed length of the item description is a multiple of 3,
	//multiply the price by 0.2 and round up to the nearest integer. The
	//result is the number of points earned.
	for _, item := range receipt.Items {
		trimmmedDescription := strings.TrimSpace(item.ShortDescription)

		if len(trimmmedDescription)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return 0, err
			}
			points += int(math.Ceil(price * 0.2))
		}
	}

	return points, nil

}

func pointsPerTwoItemsOnReceipt(receipt Receipt) (int, error) {
	points := 0
	// 5 points for every two items on receipt
	itemsCount := len(receipt.Items)
	points += int((itemsCount / 2) * 5)

	return points, nil

}

func pointsForTotalPrice(receipt Receipt) (int, error) {
	points := 0

	receiptTotal, err := strconv.ParseFloat(receipt.Total, 64)

	if err != nil {
		return 0, err
	}

	// Determine if total is round dollar amount
	if isMultipleOf(receiptTotal, 1.0) {
		points += 50
	}

	// Determine if totla multiple of 0.25
	if isMultipleOf(receiptTotal, 0.25) {
		points += 25
	}

	return points, nil
}

func pointsForRetailerName(receipt Receipt) (int, error) {
	points := 0

	// Count alphanumeric chars on retailer's name
	for _, char := range receipt.Retailer {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			points++
		}
	}

	return points, nil
}

func isMultipleOf(value, multiple float64) bool {
	remainder := math.Mod(value, multiple)
	tolerance := 1e-9
	return math.Abs(remainder) < tolerance
}
