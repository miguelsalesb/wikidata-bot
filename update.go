package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/net/publicsuffix"
)

type results struct {
	Batchcomplete string `json:"batchcomplete"`
	Query         query  `json:"query"`
}

type query struct {
	Tokens tokens `json:"tokens"`
}

type tokens struct {
	Logintoken string `json:"logintoken"`
}

type data struct {
	Data   string
	Labels labels `json:"labels"`
}

type resultsCsfr struct {
	Batchcomplete string    `json:"batchcomplete"`
	Query         queryCsfr `json:"query"`
}

type queryCsfr struct {
	Tokens tokensCsfr `json:"tokens"`
}

type tokensCsfr struct {
	Csrftoken string `json:"csrftoken"`
}

type dataCsfr struct {
	Data   string
	Labels labels `json:"labels"`
}

type labels struct {
	Pt pt `json:"pt"`
	En en `json:"en"`
}

type pt struct {
	Pt       string
	Language string `json:"language"`
	Value    string `json:"value"`
}

type en struct {
	En       string
	Language string `json:"language"`
	Value    string `json:"value"`
}

type Entity struct {
	Label       Labels       `json:"labels,omitempty"`
	Alias       *Aliases     `json:"aliases,omitempty"`
	Description Descriptions `json:"descriptions,omitempty"`
	// Claim       *Claims       `json:"claims,omitempty"`
	Id       string `json:"id,omitempty"`
	Type     string `json:"type,omitempty"`
	Datatype string `json:"datatype,omitempty"`
	// Lastrevid   int            `json:"lastrevid,omitempty"`
}

// type Claims struct {
// 	Type     string `json:"type,omitempty"`
// 	Datatype string `json:"datatype,omitempty"`
// }

type Aliases struct {
	Pt []Pts `json:"pt,omitempty"`
}

type Labels struct {
	Pt Pts `json:"pt,omitempty"`
}

type Descriptions struct {
	Pt Pts `json:"pt,omitempty"`
}

type Pts struct {
	Language string `json:"language,omitempty"`
	Value    string `json:"value,omitempty"`
}
type Ens struct {
	Language string `json:"language,omitempty"`
	Value    string `json:"value,omitempty"`
}

type object struct {
	id string
}

func main() {

	// Wikidata Bot data
	const username = "Msalesb@bot"
	const password = "u0t6u267juiomiru9idmbk3du47sb5d9"

	var entity Entity
	var replacerWiki = strings.NewReplacer("\\", "", "\"[", "[", "]\"", "]")
	var newType = "property"

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{Jar: jar}

	res, err := client.PostForm("http://127.0.0.1:8181/api.php?", url.Values{
		"action": {"query"},
		"meta":   {"tokens"},
		"type":   {"login"},
		"format": {"json"},
	})

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\ndataLogin: %s", body)

	data := results{}

	bBody := []byte(body)
	bUn := json.Unmarshal(bBody, &data)
	if bUn != nil {
		fmt.Println(bUn)
	}

	loginToken := string(data.Query.Tokens.Logintoken)
	fmt.Printf("\nloginToken: %v", loginToken)

	res, err = client.PostForm("http://127.0.0.1:8181/api.php?", url.Values{
		"action":     {"login"},
		"lgname":     {username},
		"lgpassword": {password},
		"lgtoken":    {loginToken},
		"format":     {"json"},
	})
	if err != nil {
		log.Fatal(err)
	}

	dataLogin, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	fmt.Printf("\nConnection: %v", string(dataLogin))

	resp, err := client.PostForm("http://127.0.0.1:8181/api.php?", url.Values{
		"action": {"query"},
		"meta":   {"tokens"},
		"type":   {"csrf"},
		"format": {"json"},
	})

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyCsfr, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nVER: %s", bodyCsfr)

	dataCsfr := resultsCsfr{}

	bBodyCsfr := []byte(bodyCsfr)
	bUnCsfr := json.Unmarshal(bBodyCsfr, &dataCsfr)
	if bUnCsfr != nil {
		fmt.Println(bUnCsfr)
	}

	tokenCsfr := string(dataCsfr.Query.Tokens.Csrftoken)
	fmt.Printf("\ntokenCsfr: %s", tokenCsfr)
	// wikibase-entityid
	var initialPropertiesAndEntities = [][]string{{"data de nascimento", "data na qual o sujeito nasceu", "P1", "property", "time"},
		{"data de morte", "data em que o sujeito morreu", "P2", "property", "time"},
		{"afirmado em", "a ser utilizada no campo de referências...", "P3", "property", "time"},
		{"endereço eletrónico da referência", "endereço eletrónico onde a referência está depositada", "P4", "property", "url"},
		{"instância de", "este item é uma instância deste outro item", "P5", "property", "wikibase-item"},
		{"identificador PTBNP", "identificador da Biblioteca Nacional de Portugal", "P6", "property", "wikibase-item"},
		{"data de acesso", "data em que a informação foi obtida de uma base de dados ou website", "P7", "property", "time"},
		{"obra destacada", "obras de arte ou trabalhos científicos notáveis do sujeito", "P8", "property", "wikibase-item"},
		{"data de publicação", "data em que a obra foi publicada ou lançad", "P9", "property", "time"},
		{"país de origem", "país de origem do sujeito ou obra", "P10", "property", "wikibase-item"},
		{"país de nacionalidade", "Estado soberano do qual a pessoa é nacional", "P11", "property", "wikibase-item"},
		{"ser humano", "espécie de hominídeo", "Q1", "item", ""},
		{"Biblioteca Nacional de Portugal", "biblioteca nacional de Portugal, depositária do maior património bibliográfico do país", "Q2", "item", ""},
		{"obra escrita", "qualquer obra criativa expressada por meio da escrita, como inscrições...", "Q3", "item", ""},
		{"Portugal", "país europeu", "Q4", "item", ""},
	}

	// fmt.Println("TAMANHO: ", len(initialPropertiesAndEntities[1][0]))

	for x := 0; x < len(initialPropertiesAndEntities); x++ {
		// fmt.Println("\ninitialPropertiesAndEntities[0][x]", initialPropertiesAndEntities[x])
		entity.Label.Pt =
			Pts{
				"pt",
				initialPropertiesAndEntities[x][0],
			}

		entity.Alias = nil

		entity.Description.Pt.Language = "pt"

		entity.Description.Pt.Value = initialPropertiesAndEntities[x][1]

		entity.Id = initialPropertiesAndEntities[x][2]

		entity.Type = initialPropertiesAndEntities[x][3]

		entity.Datatype = initialPropertiesAndEntities[x][4]

		t, _ := json.Marshal(entity)
		exportToWikidata := (string(t))

		// convert the result from []byte to a string
		exportToWikidata = replacerWiki.Replace((exportToWikidata))

		if x == 11 {
			newType = "item"
		}

		fmt.Printf("\n\n\nNEWTYPE: %v", newType)
		fmt.Printf("\n\n\nexportToWikidata: %v", string(exportToWikidata))

		resp, err = client.PostForm("http://127.0.0.1:8181/api.php?", url.Values{
			"action": {"wbeditentity"},
			"new":    {newType},
			"token":  {tokenCsfr},
			"bot":    {"1"},
			"data":   {exportToWikidata},
			"format": {"json"},
		})

		if err != nil {
			log.Fatal(err)
		}

		dataPost, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\n\n\nEntidade: %v", string(dataPost))

	}

	//  Resultou:
	// *********** MAS A INSTÂNCIA DE FOI ACRESCENTADA AO FINAL DA ENTIDADE
	// Entidade: {"pageinfo":{"lastrevid":146},"success":1,"claim":{"mainsnak":{"snaktype":"value","property":"P3",
	// "hash":"4f44c4bd4a2a1774d8079832e9b2098d3e23a994","datavalue":{"value":{"entity-type":"item","numeric-id":17,
	// "id":"Q17"},"type":"wikibase-entityid"},"datatype":"wikibase-item"},"type":"statement",
	// "id":"Q699$CDF17061-B5B4-4B9E-8CD0-A35C2EF7B983","rank":"normal"}}
	// `{"claims": {"type":"property","datatype":"time","id":"P1","labels":{"pt":{"language":"pt","value":"data de nascimento"}},"descriptions":{"pt":{"language":"pt","value":"data na qual o sujeito nasceu"}}}}`
	// ,
	// {"type":"property","datatype":"time","id":"P2","labels":{"pt":{"language":"pt","value":"data de morte"}},"descriptions":{"pt":{"language":"pt","value":"data em que o sujeito morreu"}},"aliases":{},"claims":{},"lastrevid":3},
	// {"type":"property","datatype":"wikibase-item","id":"P3","labels":{"pt":{"language":"pt","value":"afirmado em"}},"descriptions":{"pt":{"language":"pt","value":"a ser utilizada no campo de refer\u00eancias, para a fonte de informa\u00e7\u00e3o em que a afirma\u00e7\u00e3o \u00e9 feita; para qualificadores utilize P805"}},"aliases":{},"claims":{},"lastrevid":4},
	// {"type":"property","datatype":"url","id":"P4","labels":{"pt":{"language":"pt","value":"endere\u00e7o eletr\u00f3nico da refer\u00eancia"}},"descriptions":{"pt":{"language":"pt","value":"endere\u00e7o eletr\u00f3nico onde a refer\u00eancia est\u00e1 depositada"}},"aliases":{},"claims":{},"lastrevid":5},
	// {"type":"property","datatype":"wikibase-item","id":"P5","labels":{"pt":{"language":"pt","value":"inst\u00e2ncia de"}},"descriptions":{"pt":{"language":"pt","value":"este item \u00e9 uma inst\u00e2ncia deste outro item"}},"aliases":{},"claims":{},"lastrevid":6},
	// {"type":"property","datatype":"external-id","id":"P6","labels":{"pt":{"language":"pt","value":"identificador PTBNP"}},"descriptions":{"pt":{"language":"pt","value":"identificador da Biblioteca Nacional de Portugal"}},"aliases":{},"claims":{},"lastrevid":7},
	// {"type":"property","datatype":"time","id":"P7","labels":{"pt":{"language":"pt","value":"data de acesso"}},"descriptions":{"pt":{"language":"pt","value":"data em que a informa\u00e7\u00e3o foi obtida de uma base de dados ou website (para uso em fontes online)"}},"aliases":{},"claims":{},"lastrevid":8},
	// {"type":"property","datatype":"wikibase-item","id":"P8","labels":{"pt":{"language":"pt","value":"obra destacada"}},"descriptions":{"pt":{"language":"pt","value":"obras de arte ou trabalhos cient\u00edficos not\u00e1veis do sujeito"}},"aliases":{},"claims":{},"lastrevid":9},
	// {"type":"property","datatype":"time","id":"P9","labels":{"pt":{"language":"pt","value":"data de publica\u00e7\u00e3o"}},"descriptions":{"pt":{"language":"pt","value":"data em que a obra foi publicada ou lan\u00e7ada"}},"aliases":{},"claims":{},"lastrevid":10},
	// {"type":"property","datatype":"wikibase-item","id":"P10","labels":{"pt":{"language":"pt","value":"pa\u00eds de origem"}},"descriptions":{"pt":{"language":"pt","value":"pa\u00eds de origem do sujeito ou obra"}},"aliases":{},"claims":{},"lastrevid":11},
	// {"type":"item","id":"Q1","labels":{"pt":{"language":"pt","value":"ser humano"}},"descriptions":{"pt":{"language":"pt","value":"esp\u00e9cie de homin\u00eddeo"}},"aliases":{},"claims":{},"sitelinks":{},"lastrevid":12},
	// {"type":"item","id":"Q2","labels":{"pt":{"language":"pt","value":"Biblioteca Nacional de Portugal"}},"descriptions":{"pt":{"language":"pt","value":"biblioteca nacional de Portugal, deposit\u00e1ria do maior patrim\u00f3nio bibliogr\u00e1fico do pa\u00eds"}},"aliases":{},"claims":{},"sitelinks":{},"lastrevid":13},
	// {"type":"item","id":"Q3","labels":{"pt":{"language":"pt","value":"obra escrita"}},"descriptions":{"pt":{"language":"pt","value":"qualquer obra criativa expressada por meio da escrita, como inscri\u00e7\u00f5es, manuscritos, documentos ou mapas. Usar Q7725634 para obras liter\u00e1rias acad\u00e9micas, est\u00e9ticas ou recreativas"}},"aliases":{},"claims":{},"sitelinks":{},"lastrevid":14},
	// {"type":"item","id":"Q4","labels":{"pt":{"language":"pt","value":"Portugal"}},"descriptions":{"pt":{"language":"pt","value":"pa\u00eds europeu"}},"aliases":{},"claims":{},"sitelinks":{},"lastrevid":15},`

	// convert the struct to JSON

}
