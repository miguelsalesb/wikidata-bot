package wiki

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type C_results struct {
	C_Batchcomplete string  `json:"batchcomplete"`
	C_Query         C_query `json:"query"`
}

type C_query struct {
	C_Tokens C_tokens `json:"tokens"`
}

type C_tokens struct {
	C_Logintoken string `json:"logintoken"`
}

type C_resultsCsfr struct {
	C_Batchcomplete string      `json:"batchcomplete"`
	C_Query         C_queryCsfr `json:"query"`
}

type C_queryCsfr struct {
	C_Tokens C_tokensCsfr `json:"tokens"`
}

type C_tokensCsfr struct {
	C_Csrftoken string `json:"csrftoken"`
}

func ConnectToWikidata(client http.Client) string {
	// Wikidata Bot data
	const username = "Msalesb@bot"
	const password = "0i5cbcnrf92102r2skiq62o3a3up66oi"

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

	data := C_results{}

	bBody := []byte(body)
	bUn := json.Unmarshal(bBody, &data)
	if bUn != nil {
		fmt.Println(bUn)
	}

	loginToken := string(data.C_Query.C_Tokens.C_Logintoken)

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

	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()

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

	dataCsfr := C_resultsCsfr{}

	bBodyCsfr := []byte(bodyCsfr)
	bUnCsfr := json.Unmarshal(bBodyCsfr, &dataCsfr)
	if bUnCsfr != nil {
		fmt.Println(bUnCsfr)
	}

	tokenCsfr := string(dataCsfr.C_Query.C_Tokens.C_Csrftoken)

	return tokenCsfr
}
