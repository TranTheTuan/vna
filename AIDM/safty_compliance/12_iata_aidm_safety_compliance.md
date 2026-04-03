# IATA AIDM 25.1 — Domain 12: Safety & Compliance
## Mockup Data Model for AI Consumption

---

## Subject Area Overview

The **Safety & Compliance** domain in the IATA AIDM covers the airline's safety management framework — including the Safety Management System (SMS), hazard identification and risk assessment, safety occurrence reporting, audit and finding management, regulatory compliance tracking, IOSA (IATA Operational Safety Audit) certification, dangerous goods incident management, crew fatigue risk management, and emergency response planning. It provides the governance layer across all operational domains.

- **AIDM Version**: 25.1
- **Domain**: Governance — Safety & Compliance
- **Integrates With**: All Domains (01–11) — safety and compliance cuts across the entire enterprise

---

## Cross-Domain FK Dependencies

| Referenced Entity | PK Used | Referenced By (Safety & Compliance) |
|---|---|---|
| AIRLINE | airlineCode | SMS_PROGRAMME, SAFETY_OCCURRENCE, AUDIT, REGULATORY_REQUIREMENT, EMERGENCY_RESPONSE_PLAN |
| AIRCRAFT | aircraftID | SAFETY_OCCURRENCE |
| FLIGHT | flightID | SAFETY_OCCURRENCE |
| CREW_MEMBER | crewMemberID | SAFETY_OCCURRENCE, FATIGUE_RISK_RECORD |
| DEFECT_RECORD | defectID | SAFETY_OCCURRENCE |
| GROUND_HANDLER | groundHandlerID | AUDIT |
| AMO | amoID | AUDIT |
| WORK_ORDER | workOrderID | SAFETY_OCCURRENCE |

---

## Entities

### SMS_PROGRAMME
Defines the airline's Safety Management System programme — the overarching safety governance framework.

| Attribute | Type | Key | Description |
|---|---|---|---|
| smsProgrammeID | VARCHAR(30) | PK | Unique identifier for the SMS programme |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| programmeName | VARCHAR(100) | | SMS programme name |
| programmeVersion | VARCHAR(10) | | Current version number |
| icaoAnnex19Compliant | CHAR(1) | | Y / N — ICAO Annex 19 (Safety Management) compliance |
| approvalAuthority | VARCHAR(10) | | Regulatory authority approving the SMS (EASA/FAA/CAAC) |
| approvalReference | VARCHAR(50) | | Regulatory approval reference number |
| approvalDate | DATE | | Date SMS was approved |
| accountableManager | VARCHAR(100) | | Name of Accountable Manager (ICAO requirement) |
| safetyOfficerName | VARCHAR(100) | | Chief Safety Officer / Safety Manager name |
| lastReviewDate | DATE | | Date of last SMS review |
| nextReviewDate | DATE | | Date of next scheduled SMS review |
| status | VARCHAR(20) | | Active / UnderReview / Suspended |

---

### HAZARD_REGISTER
Records identified hazards within the airline's safety risk management framework.

| Attribute | Type | Key | Description |
|---|---|---|---|
| hazardID | VARCHAR(30) | PK | Unique identifier for the hazard |
| smsProgrammeID | VARCHAR(30) | FK | References SMS_PROGRAMME.smsProgrammeID |
| hazardTitle | VARCHAR(100) | | Short title of hazard |
| hazardDescription | VARCHAR(500) | | Detailed description of the hazard |
| hazardCategory | VARCHAR(30) | | Operational / Technical / Environmental / Human / Organisational |
| operationalArea | VARCHAR(30) | | FlightOps / GroundOps / Maintenance / Cargo / Security |
| inherentLikelihood | INTEGER | | Inherent likelihood rating (1–5 per ICAO matrix) |
| inherentSeverity | INTEGER | | Inherent severity rating (1–5 per ICAO matrix) |
| inherentRiskScore | INTEGER | | Computed: likelihood × severity |
| residualLikelihood | INTEGER | | Residual likelihood after controls |
| residualSeverity | INTEGER | | Residual severity after controls |
| residualRiskScore | INTEGER | | Computed residual risk score |
| riskAcceptanceStatus | VARCHAR(20) | | Acceptable / ALARP / Unacceptable |
| controlMeasures | VARCHAR(500) | | Description of risk control measures applied |
| hazardOwner | VARCHAR(50) | | Name or role of hazard owner |
| reviewDate | DATE | | Next review date for this hazard |
| status | VARCHAR(20) | | Open / Mitigated / Closed / Monitoring |

---

### SAFETY_OCCURRENCE
Records a safety event — incident, accident, near miss, or mandatory occurrence report (MOR).

| Attribute | Type | Key | Description |
|---|---|---|---|
| occurrenceID | VARCHAR(30) | PK | Unique identifier for the safety occurrence |
| smsProgrammeID | VARCHAR(30) | FK | References SMS_PROGRAMME.smsProgrammeID |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| aircraftID | VARCHAR(30) | FK | References AIRCRAFT.aircraftID (nullable) |
| flightID | VARCHAR(30) | FK | References FLIGHT.flightID (nullable) |
| crewMemberID | VARCHAR(30) | FK | References CREW_MEMBER.crewMemberID (nullable — reporting crew) |
| defectID | VARCHAR(30) | FK | References DEFECT_RECORD.defectID (nullable) |
| workOrderID | VARCHAR(30) | FK | References WORK_ORDER.workOrderID (nullable) |
| occurrenceType | VARCHAR(20) | | Accident / SeriousIncident / Incident / NearMiss / MOR |
| occurrenceCategory | VARCHAR(30) | | ICAO occurrence category (e.g., LOC-I, CFIT, RE, FUEL) |
| occurrenceDateTime | TIMESTAMP | | Date and time of the occurrence |
| airportCode | VARCHAR(3) | | IATA airport code where occurrence happened (nullable) |
| flightPhase | VARCHAR(20) | | Preflight / Taxi / Takeoff / Climb / Cruise / Descent / Approach / Landing |
| occurrenceDescription | VARCHAR(1000) | | Detailed narrative of the occurrence |
| injuriesCount | INTEGER | | Number of injuries (nullable) |
| fatalitiesCount | INTEGER | | Number of fatalities (nullable) |
| aircraftDamage | VARCHAR(20) | | None / Minor / Substantial / Destroyed |
| reportedToAuthority | CHAR(1) | | Y / N — whether reported to regulatory authority |
| authorityReference | VARCHAR(50) | | Authority occurrence reference number (nullable) |
| reportingSystem | VARCHAR(20) | | ASRS / ECCAIRS / MOR / Internal |
| safetyInvestigationRequired | CHAR(1) | | Y / N |
| occurrenceStatus | VARCHAR(20) | | Reported / UnderInvestigation / Closed / Corrective |

---

### SAFETY_INVESTIGATION
Records the formal investigation of a safety occurrence.

| Attribute | Type | Key | Description |
|---|---|---|---|
| investigationID | VARCHAR(30) | PK | Unique identifier for the safety investigation |
| occurrenceID | VARCHAR(30) | FK | References SAFETY_OCCURRENCE.occurrenceID |
| investigationLeadName | VARCHAR(100) | | Name of investigation lead |
| investigationTeam | VARCHAR(255) | | Comma-separated names of investigation team members |
| investigationStartDate | DATE | | Date investigation was opened |
| targetCloseDate | DATE | | Target date for investigation closure |
| actualCloseDate | DATE | | Actual closure date (nullable) |
| rootCauses | VARCHAR(500) | | Identified root cause(s) |
| contributingFactors | VARCHAR(500) | | Contributing factors identified |
| immediateActions | VARCHAR(500) | | Immediate safety actions taken |
| safetyRecommendations | VARCHAR(500) | | Safety recommendations arising |
| investigationStatus | VARCHAR(20) | | Open / InProgress / PendingApproval / Closed |
| finalReportReference | VARCHAR(50) | | Final investigation report document reference (nullable) |

---

### CORRECTIVE_ACTION
Records a corrective or preventive action arising from a safety investigation or audit finding.

| Attribute | Type | Key | Description |
|---|---|---|---|
| correctiveActionID | VARCHAR(30) | PK | Unique identifier for the corrective action |
| investigationID | VARCHAR(30) | FK | References SAFETY_INVESTIGATION.investigationID (nullable) |
| auditFindingID | VARCHAR(30) | FK | References AUDIT_FINDING.auditFindingID (nullable) |
| actionTitle | VARCHAR(100) | | Short title of corrective action |
| actionDescription | VARCHAR(500) | | Detailed description of corrective action |
| actionType | VARCHAR(20) | | Corrective / Preventive / Improvement |
| responsibleParty | VARCHAR(100) | | Name or role responsible for action |
| targetCompletionDate | DATE | | Target completion date |
| actualCompletionDate | DATE | | Actual completion date (nullable) |
| evidenceReference | VARCHAR(100) | | Reference to evidence of completion |
| actionStatus | VARCHAR(20) | | Open / InProgress / Completed / Verified / Overdue |
| verifiedBy | VARCHAR(50) | | Name of person who verified completion (nullable) |
| verifiedDate | DATE | | Date completion was verified (nullable) |

---

### AUDIT
Records a safety, quality, or compliance audit conducted against the airline or its suppliers.

| Attribute | Type | Key | Description |
|---|---|---|---|
| auditID | VARCHAR(30) | PK | Unique identifier for the audit |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| groundHandlerID | VARCHAR(30) | FK | References GROUND_HANDLER.groundHandlerID (nullable) |
| amoID | VARCHAR(30) | FK | References AMO.amoID (nullable) |
| auditType | VARCHAR(20) | | IOSA / Internal / Regulatory / Supplier / RAMP |
| auditArea | VARCHAR(30) | | FlightOps / GroundOps / Maintenance / Cargo / Security / Cabin |
| auditorName | VARCHAR(100) | | Lead auditor name |
| auditOrganisation | VARCHAR(100) | | Auditing organisation (IATA / Regulator / Internal) |
| auditStartDate | DATE | | Audit start date |
| auditEndDate | DATE | | Audit end date |
| auditLocation | VARCHAR(50) | | Location where audit was conducted |
| overallRating | VARCHAR(20) | | Satisfactory / MinorFindings / MajorFindings / Critical |
| findingsCount | INTEGER | | Total number of findings raised |
| auditStatus | VARCHAR(20) | | Planned / InProgress / Completed / Closed |
| certificateIssued | CHAR(1) | | Y / N — whether certificate was issued (e.g., IOSA) |
| certificateExpiryDate | DATE | | Certificate expiry date (nullable) |

---

### AUDIT_FINDING
Records an individual finding raised during an audit.

| Attribute | Type | Key | Description |
|---|---|---|---|
| auditFindingID | VARCHAR(30) | PK | Unique identifier for the audit finding |
| auditID | VARCHAR(30) | FK | References AUDIT.auditID |
| findingNumber | VARCHAR(20) | | Finding reference number |
| findingType | VARCHAR(20) | | Finding / Observation / Opportunity / NonConformity |
| findingSeverity | VARCHAR(10) | | Critical / Major / Minor |
| iosaStandardRef | VARCHAR(20) | | IOSA Standard reference (e.g., FLT 1.1.1) — nullable |
| findingDescription | VARCHAR(500) | | Detailed finding description |
| standardRequirement | VARCHAR(255) | | Relevant standard or regulatory requirement |
| objectiveEvidence | VARCHAR(500) | | Objective evidence supporting the finding |
| rootCause | VARCHAR(255) | | Root cause of finding (nullable) |
| findingStatus | VARCHAR(20) | | Open / UnderReview / Closed / Verified |
| targetCloseDate | DATE | | Target date for finding closure |
| actualCloseDate | DATE | | Actual closure date (nullable) |

---

### REGULATORY_REQUIREMENT
Tracks applicable regulatory requirements and the airline's compliance status against each.

| Attribute | Type | Key | Description |
|---|---|---|---|
| regulatoryReqID | VARCHAR(30) | PK | Unique identifier for the regulatory requirement |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| regulationReference | VARCHAR(50) | | Regulation reference (e.g., EASA Part-OPS ORO.FTL.205) |
| regulationTitle | VARCHAR(255) | | Full title of the regulation or standard |
| regulatoryAuthority | VARCHAR(10) | | ICAO / EASA / FAA / CAAC / TCCA |
| applicableArea | VARCHAR(30) | | FlightOps / GroundOps / Maintenance / Cargo / Security / Airworthiness |
| complianceStatus | VARCHAR(20) | | Compliant / PartiallyCompliant / NonCompliant / NotAssessed |
| complianceEvidence | VARCHAR(255) | | Reference to compliance evidence |
| lastAssessmentDate | DATE | | Date of last compliance assessment |
| nextAssessmentDate | DATE | | Date of next required assessment |
| responsibleDepartment | VARCHAR(50) | | Responsible department or role |
| notes | VARCHAR(255) | | Additional notes on compliance status |

---

### FATIGUE_RISK_RECORD
Records fatigue risk management data for a crew member — supporting FRMS compliance.

| Attribute | Type | Key | Description |
|---|---|---|---|
| fatigueRiskID | VARCHAR(30) | PK | Unique identifier for the fatigue risk record |
| crewMemberID | VARCHAR(30) | FK | References CREW_MEMBER.crewMemberID |
| recordDate | DATE | | Date of the fatigue risk record |
| reportedFatigueLevel | INTEGER | | Self-reported fatigue level (1–10 Samn-Perelli scale) |
| fatigueSource | VARCHAR(20) | | PreDuty / InFlight / PostFlight |
| dutyPeriodMins | INTEGER | | Actual duty period in minutes |
| sleepHoursLast24 | DECIMAL(4,1) | | Hours of sleep in the 24 hours before duty |
| sleepHoursLast48 | DECIMAL(4,1) | | Hours of sleep in the 48 hours before duty |
| biomathModelScore | DECIMAL(5,2) | | Fatigue bio-mathematical model score (nullable) |
| biomathModelUsed | VARCHAR(30) | | SAFE / FAID / FRMS model used (nullable) |
| frmsAlertTriggered | CHAR(1) | | Y / N — whether FRMS alert was triggered |
| actionTaken | VARCHAR(255) | | Action taken if alert triggered (nullable) |
| reportedDateTime | TIMESTAMP | | Timestamp record was submitted |

---

### EMERGENCY_RESPONSE_PLAN
Defines emergency response plans held by the airline for various emergency scenarios.

| Attribute | Type | Key | Description |
|---|---|---|---|
| erpID | VARCHAR(30) | PK | Unique identifier for the emergency response plan |
| airlineCode | VARCHAR(3) | FK | References AIRLINE.airlineCode |
| planName | VARCHAR(100) | | Name of the emergency response plan |
| emergencyType | VARCHAR(30) | | AircraftAccident / BombThreat / HijackingUnlawfulInterference / MedicalEmergency / NaturalDisaster / CyberIncident |
| planVersion | VARCHAR(10) | | Plan version number |
| approvalDate | DATE | | Date plan was approved |
| lastExerciseDate | DATE | | Date of last full-scale exercise |
| nextExerciseDate | DATE | | Date of next scheduled exercise |
| planOwner | VARCHAR(100) | | Name or role of plan owner |
| icaoAnnex13Compliant | CHAR(1) | | Y / N — ICAO Annex 13 (Aircraft Accident Investigation) compliant |
| regulatoryApprovalRef | VARCHAR(50) | | Regulatory approval reference (nullable) |
| planStatus | VARCHAR(20) | | Active / UnderReview / Superseded |

---

## Relationships

### Internal — Safety & Compliance Domain

| From Entity | To Entity | Cardinality | Description |
|---|---|---|---|
| SMS_PROGRAMME | HAZARD_REGISTER | 0..* | An SMS has zero or more registered hazards |
| SMS_PROGRAMME | SAFETY_OCCURRENCE | 0..* | An SMS captures zero or more occurrences |
| SAFETY_OCCURRENCE | SAFETY_INVESTIGATION | 0..1 | An occurrence may trigger one investigation |
| SAFETY_INVESTIGATION | CORRECTIVE_ACTION | 0..* | An investigation may produce zero or more actions |
| AUDIT | AUDIT_FINDING | 0..* | An audit has zero or more findings |
| AUDIT_FINDING | CORRECTIVE_ACTION | 0..* | A finding may produce zero or more corrective actions |

### Cross-Domain — Safety → Previous Domains

| From Entity | To Entity | Domain | Description |
|---|---|---|---|
| SMS_PROGRAMME | AIRLINE | 01 | SMS owned by airline |
| SAFETY_OCCURRENCE | AIRLINE | 01 | Occurrence attributed to airline |
| SAFETY_OCCURRENCE | AIRCRAFT | 08 | Occurrence involved aircraft |
| SAFETY_OCCURRENCE | FLIGHT | 08 | Occurrence on a specific flight |
| SAFETY_OCCURRENCE | CREW_MEMBER | 08 | Occurrence reported by crew member |
| SAFETY_OCCURRENCE | DEFECT_RECORD | 10 | Occurrence linked to a defect |
| SAFETY_OCCURRENCE | WORK_ORDER | 10 | Occurrence linked to a work order |
| AUDIT | AIRLINE | 01 | Audit of airline |
| AUDIT | GROUND_HANDLER | 07 | Audit of ground handler |
| AUDIT | AMO | 10 | Audit of AMO |
| REGULATORY_REQUIREMENT | AIRLINE | 01 | Requirement applicable to airline |
| FATIGUE_RISK_RECORD | CREW_MEMBER | 08 | Fatigue record for a crew member |
| EMERGENCY_RESPONSE_PLAN | AIRLINE | 01 | ERP held by airline |

---

## Enumerations and Code Lists

### OCCURRENCE_TYPE_CD (ICAO Annex 13)
| Code | Description |
|---|---|
| ACCIDENT | Aircraft accident — damage or injury |
| SERIOUS_INCIDENT | Serious incident — nearly an accident |
| INCIDENT | Incident — safety significance |
| NEAR_MISS | Near miss / airprox event |
| MOR | Mandatory Occurrence Report |

### ICAO_OCCURRENCE_CATEGORY_CD (Partial — ICAO ADREP)
| Code | Description |
|---|---|
| LOC-I | Loss of Control In-flight |
| CFIT | Controlled Flight Into Terrain |
| RE | Runway Excursion |
| RI | Runway Incursion |
| FUEL | Fuel Related |
| F-NI | Fire/Smoke (Non-Impact) |
| TURB | Turbulence Encounter |
| UIMC | Unintended IMC |
| WSTRW | Windshear/Thunderstorm |

### AUDIT_TYPE_CD
| Code | Description |
|---|---|
| IOSA | IATA Operational Safety Audit |
| INTERNAL | Internal airline quality audit |
| REGULATORY | Regulatory authority inspection |
| SUPPLIER | Supplier/vendor audit |
| RAMP | IATA Ramp Safety Assessment |

### RISK_ACCEPTANCE_STATUS_CD (ICAO Safety Risk Matrix)
| Code | Description |
|---|---|
| ACCEPTABLE | Risk is acceptable — no action required |
| ALARP | As Low As Reasonably Practicable — controls in place |
| UNACCEPTABLE | Risk is unacceptable — operations must cease |

### FINDING_SEVERITY_CD
| Code | Description |
|---|---|
| CRITICAL | Immediate safety risk — operations must stop |
| MAJOR | Significant non-conformity — urgent action required |
| MINOR | Minor non-conformity — action required within defined period |

### EMERGENCY_TYPE_CD
| Code | Description |
|---|---|
| AIRCRAFT_ACCIDENT | Full aircraft accident response |
| BOMB_THREAT | Bomb threat protocol |
| HIJACKING | Unlawful interference / hijacking |
| MEDICAL_EMERGENCY | In-flight or airport medical emergency |
| NATURAL_DISASTER | Natural disaster affecting operations |
| CYBER_INCIDENT | Cybersecurity incident response |

---

## Business Rules

| Rule ID | Entity | Rule Description |
|---|---|---|
| BR-SAF-001 | SAFETY_OCCURRENCE | occurrenceType = 'Accident' or 'SeriousIncident' must set reportedToAuthority = Y within 72 hours per ICAO Annex 13 |
| BR-SAF-002 | SAFETY_OCCURRENCE | safetyInvestigationRequired must be Y for all occurrences where occurrenceType = 'Accident' or 'SeriousIncident' |
| BR-SAF-003 | CORRECTIVE_ACTION | actionStatus must be set to 'Overdue' when targetCompletionDate is exceeded and actualCompletionDate is null |
| BR-SAF-004 | AUDIT | certificateExpiryDate must be populated when certificateIssued = Y |
| BR-SAF-005 | HAZARD_REGISTER | inherentRiskScore must equal inherentLikelihood multiplied by inherentSeverity |
| BR-SAF-006 | HAZARD_REGISTER | residualRiskScore must equal residualLikelihood multiplied by residualSeverity |
| BR-SAF-007 | REGULATORY_REQUIREMENT | nextAssessmentDate must be after lastAssessmentDate |
| BR-SAF-008 | SMS_PROGRAMME | accountableManager must be populated for ICAO Annex 19 compliance |
| BR-SAF-009 | FATIGUE_RISK_RECORD | reportedFatigueLevel must be between 1 and 10 (Samn-Perelli scale) |
| BR-SAF-010 | EMERGENCY_RESPONSE_PLAN | lastExerciseDate must be within 24 calendar months per ICAO Doc 9137 |
| BR-SAF-011 | AUDIT_FINDING | findingSeverity = 'Critical' must generate a CORRECTIVE_ACTION within 24 hours |

---

## Traceability to AIDM Domains

| AIDM Domain | AIDM Sub-domain | Entities |
|---|---|---|
| Safety | Safety Management System | SMS_PROGRAMME |
| Safety | Hazard & Risk Management | HAZARD_REGISTER |
| Safety | Occurrence Reporting | SAFETY_OCCURRENCE, SAFETY_INVESTIGATION |
| Safety | Corrective Actions | CORRECTIVE_ACTION |
| Safety | Audit & Quality | AUDIT, AUDIT_FINDING |
| Safety | Regulatory Compliance | REGULATORY_REQUIREMENT |
| Safety | Fatigue Risk Management | FATIGUE_RISK_RECORD |
| Safety | Emergency Response | EMERGENCY_RESPONSE_PLAN |

---

## Notes

- SMS framework follows ICAO Annex 19 (Safety Management) — Second Edition 2016.
- Safety occurrence reporting follows ICAO Annex 13 (Aircraft Accident and Incident Investigation).
- Occurrence categories follow ICAO ADREP (Accident/Incident Data Reporting) taxonomy.
- IOSA audit standard references IATA IOSA Standards Manual (ISM) — current edition.
- Hazard risk matrix follows ICAO Doc 9859 (Safety Management Manual — SMM) 5×5 risk matrix.
- Fatigue risk management follows ICAO Doc 9966 (FRMS Manual) and EASA ORO.FTL Subpart Q.
- Regulatory compliance mapping covers ICAO Annex 1 through Annex 19 applicable requirements.
- Emergency response planning follows ICAO Doc 9137 (Airport Services Manual) Part 7 and airline ERP standards.
- Corrective action tracking methodology follows ICAO continuous improvement loop (Plan-Do-Check-Act).
