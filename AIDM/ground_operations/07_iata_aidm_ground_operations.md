# IATA AIDM 25.1 — Domain 07: Ground Operations
## Mockup Data Model for AI Consumption

---

## Subject Area Overview

The **Ground Operations** domain in the IATA AIDM covers all airport-based operational activities that support a flight's turnaround — including check-in, boarding, baggage handling, ground support equipment (GSE) allocation, ramp operations, fuelling, catering, cleaning, and slot management. It is the operational bridge between the commercial domains and the airborne Flight Operations domain.

- **AIDM Version**: 25.1
- **Domain**: Operations — Ground Operations
- **Integrates With**: Network & Alliances (01), Product (03), Sales (04), Customer & Loyalty (05), Revenue Management & Pricing (06)

---

## Cross-Domain FK Dependencies

| Referenced Entity | PK Used | Referenced By (Ground Operations) |
|---|---|---|
| AIRLINE | airlineCode | FLIGHT_TURNAROUND, GROUND_HANDLER, AIRPORT_SLOT, CHECK_IN_RECORD |
| ROUTE | routeID | AIRPORT_SLOT |
| BOOKING | bookingID | CHECK_IN_RECORD, BOARDING_PASS |
| BOOKING_SEGMENT | bookingSegmentID | CHECK_IN_RECORD, BOARDING_PASS, BAGGAGE_ITEM |
| TICKET | ticketID | BOARDING_PASS |
| CUSTOMER_PROFILE | customerID | CHECK_IN_RECORD, BOARDING_PASS |
| SEAT | seatID | BOARDING_PASS |
| ANCILLARY_SERVICE | ancillaryServiceID | BAGGAGE_ITEM |
| INVENTORY_CLASS | inventoryClassID | FLIGHT_TURNAROUND |

---

## Entities

### AIRPORT_SLOT
Defines a runway or terminal slot allocated to an airline for a specific flight at an airport.

| Attribute | Type | Key | Description |
|---|---|---|---|
| slotID | VARCHAR(30) | PK | Unique identifier for the airport slot |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| routeID | VARCHAR(30) | FK | References ROUTE.routeID |
| airportCode | VARCHAR(3) | | IATA airport code where slot is held |
| slotType | VARCHAR(10) | | Departure / Arrival |
| scheduledDateTime | TIMESTAMP | | Scheduled slot date and time |
| allocatedDateTime | TIMESTAMP | | Allocated slot time (may differ from scheduled) |
| slotSeason | VARCHAR(5) | | IATA season code (e.g., S25, W25) |
| slotCoordinator | VARCHAR(50) | | Name of slot coordination authority |
| slotStatus | VARCHAR(20) | | Confirmed / Historic / Cancelled / Waived |
| historicUsageRate | DECIMAL(5,4) | | Usage rate for slot retention (must be >= 0.8000 per IATA rules) |
| effectiveDate | DATE | | Slot effective date |
| expiryDate | DATE | | Slot expiry date (nullable) |

---

### GROUND_HANDLER
Defines a ground handling service provider contracted by an airline at an airport.

| Attribute | Type | Key | Description |
|---|---|---|---|
| groundHandlerID | VARCHAR(30) | PK | Unique identifier for the ground handler |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode — contracting airline |
| handlerName | VARCHAR(100) | | Ground handler company name |
| airportCode | VARCHAR(3) | | IATA airport code where handler operates |
| handlerType | VARCHAR(30) | | Selfhandling / ThirdParty / JointVenture |
| servicesProvided | VARCHAR(255) | | Comma-separated list of services (Ramp/Fuelling/Catering/Check-in) |
| isgsMember | CHAR(1) | | Y / N — member of IATA Ground Services standard |
| contractStartDate | DATE | | Contract commencement date |
| contractEndDate | DATE | | Contract end date (nullable) |
| status | VARCHAR(20) | | Active / Suspended / Terminated |

---

### FLIGHT_TURNAROUND
The master operational record for a specific aircraft turnaround at an airport — linking all ground service events.

| Attribute | Type | Key | Description |
|---|---|---|---|
| turnaroundID | VARCHAR(30) | PK | Unique identifier for the turnaround event |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| inventoryClassID | VARCHAR(30) | FK | References INVENTORY_CLASS.inventoryClassID (flight-level link) |
| groundHandlerID | VARCHAR(30) | FK | References GROUND_HANDLER.groundHandlerID |
| flightNumber | VARCHAR(10) | | Operating flight number |
| aircraftRegistration | VARCHAR(10) | | Aircraft tail registration (e.g., G-EUPT) |
| aircraftTypeCode | VARCHAR(10) | | IATA aircraft type code |
| airportCode | VARCHAR(3) | | IATA airport code of turnaround |
| scheduledArrivalDateTime | TIMESTAMP | | Scheduled inbound arrival time |
| actualArrivalDateTime | TIMESTAMP | | Actual inbound arrival time (nullable until arrived) |
| scheduledDepartureDateTime | TIMESTAMP | | Scheduled outbound departure time |
| actualDepartureDateTime | TIMESTAMP | | Actual outbound departure time (nullable until departed) |
| standsGate | VARCHAR(10) | | Gate or stand number assigned |
| turnaroundStatus | VARCHAR(20) | | Planned / InProgress / Completed / Delayed / Diverted |
| targetTurnaroundMinutes | INTEGER | | Planned turnaround time in minutes |
| actualTurnaroundMinutes | INTEGER | | Actual turnaround time in minutes (nullable) |
| delayMinutes | INTEGER | | Total delay in minutes (nullable) |
| primaryDelayCode | VARCHAR(5) | | IATA IRROPS delay code (e.g., 93, 41, 72) |

---

### GROUND_SERVICE_EVENT
Records each individual ground service activity performed during a turnaround (e.g., fuelling, catering, cleaning).

| Attribute | Type | Key | Description |
|---|---|---|---|
| groundServiceEventID | VARCHAR(30) | PK | Unique identifier for the ground service event |
| turnaroundID | VARCHAR(30) | FK | References FLIGHT_TURNAROUND.turnaroundID |
| groundHandlerID | VARCHAR(30) | FK | References GROUND_HANDLER.groundHandlerID |
| serviceType | VARCHAR(30) | | Fuelling / Catering / Cleaning / Deicing / Ramp / Towing / WaterServicing |
| scheduledStartDateTime | TIMESTAMP | | Planned service start time |
| actualStartDateTime | TIMESTAMP | | Actual service start time (nullable) |
| scheduledEndDateTime | TIMESTAMP | | Planned service end time |
| actualEndDateTime | TIMESTAMP | | Actual service end time (nullable) |
| serviceStatus | VARCHAR(20) | | Planned / InProgress / Completed / Skipped / Delayed |
| serviceProviderRef | VARCHAR(50) | | Handler's own service reference number |
| quantity | DECIMAL(10,2) | | Quantity (e.g., litres of fuel, number of meals) |
| unit | VARCHAR(20) | | Litres / Meals / Kg / Units |
| notes | VARCHAR(255) | | Operational notes |

---

### GSE_ALLOCATION
Records the allocation of Ground Support Equipment to a turnaround.

| Attribute | Type | Key | Description |
|---|---|---|---|
| gseAllocationID | VARCHAR(30) | PK | Unique identifier for the GSE allocation |
| turnaroundID | VARCHAR(30) | FK | References FLIGHT_TURNAROUND.turnaroundID |
| equipmentType | VARCHAR(30) | | PushbackTug / BeltLoader / CateringLift / AirStart / GPU / PBB / Stairs |
| equipmentID | VARCHAR(30) | | Equipment fleet identifier |
| scheduledStartDateTime | TIMESTAMP | | Scheduled equipment usage start |
| actualStartDateTime | TIMESTAMP | | Actual equipment usage start (nullable) |
| scheduledEndDateTime | TIMESTAMP | | Scheduled equipment usage end |
| actualEndDateTime | TIMESTAMP | | Actual equipment usage end (nullable) |
| allocationStatus | VARCHAR(20) | | Allocated / InUse / Released / Unavailable |
| operatorID | VARCHAR(30) | | Ground crew operator identifier (nullable) |

---

### CHECK_IN_RECORD
Records a passenger check-in event — covering online, kiosk, and counter check-in channels.

| Attribute | Type | Key | Description |
|---|---|---|---|
| checkInRecordID | VARCHAR(30) | PK | Unique identifier for the check-in record |
| bookingID | VARCHAR(30) | FK | References BOOKING.bookingID |
| bookingSegmentID | VARCHAR(30) | FK | References BOOKING_SEGMENT.bookingSegmentID |
| customerID | VARCHAR(30) | FK | References CUSTOMER_PROFILE.customerID |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| checkInChannel | VARCHAR(20) | | Web / MobileApp / Kiosk / Counter / API |
| checkInDateTime | TIMESTAMP | | Date and time passenger checked in |
| checkInStatus | VARCHAR(20) | | CheckedIn / Standby / Denied / NoShow |
| passengerWeightKg | DECIMAL(5,1) | | Passenger weight in kg (nullable — used in weight & balance) |
| documentsVerified | CHAR(1) | | Y / N — whether travel documents were verified |
| apiTransmitted | CHAR(1) | | Y / N — whether Advance Passenger Info was transmitted |
| boardingSequence | INTEGER | | Boarding sequence number assigned at check-in |

---

### BOARDING_PASS
Represents an issued boarding pass for a passenger on a flight segment.

| Attribute | Type | Key | Description |
|---|---|---|---|
| boardingPassID | VARCHAR(30) | PK | Unique identifier for the boarding pass |
| checkInRecordID | VARCHAR(30) | FK | References CHECK_IN_RECORD.checkInRecordID |
| bookingSegmentID | VARCHAR(30) | FK | References BOOKING_SEGMENT.bookingSegmentID |
| ticketID | VARCHAR(30) | FK | References TICKET.ticketID |
| customerID | VARCHAR(30) | FK | References CUSTOMER_PROFILE.customerID |
| seatID | VARCHAR(30) | FK | References SEAT.seatID |
| boardingGroup | VARCHAR(5) | | Boarding group code (e.g., A1, B2) |
| boardingZone | VARCHAR(10) | | Boarding zone designation |
| barcodeData | VARCHAR(255) | | IATA BCBP (Bar Coded Boarding Pass) data string |
| barcodeFormat | VARCHAR(10) | | QR / Aztec / PDF417 |
| issuanceDateTime | TIMESTAMP | | Date and time boarding pass was issued |
| boardingGate | VARCHAR(10) | | Assigned boarding gate |
| boardingCutoffDateTime | TIMESTAMP | | Latest boarding time |
| boardingPassStatus | VARCHAR(20) | | Issued / Scanned / Cancelled / Upgraded |
| mobilePass | CHAR(1) | | Y / N — whether this is a mobile boarding pass |

---

### BAGGAGE_ITEM
Records an individual checked baggage item associated with a booking segment.

| Attribute | Type | Key | Description |
|---|---|---|---|
| baggageItemID | VARCHAR(30) | PK | Unique identifier for the baggage item |
| bookingSegmentID | VARCHAR(30) | FK | References BOOKING_SEGMENT.bookingSegmentID |
| ancillaryServiceID | VARCHAR(30) | FK | References ANCILLARY_SERVICE.ancillaryServiceID (nullable) |
| baggageTagNumber | VARCHAR(15) | | IATA 10-digit baggage tag number |
| baggageType | VARCHAR(20) | | Checked / Oversize / Special / Sports / Fragile |
| weightKg | DECIMAL(5,1) | | Actual weight in kilograms |
| dimensionsCm | VARCHAR(20) | | L×W×H in centimetres (e.g., 80x50x30) |
| baggageStatus | VARCHAR(20) | | CheckedIn / Loaded / InTransit / Delivered / Lost / Damaged |
| rushBaggage | CHAR(1) | | Y / N — whether bag is rush/priority baggage |
| screeningStatus | VARCHAR(20) | | Cleared / HoldForScreening / Rejected |
| loadPosition | VARCHAR(20) | | Hold position identifier (e.g., FWD-HOLD-1) |
| destinationAirportCode | VARCHAR(3) | | IATA final destination airport code |
| throughBaggage | CHAR(1) | | Y / N — whether bag is checked through to final destination |

---

### BAGGAGE_RECONCILIATION
Records the reconciliation of loaded baggage against boarded passengers (IATA BSM/BTM standard).

| Attribute | Type | Key | Description |
|---|---|---|---|
| reconciliationID | VARCHAR(30) | PK | Unique identifier for the baggage reconciliation record |
| turnaroundID | VARCHAR(30) | FK | References FLIGHT_TURNAROUND.turnaroundID |
| baggageItemID | VARCHAR(30) | FK | References BAGGAGE_ITEM.baggageItemID |
| boardingPassID | VARCHAR(30) | FK | References BOARDING_PASS.boardingPassID |
| reconciliationStatus | VARCHAR(20) | | Matched / Unmatched / OffloadRequired / Offloaded |
| scanDateTime | TIMESTAMP | | Time baggage was scanned for reconciliation |
| scanLocation | VARCHAR(30) | | Location of scan (BaggageMakeup / HoldDoor / BeltLoader) |
| offloadReason | VARCHAR(50) | | Reason for offload if applicable (nullable) |

---

## Relationships

### Internal — Ground Operations Domain

| From Entity | To Entity | Cardinality | Description |
|---|---|---|---|
| FLIGHT_TURNAROUND | GROUND_SERVICE_EVENT | 1..* | A turnaround has one or more service events |
| FLIGHT_TURNAROUND | GSE_ALLOCATION | 0..* | A turnaround has zero or more GSE allocations |
| FLIGHT_TURNAROUND | BAGGAGE_RECONCILIATION | 0..* | A turnaround has zero or more baggage reconciliation records |
| CHECK_IN_RECORD | BOARDING_PASS | 1..1 | A check-in record produces one boarding pass |
| BOARDING_PASS | BAGGAGE_RECONCILIATION | 0..* | A boarding pass may have zero or more baggage reconciliation events |
| BAGGAGE_ITEM | BAGGAGE_RECONCILIATION | 1..1 | A baggage item has one reconciliation record |

### Cross-Domain — Ground Operations → Previous Domains

| From Entity | To Entity | Domain | Description |
|---|---|---|---|
| AIRPORT_SLOT | AIRLINE | 01 | Slot held by an airline |
| AIRPORT_SLOT | ROUTE | 01 | Slot associated with a route |
| GROUND_HANDLER | AIRLINE | 01 | Handler contracted by airline |
| FLIGHT_TURNAROUND | AIRLINE | 01 | Turnaround operated by airline |
| FLIGHT_TURNAROUND | GROUND_HANDLER | 07 | Turnaround handled by a ground handler |
| FLIGHT_TURNAROUND | INVENTORY_CLASS | 06 | Turnaround linked to flight inventory |
| GROUND_SERVICE_EVENT | GROUND_HANDLER | 07 | Service event performed by handler |
| CHECK_IN_RECORD | BOOKING | 04 | Check-in linked to a booking |
| CHECK_IN_RECORD | BOOKING_SEGMENT | 04 | Check-in for a specific segment |
| CHECK_IN_RECORD | CUSTOMER_PROFILE | 05 | Check-in associated with a customer |
| CHECK_IN_RECORD | AIRLINE | 01 | Check-in at an airline's desk |
| BOARDING_PASS | BOOKING_SEGMENT | 04 | Boarding pass for a segment |
| BOARDING_PASS | TICKET | 04 | Boarding pass linked to ticket |
| BOARDING_PASS | CUSTOMER_PROFILE | 05 | Boarding pass issued to a customer |
| BOARDING_PASS | SEAT | 03 | Boarding pass assigned to a seat |
| BAGGAGE_ITEM | BOOKING_SEGMENT | 04 | Baggage belongs to a segment |
| BAGGAGE_ITEM | ANCILLARY_SERVICE | 03 | Baggage may reference an ancillary (extra bag) |

---

## Enumerations and Code Lists

### SLOT_STATUS_CD
| Code | Description |
|---|---|
| CONFIRMED | Slot confirmed by coordinator |
| HISTORIC | Historic slot with usage rights |
| CANCELLED | Slot cancelled |
| WAIVED | Slot waived without penalty |

### TURNAROUND_STATUS_CD
| Code | Description |
|---|---|
| PLANNED | Turnaround planned, not yet started |
| IN_PROGRESS | Aircraft on stand, services underway |
| COMPLETED | All services complete, aircraft departed |
| DELAYED | Turnaround has exceeded planned time |
| DIVERTED | Flight diverted — turnaround voided |

### GROUND_SERVICE_TYPE_CD
| Code | Description |
|---|---|
| FUELLING | Aircraft refuelling |
| CATERING | Meal and galley servicing |
| CLEANING | Cabin cleaning |
| DEICING | Anti-icing / de-icing operations |
| RAMP | Ramp and baggage loading/offloading |
| TOWING | Aircraft pushback / towing |
| WATER | Potable water servicing |

### GSE_EQUIPMENT_TYPE_CD
| Code | Description |
|---|---|
| PUSHBACK_TUG | Pushback tractor |
| BELT_LOADER | Baggage belt loader |
| CATERING_LIFT | High-loader for catering |
| AIR_START | Aircraft engine air start unit |
| GPU | Ground Power Unit |
| PBB | Passenger Boarding Bridge (jetway) |
| STAIRS | Passenger boarding stairs |

### CHECK_IN_CHANNEL_CD
| Code | Description |
|---|---|
| WEB | Online web check-in |
| MOBILE_APP | Airline mobile app |
| KIOSK | Self-service kiosk at airport |
| COUNTER | Airport check-in counter |
| API | Third-party API integration |

### BAGGAGE_STATUS_CD
| Code | Description |
|---|---|
| CHECKED_IN | Bag accepted at check-in |
| LOADED | Bag loaded onto aircraft |
| IN_TRANSIT | Bag in transit between connections |
| DELIVERED | Bag delivered to passenger at destination |
| LOST | Bag not located |
| DAMAGED | Bag received in damaged condition |

### IATA_DELAY_CODE_CD (Partial — IATA Standard AHM 780)
| Code | Description |
|---|---|
| 11 | Late check-in — congestion |
| 21 | Documentation, passenger processing |
| 41 | Aircraft rotation — late arrival |
| 72 | Late fuelling / fuel supplier |
| 93 | ATC restrictions en-route |

---

## Business Rules

| Rule ID | Entity | Rule Description |
|---|---|---|
| BR-GRD-001 | AIRPORT_SLOT | historicUsageRate must be >= 0.8000 to retain slot under IATA Worldwide Slot Guidelines |
| BR-GRD-002 | FLIGHT_TURNAROUND | actualTurnaroundMinutes must be populated before turnaroundStatus = 'Completed' |
| BR-GRD-003 | BAGGAGE_RECONCILIATION | reconciliationStatus must be 'Matched' for all BAGGAGE_ITEM records before flight departure |
| BR-GRD-004 | BOARDING_PASS | barcodeData must conform to IATA BCBP Resolution 792 format |
| BR-GRD-005 | CHECK_IN_RECORD | checkInDateTime must be before the flight's scheduledDepartureDateTime |
| BR-GRD-006 | BAGGAGE_ITEM | baggageTagNumber must be exactly 10 numeric digits per IATA Resolution 740 |
| BR-GRD-007 | FLIGHT_TURNAROUND | actualDepartureDateTime must be after actualArrivalDateTime |
| BR-GRD-008 | GSE_ALLOCATION | actualEndDateTime must be after actualStartDateTime when both are populated |
| BR-GRD-009 | GROUND_SERVICE_EVENT | actualEndDateTime must be after actualStartDateTime when both are populated |
| BR-GRD-010 | BOARDING_PASS | boardingCutoffDateTime must be before the flight's scheduledDepartureDateTime |

---

## Traceability to AIDM Domains

| AIDM Domain | AIDM Sub-domain | Entities |
|---|---|---|
| Operations | Airport Slot Management | AIRPORT_SLOT |
| Operations | Ground Handling | GROUND_HANDLER, GROUND_SERVICE_EVENT, GSE_ALLOCATION |
| Operations | Aircraft Turnaround | FLIGHT_TURNAROUND |
| Operations | Passenger Processing | CHECK_IN_RECORD, BOARDING_PASS |
| Operations | Baggage Management | BAGGAGE_ITEM, BAGGAGE_RECONCILIATION |

---

## Notes

- Slot management follows IATA Worldwide Slot Guidelines (WSG) — 80/20 grandfather rule for slot retention.
- Baggage tag numbering follows IATA Resolution 740 (10-digit licence plate tag).
- Boarding pass barcode format follows IATA Resolution 792 (BCBP — Bar Coded Boarding Pass).
- Baggage reconciliation follows IATA Resolution 753 (Baggage Tracking) and IATA BSM/BTM message standards.
- Advance Passenger Information (API) transmission flag aligns with ICAO Annex 9 (Facilitation) requirements.
- Delay codes follow IATA AHM 780 (IATA Standard Delay Codes).
- Ground handling service standards reference IATA IGOM (IATA Ground Operations Manual).
