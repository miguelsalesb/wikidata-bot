package wiki

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/net/publicsuffix"
)

type C_object struct {
	Entity  Entity `json:"entity"`
	Success int    `json:"success"`
}

type Entity struct {
	Label    Labels `json:"labels,omitempty"`
	Datatype string `json:"datatype,omitempty"`
	Id       string `json:"id,omitempty"`
}

type Aliases struct {
	Pt []Pts `json:"pt,omitempty"`
}

type Labels struct {
	Pt Pts `json:"pt,omitempty"`
}

type Descriptions struct {
	Pt Pts `json:"pt,omitempty"`
	En Ens `json:"en,omitempty"`
}

type Pts struct {
	Language string `json:"language,omitempty"`
	Value    string `json:"value,omitempty"`
}
type Ens struct {
	Language string `json:"language,omitempty"`
	Value    string `json:"value,omitempty"`
}

func CreateFirst() {

	var entity Entity
	var tokenCsfr string
	var replacerWiki = strings.NewReplacer("\\", "", "\"[", "[", "]\"", "]")
	// var newType string

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{Jar: jar}

	tokenCsfr = ConnectToWikidata(client)

	var initialPropertiesAndEntities = [][]string{
		{"data de nascimento", "data na qual o sujeito nasceu", "P1", "property", "time"},
		{"data de morte", "data em que o sujeito morreu", "P2", "property", "time"},
		{"afirmado em", "a ser utilizada no campo de referências...", "P3", "property", "wikibase-item"},
		{"endereço eletrónico da referência", "endereço eletrónico onde a referência está depositada", "P4", "property", "url"},
		{"instância de", "este item é uma instância deste outro item", "P5", "property", "wikibase-item"},
		{"identificador PTBNP", "identificador da Biblioteca Nacional de Portugal", "P6", "property", "external-id"},
		{"data de acesso", "data em que a informação foi obtida de uma base de dados ou website", "P7", "property", "time"},
		{"obra destacada", "obras de arte ou trabalhos científicos notáveis do sujeito", "P8", "property", "wikibase-item"},
		{"data de publicação", "data em que a obra foi publicada ou lançado", "P9", "property", "time"},
		{"país de origem", "país de origem do sujeito ou obra", "P10", "property", "wikibase-item"},
		{"país de nacionalidade", "Estado soberano do qual a pessoa é nacional", "P11", "property", "wikibase-item"},
		{"ser humano", "espécie de hominídeo", "Q1", "item"},
		{"Biblioteca Nacional de Portugal", "biblioteca nacional de Portugal, depositária do maior património bibliográfico do país", "Q2", "item"},
		{"obra escrita", "qualquer obra criativa expressada por meio da escrita, como inscrições...", "Q3", "item"},
		{"Portugal", "país europeu", "Q4", "item"},
	}

	fmt.Println(`
creating initial entities and properties... 

P1 = date of birth (P569)
P2 = date of death (P570)
P3 afirmado em - stated in (P248)
P4 endereço eletrónico da referência - reference URL (P854)
P5 instância de (P5) - instance of (P31)
P6 identificador PTBNP - Portuguese National Library ID (P1005)
P7 data de acesso (P7) - retrieved (P813)
P8 obra destacada - notable work (P800)
P9 data de publicação - publication date (P577)
P10 país de origem - country of origin (P495)
P11 país de nacionalidade - country of citizenship (P27)
Q1 ser humano (Q1) - human (Q5)
Q2 BNP - National Library of Portugal (Q245966)
Q3 obra escrita - written work (Q47461344)
Q4 Portugal - Portugal (Q45)
	`)

	time.Sleep(2 * time.Second)

	for x := 0; x < len(initialPropertiesAndEntities); x++ {

		time.Sleep(500 * time.Millisecond)

		entity.Label.Pt =
			Pts{
				"pt",
				initialPropertiesAndEntities[x][0],
			}

		if initialPropertiesAndEntities[x][3] == "property" {
			// 	// newType = "property"
			entity.Datatype = initialPropertiesAndEntities[x][4]
		}

		t, _ := json.Marshal(entity)
		exportToWikidata := (string(t))

		// convert the result from []byte to a string
		exportToWikidata = replacerWiki.Replace((exportToWikidata))

		resp, err := client.PostForm("http://127.0.0.1:8181/api.php?", url.Values{
			"action": {"wbeditentity"},
			"new":    {initialPropertiesAndEntities[x][3]},
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

		data := C_object{}

		errP := json.Unmarshal(dataPost, &data) // decode the JSON data
		if errP != nil {
			fmt.Println(errP)
		}

		QID := data.Entity.Id
		fmt.Printf("\n\nCreated: %v", QID)

	}
}
