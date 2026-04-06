# IATA AIDM 25.1 — Domain 10: Maintenance, Repair & Overhaul (MRO)
## Mockup Data Model for AI Consumption

---

## Subject Area Overview

The **Maintenance, Repair & Overhaul (MRO)** domain in the IATA AIDM covers the full lifecycle of aircraft airworthiness management — including maintenance programme definition, work order management, component tracking, defect recording, airworthiness directives (ADs), service bulletin compliance, approved maintenance organisation (AMO) management, parts inventory, and Certificate of Release to Service (CRS) management. It is the safety and regulatory backbone of airline technical operations.

- **AIDM Version**: 25.1
- **Domain**: Operations — Maintenance, Repair & Overhaul
- **Integrates With**: Network & Alliances (01), Flight Operations (08)

---

## Cross-Domain FK Dependencies

| Referenced Entity | PK Used | Referenced By (MRO) |
|---|---|---|
| AIRLINE | airlineCode | MAINTENANCE_PROGRAMME, WORK_ORDER, PARTS_INVENTORY, AMO |
| AIRCRAFT | aircraftID | AIRCRAFT_MAINTENANCE_RECORD, WORK_ORDER, DEFECT_RECORD, COMPONENT |
| AIRCRAFT_UTILISATION | utilisationID | MAINTENANCE_TASK |
| CREW_MEMBER | crewMemberID | WORK_ORDER (certifying engineer link) |
| FLIGHT | flightID | DEFECT_RECORD |

---

## Entities

### MAINTENANCE_PROGRAMME
Defines the approved maintenance programme for an aircraft type — containing all scheduled task intervals.

| Attribute | Type | Key | Description |
|---|---|---|---|
| maintenanceProgrammeID | VARCHAR(30) | PK | Unique identifier for the maintenance programme |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| aircraftTypeCode | VARCHAR(10) | | IATA aircraft type code programme applies to |
| programmeName | VARCHAR(100) | | Programme name (e.g., BA B787 AMP Rev 15) |
| programmeRevision | VARCHAR(10) | | Revision number or letter |
| approvalAuthority | VARCHAR(10) | | Regulatory authority approving programme (EASA/FAA/CAAC) |
| approvalReference | VARCHAR(50) | | Regulatory approval reference number |
| approvalDate | DATE | | Date programme was approved |
| revisionDate | DATE | | Date of latest revision |
| status | VARCHAR(20) | | Active / Superseded / Draft |

---

### MAINTENANCE_TASK
Defines an individual scheduled maintenance task within a maintenance programme.

| Attribute | Type | Key | Description |
|---|---|---|---|
| maintenanceTaskID | VARCHAR(30) | PK | Unique identifier for the maintenance task |
| maintenanceProgrammeID | VARCHAR(30) | FK | References MAINTENANCE_PROGRAMME.maintenanceProgrammeID |
| taskCardNumber | VARCHAR(30) | | Task card / MPD reference number |
| taskDescription | VARCHAR(255) | | Description of the maintenance task |
| taskType | VARCHAR(30) | | Check / Inspection / Lubrication / Replacement / Overhaul / Functional |
| checkLevel | VARCHAR(5) | | A / B / C / D / Line / Transit / PreFlight |
| intervalFH | DECIMAL(8,2) | | Flight hours interval (nullable) |
| intervalFC | INTEGER | | Flight cycles interval (nullable) |
| intervalDays | INTEGER | | Calendar days interval (nullable) |
| tolerancePct | DECIMAL(5,2) | | Permitted tolerance as % of interval |
| estimatedManHours | DECIMAL(6,2) | | Estimated man-hours to complete task |
| zoneCode | VARCHAR(10) | | Aircraft zone code (ATA chapter/zone) |
| ataChapter | VARCHAR(5) | | ATA 100 chapter code |
| accessRequired | VARCHAR(50) | | Access panels or equipment required |

---

### AIRCRAFT_MAINTENANCE_RECORD
Tracks the cumulative maintenance status and next-due intervals for each task on a specific aircraft.

| Attribute | Type | Key | Description |
|---|---|---|---|
| amrID | VARCHAR(30) | PK | Unique identifier for the aircraft maintenance record |
| aircraftID | VARCHAR(30) | FK | References AIRCRAFT.aircraftID |
| maintenanceTaskID | VARCHAR(30) | FK | References MAINTENANCE_TASK.maintenanceTaskID |
| lastCompletedDate | DATE | | Date task was last completed |
| lastCompletedFH | DECIMAL(10,2) | | Flight hours at last completion |
| lastCompletedFC | INTEGER | | Flight cycles at last completion |
| nextDueDate | DATE | | Next due date (calendar interval) |
| nextDueFH | DECIMAL(10,2) | | Next due flight hours |
| nextDueFC | INTEGER | | Next due flight cycles |
| overdueFlag | CHAR(1) | | Y / N — whether task is currently overdue |
| deferralAllowed | CHAR(1) | | Y / N — whether task may be deferred under MEL/CDL |
| deferralExpiryDate | DATE | | Deferral expiry date if deferred (nullable) |
| lastUpdatedDateTime | TIMESTAMP | | Timestamp of last record update |

---

### WORK_ORDER
Records a maintenance work order — authorising and tracking execution of one or more maintenance tasks.

| Attribute | Type | Key | Description |
|---|---|---|---|
| workOrderID | VARCHAR(30) | PK | Unique identifier for the work order |
| aircraftID | VARCHAR(30) | FK | References AIRCRAFT.aircraftID |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| amoID | VARCHAR(30) | FK | References AMO.amoID — performing organisation |
| certifyingEngineerID | VARCHAR(30) | FK | References CREW_MEMBER.crewMemberID — certifying engineer (nullable) |
| workOrderType | VARCHAR(20) | | Scheduled / Unscheduled / AOG / Modification |
| checkLevel | VARCHAR(5) | | A / B / C / D / Line |
| workOrderStatus | VARCHAR(20) | | Open / InProgress / Completed / Closed / Cancelled |
| openedDateTime | TIMESTAMP | | Date and time work order was opened |
| plannedStartDateTime | TIMESTAMP | | Planned start date and time |
| plannedEndDateTime | TIMESTAMP | | Planned completion date and time |
| actualStartDateTime | TIMESTAMP | | Actual start date and time (nullable) |
| actualEndDateTime | TIMESTAMP | | Actual completion date and time (nullable) |
| totalManHoursActual | DECIMAL(8,2) | | Total actual man-hours expended |
| maintenanceLocation | VARCHAR(50) | | Hangar / Line / Remote / AMO facility |
| airportCode | VARCHAR(3) | | IATA airport code where work is performed |
| crsIssued | CHAR(1) | | Y / N — Certificate of Release to Service issued |
| crsReference | VARCHAR(50) | | CRS document reference number (nullable) |
| crsIssuedDateTime | TIMESTAMP | | Date and time CRS was issued (nullable) |

---

### WORK_ORDER_TASK
Links individual maintenance tasks to a work order and records their completion status.

| Attribute | Type | Key | Description |
|---|---|---|---|
| workOrderTaskID | VARCHAR(30) | PK | Unique identifier for the work order task |
| workOrderID | VARCHAR(30) | FK | References WORK_ORDER.workOrderID |
| maintenanceTaskID | VARCHAR(30) | FK | References MAINTENANCE_TASK.maintenanceTaskID |
| taskStatus | VARCHAR(20) | | Pending / InProgress / Completed / Deferred / Carried-Forward |
| assignedEngineerID | VARCHAR(30) | | Engineer employee ID assigned to task |
| actualManHours | DECIMAL(6,2) | | Actual man-hours for this task |
| completedDateTime | TIMESTAMP | | Task completion timestamp (nullable) |
| findingsNarrative | VARCHAR(500) | | Technical findings during task execution |
| partsUsed | VARCHAR(255) | | Comma-separated part numbers used (nullable) |

---

### DEFECT_RECORD
Records a technical defect found during flight, line checks, or maintenance inspections.

| Attribute | Type | Key | Description |
|---|---|---|---|
| defectID | VARCHAR(30) | PK | Unique identifier for the defect record |
| aircraftID | VARCHAR(30) | FK | References AIRCRAFT.aircraftID |
| flightID | VARCHAR(30) | FK | References FLIGHT.flightID (nullable — if found in flight) |
| workOrderID | VARCHAR(30) | FK | References WORK_ORDER.workOrderID (nullable — if actioned) |
| defectSource | VARCHAR(20) | | Pilot / CabinCrew / Engineer / Maintenance / ATC |
| ataChapter | VARCHAR(5) | | ATA 100 chapter code |
| defectDescription | VARCHAR(500) | | Full technical description of defect |
| defectCategory | VARCHAR(20) | | MEL / CDL / Cosmetic / Safety / Performance |
| melReference | VARCHAR(20) | | Minimum Equipment List reference (nullable) |
| defectStatus | VARCHAR(20) | | Open / Deferred / Rectified / Closed / Monitoring |
| deferralCategory | VARCHAR(5) | | A / B / C / D — MEL deferral category (nullable) |
| deferralExpiryDate | DATE | | MEL deferral expiry date (nullable) |
| reportedDateTime | TIMESTAMP | | Date and time defect was reported |
| rectifiedDateTime | TIMESTAMP | | Date and time defect was rectified (nullable) |
| airworthyStatus | CHAR(1) | | Y / N — aircraft airworthy with this defect open |

---

### AIRWORTHINESS_DIRECTIVE
Records an Airworthiness Directive (AD) issued by a regulatory authority that must be complied with.

| Attribute | Type | Key | Description |
|---|---|---|---|
| adID | VARCHAR(30) | PK | Unique identifier for the AD record |
| adNumber | VARCHAR(30) | | Regulatory AD number (e.g., EASA AD 2024-0123) |
| issuingAuthority | VARCHAR(10) | | EASA / FAA / CAAC / TCCA |
| applicableAircraftType | VARCHAR(10) | | IATA aircraft type code |
| adTitle | VARCHAR(255) | | Full AD title |
| adDescription | TEXT | | Full AD description and compliance instructions |
| complianceType | VARCHAR(20) | | Recurring / OneTime / OnCondition |
| complianceDeadlineDate | DATE | | Mandatory compliance deadline date (nullable) |
| complianceFH | DECIMAL(8,2) | | Compliance deadline in flight hours (nullable) |
| complianceFC | INTEGER | | Compliance deadline in flight cycles (nullable) |
| recurringIntervalFH | DECIMAL(8,2) | | Recurring interval in flight hours (nullable) |
| recurringIntervalDays | INTEGER | | Recurring interval in calendar days (nullable) |
| effectiveDate | DATE | | AD effective date |
| status | VARCHAR(20) | | Active / Superseded / Cancelled |

---

### AD_COMPLIANCE_RECORD
Tracks compliance status of each AD against each applicable aircraft.

| Attribute | Type | Key | Description |
|---|---|---|---|
| adComplianceID | VARCHAR(30) | PK | Unique identifier for the AD compliance record |
| adID | VARCHAR(30) | FK | References AIRWORTHINESS_DIRECTIVE.adID |
| aircraftID | VARCHAR(30) | FK | References AIRCRAFT.aircraftID |
| workOrderID | VARCHAR(30) | FK | References WORK_ORDER.workOrderID (nullable) |
| complianceStatus | VARCHAR(20) | | Compliant / Due / Overdue / NotApplicable |
| lastComplianceDate | DATE | | Date of last compliance (nullable) |
| lastComplianceFH | DECIMAL(10,2) | | Flight hours at last compliance (nullable) |
| lastComplianceFC | INTEGER | | Flight cycles at last compliance (nullable) |
| nextDueDate | DATE | | Next compliance due date (nullable) |
| nextDueFH | DECIMAL(10,2) | | Next due flight hours (nullable) |
| nextDueFC | INTEGER | | Next due flight cycles (nullable) |
| complianceNotes | VARCHAR(255) | | Notes on compliance method or findings |

---

### COMPONENT
Tracks a tracked/serialised aircraft component (engine, APU, landing gear, etc.).

| Attribute | Type | Key | Description |
|---|---|---|---|
| componentID | VARCHAR(30) | PK | Unique identifier for the component |
| aircraftID | VARCHAR(30) | FK | References AIRCRAFT.aircraftID — current installation (nullable if in pool) |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode — owner |
| partNumber | VARCHAR(30) | | Manufacturer part number |
| serialNumber | VARCHAR(30) | | Component serial number |
| componentType | VARCHAR(30) | | Engine / APU / LandingGear / Avionics / Propeller / LifeVest |
| componentDescription | VARCHAR(100) | | Description of component |
| manufacturerName | VARCHAR(50) | | Component manufacturer name |
| installationDate | DATE | | Date component was installed on aircraft (nullable) |
| installationFH | DECIMAL(10,2) | | Aircraft flight hours at installation |
| installationFC | INTEGER | | Aircraft flight cycles at installation |
| componentFHTSN | DECIMAL(10,2) | | Component total flight hours since new |
| componentFCTSN | INTEGER | | Component total flight cycles since new |
| componentFHTSO | DECIMAL(10,2) | | Component flight hours since last overhaul |
| componentFCTSO | INTEGER | | Component flight cycles since last overhaul |
| hardLifeLimitFH | DECIMAL(10,2) | | Hard life limit in flight hours (nullable) |
| hardLifeLimitFC | INTEGER | | Hard life limit in flight cycles (nullable) |
| componentStatus | VARCHAR(20) | | Installed / Serviceable / Unserviceable / Scrapped / InShop |
| locationCode | VARCHAR(30) | | Current location (aircraft reg or shop/warehouse code) |

---

### AMO
Defines an Approved Maintenance Organisation authorised to perform maintenance on airline aircraft.

| Attribute | Type | Key | Description |
|---|---|---|---|
| amoID | VARCHAR(30) | PK | Unique identifier for the AMO |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode — contracting airline |
| amoName | VARCHAR(100) | | Full legal name of the AMO |
| approvalNumber | VARCHAR(30) | | Regulatory approval reference (e.g., EASA.145.0001) |
| approvalAuthority | VARCHAR(10) | | EASA / FAA / CAAC / TCCA |
| amoType | VARCHAR(20) | | Airframe / Engine / Components / Avionics / Line |
| countryCode | VARCHAR(3) | | ISO 3166-1 alpha-3 country of operation |
| airportCode | VARCHAR(3) | | IATA airport code of AMO facility (nullable) |
| approvalExpiryDate | DATE | | Approval expiry date |
| status | VARCHAR(20) | | Active / Suspended / Revoked |

---

### PARTS_INVENTORY
Tracks consumable and rotable parts inventory held at airline or AMO warehouses.

| Attribute | Type | Key | Description |
|---|---|---|---|
| partsInventoryID | VARCHAR(30) | PK | Unique identifier for the parts inventory record |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| partNumber | VARCHAR(30) | | Manufacturer part number |
| partDescription | VARCHAR(100) | | Part description |
| partCategory | VARCHAR(20) | | Rotable / Consumable / Expendable / Chemical |
| quantityOnHand | INTEGER | | Current stock quantity |
| quantityOnOrder | INTEGER | | Outstanding purchase order quantity |
| minimumStockLevel | INTEGER | | Reorder trigger level |
| unitCost | DECIMAL(12,2) | | Unit cost of part |
| currency | CHAR(3) | | ISO 4217 currency code |
| warehouseLocation | VARCHAR(50) | | Warehouse and bin location code |
| airportCode | VARCHAR(3) | | IATA airport code of warehouse |
| certificationRef | VARCHAR(30) | | EASA Form 1 / FAA 8130-3 certificate reference |
| expiryDate | DATE | | Part shelf-life expiry (nullable for non-life-limited) |
| conditionCode | VARCHAR(10) | | New / Overhauled / Serviceable / Unserviceable / Scrap |

---

## Relationships

### Internal — MRO Domain

| From Entity | To Entity | Cardinality | Description |
|---|---|---|---|
| MAINTENANCE_PROGRAMME | MAINTENANCE_TASK | 1..* | A programme contains one or more tasks |
| MAINTENANCE_TASK | AIRCRAFT_MAINTENANCE_RECORD | 0..* | A task has zero or more aircraft-level records |
| MAINTENANCE_TASK | WORK_ORDER_TASK | 0..* | A task appears in zero or more work orders |
| WORK_ORDER | WORK_ORDER_TASK | 1..* | A work order contains one or more tasks |
| WORK_ORDER | DEFECT_RECORD | 0..* | A work order may address zero or more defects |
| AIRWORTHINESS_DIRECTIVE | AD_COMPLIANCE_RECORD | 1..* | An AD has one or more compliance records |
| AMO | WORK_ORDER | 0..* | An AMO performs zero or more work orders |

### Cross-Domain — MRO → Previous Domains

| From Entity | To Entity | Domain | Description |
|---|---|---|---|
| MAINTENANCE_PROGRAMME | AIRLINE | 01 | Programme owned by airline |
| WORK_ORDER | AIRLINE | 01 | Work order raised by airline |
| WORK_ORDER | AIRCRAFT | 08 | Work order performed on aircraft |
| WORK_ORDER | AMO | 10 | Work order executed by AMO |
| WORK_ORDER | CREW_MEMBER | 08 | Work order certified by engineer |
| AIRCRAFT_MAINTENANCE_RECORD | AIRCRAFT | 08 | Maintenance record for an aircraft |
| DEFECT_RECORD | AIRCRAFT | 08 | Defect raised against aircraft |
| DEFECT_RECORD | FLIGHT | 08 | Defect found during a flight |
| DEFECT_RECORD | WORK_ORDER | 10 | Defect actioned by work order |
| AD_COMPLIANCE_RECORD | AIRCRAFT | 08 | Compliance tracked per aircraft |
| AD_COMPLIANCE_RECORD | WORK_ORDER | 10 | Compliance performed via work order |
| COMPONENT | AIRCRAFT | 08 | Component installed on aircraft |
| COMPONENT | AIRLINE | 01 | Component owned by airline |
| AMO | AIRLINE | 01 | AMO contracted by airline |
| PARTS_INVENTORY | AIRLINE | 01 | Parts inventory held by airline |

---

## Enumerations and Code Lists

### CHECK_LEVEL_CD
| Code | Description |
|---|---|
| TRANSIT | Transit check between flights |
| PRE_FLIGHT | Pre-flight walkaround |
| A | A-Check — light line maintenance |
| B | B-Check — intermediate check |
| C | C-Check — heavy structural check |
| D | D-Check — full overhaul / heavy maintenance |

### WORK_ORDER_TYPE_CD
| Code | Description |
|---|---|
| SCHEDULED | Planned scheduled maintenance |
| UNSCHEDULED | Reactive maintenance — defect rectification |
| AOG | Aircraft on Ground — urgent critical repair |
| MODIFICATION | Modification / Service Bulletin embodiment |

### DEFECT_CATEGORY_CD
| Code | Description |
|---|---|
| MEL | Minimum Equipment List item — deferrable |
| CDL | Configuration Deviation List item |
| SAFETY | Safety-critical — aircraft must be grounded |
| PERFORMANCE | Performance-affecting defect |
| COSMETIC | Cosmetic — no airworthiness impact |

### MEL_DEFERRAL_CATEGORY_CD (IATA/ICAO Standard)
| Code | Description |
|---|---|
| A | Repair interval defined by airline — usually < 3 days |
| B | Repair within 3 consecutive calendar days |
| C | Repair within 10 consecutive calendar days |
| D | Repair within 120 consecutive calendar days |

### COMPONENT_STATUS_CD
| Code | Description |
|---|---|
| INSTALLED | Installed and serviceable on aircraft |
| SERVICEABLE | Removed and available for installation |
| UNSERVICEABLE | Removed — requires repair or overhaul |
| SCRAPPED | Beyond economic repair — condemned |
| IN_SHOP | Currently in MRO shop for overhaul |

### PART_CATEGORY_CD
| Code | Description |
|---|---|
| ROTABLE | Repairable component — returned to service after overhaul |
| CONSUMABLE | Single-use item — expended on use |
| EXPENDABLE | Non-repairable — discarded on removal |
| CHEMICAL | Lubricant, sealant, cleaning agent |

---

## Business Rules

| Rule ID | Entity | Rule Description |
|---|---|---|
| BR-MRO-001 | WORK_ORDER | crsIssued must be Y before AIRCRAFT.status can return to 'Active' from 'Grounded' |
| BR-MRO-002 | AIRCRAFT_MAINTENANCE_RECORD | overdueFlag must be set to Y when current date exceeds nextDueDate or current FH exceeds nextDueFH |
| BR-MRO-003 | DEFECT_RECORD | airworthyStatus must be N for defectCategory = 'Safety' until rectifiedDateTime is populated |
| BR-MRO-004 | AD_COMPLIANCE_RECORD | complianceStatus must not be 'Overdue' for any AD for aircraft to be released to service |
| BR-MRO-005 | COMPONENT | componentFHTSN must be less than hardLifeLimitFH when hardLifeLimitFH is not null |
| BR-MRO-006 | COMPONENT | componentFCTSN must be less than hardLifeLimitFC when hardLifeLimitFC is not null |
| BR-MRO-007 | PARTS_INVENTORY | quantityOnHand must trigger purchase order when it falls below minimumStockLevel |
| BR-MRO-008 | AMO | approvalExpiryDate must be in the future for status = 'Active' |
| BR-MRO-009 | MAINTENANCE_TASK | At least one of intervalFH, intervalFC, or intervalDays must be non-null |
| BR-MRO-010 | WORK_ORDER | actualEndDateTime must be after actualStartDateTime when both are populated |
| BR-MRO-011 | DEFECT_RECORD | deferralExpiryDate must be populated when deferralCategory is A, B, C, or D |

---

## Traceability to AIDM Domains

| AIDM Domain | AIDM Sub-domain | Entities |
|---|---|---|
| Operations | Maintenance Programme | MAINTENANCE_PROGRAMME, MAINTENANCE_TASK |
| Operations | Maintenance Tracking | AIRCRAFT_MAINTENANCE_RECORD |
| Operations | Work Order Management | WORK_ORDER, WORK_ORDER_TASK |
| Operations | Defect Management | DEFECT_RECORD |
| Operations | Airworthiness | AIRWORTHINESS_DIRECTIVE, AD_COMPLIANCE_RECORD |
| Operations | Component Tracking | COMPONENT |
| Operations | AMO Management | AMO |
| Operations | Parts & Inventory | PARTS_INVENTORY |

---

## Notes

- Maintenance programme approval follows EASA Part-M Subpart G (Continuing Airworthiness Management) and FAA AC 120-16 standards.
- ATA chapter codes follow ATA iSpec 2200 / ATA 100 chapter numbering.
- MEL/CDL deferral categories follow IATA MEL Policy Guidelines and ICAO Doc 9760.
- Airworthiness Directive numbering follows EASA AD format and FAA AD format (e.g., 2024-XX-XX).
- Component tracking (TSN/TSO) follows EASA Part-145 and FAA Part-145 maintenance records requirements.
- Certificate of Release to Service (CRS) follows EASA Form 1 / FAA 8130-3 dual-release format.
- Parts certification references follow EASA Part-145.A.50 release documentation requirements.
- Hard life limits for life-limited parts follow EASA CS-25 Airworthiness Limitations Section (ALS).
