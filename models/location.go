package models

//COPY locations FROM '/Users/JonathanThom/Desktop/IP2LOCATION-LITE-DB11.CSV/IP2LOCATION-LITE-DB11.CSV' WITH CSV QUOTE AS '"' HEADER;

type Location struct {
	IpFrom      int     `json:"-"`
	IpTo        int     `json:"-"`
	CountryCode string  `json:"-"`
	CountryName string  `json:"country"`
	RegionName  string  `json:"-"`
	CityName    string  `json:"city"`
	Latitude    float64 `json:"-"`
	Longitude   float64 `json:"-"`
	ZipCode     string  `json:"-"`
	TimeZone    string  `json:"-"`
}
