package structs

// blob from /retrieve path
type StartramRetrieve struct {
	Action     string `json:"action"`
	Conf       string `json:"conf"`
	Debug      any    `json:"debug"`
	Error      int    `json:"error"`
	Lease      string `json:"lease"`
	Ongoing    int    `json:"ongoing"`
	Pubkey     string `json:"pubkey"`
	Region     string `json:"region"`
	Status     string `json:"status"`
	Subdomains []struct {
		Alias   string `json:"alias"`
		Port    int    `json:"port"`
		Status  string `json:"status"`
		SvcType string `json:"svc_type"`
		URL     string `json:"url"`
	} `json:"subdomains"`
}

// startram region server subobject
type StartramRegion struct {
	Country string `json:"country"`
	Desc    string `json:"desc"`
}

// register a pubkey
type StartramRegister struct {
	RegCode string `json:"reg_code"`
	Pubkey  string `json:"pubkey"`
	Region  string `json:"region"`
}

type StartramRegisterResp struct {
	Action string
	Debug  any    `json:"debug"`
	Error  int    `json:"error"`
	Pubkey string `json:"pubkey"`
	Lease  string `json:"lease"`
	Region string `json:"region"`
}

// for create or delete
type StartramSvc struct {
	Subdomain string `json:"subdomain"`
	Pubkey    string `json:"pubkey"`
	SvcType   string `json:"svc_type"`
}

type StartramSvcResp struct {
	Action    string `json:"action"`
	Debug     any    `json:"debug"`
	Error     int    `json:"error"`
	Subdomain string `json:"subdomain"`
	SvcType   string `json:"svc_type"`
	Pubkey    string `json:"pubkey"`
	Status    string `json:"status"`
	Lease     string `json:"lease"`
}

// for create or delete
type StartramAlias struct {
	Subdomain string `json:"subdomain"`
	Pubkey    string `json:"pubkey"`
	Alias     string `json:"alias"`
}

type StartramAliasResp struct {
	Action    string `json:"action"`
	Debug     any    `json:"debug"`
	Error     int    `json:"error"`
	Subdomain string `json:"subdomain"`
	Alias     string `json:"alias"`
	Pubkey    string `json:"pubkey"`
	Lease     string `json:"lease"`
}
