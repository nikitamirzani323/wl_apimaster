package entities

type Model_company struct {
	Company_idcomp      string `json:"company_idcompany"`
	Company_startjoin   string `json:"company_startjoin"`
	Company_endjoin     string `json:"company_endjoin"`
	Company_idcurr      string `json:"company_idcurr"`
	Company_nmcompany   string `json:"company_nmcompany"`
	Company_nmowner     string `json:"company_nmowner"`
	Company_phoneowner  string `json:"company_phoneowner"`
	Company_emailowner  string `json:"company_emailowner"`
	Company_urlendpoint string `json:"company_urlendpoint"`
	Company_status      string `json:"company_status"`
	Company_create      string `json:"company_create"`
	Company_update      string `json:"company_update"`
}
type Model_compcurr struct {
	Curr_idcurr string `json:"curr_idcurr"`
}
type Controller_company struct {
	Master string `json:"master" validate:"required"`
}

type Controller_companysave struct {
	Sdata               string `json:"sdata" validate:"required"`
	Page                string `json:"page" validate:"required"`
	Company_idcomp      string `json:"company_idcompany" validate:"required"`
	Company_idcurr      string `json:"company_idcurr" validate:"required"`
	Company_nmcompany   string `json:"company_nmcompany" validate:"required"`
	Company_nmowner     string `json:"company_nmowner" validate:"required"`
	Company_phoneowner  string `json:"company_phoneowner" validate:"required"`
	Company_emailowner  string `json:"company_emailowner" `
	Company_urlendpoint string `json:"company_urlendpoint" validate:"required"`
	Company_status      string `json:"company_status" validate:"required"`
}
