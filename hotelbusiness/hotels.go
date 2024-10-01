//go:build !solution

package hotelbusiness

import "sort"

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	if len(guests) == 0 {
		return []Load{}
	}

	delta := make(map[int]int)
	for _, guest := range guests {
		delta[guest.CheckInDate]++
		delta[guest.CheckOutDate]--
	}

	dates := make([]int, 0, len(delta))
	for date, change := range delta {
		if change != 0 {
			dates = append(dates, date)
		}
	}
	sort.Ints(dates)

	result := []Load{}
	var sum int

	for _, date := range dates {
		sum += delta[date]
		result = append(result, Load{
			StartDate:  date,
			GuestCount: sum,
		})
	}

	return result
}
