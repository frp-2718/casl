package entities

import "testing"

func TestValid(t *testing.T) {
	validLocations := provideValidAlmaLocations()
	invalidLocations := provideInvalidAlmaLocations()

	for _, l := range validLocations {
		t.Run(l.Library_name, func(t *testing.T) {
			if !ValidLocation(l) {
				t.Errorf("got false; want true")
			}
		})
	}
	for _, l := range invalidLocations {
		t.Run(l.Library_name, func(t *testing.T) {
			if ValidLocation(l) {
				t.Errorf("got true; want false")
			}
		})
	}
}

func provideInvalidAlmaLocations() []AlmaLocation {
	var res []AlmaLocation

	validItem1 := AlmaItem{Process_name: "Missing", Process_code: "MISSING", Status: "Item in place"}
	invalidItem1 := AlmaItem{Process_name: "Acquisition", Process_code: "ACQ", Status: "Item not in place"}
	invalidItem2 := AlmaItem{Process_name: "Acquisition", Process_code: "ACQ", Status: "Item in place"}

	invalidLocation1 := AlmaLocation{Library_name: "no discovery", NoDiscovery: true, Items: []*AlmaItem{&validItem1}}
	res = append(res, invalidLocation1)

	invalidLocation2 := AlmaLocation{Library_name: "no item", NoDiscovery: false}
	res = append(res, invalidLocation2)

	invalidLocation3 := AlmaLocation{Library_name: "empty items", NoDiscovery: false, Items: []*AlmaItem{}}
	res = append(res, invalidLocation3)

	invalidLocation4 := AlmaLocation{Library_name: "one invalid", NoDiscovery: false, Items: []*AlmaItem{&invalidItem1}}
	res = append(res, invalidLocation4)

	invalidLocation5 := AlmaLocation{Library_name: "two invalid", NoDiscovery: false, Items: []*AlmaItem{&invalidItem1, &invalidItem2}}
	res = append(res, invalidLocation5)

	return res
}

func provideValidAlmaLocations() []AlmaLocation {
	var res []AlmaLocation

	validItem1 := AlmaItem{Process_name: "Missing", Process_code: "MISSING", Status: "Item in place"}
	validItem2 := AlmaItem{Status: "Item not in place"}
	invalidItem1 := AlmaItem{Process_name: "Acquisition", Process_code: "ACQ", Status: "Item not in place"}

	validLocation1 := AlmaLocation{Library_name: "one valid", NoDiscovery: false, Items: []*AlmaItem{&validItem1}}
	res = append(res, validLocation1)

	validLocation2 := AlmaLocation{Library_name: "two valid", NoDiscovery: false, Items: []*AlmaItem{&validItem1, &validItem2}}
	res = append(res, validLocation2)

	validLocation3 := AlmaLocation{Library_name: "valid/invalid", NoDiscovery: false, Items: []*AlmaItem{&validItem1, &invalidItem1}}
	res = append(res, validLocation3)

	return res
}
