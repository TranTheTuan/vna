# IATA AIDM 25.1 — Domain 05: Customer & Loyalty
## Mockup Data Model for AI Consumption

---

## Subject Area Overview

The **Customer & Loyalty** domain in the IATA AIDM manages the full lifecycle of passenger identity, customer profiles, loyalty programme membership, miles/points accrual and redemption, tier qualification tracking, and customer communication preferences. It is the master source of truth for all customer-centric data across the airline enterprise.

- **AIDM Version**: 25.1
- **Domain**: Commercial — Customer & Loyalty
- **Integrates With**: Network & Alliances (01), Brand & Marketing (02), Product (03), Sales (04)

---

## Cross-Domain FK Dependencies

| Referenced Entity | PK Used | Referenced By (Customer & Loyalty) |
|---|---|---|
| AIRLINE | airlineCode | CUSTOMER_PROFILE, LOYALTY_ACCOUNT, MILES_TRANSACTION |
| LOYALTY_PROGRAMME | loyaltyProgrammeID | LOYALTY_ACCOUNT, TIER_QUALIFICATION |
| LOYALTY_TIER | loyaltyTierID | LOYALTY_ACCOUNT, TIER_QUALIFICATION |
| LOYALTY_TIER_BENEFIT | benefitID | CUSTOMER_BENEFIT_REDEMPTION |
| FREQUENT_FLYER_PARTNER | ffpPartnerID | PARTNER_MILES_TRANSACTION |
| BOOKING | bookingID | MILES_TRANSACTION |
| TICKET | ticketID | MILES_TRANSACTION |
| BOOKING_SEGMENT | bookingSegmentID | MILES_TRANSACTION |
| OFFER | offerID | CUSTOMER_OFFER_ELIGIBILITY |
| FARE_FAMILY | fareFamilyID | TIER_QUALIFICATION |

---

## Entities

### CUSTOMER_PROFILE
The master record of an individual passenger/customer. Single customer identity across all touchpoints.

| Attribute | Type | Key | Description |
|---|---|---|---|
| customerID | VARCHAR(30) | PK | Unique internal customer identifier |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode — primary airline relationship |
| title | VARCHAR(10) | | Mr / Mrs / Ms / Dr / Prof |
| firstName | VARCHAR(50) | | Customer first name |
| lastName | VARCHAR(50) | | Customer last name |
| dateOfBirth | DATE | | Date of birth (ISO 8601) |
| gender | CHAR(1) | | M / F / U (Undisclosed) |
| nationality | VARCHAR(3) | | ISO 3166-1 alpha-3 nationality code |
| countryOfResidence | VARCHAR(3) | | ISO 3166-1 alpha-3 country of residence |
| emailAddress | VARCHAR(255) | | Primary email address |
| phoneNumber | VARCHAR(20) | | Primary phone number (E.164 format) |
| preferredLanguage | VARCHAR(5) | | ISO 639-1 language code (e.g., en, fr, zh) |
| preferredCurrency | CHAR(3) | | ISO 4217 preferred currency |
| customerSegment | VARCHAR(30) | | Leisure / Business / FrequentFlyer / VIP / Corporate |
| kycStatus | VARCHAR(20) | | Verified / Pending / Rejected / NotRequired |
| gdprConsentDate | DATE | | Date GDPR / data consent was given (nullable) |
| marketingOptIn | CHAR(1) | | Y / N — whether customer consents to marketing |
| status | VARCHAR(20) | | Active / Inactive / Blocked / Deceased |
| createdDateTime | TIMESTAMP | | Record creation timestamp |
| lastModifiedDateTime | TIMESTAMP | | Last modification timestamp |

---

### CUSTOMER_IDENTITY_DOCUMENT
Stores travel documents associated with a customer profile (passport, national ID, visa).

| Attribute | Type | Key | Description |
|---|---|---|---|
| documentID | VARCHAR(30) | PK | Unique identifier for the document record |
| customerID | VARCHAR(30) | FK | References CUSTOMER_PROFILE.customerID |
| documentType | VARCHAR(20) | | Passport / NationalID / DriversLicence / ResidencePermit |
| documentNumber | VARCHAR(30) | | Document number as printed |
| issuingCountry | VARCHAR(3) | | ISO 3166-1 alpha-3 country of issue |
| issuingAuthority | VARCHAR(100) | | Issuing authority name (nullable) |
| issueDate | DATE | | Document issue date |
| expiryDate | DATE | | Document expiry date |
| primaryDocument | CHAR(1) | | Y / N — whether this is the primary travel document |
| verificationStatus | VARCHAR(20) | | Verified / Pending / Failed / Expired |

---

### CUSTOMER_ADDRESS
Physical and correspondence addresses for a customer.

| Attribute | Type | Key | Description |
|---|---|---|---|
| addressID | VARCHAR(30) | PK | Unique identifier for the address record |
| customerID | VARCHAR(30) | FK | References CUSTOMER_PROFILE.customerID |
| addressType | VARCHAR(20) | | Home / Business / Billing / Mailing |
| addressLine1 | VARCHAR(100) | | First line of address |
| addressLine2 | VARCHAR(100) | | Second line of address (nullable) |
| city | VARCHAR(50) | | City name |
| stateProvince | VARCHAR(50) | | State or province (nullable) |
| postalCode | VARCHAR(20) | | Postal / ZIP code |
| country | VARCHAR(3) | | ISO 3166-1 alpha-3 country code |
| primaryAddress | CHAR(1) | | Y / N — whether this is the primary address |
| effectiveDate | DATE | | Date address became valid |
| expiryDate | DATE | | Date address became invalid (nullable) |

---

### CUSTOMER_PREFERENCE
Stores travel and service preferences for personalisation and operational use.

| Attribute | Type | Key | Description |
|---|---|---|---|
| preferenceID | VARCHAR(30) | PK | Unique identifier for the preference record |
| customerID | VARCHAR(30) | FK | References CUSTOMER_PROFILE.customerID |
| preferenceCategory | VARCHAR(30) | | Seat / Meal / Communication / Service / Ancillary |
| preferenceCode | VARCHAR(10) | | IATA SSR/OSI code or internal code (e.g., AISLE, VGML) |
| preferenceValue | VARCHAR(100) | | Preference value or description |
| preferenceSource | VARCHAR(20) | | CustomerSet / AIInferred / AgentSet / SystemDefault |
| confidenceScore | DECIMAL(4,3) | | AI/ML confidence score for inferred preferences (0.000–1.000) |
| effectiveDate | DATE | | Date preference was recorded |
| expiryDate | DATE | | Date preference expires (nullable) |

---

### LOYALTY_ACCOUNT
The loyalty membership account linking a customer to a loyalty programme.

| Attribute | Type | Key | Description |
|---|---|---|---|
| loyaltyAccountID | VARCHAR(30) | PK | Unique identifier for the loyalty account |
| customerID | VARCHAR(30) | FK | References CUSTOMER_PROFILE.customerID |
| loyaltyProgrammeID | VARCHAR(30) | FK | References LOYALTY_PROGRAMME.loyaltyProgrammeID |
| loyaltyTierID | VARCHAR(30) | FK | References LOYALTY_TIER.loyaltyTierID — current tier |
| membershipNumber | VARCHAR(20) | | Public-facing membership number |
| enrollmentDate | DATE | | Date customer enrolled in the programme |
| tierAchievedDate | DATE | | Date current tier was achieved |
| tierReviewDate | DATE | | Date current tier is next reviewed |
| lifetimeMiles | INTEGER | | Cumulative lifetime qualifying miles earned |
| currentMilesBalance | INTEGER | | Current redeemable miles balance |
| qualifyingMilesYTD | INTEGER | | Qualifying miles earned in current programme year |
| qualifyingSegmentsYTD | INTEGER | | Qualifying segments flown in current programme year |
| accountStatus | VARCHAR(20) | | Active / Suspended / Closed / Merged |
| pinHash | VARCHAR(255) | | Hashed account PIN (nullable) |

---

### TIER_QUALIFICATION
Tracks progress and history of a customer's tier qualification within a loyalty programme.

| Attribute | Type | Key | Description |
|---|---|---|---|
| tierQualificationID | VARCHAR(30) | PK | Unique identifier for the tier qualification record |
| loyaltyAccountID | VARCHAR(30) | FK | References LOYALTY_ACCOUNT.loyaltyAccountID |
| loyaltyProgrammeID | VARCHAR(30) | FK | References LOYALTY_PROGRAMME.loyaltyProgrammeID |
| loyaltyTierID | VARCHAR(30) | FK | References LOYALTY_TIER.loyaltyTierID — target tier |
| fareFamilyID | VARCHAR(30) | FK | References FARE_FAMILY.fareFamilyID (nullable) |
| qualificationPeriodStart | DATE | | Start of the qualification measurement period |
| qualificationPeriodEnd | DATE | | End of the qualification measurement period |
| milesAtPeriodStart | INTEGER | | Qualifying miles balance at period start |
| milesAtPeriodEnd | INTEGER | | Qualifying miles balance at period end |
| segmentsAtPeriodStart | INTEGER | | Qualifying segments at period start |
| segmentsAtPeriodEnd | INTEGER | | Qualifying segments at period end |
| qualificationResult | VARCHAR(20) | | Achieved / NotAchieved / Retained / Downgraded |
| processedDateTime | TIMESTAMP | | Timestamp qualification was processed |

---

### MILES_TRANSACTION
Records every miles/points earn and redemption event on a loyalty account.

| Attribute | Type | Key | Description |
|---|---|---|---|
| milesTransactionID | VARCHAR(30) | PK | Unique identifier for the miles transaction |
| loyaltyAccountID | VARCHAR(30) | FK | References LOYALTY_ACCOUNT.loyaltyAccountID |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode — earning carrier |
| bookingID | VARCHAR(30) | FK | References BOOKING.bookingID (nullable) |
| ticketID | VARCHAR(30) | FK | References TICKET.ticketID (nullable) |
| bookingSegmentID | VARCHAR(30) | FK | References BOOKING_SEGMENT.bookingSegmentID (nullable) |
| transactionType | VARCHAR(20) | | Earn / Redeem / Adjust / Expire / Reinstate / Transfer |
| transactionSource | VARCHAR(30) | | Flight / Partner / Promotion / ManualAdjust / Expiry |
| milesAmount | INTEGER | | Miles amount (positive = earn, negative = redeem/expire) |
| qualifyingMiles | INTEGER | | Portion of transaction that counts as qualifying miles |
| bonusMiles | INTEGER | | Bonus miles component (e.g., tier multiplier) |
| transactionDateTime | TIMESTAMP | | Date and time of transaction |
| expiryDate | DATE | | Date miles from this transaction expire (nullable) |
| transactionStatus | VARCHAR(20) | | Posted / Pending / Reversed / Disputed |
| referenceCode | VARCHAR(50) | | External reference (e.g., partner transaction ID) |

---

### PARTNER_MILES_TRANSACTION
Records miles earned or redeemed via frequent flyer partner activities (hotels, car hire, retail).

| Attribute | Type | Key | Description |
|---|---|---|---|
| partnerTransactionID | VARCHAR(30) | PK | Unique identifier for the partner transaction |
| loyaltyAccountID | VARCHAR(30) | FK | References LOYALTY_ACCOUNT.loyaltyAccountID |
| ffpPartnerID | VARCHAR(30) | FK | References FREQUENT_FLYER_PARTNER.ffpPartnerID |
| partnerTransactionRef | VARCHAR(50) | | Partner's own transaction reference number |
| partnerActivityType | VARCHAR(30) | | Hotel / CarHire / Retail / Dining / Finance / Other |
| partnerName | VARCHAR(100) | | Name of partner organisation |
| activityDate | DATE | | Date the partner activity occurred |
| activityValue | DECIMAL(12,2) | | Monetary value of the partner activity |
| currency | CHAR(3) | | ISO 4217 currency code |
| milesAwarded | INTEGER | | Miles awarded for this partner activity |
| processingStatus | VARCHAR(20) | | Posted / Pending / Rejected / Reversed |
| processedDateTime | TIMESTAMP | | Timestamp transaction was processed |

---

### CUSTOMER_BENEFIT_REDEMPTION
Records an instance of a customer redeeming a tier benefit.

| Attribute | Type | Key | Description |
|---|---|---|---|
| benefitRedemptionID | VARCHAR(30) | PK | Unique identifier for the benefit redemption event |
| loyaltyAccountID | VARCHAR(30) | FK | References LOYALTY_ACCOUNT.loyaltyAccountID |
| benefitID | VARCHAR(30) | FK | References LOYALTY_TIER_BENEFIT.benefitID |
| redemptionDateTime | TIMESTAMP | | Date and time benefit was redeemed |
| redemptionChannel | VARCHAR(30) | | Airport / Lounge / Digital / Inflight / CallCentre |
| redemptionLocation | VARCHAR(100) | | Location descriptor (e.g., airport code, lounge name) |
| redemptionStatus | VARCHAR(20) | | Used / Cancelled / Pending / Denied |
| unitsRedeemed | INTEGER | | Number of benefit units consumed |
| notes | VARCHAR(255) | | Optional agent or system notes |

---

### CUSTOMER_OFFER_ELIGIBILITY
Records which marketing offers a customer is eligible for, enabling personalisation.

| Attribute | Type | Key | Description |
|---|---|---|---|
| offerEligibilityID | VARCHAR(30) | PK | Unique identifier for the eligibility record |
| customerID | VARCHAR(30) | FK | References CUSTOMER_PROFILE.customerID |
| offerID | VARCHAR(30) | FK | References OFFER.offerID |
| eligibilitySource | VARCHAR(30) | | RuleEngine / AIModel / ManualAssign / LoyaltyTier |
| eligibilityScore | DECIMAL(5,4) | | Propensity score (0.0000–1.0000, nullable for rule-based) |
| eligibilityStatus | VARCHAR(20) | | Eligible / Ineligible / Redeemed / Expired |
| assignedDateTime | TIMESTAMP | | Date and time eligibility was determined |
| expiryDateTime | TIMESTAMP | | Date and time eligibility expires |
| redemptionCount | INTEGER | | Number of times customer has redeemed this offer |
| maxRedemptions | INTEGER | | Maximum allowed redemptions per customer |

---

## Relationships

### Internal — Customer & Loyalty Domain

| From Entity | To Entity | Cardinality | Description |
|---|---|---|---|
| CUSTOMER_PROFILE | CUSTOMER_IDENTITY_DOCUMENT | 1..* | A customer has one or more identity documents |
| CUSTOMER_PROFILE | CUSTOMER_ADDRESS | 1..* | A customer has one or more addresses |
| CUSTOMER_PROFILE | CUSTOMER_PREFERENCE | 0..* | A customer may have zero or more preferences |
| CUSTOMER_PROFILE | LOYALTY_ACCOUNT | 0..* | A customer may hold zero or more loyalty accounts |
| CUSTOMER_PROFILE | CUSTOMER_OFFER_ELIGIBILITY | 0..* | A customer may be eligible for zero or more offers |
| LOYALTY_ACCOUNT | TIER_QUALIFICATION | 0..* | An account has zero or more tier qualification records |
| LOYALTY_ACCOUNT | MILES_TRANSACTION | 0..* | An account has zero or more miles transactions |
| LOYALTY_ACCOUNT | PARTNER_MILES_TRANSACTION | 0..* | An account has zero or more partner transactions |
| LOYALTY_ACCOUNT | CUSTOMER_BENEFIT_REDEMPTION | 0..* | An account has zero or more benefit redemptions |

### Cross-Domain — Customer & Loyalty → Previous Domains

| From Entity | To Entity | Domain | Description |
|---|---|---|---|
| CUSTOMER_PROFILE | AIRLINE | 01 | Customer primary airline relationship |
| LOYALTY_ACCOUNT | LOYALTY_PROGRAMME | 02 | Account belongs to a loyalty programme |
| LOYALTY_ACCOUNT | LOYALTY_TIER | 02 | Account holds a current tier |
| TIER_QUALIFICATION | LOYALTY_PROGRAMME | 02 | Qualification tracked per programme |
| TIER_QUALIFICATION | LOYALTY_TIER | 02 | Qualification targets a specific tier |
| TIER_QUALIFICATION | FARE_FAMILY | 03 | Qualification may be scoped to a fare family |
| MILES_TRANSACTION | AIRLINE | 01 | Miles earned on a specific carrier |
| MILES_TRANSACTION | BOOKING | 04 | Transaction linked to a booking |
| MILES_TRANSACTION | TICKET | 04 | Transaction linked to a ticket |
| MILES_TRANSACTION | BOOKING_SEGMENT | 04 | Transaction linked to a segment |
| PARTNER_MILES_TRANSACTION | FREQUENT_FLYER_PARTNER | 01 | Partner activity via FFP partnership |
| CUSTOMER_BENEFIT_REDEMPTION | LOYALTY_TIER_BENEFIT | 02 | Benefit redeemed from tier entitlement |
| CUSTOMER_OFFER_ELIGIBILITY | OFFER | 02 | Eligibility for a marketing offer |

---

## Enumerations and Code Lists

### CUSTOMER_SEGMENT_CD
| Code | Description |
|---|---|
| LEISURE | Leisure / holiday traveller |
| BUSINESS | Corporate / business traveller |
| FREQUENT_FLYER | High-frequency traveller enrolled in loyalty |
| VIP | VIP or ultra-high-value customer |
| CORPORATE | Managed corporate account traveller |

### KYC_STATUS_CD
| Code | Description |
|---|---|
| VERIFIED | Identity fully verified |
| PENDING | Verification in progress |
| REJECTED | Verification failed |
| NOT_REQUIRED | KYC not required for this customer type |

### MILES_TRANSACTION_TYPE_CD
| Code | Description |
|---|---|
| EARN | Miles credited to account |
| REDEEM | Miles debited for redemption |
| ADJUST | Manual adjustment by agent or system |
| EXPIRE | Miles expired per programme policy |
| REINSTATE | Expired miles reinstated |
| TRANSFER | Miles transferred to/from another account |

### PREFERENCE_CATEGORY_CD
| Code | Description |
|---|---|
| SEAT | Seat location preference (Window / Aisle / etc.) |
| MEAL | Meal preference (VGML / KSML / AVML etc.) |
| COMMUNICATION | Preferred contact channel (Email / SMS / Push) |
| SERVICE | Special service requirement |
| ANCILLARY | Preferred ancillary services |

### ACCOUNT_STATUS_CD
| Code | Description |
|---|---|
| ACTIVE | Account is active and operational |
| SUSPENDED | Account is temporarily suspended |
| CLOSED | Account has been closed |
| MERGED | Account was merged into another account |

### QUALIFICATION_RESULT_CD
| Code | Description |
|---|---|
| ACHIEVED | Customer achieved the target tier |
| NOT_ACHIEVED | Customer did not reach the target tier |
| RETAINED | Customer retained their existing tier |
| DOWNGRADED | Customer was downgraded to a lower tier |

---

## Business Rules

| Rule ID | Entity | Rule Description |
|---|---|---|
| BR-CUS-001 | CUSTOMER_PROFILE | emailAddress must be unique across all active customer records |
| BR-CUS-002 | LOYALTY_ACCOUNT | A customer may hold only one ACTIVE account per LOYALTY_PROGRAMME |
| BR-CUS-003 | MILES_TRANSACTION | milesAmount must be negative for REDEEM and EXPIRE transactions |
| BR-CUS-004 | MILES_TRANSACTION | milesAmount must be positive for EARN and REINSTATE transactions |
| BR-CUS-005 | LOYALTY_ACCOUNT | currentMilesBalance must equal the sum of all posted MILES_TRANSACTION.milesAmount for the account |
| BR-CUS-006 | CUSTOMER_IDENTITY_DOCUMENT | expiryDate must be in the future for documents with verificationStatus = 'Verified' |
| BR-CUS-007 | TIER_QUALIFICATION | qualificationPeriodEnd must be later than qualificationPeriodStart |
| BR-CUS-008 | CUSTOMER_OFFER_ELIGIBILITY | redemptionCount must not exceed maxRedemptions |
| BR-CUS-009 | CUSTOMER_PREFERENCE | confidenceScore must be between 0.000 and 1.000 when populated |
| BR-CUS-010 | PARTNER_MILES_TRANSACTION | milesAwarded must be greater than zero |

---

## Traceability to AIDM Domains

| AIDM Domain | AIDM Sub-domain | Entities |
|---|---|---|
| Customer | Identity Management | CUSTOMER_PROFILE, CUSTOMER_IDENTITY_DOCUMENT, CUSTOMER_ADDRESS |
| Customer | Personalisation | CUSTOMER_PREFERENCE, CUSTOMER_OFFER_ELIGIBILITY |
| Customer | Loyalty Membership | LOYALTY_ACCOUNT, TIER_QUALIFICATION |
| Customer | Miles Ledger | MILES_TRANSACTION, PARTNER_MILES_TRANSACTION |
| Customer | Benefit Fulfilment | CUSTOMER_BENEFIT_REDEMPTION |

---

## Notes

- Customer identity structures align with IATA Resolution 830d (Passenger Name Record) and GDPR Article 17 (Right to Erasure) compliance requirements.
- Miles/points ledger design follows double-entry accounting principles — every earn has a matching liability on the airline's balance sheet.
- Preference codes align with IATA PADIS SSR and OSI code tables.
- KYC (Know Your Customer) status supports AML/CTF regulatory compliance requirements.
- confidenceScore in CUSTOMER_PREFERENCE supports IATA ONE Order AI/ML personalisation frameworks.
