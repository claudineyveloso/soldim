package types

type Product struct {
	ID             int64       `json:"id"`
	Nome           string      `json:"nome"`
	Codigo         string      `json:"codigo"`
	Preco          interface{} `json:"preco"`
	ImagemURL      string      `json:"imagemURL"`
	Tipo           string      `json:"tipo"`
	Situacao       string      `json:"situacao"`
	Formato        string      `json:"formato"`
	DescricaoCurta string      `json:"descricaoCurta"`
}

type ProductWrapper struct {
	Produto Product `json:"produto"`
}

//	type ProductResponse struct {
//		Retorno struct {
//			Products []ProductWrapper `json:"produtos"`
//			Total    int              `json:"total"`
//			Limit    int              `json:"limit"`
//		} `json:"retorno"`
//	}
type ProductResponse struct {
	Data  []Product `json:"data"`
	Total int       `json:"total"`
	Limit int       `json:"limit"`
}
