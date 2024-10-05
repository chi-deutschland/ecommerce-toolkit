package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

	ShipperName     ColumnMapping `json:"shipperName"`
	ShipperAddress  ColumnMapping `json:"shipperAddress"`
	ShipperState    ColumnMapping `json:"shipperState"`
	ShipperPostcode ColumnMapping `json:"shipperPostcode"`
	ShipperCountry  ColumnMapping `json:"shipperCountry"`

	RecipientName     ColumnMapping `json:"recipientName"`
	RecipientAddress  ColumnMapping `json:"recipientAddress"`
	RecipientState    ColumnMapping `json:"recipientState"`
	RecipientPostcode ColumnMapping `json:"recipientPostcode"`
	RecipientCountry  ColumnMapping `json:"recipientCountry"`

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

	http.Handle("/schema", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

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

	http.Handle("/pipeline", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

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

	http.Handle("/pipelines/{pipeline}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

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

	http.Handle("/pipelines/{pipeline}/input", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

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

		log.Debug().Interface("waybill", waybill).Msg("masterwaybill")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "message": "Hello Tim :)" }`))
	}))

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal().Err(err).Msg("http server")
	}
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
	waybill := NewMasterWaybill()
	for rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		if len(columns) == 0 {
			continue
		}

		if pipeline.MasterWaybillNumber.Column != nil {
			masterNumber := SanitizeMawb(columns[*pipeline.MasterWaybillNumber.Column])
			prefix, suffix := SplitMawb(masterNumber)
			waybill.WaybillNumber = suffix
			waybill.WaybillPrefix = prefix
		} else if pipeline.MasterWaybillNumber.UseFilename != nil {
			masterNumber := SanitizeMawb(strings.TrimSuffix(filename, ".xlsx"))
			prefix, suffix := SplitMawb(masterNumber)
			waybill.WaybillNumber = suffix
			waybill.WaybillPrefix = prefix
		} else if pipeline.MasterWaybillNumber.Constant != nil {
			prefix, suffix := SplitMawb(SanitizeMawb(*pipeline.MasterWaybillNumber.Constant))
			waybill.WaybillNumber = suffix
			waybill.WaybillPrefix = prefix
		}

		houseWaybill := &Waybill{
			WaybillType: WaybillTypeHouse,
			Type:        "cargo:Waybill",
			Shipment: &Shipment{
				Type: "Shipment",
			},
		}
		if pipeline.HouseWaybillNumber.Column != nil {
			houseWaybill.WaybillNumber = columns[*pipeline.HouseWaybillNumber.Column]
		} else if pipeline.HouseWaybillNumber.Constant != nil {
			houseWaybill.WaybillNumber = *pipeline.HouseWaybillNumber.Constant
		}

		if pipeline.ShippingReference.Column != nil {
			houseWaybill.ShippingRef = columns[*pipeline.ShippingReference.Column]
		}

		if pipeline.TotalShipmentGrossWeight.Column != nil {
			houseWaybill.Shipment.TotalGrossWeight = &Value{
				Type:           "Value",
				NumericalValue: columns[*pipeline.TotalShipmentGrossWeight.Column],
				Unit: &Unit{
					Type:              "CodeListElement",
					Code:              "KGM",
					CodeListReference: "https://vocabulary.uncefact.org/WeightUnitMeasureCode",
				},
			}
		}

		waybill.HouseWaybills = append(waybill.HouseWaybills, houseWaybill)
	}
	return waybill, nil
}

func sendWaybill(waybill *Waybill) error {
	body, err := json.Marshal(waybill)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "TODO: URl", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header = map[string][]string{
		"Accept":        {"application/ld+json"},
		"Authorization": {"Bearer {{neone_server_access_token}}"},
		"Content-Type":  {"application/ld+json"},
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Status %s, Body: %s", resp.Status, string(body))
	}

	return nil
}

type Waybill struct {
	WaybillNumber string      `json:"cargo:waybillNumber,omitempty"`
	WaybillPrefix string      `json:"cargo:waybillPrefix,omitempty"`
	WaybillType   WaybillType `json:"cargo:waybillType,omitempty"`
	HouseWaybills []*Waybill  `json:"cargo:houseWaybills,omitempty"`
	ShippingRef   string      `json:"cargo:shippingRefNo,omitempty"`
	Shipment      *Shipment   `json:"cargo:shipment,omitempty"`

	// JSON-LD stuff
	Context Context `json:"@context"`
	Type    string  `json:"@type"`
}

type Context struct {
	Cargo string `json:"cargo"`
}

func NewMasterWaybill() *Waybill {
	return &Waybill{
		Context: Context{
			Cargo: "https://onerecord.iata.org/ns/cargo#",
		},
		Type: "cargo:waybill",
	}
}

type Shipment struct {
	Type             string `json:"@type"`
	TotalGrossWeight *Value `json:"totalGrossWeight,omitempty"`
}

type Value struct {
	Type           string `json:"@type"`
	NumericalValue string `json:"numericalValue,omitempty"`
	Unit           *Unit  `json:"unit,omitempty"`
}

type Unit struct {
	Type              string `json:"@type"`
	Code              string `json:"code,omitempty"`
	CodeListReference string `json:"codeListReference,omitempty"`
}

type Piece struct{}

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
			ShipperAddress: ColumnMapping{
				Title:   "Shipper Address Line 1",
				Content: "$E:SHIPPER ADD 1",
				Column:  Ptr(AlphaToIndex("E")),
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
			RecipientAddress: ColumnMapping{
				Title:   "Recipient Address Line 1",
				Content: "$M:RECEIPIENT ADD 1",
				Column:  Ptr(AlphaToIndex("M")),
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
			RecipientCountry: ColumnMapping{
				Title:   "Recipient State/County/Province",
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
	for _, r := range b {
		index += int('A') - int(r)
	}
	return index
}

func Ptr[T any](v T) *T {
	return &v
}
