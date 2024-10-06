"use client";
import React, { useState } from 'react';

import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
  } from "@/components/ui/dialog"

import { Button } from '@/components/ui/button';
import ReactJson from 'react-json-view';

const jsonData = {
  "@context": {
    "cargo": "https://onerecord.iata.org/ns/cargo#"
  },
  "@type": "cargo:Waybill",
  "cargo:houseWaybills": {
    "H0483A0710460500": {
      "@type": "cargo:Waybill",
      "cargo:arrivalLocation": {
        "@type": "cargo:Location",
        "cargo:Address": {
          "@type": "cargo:Address",
          "cargo:country": {
            "@type": "cargo:CodeListElement",
            "cargo:code": "DE",
            "cargo:codeListReference": "https://vocabulary.uncefact.org/CountryId"
          },
          "cargo:regionCode": {
            "@type": "cargo:CodeListElement",
            "cargo:code": "Frankfurt",
            "cargo:codeListReference": "https://www.iso.org/obp/ui/#iso:code:3166",
            "cargo:codeListVersion": "3166-2"
          },
          "cargo:streetAddressLines": [
            "115 Broadway",
            "",
            "",
            "Exeter",
            "60000"
          ]
        }
      },
      "cargo:departureLocation": {
        "@type": "cargo:Location",
        "cargo:Address": {
          "@type": "cargo:Address",
          "cargo:country": {
            "@type": "cargo:CodeListElement",
            "cargo:code": "CN",
            "cargo:codeListReference": "https://vocabulary.uncefact.org/CountryId"
          },
          "cargo:regionCode": {
            "@type": "cargo:CodeListElement",
            "cargo:code": "GD",
            "cargo:codeListReference": "https://www.iso.org/obp/ui/#iso:code:3166",
            "cargo:codeListVersion": "3166-2"
          },
          "cargo:streetAddressLines": [
            "North side of Chuangxin street, Sihui City",
            "",
            "",
            "Istanbul",
            "526200"
          ]
        }
      },
      "cargo:involvedParties": [
        {
          "@type": "cargo:Party",
          "cargo:partyDetails": {
            "@type": "cargo:LogisticsAgent",
            "cargo:name": "ZQ01"
          },
          "cargo:partyRole": {
            "@type": "cargo:CodeListElement",
            "cargo:code": "SHP",
            "cargo:codeListReference": "https://onerecord.iata.org/ns/coreCodeLists",
            "cargo:codeListVersion": "1.0.0"
          }
        },
        {
          "@type": "cargo:Party",
          "cargo:partyDetails": {
            "@type": "cargo:LogisticsAgent",
            "cargo:contactRole": "CUSTOMER_CONTACT",
            "cargo:name": "Marie Seco"
          }
        }
      ],
      "cargo:shipment": {
        "@type": "cargo:Shipment",
        "cargo:pieces": [
          {
            "@type": "cargo:Piece",
            "cargo:containedItems": [
              {
                "@type": "cargo:Item",
                "cargo:itemQuantity": {
                  "@type": "cargo:Value",
                  "cargo:numericalValue": "1",
                  "cargo:unit": {
                    "@type": "cargo:CodeListElement",
                    "cargo:code": "H87",
                    "cargo:codeListReference": "https://docs.peppol.eu/poacc/billing/3.0/codelist/UNECERec20/",
                    "cargo:codeListVersion": "Revision 11e"
                  }
                },
                "cargo:ofProduct": {
                  "@type": "cargo:Product",
                  "cargo:hsCode": {
                    "@type": "cargo:CodeListElement",
                    "cargo:code": "3926400000",
                    "cargo:codeListReference": "www.tariffnumber.com",
                    "cargo:codeListVersion": "2024"
                  },
                  "cargo:hsType": "UN Standard International Trade Classification",
                  "cargo:otherIdentifiers": [
                    {
                      "@type": "cargo:OtherIdentifier",
                      "cargo:otherIdentifierType": "SKU"
                    }
                  ]
                }
              }
            ],
            "cargo:goodsDescription": "wall hanging",
            "cargo:otherIdentifiers": [
              {
                "@type": "cargo:OtherIdentifier",
                "cargo:otherIdentifierType": "Box Number",
                "cargo:textualValue": "UKBA240828843714"
              }
            ]
          }
        ],
        "cargo:totalGrossWeight": {
          "@type": "cargo:Value",
          "cargo:numericalValue": "0.475",
          "cargo:unit": {
            "@type": "cargo:CodeListElement",
            "cargo:code": "KGM",
            "cargo:codeListReference": "https://vocabulary.uncefact.org/WeightUnitMeasureCode"
          }
        }
      },
      "cargo:shippingRefNo": "BG-2408245623JCYGJE",
      "cargo:waybillType": "HOUSE"
    }
  },
  "cargo:waybillNumber": "93976750",
  "cargo:waybillPrefix": "250",
  "cargo:waybillType": "MASTER"
};

export default function IntegrationPage() {
  const [showAlert, setShowAlert] = useState(false);

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
    setShowAlert(true);
    setTimeout(() => setShowAlert(false), 2000);
  };

  return (
    <div>
      <div className="flex items-center justify-center mt-24">
        <div className="flex items-center">
          <div className="flex items-center justify-center w-10 h-10 rounded-full bg-green-600 text-xl">✓</div>
          <div className="w-32 h-1 bg-green-600 mx-2"></div>
          <div className="flex items-center justify-center w-10 h-10 rounded-full bg-green-600 text-xl">✓</div>
          <div className="w-32 h-1 bg-green-600 mx-2"></div>
          <div className="flex items-center justify-center w-10 h-10 rounded-full bg-green-600 text-xl">✓</div>
        </div>
      </div>
      <div className="flex justify-center items-center min-h-screen position-relative">
        <div className="flex flex-col items-center gap-5">
          <div className="container p-5 inline-block rounded-lg">
            <div className="flex items-center mb-16 mt-8">
              <div className="flex justify-center items-center rounded-full border-4 border-[#36b47E] h-[75px] w-[75px]">
                  <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-16 w-16"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                  strokeWidth={2}
                  style={{ color: '#36b47E' }}
                  >
                  <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                  </svg>
              </div>
              <h1 className="ml-4 text-3xl">Your new pipeline Temu FRA has successfully been setup</h1>
            </div>
            <div className="flex items-top">
              <p className="mt-5 ml-4">For all future data transfer you can access your pipline via</p>
              <div className="flex flex-col">
                <div className="relative">
                  <Badge
                      variant="secondary"
                      className="text-xl ml-4 my-4 max-w-xs cursor-pointer"
                      onClick={() => copyToClipboard('temu-fra@ecom-pipeline.net')}
                    >
                      temu-fra@ecom-pipeline.net
                  </Badge>
                  {showAlert && (
                    <Alert className="absolute w-[180px] top-0 end-48">
                      <AlertDescription>
                        Email address copied to clipboard.
                      </AlertDescription>
                    </Alert>
                  )}
                </div>
              </div>
            </div>
            <Dialog>
              <DialogTrigger className="mt-8 ml-4">
                <Button variant="outline">View your One Record Data Structure</Button>
              </DialogTrigger>
              <DialogContent className="max-h-[80vh] max-w-[90vw] overflow-y-auto overflow-x-auto">
                <DialogHeader>
                  <DialogTitle>One Record as Json</DialogTitle>
                  <DialogDescription>
                    <ReactJson src={jsonData} theme="monokai" />
                  </DialogDescription>
                </DialogHeader>
              </DialogContent>
            </Dialog>
          </div>
        </div>
      </div>
    </div>
  );
};