package idscan

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	ServiceDLDV            = "F1366C39-D0CD-458C-81F8-EFD7B0753A5B"
	ServiceIdentiFraud     = "211B428E-B8B5-4EB9-973D-106F013C371F"
	ServiceSexOffender     = "3ACCF2A9-2E52-4675-B546-7C3676897C3C"
	ServiceCrimeRecord     = "80C1111D-6ACE-4820-A88C-284982013C33"
	ServicePEP             = "4DF86C8E-B53E-49D6-B75D-C1F2FFCBC1E5"
	ServiceOFAC            = "D0F7EE7A-4CCA-4807-A779-F77EB1501EED"
	ServiceEveryPolitician = "C18ECBD7-907C-47AF-8871-1C6BB19838CA"
)

type ScreeningAPI struct {
	token string
}

type ScreeningAPIRequest struct {
	FirstName               string   `json:"firstName,omitempty"`
	MiddleName              string   `json:"middleName,omitempty"`
	LastName                string   `json:"lastName,omitempty"`
	DateOfBirths            string   `json:"dateOfBirths,omitempty"`
	Sex                     string   `json:"sex,omitempty"`
	IDType                  string   `json:"idType,omitempty"`
	IDNumber                string   `json:"idNumber,omitempty"`
	Address                 string   `json:"address,omitempty"`
	AddressLine2            string   `json:"addressLine2,omitempty"`
	City                    string   `json:"city,omitempty"`
	State                   string   `json:"state,omitempty"`
	ZIP                     string   `json:"zip,omitempty"`
	County                  string   `json:"county,omitempty"`
	SSN                     string   `json:"ssn,omitempty"`
	Services                []string `json:"services"`
	DriverLicenseNumber     string   `json:"driverLicenseNumber,omitempty"`
	DocumentCategoryCode    uint     `json:"documentCategoryCode"`
	DriverLicenseIssueDate  string   `json:"driverLicenseIssueDate,omitempty"`
	DriverLicenseExpireDate string   `json:"driverLicenseExpireDate,omitempty"`
	ReferenceID             string   `json:"referenceId,omitempty"`
}

type ScreeningAPIResult struct {
	ServiceID          string                `json:"serviceID"`
	ServiceName        string                `json:"serviceName"`
	ServiceDescription string                `json:"serviceDescription"`
	Error              *string               `json:"error"`
	Success            bool                  `json:"success"`
	Profiles           []ScreeningAPIProfile `json:"profiles"`
}

type ScreeningAPIProfile struct {
	InternalID                       *string                         `json:"internalId"`
	FirstName                        *string                         `json:"firstName"`
	MiddleName                       *string                         `json:"middleName"`
	LastName                         *string                         `json:"lastName"`
	Aliases                          []string                        `json:"aliases,omitempty"`
	DateOfBirths                     *string                         `json:"dateOfBirths"`
	Address                          *string                         `json:"address"`
	CountryName                      *string                         `json:"countryName"`
	CountryCode                      *string                         `json:"countryCode"`
	Street1                          *string                         `json:"street1"`
	Street2                          *string                         `json:"street2"`
	City                             *string                         `json:"city"`
	State                            *string                         `json:"state"`
	ZIPCode                          *string                         `json:"zipCode"`
	County                           *string                         `json:"county"`
	ConvictionType                   *string                         `json:"convictionType"`
	Offenses                         []ScreeningAPIOffense           `json:"offenses"`
	PhotoURL                         *string                         `json:"photoUrl"`
	Source                           *string                         `json:"source"`
	VerificationResult               *ScreeningAPIVerificationResult `json:"verificationResult"`
	DriversLicenseVerificationResult *ScreeningAPIDLVResult          `json:"driversLicenseVerificationResult"`
}

type ScreeningAPIOffense struct {
	Title                *string
	Class                *string
	Code                 *string
	Section              *string
	Description          *string
	CaseNumber           *string
	Jurisdiction         *string
	AgeOfVictim          *string
	AdmissionDate        *string
	ArrestingAgency      *string
	Category             *string
	ChargeFilingDate     *string
	ClosedDate           *string
	Counts               *string
	Court                *string
	DateConvicted        *string
	DateOfCrime          *string
	DateOfWarrant        *string
	Disposition          *string
	DispositionDate      *string
	Facility             *string
	PrisonerNumber       *string
	RelationshipToVictim *string
	ReleaseDate          *string
	Sentence             *string
	SentenceDate         *string
	SexOfVictim          *string
	Subsection           *string
	WarrantDate          *string
	WarrantNumber        *string
	WeaponsUsed          *string
}

type ScreeningAPIVerificationResult struct {
	Verified bool                         `json:"verified"`
	Data     ScreeningAPIVerificationData `json:"data"`
}

type ScreeningAPIVerificationData struct {
	WorkflowOutcome           ScreeningAPIMessage     `json:"workflowOutcome"`
	PrimaryResult             ScreeningAPIMessage     `json:"primaryResult"`
	TransactionDetail         ScreeningAPITransaction `json:"transactionDetail"`
	AddressVerificationResult ScreeningAPIMessage     `json:"addressVerificationResult"`
	AddressUnitMismatchResult ScreeningAPIMessage     `json:"addressUnitMismatchResult"`
	AddressTypeResult         ScreeningAPIMessage     `json:"addressTypeResult"`
	AddressHighRiskResult     ScreeningAPIMessage     `json:"addressHighRiskResult"`
	DriverLicenseResult       ScreeningAPIMessage     `json:"driverLicenseResult"`
	SSNResult                 ScreeningAPIMessage     `json:"ssnResult"`
	SSNDetail                 ScreeningAPISSNDetail   `json:"ssnDetail"`
	SSNFinderDetails          []string                `json:"ssnFinderDetails"`
	DateOfBirthResult         ScreeningAPIMessage     `json:"dateOfBirthResult"`
	DateOfBirth               *string                 `json:"dateOfBirth"`
	OFACValidation            ScreeningAPIOFACResult  `json:"ofacValidation"`
	Questions                 []ScreeningAPIQuestion  `json:"questions"`
}

type ScreeningAPIMessage struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type ScreeningAPITransaction struct {
	ID                string   `json:"id"`
	Date              string   `json:"date"`
	CustomerReference string   `json:"customerReference"`
	Errors            []string `json:"errors"`
	Warnings          []string `json:"warnings"`
}

type ScreeningAPISSNDetail struct {
	FirstName         string              `json:"firstName"`
	MiddleInitial     string              `json:"middleInitial"`
	LastName          string              `json:"lastName"`
	Street            string              `json:"street"`
	City              string              `json:"city"`
	State             string              `json:"state"`
	ZIPCode           string              `json:"zipCode"`
	ZIPPlusFour       string              `json:"zipPlusFour"`
	AreaCode          string              `json:"areaCode"`
	Phone             string              `json:"phone"`
	DateOfBirth       *string             `json:"dateOfBirth"`
	DateOfBirthResult ScreeningAPIMessage `json:"dateOfBirthResult"`
	ReportedDate      ScreeningAPIDate    `json:"reportedDate"`
	LastTouchedDate   ScreeningAPIDate    `json:"lastTouchedDate"`
}

type ScreeningAPIDate struct {
	Day   string `json:"day"`
	Month string `json:"month"`
	Year  string `json:"year"`
}

type ScreeningAPIOFACResult struct {
	OFACValidationResult ScreeningAPIMessage `json:"ofacValidationResult"`
	OFACRecord           string              `json:"ofacRecord"`
}

type ScreeningAPIQuestion struct {
	Text         string               `json:"text"`
	QuestionType uint                 `json:"questionType"`
	Answers      []ScreeningAPIAnswer `json:"answers"`
}

type ScreeningAPIAnswer struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
}

type ScreeningAPIDLVResult struct {
	DocumentCategoryMatch               *bool `json:"documentCategoryMatch"`
	PersonLastNameExactMatch            *bool `json:"personLastNameExactMatch"`
	PersonLastNameFuzzyPrimaryMatch     *bool `json:"personLastNameFuzzyPrimaryMatch"`
	PersonLastNameFuzzyAlternateMatch   *bool `json:"personLastNameFuzzyAlternateMatch"`
	PersonFirstNameExactMatch           *bool `json:"personFirstNameExactMatch"`
	PersonFirstNameFuzzyPrimaryMatch    *bool `json:"personFirstNameFuzzyPrimaryMatch"`
	PersonFirstNameFuzzyAlternateMatch  *bool `json:"personFirstNameFuzzyAlternateMatch"`
	PersonMiddleNameExactMatch          *bool `json:"personMiddleNameExactMatch"`
	PersonMiddleNameFuzzyPrimaryMatch   *bool `json:"personMiddleNameFuzzyPrimaryMatch"`
	PersonMiddleNameFuzzyAlternateMatch *bool `json:"personMiddleNameFuzzyAlternateMatch"`
	PersonMiddleInitialMatch            *bool `json:"personMiddleInitialMatch"`
	PersonBirthDateMatch                *bool `json:"personBirthDateMatch"`
	DriverLicenseIssueDateMatch         *bool `json:"driverLicenseIssueDateMatch"`
	DriverLicenseExpirationDateMatch    *bool `json:"driverLicenseExpirationDateMatch"`
	DriverLicenseNumberMatch            *bool `json:"driverLicenseNumberMatch"`
	AddressLine1Match                   *bool `json:"addressLine1Match"`
	AddressLine2Match                   *bool `json:"addressLine2Match"`
	AddressCityMatch                    *bool `json:"addressCityMatch"`
	AddressStateCodeMatch               *bool `json:"addressStateCodeMatch"`
	AddressZIP5Match                    *bool `json:"addressZIP5Match"`
	AddressZIP4Match                    *bool `json:"addressZIP4Match"`
	PersonSexCodeMatch                  *bool `json:"personSexCodeMatch"`
}

func NewScreeningAPI(token string) (ScreeningAPI, error) {
	if token == "" {
		return ScreeningAPI{}, errors.New("token can't be empty")
	}

	return ScreeningAPI{
		token: token,
	}, nil
}

// ACTIONS

func (s *ScreeningAPI) ScreenDL(state, licenseNumber, firstName, lastName string) (ScreeningAPIResult, error) {
	if state == "" || licenseNumber == "" || firstName == "" || lastName == "" {
		return ScreeningAPIResult{}, errors.New("all arguments are required to have a non-empty value")
	}

	result, err := s.Screen(ScreeningAPIRequest{
		State:               state,
		DriverLicenseNumber: licenseNumber,
		FirstName:           firstName,
		LastName:            lastName,
		Services: []string{
			ServiceDLDV,
		},
	})

	if err != nil {
		return ScreeningAPIResult{}, err
	}
	if !result[0].Success && !strings.Contains(*result[0].Error, "is not yet supported") {
		err = errors.New(*result[0].Error)
	}

	return result[0], err
}

func (s *ScreeningAPI) ScreenIdentiFraud(firstName, lastName, address, city, state, zip string) (ScreeningAPIResult, error) {
	if firstName == "" || lastName == "" || address == "" || city == "" || state == "" || zip == "" {
		return ScreeningAPIResult{}, errors.New("all arguments are required to have a non-empty value")
	}

	result, err := s.Screen(ScreeningAPIRequest{
		FirstName: firstName,
		LastName:  lastName,
		Address:   address,
		City:      city,
		State:     state,
		ZIP:       zip,
		Services: []string{
			ServiceIdentiFraud,
		},
	})

	if err != nil {
		return ScreeningAPIResult{}, err
	}
	if !result[0].Success {
		err = errors.New(*result[0].Error)
	}

	return result[0], err
}

func (s *ScreeningAPI) ScreenSexOffender(firstName, lastName, dob string) ([]ScreeningAPIResult, error) {
	if firstName == "" || lastName == "" || dob == "" {
		return []ScreeningAPIResult{}, errors.New("all arguments are required to have a non-empty value")
	}

	return s.Screen(ScreeningAPIRequest{
		FirstName:    firstName,
		LastName:     lastName,
		DateOfBirths: dob,
		Services: []string{
			ServiceSexOffender,
		},
	})
}

func (s *ScreeningAPI) ScreenCrimeRecord(firstName, lastName, dob string) ([]ScreeningAPIResult, error) {
	if firstName == "" || lastName == "" || dob == "" {
		return []ScreeningAPIResult{}, errors.New("all arguments are required to have a non-empty value")
	}

	return s.Screen(ScreeningAPIRequest{
		FirstName:    firstName,
		LastName:     lastName,
		DateOfBirths: dob,
		Services: []string{
			ServiceCrimeRecord,
		},
	})
}

func (s *ScreeningAPI) ScreenPEP(firstName, lastName, dob string) ([]ScreeningAPIResult, error) {
	if firstName == "" || lastName == "" || dob == "" {
		return []ScreeningAPIResult{}, errors.New("all arguments are required to have a non-empty value")
	}

	return s.Screen(ScreeningAPIRequest{
		FirstName:    firstName,
		LastName:     lastName,
		DateOfBirths: dob,
		Services: []string{
			ServicePEP,
		},
	})
}

func (s *ScreeningAPI) ScreenOFAC(firstName, lastName, dob string) ([]ScreeningAPIResult, error) {
	if firstName == "" || lastName == "" || dob == "" {
		return []ScreeningAPIResult{}, errors.New("all arguments are required to have a non-empty value")
	}

	return s.Screen(ScreeningAPIRequest{
		FirstName:    firstName,
		LastName:     lastName,
		DateOfBirths: dob,
		Services: []string{
			ServiceOFAC,
		},
	})
}

func (s *ScreeningAPI) ScreenEveryPolitician(firstName, lastName, dob string) ([]ScreeningAPIResult, error) {
	if firstName == "" || lastName == "" || dob == "" {
		return []ScreeningAPIResult{}, errors.New("all arguments are required to have a non-empty value")
	}

	return s.Screen(ScreeningAPIRequest{
		FirstName:    firstName,
		LastName:     lastName,
		DateOfBirths: dob,
		Services: []string{
			ServiceEveryPolitician,
		},
	})
}

func (s *ScreeningAPI) Screen(request ScreeningAPIRequest) ([]ScreeningAPIResult, error) {
	payload, _ := json.Marshal(request)

	api, _ := http.NewRequest("POST", "https://screening.idware.net/api/Check", bytes.NewBuffer(payload))
	api.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.token))
	if response, err := http.DefaultClient.Do(api); err != nil {
		return []ScreeningAPIResult{}, err
	} else {
		var result []ScreeningAPIResult

		body, _ := io.ReadAll(response.Body)
		json.Unmarshal(body, &result)

		return result, nil
	}
}
