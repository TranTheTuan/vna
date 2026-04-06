# IATA AIDM 25.1 — Domain 09: Cargo
## Mockup Data Model for AI Consumption

---

## Subject Area Overview

The **Cargo** domain in the IATA AIDM covers the end-to-end lifecycle of air freight operations — including shipment booking, airwaybill management, cargo product catalogue, unit load device (ULD) management, cargo acceptance and build-up, dangerous goods compliance, freight rating and pricing, cargo tracking, and revenue accounting for freight. It supports both belly-hold cargo on passenger aircraft and dedicated freighter operations.

- **AIDM Version**: 25.1
- **Domain**: Operations & Commercial — Cargo
- **Integrates With**: Network & Alliances (01), Product (03), Sales (04), Revenue Management & Pricing (06), Ground Operations (07), Flight Operations (08)

---

## Cross-Domain FK Dependencies

| Referenced Entity | PK Used | Referenced By (Cargo) |
|---|---|---|
| AIRLINE | airlineCode | CARGO_SHIPMENT, AIRWAYBILL, CARGO_FLIGHT_BOOKING, CARGO_RATE |
| ROUTE | routeID | CARGO_RATE, CARGO_FLIGHT_BOOKING |
| FLIGHT | flightID | CARGO_FLIGHT_BOOKING, ULD_FLIGHT_ASSIGNMENT |
| AIRCRAFT | aircraftID | CARGO_CAPACITY |
| ROUTE_NETWORK | routeNetworkID | CARGO_RATE |
| REVENUE_RECORD | revenueRecordID | CARGO_REVENUE_RECORD |
| GROUND_HANDLER | groundHandlerID | CARGO_ACCEPTANCE |

---

## Entities

### CARGO_PRODUCT
Defines cargo service products offered by an airline (e.g., express, perishable, general freight).

| Attribute | Type | Key | Description |
|---|---|---|---|
| cargoProductID | VARCHAR(30) | PK | Unique identifier for the cargo product |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| productCode | VARCHAR(10) | | Internal product code |
| productName | VARCHAR(100) | | Commercial product name (e.g., Priority Freight, Cool Chain) |
| productCategory | VARCHAR(30) | | General / Express / Perishable / Pharmaceutical / LiveAnimals / Valuable |
| handlingCode | VARCHAR(5) | | IATA special handling code (e.g., PER, VAL, AVI, ICE) |
| temperatureControlled | CHAR(1) | | Y / N — whether product requires temperature control |
| minTempCelsius | DECIMAL(5,2) | | Minimum temperature in Celsius (nullable) |
| maxTempCelsius | DECIMAL(5,2) | | Maximum temperature in Celsius (nullable) |
| dgr | CHAR(1) | | Y / N — whether product may contain dangerous goods |
| transitTimeHours | INTEGER | | Standard transit time commitment in hours |
| status | VARCHAR(20) | | Active / Inactive / Seasonal |
| effectiveDate | DATE | | Product effective date |
| expiryDate | DATE | | Product expiry date (nullable) |

---

### CARGO_RATE
Defines freight rate tariffs for a cargo product on a specific route or network.

| Attribute | Type | Key | Description |
|---|---|---|---|
| cargoRateID | VARCHAR(30) | PK | Unique identifier for the cargo rate |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| cargoProductID | VARCHAR(30) | FK | References CARGO_PRODUCT.cargoProductID |
| routeID | VARCHAR(30) | FK | References ROUTE.routeID (nullable for network-level rates) |
| routeNetworkID | VARCHAR(30) | FK | References ROUTE_NETWORK.routeNetworkID |
| rateType | VARCHAR(20) | | General / Specific / ClassRate / ULD / Pivot |
| weightBreakpoint | DECIMAL(8,2) | | Weight breakpoint in kg (e.g., 45, 100, 300, 500) |
| ratePerKg | DECIMAL(10,4) | | Rate per kilogram |
| minimumCharge | DECIMAL(10,2) | | Minimum charge for this rate |
| currency | CHAR(3) | | ISO 4217 currency code |
| iataRateClass | VARCHAR(5) | | IATA rate class code (e.g., N, B, Q, E, U, C) |
| effectiveDate | DATE | | Rate effective date |
| expiryDate | DATE | | Rate expiry date (nullable) |
| status | VARCHAR(20) | | Active / Inactive / Filed |

---

### CARGO_SHIPPER
Master record for a cargo shipper (exporter/consignor).

| Attribute | Type | Key | Description |
|---|---|---|---|
| shipperID | VARCHAR(30) | PK | Unique identifier for the shipper |
| shipperName | VARCHAR(100) | | Legal name of shipper |
| shipperType | VARCHAR(20) | | Forwarder / DirectShipper / Consolidator / Broker |
| iataAgentCode | VARCHAR(7) | | IATA cargo agent code (nullable for direct shippers) |
| addressLine1 | VARCHAR(100) | | Shipper address line 1 |
| city | VARCHAR(50) | | Shipper city |
| country | VARCHAR(3) | | ISO 3166-1 alpha-3 country code |
| contactEmail | VARCHAR(255) | | Primary contact email |
| contactPhone | VARCHAR(20) | | Primary contact phone (E.164 format) |
| accountStatus | VARCHAR(20) | | Active / Suspended / Blacklisted |
| knownShipperStatus | CHAR(1) | | Y / N — IATA known shipper status |
| knownShipperVerifiedDate | DATE | | Date known shipper status was verified (nullable) |

---

### CARGO_SHIPMENT
The core entity representing a single freight shipment from origin to destination.

| Attribute | Type | Key | Description |
|---|---|---|---|
| shipmentID | VARCHAR(30) | PK | Unique identifier for the cargo shipment |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode — issuing carrier |
| shipperID | VARCHAR(30) | FK | References CARGO_SHIPPER.shipperID |
| cargoProductID | VARCHAR(30) | FK | References CARGO_PRODUCT.cargoProductID |
| originAirportCode | VARCHAR(3) | | IATA origin airport code |
| destinationAirportCode | VARCHAR(3) | | IATA destination airport code |
| shipmentStatus | VARCHAR(20) | | Booked / Accepted / BuildUp / Loaded / InTransit / Delivered / OnHold |
| commodityCode | VARCHAR(10) | | IATA commodity classification code |
| commodityDescription | VARCHAR(255) | | Description of goods |
| totalPieces | INTEGER | | Total number of pieces in shipment |
| totalWeightKg | DECIMAL(10,3) | | Actual gross weight in kg |
| chargeableWeightKg | DECIMAL(10,3) | | Chargeable weight (greater of actual/volumetric) |
| volumeCbm | DECIMAL(8,4) | | Total volume in cubic metres |
| shipmentValue | DECIMAL(14,2) | | Declared value of goods |
| shipmentCurrency | CHAR(3) | | ISO 4217 currency code for declared value |
| specialHandlingCodes | VARCHAR(50) | | Pipe-separated IATA SHC codes (e.g., PER|RCL|VAL) |
| bookingDateTime | TIMESTAMP | | Date and time shipment was booked |
| requiredDeliveryDate | DATE | | Customer-required delivery date (nullable) |

---

### AIRWAYBILL
Represents a Master Air Waybill (MAWB) or House Air Waybill (HAWB) issued for a shipment.

| Attribute | Type | Key | Description |
|---|---|---|---|
| airwaybillID | VARCHAR(30) | PK | Unique identifier for the airwaybill record |
| shipmentID | VARCHAR(30) | FK | References CARGO_SHIPMENT.shipmentID |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode — issuing carrier |
| awbNumber | VARCHAR(12) | | IATA 11-digit AWB number (e.g., 125-12345678) |
| awbType | VARCHAR(5) | | MAWB / HAWB |
| parentAWBID | VARCHAR(30) | FK | References AIRWAYBILL.airwaybillID — parent MAWB for HAWB (nullable) |
| shipperName | VARCHAR(100) | | Shipper name as printed on AWB |
| consigneeName | VARCHAR(100) | | Consignee name as printed on AWB |
| originAirportCode | VARCHAR(3) | | IATA origin airport |
| destinationAirportCode | VARCHAR(3) | | IATA destination airport |
| pieces | INTEGER | | Number of pieces on this AWB |
| grossWeightKg | DECIMAL(10,3) | | Gross weight in kg |
| chargeableWeightKg | DECIMAL(10,3) | | Chargeable weight in kg |
| ratePerKg | DECIMAL(10,4) | | Applied rate per kg |
| totalCharges | DECIMAL(12,2) | | Total freight charges |
| currency | CHAR(3) | | ISO 4217 currency code |
| chargesPrepaid | CHAR(1) | | P=Prepaid / C=Collect |
| issuanceDateTime | TIMESTAMP | | Date and time AWB was issued |
| awbStatus | VARCHAR(20) | | Draft / Issued / Manifested / Delivered / Voided |

---

### CARGO_FLIGHT_BOOKING
Records the booking of a shipment onto a specific flight segment.

| Attribute | Type | Key | Description |
|---|---|---|---|
| cargoFlightBookingID | VARCHAR(30) | PK | Unique identifier for the cargo flight booking |
| shipmentID | VARCHAR(30) | FK | References CARGO_SHIPMENT.shipmentID |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| flightID | VARCHAR(30) | FK | References FLIGHT.flightID |
| routeID | VARCHAR(30) | FK | References ROUTE.routeID |
| bookedWeightKg | DECIMAL(10,3) | | Weight booked on this flight |
| bookedVolumeCbm | DECIMAL(8,4) | | Volume booked on this flight |
| bookedPieces | INTEGER | | Pieces booked on this flight |
| bookingStatus | VARCHAR(20) | | Confirmed / Waitlisted / Offloaded / Carried |
| bookingClass | VARCHAR(5) | | Cargo booking class (e.g., Q, C, G) |
| sequenceNumber | INTEGER | | Leg sequence in multi-leg routing |

---

### ULD
Master record for a Unit Load Device (container or pallet) used in cargo operations.

| Attribute | Type | Key | Description |
|---|---|---|---|
| uldID | VARCHAR(30) | PK | Unique identifier for the ULD |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode — owning airline |
| uldSerialNumber | VARCHAR(10) | | IATA ULD serial number (e.g., AKE12345BA) |
| uldType | VARCHAR(5) | | IATA ULD type code (e.g., AKE, PMC, PAG, AAP) |
| uldCategory | VARCHAR(10) | | Container / Pallet |
| maxGrossWeightKg | DECIMAL(8,2) | | Maximum gross weight limit in kg |
| tarWeightKg | DECIMAL(6,2) | | Tare weight of empty ULD in kg |
| volumeCbm | DECIMAL(6,3) | | Internal usable volume in cubic metres |
| ownerCode | VARCHAR(3) | | IATA airline owner code |
| uldStatus | VARCHAR(20) | | Available / BuildUp / Loaded / InTransit / Damaged / Condemned |
| currentLocationCode | VARCHAR(3) | | IATA airport code of current location |
| lastInspectionDate | DATE | | Date of last airworthiness inspection |

---

### ULD_FLIGHT_ASSIGNMENT
Records the loading of a ULD onto a specific flight.

| Attribute | Type | Key | Description |
|---|---|---|---|
| uldFlightAssignmentID | VARCHAR(30) | PK | Unique identifier for the ULD-flight assignment |
| uldID | VARCHAR(30) | FK | References ULD.uldID |
| flightID | VARCHAR(30) | FK | References FLIGHT.flightID |
| loadPosition | VARCHAR(10) | | Hold position identifier (e.g., 11L, 22R, FWD-BULK) |
| loadedWeightKg | DECIMAL(10,2) | | Actual weight loaded in this ULD position |
| loadedDateTime | TIMESTAMP | | Timestamp ULD was loaded |
| offloadedDateTime | TIMESTAMP | | Timestamp ULD was offloaded (nullable) |
| assignmentStatus | VARCHAR(20) | | Planned / Loaded / Offloaded / Damaged |

---

### CARGO_ACCEPTANCE
Records the physical acceptance of a cargo shipment at the warehouse/ramp.

| Attribute | Type | Key | Description |
|---|---|---|---|
| acceptanceID | VARCHAR(30) | PK | Unique identifier for the acceptance record |
| shipmentID | VARCHAR(30) | FK | References CARGO_SHIPMENT.shipmentID |
| groundHandlerID | VARCHAR(30) | FK | References GROUND_HANDLER.groundHandlerID |
| acceptanceDateTime | TIMESTAMP | | Date and time shipment was physically accepted |
| acceptanceLocation | VARCHAR(50) | | Warehouse or cargo terminal identifier |
| actualPiecesReceived | INTEGER | | Number of pieces physically received |
| actualWeightKg | DECIMAL(10,3) | | Actual weight measured at acceptance |
| dimensionsVerified | CHAR(1) | | Y / N — whether dimensions were verified |
| screeningMethod | VARCHAR(30) | | XRay / ETD / Physical / CaninUnit / KnownShipper |
| screeningStatus | VARCHAR(20) | | Cleared / HoldForScreening / Rejected |
| screeningDateTime | TIMESTAMP | | Date and time screening was completed |
| acceptanceStatus | VARCHAR(20) | | Accepted / Refused / HeldForQuery |
| refusalReason | VARCHAR(100) | | Reason for refusal (nullable) |

---

### DGR_DECLARATION
Records dangerous goods compliance declaration for a shipment.

| Attribute | Type | Key | Description |
|---|---|---|---|
| dgrDeclarationID | VARCHAR(30) | PK | Unique identifier for the DGR declaration |
| shipmentID | VARCHAR(30) | FK | References CARGO_SHIPMENT.shipmentID |
| unNumber | VARCHAR(6) | | UN number of dangerous goods (e.g., UN1234) |
| properShippingName | VARCHAR(255) | | IATA DGR proper shipping name |
| hazardClass | VARCHAR(5) | | IATA hazard class / division (e.g., 3, 8, 6.1) |
| packingGroup | VARCHAR(5) | | IATA packing group (I / II / III) |
| netQuantityKg | DECIMAL(8,3) | | Net quantity of DG in kg |
| packagingType | VARCHAR(50) | | IATA packaging specification |
| emergencyContact | VARCHAR(20) | | 24-hour emergency contact phone number |
| shipperDeclarationRef | VARCHAR(50) | | Shipper's DGR declaration reference number |
| regulatoryAuthority | VARCHAR(10) | | Applicable regulation (IATA-DGR / ICAO-TI) |
| complianceStatus | VARCHAR(20) | | Compliant / NonCompliant / PendingReview |
| declarationDateTime | TIMESTAMP | | Date and time declaration was submitted |

---

### CARGO_REVENUE_RECORD
Records freight revenue recognised for a carried shipment — feeds into airline revenue accounting.

| Attribute | Type | Key | Description |
|---|---|---|---|
| cargoRevenueID | VARCHAR(30) | PK | Unique identifier for the cargo revenue record |
| airwaybillID | VARCHAR(30) | FK | References AIRWAYBILL.airwaybillID |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| freightCharges | DECIMAL(12,2) | | Base freight charges |
| fuelSurcharge | DECIMAL(10,2) | | Fuel surcharge (YQ equivalent for cargo) |
| securitySurcharge | DECIMAL(10,2) | | Security surcharge |
| otherCharges | DECIMAL(10,2) | | Miscellaneous other charges |
| totalRevenue | DECIMAL(12,2) | | Total revenue recognised |
| currency | CHAR(3) | | ISO 4217 currency code |
| revenueDate | DATE | | Date revenue is recognised (flight date) |
| accountingPeriod | VARCHAR(7) | | Accounting period (YYYY-MM) |
| revenuePerKg | DECIMAL(10,4) | | Yield: revenue per chargeable kg |
| postedDateTime | TIMESTAMP | | Timestamp revenue was posted to accounting |

---

## Relationships

### Internal — Cargo Domain

| From Entity | To Entity | Cardinality | Description |
|---|---|---|---|
| CARGO_PRODUCT | CARGO_RATE | 0..* | A product has zero or more rate entries |
| CARGO_PRODUCT | CARGO_SHIPMENT | 0..* | A product governs zero or more shipments |
| CARGO_SHIPPER | CARGO_SHIPMENT | 0..* | A shipper has zero or more shipments |
| CARGO_SHIPMENT | AIRWAYBILL | 1..* | A shipment has one or more AWBs |
| CARGO_SHIPMENT | CARGO_FLIGHT_BOOKING | 1..* | A shipment has one or more flight bookings |
| CARGO_SHIPMENT | CARGO_ACCEPTANCE | 0..1 | A shipment has zero or one acceptance record |
| CARGO_SHIPMENT | DGR_DECLARATION | 0..* | A shipment may have zero or more DGR declarations |
| AIRWAYBILL | AIRWAYBILL | 0..* | A MAWB may have zero or more HAWB children |
| AIRWAYBILL | CARGO_REVENUE_RECORD | 0..1 | An AWB generates zero or one revenue record |
| ULD | ULD_FLIGHT_ASSIGNMENT | 0..* | A ULD has zero or more flight assignments |

### Cross-Domain — Cargo → Previous Domains

| From Entity | To Entity | Domain | Description |
|---|---|---|---|
| CARGO_PRODUCT | AIRLINE | 01 | Product owned by airline |
| CARGO_RATE | AIRLINE | 01 | Rate filed by airline |
| CARGO_RATE | ROUTE | 01 | Rate applicable to a route |
| CARGO_RATE | ROUTE_NETWORK | 01 | Rate scoped to a route network |
| CARGO_SHIPMENT | AIRLINE | 01 | Shipment carried by airline |
| AIRWAYBILL | AIRLINE | 01 | AWB issued by airline |
| CARGO_FLIGHT_BOOKING | AIRLINE | 01 | Booking with airline |
| CARGO_FLIGHT_BOOKING | FLIGHT | 08 | Booking on a specific flight |
| CARGO_FLIGHT_BOOKING | ROUTE | 01 | Booking for a route |
| ULD | AIRLINE | 01 | ULD owned by airline |
| ULD_FLIGHT_ASSIGNMENT | FLIGHT | 08 | ULD loaded on a flight |
| CARGO_ACCEPTANCE | GROUND_HANDLER | 07 | Acceptance handled by ground handler |
| CARGO_REVENUE_RECORD | AIRLINE | 01 | Revenue attributed to airline |

---

## Enumerations and Code Lists

### CARGO_PRODUCT_CATEGORY_CD
| Code | Description |
|---|---|
| GENERAL | General cargo — standard freight |
| EXPRESS | Time-definite express freight |
| PERISHABLE | Perishable goods (temperature sensitive) |
| PHARMACEUTICAL | Pharmaceutical / GDP-compliant cold chain |
| LIVE_ANIMALS | Live animals (IATA LAR) |
| VALUABLE | High-value cargo (jewellery, currency, art) |

### IATA_SPECIAL_HANDLING_CODE_CD (Partial)
| Code | Description |
|---|---|
| PER | Perishable cargo |
| VAL | Valuable cargo |
| AVI | Live animals |
| ICE | Dry ice / carbon dioxide solid |
| RCL | Radioactive material — Category I |
| EAT | Foodstuff — perishable |
| PIL | Pharmaceutical in-transit |

### AWB_TYPE_CD
| Code | Description |
|---|---|
| MAWB | Master Air Waybill |
| HAWB | House Air Waybill |

### SHIPMENT_STATUS_CD
| Code | Description |
|---|---|
| BOOKED | Shipment booked, not yet physical |
| ACCEPTED | Physically accepted at cargo terminal |
| BUILD_UP | Being loaded into ULD |
| LOADED | Loaded onto aircraft |
| IN_TRANSIT | En route to destination |
| DELIVERED | Delivered to consignee |
| ON_HOLD | Held pending resolution |

### ULD_TYPE_CD (IATA Standard)
| Code | Description |
|---|---|
| AKE | LD3 — Lower deck container |
| PMC | P6P — Main deck pallet |
| PAG | P1P — Lower deck pallet |
| AAP | LD3 — Lower deck container (half-height) |
| PGA | P6P — Main deck pallet (with net) |

### DGR_HAZARD_CLASS_CD (IATA DGR)
| Code | Description |
|---|---|
| 1 | Explosives |
| 2 | Gases |
| 3 | Flammable Liquids |
| 4 | Flammable Solids |
| 5 | Oxidising Substances |
| 6 | Toxic and Infectious Substances |
| 7 | Radioactive Material |
| 8 | Corrosives |
| 9 | Miscellaneous Dangerous Goods |

---

## Business Rules

| Rule ID | Entity | Rule Description |
|---|---|---|
| BR-CGO-001 | AIRWAYBILL | awbNumber must be exactly 11 digits in format NNN-NNNNNNNN (IATA Resolution 600b) |
| BR-CGO-002 | CARGO_SHIPMENT | chargeableWeightKg must equal MAX(totalWeightKg, volumeCbm × 167) — IATA volumetric weight factor |
| BR-CGO-003 | CARGO_ACCEPTANCE | screeningStatus must be 'Cleared' before acceptanceStatus = 'Accepted' |
| BR-CGO-004 | DGR_DECLARATION | complianceStatus must be 'Compliant' before shipment can be loaded onto aircraft |
| BR-CGO-005 | ULD | uldSerialNumber must follow IATA ULD naming convention: 3-letter type + 5-digit serial + 2-letter owner code |
| BR-CGO-006 | ULD_FLIGHT_ASSIGNMENT | loadedWeightKg must not exceed ULD.maxGrossWeightKg |
| BR-CGO-007 | CARGO_REVENUE_RECORD | totalRevenue must equal freightCharges + fuelSurcharge + securitySurcharge + otherCharges |
| BR-CGO-008 | CARGO_SHIPMENT | knownShipperStatus must be Y OR screeningMethod must be XRay/ETD for acceptance |
| BR-CGO-009 | AIRWAYBILL | chargesPrepaid must be either 'P' (Prepaid) or 'C' (Collect) |
| BR-CGO-010 | CARGO_FLIGHT_BOOKING | bookedWeightKg must be greater than zero |

---

## Traceability to AIDM Domains

| AIDM Domain | AIDM Sub-domain | Entities |
|---|---|---|
| Cargo | Product & Rating | CARGO_PRODUCT, CARGO_RATE |
| Cargo | Shipper Management | CARGO_SHIPPER |
| Cargo | Shipment Management | CARGO_SHIPMENT, AIRWAYBILL |
| Cargo | Flight Booking | CARGO_FLIGHT_BOOKING |
| Cargo | ULD Management | ULD, ULD_FLIGHT_ASSIGNMENT |
| Cargo | Acceptance & Security | CARGO_ACCEPTANCE, DGR_DECLARATION |
| Cargo | Revenue Accounting | CARGO_REVENUE_RECORD |

---

## Notes

- AWB numbering follows IATA Resolution 600b (Air Waybill specifications).
- Volumetric weight conversion factor of 1 CBM = 167 kg follows IATA standard cargo practice.
- ULD type codes follow IATA ULD Regulations (ULD Reg) — Section 4 ULD identification.
- Dangerous goods compliance follows IATA Dangerous Goods Regulations (DGR) — 65th Edition.
- Special Handling Codes follow IATA TACT (The Air Cargo Tariff) Rules.
- Cargo screening standards align with ICAO Annex 17 (Security) and EC Regulation 300/2008.
- Known shipper programme follows IATA CEIV (Centre of Excellence for Independent Validators) Pharma/Fresh standards.
- Cargo rate classes follow IATA Resolutions 010e and 502 (cargo tariff construction rules).
