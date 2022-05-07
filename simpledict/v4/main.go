package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type DictRequest struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
	UserID    string `json:"user_id"`
}

type DictResponse struct {
	Rc   int `json:"rc"`
	Wiki struct {
		KnownInLaguages int `json:"known_in_laguages"`
		Description     struct {
			Source string      `json:"source"`
			Target interface{} `json:"target"`
		} `json:"description"`
		ID   string `json:"id"`
		Item struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"item"`
		ImageURL  string `json:"image_url"`
		IsSubject string `json:"is_subject"`
		Sitelink  string `json:"sitelink"`
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []string      `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}

type DictResponseWithBaidu struct {
	TransResult struct {
		Data []struct {
			Dst        string          `json:"dst"`
			PrefixWrap int             `json:"prefixWrap"`
			Result     [][]interface{} `json:"result"`
			Src        string          `json:"src"`
		} `json:"data"`
		From     string `json:"from"`
		Status   int    `json:"status"`
		To       string `json:"to"`
		Type     int    `json:"type"`
		Phonetic []struct {
			SrcStr string `json:"src_str"`
			TrgStr string `json:"trg_str"`
		} `json:"phonetic"`
	} `json:"trans_result"`
	DictResult struct {
		Edict struct {
			Item []struct {
				TrGroup []struct {
					Tr          []string `json:"tr"`
					Example     []string `json:"example"`
					SimilarWord []string `json:"similar_word"`
				} `json:"tr_group"`
				Pos string `json:"pos"`
			} `json:"item"`
			Word string `json:"word"`
		} `json:"edict"`
		Collins struct {
			Entry []struct {
				EntryID string `json:"entry_id"`
				Type    string `json:"type"`
				Value   []struct {
					MeanType []struct {
						InfoType string `json:"info_type"`
						InfoID   string `json:"info_id"`
						Example  []struct {
							ExampleID string `json:"example_id"`
							TtsSize   string `json:"tts_size"`
							Tran      string `json:"tran"`
							Ex        string `json:"ex"`
							TtsMp3    string `json:"tts_mp3"`
						} `json:"example,omitempty"`
						Posc []struct {
							Tran    string `json:"tran"`
							PoscID  string `json:"posc_id"`
							Example []struct {
								ExampleID string `json:"example_id"`
								Tran      string `json:"tran"`
								Ex        string `json:"ex"`
								TtsMp3    string `json:"tts_mp3"`
							} `json:"example"`
							Def string `json:"def"`
						} `json:"posc,omitempty"`
					} `json:"mean_type"`
					Gramarinfo []struct {
						Tran  string `json:"tran"`
						Type  string `json:"type"`
						Label string `json:"label"`
					} `json:"gramarinfo"`
					Tran   string `json:"tran"`
					Def    string `json:"def"`
					MeanID string `json:"mean_id"`
					Posp   []struct {
						Label string `json:"label"`
					} `json:"posp"`
				} `json:"value"`
			} `json:"entry"`
			WordName      string `json:"word_name"`
			Frequence     string `json:"frequence"`
			WordEmphasize string `json:"word_emphasize"`
			WordID        string `json:"word_id"`
		} `json:"collins"`
		From        string `json:"from"`
		SimpleMeans struct {
			WordName  string   `json:"word_name"`
			From      string   `json:"from"`
			WordMeans []string `json:"word_means"`
			Exchange  struct {
				WordPl []string `json:"word_pl"`
			} `json:"exchange"`
			Tags struct {
				Core  []string `json:"core"`
				Other []string `json:"other"`
			} `json:"tags"`
			Symbols []struct {
				PhEn  string `json:"ph_en"`
				PhAm  string `json:"ph_am"`
				Parts []struct {
					Part  string   `json:"part"`
					Means []string `json:"means"`
				} `json:"parts"`
				PhOther string `json:"ph_other"`
			} `json:"symbols"`
		} `json:"simple_means"`
		Lang   string `json:"lang"`
		Oxford struct {
			Entry []struct {
				Tag  string `json:"tag"`
				Name string `json:"name"`
				Data []struct {
					Tag  string `json:"tag"`
					Data []struct {
						Tag  string `json:"tag"`
						Data []struct {
							Tag  string `json:"tag"`
							Data []struct {
								Tag  string `json:"tag"`
								Data []struct {
									Tag    string `json:"tag"`
									EnText string `json:"enText,omitempty"`
									ChText string `json:"chText,omitempty"`
									G      string `json:"g,omitempty"`
									Data   []struct {
										Text      string `json:"text"`
										HoverText string `json:"hoverText"`
									} `json:"data,omitempty"`
								} `json:"data"`
							} `json:"data"`
						} `json:"data,omitempty"`
						P     string `json:"p,omitempty"`
						PText string `json:"p_text,omitempty"`
						N     string `json:"n,omitempty"`
						Xt    string `json:"xt,omitempty"`
					} `json:"data"`
				} `json:"data"`
			} `json:"entry"`
			Unbox []struct {
				Tag  string `json:"tag"`
				Type string `json:"type"`
				Name string `json:"name"`
				Data []struct {
					Tag     string `json:"tag"`
					Text    string `json:"text,omitempty"`
					Words   string `json:"words,omitempty"`
					Outdent string `json:"outdent,omitempty"`
					Data    []struct {
						Tag    string `json:"tag"`
						EnText string `json:"enText"`
						ChText string `json:"chText"`
					} `json:"data,omitempty"`
				} `json:"data"`
			} `json:"unbox"`
		} `json:"oxford"`
		BaiduPhrase []struct {
			Tit   []string `json:"tit"`
			Trans []string `json:"trans"`
		} `json:"baidu_phrase"`
	} `json:"dict_result"`
	LijuResult struct {
		Double string   `json:"double"`
		Tag    []string `json:"tag"`
		Single string   `json:"single"`
	} `json:"liju_result"`
	Logid int64 `json:"logid"`
}

func query(word string) {
	client := &http.Client{}
	request := DictRequest{TransType: "en2zh", Source: word}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")
	req.Header.Set("os-version", "")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
	req.Header.Set("app-name", "xy")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("device-id", "")
	req.Header.Set("os-type", "web")
	req.Header.Set("X-Authorization", "token:qgemv4jr1y38jyq6vhvi")
	req.Header.Set("Origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", "_ym_uid=16456948721020430059; _ym_d=1645694872")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}
	var dictResponse DictResponse
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(word, "UK:", dictResponse.Dictionary.Prons.En, "US:", dictResponse.Dictionary.Prons.EnUs)
	for _, item := range dictResponse.Dictionary.Explanations {
		fmt.Println(item)
	}
}

func queryWithBaidu(word string) {
	client := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf(`from=en&to=zh&query=%s&transtype=realtime&simple_means_flag=3&sign=54706.276099&token=db0e12e7028abbe6b85fa5d468713b54&domain=common`, word))
	req, err := http.NewRequest("POST", "https://fanyi.baidu.com/v2transapi?from=en&to=zh", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="99", "Google Chrome";v="99"`)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("Origin", "https://fanyi.baidu.com")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://fanyi.baidu.com/translate?aldtype=16047&query=&keyfrom=baidu&smartresult=dict&lang=auto2zh")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", "BIDUPSID=3F29A27AD86E55057CD97637A802E1E6; PSTM=1647253296; BAIDUID=3F29A27AD86E5505706D23E2DA4011D8:FG=1; REALTIME_TRANS_SWITCH=1; FANYI_WORD_SWITCH=1; HISTORY_SWITCH=1; SOUND_SPD_SWITCH=1; SOUND_PREFER_SWITCH=1; APPGUIDE_10_0_2=1; BDSFRCVID=J-FOJeC624dLGOnD_vjJupsQxxvg0f5TH6aozZ6YD1l0A_6Yu0TYEG0P-M8g0Ku-KA06ogKK0eOTHkCF_2uxOjjg8UtVJeC6EG0Ptf8g0f5; H_BDCLCKID_SF=tR-tVCtatCI3HnRv5t8_5-LH-UoX-I62aKDsLRI2BhcqEIL4hjjoej5yQ-PfK-5t-T7IXUcNB-TpSMbSj4Qo24POhxQn-hjIWnILbhRT5p5nhMJN3j7JDMP0-xPfa5Oy523ion5vQpnOEpQ3DRoWXPIqbN7P-p5Z5mAqKl0MLPbtbb0xXj_0Djb-Datttjna--oa3RTeb6rjDnCr-UDVXUI82h5y05OO3JrNKpA55bQhq-Oh2-cvynKZDnORXx74B5vvbPOMthRnOlRKbpJ8DUL1Db3J2-ox5TTdsR7yfp5oepvoD-oc3MkfLPjdJJQOBKQB0KnGbUQkeq8CQft20b0EeMtjW6LEtRk8oI0aJDvDqTrP-trf5DCShUFs0fCJB2Q-XPoO3KJWsCo-QMPb3UD0KhbIhPriW5cpoMbgylRp8P3y0bb2DUA1y4vpKhbBt2TxoUJ2abjne-53qtnWeMLebPRiJPQ9QgbW5hQ7tt5W8ncFbT7l5hKpbt-q0x-jLTnhVn0MBCK0hI_GjTL2j6QMMhKX2tRfKKOb04_8Kb7VbnAwQxnkbfJBDxcUX6bfJ2neahONb4J5VtP6LT_Vytt7yajK2MvbLarnKloY5Un_VqOw0bJpQT8rKn_OK5OibCrQMKTzab3vOIJNXpO1MUtzBN5thURB2DkO-4bCWJ5TMl5jDh3Mb6ksDMDtqj-etJCe_K-Qb-3bK4TYhR7E-tCsqxby26nZHmc9aJ5nJD_MehRjXPTUBnKqylojbhoOMTcMonLaQpP-HJ7uW6jZQ5_jD-QdtqcttNnkKl0MLT6Ybb0xyn_VyUoQjxnMBMPj5mOnanvn3fAKftnOM46JehL3346-35543bRTLnLy5KJtMDcnK4-XjT3QDM5; BDSFRCVID_BFESS=J-FOJeC624dLGOnD_vjJupsQxxvg0f5TH6aozZ6YD1l0A_6Yu0TYEG0P-M8g0Ku-KA06ogKK0eOTHkCF_2uxOjjg8UtVJeC6EG0Ptf8g0f5; H_BDCLCKID_SF_BFESS=tR-tVCtatCI3HnRv5t8_5-LH-UoX-I62aKDsLRI2BhcqEIL4hjjoej5yQ-PfK-5t-T7IXUcNB-TpSMbSj4Qo24POhxQn-hjIWnILbhRT5p5nhMJN3j7JDMP0-xPfa5Oy523ion5vQpnOEpQ3DRoWXPIqbN7P-p5Z5mAqKl0MLPbtbb0xXj_0Djb-Datttjna--oa3RTeb6rjDnCr-UDVXUI82h5y05OO3JrNKpA55bQhq-Oh2-cvynKZDnORXx74B5vvbPOMthRnOlRKbpJ8DUL1Db3J2-ox5TTdsR7yfp5oepvoD-oc3MkfLPjdJJQOBKQB0KnGbUQkeq8CQft20b0EeMtjW6LEtRk8oI0aJDvDqTrP-trf5DCShUFs0fCJB2Q-XPoO3KJWsCo-QMPb3UD0KhbIhPriW5cpoMbgylRp8P3y0bb2DUA1y4vpKhbBt2TxoUJ2abjne-53qtnWeMLebPRiJPQ9QgbW5hQ7tt5W8ncFbT7l5hKpbt-q0x-jLTnhVn0MBCK0hI_GjTL2j6QMMhKX2tRfKKOb04_8Kb7VbnAwQxnkbfJBDxcUX6bfJ2neahONb4J5VtP6LT_Vytt7yajK2MvbLarnKloY5Un_VqOw0bJpQT8rKn_OK5OibCrQMKTzab3vOIJNXpO1MUtzBN5thURB2DkO-4bCWJ5TMl5jDh3Mb6ksDMDtqj-etJCe_K-Qb-3bK4TYhR7E-tCsqxby26nZHmc9aJ5nJD_MehRjXPTUBnKqylojbhoOMTcMonLaQpP-HJ7uW6jZQ5_jD-QdtqcttNnkKl0MLT6Ybb0xyn_VyUoQjxnMBMPj5mOnanvn3fAKftnOM46JehL3346-35543bRTLnLy5KJtMDcnK4-XjT3QDM5; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; JAPAN_PINYIN_SWITCH=1; delPer=0; PSINO=5; H_PS_PSSID=36309_31254_36004_35910_36167_34584_35979_36074_36235_26350_36303_36312_36061; BA_HECTOR=20al85ak21ak0584ab1h7cb820q; Hm_lvt_64ecd82404c51e03dc91cb9e8c025574=1651912098; Hm_lpvt_64ecd82404c51e03dc91cb9e8c025574=1651912098; ab_sr=1.0.1_ZjRkMTE2ZDEyYjYzYmExODk0ZDAwYmY3NTBjOWQ4MjdmYmFlODM1NjcxNjA3YTYxYTA2N2ExM2U0YzcxZTdmNTM1ZmYzNGUyNDE2MGRlZmM3MmFmNWYxYmI4NmE4ODMxMDc0OGUwZmY2OGQzYzI2ODc4MmFmOGNiMzJjYTA1N2VjMTlmZTAwZDhmOTY2MGFmNTEzMTk3MGExOWZiMDRkMQ==")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("%s\n", bodyText)
	if resp.StatusCode != 200 {
		log.Fatal("bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}
	var dictResponseWithBaidu DictResponseWithBaidu
	err = json.Unmarshal(bodyText, &dictResponseWithBaidu)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range dictResponseWithBaidu.DictResult.Edict.Item {
		for _, number := range item.TrGroup {
			fmt.Println("by baidu", number.Example)
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, `usage: simpleDict WORD
example: simpleDict hello
		`)
		os.Exit(1)
	}
	word := os.Args[1]
	query(word)
	queryWithBaidu(word)
}
