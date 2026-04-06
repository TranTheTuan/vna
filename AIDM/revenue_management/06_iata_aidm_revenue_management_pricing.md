# IATA AIDM 25.1 — Domain 06: Revenue Management & Pricing
## Mockup Data Model for AI Consumption

---

## Subject Area Overview

The **Revenue Management & Pricing** domain in the IATA AIDM covers fare construction, inventory control, demand forecasting, yield management, dynamic pricing, overbooking management, and revenue accounting. It is the analytical and operational engine that maximises airline revenue per available seat kilometre (RASK) by controlling which products are sold at which price points at any given time.

- **AIDM Version**: 25.1
- **Domain**: Commercial — Revenue Management & Pricing
- **Integrates With**: Network & Alliances (01), Brand & Marketing (02), Product (03), Sales (04), Customer & Loyalty (05)

---

## Cross-Domain FK Dependencies

| Referenced Entity | PK Used | Referenced By (RM & Pricing) |
|---|---|---|
| AIRLINE | airlineCode | FARE, INVENTORY_CLASS, FARE_RULE_SET, REVENUE_RECORD |
| ROUTE | routeID | FARE, INVENTORY_CLASS, DEMAND_FORECAST |
| ROUTE_NETWORK | routeNetworkID | REVENUE_RECORD |
| FARE_FAMILY | fareFamilyID | FARE |
| CABIN_CLASS | cabinClassID | INVENTORY_CLASS |
| BOOKING_SEGMENT | bookingSegmentID | INVENTORY_BOOKING_ACTIVITY |
| BOOKING | bookingID | REVENUE_RECORD |
| TICKET | ticketID | REVENUE_RECORD |
| CUSTOMER_PROFILE | customerID | DYNAMIC_PRICE_OFFER |
| LOYALTY_ACCOUNT | loyaltyAccountID | DYNAMIC_PRICE_OFFER |
| SEAT_MAP | seatMapID | INVENTORY_CLASS |

---

## Entities

### FARE
Defines a published fare for a specific origin-destination market, booking class, and carrier.

| Attribute | Type | Key | Description |
|---|---|---|---|
| fareID | VARCHAR(30) | PK | Unique identifier for the fare record |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode — filing carrier |
| routeID | VARCHAR(30) | FK | References ROUTE.routeID |
| fareFamilyID | VARCHAR(30) | FK | References FARE_FAMILY.fareFamilyID |
| fareBasisCode | VARCHAR(20) | | ATPCO fare basis code (e.g., YOWUS, BAPXEUS) |
| bookingClass | VARCHAR(2) | | RBD booking class (e.g., Y, B, M, H, Q, V) |
| cabinCode | VARCHAR(1) | | IATA cabin code (Y/W/C/F) |
| fareAmount | DECIMAL(12,2) | | Published one-way fare amount |
| currency | CHAR(3) | | ISO 4217 currency code |
| fareType | VARCHAR(20) | | Published / Negotiated / Corporate / Group / Net |
| owrt | CHAR(1) | | O=One-Way / R=Round-Trip / X=Either |
| globalIndicator | VARCHAR(2) | | IATA global indicator (AT/PA/PN/SA/AP/EH/TS/PO/FE) |
| routingNumber | VARCHAR(10) | | ATPCO routing number |
| mileageIndicator | CHAR(1) | | M=Mileage / R=Routing |
| effectiveDate | DATE | | Fare effective date |
| discontinuationDate | DATE | | Fare discontinuation date (nullable) |
| tariffNumber | VARCHAR(10) | | ATPCO tariff number |
| ruleNumber | VARCHAR(10) | | ATPCO rule number |
| status | VARCHAR(20) | | Active / Inactive / Filed / Withdrawn |

---

### FARE_RULE_SET
Defines the structured rule conditions attached to a fare (ATPCO Categories 1–35).

| Attribute | Type | Key | Description |
|---|---|---|---|
| fareRuleSetID | VARCHAR(30) | PK | Unique identifier for the fare rule set |
| fareID | VARCHAR(30) | FK | References FARE.fareID |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| atpcoCategory | INTEGER | | ATPCO rule category number (1–35) |
| categoryName | VARCHAR(50) | | Category description (e.g., Eligibility, Stopovers, Combinability) |
| ruleText | TEXT | | Full rule text as filed |
| minStay | VARCHAR(20) | | Minimum stay requirement (e.g., SU=Sunday) |
| maxStay | VARCHAR(20) | | Maximum stay requirement (e.g., 1M=1 Month) |
| advancePurchaseDays | INTEGER | | Days before departure advance purchase required |
| blackoutDates | TEXT | | JSON array of blacked-out date ranges |
| permittedDaysOfWeek | VARCHAR(7) | | Bitmask Mon–Sun permitted travel days |
| effectiveDate | DATE | | Rule set effective date |
| discontinuationDate | DATE | | Rule set discontinuation date (nullable) |

---

### TAX_FEE_CHARGE
Defines taxes, fees, and charges applicable to fares for specific markets.

| Attribute | Type | Key | Description |
|---|---|---|---|
| taxFeeChargeID | VARCHAR(30) | PK | Unique identifier for the tax/fee/charge record |
| fareID | VARCHAR(30) | FK | References FARE.fareID |
| taxCode | VARCHAR(10) | | IATA/ATPCO tax code (e.g., YQ, YR, GB, US, XY) |
| taxType | VARCHAR(20) | | GovernmentTax / AirportFee / CarrierSurcharge / Levy |
| taxName | VARCHAR(100) | | Full descriptive name of the tax/fee |
| calculationMethod | VARCHAR(20) | | Flat / PercentageOfFare / PerSegment / PerCoupon |
| taxAmount | DECIMAL(10,2) | | Fixed tax amount (nullable if percentage-based) |
| taxPercentage | DECIMAL(5,2) | | Percentage rate (nullable if flat amount) |
| currency | CHAR(3) | | ISO 4217 currency code |
| applicableCountry | VARCHAR(3) | | ISO 3166-1 alpha-3 country where tax applies |
| effectiveDate | DATE | | Tax effective date |
| discontinuationDate | DATE | | Tax discontinuation date (nullable) |

---

### INVENTORY_CLASS
Defines a fare class inventory bucket for a specific route and aircraft configuration.

| Attribute | Type | Key | Description |
|---|---|---|---|
| inventoryClassID | VARCHAR(30) | PK | Unique identifier for the inventory class record |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| routeID | VARCHAR(30) | FK | References ROUTE.routeID |
| cabinClassID | VARCHAR(30) | FK | References CABIN_CLASS.cabinClassID |
| seatMapID | VARCHAR(30) | FK | References SEAT_MAP.seatMapID |
| bookingClass | VARCHAR(2) | | RBD booking class code |
| cabinCode | VARCHAR(1) | | IATA cabin code (Y/W/C/F) |
| totalCapacity | INTEGER | | Total seats allocated to this class |
| authorisedCapacity | INTEGER | | Seats released for sale (may exceed physical for overbooking) |
| bookedCount | INTEGER | | Current number of confirmed bookings |
| waitlistCount | INTEGER | | Current number of waitlisted bookings |
| availableCount | INTEGER | | Computed available seats (authorised − booked) |
| overbookingFactor | DECIMAL(5,3) | | Overbooking multiplier (e.g., 1.050 = 5% overbooked) |
| closeOutLevel | INTEGER | | Booking count at which class is closed |
| classStatus | VARCHAR(10) | | Open / Closed / Waitlist / Request |
| protectedCapacity | INTEGER | | Seats protected for upgrades or operations |
| lastUpdatedDateTime | TIMESTAMP | | Last inventory update timestamp |

---

### INVENTORY_BOOKING_ACTIVITY
Records each booking action against an inventory class (sale, cancellation, upgrade, etc.).

| Attribute | Type | Key | Description |
|---|---|---|---|
| inventoryActivityID | VARCHAR(30) | PK | Unique identifier for the inventory activity record |
| inventoryClassID | VARCHAR(30) | FK | References INVENTORY_CLASS.inventoryClassID |
| bookingSegmentID | VARCHAR(30) | FK | References BOOKING_SEGMENT.bookingSegmentID |
| activityType | VARCHAR(20) | | Sale / Cancellation / Upgrade / Downgrade / NoShow / Reinstate |
| seatsAffected | INTEGER | | Number of seats affected by this activity |
| activityDateTime | TIMESTAMP | | Date and time of activity |
| sourceSystem | VARCHAR(30) | | System that triggered the activity (PSS/NDC/DCS) |

---

### DEMAND_FORECAST
Stores demand forecast data for a route and departure date, used by RM optimisation systems.

| Attribute | Type | Key | Description |
|---|---|---|---|
| demandForecastID | VARCHAR(30) | PK | Unique identifier for the demand forecast record |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| routeID | VARCHAR(30) | FK | References ROUTE.routeID |
| forecastDate | DATE | | The departure date being forecast |
| forecastRunDateTime | TIMESTAMP | | When the forecast was generated |
| forecastModel | VARCHAR(50) | | Model identifier (e.g., EMSR-b, PickUp, HybridML) |
| bookingClass | VARCHAR(2) | | RBD booking class being forecast |
| forecastDemand | DECIMAL(8,2) | | Forecast number of passengers demanding this class |
| forecastRevenue | DECIMAL(12,2) | | Forecast revenue for this class on this flight |
| currency | CHAR(3) | | ISO 4217 currency code |
| loadFactorForecast | DECIMAL(5,4) | | Forecast load factor (0.0000–1.0000) |
| confidenceInterval | DECIMAL(5,4) | | Statistical confidence interval |
| forecastHorizonDays | INTEGER | | Days before departure this forecast covers |

---

### PRICING_ACTION
Records a revenue management pricing decision or fare availability change made by the RM system or analyst.

| Attribute | Type | Key | Description |
|---|---|---|---|
| pricingActionID | VARCHAR(30) | PK | Unique identifier for the pricing action |
| inventoryClassID | VARCHAR(30) | FK | References INVENTORY_CLASS.inventoryClassID |
| fareID | VARCHAR(30) | FK | References FARE.fareID |
| actionType | VARCHAR(20) | | Open / Close / Adjust / Protect / Release |
| triggeredBy | VARCHAR(20) | | Automated / Analyst / Rule / Competitor |
| previousStatus | VARCHAR(10) | | Inventory class status before action |
| newStatus | VARCHAR(10) | | Inventory class status after action |
| previousCapacity | INTEGER | | Authorised capacity before action |
| newCapacity | INTEGER | | Authorised capacity after action |
| actionDateTime | TIMESTAMP | | Timestamp of pricing action |
| analystID | VARCHAR(30) | | Agent/analyst identifier (nullable for automated) |
| notes | VARCHAR(255) | | Optional notes on the pricing decision |

---

### DYNAMIC_PRICE_OFFER
Records a real-time personalised price offer generated for a specific customer interaction.

| Attribute | Type | Key | Description |
|---|---|---|---|
| dynamicPriceOfferID | VARCHAR(30) | PK | Unique identifier for the dynamic price offer |
| fareID | VARCHAR(30) | FK | References FARE.fareID — base fare |
| customerID | VARCHAR(30) | FK | References CUSTOMER_PROFILE.customerID (nullable for anonymous) |
| loyaltyAccountID | VARCHAR(30) | FK | References LOYALTY_ACCOUNT.loyaltyAccountID (nullable) |
| sessionID | VARCHAR(50) | | Browser/app session identifier |
| offerChannel | VARCHAR(30) | | Web / Mobile / NDC / GDS / Kiosk |
| basePrice | DECIMAL(12,2) | | Published base fare amount |
| dynamicAdjustment | DECIMAL(10,2) | | Price adjustment applied (positive or negative) |
| finalPrice | DECIMAL(12,2) | | Final presented price |
| currency | CHAR(3) | | ISO 4217 currency code |
| pricingModel | VARCHAR(50) | | Model used (e.g., WTP_ML, RealTimeRM, PersonalisedFare) |
| willingnessToPay | DECIMAL(12,2) | | Predicted willingness-to-pay estimate |
| offerValidUntil | TIMESTAMP | | Offer expiry timestamp |
| conversionStatus | VARCHAR(20) | | Presented / Booked / Abandoned / Expired |
| offerGeneratedDateTime | TIMESTAMP | | Timestamp offer was generated |

---

### OVERBOOKING_POLICY
Defines overbooking parameters per route, season, and cabin class.

| Attribute | Type | Key | Description |
|---|---|---|---|
| overbookingPolicyID | VARCHAR(30) | PK | Unique identifier for the overbooking policy |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| routeID | VARCHAR(30) | FK | References ROUTE.routeID |
| cabinCode | VARCHAR(1) | | IATA cabin code (Y/W/C/F) |
| seasonCode | VARCHAR(5) | | IATA season code (e.g., S25, W25) |
| maxOverbookingFactor | DECIMAL(5,3) | | Maximum permitted overbooking factor |
| targetOverbookingFactor | DECIMAL(5,3) | | Target overbooking factor for optimisation |
| noShowRate | DECIMAL(5,4) | | Historical no-show rate for this market |
| cancellationRate | DECIMAL(5,4) | | Historical cancellation rate for this market |
| voluntaryDenialBudget | INTEGER | | Target number of volunteer denied boardings per month |
| involuntaryDenialThreshold | INTEGER | | Threshold triggering involuntary denied boarding process |
| effectiveDate | DATE | | Policy effective date |
| expiryDate | DATE | | Policy expiry date (nullable) |

---

### REVENUE_RECORD
Stores the final revenue accounting record for a flown booking segment — input to revenue accounting.

| Attribute | Type | Key | Description |
|---|---|---|---|
| revenueRecordID | VARCHAR(30) | PK | Unique identifier for the revenue record |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| bookingID | VARCHAR(30) | FK | References BOOKING.bookingID |
| ticketID | VARCHAR(30) | FK | References TICKET.ticketID |
| routeNetworkID | VARCHAR(30) | FK | References ROUTE_NETWORK.routeNetworkID |
| fareID | VARCHAR(30) | FK | References FARE.fareID |
| revenueType | VARCHAR(20) | | Passenger / Ancillary / Cargo / Other |
| grossRevenue | DECIMAL(12,2) | | Total gross revenue before deductions |
| commissionPaid | DECIMAL(10,2) | | Commission paid to agency/channel |
| taxesCollected | DECIMAL(10,2) | | Total taxes and fees collected |
| netRevenue | DECIMAL(12,2) | | Net revenue after commission and taxes |
| currency | CHAR(3) | | ISO 4217 currency code |
| revenueDate | DATE | | Date revenue is recognised (flight date) |
| raskContribution | DECIMAL(10,6) | | Revenue per Available Seat Kilometre contribution |
| yieldPerMile | DECIMAL(10,6) | | Revenue yield per revenue passenger mile |
| accountingPeriod | VARCHAR(7) | | Accounting period (YYYY-MM) |
| postedDateTime | TIMESTAMP | | Timestamp revenue was posted to accounting |

---

## Relationships

### Internal — Revenue Management & Pricing Domain

| From Entity | To Entity | Cardinality | Description |
|---|---|---|---|
| FARE | FARE_RULE_SET | 1..* | A fare has one or more rule sets |
| FARE | TAX_FEE_CHARGE | 0..* | A fare may have zero or more taxes/fees |
| FARE | DYNAMIC_PRICE_OFFER | 0..* | A fare may generate zero or more dynamic offers |
| FARE | PRICING_ACTION | 0..* | A fare may have zero or more pricing actions |
| FARE | REVENUE_RECORD | 0..* | A fare generates zero or more revenue records |
| INVENTORY_CLASS | INVENTORY_BOOKING_ACTIVITY | 0..* | An inventory class has zero or more booking activities |
| INVENTORY_CLASS | PRICING_ACTION | 0..* | An inventory class has zero or more pricing actions |
| INVENTORY_CLASS | DEMAND_FORECAST | 0..* | An inventory class has zero or more forecasts |

### Cross-Domain — RM & Pricing → Previous Domains

| From Entity | To Entity | Domain | Description |
|---|---|---|---|
| FARE | AIRLINE | 01 | Fare filed by an airline |
| FARE | ROUTE | 01 | Fare applicable to a route |
| FARE | FARE_FAMILY | 03 | Fare belongs to a fare family |
| INVENTORY_CLASS | AIRLINE | 01 | Inventory class owned by airline |
| INVENTORY_CLASS | ROUTE | 01 | Inventory class for a route |
| INVENTORY_CLASS | CABIN_CLASS | 03 | Inventory class maps to a cabin |
| INVENTORY_CLASS | SEAT_MAP | 03 | Inventory class uses a seat map config |
| INVENTORY_BOOKING_ACTIVITY | BOOKING_SEGMENT | 04 | Activity driven by a booking segment |
| DEMAND_FORECAST | AIRLINE | 01 | Forecast for an airline's route |
| DEMAND_FORECAST | ROUTE | 01 | Forecast scoped to a specific route |
| DYNAMIC_PRICE_OFFER | CUSTOMER_PROFILE | 05 | Offer personalised for a customer |
| DYNAMIC_PRICE_OFFER | LOYALTY_ACCOUNT | 05 | Offer influenced by loyalty status |
| OVERBOOKING_POLICY | AIRLINE | 01 | Policy defined by an airline |
| OVERBOOKING_POLICY | ROUTE | 01 | Policy scoped to a route |
| REVENUE_RECORD | AIRLINE | 01 | Revenue attributed to an airline |
| REVENUE_RECORD | BOOKING | 04 | Revenue from a booking |
| REVENUE_RECORD | TICKET | 04 | Revenue evidenced by a ticket |
| REVENUE_RECORD | ROUTE_NETWORK | 01 | Revenue attributed to a network |

---

## Enumerations and Code Lists

### FARE_TYPE_CD
| Code | Description |
|---|---|
| PUBLISHED | Standard ATPCO published fare |
| NEGOTIATED | Private negotiated fare |
| CORPORATE | Corporate account fare |
| GROUP | Group booking fare |
| NET | Net/wholesale fare |

### OWRT_CD
| Code | Description |
|---|---|
| O | One-Way fare |
| R | Round-Trip fare |
| X | Either direction |

### GLOBAL_INDICATOR_CD (IATA Standard)
| Code | Description |
|---|---|
| AT | Atlantic |
| PA | Pacific |
| PN | Pacific/North America |
| SA | South Atlantic |
| AP | Atlantic/Pacific |
| EH | Eastern Hemisphere |
| TS | Trans-Siberian |
| PO | Polar |
| FE | Far East |

### INVENTORY_CLASS_STATUS_CD
| Code | Description |
|---|---|
| OPEN | Class is open for sale |
| CLOSED | Class is closed — no further sales |
| WAITLIST | Class is on waitlist |
| REQUEST | Class requires approval to book |

### PRICING_ACTION_TYPE_CD
| Code | Description |
|---|---|
| OPEN | Open a previously closed booking class |
| CLOSE | Close a booking class to further sales |
| ADJUST | Adjust authorised capacity |
| PROTECT | Protect capacity from lower-class sales |
| RELEASE | Release protected capacity |

### REVENUE_TYPE_CD
| Code | Description |
|---|---|
| PASSENGER | Passenger flight revenue |
| ANCILLARY | Ancillary service revenue |
| CARGO | Cargo belly-hold revenue |
| OTHER | Miscellaneous revenue |

---

## Business Rules

| Rule ID | Entity | Rule Description |
|---|---|---|
| BR-RMP-001 | INVENTORY_CLASS | availableCount must equal authorisedCapacity minus bookedCount |
| BR-RMP-002 | FARE | fareAmount must be greater than zero |
| BR-RMP-003 | FARE_RULE_SET | atpcoCategory must be an integer between 1 and 35 |
| BR-RMP-004 | OVERBOOKING_POLICY | maxOverbookingFactor must be greater than or equal to targetOverbookingFactor |
| BR-RMP-005 | DYNAMIC_PRICE_OFFER | finalPrice must equal basePrice plus dynamicAdjustment |
| BR-RMP-006 | REVENUE_RECORD | netRevenue must equal grossRevenue minus commissionPaid minus taxesCollected |
| BR-RMP-007 | DEMAND_FORECAST | loadFactorForecast must be between 0.0000 and 1.0000 |
| BR-RMP-008 | INVENTORY_CLASS | overbookingFactor must be between 1.000 and 2.000 |
| BR-RMP-009 | FARE | effectiveDate must precede discontinuationDate when discontinuationDate is not null |
| BR-RMP-010 | PRICING_ACTION | actionDateTime must be after the inventory class lastUpdatedDateTime |

---

## Traceability to AIDM Domains

| AIDM Domain | AIDM Sub-domain | Entities |
|---|---|---|
| Commercial | Fare Management | FARE, FARE_RULE_SET, TAX_FEE_CHARGE |
| Commercial | Inventory Management | INVENTORY_CLASS, INVENTORY_BOOKING_ACTIVITY |
| Commercial | Demand Forecasting | DEMAND_FORECAST |
| Commercial | Pricing Optimisation | PRICING_ACTION, DYNAMIC_PRICE_OFFER |
| Commercial | Overbooking | OVERBOOKING_POLICY |
| Commercial | Revenue Accounting | REVENUE_RECORD |

---

## Notes

- Fare structures follow ATPCO (Airline Tariff Publishing Company) Category 1–35 rule framework.
- Global indicators follow IATA Resolution 045e (Fare Construction Principles).
- Tax codes follow IATA ATC (Airlines Tax Code) standard.
- Inventory management follows IATA ADRM RM standards and EMSR-b (Expected Marginal Seat Revenue) optimisation model conventions.
- Dynamic pricing model identifiers align with IATA NDC Offer & Order real-time pricing capability framework.
- Revenue recognition follows IATA Resolution 010a (Revenue Accounting Manual) principles.
- RASK and yield metrics align with IATA Industry Statistics definitions.
