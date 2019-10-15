package obs

type Issue struct {
	Desc    string `xml:"issue"`
	Id      string `xml:"id,attr"`
	Tracker string `xml:"tracker,attr"`
}
type Patchinfo struct {
	Incident    string  `xml:"incident,attr"`
	Issues      []Issue `xml:"issue,attr"`
	Category    string  `xml:"category"`
	Rating      string  `xml:"rating"`
	Packager    string  `xml:"packager"`
	Description string  `xml:"description"`
	Summary     string  `xml:"summary"`
}
