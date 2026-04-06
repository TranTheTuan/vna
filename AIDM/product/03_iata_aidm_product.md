# IATA AIDM 25.1 — Domain 03: Product
## Mockup Data Model for AI Consumption

---

## Subject Area Overview

The **Product** domain in the IATA AIDM defines the structured catalogue of airline products offered to customers — from cabin classes and seat configurations to ancillary services, fare families, and bundled product offers. It sits between the commercial (Brand & Marketing, Network & Alliances) and transactional (Sales, Revenue Management) domains.

- **AIDM Version**: 25.1
- **Domain**: Commercial — Product
- **Integrates With**: Network & Alliances (Domain 01), Brand & Marketing (Domain 02)

---

## Cross-Domain FK Dependencies

| Referenced Entity | PK Used | Referenced By (Product) |
|---|---|---|
| AIRLINE | airlineCode | PRODUCT_CATALOGUE, FARE_FAMILY, ANCILLARY_SERVICE, SEAT_MAP |
| ROUTE | routeID | PRODUCT_ROUTE_APPLICABILITY |
| ROUTE_NETWORK | routeNetworkID | PRODUCT_ROUTE_APPLICABILITY |
| BRAND | brandID | PRODUCT_CATALOGUE |
| OFFER | offerID | OFFER_PRODUCT_BUNDLE |

---

## Entities

### PRODUCT_CATALOGUE
The master catalogue of all products offered by an airline. Acts as the root entity for all product types.

| Attribute | Type | Key | Description |
|---|---|---|---|
| productCatalogueID | VARCHAR(30) | PK | Unique identifier for the product catalogue record |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| brandID | VARCHAR(30) | FK | References BRAND.brandID (nullable) |
| productCode | VARCHAR(20) | | Internal product code |
| productName | VARCHAR(100) | | Commercial product name |
| productCategory | VARCHAR(30) | | Flight / Ancillary / Bundle / Lounge / Insurance |
| productType | VARCHAR(30) | | Fare / Seat / Baggage / Meal / Priority / WiFi / Other |
| productDescription | TEXT | | Full product description for display |
| status | VARCHAR(20) | | Active / Inactive / Draft / Retired |
| effectiveDate | DATE | | Date product becomes available |
| expiryDate | DATE | | Date product is retired (nullable) |

---

### FARE_FAMILY
Defines a branded fare family (e.g., Basic, Standard, Flex, Business) within an airline's product strategy.

| Attribute | Type | Key | Description |
|---|---|---|---|
| fareFamilyID | VARCHAR(30) | PK | Unique identifier for the fare family |
| productCatalogueID | VARCHAR(30) | FK | References PRODUCT_CATALOGUE.productCatalogueID |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| fareFamilyCode | VARCHAR(10) | | Short code (e.g., BASIC, FLEX, BIZ) |
| fareFamilyName | VARCHAR(50) | | Display name (e.g., Basic Economy, Flexible) |
| cabinClass | VARCHAR(10) | | Economy / PremiumEconomy / Business / First |
| refundable | CHAR(1) | | Y / N — whether fare is refundable |
| changeable | CHAR(1) | | Y / N — whether itinerary changes are permitted |
| seatSelectionIncluded | CHAR(1) | | Y / N — whether seat selection is included |
| baggageAllowanceKg | INTEGER | | Included checked baggage allowance in kg (0 if none) |
| priorityBoardingIncluded | CHAR(1) | | Y / N — whether priority boarding is included |
| milesAccrualRate | DECIMAL(5,2) | | Miles earning multiplier (e.g., 0.5, 1.0, 1.5) |
| effectiveDate | DATE | | Fare family effective date |
| expiryDate | DATE | | Fare family expiry date (nullable) |

---

### FARE_FAMILY_RULE
Defines the specific conditions, restrictions, and penalties attached to a fare family.

| Attribute | Type | Key | Description |
|---|---|---|---|
| fareRuleID | VARCHAR(30) | PK | Unique identifier for the fare rule |
| fareFamilyID | VARCHAR(30) | FK | References FARE_FAMILY.fareFamilyID |
| ruleCategory | VARCHAR(30) | | Cancellation / Change / Refund / NoShow / Upgrade |
| ruleDescription | VARCHAR(255) | | Human-readable description of the rule |
| penaltyType | VARCHAR(20) | | Fixed / Percentage / NotPermitted / Free |
| penaltyAmount | DECIMAL(10,2) | | Penalty amount (if Fixed type) |
| penaltyPercentage | DECIMAL(5,2) | | Penalty percentage of base fare (if Percentage type) |
| currency | CHAR(3) | | ISO 4217 currency code |
| minHoursBeforeDeparture | INTEGER | | Minimum hours before departure rule applies (nullable) |
| maxHoursBeforeDeparture | INTEGER | | Maximum hours before departure rule applies (nullable) |
| effectiveDate | DATE | | Rule effective date |
| expiryDate | DATE | | Rule expiry date (nullable) |

---

### CABIN_CLASS
Defines a cabin class offered on a flight, including service level and seat configuration reference.

| Attribute | Type | Key | Description |
|---|---|---|---|
| cabinClassID | VARCHAR(30) | PK | Unique identifier for the cabin class record |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| cabinCode | VARCHAR(5) | | IATA cabin code (Y=Economy, W=PremiumEconomy, C=Business, F=First) |
| cabinName | VARCHAR(50) | | Full cabin name (e.g., World Business Class) |
| cabinDescription | TEXT | | Description of the cabin experience |
| pitchInches | INTEGER | | Seat pitch in inches |
| widthInches | DECIMAL(4,1) | | Seat width in inches |
| lieFlat | CHAR(1) | | Y / N — whether seats are lie-flat |
| personalIFE | CHAR(1) | | Y / N — whether personal in-flight entertainment is provided |
| powerOutlet | CHAR(1) | | Y / N — whether power outlets are available |
| wifiAvailable | CHAR(1) | | Y / N — whether WiFi is available |
| status | VARCHAR(20) | | Active / Inactive |

---

### SEAT_MAP
Defines the seat map configuration for an aircraft type used by an airline.

| Attribute | Type | Key | Description |
|---|---|---|---|
| seatMapID | VARCHAR(30) | PK | Unique identifier for the seat map |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| aircraftTypeCode | VARCHAR(10) | | IATA aircraft type code (e.g., 789, 77W, 320) |
| configurationCode | VARCHAR(10) | | Airline-specific configuration code |
| totalSeats | INTEGER | | Total number of seats on this configuration |
| economySeats | INTEGER | | Number of Economy seats |
| premiumEconomySeats | INTEGER | | Number of Premium Economy seats (0 if not offered) |
| businessSeats | INTEGER | | Number of Business class seats (0 if not offered) |
| firstClassSeats | INTEGER | | Number of First class seats (0 if not offered) |
| deckConfiguration | VARCHAR(10) | | Single / Double |
| effectiveDate | DATE | | Date this seat map configuration is effective |
| expiryDate | DATE | | Date this seat map configuration expires (nullable) |

---

### SEAT
Defines an individual seat within a seat map.

| Attribute | Type | Key | Description |
|---|---|---|---|
| seatID | VARCHAR(30) | PK | Unique identifier for the seat record |
| seatMapID | VARCHAR(30) | FK | References SEAT_MAP.seatMapID |
| cabinClassID | VARCHAR(30) | FK | References CABIN_CLASS.cabinClassID |
| seatNumber | VARCHAR(5) | | Seat designator (e.g., 12A, 34F) |
| seatRow | INTEGER | | Row number |
| seatColumn | CHAR(1) | | Column letter (A–K) |
| seatType | VARCHAR(20) | | Window / Middle / Aisle / BulkHead / ExitRow / Solo |
| extraLegroom | CHAR(1) | | Y / N — whether seat has extra legroom |
| reclineRestricted | CHAR(1) | | Y / N — whether seat recline is restricted |
| bassinet | CHAR(1) | | Y / N — whether seat is bassinet-compatible |
| chargeable | CHAR(1) | | Y / N — whether seat selection incurs a fee |
| status | VARCHAR(20) | | Available / BlockedCrew / BlockedMaintenance / Decommissioned |

---

### ANCILLARY_SERVICE
Defines ancillary services available for purchase (e.g., extra baggage, meals, lounge access, WiFi).

| Attribute | Type | Key | Description |
|---|---|---|---|
| ancillaryServiceID | VARCHAR(30) | PK | Unique identifier for the ancillary service |
| productCatalogueID | VARCHAR(30) | FK | References PRODUCT_CATALOGUE.productCatalogueID |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| serviceCode | VARCHAR(10) | | IATA SSIM / PADIS standard service code (e.g., PETC, SPEQ) |
| serviceName | VARCHAR(100) | | Display name of the service |
| serviceCategory | VARCHAR(30) | | Baggage / Meal / Seat / Lounge / WiFi / Insurance / Priority |
| deliveryChannel | VARCHAR(30) | | Booking / CheckIn / OnBoard / Airport |
| unitType | VARCHAR(20) | | PerPax / PerBag / PerFlight / PerKg / PerItem |
| maxUnitsPerPax | INTEGER | | Maximum units a single passenger can purchase |
| chargeType | VARCHAR(20) | | OneWay / Return / PerSector |
| basePrice | DECIMAL(10,2) | | Base price before taxes |
| currency | CHAR(3) | | ISO 4217 currency code |
| taxable | CHAR(1) | | Y / N — whether service is subject to tax |
| status | VARCHAR(20) | | Active / Inactive / Seasonal |
| effectiveDate | DATE | | Date service becomes available |
| expiryDate | DATE | | Date service is retired (nullable) |

---

### PRODUCT_ROUTE_APPLICABILITY
Resolves which products (fare families, ancillaries) are available on which routes or route networks.

| Attribute | Type | Key | Description |
|---|---|---|---|
| productRouteID | VARCHAR(30) | PK | Unique identifier for the product-route applicability record |
| productCatalogueID | VARCHAR(30) | FK | References PRODUCT_CATALOGUE.productCatalogueID |
| routeNetworkID | VARCHAR(30) | FK | References ROUTE_NETWORK.routeNetworkID |
| routeID | VARCHAR(30) | FK | References ROUTE.routeID (nullable — for route-specific applicability) |
| applicabilityType | VARCHAR(20) | | Mandatory / Optional / UpgradeOnly |
| effectiveDate | DATE | | Date applicability comes into force |
| expiryDate | DATE | | Date applicability expires (nullable) |

---

### OFFER_PRODUCT_BUNDLE
Junction entity linking Brand & Marketing OFFER records to specific Product Catalogue items, enabling bundle construction.

| Attribute | Type | Key | Description |
|---|---|---|---|
| offerProductID | VARCHAR(30) | PK | Unique identifier for the offer-product link |
| offerID | VARCHAR(30) | FK | References OFFER.offerID |
| productCatalogueID | VARCHAR(30) | FK | References PRODUCT_CATALOGUE.productCatalogueID |
| quantity | INTEGER | | Number of units of this product included in the bundle |
| sequenceNumber | INTEGER | | Ordering of product within the bundle display |
| overridePrice | DECIMAL(10,2) | | Bundle-specific price override (nullable — uses base price if null) |
| currency | CHAR(3) | | ISO 4217 currency code |

---

## Relationships

### Internal — Product Domain

| From Entity | To Entity | Cardinality | Description |
|---|---|---|---|
| PRODUCT_CATALOGUE | FARE_FAMILY | 0..* | A catalogue entry may define zero or more fare families |
| PRODUCT_CATALOGUE | ANCILLARY_SERVICE | 0..* | A catalogue entry may define zero or more ancillary services |
| PRODUCT_CATALOGUE | PRODUCT_ROUTE_APPLICABILITY | 0..* | A product may be applicable on zero or more routes |
| FARE_FAMILY | FARE_FAMILY_RULE | 1..* | A fare family must have one or more rules |
| SEAT_MAP | SEAT | 1..* | A seat map contains one or more seats |
| CABIN_CLASS | SEAT | 1..* | A cabin class contains one or more seats |

### Cross-Domain — Product → Network & Alliances / Brand & Marketing

| From Entity (Product) | To Entity | Cardinality | Description |
|---|---|---|---|
| PRODUCT_CATALOGUE | AIRLINE | 0..* → 1 | Product belongs to one airline |
| PRODUCT_CATALOGUE | BRAND | 0..* → 0..1 | Product may be associated with a brand identity |
| FARE_FAMILY | AIRLINE | 0..* → 1 | Fare family belongs to one airline |
| SEAT_MAP | AIRLINE | 0..* → 1 | Seat map is owned by one airline |
| CABIN_CLASS | AIRLINE | 0..* → 1 | Cabin class is defined by one airline |
| ANCILLARY_SERVICE | AIRLINE | 0..* → 1 | Ancillary service belongs to one airline |
| PRODUCT_ROUTE_APPLICABILITY | ROUTE_NETWORK | 0..* → 1 | Applicability is scoped to a route network |
| PRODUCT_ROUTE_APPLICABILITY | ROUTE | 0..* → 0..1 | Applicability may narrow to a specific route |
| OFFER_PRODUCT_BUNDLE | OFFER | 0..* → 1 | Bundle links back to marketing offer |

---

## Enumerations and Code Lists

### PRODUCT_CATEGORY_CD
| Code | Description |
|---|---|
| FLIGHT | Core flight product |
| ANCILLARY | Ancillary add-on service |
| BUNDLE | Bundled product package |
| LOUNGE | Airport lounge product |
| INSURANCE | Travel insurance product |

### CABIN_CODE_CD (IATA Standard)
| Code | Description |
|---|---|
| F | First Class |
| C | Business Class |
| W | Premium Economy |
| Y | Economy |

### SEAT_TYPE_CD
| Code | Description |
|---|---|
| WINDOW | Window seat |
| MIDDLE | Middle seat |
| AISLE | Aisle seat |
| BULKHEAD | Bulkhead row seat |
| EXIT_ROW | Emergency exit row seat |
| SOLO | Single solo seat (suite/pod) |

### ANCILLARY_CATEGORY_CD
| Code | Description |
|---|---|
| BAGGAGE | Checked / excess baggage |
| MEAL | Pre-ordered meal / special meal |
| SEAT | Seat selection / upgrade |
| LOUNGE | Lounge day pass |
| WIFI | In-flight WiFi |
| INSURANCE | Travel insurance |
| PRIORITY | Priority boarding / check-in / security |

### PENALTY_TYPE_CD
| Code | Description |
|---|---|
| FIXED | Fixed monetary penalty |
| PERCENTAGE | Percentage of base fare |
| NOT_PERMITTED | Action is not permitted |
| FREE | No penalty applies |

---

## Business Rules

| Rule ID | Entity | Rule Description |
|---|---|---|
| BR-PRD-001 | FARE_FAMILY | cabinClass must match a valid CABIN_CODE_CD value |
| BR-PRD-002 | FARE_FAMILY_RULE | penaltyAmount must be NULL when penaltyType = 'PERCENTAGE'; penaltyPercentage must be NULL when penaltyType = 'FIXED' |
| BR-PRD-003 | SEAT | seatMapID must reference a SEAT_MAP belonging to the same airlineCode as the SEAT record's parent CABIN_CLASS |
| BR-PRD-004 | PRODUCT_ROUTE_APPLICABILITY | When routeID is populated, it must belong to the referenced routeNetworkID |
| BR-PRD-005 | OFFER_PRODUCT_BUNDLE | overridePrice, if populated, must be greater than zero |
| BR-PRD-006 | SEAT_MAP | economySeats + premiumEconomySeats + businessSeats + firstClassSeats must equal totalSeats |
| BR-PRD-007 | ANCILLARY_SERVICE | basePrice must be >= 0 (zero is valid for complimentary services) |
| BR-PRD-008 | PRODUCT_CATALOGUE | effectiveDate must precede expiryDate when expiryDate is not null |

---

## Traceability to AIDM Domains

| AIDM Domain | AIDM Sub-domain | Entities |
|---|---|---|
| Commercial | Product Catalogue | PRODUCT_CATALOGUE, FARE_FAMILY, FARE_FAMILY_RULE |
| Commercial | Cabin & Seat | CABIN_CLASS, SEAT_MAP, SEAT |
| Commercial | Ancillary Services | ANCILLARY_SERVICE |
| Commercial | Product Distribution | PRODUCT_ROUTE_APPLICABILITY, OFFER_PRODUCT_BUNDLE |

---

## Cross-Domain References Summary

| Entity (Product) | Referenced Entity | Referenced Domain |
|---|---|---|
| PRODUCT_CATALOGUE.airlineCode | AIRLINE.airlineCode | Network & Alliances (01) |
| PRODUCT_CATALOGUE.brandID | BRAND.brandID | Brand & Marketing (02) |
| FARE_FAMILY.airlineCode | AIRLINE.airlineCode | Network & Alliances (01) |
| CABIN_CLASS.airlineCode | AIRLINE.airlineCode | Network & Alliances (01) |
| SEAT_MAP.airlineCode | AIRLINE.airlineCode | Network & Alliances (01) |
| ANCILLARY_SERVICE.airlineCode | AIRLINE.airlineCode | Network & Alliances (01) |
| PRODUCT_ROUTE_APPLICABILITY.routeNetworkID | ROUTE_NETWORK.routeNetworkID | Network & Alliances (01) |
| PRODUCT_ROUTE_APPLICABILITY.routeID | ROUTE.routeID | Network & Alliances (01) |
| OFFER_PRODUCT_BUNDLE.offerID | OFFER.offerID | Brand & Marketing (02) |

---

## Notes

- Ancillary service codes align with IATA PADIS (Passenger and Airport Data Interchange Standards) and ATPCO (Airline Tariff Publishing Company) Optional Services standards.
- Cabin codes follow IATA Resolution 728 cabin class definitions.
- Seat map structure is compatible with IATA SSIM Chapter 9 (Aircraft Seat Map) and OTA (OpenTravel Alliance) seat map schema.
- Fare family structures align with IATA NDC (New Distribution Capability) Offer & Order FareComponent model.
