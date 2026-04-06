# IATA AIDM 25.1 — Domain 08: Flight Operations
## Mockup Data Model for AI Consumption

---

## Subject Area Overview

The **Flight Operations** domain in the IATA AIDM covers all airborne and pre-flight operational management activities — including flight planning, crew rostering and qualification, aircraft fleet management, weight & balance, fuel management, operational flight plans (OFP), Air Traffic Control (ATC) communication, ACARS/OOOI event tracking, and flight disruption/IRROPS management. It is the operational core of the airline enterprise.

- **AIDM Version**: 25.1
- **Domain**: Operations — Flight Operations
- **Integrates With**: Network & Alliances (01), Product (03), Sales (04), Revenue Management & Pricing (06), Ground Operations (07)

---

## Cross-Domain FK Dependencies

| Referenced Entity | PK Used | Referenced By (Flight Operations) |
|---|---|---|
| AIRLINE | airlineCode | AIRCRAFT, FLIGHT, CREW_MEMBER, OPERATIONAL_FLIGHT_PLAN |
| ROUTE | routeID | FLIGHT |
| AIRPORT_SLOT | slotID | FLIGHT |
| FLIGHT_TURNAROUND | turnaroundID | FLIGHT |
| INVENTORY_CLASS | inventoryClassID | FLIGHT |
| BOOKING_SEGMENT | bookingSegmentID | FLIGHT_PAX_MANIFEST |

---

## Entities

### AIRCRAFT
The master record for an individual aircraft asset in an airline's operating fleet.

| Attribute | Type | Key | Description |
|---|---|---|---|
| aircraftID | VARCHAR(30) | PK | Unique identifier for the aircraft |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode — operating airline |
| registration | VARCHAR(10) | | Tail registration (e.g., G-EUPT, N12345) |
| aircraftTypeCode | VARCHAR(10) | | IATA aircraft type code (e.g., 789, 77W, 320) |
| icaoTypeCode | VARCHAR(4) | | ICAO aircraft type designator (e.g., B789, B77W, A320) |
| manufacturer | VARCHAR(50) | | Airframe manufacturer (Boeing / Airbus / Embraer) |
| modelSeries | VARCHAR(30) | | Model series (e.g., 737-800, A321neo) |
| engineType | VARCHAR(50) | | Engine model (e.g., GE90-115B, LEAP-1A) |
| engineCount | INTEGER | | Number of engines |
| maxTakeoffWeightKg | DECIMAL(10,2) | | Maximum Takeoff Weight (MTOW) in kg |
| operatingEmptyWeightKg | DECIMAL(10,2) | | Operating Empty Weight (OEW) in kg |
| maxFuelCapacityLitres | DECIMAL(10,2) | | Maximum fuel capacity in litres |
| maxRangeKm | INTEGER | | Maximum range in kilometres |
| deliveryDate | DATE | | Date aircraft was delivered to airline |
| aircraftAge | DECIMAL(5,2) | | Aircraft age in years (computed) |
| ownershipType | VARCHAR(20) | | Owned / Leased / ACMI / Wet-Lease |
| lessorName | VARCHAR(100) | | Lessor name if leased (nullable) |
| status | VARCHAR(20) | | Active / Grounded / Storage / Retired / AOG |
| baseAirportCode | VARCHAR(3) | | IATA base station airport code |

---

### AIRCRAFT_UTILISATION
Tracks cumulative flight hours and cycles for an individual aircraft — feeds into maintenance scheduling.

| Attribute | Type | Key | Description |
|---|---|---|---|
| utilisationID | VARCHAR(30) | PK | Unique identifier for the utilisation record |
| aircraftID | VARCHAR(30) | FK | References AIRCRAFT.aircraftID |
| recordDate | DATE | | Date of the utilisation snapshot |
| totalFlightHours | DECIMAL(10,2) | | Cumulative total airframe flight hours |
| totalCycles | INTEGER | | Cumulative total flight cycles (landings) |
| flightHoursThisMonth | DECIMAL(8,2) | | Flight hours in current calendar month |
| cyclesThisMonth | INTEGER | | Cycles in current calendar month |
| dailyUtilisationHours | DECIMAL(5,2) | | Average daily utilisation hours |
| blockHoursYTD | DECIMAL(10,2) | | Block hours year-to-date |
| lastUpdatedDateTime | TIMESTAMP | | Timestamp of last utilisation update |

---

### CREW_MEMBER
Master record for a crew member — both flight deck (pilots) and cabin crew.

| Attribute | Type | Key | Description |
|---|---|---|---|
| crewMemberID | VARCHAR(30) | PK | Unique identifier for the crew member |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| employeeID | VARCHAR(20) | | HR employee number |
| firstName | VARCHAR(50) | | Crew member first name |
| lastName | VARCHAR(50) | | Crew member last name |
| crewType | VARCHAR(20) | | FlightDeck / CabinCrew |
| crewRank | VARCHAR(30) | | Captain / FirstOfficer / SeniorCabinCrew / CabinCrew / Purser |
| baseAirportCode | VARCHAR(3) | | IATA crew base airport code |
| licenceNumber | VARCHAR(30) | | Pilot licence number (ATPL/CPL) or cabin crew certificate number |
| licenceIssuingAuthority | VARCHAR(10) | | ICAO state of licence issue (e.g., EASA, FAA, CAAC) |
| licenceExpiryDate | DATE | | Licence expiry date |
| medicalCertExpiry | DATE | | Aviation medical certificate expiry date |
| nationality | VARCHAR(3) | | ISO 3166-1 alpha-3 nationality |
| status | VARCHAR(20) | | Active / Suspended / OnLeave / Retired |
| hireDate | DATE | | Date crew member was hired |

---

### CREW_QUALIFICATION
Records type ratings, endorsements, and route qualifications for a crew member.

| Attribute | Type | Key | Description |
|---|---|---|---|
| qualificationID | VARCHAR(30) | PK | Unique identifier for the qualification |
| crewMemberID | VARCHAR(30) | FK | References CREW_MEMBER.crewMemberID |
| qualificationType | VARCHAR(30) | | TypeRating / RouteCheck / LineCheck / SEP / DGR / CRM |
| aircraftTypeCode | VARCHAR(10) | | IATA aircraft type code qualification applies to (nullable) |
| issuingAuthority | VARCHAR(10) | | ICAO regulatory authority (e.g., EASA, FAA) |
| issueDate | DATE | | Qualification issue date |
| expiryDate | DATE | | Qualification expiry date |
| qualificationStatus | VARCHAR(20) | | Valid / Expired / Suspended |
| lastProficiencyCheckDate | DATE | | Date of last simulator or line check |
| nextProficiencyCheckDate | DATE | | Date of next required check |

---

### FLIGHT
The core operational record for a scheduled or ad-hoc flight operation.

| Attribute | Type | Key | Description |
|---|---|---|---|
| flightID | VARCHAR(30) | PK | Unique identifier for the flight |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| routeID | VARCHAR(30) | FK | References ROUTE.routeID |
| slotID | VARCHAR(30) | FK | References AIRPORT_SLOT.slotID (nullable) |
| turnaroundID | VARCHAR(30) | FK | References FLIGHT_TURNAROUND.turnaroundID (nullable) |
| inventoryClassID | VARCHAR(30) | FK | References INVENTORY_CLASS.inventoryClassID |
| aircraftID | VARCHAR(30) | FK | References AIRCRAFT.aircraftID |
| flightNumber | VARCHAR(10) | | IATA/ICAO flight number (e.g., BA0001) |
| operatingCarrierCode | VARCHAR(3) | | IATA code of operating carrier |
| marketingCarrierCode | VARCHAR(3) | | IATA code of marketing carrier (if codeshare) |
| originAirportCode | VARCHAR(3) | | IATA departure airport code |
| destinationAirportCode | VARCHAR(3) | | IATA arrival airport code |
| scheduledDepartureDateTime | TIMESTAMP | | Scheduled off-block time (SOBT) |
| scheduledArrivalDateTime | TIMESTAMP | | Scheduled in-block time (SIBT) |
| estimatedDepartureDateTime | TIMESTAMP | | Estimated off-block time (EOBT) |
| estimatedArrivalDateTime | TIMESTAMP | | Estimated in-block time (EIBT) |
| actualDepartureDateTime | TIMESTAMP | | Actual off-block time (AOBT) |
| actualArrivalDateTime | TIMESTAMP | | Actual in-block time (AIBT) |
| flightStatus | VARCHAR(20) | | Scheduled / Boarding / Departed / Airborne / Landed / Cancelled / Diverted |
| codeShareFlight | CHAR(1) | | Y / N — whether flight is a codeshare |
| flightRules | VARCHAR(5) | | IFR / VFR |
| flightType | VARCHAR(10) | | Scheduled / Charter / Ferry / Training |
| totalPaxBoarded | INTEGER | | Total passengers boarded |
| totalBaggageKg | DECIMAL(8,2) | | Total checked baggage weight in kg |
| totalCargoKg | DECIMAL(8,2) | | Total cargo weight in kg |

---

### OPERATIONAL_FLIGHT_PLAN
The filed operational flight plan for a specific flight — containing routing, fuel, and weather data.

| Attribute | Type | Key | Description |
|---|---|---|---|
| ofpID | VARCHAR(30) | PK | Unique identifier for the OFP |
| flightID | VARCHAR(30) | FK | References FLIGHT.flightID |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| planningSystem | VARCHAR(50) | | System used to generate OFP (e.g., Lido, Jeppesen, SITA) |
| ofpVersion | INTEGER | | Version number (increments with revisions) |
| filedRoute | TEXT | | Filed ATC route string |
| plannedAltitudeFt | INTEGER | | Initial planned cruise altitude in feet |
| plannedTASKts | INTEGER | | Planned True Air Speed in knots |
| plannedFlightTimeMins | INTEGER | | Planned block-to-block flight time in minutes |
| tripFuelKg | DECIMAL(8,2) | | Planned trip fuel in kg |
| contingencyFuelKg | DECIMAL(8,2) | | Contingency fuel (5% or fixed) in kg |
| alternateFuelKg | DECIMAL(8,2) | | Fuel to alternate airport in kg |
| finalReserveFuelKg | DECIMAL(8,2) | | Final reserve fuel (30/45 min) in kg |
| additionalFuelKg | DECIMAL(8,2) | | Additional discretionary fuel in kg |
| totalFuelOnBoardKg | DECIMAL(8,2) | | Total fuel on board at departure in kg |
| takeoffWeightKg | DECIMAL(10,2) | | Planned takeoff weight in kg |
| landingWeightKg | DECIMAL(10,2) | | Planned landing weight in kg |
| alternateAirportCode | VARCHAR(3) | | IATA alternate destination airport code |
| atsFiledDateTime | TIMESTAMP | | Date and time OFP was filed with ATS |
| ofpStatus | VARCHAR(20) | | Draft / Filed / Active / Superseded / Closed |

---

### CREW_ASSIGNMENT
Records the assignment of a crew member to a specific flight.

| Attribute | Type | Key | Description |
|---|---|---|---|
| crewAssignmentID | VARCHAR(30) | PK | Unique identifier for the crew assignment |
| flightID | VARCHAR(30) | FK | References FLIGHT.flightID |
| crewMemberID | VARCHAR(30) | FK | References CREW_MEMBER.crewMemberID |
| crewPosition | VARCHAR(30) | | Captain / FirstOfficer / CabinCrewLead / CabinCrew / Purser |
| assignmentType | VARCHAR(20) | | Rostered / Standby / Deadhead / Training |
| dutyStartDateTime | TIMESTAMP | | Crew duty period start (FDP start) |
| dutyEndDateTime | TIMESTAMP | | Crew duty period end |
| flightDutyPeriodMins | INTEGER | | Total flight duty period in minutes |
| restPeriodBeforeMins | INTEGER | | Rest period before duty in minutes |
| assignmentStatus | VARCHAR(20) | | Confirmed / Standby / Replaced / Cancelled |
| fdpLimitMins | INTEGER | | Maximum permitted FDP per regulations (EASA/FAA) |
| cumulativeFlightHrs28Day | DECIMAL(6,2) | | Cumulative flight hours in last 28 days |
| cumulativeFlightHrsYear | DECIMAL(7,2) | | Cumulative flight hours in last 365 days |

---

### FLIGHT_OOOI_EVENT
Records the four key OOOI gate/air events for a flight (Out, Off, On, In) transmitted via ACARS.

| Attribute | Type | Key | Description |
|---|---|---|---|
| oooi_EventID | VARCHAR(30) | PK | Unique identifier for the OOOI event |
| flightID | VARCHAR(30) | FK | References FLIGHT.flightID |
| eventType | VARCHAR(5) | | OUT / OFF / ON / IN |
| eventDateTime | TIMESTAMP | | Timestamp of the OOOI event |
| airportCode | VARCHAR(3) | | IATA airport code where event occurred |
| fuelOnBoardKg | DECIMAL(8,2) | | Fuel on board at event time (nullable) |
| acarsMessageID | VARCHAR(30) | | ACARS message identifier |
| transmissionSource | VARCHAR(20) | | ACARS / Manual / ADS-B / ATIS |
| latitude | DECIMAL(9,6) | | Latitude at event time (nullable) |
| longitude | DECIMAL(9,6) | | Longitude at event time (nullable) |

---

### FLIGHT_DISRUPTION
Records an irregular operation (IRROPS) event affecting a flight — cancellation, diversion, delay, or swap.

| Attribute | Type | Key | Description |
|---|---|---|---|
| disruptionID | VARCHAR(30) | PK | Unique identifier for the disruption record |
| flightID | VARCHAR(30) | FK | References FLIGHT.flightID |
| disruptionType | VARCHAR(20) | | Cancellation / Diversion / Delay / AircraftSwap / CrewChange |
| disruptionCause | VARCHAR(30) | | Weather / ATC / Technical / Operational / Security / Medical |
| primaryDelayCode | VARCHAR(5) | | IATA AHM 780 delay code |
| secondaryDelayCode | VARCHAR(5) | | Secondary delay code (nullable) |
| delayMinutes | INTEGER | | Total delay in minutes (nullable) |
| diversionAirportCode | VARCHAR(3) | | IATA diversion airport code (nullable) |
| cancellationDateTime | TIMESTAMP | | Cancellation timestamp (nullable) |
| disruptionDescription | VARCHAR(255) | | Free-text description of disruption event |
| passengerImpactCount | INTEGER | | Number of passengers impacted |
| recoveryFlightID | VARCHAR(30) | FK | References FLIGHT.flightID — recovery flight (nullable) |
| reportedDateTime | TIMESTAMP | | Timestamp disruption was reported |
| resolvedDateTime | TIMESTAMP | | Timestamp disruption was resolved (nullable) |

---

### WEIGHT_AND_BALANCE
Records the actual weight and balance calculation performed before departure.

| Attribute | Type | Key | Description |
|---|---|---|---|
| wbRecordID | VARCHAR(30) | PK | Unique identifier for the W&B record |
| flightID | VARCHAR(30) | FK | References FLIGHT.flightID |
| aircraftID | VARCHAR(30) | FK | References AIRCRAFT.aircraftID |
| dryOperatingWeightKg | DECIMAL(10,2) | | Dry operating weight including crew and equipment |
| passengerWeightKg | DECIMAL(10,2) | | Total passenger weight (standard or actual) |
| baggageWeightKg | DECIMAL(10,2) | | Total checked baggage weight |
| cargoWeightKg | DECIMAL(10,2) | | Total cargo and mail weight |
| fuelWeightKg | DECIMAL(10,2) | | Total fuel weight at takeoff |
| actualTakeoffWeightKg | DECIMAL(10,2) | | Actual takeoff weight |
| actualLandingWeightKg | DECIMAL(10,2) | | Actual landing weight |
| maximumTakeoffWeightKg | DECIMAL(10,2) | | Structural MTOW limit |
| maximumLandingWeightKg | DECIMAL(10,2) | | Structural MLW limit |
| centreOfGravityPct | DECIMAL(5,2) | | CG position as % of Mean Aerodynamic Chord (%MAC) |
| cgWithinLimits | CHAR(1) | | Y / N — whether CG is within certified limits |
| preparedByCrewID | VARCHAR(30) | FK | References CREW_MEMBER.crewMemberID — dispatcher/crew |
| preparedDateTime | TIMESTAMP | | Timestamp W&B was finalised |

---

## Relationships

### Internal — Flight Operations Domain

| From Entity | To Entity | Cardinality | Description |
|---|---|---|---|
| AIRCRAFT | AIRCRAFT_UTILISATION | 0..* | An aircraft has zero or more utilisation snapshots |
| AIRCRAFT | FLIGHT | 0..* | An aircraft operates zero or more flights |
| CREW_MEMBER | CREW_QUALIFICATION | 1..* | A crew member holds one or more qualifications |
| CREW_MEMBER | CREW_ASSIGNMENT | 0..* | A crew member has zero or more assignments |
| FLIGHT | OPERATIONAL_FLIGHT_PLAN | 1..* | A flight has one or more OFP versions |
| FLIGHT | CREW_ASSIGNMENT | 1..* | A flight has one or more crew assignments |
| FLIGHT | FLIGHT_OOOI_EVENT | 0..4 | A flight has zero to four OOOI events |
| FLIGHT | FLIGHT_DISRUPTION | 0..* | A flight may have zero or more disruption records |
| FLIGHT | WEIGHT_AND_BALANCE | 0..1 | A flight has zero or one W&B record |

### Cross-Domain — Flight Operations → Previous Domains

| From Entity | To Entity | Domain | Description |
|---|---|---|---|
| AIRCRAFT | AIRLINE | 01 | Aircraft owned/operated by airline |
| CREW_MEMBER | AIRLINE | 01 | Crew member employed by airline |
| FLIGHT | AIRLINE | 01 | Flight operated by airline |
| FLIGHT | ROUTE | 01 | Flight operates on a route |
| FLIGHT | AIRPORT_SLOT | 07 | Flight uses an airport slot |
| FLIGHT | FLIGHT_TURNAROUND | 07 | Flight associated with a ground turnaround |
| FLIGHT | INVENTORY_CLASS | 06 | Flight linked to RM inventory |
| FLIGHT | AIRCRAFT | 08 | Flight operated by a specific aircraft |
| FLIGHT_DISRUPTION | FLIGHT | 08 | Disruption affects a flight |
| WEIGHT_AND_BALANCE | CREW_MEMBER | 08 | W&B prepared by dispatcher/crew |

---

## Enumerations and Code Lists

### AIRCRAFT_STATUS_CD
| Code | Description |
|---|---|
| ACTIVE | Aircraft in active revenue service |
| GROUNDED | Temporarily grounded — technical/regulatory |
| STORAGE | Aircraft in long-term storage |
| RETIRED | Aircraft permanently withdrawn from service |
| AOG | Aircraft on Ground — awaiting critical repair |

### FLIGHT_STATUS_CD
| Code | Description |
|---|---|
| SCHEDULED | Flight is scheduled |
| BOARDING | Boarding in progress |
| DEPARTED | Flight has pushed back |
| AIRBORNE | Flight is airborne |
| LANDED | Flight has landed |
| CANCELLED | Flight has been cancelled |
| DIVERTED | Flight has been diverted to alternate airport |

### OOOI_EVENT_TYPE_CD (IATA Standard)
| Code | Description |
|---|---|
| OUT | Aircraft off-blocks (pushback) |
| OFF | Aircraft airborne (wheels off) |
| ON | Aircraft touchdown (wheels on) |
| IN | Aircraft on-blocks (engines off) |

### CREW_RANK_CD
| Code | Description |
|---|---|
| CAPTAIN | Aircraft Commander (PIC) |
| FIRST_OFFICER | Co-pilot (SIC) |
| SENIOR_CABIN_CREW | Senior cabin crew / In-charge |
| CABIN_CREW | Standard cabin crew |
| PURSER | Purser (long-haul senior crew) |

### DISRUPTION_TYPE_CD
| Code | Description |
|---|---|
| CANCELLATION | Flight cancelled entirely |
| DIVERSION | Flight diverted to alternate airport |
| DELAY | Flight delayed beyond scheduled time |
| AIRCRAFT_SWAP | Aircraft substituted with different type/registration |
| CREW_CHANGE | Crew replacement due to availability/compliance |

### OWNERSHIP_TYPE_CD
| Code | Description |
|---|---|
| OWNED | Airline-owned asset |
| LEASED | Operating lease from lessor |
| ACMI | Aircraft, Crew, Maintenance & Insurance lease |
| WET_LEASE | Wet lease including crew |

---

## Business Rules

| Rule ID | Entity | Rule Description |
|---|---|---|
| BR-FLT-001 | FLIGHT | actualDepartureDateTime must be after scheduledDepartureDateTime for delayed flights |
| BR-FLT-002 | CREW_ASSIGNMENT | flightDutyPeriodMins must not exceed fdpLimitMins per applicable regulations |
| BR-FLT-003 | CREW_ASSIGNMENT | restPeriodBeforeMins must meet minimum rest requirements (EASA ORO.FTL: 720 mins minimum) |
| BR-FLT-004 | WEIGHT_AND_BALANCE | actualTakeoffWeightKg must not exceed maximumTakeoffWeightKg |
| BR-FLT-005 | WEIGHT_AND_BALANCE | cgWithinLimits must be Y before flight departure is permitted |
| BR-FLT-006 | OPERATIONAL_FLIGHT_PLAN | totalFuelOnBoardKg must equal tripFuelKg + contingencyFuelKg + alternateFuelKg + finalReserveFuelKg + additionalFuelKg |
| BR-FLT-007 | CREW_QUALIFICATION | qualificationStatus must be 'Valid' for all assigned crew before CREW_ASSIGNMENT.assignmentStatus = 'Confirmed' |
| BR-FLT-008 | CREW_QUALIFICATION | licenceExpiryDate must be in the future for active crew |
| BR-FLT-009 | AIRCRAFT_UTILISATION | totalCycles must be a non-negative integer |
| BR-FLT-010 | FLIGHT_OOOI_EVENT | eventType sequence must be OUT → OFF → ON → IN for a complete flight |
| BR-FLT-011 | FLIGHT | totalPaxBoarded must not exceed INVENTORY_CLASS.authorisedCapacity |

---

## Traceability to AIDM Domains

| AIDM Domain | AIDM Sub-domain | Entities |
|---|---|---|
| Operations | Fleet Management | AIRCRAFT, AIRCRAFT_UTILISATION |
| Operations | Crew Management | CREW_MEMBER, CREW_QUALIFICATION, CREW_ASSIGNMENT |
| Operations | Flight Scheduling | FLIGHT |
| Operations | Flight Planning | OPERATIONAL_FLIGHT_PLAN, WEIGHT_AND_BALANCE |
| Operations | Flight Tracking | FLIGHT_OOOI_EVENT |
| Operations | IRROPS Management | FLIGHT_DISRUPTION |

---

## Notes

- OOOI event tracking follows IATA AHM (Airport Handling Manual) Chapter 9 and ACARS OOOI message standards.
- Flight Duty Period (FDP) limits follow EASA ORO.FTL Subpart Q and FAA Part 117 fatigue risk management rules.
- Fuel planning follows EASA SPA.FUEL / EU-OPS 1.255 minimum fuel requirements.
- Weight & Balance calculation follows EASA CS-25 and manufacturer Airplane Flight Manual (AFM) limitations.
- OFP routing follows ICAO Doc 4444 (PANS-ATM) flight plan format.
- Aircraft type codes follow IATA DOC 8643 aircraft type designators.
- Delay codes follow IATA AHM 780 standard delay code list.
- Crew qualification records support ICAO Annex 1 (Personnel Licensing) compliance.
