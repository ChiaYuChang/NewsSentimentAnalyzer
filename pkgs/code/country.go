package code

type CountryCode string

func (cc CountryCode) IsEmpty() bool {
	return cc == ""
}

const (
	CCCanada        CountryCode = "ca"
	CCChina         CountryCode = "cn"
	CCJapan         CountryCode = "jp"
	CCKorea         CountryCode = "kr"
	CCTaiwan        CountryCode = "tw"
	CCUnitedStates  CountryCode = "us"
	CCUnitedKingdom CountryCode = "gb"
)
