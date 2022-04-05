package bib

func in(s string, ens []string) bool {
	for _, elem := range ens {
		if s == elem {
			return true
		}
	}
	return false
}

func equalBibRecord(br1, br2 BibRecord) bool {
	if br1.ppn != br2.ppn || br1.mms != br2.mms {
		return false
	}
	if len(br1.almaLocations) != len(br2.almaLocations) {
		return false
	}
	if len(br1.sudocLocations) != len(br2.sudocLocations) {
		return false
	}

	for i, sl := range br1.sudocLocations {
		if sl.iln != br2.sudocLocations[i].iln ||
			sl.rcr != br2.sudocLocations[i].rcr ||
			sl.name != br2.sudocLocations[i].name {
			return false
		}
	}

	for i, al := range br1.almaLocations {
		if al.collection != br2.almaLocations[i].collection ||
			al.ownerCode != br2.almaLocations[i].ownerCode ||
			al.rcr != br2.almaLocations[i].rcr {
			return false
		}
	}
	return true
}

func equalBibRecords(b1, b2 []BibRecord) bool {
	if len(b1) != len(b2) {
		return false
	}
	for i, record := range b1 {
		if !equalBibRecord(record, b2[i]) {
			return false
		}
	}
	return true
}

func equalStrings(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, s := range s1 {
		if s != s2[i] {
			return false
		}
	}
	return true
}

func equalCRecord(c1, c2 CRecord) bool {
	if c1.PPN == c2.PPN && c1.AlmaLibrary == c2.AlmaLibrary &&
		c1.SUDOCLibrary == c2.SUDOCLibrary && c1.ILN == c2.ILN &&
		c1.RCR == c2.RCR && c1.InAlma == c2.InAlma && c1.InSUDOC == c2.InSUDOC {
		return true
	}
	return false
}

func equalCRecords(c1, c2 []CRecord) bool {
	if len(c1) != len(c2) {
		return false
	}
	for i, record := range c1 {
		if !equalCRecord(record, c2[i]) {
			return false
		}
	}
	return true
}
