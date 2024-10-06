package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

type Pipeline struct {
	Name    string  `json:"name"`
	Mapping *Schema `json:"mapping"`
}

type SchemaSuggestion struct {
	Schema
	Headers []string `json:"headers"`
}

type Schema struct {
	MasterWaybillNumber ColumnMapping `json:"waybill"`
	HouseWaybillNumber  ColumnMapping `json:"houseWaybill"`
	ShippingReference   ColumnMapping `json:"shippingReference"`
	BoxNumber           ColumnMapping `json:"boxNumber"`

	ShipperName         ColumnMapping `json:"shipperName"`
	ShipperAddressLine1 ColumnMapping `json:"shipperAddressLine1"`
	ShipperAddressLine2 ColumnMapping `json:"shipperAddressLine2"`
	ShipperAddressLine3 ColumnMapping `json:"shipperAddressLine3"`
	ShipperCity         ColumnMapping `json:"shipperCity"`
	ShipperState        ColumnMapping `json:"shipperState"`
	ShipperPostcode     ColumnMapping `json:"shipperPostcode"`
	ShipperCountry      ColumnMapping `json:"shipperCountry"`

	RecipientName         ColumnMapping `json:"recipientName"`
	RecipientAddressLine1 ColumnMapping `json:"recipientAddressLine1"`
	RecipientAddressLine2 ColumnMapping `json:"recipientAddressLine2"`
	RecipientAddressLine3 ColumnMapping `json:"recipientAddressLine3"`
	RecipientCity         ColumnMapping `json:"recipientCity"`
	RecipientState        ColumnMapping `json:"recipientState"`
	RecipientPostcode     ColumnMapping `json:"recipientPostcode"`
	RecipientCounty       ColumnMapping `json:"recipientCounty"`
	RecipientCountry      ColumnMapping `json:"recipientCountry"`

	TotalShipmentGrossWeight ColumnMapping `json:"totalShipmentGrossWeight"`
	ItemUnitPriceConcurrency ColumnMapping `json:"itemUnitPriceConcurrency"`
	ProductSKU               ColumnMapping `json:"productSku"`
	ShipmentGoodsDescription ColumnMapping `json:"shipmentGoodsDescription"`
	ProductHSCode            ColumnMapping `json:"productHSCode"`
	ItemQuantity             ColumnMapping `json:"itemQuantity"`
	ItemPrice                ColumnMapping `json:"itemPrice"`
}

type ColumnMapping struct {
	Title   string `json:"title"`
	Content string `json:"content"`

	Column      *int    `json:"column"`
	Constant    *string `json:"constant"`
	UseFilename *bool   `json:"useFilename"`
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	mux := http.NewServeMux()

	mux.Handle("/schema", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := excelize.OpenReader(r.Body)
		if err != nil {
			log.Err(err).Msg("read excel")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sheetList := file.GetSheetList()
		if len(sheetList) == 0 {
			log.Err(err).Msg("get sheet list")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		firstSheet := sheetList[0]

		rows, err := file.Rows(firstSheet)
		if err != nil {
			log.Err(err).Msg("create row iterator")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !rows.Next() {
			log.Err(err).Msg("get sheet list")
			w.WriteHeader(http.StatusBadRequest)
			return

		}

		headers, err := rows.Columns()
		if err != nil {
			log.Err(err).Msg("get sheet list")
			w.WriteHeader(http.StatusBadRequest)
			return

		}

		mapping := headersToMapping(headers)

		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(mapping); err != nil {
			log.Err(err).Msg("get sheet list")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}))

	mux.Handle("/pipeline", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Err(err).Msg("error")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var pipeline Pipeline
		if err := json.Unmarshal(body, &pipeline); err != nil {
			log.Err(err).Msg("error")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		filename := fmt.Sprintf("%s/%s.json", PIPELINE_DIR, pipeline.Name)

		if err := os.WriteFile(filename, body, 0644); err != nil {
			log.Err(err).Msg("error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}))

	mux.Handle("/pipelines", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		files, err := os.ReadDir(PIPELINE_DIR)
		if err != nil {
			log.Err(err).Msg("error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		pipelines := make([]*Pipeline, 0, len(files))
		for _, file := range files {
			f, err := os.Open(fmt.Sprintf("%s/%s", PIPELINE_DIR, file.Name()))
			if err != nil {
				log.Err(err).Msg("read pipelines")
				continue
			}
			dec := json.NewDecoder(f)
			var pipeline Pipeline
			if err := dec.Decode(&pipeline); err != nil {
				log.Err(err).Msg("decode pipelines")
				continue
			}
			pipelines = append(pipelines, &pipeline)
		}

		enc := json.NewEncoder(w)
		if err := enc.Encode(&pipelines); err != nil {
			log.Err(err).Msg("write pipelines")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}))

	mux.Handle("/pipelines/{pipeline}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		pipelineName := r.PathValue("pipeline")
		pipeline, err := readPipeline(pipelineName)
		if err != nil {
			log.Err(err).Msg("error")
			if errors.Is(err, os.ErrExist) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		enc := json.NewEncoder(w)
		if err := enc.Encode(pipeline); err != nil {
			log.Err(err).Msg("error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}))

	mux.Handle("/pipelines/{pipeline}/input", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pipelineName := r.PathValue("pipeline")
		pipeline, err := readPipeline(pipelineName)
		if err != nil {
			log.Err(err).Msg("error")
			if errors.Is(err, os.ErrExist) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		file, err := excelize.OpenReader(r.Body)
		if err != nil {
			log.Err(err).Msg("read excel")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sheetList := file.GetSheetList()
		if len(sheetList) == 0 {
			log.Err(err).Msg("get sheet list")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		firstSheet := sheetList[0]

		rows, err := file.Rows(firstSheet)
		if err != nil {
			log.Err(err).Msg("create row iterator")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		filename := r.Header.Get("Content-Disposition")

		waybill, err := excelToOneRecord(pipeline.Mapping, rows, filename)
		if err != nil {
			log.Err(err).Msg("transform")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		logisticsObjectUrl, err := sendWaybill(waybill)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(err.Error())); err != nil {
				log.Err(err)
			}
			return
		}

		if _, err := fmt.Fprintf(w, "Logistics Object URL: %s\n", logisticsObjectUrl); err != nil {
			log.Err(err)
		}

	}))

	mux.Handle("/ai/{hscode}/{term}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hscode := r.PathValue("hscode")
		term := r.PathValue("term")
		ok, err := verifyHSCodeWithAI(hscode, term)
		if err != nil {
			log.Err(err).Msg("ai")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "%t", ok)
	}))

	if err := http.ListenAndServe(":80", logMiddleware(CORSMiddleware(mux))); err != nil {
		log.Fatal().Err(err).Msg("http server")
	}
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Str("method", r.Method).Str("url", r.URL.String()).Msg("HTTP request")
		next.ServeHTTP(w, r)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // change this never :P
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")

		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}

		next.ServeHTTP(w, r)
	})
}

const PIPELINE_DIR = "pipelines"

func readPipeline(pipelineName string) (*Pipeline, error) {
	filename := fmt.Sprintf("%s/%s.json", PIPELINE_DIR, pipelineName)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(file)
	pipeline := &Pipeline{}
	if err := dec.Decode(pipeline); err != nil {
		return nil, err
	}
	return pipeline, err
}

func excelToOneRecord(pipeline *Schema, rows *excelize.Rows, filename string) (*Waybill, error) {
	masterWaybill := NewMasterWaybill()

	i := 0
	for rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		if len(columns) == 0 {
			continue
		}

		houseWaybillNumber := columns[*pipeline.HouseWaybillNumber.Column]
		if len(houseWaybillNumber) < 3 || houseWaybillNumber == pipeline.HouseWaybillNumber.Content[3:] {
			// skip header line
			continue
		}

		if pipeline.MasterWaybillNumber.Column != nil {
			masterNumber := SanitizeMawb(columns[*pipeline.MasterWaybillNumber.Column])
			prefix, suffix := SplitMawb(masterNumber)
			masterWaybill.WaybillNumber = suffix
			masterWaybill.WaybillPrefix = prefix
		} else if pipeline.MasterWaybillNumber.UseFilename != nil {
			masterNumber := SanitizeMawb(strings.TrimSuffix(filename, ".xlsx"))
			prefix, suffix := SplitMawb(masterNumber)
			masterWaybill.WaybillNumber = suffix
			masterWaybill.WaybillPrefix = prefix
		} else if pipeline.MasterWaybillNumber.Constant != nil {
			// TODO: review. There is no reason to have a constant master waybill number throughout the pipeline,
			//	but there might be a reason to have a constant value such as currency and destination country
			prefix, suffix := SplitMawb(SanitizeMawb(*pipeline.MasterWaybillNumber.Constant))
			masterWaybill.WaybillNumber = suffix
			masterWaybill.WaybillPrefix = prefix
		}

		var houseWaybill *Waybill

		item := newItem(
			columns[*pipeline.ProductSKU.Column],
			columns[*pipeline.ProductHSCode.Column],
			columns[*pipeline.ItemQuantity.Column],
			columns[*pipeline.ItemPrice.Column],
			columns[*pipeline.ItemUnitPriceConcurrency.Column],
		)

		piece := newPiece(
			[]*Item{item},
			columns[*pipeline.BoxNumber.Column],
			columns[*pipeline.ShipmentGoodsDescription.Column],
		)

		if masterWaybill.HouseWaybills[HouseWaybillNumber(houseWaybillNumber)] != nil {
			houseWaybill = masterWaybill.HouseWaybills[HouseWaybillNumber(houseWaybillNumber)]
			houseWaybill.Shipment.Pieces = append(houseWaybill.Shipment.Pieces, piece)
		} else {
			houseWaybill = newHouseWaybill()
			if pipeline.ShippingReference.Column != nil {
				houseWaybill.ShippingRef = columns[*pipeline.ShippingReference.Column]
			}

			arrivalCountryCode := "DE" // TODO: implement dynamic arrival country code in pipeline
			arrivalRegionCode := columns[*pipeline.RecipientCounty.Column]
			arrivalStreetAddressLines := []string{
				columns[*pipeline.RecipientAddressLine1.Column],
				columns[*pipeline.RecipientAddressLine2.Column],
				columns[*pipeline.RecipientAddressLine3.Column],
				columns[*pipeline.RecipientCity.Column],
				columns[*pipeline.RecipientPostcode.Column],
			}
			houseWaybill.ArrivalLocation = newLocation(arrivalCountryCode, arrivalRegionCode, arrivalStreetAddressLines)

			departureCountryCode := columns[*pipeline.ShipperCountry.Column]
			departureRegionCode := columns[*pipeline.ShipperState.Column]
			departureStreetAddressLines := []string{
				columns[*pipeline.ShipperAddressLine1.Column],
				columns[*pipeline.ShipperAddressLine2.Column],
				columns[*pipeline.ShipperAddressLine3.Column],
				columns[*pipeline.ShipperCity.Column],
				columns[*pipeline.ShipperPostcode.Column],
			}
			houseWaybill.DepartureLocation = newLocation(departureCountryCode, departureRegionCode, departureStreetAddressLines)

			shipper := newShipper(columns[*pipeline.ShipperName.Column])
			customer := newCustomer(columns[*pipeline.RecipientName.Column])

			houseWaybill.InvolvedParties = []*Party{shipper, customer}

			houseWaybill.Shipment = newShipment(
				[]*Piece{piece},
				columns[*pipeline.TotalShipmentGrossWeight.Column],
			)

			if masterWaybill.HouseWaybills == nil {
				masterWaybill.HouseWaybills = make(map[HouseWaybillNumber]*Waybill)
			}
			masterWaybill.HouseWaybills[HouseWaybillNumber(houseWaybillNumber)] = houseWaybill
		}

		i++
	}
	return masterWaybill, nil
}

func newShipment(pieces []*Piece, totalGrossWeight string) *Shipment {
	return &Shipment{
		Type:             "cargo:Shipment",
		TotalGrossWeight: newTotalGrossWeight(totalGrossWeight),
		Pieces:           pieces,
	}
}

func newTotalGrossWeight(totalGrossWeight string) *Value {
	return newValue(totalGrossWeight, newCodeListElement("KGM", weightUnitCodeListReference, ""))
}

func newCustomer(customerName string) *Party {
	customerContact := "CUSTOMER_CONTACT"
	return newParty(customerName, &customerContact, nil)
}

type LogisticsAgent struct {
	Name        string  `json:"cargo:name,omitempty"`
	ContactRole *string `json:"cargo:contactRole,omitempty"`
	Type        string  `json:"@type"`
}

type Party struct {
	PartyDetails *LogisticsAgent  `json:"cargo:partyDetails,omitempty"`
	PartyRole    *CodeListElement `json:"cargo:partyRole,omitempty"` // TODO: bug: does not show up in NE-ONE Play. String also does not work.
	Type         string           `json:"@type"`
}

func newShipper(shipperName string) *Party {
	return newParty(
		shipperName,
		nil,
		newCodeListElement("SHP", iataCoreCodeListReference, iataCoreCodeListVersion),
	)
}

func newParty(logisticsAgentName string, contactRole *string, partyRole *CodeListElement) *Party {
	return &Party{
		PartyDetails: newLogisticsAgent(logisticsAgentName, contactRole),
		PartyRole:    partyRole,
		Type:         "cargo:Party",
	}
}

func newLogisticsAgent(name string, contactRole *string) *LogisticsAgent {
	return &LogisticsAgent{
		Name:        name,
		ContactRole: contactRole,
		Type:        "cargo:LogisticsAgent",
	}
}

func newLocation(countryCode, regionCode string, streetAddressLines []string) *Location {
	return &Location{
		Address: newAddress(countryCode, regionCode, streetAddressLines),
		Type:    "cargo:Location",
	}
}

func newAddress(countryCode, regionCode string, streetAddressLines []string) *Address {
	return &Address{
		Country:            newCountry(countryCode),
		RegionCode:         newRegionCode(regionCode),
		StreetAddressLines: streetAddressLines,
		Type:               "cargo:Address",
	}
}

func newRegionCode(regionCode string) *CodeListElement {
	return newCodeListElement(regionCode, regionCodeListReference, regionCodeListVersion)
}

// TODO: consider using specific type instead of strings
const (
	regionCodeListReference = "https://www.iso.org/obp/ui/#iso:code:3166"
	regionCodeListVersion   = "3166-2"

	countryCodeListReference = "https://vocabulary.uncefact.org/CountryId"

	iataCoreCodeListReference = "https://onerecord.iata.org/ns/coreCodeLists"
	iataCoreCodeListVersion   = "1.0.0"

	hsCodeType = "UN Standard International Trade Classification"

	hsCodeListReference = "www.tariffnumber.com"
	hsCodeListVersion   = "2024"

	currencyCodeListReference = "https://vocabulary.uncefact.org/RevisedCurrencyCode"

	pieceUnitCode              = "H87"
	pieceUnitCodeListReference = "https://docs.peppol.eu/poacc/billing/3.0/codelist/UNECERec20/"
	pieceUnitCodeListVersion   = "Revision 11e"

	weightUnitCodeListReference = "https://vocabulary.uncefact.org/WeightUnitMeasureCode"
)

func newCodeListElement(code, reference, version string) *CodeListElement {
	return &CodeListElement{
		Type:              "cargo:CodeListElement",
		Code:              code,
		CodeListVersion:   version,
		CodeListReference: reference,
	}
}

func newCountry(countryCode string) *CodeListElement {
	return newCodeListElement(countryCode, countryCodeListReference, "")
}

func newHouseWaybill() *Waybill {
	return &Waybill{
		Context:     nil,
		WaybillType: WaybillTypeHouse,
		Type:        "cargo:Waybill",
	}
}

func sendWaybill(waybill *Waybill) (string, error) {
	body, err := json.Marshal(waybill)
	if err != nil {
		return "", err
	}

	logisticsObjectEndpoint := "http://your_ne_one_server_here/logistics-objects"

	req, err := http.NewRequest(http.MethodPost, logisticsObjectEndpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	auth := "Bearer your_token_here"

	req.Header = map[string][]string{
		"Accept":        {"application/ld+json"},
		"Authorization": {auth},
		"Content-Type":  {"application/ld+json"},
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("Status %s, Body: %s", resp.Status, string(body))
	}

	return resp.Header.Get("Location"), nil
}

type HouseWaybillNumber string

type CodeListElement struct {
	Code              string `json:"cargo:code,omitempty"`
	CodeListReference string `json:"cargo:codeListReference,omitempty"`
	CodeListVersion   string `json:"cargo:codeListVersion,omitempty"`
	Type              string `json:"@type"`
}

type Address struct {
	Country            *CodeListElement `json:"cargo:country,omitempty"`
	RegionCode         *CodeListElement `json:"cargo:regionCode,omitempty"`
	StreetAddressLines []string         `json:"cargo:streetAddressLines,omitempty"`
	Type               string           `json:"@type"`
}

type Location struct {
	Address *Address `json:"cargo:Address,omitempty"`
	Type    string   `json:"@type"`
}

type Waybill struct {
	ArrivalLocation   *Location                       `json:"cargo:arrivalLocation,omitempty"`
	DepartureLocation *Location                       `json:"cargo:departureLocation,omitempty"`
	WaybillNumber     string                          `json:"cargo:waybillNumber,omitempty"`
	WaybillPrefix     string                          `json:"cargo:waybillPrefix,omitempty"`
	WaybillType       WaybillType                     `json:"cargo:waybillType,omitempty"`
	HouseWaybills     map[HouseWaybillNumber]*Waybill `json:"cargo:houseWaybills,omitempty"`
	InvolvedParties   []*Party                        `json:"cargo:involvedParties,omitempty"`
	ShippingRef       string                          `json:"cargo:shippingRefNo,omitempty"`
	Shipment          *Shipment                       `json:"cargo:shipment,omitempty"`

	// JSON-LD stuff
	Context *Context `json:"@context,omitempty"`
	Type    string   `json:"@type"`
}

func (w *Waybill) MarshalJSON() ([]byte, error) {
	type Alias Waybill
	aux := &struct {
		Context       *Context   `json:"@context,omitempty"`
		HouseWaybills []*Waybill `json:"cargo:houseWaybills,omitempty"`
		*Alias
	}{
		Alias:   (*Alias)(w),
		Context: w.Context,
	}

	if w.HouseWaybills != nil {
		aux.HouseWaybills = []*Waybill{}
		for key, value := range w.HouseWaybills {
			value.WaybillNumber = string(key)
			aux.HouseWaybills = append(aux.HouseWaybills, value)
		}
	}

	return json.Marshal(aux)
}

type Context struct {
	Cargo string `json:"cargo,omitempty"`
}

func NewMasterWaybill() *Waybill {
	return &Waybill{
		Context: &Context{
			Cargo: "https://onerecord.iata.org/ns/cargo#",
		},
		Type:        "cargo:Waybill",
		WaybillType: WaybillTypeMaster,
	}
}

type Shipment struct {
	Pieces           []*Piece `json:"cargo:pieces,omitempty"`
	Type             string   `json:"@type"`
	TotalGrossWeight *Value   `json:"cargo:totalGrossWeight,omitempty"`
}

type Value struct {
	Type           string           `json:"@type"`
	NumericalValue string           `json:"cargo:numericalValue,omitempty"`
	Unit           *CodeListElement `json:"cargo:unit,omitempty"` // TODO: some CodeListElement units on the ne one server string may work
}

type OtherIdentifier struct {
	Type                string `json:"@type"`
	OtherIdentifierType string `json:"cargo:otherIdentifierType,omitempty"`
	TextualValue        string `json:"cargo:textualValue,omitempty"`
}

func NewOtherIdentifier(otherIdentifierType, textualValue string) *OtherIdentifier {
	return &OtherIdentifier{
		Type:                "cargo:OtherIdentifier",
		OtherIdentifierType: otherIdentifierType,
		TextualValue:        textualValue,
	}
}

type Product struct {
	OtherIdentifiers []*OtherIdentifier `json:"cargo:otherIdentifiers,omitempty"`
	HsCode           *CodeListElement   `json:"cargo:hsCode,omitempty"`
	HsType           string             `json:"cargo:hsType,omitempty"`
	Type             string             `json:"@type"`
}

func NewProduct(skuNumber, hsCode string) *Product {
	var otherIdentifiers []*OtherIdentifier
	if skuNumber != "" {
		otherIdentifiers = []*OtherIdentifier{NewOtherIdentifier("SKU", skuNumber)}
	}

	return &Product{
		OtherIdentifiers: otherIdentifiers,
		HsCode:           newHsCode(hsCode),
		HsType:           hsCodeType,
		Type:             "cargo:Product",
	}
}

func newHsCode(hsCode string) *CodeListElement {
	return newCodeListElement(hsCode, hsCodeListReference, hsCodeListVersion)
}

type Item struct {
	OfProduct    *Product `json:"cargo:ofProduct,omitempty"`
	ItemQuantity *Value   `json:"cargo:itemQuantity,omitempty"`
	UnitPrice    *Value   `json:"cargo:unitPrice,omitempty"`
	Type         string   `json:"@type"`
}

func newItem(skuNumber, hsCode, itemQuantity, itemPrice, currency string) *Item {
	return &Item{
		OfProduct:    NewProduct(skuNumber, hsCode),
		ItemQuantity: newValue(itemQuantity, newPieceUnit()),
		UnitPrice:    newValue(itemPrice, newCodeListElement(currency, currencyCodeListReference, "")),
		Type:         "cargo:Item",
	}
}

func newValue(numericalValue string, unitCodeList *CodeListElement) *Value {
	return &Value{
		Type:           "cargo:Value",
		NumericalValue: numericalValue,
		Unit:           unitCodeList,
	}
}

func newPieceUnit() *CodeListElement {
	return newCodeListElement(pieceUnitCode, pieceUnitCodeListReference, pieceUnitCodeListVersion)
}

type Piece struct {
	Type             string             `json:"@type"`
	ContainedItems   []*Item            `json:"cargo:containedItems,omitempty"`
	OtherIdentifiers []*OtherIdentifier `json:"cargo:otherIdentifiers,omitempty"`
	GoodsDescription string             `json:"cargo:goodsDescription,omitempty"`
}

func newPiece(items []*Item, boxNumber, goodsDescription string) *Piece {
	return &Piece{
		Type:             "cargo:Piece",
		ContainedItems:   items,
		OtherIdentifiers: []*OtherIdentifier{NewOtherIdentifier("Box Number", boxNumber)},
		GoodsDescription: goodsDescription,
	}
}

type WaybillType string

func (t *WaybillType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "MASTER":
		*t = WaybillTypeMaster
	case "HOUSE":
		*t = WaybillTypeHouse
	case "DIRECT":
		*t = WaybillTypeDirect
	default:
		return fmt.Errorf("unknown waybill type: %s", s)
	}
	return nil
}

const (
	WaybillTypeMaster WaybillType = "MASTER"
	WaybillTypeHouse  WaybillType = "HOUSE"
	WaybillTypeDirect WaybillType = "DIRECT"
)

func SanitizeMawb(number string) string {
	sanitizedMawb := strings.ReplaceAll(number, "-", "")
	sanitizedMawb = strings.ReplaceAll(sanitizedMawb, " ", "")
	if utf8.RuneCountInString(sanitizedMawb) < 11 {
		sanitizedMawb = fmt.Sprintf("%011s", sanitizedMawb)
	}
	return sanitizedMawb
}

func SplitMawb(number string) (string, string) {
	return number[:3], number[3:]
}

func headersToMapping(headers []string) *SchemaSuggestion {
	// AI magic here
	return &SchemaSuggestion{
		Schema: Schema{
			MasterWaybillNumber: ColumnMapping{
				Title:       "Master Waybill Number",
				Content:     "Filename",
				UseFilename: Ptr(true),
			},
			HouseWaybillNumber: ColumnMapping{
				Title:   "House Waybill Number",
				Content: "$A:TRACKING NO.",
				Column:  Ptr(AlphaToIndex("A")),
			},
			ShippingReference: ColumnMapping{
				Title:   "HAWB Shipping Reference Number",
				Content: "$B:CUSTOMER REF",
				Column:  Ptr(AlphaToIndex("B")),
			},
			BoxNumber: ColumnMapping{
				Title:   "Box Number",
				Content: "$C:Box Number",
				Column:  Ptr(AlphaToIndex("C")),
			},
			ShipperName: ColumnMapping{
				Title:   "Shipper Name",
				Content: "$D:SENDER NAME",
				Column:  Ptr(AlphaToIndex("D")),
			},
			ShipperAddressLine1: ColumnMapping{
				Title:   "Shipper Address Line 1",
				Content: "$E:SHIPPER ADD 1",
				Column:  Ptr(AlphaToIndex("E")),
			},
			ShipperAddressLine2: ColumnMapping{
				Title:   "Shipper Address Line 2",
				Content: "$F:SHIPPER ADD 2",
				Column:  Ptr(AlphaToIndex("F")),
			},
			ShipperAddressLine3: ColumnMapping{
				Title:   "Shipper Address Line 3",
				Content: "$G:SHIPPER ADD 3",
				Column:  Ptr(AlphaToIndex("G")),
			},
			ShipperCity: ColumnMapping{
				Title:   "Shipper City",
				Content: "$H:SENDER CITY",
				Column:  Ptr(AlphaToIndex("H")),
			},
			ShipperState: ColumnMapping{
				Title:   "Shipper State/County/Province",
				Content: "$I:Ship State",
				Column:  Ptr(AlphaToIndex("I")),
			},
			ShipperPostcode: ColumnMapping{
				Title:   "Shipper Address Line 5",
				Content: "$J:SENDER POSTCODE",
				Column:  Ptr(AlphaToIndex("J")),
			},
			ShipperCountry: ColumnMapping{
				Title:   "Shipper Country",
				Content: "$K:SENDER COUNTRY",
				Column:  Ptr(AlphaToIndex("K")),
			},
			RecipientName: ColumnMapping{
				Title:   "Recipient Name",
				Content: "$L:RECEIPIENT NAME",
				Column:  Ptr(AlphaToIndex("L")),
			},
			RecipientAddressLine1: ColumnMapping{
				Title:   "Recipient Address Line 1",
				Content: "$M:RECEIPIENT ADD 1",
				Column:  Ptr(AlphaToIndex("M")),
			},
			RecipientAddressLine2: ColumnMapping{
				Title:   "Recipient Address Line 2",
				Content: "$N:RECEIPIENT ADD 2",
				Column:  Ptr(AlphaToIndex("N")),
			},
			RecipientAddressLine3: ColumnMapping{
				Title:   "Recipient Address Line 3",
				Content: "$O:RECEIPIENT ADD 3",
				Column:  Ptr(AlphaToIndex("O")),
			},
			RecipientCity: ColumnMapping{
				Title:   "Recipient City",
				Content: "$P:RECEIPIENT CITY",
				Column:  Ptr(AlphaToIndex("P")),
			},
			RecipientState: ColumnMapping{
				Title:   "Recipient State/County/Province",
				Content: "$Q:RECEIPIENT COUNTY",
				Column:  Ptr(AlphaToIndex("Q")),
			},
			RecipientPostcode: ColumnMapping{
				Title:   "Recipient Address Line 5",
				Content: "$R:RECEIPIENT POSTCODE",
				Column:  Ptr(AlphaToIndex("R")),
			},
			RecipientCounty: ColumnMapping{
				Title:   "Recipient County",
				Content: "$Q:RECEIPIENT COUNTY",
				Column:  Ptr(AlphaToIndex("Q")),
			},
			TotalShipmentGrossWeight: ColumnMapping{
				Title:   "Total Shipment Gross Weight",
				Content: "$U:GROSS WEIGHT (KG)",
				Column:  Ptr(AlphaToIndex("U")),
			},
			ItemUnitPriceConcurrency: ColumnMapping{
				Title:   "Item Unit Price Currency",
				Content: "$W:Currency",
				Column:  Ptr(AlphaToIndex("W")),
			},
			ProductSKU: ColumnMapping{
				Title:   "Product SKU",
				Content: "$X:SKU NUMBER",
				Column:  Ptr(AlphaToIndex("X")),
			},
			ShipmentGoodsDescription: ColumnMapping{
				Title:   "Shipment Goods Description",
				Content: "$Y:ITEM DESCRIPTION1",
				Column:  Ptr(AlphaToIndex("Y")),
			},
			ProductHSCode: ColumnMapping{
				Title:   "Product HS Code",
				Content: "$Z:ITEM HS CODE",
				Column:  Ptr(AlphaToIndex("Z")),
			},
			ItemQuantity: ColumnMapping{
				Title:   "Item Quantity",
				Content: "$AA:ITEM QUANTITY",
				Column:  Ptr(AlphaToIndex("AA")),
			},
			ItemPrice: ColumnMapping{
				Title:   "Item Price",
				Content: "$AB:UNIT VALUE",
				Column:  Ptr(AlphaToIndex("AB")),
			},
		},
		Headers: headers,
	}
}

func AlphaToIndex(b string) int {
	index := 0
	for i, r := range b {
		index += int(r) - int('A') + i*26
	}
	return index
}

func Ptr[T any](v T) *T {
	return &v
}

type AIAnswer struct {
	Suggestions []AISuggestion `json:"suggestions"`
}

type AISuggestion struct {
	Code  string
	Score float32
}

func verifyHSCodeWithAI(hscode, description string) (bool, error) {
	url := fmt.Sprintf("https://www.tariffnumber.com/api/v2/cnSuggest?term=%s&lang=en", url.PathEscape(description))
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	var answer AIAnswer
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&answer); err != nil {
		return false, err
	}

	for len(hscode) > 8 {
		hscode = hscode[:8]

	}
	for _, s := range answer.Suggestions {
		if strings.HasPrefix(s.Code, hscode) {
			return true, nil
		}
	}
	return false, nil
}
