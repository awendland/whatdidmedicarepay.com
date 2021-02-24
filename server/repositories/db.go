package repositories

import (
	"context"
	"database/sql"

	"github.com/awendland/whatdidmedicarepay.com/server/config"
)

// PUPDEntry is a model representing rows in the Provider Utilization and Payment Data
// Public Use File provided by CMS.
type PUPDEntry struct {
	NPI                            string
	NPPESProviderLastOrgName       string
	NPPESProviderFirstName         string
	NPPESProviderMi                string
	NPPESCredentials               string
	NPPESProviderGender            string
	NPPESEntityCode                string
	NPPESProviderStreet1           string
	NPPESProviderStreet2           string
	NPPESProviderCity              string
	NPPESProviderZip               string
	NPPESProviderState             string
	NPPESProviderCountry           string
	ProviderType                   string
	MedicareParticipationIndicator string
	PlaceOfService                 string
	HCPCSCode                      string
	HCPCSDescription               string
	HCPCSDrugIndicator             string
	LineSrvcCnt                    string
	BeneUniqueCnt                  string
	BeneDaySrvcCnt                 string
	AverageMedicareAllowedAmt      string
	AverageSubmittedChrgAmt        string
	AverageMedicarePaymentAmt      string
	AverageMedicareStandardAmt     string
}

// SearchPUPDEntries queries the DB using a full-text search MATCH operator to
// find relevant Provider Utilization and Payment Data entries.
//
// See the schema specified by `:data/prepare.py` to understand how the full-text search DB is
// structured. Any valid fts query can be submitted to this method.
func SearchPUPDEntries(ctx context.Context, config *config.Config, query string, limit int) ([]PUPDEntry, error) {
	rows, err := config.DB.QueryContext(ctx, `
	SELECT
		PUP_data."npi",
		PUP_data."nppes_provider_last_org_name",
		PUP_data."nppes_provider_first_name",
		PUP_data."nppes_provider_mi",
		PUP_data."nppes_credentials",
		PUP_data."nppes_provider_gender",
		PUP_data."nppes_entity_code",
		PUP_data."nppes_provider_street1",
		PUP_data."nppes_provider_street2",
		PUP_data."nppes_provider_city",
		PUP_data."nppes_provider_zip",
		PUP_data."nppes_provider_state",
		PUP_data."nppes_provider_country",
		PUP_data."provider_type",
		PUP_data."medicare_participation_indicator",
		PUP_data."place_of_service",
		PUP_data."hcpcs_code",
		PUP_data."hcpcs_description",
		PUP_data."hcpcs_drug_indicator",
		PUP_data."line_srvc_cnt",
		PUP_data."bene_unique_cnt",
		PUP_data."bene_day_srvc_cnt",
		PUP_data."average_Medicare_allowed_amt",
		PUP_data."average_submitted_chrg_amt",
		PUP_data."average_Medicare_payment_amt",
		PUP_data."average_Medicare_standard_amt"
	FROM PUP_data
	JOIN PUP_data_fts ON PUP_data_fts.rowid = PUP_data.rowid
	WHERE PUP_data_fts MATCH :query
	LIMIT :limit
	`, sql.Named("query", query), sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	var entries []PUPDEntry
	for rows.Next() {
		var p PUPDEntry
		if err = rows.Scan(
			&p.NPI,
			&p.NPPESProviderLastOrgName,
			&p.NPPESProviderFirstName,
			&p.NPPESProviderMi,
			&p.NPPESCredentials,
			&p.NPPESProviderGender,
			&p.NPPESEntityCode,
			&p.NPPESProviderStreet1,
			&p.NPPESProviderStreet2,
			&p.NPPESProviderCity,
			&p.NPPESProviderZip,
			&p.NPPESProviderState,
			&p.NPPESProviderCountry,
			&p.ProviderType,
			&p.MedicareParticipationIndicator,
			&p.PlaceOfService,
			&p.HCPCSCode,
			&p.HCPCSDescription,
			&p.HCPCSDrugIndicator,
			&p.LineSrvcCnt,
			&p.BeneUniqueCnt,
			&p.BeneDaySrvcCnt,
			&p.AverageMedicareAllowedAmt,
			&p.AverageSubmittedChrgAmt,
			&p.AverageMedicarePaymentAmt,
			&p.AverageMedicareStandardAmt,
		); err != nil {
			break
		}
		entries = append(entries, p)
	}
	return entries, nil
}
