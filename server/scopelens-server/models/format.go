package models

// Supported formats

var formats = []string{
	"VGC 2020",
	"VGC 2019 Ultra Series",
	"VGC 2019 Moon Series",
	"VGC 2019 Sun Series",
	"VGC 2018",
	"VGC 2017",
	"VGC 2016",
	"VGC 2015",
	"VGC 2014",
	"VGC 2013",
	"VGC 2012",
	"VGC 2011",
	"VGC 2010",
	"[Gen 8] Battle Stadium Singles",
	"[Gen 8] Ubers",
	"[Gen 8] OU",
	"[Gen 8] UU",
	"[Gen 8] RU",
	"[Gen 8] NU",
	"[Gen 8] LC",
	"[Gen 8] Anything Goes",
	"[Gen 8] Doubles OU",
	"[Gen 8] Doubles UU",
	"[Gen 8] Battle Stadium Doubles",
	"Others",
}

// Get supported formats
func GetFormats() []string {
	return formats
}