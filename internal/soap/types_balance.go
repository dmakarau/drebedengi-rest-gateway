package soap

import "encoding/xml"

type balanceKV struct {
	Key   string `xml:"key"`
	Value string `xml:"value"`
}

type balanceMapEntry struct {
	Items []balanceKV `xml:"item"`
}

type getBalanceResponse struct {
	XMLName xml.Name          `xml:"getBalanceResponse"`
	Return  []balanceMapEntry `xml:"getBalanceReturn>item"`
}
