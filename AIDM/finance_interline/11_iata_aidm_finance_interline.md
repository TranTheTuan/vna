# IATA AIDM 25.1 — Domain 11: Finance & Interline
## Mockup Data Model for AI Consumption

---

## Subject Area Overview

The **Finance & Interline** domain in the IATA AIDM covers airline financial management — including general ledger accounts, cost centre allocation, interline billing and settlement (BSP/CASS/ICH), proration of interline revenue, accounts payable and receivable, budget and forecast management, currency exchange rate management, and financial reporting. It is the financial backbone that consolidates revenue and cost data from all operational and commercial domains.

- **AIDM Version**: 25.1
- **Domain**: Finance — Finance & Interline
- **Integrates With**: Network & Alliances (01), Sales (04), Revenue Management & Pricing (06), Ground Operations (07), Flight Operations (08), Cargo (09), MRO (10)

---

## Cross-Domain FK Dependencies

| Referenced Entity | PK Used | Referenced By (Finance & Interline) |
|---|---|---|
| AIRLINE | airlineCode | GL_ACCOUNT, COST_CENTRE, INTERLINE_AGREEMENT, BSP_SETTLEMENT, CASS_SETTLEMENT, BUDGET |
| ROUTE | routeID | PRORATION_RULE, INTERLINE_REVENUE_PRORATION |
| BOOKING | bookingID | INTERLINE_REVENUE_PRORATION |
| TICKET | ticketID | INTERLINE_BILLING_COUPON, INTERLINE_REVENUE_PRORATION |
| TICKET_COUPON | ticketCouponID | INTERLINE_BILLING_COUPON |
| REVENUE_RECORD | revenueRecordID | GL_POSTING |
| CARGO_REVENUE_RECORD | cargoRevenueID | GL_POSTING |
| WORK_ORDER | workOrderID | ACCOUNTS_PAYABLE_INVOICE |
| GROUND_HANDLER | groundHandlerID | ACCOUNTS_PAYABLE_INVOICE |

---

## Entities

### GL_ACCOUNT
Defines a General Ledger account within the airline's chart of accounts.

| Attribute | Type | Key | Description |
|---|---|---|---|
| glAccountID | VARCHAR(30) | PK | Unique identifier for the GL account |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| accountCode | VARCHAR(20) | | GL account code (e.g., 4100, 5200-01) |
| accountName | VARCHAR(100) | | Full account name |
| accountType | VARCHAR(20) | | Asset / Liability / Equity / Revenue / Expense |
| accountSubType | VARCHAR(30) | | PassengerRevenue / CargoRevenue / FuelExpense / CrewCost etc. |
| normalBalance | VARCHAR(6) | | Debit / Credit |
| parentAccountCode | VARCHAR(20) | | Parent account for hierarchical grouping (nullable) |
| currency | CHAR(3) | | Functional currency for this account |
| iataAccountCode | VARCHAR(10) | | IATA Uniform System of Accounts (USA) code |
| costCentreID | VARCHAR(30) | FK | References COST_CENTRE.costCentreID (nullable) |
| status | VARCHAR(20) | | Active / Inactive |
| effectiveDate | DATE | | Account effective date |

---

### COST_CENTRE
Defines an organisational cost centre for expense allocation and management reporting.

| Attribute | Type | Key | Description |
|---|---|---|---|
| costCentreID | VARCHAR(30) | PK | Unique identifier for the cost centre |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| costCentreCode | VARCHAR(20) | | Cost centre code |
| costCentreName | VARCHAR(100) | | Full name of cost centre |
| costCentreType | VARCHAR(30) | | FlightOps / Maintenance / Commercial / Corporate / Airport |
| parentCostCentreID | VARCHAR(30) | FK | References COST_CENTRE.costCentreID — parent (nullable) |
| managerEmployeeID | VARCHAR(30) | | HR employee ID of cost centre manager |
| airportCode | VARCHAR(3) | | IATA airport code if station-based (nullable) |
| status | VARCHAR(20) | | Active / Inactive |

---

### GL_POSTING
Records an individual debit or credit posting to a GL account.

| Attribute | Type | Key | Description |
|---|---|---|---|
| glPostingID | VARCHAR(30) | PK | Unique identifier for the GL posting |
| glAccountID | VARCHAR(30) | FK | References GL_ACCOUNT.glAccountID |
| costCentreID | VARCHAR(30) | FK | References COST_CENTRE.costCentreID (nullable) |
| revenueRecordID | VARCHAR(30) | FK | References REVENUE_RECORD.revenueRecordID (nullable) |
| cargoRevenueID | VARCHAR(30) | FK | References CARGO_REVENUE_RECORD.cargoRevenueID (nullable) |
| postingType | VARCHAR(20) | | Debit / Credit |
| postingAmount | DECIMAL(14,2) | | Amount of the posting |
| currency | CHAR(3) | | ISO 4217 transaction currency |
| functionalAmount | DECIMAL(14,2) | | Amount in airline functional currency |
| functionalCurrency | CHAR(3) | | Airline functional currency |
| exchangeRate | DECIMAL(12,6) | | Exchange rate applied |
| postingDate | DATE | | Accounting posting date |
| accountingPeriod | VARCHAR(7) | | Accounting period (YYYY-MM) |
| journalReference | VARCHAR(50) | | Journal batch or document reference |
| postingDescription | VARCHAR(255) | | Description of the posting |
| sourceSystem | VARCHAR(30) | | Source system generating the posting |
| postedDateTime | TIMESTAMP | | Timestamp posting was created |

---

### EXCHANGE_RATE
Stores daily currency exchange rates used for financial conversions.

| Attribute | Type | Key | Description |
|---|---|---|---|
| exchangeRateID | VARCHAR(30) | PK | Unique identifier for the exchange rate record |
| fromCurrency | CHAR(3) | | ISO 4217 source currency code |
| toCurrency | CHAR(3) | | ISO 4217 target currency code |
| rateType | VARCHAR(20) | | IATA / BSP / Central / Spot / Budget |
| exchangeRate | DECIMAL(14,6) | | Exchange rate (1 unit of fromCurrency = N toCurrency) |
| effectiveDate | DATE | | Rate effective date |
| expiryDate | DATE | | Rate expiry date (nullable) |
| rateSource | VARCHAR(50) | | Source of rate (e.g., IATA Rate, ECB, Reuters) |

---

### INTERLINE_AGREEMENT
Defines a bilateral or multilateral interline traffic agreement between two airlines.

| Attribute | Type | Key | Description |
|---|---|---|---|
| interlineAgreementID | VARCHAR(30) | PK | Unique identifier for the interline agreement |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode — home carrier |
| partnerAirlineCode | VARCHAR(3) | | IATA code of interline partner airline |
| agreementType | VARCHAR(20) | | Prorate / Interline / SPA / MIA / Codeshare |
| settlementMethod | VARCHAR(10) | | ICH / Bilateral / BSP |
| defaultProrateMethod | VARCHAR(20) | | Mileage / Equal / Published / Negotiated |
| billingCurrency | CHAR(3) | | ISO 4217 billing currency |
| effectiveDate | DATE | | Agreement effective date |
| expiryDate | DATE | | Agreement expiry date (nullable) |
| status | VARCHAR(20) | | Active / Suspended / Terminated |
| iataAgreementRef | VARCHAR(30) | | IATA agreement registration reference |

---

### PRORATION_RULE
Defines the revenue proration rules applied to an interline journey.

| Attribute | Type | Key | Description |
|---|---|---|---|
| prorationRuleID | VARCHAR(30) | PK | Unique identifier for the proration rule |
| interlineAgreementID | VARCHAR(30) | FK | References INTERLINE_AGREEMENT.interlineAgreementID |
| routeID | VARCHAR(30) | FK | References ROUTE.routeID (nullable — may be market-level) |
| prorationMethod | VARCHAR(20) | | Mileage / Equal / Sector / Published / IATA-RPDG |
| proratePercentage | DECIMAL(6,4) | | Fixed prorate percentage (nullable if mileage-based) |
| minimumProrate | DECIMAL(10,2) | | Minimum prorate amount |
| currency | CHAR(3) | | ISO 4217 currency |
| effectiveDate | DATE | | Rule effective date |
| expiryDate | DATE | | Rule expiry date (nullable) |

---

### INTERLINE_BILLING_COUPON
Records an individual interline billing coupon — a segment of a ticket flown on a partner carrier.

| Attribute | Type | Key | Description |
|---|---|---|---|
| billingCouponID | VARCHAR(30) | PK | Unique identifier for the billing coupon |
| ticketID | VARCHAR(30) | FK | References TICKET.ticketID |
| ticketCouponID | VARCHAR(30) | FK | References TICKET_COUPON.ticketCouponID |
| interlineAgreementID | VARCHAR(30) | FK | References INTERLINE_AGREEMENT.interlineAgreementID |
| issuingCarrierCode | VARCHAR(3) | | IATA code of issuing/validating carrier |
| operatingCarrierCode | VARCHAR(3) | | IATA code of operating carrier |
| originAirportCode | VARCHAR(3) | | IATA origin airport of coupon |
| destinationAirportCode | VARCHAR(3) | | IATA destination airport of coupon |
| couponFaceValue | DECIMAL(12,2) | | Face value of coupon in ticket currency |
| proratedAmount | DECIMAL(12,2) | | Prorated amount due to operating carrier |
| billingCurrency | CHAR(3) | | ISO 4217 billing currency |
| billingStatus | VARCHAR(20) | | Unbilled / Billed / Settled / Disputed |
| flightDate | DATE | | Date of travel on this coupon |
| billingPeriod | VARCHAR(7) | | Billing period (YYYY-MM) |

---

### INTERLINE_REVENUE_PRORATION
Records the revenue proration calculation for an interline booking.

| Attribute | Type | Key | Description |
|---|---|---|---|
| prorationID | VARCHAR(30) | PK | Unique identifier for the proration record |
| bookingID | VARCHAR(30) | FK | References BOOKING.bookingID |
| ticketID | VARCHAR(30) | FK | References TICKET.ticketID |
| prorationRuleID | VARCHAR(30) | FK | References PRORATION_RULE.prorationRuleID |
| routeID | VARCHAR(30) | FK | References ROUTE.routeID |
| totalFareAmount | DECIMAL(12,2) | | Total fare to be prorated |
| homeCarrierProrate | DECIMAL(12,2) | | Amount allocated to home carrier |
| partnerCarrierProrate | DECIMAL(12,2) | | Amount allocated to interline partner |
| prorateMethod | VARCHAR(20) | | Method used (Mileage / Equal / Sector) |
| homeCarrierMileage | INTEGER | | Mileage on home carrier segment |
| partnerCarrierMileage | INTEGER | | Mileage on partner carrier segment |
| totalMileage | INTEGER | | Total journey mileage |
| currency | CHAR(3) | | ISO 4217 currency |
| calculationDateTime | TIMESTAMP | | Timestamp proration was calculated |

---

### BSP_SETTLEMENT
Records a Billing and Settlement Plan settlement batch between an airline and IATA BSP.

| Attribute | Type | Key | Description |
|---|---|---|---|
| bspSettlementID | VARCHAR(30) | PK | Unique identifier for the BSP settlement record |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| bspCountryCode | VARCHAR(3) | | ISO 3166-1 alpha-3 BSP market country |
| settlementPeriodStart | DATE | | Settlement reporting period start |
| settlementPeriodEnd | DATE | | Settlement reporting period end |
| settlementDate | DATE | | Date of actual financial settlement |
| totalSales | DECIMAL(14,2) | | Total ticket sales in period |
| totalRefunds | DECIMAL(12,2) | | Total refunds processed in period |
| totalCommissions | DECIMAL(12,2) | | Total commissions paid to agents |
| totalTaxes | DECIMAL(12,2) | | Total taxes collected |
| netSettlementAmount | DECIMAL(14,2) | | Net amount settled with IATA BSP |
| currency | CHAR(3) | | ISO 4217 settlement currency |
| settlementStatus | VARCHAR(20) | | Pending / Submitted / Settled / Disputed |
| iataReferenceNumber | VARCHAR(30) | | IATA BSP settlement reference number |

---

### CASS_SETTLEMENT
Records a Cargo Accounts Settlement System settlement batch for cargo interline billing.

| Attribute | Type | Key | Description |
|---|---|---|---|
| cassSettlementID | VARCHAR(30) | PK | Unique identifier for the CASS settlement record |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| cassCountryCode | VARCHAR(3) | | ISO 3166-1 alpha-3 CASS market country |
| settlementPeriodStart | DATE | | Settlement period start date |
| settlementPeriodEnd | DATE | | Settlement period end date |
| settlementDate | DATE | | Date of actual financial settlement |
| totalFreightCharges | DECIMAL(14,2) | | Total freight charges in period |
| totalSurcharges | DECIMAL(12,2) | | Total surcharges collected |
| totalCommissions | DECIMAL(12,2) | | Total commissions paid |
| netSettlementAmount | DECIMAL(14,2) | | Net amount settled with IATA CASS |
| currency | CHAR(3) | | ISO 4217 settlement currency |
| settlementStatus | VARCHAR(20) | | Pending / Submitted / Settled / Disputed |
| iataReferenceNumber | VARCHAR(30) | | IATA CASS settlement reference number |

---

### ACCOUNTS_PAYABLE_INVOICE
Records a supplier invoice for airline payables — covering ground handlers, MRO, fuel, catering.

| Attribute | Type | Key | Description |
|---|---|---|---|
| apInvoiceID | VARCHAR(30) | PK | Unique identifier for the AP invoice |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| supplierID | VARCHAR(30) | | Supplier master identifier |
| supplierName | VARCHAR(100) | | Supplier legal name |
| supplierType | VARCHAR(30) | | GroundHandler / MRO / Fuel / Catering / Airport / Lessor |
| workOrderID | VARCHAR(30) | FK | References WORK_ORDER.workOrderID (nullable) |
| groundHandlerID | VARCHAR(30) | FK | References GROUND_HANDLER.groundHandlerID (nullable) |
| invoiceNumber | VARCHAR(50) | | Supplier invoice number |
| invoiceDate | DATE | | Invoice date |
| dueDate | DATE | | Payment due date |
| invoiceAmount | DECIMAL(14,2) | | Total invoice amount |
| taxAmount | DECIMAL(10,2) | | Tax amount on invoice |
| currency | CHAR(3) | | ISO 4217 invoice currency |
| paymentStatus | VARCHAR(20) | | Unpaid / PartiallyPaid / Paid / Disputed / Cancelled |
| paymentDate | DATE | | Date payment was made (nullable) |
| glAccountID | VARCHAR(30) | FK | References GL_ACCOUNT.glAccountID — expense account |
| costCentreID | VARCHAR(30) | FK | References COST_CENTRE.costCentreID |

---

### BUDGET
Defines annual and periodic financial budgets per cost centre and GL account.

| Attribute | Type | Key | Description |
|---|---|---|---|
| budgetID | VARCHAR(30) | PK | Unique identifier for the budget record |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| costCentreID | VARCHAR(30) | FK | References COST_CENTRE.costCentreID |
| glAccountID | VARCHAR(30) | FK | References GL_ACCOUNT.glAccountID |
| budgetYear | INTEGER | | Budget fiscal year |
| budgetMonth | INTEGER | | Budget month (1–12, nullable for annual totals) |
| budgetAmount | DECIMAL(14,2) | | Budgeted amount |
| forecastAmount | DECIMAL(14,2) | | Latest forecast amount |
| actualAmount | DECIMAL(14,2) | | Actual amount posted (nullable until period closes) |
| currency | CHAR(3) | | ISO 4217 currency |
| budgetType | VARCHAR(20) | | Original / Revised / Forecast / Actuals |
| approvedBy | VARCHAR(50) | | Name or ID of approver |
| approvedDate | DATE | | Date budget was approved |

---

## Relationships

### Internal — Finance & Interline Domain

| From Entity | To Entity | Cardinality | Description |
|---|---|---|---|
| GL_ACCOUNT | GL_POSTING | 0..* | A GL account has zero or more postings |
| GL_ACCOUNT | ACCOUNTS_PAYABLE_INVOICE | 0..* | A GL account is charged by zero or more invoices |
| GL_ACCOUNT | BUDGET | 0..* | A GL account has zero or more budget lines |
| COST_CENTRE | GL_ACCOUNT | 0..* | A cost centre owns zero or more GL accounts |
| COST_CENTRE | GL_POSTING | 0..* | A cost centre has zero or more postings |
| COST_CENTRE | BUDGET | 0..* | A cost centre has zero or more budget lines |
| COST_CENTRE | COST_CENTRE | 0..* | A cost centre may have child cost centres |
| INTERLINE_AGREEMENT | PRORATION_RULE | 0..* | An agreement has zero or more proration rules |
| INTERLINE_AGREEMENT | INTERLINE_BILLING_COUPON | 0..* | An agreement governs zero or more billing coupons |
| PRORATION_RULE | INTERLINE_REVENUE_PRORATION | 0..* | A proration rule governs zero or more proration records |

### Cross-Domain — Finance → Previous Domains

| From Entity | To Entity | Domain | Description |
|---|---|---|---|
| GL_ACCOUNT | AIRLINE | 01 | Account belongs to airline chart of accounts |
| GL_ACCOUNT | COST_CENTRE | 11 | Account linked to a cost centre |
| COST_CENTRE | AIRLINE | 01 | Cost centre belongs to airline |
| GL_POSTING | REVENUE_RECORD | 06 | Posting driven by passenger revenue record |
| GL_POSTING | CARGO_REVENUE_RECORD | 09 | Posting driven by cargo revenue record |
| INTERLINE_AGREEMENT | AIRLINE | 01 | Agreement held by airline |
| INTERLINE_BILLING_COUPON | TICKET | 04 | Coupon from a ticket |
| INTERLINE_BILLING_COUPON | TICKET_COUPON | 04 | Coupon at segment level |
| INTERLINE_REVENUE_PRORATION | BOOKING | 04 | Proration for a booking |
| INTERLINE_REVENUE_PRORATION | TICKET | 04 | Proration for a ticket |
| INTERLINE_REVENUE_PRORATION | ROUTE | 01 | Proration for a route |
| BSP_SETTLEMENT | AIRLINE | 01 | Settlement by airline with IATA BSP |
| CASS_SETTLEMENT | AIRLINE | 01 | Cargo settlement by airline with IATA CASS |
| ACCOUNTS_PAYABLE_INVOICE | AIRLINE | 01 | Invoice payable by airline |
| ACCOUNTS_PAYABLE_INVOICE | WORK_ORDER | 10 | Invoice for MRO work order |
| ACCOUNTS_PAYABLE_INVOICE | GROUND_HANDLER | 07 | Invoice from ground handler |
| BUDGET | AIRLINE | 01 | Budget set by airline |

---

## Enumerations and Code Lists

### GL_ACCOUNT_TYPE_CD
| Code | Description |
|---|---|
| ASSET | Balance sheet asset account |
| LIABILITY | Balance sheet liability account |
| EQUITY | Shareholders equity account |
| REVENUE | Income / revenue account |
| EXPENSE | Expenditure / cost account |

### INTERLINE_AGREEMENT_TYPE_CD
| Code | Description |
|---|---|
| PRORATE | Standard bilateral prorate agreement |
| INTERLINE | General interline traffic agreement |
| SPA | Special Prorate Agreement |
| MIA | Multilateral Interline Agreement |
| CODESHARE | Commercial codeshare arrangement |

### PRORATION_METHOD_CD
| Code | Description |
|---|---|
| MILEAGE | Pro-rate based on sector mileage |
| EQUAL | Equal split between carriers |
| SECTOR | Fixed sector-based proration |
| PUBLISHED | Published tariff proration |
| IATA_RPDG | IATA Revenue Passenger Division Guide |

### SETTLEMENT_STATUS_CD
| Code | Description |
|---|---|
| PENDING | Awaiting submission to IATA |
| SUBMITTED | Submitted to IATA BSP/CASS |
| SETTLED | Payment settled |
| DISPUTED | Settlement under dispute |

### SUPPLIER_TYPE_CD
| Code | Description |
|---|---|
| GROUND_HANDLER | Airport ground handling provider |
| MRO | Maintenance, Repair & Overhaul provider |
| FUEL | Fuel supplier |
| CATERING | Inflight catering provider |
| AIRPORT | Airport authority / charges |
| LESSOR | Aircraft lessor |

---

## Business Rules

| Rule ID | Entity | Rule Description |
|---|---|---|
| BR-FIN-001 | GL_POSTING | Sum of all Debit postings must equal sum of all Credit postings per journal (double-entry) |
| BR-FIN-002 | BSP_SETTLEMENT | netSettlementAmount must equal totalSales minus totalRefunds minus totalCommissions |
| BR-FIN-003 | CASS_SETTLEMENT | netSettlementAmount must equal totalFreightCharges plus totalSurcharges minus totalCommissions |
| BR-FIN-004 | INTERLINE_REVENUE_PRORATION | homeCarrierProrate plus partnerCarrierProrate must equal totalFareAmount |
| BR-FIN-005 | BUDGET | budgetAmount must be greater than zero |
| BR-FIN-006 | ACCOUNTS_PAYABLE_INVOICE | dueDate must be on or after invoiceDate |
| BR-FIN-007 | GL_ACCOUNT | accountCode must be unique within an airline's chart of accounts |
| BR-FIN-008 | EXCHANGE_RATE | exchangeRate must be greater than zero |
| BR-FIN-009 | INTERLINE_BILLING_COUPON | proratedAmount must be less than or equal to couponFaceValue |
| BR-FIN-010 | COST_CENTRE | A cost centre must not reference itself as parentCostCentreID |

---

## Traceability to AIDM Domains

| AIDM Domain | AIDM Sub-domain | Entities |
|---|---|---|
| Finance | Chart of Accounts | GL_ACCOUNT, COST_CENTRE |
| Finance | General Ledger | GL_POSTING, EXCHANGE_RATE |
| Finance | Interline & Proration | INTERLINE_AGREEMENT, PRORATION_RULE, INTERLINE_BILLING_COUPON, INTERLINE_REVENUE_PRORATION |
| Finance | Settlement | BSP_SETTLEMENT, CASS_SETTLEMENT |
| Finance | Accounts Payable | ACCOUNTS_PAYABLE_INVOICE |
| Finance | Budgeting | BUDGET |

---

## Notes

- General Ledger structure follows IATA Uniform System of Accounts (USA) for Airlines.
- BSP settlement follows IATA Resolution 850 (BSP Operations Procedures Manual).
- CASS settlement follows IATA CASS Operations Procedures Manual.
- Interline billing follows IATA Resolution 780 (Interline Traffic Agreements — Passenger) and Resolution 850m (Interline Billing).
- Revenue proration follows IATA Revenue Accounting Manual (RAM) and IATA RPDG (Revenue Proration Division Guide).
- Exchange rates follow IATA Resolution 024 (IATA Rate of Exchange — IROE) for international settlements.
- ICH (Interline Clearing House) settlement references IATA ICH Operations Manual.
- AP invoice processing follows airline internal procurement and three-way match (PO/GR/Invoice) principles.
