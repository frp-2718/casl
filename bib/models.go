package bib

type BibRecord struct {
	ppn            string
	mms            string
	sudocLocations []sudocLocation
	almaLocations  []almaLocation
}

type sudocLocation struct {
	iln  string
	rcr  string
	name string
}

type almaLocation struct {
	collection string
	ownerCode  string
	rcr        string
}

// CRecord is the fusion of a SUDOC record and an Alma record.
type CRecord struct {
	PPN          string
	AlmaLibrary  string
	SUDOCLibrary string
	ILN          string
	RCR          string
	InAlma       bool
	InSUDOC      bool
}
