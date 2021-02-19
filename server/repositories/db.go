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
		puf_2017."npi",
		puf_2017."nppes_provider_last_org_name",
		puf_2017."nppes_provider_first_name",
		puf_2017."nppes_provider_mi",
		puf_2017."nppes_credentials",
		puf_2017."nppes_provider_gender",
		puf_2017."nppes_entity_code",
		puf_2017."nppes_provider_street1",
		puf_2017."nppes_provider_street2",
		puf_2017."nppes_provider_city",
		puf_2017."nppes_provider_zip",
		puf_2017."nppes_provider_state",
		puf_2017."nppes_provider_country",
		puf_2017."provider_type",
		puf_2017."medicare_participation_indicator",
		puf_2017."place_of_service",
		puf_2017."hcpcs_code",
		puf_2017."hcpcs_description",
		puf_2017."hcpcs_drug_indicator",
		puf_2017."line_srvc_cnt",
		puf_2017."bene_unique_cnt",
		puf_2017."bene_day_srvc_cnt",
		puf_2017."average_Medicare_allowed_amt",
		puf_2017."average_submitted_chrg_amt",
		puf_2017."average_Medicare_payment_amt",
		puf_2017."average_Medicare_standard_amt"
	FROM puf_2017
	JOIN data ON data.rowid = puf_2017.rowid
	WHERE data MATCH :query
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
