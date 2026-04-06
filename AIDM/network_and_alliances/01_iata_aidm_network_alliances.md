# IATA AIDM 25.1 — Network & Alliances
## Mockup Data Model for AI Consumption

---

## Subject Area Overview

The **Network & Alliances** subject area in the IATA AIDM (Aviation Industry Data Model) covers the commercial and operational agreements between airlines — including alliances, codeshare agreements, interline agreements, and prorate arrangements.

- **AIDM Version**: 25.1
- **Domain**: Commercial
- **Subject Area**: Network & Alliances
- **Reference Node**: 6:2:3:4:1:14378
- **Source**: https://airtechzone.iata.org/aidm_model/25.1/index.htm

---

## Entities

### AIRLINE
Represents an air carrier participating in network and alliance arrangements.

| Attribute | Type | Key | Description |
|---|---|---|---|
| airlineCode | VARCHAR(3) | PK | IATA 2-letter or ICAO 3-letter carrier code |
| airlineName | VARCHAR(100) | | Full legal name of the airline |
| IATA_code | VARCHAR(2) | | IATA 2-letter designator |
| ICAO_code | VARCHAR(3) | | ICAO 3-letter designator |
| country | VARCHAR(3) | | ISO 3166-1 alpha-3 country code |
| homeBase | VARCHAR(3) | FK | IATA airport code of primary hub |
| allianceMember | BOOLEAN | | Indicates if airline belongs to a global alliance |

---

### ALLIANCE
Represents a global airline alliance grouping.

| Attribute | Type | Key | Description |
|---|---|---|---|
| allianceCode | VARCHAR(10) | PK | Unique code for the alliance (e.g., OW, ST, SA) |
| allianceName | VARCHAR(100) | | Full name (e.g., oneworld, Star Alliance, SkyTeam) |
| foundedDate | DATE | | Date the alliance was established |
| headquarters | VARCHAR(100) | | City/country of alliance headquarters |
| websiteURL | VARCHAR(255) | | Official website |

---

### ALLIANCE_MEMBERSHIP
Junction entity recording an airline's membership in an alliance.

| Attribute | Type | Key | Description |
|---|---|---|---|
| membershipID | VARCHAR(20) | PK | Unique identifier for the membership record |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| allianceCode | VARCHAR(10) | FK | References ALLIANCE.allianceCode |
| membershipType | VARCHAR(20) | | Full Member / Affiliate / Connect Partner |
| joinDate | DATE | | Date the airline joined the alliance |
| status | VARCHAR(20) | | Active / Suspended / Terminated |

---

### CODESHARE_AGREEMENT
Records a bilateral codeshare agreement between an operating and a marketing carrier.

| Attribute | Type | Key | Description |
|---|---|---|---|
| codeshareAgreementID | VARCHAR(30) | PK | Unique identifier for the codeshare agreement |
| operatingCarrierCode | VARCHAR(3) | FK | References AIRLINE — carrier operating the flight |
| marketingCarrierCode | VARCHAR(3) | FK | References AIRLINE — carrier selling the flight |
| agreementType | VARCHAR(20) | | Free Sale / Block Space / Hard Block / Soft Block |
| effectiveDate | DATE | | Date the agreement comes into force |
| discontinuationDate | DATE | | Date the agreement ends (nullable) |
| status | VARCHAR(20) | | Active / Terminated / Pending / Under Review |
| approvalAuthority | VARCHAR(50) | | Regulatory body (e.g., DOT, EC, ACCC) |
| blockedSeats | INTEGER | | Number of seats blocked (for Block Space type) |

---

### CODESHARE_FLIGHT
Represents a specific flight covered under a codeshare agreement.

| Attribute | Type | Key | Description |
|---|---|---|---|
| codeshareFlightID | VARCHAR(30) | PK | Unique identifier for the codeshare flight record |
| codeshareAgreementID | VARCHAR(30) | FK | References CODESHARE_AGREEMENT |
| operatingFlightNumber | VARCHAR(10) | FK | Flight number of the operating carrier |
| marketingFlightNumber | VARCHAR(10) | | Flight number used by the marketing carrier |
| marketingCarrierCode | VARCHAR(3) | FK | References AIRLINE |
| origin | VARCHAR(3) | | IATA airport code — departure |
| destination | VARCHAR(3) | | IATA airport code — arrival |
| effectiveDate | DATE | | Start date for this flight-level codeshare |
| discontinuationDate | DATE | | End date for this flight-level codeshare (nullable) |
| cabinClassMapping | TEXT | | JSON or reference to cabin class equivalence table |

---

### INTERLINE_AGREEMENT
Records a bilateral interline agreement between two carriers for ticketing and/or baggage.

| Attribute | Type | Key | Description |
|---|---|---|---|
| interlineAgreementID | VARCHAR(30) | PK | Unique identifier for the interline agreement |
| airline1Code | VARCHAR(3) | FK | References AIRLINE — first party |
| airline2Code | VARCHAR(3) | FK | References AIRLINE — second party (must differ from airline1) |
| agreementScope | VARCHAR(20) | | Ticketing / Baggage / Both |
| effectiveDate | DATE | | Date the agreement comes into force |
| expiryDate | DATE | | Date the agreement expires (nullable) |
| status | VARCHAR(20) | | Active / Expired / Pending |
| autoTicketing | CHAR(1) | | Y / N — whether auto-ticketing is enabled |
| selfConnect | CHAR(1) | | Y / N — whether self-connect itineraries are allowed |

---

### PRORATE_AGREEMENT
Defines the revenue-sharing / financial settlement terms between two interline carriers.

| Attribute | Type | Key | Description |
|---|---|---|---|
| prorateAgreementID | VARCHAR(30) | PK | Unique identifier for the prorate agreement |
| interlineAgreementID | VARCHAR(30) | FK | References INTERLINE_AGREEMENT |
| airline1Code | VARCHAR(3) | FK | References AIRLINE — first party |
| airline2Code | VARCHAR(3) | FK | References AIRLINE — second party |
| prorateType | VARCHAR(20) | | MITA / BITA / Special Prorate Agreement |
| currency | CHAR(3) | | ISO 4217 currency code |
| effectiveDate | DATE | | Date the prorate terms come into force |
| expiryDate | DATE | | Date the prorate terms expire (nullable) |
| rulesetReference | VARCHAR(50) | | Reference to governing IATA ruleset or tariff |

---

### SPECIAL_PRORATE_AGREEMENT (SPA)
City-pair-level prorate overrides applied on top of a general prorate agreement.

| Attribute | Type | Key | Description |
|---|---|---|---|
| spaID | VARCHAR(30) | PK | Unique identifier for the SPA record |
| prorateAgreementID | VARCHAR(30) | FK | References PRORATE_AGREEMENT |
| origin | VARCHAR(3) | | IATA city/airport code — origin of city pair |
| destination | VARCHAR(3) | | IATA city/airport code — destination of city pair |
| carrierSequence | VARCHAR(50) | | Ordered list of carriers in the itinerary |
| fareClass | VARCHAR(5) | | Booking class or fare basis code |
| prorateAmount | DECIMAL(10,2) | | Fixed prorate amount (if fixed method) |
| currency | CHAR(3) | | ISO 4217 currency code |
| calculationMethod | VARCHAR(30) | | Fixed / Percentage of Through Fare |
| effectiveDate | DATE | | Start date of SPA applicability |
| expiryDate | DATE | | End date of SPA applicability (nullable) |

---

### ROUTE_NETWORK
Represents an airline's defined route network for a given IATA season.

| Attribute | Type | Key | Description |
|---|---|---|---|
| routeNetworkID | VARCHAR(30) | PK | Unique identifier for the route network record |
| airlineCode | VARCHAR(3) | FK | References AIRLINE |
| routeType | VARCHAR(20) | | Domestic / International / Regional / Long-haul |
| hubAirport | VARCHAR(3) | | IATA code of the primary hub airport |
| version | INTEGER | | Version number of the network plan |
| seasonCode | VARCHAR(5) | | IATA season code (e.g., S25, W25) |

---

### ROUTE
Represents an individual city-pair route within a route network.

| Attribute | Type | Key | Description |
|---|---|---|---|
| routeID | VARCHAR(30) | PK | Unique identifier for the route |
| routeNetworkID | VARCHAR(30) | FK | References ROUTE_NETWORK |
| originAirportCode | VARCHAR(3) | FK | IATA airport code — origin |
| destinationAirportCode | VARCHAR(3) | FK | IATA airport code — destination |
| distance | DECIMAL(8,2) | | Route distance in kilometres |
| distanceType | VARCHAR(10) | | GCD (Great Circle Distance) / Ticketed |
| routeStatus | VARCHAR(20) | | Active / Planned / Suspended / Codeshare-Only |
| onlineStatus | VARCHAR(15) | | Online / Interline |

---

### FREQUENT_FLYER_PARTNER
Records a bilateral frequent flyer programme partnership between two airlines.

| Attribute | Type | Key | Description |
|---|---|---|---|
| ffpPartnerID | VARCHAR(30) | PK | Unique identifier for the FFP partnership |
| sponsorAirlineCode | VARCHAR(3) | FK | References AIRLINE — programme owner |
| partnerAirlineCode | VARCHAR(3) | FK | References AIRLINE — partner carrier |
| programName | VARCHAR(100) | | Name of the FFP programme (e.g., Executive Club) |
| earnEnabled | CHAR(1) | | Y / N — miles/points can be earned on partner flights |
| redeemEnabled | CHAR(1) | | Y / N — miles/points can be redeemed on partner flights |
| statusRecognition | CHAR(1) | | Y / N — elite status is recognised by partner |
| conversionRate | DECIMAL(6,4) | | Miles conversion ratio between programmes |
| effectiveDate | DATE | | Partnership start date |
| expiryDate | DATE | | Partnership end date (nullable) |

---

## Relationships

| From Entity | To Entity | Cardinality | Relationship Description |
|---|---|---|---|
| AIRLINE | ALLIANCE_MEMBERSHIP | 0..* | An airline may have zero or more alliance memberships |
| ALLIANCE | ALLIANCE_MEMBERSHIP | 1..* | An alliance must have at least one member airline |
| AIRLINE | CODESHARE_AGREEMENT (as Operating) | 0..* | An airline may operate flights under multiple codeshare agreements |
| AIRLINE | CODESHARE_AGREEMENT (as Marketing) | 0..* | An airline may market flights under multiple codeshare agreements |
| CODESHARE_AGREEMENT | CODESHARE_FLIGHT | 1..* | An agreement covers one or more specific flights |
| AIRLINE | INTERLINE_AGREEMENT | 0..* | An airline may have multiple bilateral interline agreements |
| INTERLINE_AGREEMENT | PRORATE_AGREEMENT | 0..* | An interline agreement may have attached prorate terms |
| PRORATE_AGREEMENT | SPECIAL_PRORATE_AGREEMENT | 0..* | A prorate agreement may have city-pair-level SPA overrides |
| AIRLINE | ROUTE_NETWORK | 1..* | An airline owns one route network per season |
| ROUTE_NETWORK | ROUTE | 1..* | A network consists of one or more routes |
| AIRLINE | FREQUENT_FLYER_PARTNER | 0..* | An airline may have multiple FFP bilateral partnerships |

---

## Enumerations and Code Lists

### AGREEMENT_TYPE_CD
| Code | Description |
|---|---|
| FREE_SALE | Marketing carrier sells freely without seat allocation |
| BLOCK_SPACE | Operating carrier allocates a fixed block of seats |
| HARD_BLOCK | Fixed block — unused seats returned at no cost |
| SOFT_BLOCK | Fixed block — unused seats returned with penalty |

### INTERLINE_SCOPE_CD
| Code | Description |
|---|---|
| TICKETING | Agreement covers joint ticketing only |
| BAGGAGE | Agreement covers through-checked baggage only |
| BOTH | Agreement covers ticketing and baggage |

### PRORATE_TYPE_CD
| Code | Description |
|---|---|
| MITA | Multilateral Interline Traffic Agreements — standard IATA prorate |
| BITA | Bilateral Interline Traffic Agreements — bilateral override of MITA |
| SPA | Special Prorate Agreement — city-pair-specific prorate terms |

### ALLIANCE_MEMBERSHIP_TYPE_CD
| Code | Description |
|---|---|
| FULL | Full alliance member |
| AFFILIATE | Affiliate member with partial benefits |
| CONNECT | Connect partner with limited agreement scope |

### ROUTE_STATUS_CD
| Code | Description |
|---|---|
| ACTIVE | Route is currently operated |
| PLANNED | Route is approved but not yet operated |
| SUSPENDED | Route is temporarily not operated |
| CODESHARE_ONLY | Route is served only via codeshare arrangement |

### SEASON_CD
| Format | Example | Description |
|---|---|---|
| S + 2-digit year | S25 | IATA Northern Summer season (late March – late October) |
| W + 2-digit year | W25 | IATA Northern Winter season (late October – late March) |

### AGREEMENT_STATUS_CD
| Code | Description |
|---|---|
| ACTIVE | Agreement is currently in force |
| TERMINATED | Agreement has been ended |
| PENDING | Agreement is awaiting regulatory or partner approval |
| UNDER_REVIEW | Agreement is being renegotiated or audited |

---

## Business Rules

| Rule ID | Entity | Rule Description |
|---|---|---|
| BR-NET-001 | CODESHARE_FLIGHT | Must reference a CODESHARE_AGREEMENT where status = 'ACTIVE' |
| BR-NET-002 | INTERLINE_AGREEMENT | airline1Code and airline2Code must be distinct values |
| BR-NET-003 | SPECIAL_PRORATE_AGREEMENT | Can only exist if a parent PRORATE_AGREEMENT record exists |
| BR-NET-004 | ALLIANCE_MEMBERSHIP | An airline may hold only one active membership per alliance at any time |
| BR-NET-005 | All agreement entities | effectiveDate must be earlier than discontinuationDate / expiryDate |
| BR-NET-006 | ROUTE | A route with onlineStatus = 'Online' must have a matching active flight in the operating carrier's schedule |
| BR-NET-007 | CODESHARE_AGREEMENT | operatingCarrierCode and marketingCarrierCode must be distinct values |
| BR-NET-008 | FREQUENT_FLYER_PARTNER | sponsorAirlineCode and partnerAirlineCode must be distinct values |

---

## Traceability to AIDM Domains

| AIDM Domain | AIDM Sub-domain | Entities |
|---|---|---|
| Commercial | Network Planning | ROUTE, ROUTE_NETWORK |
| Commercial | Partnership Management | ALLIANCE, ALLIANCE_MEMBERSHIP, CODESHARE_AGREEMENT, CODESHARE_FLIGHT |
| Commercial | Revenue Accounting | PRORATE_AGREEMENT, SPECIAL_PRORATE_AGREEMENT, INTERLINE_AGREEMENT |
| Commercial | Loyalty | FREQUENT_FLYER_PARTNER |
| Operations | Flight Operations | CODESHARE_FLIGHT (cross-reference to Flight domain) |

---

## Cross-Domain References

| Entity (Network & Alliances) | Referenced Entity | Referenced Domain |
|---|---|---|
| ROUTE.originAirportCode | AIRPORT.airportCode | Airport & Infrastructure |
| ROUTE.destinationAirportCode | AIRPORT.airportCode | Airport & Infrastructure |
| CODESHARE_FLIGHT.operatingFlightNumber | FLIGHT_DESIGNATOR.flightNumber | Flight Operations |
| ROUTE_NETWORK.hubAirport | AIRPORT.airportCode | Airport & Infrastructure |
| AIRLINE.homeBase | AIRPORT.airportCode | Airport & Infrastructure |

---

## Notes

- This model is based on IATA AIDM version 25.1, node reference `6:2:3:4:1:14378`.
- The source portal (https://airtechzone.iata.org) requires authenticated IATA member access.
- Prorate rules follow IATA Resolution 850 (MITA) and bilateral SPA frameworks.
- Codeshare regulatory approval references apply to filings with DOT (USA), EC (Europe), ACCC (Australia), and equivalent authorities.
- IATA season codes follow the standard IATA Scheduling Guidelines.
- All airport and city codes follow IATA SSIM (Standard Schedules Information Manual) conventions.
