# IATA AIDM 25.1 — Brand & Marketing
## Mockup Data Model for AI Consumption

---

## Subject Area Overview

The **Brand & Marketing** subject area in the IATA AIDM covers customer loyalty programme management, brand identity governance, marketing campaign execution, personalisation attributes, brand touchpoints at alliance level, and commercial offer management.

- **AIDM Version**: 25.1
- **Domain**: Commercial
- **Subject Area**: Brand & Marketing
- **Integrates With**: Network & Alliances (iata_aidm_network_alliances.md)

---

## Cross-Domain FK Dependencies (Network & Alliances → Brand & Marketing)

| Referenced Entity (N&A) | PK Used | Referenced By (B&M) |
|---|---|---|
| AIRLINE | airlineCode | BRAND, LOYALTY_PROGRAMME, LOYALTY_TIER_BENEFIT, MARKETING_CAMPAIGN, OFFER |
| ALLIANCE | allianceCode | BRAND_ALLIANCE_TOUCHPOINT |
| ALLIANCE_MEMBERSHIP | membershipID | LOYALTY_TIER_BENEFIT |
| ROUTE_NETWORK | routeNetworkID | CAMPAIGN_ROUTE_TARGET |
| ROUTE | routeID | CAMPAIGN_ROUTE_TARGET, OFFER_ROUTE_APPLICABILITY |
| FREQUENT_FLYER_PARTNER | ffpPartnerID | LOYALTY_TIER_BENEFIT |

---

## Entities

### BRAND
Core brand identity record for an airline. Governs visual identity, brand tier classification, and guidelines used across all customer touchpoints.

| Attribute | Type | Key | Description |
|---|---|---|---|
| brandID | VARCHAR(30) | PK | Unique identifier for the brand record |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| brandName | VARCHAR(100) | | Commercial brand name (may differ from airline legal name) |
| brandTier | VARCHAR(20) | | Premium / Economy / Regional / Low-Cost |
| brandGuidelines | VARCHAR(255) | | URL to official brand guidelines document |
| primaryColourHex | CHAR(7) | | Primary brand colour in HEX format (e.g., #1A2B3C) |
| secondaryColourHex | CHAR(7) | | Secondary brand colour in HEX format |
| logoURL | VARCHAR(255) | | URL to master logo asset |
| effectiveDate | DATE | | Date the brand identity comes into force |
| expiryDate | DATE | | Date the brand identity is retired (nullable) |

---

### BRAND_ALLIANCE_TOUCHPOINT
Maps how an airline's brand is presented within a specific alliance context (e.g., co-branded lounge, shared check-in, digital co-presentation).

| Attribute | Type | Key | Description |
|---|---|---|---|
| touchpointID | VARCHAR(30) | PK | Unique identifier for the touchpoint record |
| brandID | VARCHAR(30) | FK | References BRAND.brandID |
| allianceCode | VARCHAR(10) | FK | References ALLIANCE.allianceCode |
| touchpointType | VARCHAR(30) | | Lounge / Check-in / Digital / Inflight / Co-brand |
| displayPriority | INTEGER | | Rendering order when multiple brands are co-presented (1=highest) |
| cobranded | CHAR(1) | | Y / N — whether this touchpoint uses co-branded assets |
| effectiveDate | DATE | | Date touchpoint agreement comes into force |
| expiryDate | DATE | | Date touchpoint agreement expires (nullable) |

---

### LOYALTY_PROGRAMME
Defines a frequent flyer or loyalty programme owned and operated by an airline.

| Attribute | Type | Key | Description |
|---|---|---|---|
| loyaltyProgrammeID | VARCHAR(30) | PK | Unique identifier for the loyalty programme |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode — programme owner |
| programmeName | VARCHAR(100) | | Name of the programme (e.g., Executive Club, SkyMiles) |
| programmeType | VARCHAR(20) | | Miles / Points / Cashback |
| currency | CHAR(3) | | ISO 4217 currency code for cashback programmes |
| expiryPolicy | VARCHAR(50) | | Miles/points expiry rule description |
| status | VARCHAR(20) | | Active / Inactive / Suspended |
| launchDate | DATE | | Date the programme was launched |

---

### LOYALTY_TIER
Defines a membership tier within a loyalty programme (e.g., Silver, Gold, Platinum).

| Attribute | Type | Key | Description |
|---|---|---|---|
| loyaltyTierID | VARCHAR(30) | PK | Unique identifier for the loyalty tier |
| loyaltyProgrammeID | VARCHAR(30) | FK | References LOYALTY_PROGRAMME.loyaltyProgrammeID |
| tierName | VARCHAR(50) | | Tier label (e.g., Silver / Gold / Platinum / Elite) |
| tierLevel | INTEGER | | Numeric level — 1 = base tier, higher = more premium |
| minMilesRequired | INTEGER | | Minimum qualifying miles/points to achieve this tier |
| minSegmentsRequired | INTEGER | | Minimum qualifying flight segments to achieve this tier |
| tierStatus | VARCHAR(20) | | Active / Inactive / Deprecated |
| effectiveDate | DATE | | Date the tier definition comes into force |
| expiryDate | DATE | | Date the tier is retired (nullable) |

---

### LOYALTY_TIER_BENEFIT
Specific benefits granted to members of a loyalty tier. Optionally linked to an alliance membership or FFP partner for partner-specific benefit delivery.

| Attribute | Type | Key | Description |
|---|---|---|---|
| benefitID | VARCHAR(30) | PK | Unique identifier for the benefit record |
| loyaltyTierID | VARCHAR(30) | FK | References LOYALTY_TIER.loyaltyTierID |
| ffpPartnerID | VARCHAR(30) | FK | References FREQUENT_FLYER_PARTNER.ffpPartnerID (nullable) |
| membershipID | VARCHAR(20) | FK | References ALLIANCE_MEMBERSHIP.membershipID (nullable) |
| benefitType | VARCHAR(30) | | Lounge / Upgrade / ExtraBaggage / FastTrack / BonusMiles |
| benefitValue | VARCHAR(100) | | Quantified benefit description (e.g., "2 x lounge visits", "+50% miles") |
| applicableCarrier | VARCHAR(3) | FK | References AIRLINE.airlineCode — carrier granting benefit (nullable) |
| effectiveDate | DATE | | Date the benefit comes into force |
| expiryDate | DATE | | Date the benefit expires (nullable) |

---

### MARKETING_CAMPAIGN
Tracks a marketing campaign executed by an airline across one or more channels, targeting specific customer segments.

| Attribute | Type | Key | Description |
|---|---|---|---|
| campaignID | VARCHAR(30) | PK | Unique identifier for the campaign |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode — campaign owner |
| campaignName | VARCHAR(100) | | Descriptive name of the campaign |
| campaignType | VARCHAR(30) | | Email / Social / Display / Metasearch / Retargeting / OTA |
| targetSegment | VARCHAR(50) | | Customer segment targeted (e.g., Leisure / Business / FFP-Gold) |
| channel | VARCHAR(30) | | Primary distribution channel |
| budget | DECIMAL(12,2) | | Total campaign budget |
| currency | CHAR(3) | | ISO 4217 currency code for budget |
| startDate | DATE | | Campaign start date |
| endDate | DATE | | Campaign end date |
| status | VARCHAR(20) | | Draft / Active / Paused / Completed / Cancelled |
| KPI_target | VARCHAR(100) | | Key performance indicator target description |

---

### CAMPAIGN_ROUTE_TARGET
Junction entity linking a marketing campaign to specific route networks or individual routes being promoted.

| Attribute | Type | Key | Description |
|---|---|---|---|
| campaignRouteID | VARCHAR(30) | PK | Unique identifier for the campaign-route link |
| campaignID | VARCHAR(30) | FK | References MARKETING_CAMPAIGN.campaignID |
| routeNetworkID | VARCHAR(30) | FK | References ROUTE_NETWORK.routeNetworkID |
| routeID | VARCHAR(30) | FK | References ROUTE.routeID — optional route-level targeting (nullable) |
| promotionType | VARCHAR(30) | | Flash Sale / Seasonal / Launch / Milestone |
| discountType | VARCHAR(20) | | Percentage / Fixed / Miles-Bonus |
| discountValue | DECIMAL(8,2) | | Discount amount or percentage value |
| effectiveDate | DATE | | Promotion start date |
| expiryDate | DATE | | Promotion end date |

---

### OFFER
A commercial offer made available to customers, optionally produced as part of a marketing campaign and tied to a fare class or ancillary product.

| Attribute | Type | Key | Description |
|---|---|---|---|
| offerID | VARCHAR(30) | PK | Unique identifier for the offer |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode — offer owner |
| campaignID | VARCHAR(30) | FK | References MARKETING_CAMPAIGN.campaignID (nullable) |
| offerType | VARCHAR(30) | | Fare / Ancillary / Bundle / Upgrade / Miles-Redemption |
| offerName | VARCHAR(100) | | Descriptive name of the offer |
| offerDescription | TEXT | | Full description for customer-facing display |
| currency | CHAR(3) | | ISO 4217 currency code |
| offerPrice | DECIMAL(12,2) | | Base price of the offer |
| fareClass | VARCHAR(5) | | Booking class or fare basis code applicable |
| effectiveDate | DATE | | Offer start date |
| expiryDate | DATE | | Offer end date |
| status | VARCHAR(20) | | Draft / Active / Expired / Withdrawn |

---

### OFFER_ROUTE_APPLICABILITY
Resolves the specific routes, cabin classes, and travel-day conditions under which an offer is valid.

| Attribute | Type | Key | Description |
|---|---|---|---|
| offerRouteID | VARCHAR(30) | PK | Unique identifier for the offer-route applicability record |
| offerID | VARCHAR(30) | FK | References OFFER.offerID |
| routeID | VARCHAR(30) | FK | References ROUTE.routeID |
| cabinClass | VARCHAR(10) | | Applicable cabin class (Economy / Business / First) |
| dayOfWeekMask | VARCHAR(7) | | 7-char bitmask Mon–Sun (e.g., 1111100 = Mon–Fri only) |
| minAdvancePurchase | INTEGER | | Minimum days before departure required to use offer |
| maxAdvancePurchase | INTEGER | | Maximum days before departure allowed to use offer |
| effectiveDate | DATE | | Start date of applicability |
| expiryDate | DATE | | End date of applicability |

---

## Relationships

### Internal — Brand & Marketing

| From Entity | To Entity | Cardinality | Description |
|---|---|---|---|
| AIRLINE | BRAND | 1..* | An airline owns one or more brand identities |
| BRAND | BRAND_ALLIANCE_TOUCHPOINT | 0..* | A brand may be presented at zero or more alliance touchpoints |
| AIRLINE | LOYALTY_PROGRAMME | 1..* | An airline operates one or more loyalty programmes |
| LOYALTY_PROGRAMME | LOYALTY_TIER | 1..* | A programme has one or more membership tiers |
| LOYALTY_TIER | LOYALTY_TIER_BENEFIT | 0..* | A tier grants zero or more specific benefits |
| AIRLINE | MARKETING_CAMPAIGN | 0..* | An airline runs zero or more marketing campaigns |
| MARKETING_CAMPAIGN | CAMPAIGN_ROUTE_TARGET | 0..* | A campaign may target zero or more route networks/routes |
| MARKETING_CAMPAIGN | OFFER | 0..* | A campaign may produce zero or more offers |
| OFFER | OFFER_ROUTE_APPLICABILITY | 0..* | An offer may apply to zero or more routes |

### Cross-Domain — Brand & Marketing → Network & Alliances

| From Entity (B&M) | To Entity (N&A) | Cardinality | Description |
|---|---|---|---|
| BRAND_ALLIANCE_TOUCHPOINT | ALLIANCE | 0..* → 1 | A touchpoint is hosted by exactly one alliance |
| LOYALTY_TIER_BENEFIT | FREQUENT_FLYER_PARTNER | 0..* → 0..1 | Benefit may be enabled by an FFP partnership |
| LOYALTY_TIER_BENEFIT | ALLIANCE_MEMBERSHIP | 0..* → 0..1 | Benefit may require a specific alliance membership |
| CAMPAIGN_ROUTE_TARGET | ROUTE_NETWORK | 0..* → 1 | Target always references a route network |
| CAMPAIGN_ROUTE_TARGET | ROUTE | 0..* → 0..1 | Target may optionally narrow to a specific route |
| OFFER_ROUTE_APPLICABILITY | ROUTE | 0..* → 1 | Applicability references exactly one route |

---

## Enumerations and Code Lists

### BRAND_TIER_CD
| Code | Description |
|---|---|
| PREMIUM | Full-service premium airline brand |
| ECONOMY | Full-service economy-focused brand |
| REGIONAL | Regional/feeder carrier brand |
| LOW_COST | Low-cost carrier brand |

### TOUCHPOINT_TYPE_CD
| Code | Description |
|---|---|
| LOUNGE | Airport lounge co-branding |
| CHECK_IN | Check-in desk co-branding |
| DIGITAL | Website / app digital co-presentation |
| INFLIGHT | Inflight service co-branding |
| CO_BRAND | General co-branded marketing material |

### PROGRAMME_TYPE_CD
| Code | Description |
|---|---|
| MILES | Accrual based on flown miles |
| POINTS | Accrual based on spend or activity points |
| CASHBACK | Monetary cashback on purchases |

### BENEFIT_TYPE_CD
| Code | Description |
|---|---|
| LOUNGE | Airport lounge access |
| UPGRADE | Complimentary or discounted cabin upgrade |
| EXTRA_BAGGAGE | Additional baggage allowance |
| FAST_TRACK | Priority security / boarding / check-in lane |
| BONUS_MILES | Accelerated miles/points earning rate |

### CAMPAIGN_TYPE_CD
| Code | Description |
|---|---|
| EMAIL | Direct email marketing |
| SOCIAL | Social media advertising |
| DISPLAY | Display / banner advertising |
| METASEARCH | Metasearch engine promotion (Skyscanner, Google Flights) |
| RETARGETING | Retargeting / remarketing to past visitors |
| OTA | Online travel agency co-marketing |

### OFFER_TYPE_CD
| Code | Description |
|---|---|
| FARE | Discounted fare offer |
| ANCILLARY | Ancillary product offer (seat, meal, bag) |
| BUNDLE | Combined fare + ancillary bundle |
| UPGRADE | Upgrade offer (paid or complimentary) |
| MILES_REDEMPTION | Miles/points redemption offer |

### PROMOTION_TYPE_CD
| Code | Description |
|---|---|
| FLASH_SALE | Limited-time flash sale (typically 24–72 hours) |
| SEASONAL | Seasonal promotion tied to IATA season |
| LAUNCH | New route or product launch promotion |
| MILESTONE | Anniversary or milestone celebration offer |

### DISCOUNT_TYPE_CD
| Code | Description |
|---|---|
| PERCENTAGE | Percentage discount off standard fare |
| FIXED | Fixed monetary amount discount |
| MILES_BONUS | Bonus miles/points awarded instead of price reduction |

### CAMPAIGN_STATUS_CD
| Code | Description |
|---|---|
| DRAFT | Campaign is being planned, not yet active |
| ACTIVE | Campaign is currently running |
| PAUSED | Campaign is temporarily suspended |
| COMPLETED | Campaign has run its full course |
| CANCELLED | Campaign was cancelled before completion |

---

## Business Rules

| Rule ID | Entity | Rule Description |
|---|---|---|
| BR-BM-001 | BRAND_ALLIANCE_TOUCHPOINT | A BRAND_ALLIANCE_TOUCHPOINT can only exist for an AIRLINE where allianceMember = TRUE |
| BR-BM-002 | LOYALTY_TIER | tierLevel values within a single LOYALTY_PROGRAMME must be unique |
| BR-BM-003 | LOYALTY_TIER_BENEFIT | ffpPartnerID and membershipID are both nullable, but at least one of (ffpPartnerID, membershipID, applicableCarrier) must be populated |
| BR-BM-004 | OFFER | effectiveDate and expiryDate must fall within the parent MARKETING_CAMPAIGN startDate–endDate range when campaignID is not null |
| BR-BM-005 | OFFER_ROUTE_APPLICABILITY | routeID must reference a ROUTE with routeStatus = 'ACTIVE' |
| BR-BM-006 | CAMPAIGN_ROUTE_TARGET | When routeID is populated, that ROUTE must belong to the referenced routeNetworkID |
| BR-BM-007 | LOYALTY_TIER | minMilesRequired must increase with each higher tierLevel within the same programme |
| BR-BM-008 | MARKETING_CAMPAIGN | endDate must be later than startDate |

---

## Traceability to AIDM Domains

| AIDM Domain | AIDM Sub-domain | Entities |
|---|---|---|
| Commercial | Brand Management | BRAND, BRAND_ALLIANCE_TOUCHPOINT |
| Commercial | Loyalty Management | LOYALTY_PROGRAMME, LOYALTY_TIER, LOYALTY_TIER_BENEFIT |
| Commercial | Campaign Management | MARKETING_CAMPAIGN, CAMPAIGN_ROUTE_TARGET |
| Commercial | Offer Management | OFFER, OFFER_ROUTE_APPLICABILITY |

---

## Cross-Domain References Summary

| Entity (Brand & Marketing) | Referenced Entity | Referenced Domain |
|---|---|---|
| BRAND.airlineCode | AIRLINE.airlineCode | Network & Alliances |
| BRAND_ALLIANCE_TOUCHPOINT.allianceCode | ALLIANCE.allianceCode | Network & Alliances |
| LOYALTY_PROGRAMME.airlineCode | AIRLINE.airlineCode | Network & Alliances |
| LOYALTY_TIER_BENEFIT.ffpPartnerID | FREQUENT_FLYER_PARTNER.ffpPartnerID | Network & Alliances |
| LOYALTY_TIER_BENEFIT.membershipID | ALLIANCE_MEMBERSHIP.membershipID | Network & Alliances |
| LOYALTY_TIER_BENEFIT.applicableCarrier | AIRLINE.airlineCode | Network & Alliances |
| MARKETING_CAMPAIGN.airlineCode | AIRLINE.airlineCode | Network & Alliances |
| CAMPAIGN_ROUTE_TARGET.routeNetworkID | ROUTE_NETWORK.routeNetworkID | Network & Alliances |
| CAMPAIGN_ROUTE_TARGET.routeID | ROUTE.routeID | Network & Alliances |
| OFFER.airlineCode | AIRLINE.airlineCode | Network & Alliances |
| OFFER_ROUTE_APPLICABILITY.routeID | ROUTE.routeID | Network & Alliances |

---

## Notes

- This model is based on IATA AIDM version 25.1 and is a domain extension of iata_aidm_network_alliances.md.
- Loyalty programme structures align with IATA Resolution 787 (Frequent Flyer Programmes) principles.
- Campaign and offer management structures are compatible with IATA NDC (New Distribution Capability) Offer & Order framework.
- All date fields follow ISO 8601 (YYYY-MM-DD) format.
- All currency codes follow ISO 4217 standard (CHAR 3).
- Colour hex codes follow CSS/HTML hex format including leading `#` character (CHAR 7).
