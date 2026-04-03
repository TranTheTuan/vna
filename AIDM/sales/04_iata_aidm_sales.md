# IATA AIDM 25.1 — Domain 04: Sales
## Mockup Data Model for AI Consumption

---

## Subject Area Overview

The **Sales** domain in the IATA AIDM covers the end-to-end commercial transaction lifecycle — from booking and ticketing through payment, distribution channels, travel agency management, and commission structures. It is the primary transactional bridge between the Product domain and the customer-facing revenue stream.

- **AIDM Version**: 25.1
- **Domain**: Commercial — Sales
- **Integrates With**: Network & Alliances (01), Brand & Marketing (02), Product (03)

---

## Cross-Domain FK Dependencies

| Referenced Entity | PK Used | Referenced By (Sales) |
|---|---|---|
| AIRLINE | airlineCode | BOOKING, TICKET, SALES_CHANNEL, TRAVEL_AGENCY |
| ROUTE | routeID | BOOKING_SEGMENT |
| FARE_FAMILY | fareFamilyID | BOOKING_SEGMENT |
| ANCILLARY_SERVICE | ancillaryServiceID | BOOKING_ANCILLARY |
| OFFER | offerID | BOOKING |
| PRODUCT_CATALOGUE | productCatalogueID | BOOKING_ANCILLARY |
| LOYALTY_PROGRAMME | loyaltyProgrammeID | BOOKING |
| FREQUENT_FLYER_PARTNER | ffpPartnerID | TICKET |

---

## Entities

### SALES_CHANNEL
Defines the distribution channels through which an airline's products are sold.

| Attribute | Type | Key | Description |
|---|---|---|---|
| salesChannelID | VARCHAR(30) | PK | Unique identifier for the sales channel |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| channelCode | VARCHAR(10) | | Internal channel code |
| channelName | VARCHAR(100) | | Display name (e.g., Direct Web, GDS, NDC, OTA, Airport) |
| channelType | VARCHAR(30) | | Direct / GDS / NDC / OTA / TMC / Airport / CallCentre |
| distributionType | VARCHAR(20) | | Retail / Wholesale / Corporate / Group |
| commissionable | CHAR(1) | | Y / N — whether this channel earns commission |
| ndc_enabled | CHAR(1) | | Y / N — whether channel supports IATA NDC |
| status | VARCHAR(20) | | Active / Inactive / Pilot |
| effectiveDate | DATE | | Channel activation date |
| expiryDate | DATE | | Channel deactivation date (nullable) |

---

### TRAVEL_AGENCY
Represents a travel agency or booking intermediary authorised to sell airline products.

| Attribute | Type | Key | Description |
|---|---|---|---|
| travelAgencyID | VARCHAR(30) | PK | Unique identifier for the travel agency |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode — contracting airline |
| salesChannelID | VARCHAR(30) | FK | References SALES_CHANNEL.salesChannelID |
| iataAgencyCode | VARCHAR(8) | | IATA 8-digit numeric agency code |
| agencyName | VARCHAR(100) | | Full legal agency name |
| agencyType | VARCHAR(20) | | GDS / OTA / TMC / Consolidator / Retail |
| country | VARCHAR(3) | | ISO 3166-1 alpha-3 country code |
| city | VARCHAR(50) | | City of agency location |
| gdsCode | VARCHAR(10) | | GDS identifier (Amadeus, Sabre, Travelport) |
| accreditationStatus | VARCHAR(20) | | Accredited / Suspended / Terminated |
| effectiveDate | DATE | | Accreditation start date |
| expiryDate | DATE | | Accreditation expiry date (nullable) |

---

### COMMISSION_AGREEMENT
Defines commission rates and override terms agreed between an airline and a travel agency or channel.

| Attribute | Type | Key | Description |
|---|---|---|---|
| commissionAgreementID | VARCHAR(30) | PK | Unique identifier for the commission agreement |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| travelAgencyID | VARCHAR(30) | FK | References TRAVEL_AGENCY.travelAgencyID (nullable for channel-level) |
| salesChannelID | VARCHAR(30) | FK | References SALES_CHANNEL.salesChannelID (nullable for agency-level) |
| commissionType | VARCHAR(20) | | Standard / Override / Incentive / Net |
| commissionBasis | VARCHAR(20) | | BaseFare / TotalFare / Segment / Flat |
| commissionRate | DECIMAL(5,2) | | Commission percentage rate |
| flatAmount | DECIMAL(10,2) | | Fixed flat commission amount (nullable if rate-based) |
| currency | CHAR(3) | | ISO 4217 currency code |
| effectiveDate | DATE | | Agreement effective date |
| expiryDate | DATE | | Agreement expiry date (nullable) |
| status | VARCHAR(20) | | Active / Expired / Terminated |

---

### BOOKING
The core transactional entity representing a passenger booking / Passenger Name Record (PNR).

| Attribute | Type | Key | Description |
|---|---|---|---|
| bookingID | VARCHAR(30) | PK | Unique internal booking identifier |
| pnrCode | VARCHAR(6) | | IATA 6-character alphanumeric PNR locator |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode — operating/validating carrier |
| salesChannelID | VARCHAR(30) | FK | References SALES_CHANNEL.salesChannelID |
| travelAgencyID | VARCHAR(30) | FK | References TRAVEL_AGENCY.travelAgencyID (nullable for direct) |
| offerID | VARCHAR(30) | FK | References OFFER.offerID (nullable) |
| loyaltyProgrammeID | VARCHAR(30) | FK | References LOYALTY_PROGRAMME.loyaltyProgrammeID (nullable) |
| bookingStatus | VARCHAR(20) | | Confirmed / OnHold / Cancelled / NoShow / Flown |
| bookingDateTime | TIMESTAMP | | Date and time booking was created |
| passengerCount | INTEGER | | Total number of passengers on PNR |
| totalFareAmount | DECIMAL(12,2) | | Total fare before taxes |
| totalTaxAmount | DECIMAL(12,2) | | Total tax amount |
| totalAmount | DECIMAL(12,2) | | Grand total charged |
| currency | CHAR(3) | | ISO 4217 currency code |
| lastModifiedDateTime | TIMESTAMP | | Last modification timestamp |

---

### BOOKING_SEGMENT
Represents an individual flight segment within a booking (one leg of the itinerary).

| Attribute | Type | Key | Description |
|---|---|---|---|
| bookingSegmentID | VARCHAR(30) | PK | Unique identifier for the booking segment |
| bookingID | VARCHAR(30) | FK | References BOOKING.bookingID |
| routeID | VARCHAR(30) | FK | References ROUTE.routeID |
| fareFamilyID | VARCHAR(30) | FK | References FARE_FAMILY.fareFamilyID |
| flightNumber | VARCHAR(10) | | Operating flight number |
| operatingCarrierCode | VARCHAR(3) | | IATA code of operating carrier |
| marketingCarrierCode | VARCHAR(3) | | IATA code of marketing carrier |
| departureDateTime | TIMESTAMP | | Scheduled departure date and time |
| arrivalDateTime | TIMESTAMP | | Scheduled arrival date and time |
| originAirportCode | VARCHAR(3) | | IATA departure airport code |
| destinationAirportCode | VARCHAR(3) | | IATA arrival airport code |
| bookingClass | VARCHAR(2) | | RBD booking class (e.g., Y, B, M, J, C) |
| cabinCode | VARCHAR(1) | | IATA cabin code (Y/W/C/F) |
| segmentStatus | VARCHAR(10) | | HK / KL / UN / NO / WK (IATA status codes) |
| fareBasisCode | VARCHAR(20) | | Fare basis code for this segment |
| segmentFareAmount | DECIMAL(10,2) | | Fare amount for this segment |
| currency | CHAR(3) | | ISO 4217 currency code |
| sequenceNumber | INTEGER | | Segment order within the itinerary |

---

### TICKET
Represents an issued electronic ticket (ET) or EMD linked to a booking.

| Attribute | Type | Key | Description |
|---|---|---|---|
| ticketID | VARCHAR(30) | PK | Unique internal ticket identifier |
| bookingID | VARCHAR(30) | FK | References BOOKING.bookingID |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode — validating carrier |
| ffpPartnerID | VARCHAR(30) | FK | References FREQUENT_FLYER_PARTNER.ffpPartnerID (nullable) |
| ticketNumber | VARCHAR(14) | | 13-digit IATA airline ticket number |
| ticketType | VARCHAR(10) | | ET (Electronic Ticket) / EMD-A / EMD-S |
| passengerName | VARCHAR(100) | | Passenger name as on ticket (Last/First) |
| passengerType | VARCHAR(5) | | ADT / CHD / INF / YTH (IATA passenger type codes) |
| issuanceDateTime | TIMESTAMP | | Date and time ticket was issued |
| issuanceCountry | VARCHAR(3) | | ISO 3166-1 alpha-3 country of issuance |
| issuanceChannel | VARCHAR(30) | | Channel through which ticket was issued |
| ticketStatus | VARCHAR(20) | | Open / Used / Refunded / Voided / Exchanged |
| baseFareAmount | DECIMAL(12,2) | | Base fare amount on ticket |
| taxAmount | DECIMAL(12,2) | | Total tax amount on ticket |
| totalAmount | DECIMAL(12,2) | | Total ticket value |
| currency | CHAR(3) | | ISO 4217 currency code |
| couponCount | INTEGER | | Number of flight coupons on ticket (1–4) |

---

### TICKET_COUPON
Represents individual flight coupons within an issued ticket (one coupon per segment).

| Attribute | Type | Key | Description |
|---|---|---|---|
| couponID | VARCHAR(30) | PK | Unique identifier for the ticket coupon |
| ticketID | VARCHAR(30) | FK | References TICKET.ticketID |
| bookingSegmentID | VARCHAR(30) | FK | References BOOKING_SEGMENT.bookingSegmentID |
| couponNumber | INTEGER | | Coupon sequence number (1–4) |
| couponStatus | VARCHAR(10) | | O (Open) / A (Airport Ctrl) / C (Checked-In) / F (Flown) / R (Refunded) |
| flightCouponValue | DECIMAL(10,2) | | Value attributed to this coupon |
| currency | CHAR(3) | | ISO 4217 currency code |
| boardingPassIssued | CHAR(1) | | Y / N — whether boarding pass has been issued |

---

### PAYMENT
Records a payment transaction associated with a booking.

| Attribute | Type | Key | Description |
|---|---|---|---|
| paymentID | VARCHAR(30) | PK | Unique identifier for the payment record |
| bookingID | VARCHAR(30) | FK | References BOOKING.bookingID |
| paymentMethod | VARCHAR(20) | | CreditCard / DebitCard / Cash / Miles / Voucher / BankTransfer |
| paymentStatus | VARCHAR(20) | | Authorised / Captured / Refunded / Failed / Pending |
| paymentAmount | DECIMAL(12,2) | | Amount of this payment transaction |
| currency | CHAR(3) | | ISO 4217 currency code |
| transactionReference | VARCHAR(50) | | Payment gateway or bank transaction reference |
| paymentDateTime | TIMESTAMP | | Date and time payment was processed |
| cardType | VARCHAR(20) | | Visa / MasterCard / Amex / UnionPay (nullable) |
| cardLastFour | CHAR(4) | | Last four digits of payment card (nullable) |
| milesUsed | INTEGER | | Miles redeemed as part of payment (nullable) |
| loyaltyProgrammeID | VARCHAR(30) | FK | References LOYALTY_PROGRAMME.loyaltyProgrammeID (nullable) |

---

### BOOKING_ANCILLARY
Records ancillary services added to a booking segment.

| Attribute | Type | Key | Description |
|---|---|---|---|
| bookingAncillaryID | VARCHAR(30) | PK | Unique identifier for the booking ancillary record |
| bookingSegmentID | VARCHAR(30) | FK | References BOOKING_SEGMENT.bookingSegmentID |
| ancillaryServiceID | VARCHAR(30) | FK | References ANCILLARY_SERVICE.ancillaryServiceID |
| productCatalogueID | VARCHAR(30) | FK | References PRODUCT_CATALOGUE.productCatalogueID (nullable) |
| quantity | INTEGER | | Number of units purchased |
| unitPrice | DECIMAL(10,2) | | Price per unit at time of purchase |
| totalPrice | DECIMAL(10,2) | | Total price for all units |
| currency | CHAR(3) | | ISO 4217 currency code |
| ancillaryStatus | VARCHAR(20) | | Confirmed / Requested / Cancelled / Flown |
| ssrCode | VARCHAR(4) | | IATA SSR code (e.g., AVIH, WCHR, VGML) |
| deliveryStatus | VARCHAR(20) | | Pending / Delivered / NotAvailable |

---

## Relationships

### Internal — Sales Domain

| From Entity | To Entity | Cardinality | Description |
|---|---|---|---|
| SALES_CHANNEL | TRAVEL_AGENCY | 0..* | A channel may have zero or more associated agencies |
| TRAVEL_AGENCY | COMMISSION_AGREEMENT | 0..* | An agency may have zero or more commission agreements |
| SALES_CHANNEL | COMMISSION_AGREEMENT | 0..* | A channel may have zero or more commission agreements |
| BOOKING | BOOKING_SEGMENT | 1..* | A booking contains one or more flight segments |
| BOOKING | TICKET | 1..* | A booking has one or more issued tickets |
| BOOKING | PAYMENT | 1..* | A booking has one or more payment records |
| TICKET | TICKET_COUPON | 1..* | A ticket contains one or more flight coupons |
| BOOKING_SEGMENT | TICKET_COUPON | 0..* | A segment may have one or more coupons |
| BOOKING_SEGMENT | BOOKING_ANCILLARY | 0..* | A segment may have zero or more ancillary services |

### Cross-Domain — Sales → Previous Domains

| From Entity (Sales) | To Entity | Domain | Description |
|---|---|---|---|
| SALES_CHANNEL | AIRLINE | 01 | Channel belongs to an airline |
| TRAVEL_AGENCY | AIRLINE | 01 | Agency contracted with an airline |
| COMMISSION_AGREEMENT | AIRLINE | 01 | Agreement governed by an airline |
| BOOKING | AIRLINE | 01 | Booking made with an airline |
| BOOKING | OFFER | 02 | Booking may reference a marketing offer |
| BOOKING | LOYALTY_PROGRAMME | 02 | Booking may accrue miles to a loyalty programme |
| BOOKING_SEGMENT | ROUTE | 01 | Segment operates on a defined route |
| BOOKING_SEGMENT | FARE_FAMILY | 03 | Segment is sold under a fare family |
| BOOKING_ANCILLARY | ANCILLARY_SERVICE | 03 | Ancillary references the product catalogue |
| BOOKING_ANCILLARY | PRODUCT_CATALOGUE | 03 | Optional deeper product reference |
| TICKET | AIRLINE | 01 | Ticket issued by validating carrier |
| TICKET | FREQUENT_FLYER_PARTNER | 01 | Miles credited via FFP partner |
| PAYMENT | LOYALTY_PROGRAMME | 02 | Miles redemption payment reference |

---

## Enumerations and Code Lists

### CHANNEL_TYPE_CD
| Code | Description |
|---|---|
| DIRECT | Airline direct website or app |
| GDS | Global Distribution System (Amadeus, Sabre, Travelport) |
| NDC | IATA New Distribution Capability channel |
| OTA | Online Travel Agency (Expedia, Booking.com) |
| TMC | Travel Management Company |
| AIRPORT | Airport ticket office / check-in desk |
| CALL_CENTRE | Telephone reservations |

### BOOKING_STATUS_CD
| Code | Description |
|---|---|
| CONFIRMED | Booking is confirmed and active |
| ON_HOLD | Booking is on hold pending payment |
| CANCELLED | Booking has been cancelled |
| NO_SHOW | Passenger did not appear for departure |
| FLOWN | All segments have been flown |

### TICKET_TYPE_CD
| Code | Description |
|---|---|
| ET | Electronic Ticket (flight coupon) |
| EMD_A | Electronic Miscellaneous Document — Associated |
| EMD_S | Electronic Miscellaneous Document — Standalone |

### PASSENGER_TYPE_CD (IATA Standard)
| Code | Description |
|---|---|
| ADT | Adult (12 years and over) |
| CHD | Child (2–11 years) |
| INF | Infant (under 2 years, no seat) |
| YTH | Youth (12–25 years, specific fare) |

### PAYMENT_METHOD_CD
| Code | Description |
|---|---|
| CREDIT_CARD | Credit card payment |
| DEBIT_CARD | Debit card payment |
| CASH | Cash payment at airport or agency |
| MILES | Loyalty miles/points redemption |
| VOUCHER | Airline-issued voucher or credit |
| BANK_TRANSFER | Direct bank transfer |

### SEGMENT_STATUS_CD (IATA Standard)
| Code | Description |
|---|---|
| HK | Holds Confirmed |
| KL | Confirmed from Waitlist |
| UN | Unable to Confirm |
| NO | No Action Taken |
| WK | Was Confirmed — Now Waitlisted |

---

## Business Rules

| Rule ID | Entity | Rule Description |
|---|---|---|
| BR-SAL-001 | BOOKING | pnrCode must be exactly 6 alphanumeric characters (IATA standard) |
| BR-SAL-002 | TICKET | ticketNumber must be exactly 13 numeric digits prefixed by the 3-digit airline code |
| BR-SAL-003 | TICKET | couponCount must equal the number of TICKET_COUPON records linked to this ticket |
| BR-SAL-004 | BOOKING | totalAmount must equal totalFareAmount + totalTaxAmount |
| BR-SAL-005 | BOOKING_SEGMENT | departureDateTime must be earlier than arrivalDateTime |
| BR-SAL-006 | COMMISSION_AGREEMENT | Either travelAgencyID or salesChannelID must be populated — not both null simultaneously |
| BR-SAL-007 | BOOKING_ANCILLARY | totalPrice must equal quantity × unitPrice |
| BR-SAL-008 | PAYMENT | paymentAmount must be greater than zero |
| BR-SAL-009 | TICKET_COUPON | couponNumber must be between 1 and 4 (IATA maximum coupons per ticket) |
| BR-SAL-010 | BOOKING_SEGMENT | sequenceNumber must be unique within a BOOKING |

---

## Traceability to AIDM Domains

| AIDM Domain | AIDM Sub-domain | Entities |
|---|---|---|
| Commercial | Distribution | SALES_CHANNEL, TRAVEL_AGENCY, COMMISSION_AGREEMENT |
| Commercial | Booking & Reservation | BOOKING, BOOKING_SEGMENT |
| Commercial | Ticketing | TICKET, TICKET_COUPON |
| Commercial | Payment | PAYMENT |
| Commercial | Ancillary Sales | BOOKING_ANCILLARY |

---

## Notes

- PNR structure follows IATA Resolution 722 (Passenger Name Record) standards.
- Ticket number format follows IATA Resolution 722b (13-digit electronic ticket number).
- Segment status codes follow IATA ADRM (Airport Development Reference Manual) and IATA Reservations Interline Message Procedures (IIMP).
- SSR codes in BOOKING_ANCILLARY follow IATA PADIS SSR code list.
- NDC channel flag aligns with IATA NDC Standard (21.3+).
