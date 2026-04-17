package soap

import "encoding/xml"

// Account SOAP response types

type getAccessStatusResponse struct {
	XMLName xml.Name `xml:"getAccessStatusResponse"`
	Return  int      `xml:"getAccessStatusReturn"`
}

type getCurrentRevisionResponse struct {
	XMLName xml.Name `xml:"getCurrentRevisionResponse"`
	Return  int      `xml:"getCurrentRevisionReturn"`
}

type getExpireDateResponse struct {
	XMLName xml.Name `xml:"getExpireDateResponse"`
	Return  string   `xml:"getExpireDateReturn"`
}

type getSubscriptionStatusResponse struct {
	XMLName xml.Name `xml:"getSubscriptionStatusResponse"`
	Return  string   `xml:"getSubscriptionStatusReturn"`
}

type getRightAccessResponse struct {
	XMLName xml.Name `xml:"getRightAccessResponse"`
	Return  string   `xml:"getRightAccessReturn"`
}

type getUserIdByLoginResponse struct {
	XMLName xml.Name `xml:"getUserIdByLoginResponse"`
	Return  string   `xml:"getUserIdByLoginReturn"`
}
