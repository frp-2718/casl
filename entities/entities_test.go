package entities

import "testing"

//	type AlmaLocation struct {
//		Library_name  string
//		Library_code  string
//		Location_name string
//		Location_code string
//		Call_number   string
//		NoDiscovery   bool
//		Items         []*AlmaItem
//	}
func provideAlmaLocations() []AlmaLocation {
	var res []AlmaLocation

	validItem1 := AlmaItem{Process_name: "Missing", Process_code: "MISSING", Status: "Item in place"}
	validItem2 := AlmaItem{Status: "Item not in place"}
	invalidItem1 := AlmaItem{Process_name: "Acquisition", Process_code: "ACQ", Status: "Item not in place"}
	invalidItem2 := AlmaItem{Process_name: "Acquisition", Process_code: "ACQ", Status: "Item in place"}

	validLocation1 := AlmaLocation{NoDiscovery: false, Items: []*AlmaItem{&validItem1}}
	res = append(res, validLocation1)

	validLocation2 := AlmaLocation{NoDiscovery: false, Items: []*AlmaItem{&validItem1, &validItem2}}
	res = append(res, validLocation2)

	invalidLocation1 := AlmaLocation{NoDiscovery: true, Items: []*AlmaItem{&validItem1}}
	res = append(res, invalidLocation1)

	invalidLocation2 := AlmaLocation{NoDiscovery: false}
	res = append(res, invalidLocation2)

	invalidLocation3 := AlmaLocation{NoDiscovery: false, Items: []*AlmaItem{}}
	res = append(res, invalidLocation3)

	invalidLocation4 := AlmaLocation{NoDiscovery: false, Items: []*AlmaItem{&invalidItem1}}
	res = append(res, invalidLocation4)

	invalidLocation5 := AlmaLocation{NoDiscovery: true, Items: []*AlmaItem{&invalidItem1, &invalidItem2}}
	res = append(res, invalidLocation5)

	return res
}

func TestValid(t *testing.T) {
}
