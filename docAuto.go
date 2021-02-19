package docauto

import (
	"encoding/json"
	"fmt"
	"github.com/Titanarthas/docauto/models"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
)

const (
	tagName = "docComment"
	tagJson = "json"
)

func GenerateDocStruct(methodType, urlPath  string, interfaceCom string, req, rsp interface{}) {
	val := reflect.ValueOf(req)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	// we only accept structs
	if val.Kind() != reflect.Struct {
		// return
	}

	reqBody, _ := json.Marshal(req)
	rspBody, _ := json.Marshal(rsp)
	apiCall := models.ApiCall{}
	apiCall.MethodType = methodType
	apiCall.CurrentPath = urlPath + " " + interfaceCom
	apiCall.RequestBody = string(reqBody)
	apiCall.ResponseBody = string(rspBody)
	apiCall.ResponseCode = 200
	apiCall.RequestComment = make(map[string][]string)

	if val.Kind() == reflect.Struct {
		generateVar(val, &apiCall, "")
	}

	GenerateHtml(&apiCall)
}

func generateVar(val reflect.Value, apiCall *models.ApiCall, parName string) {
	if val.Kind() == reflect.Slice {
		if val.Len() > 0 {
			ti := val.Index(0)
			generateVar(ti, apiCall, parName)
		}
		return
	}

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		if typeField.PkgPath != "" {
			continue // Private field
		}

		t := val.Field(i)
		tt := t.Type().Kind()

		varName := typeField.Name
		jsonName := typeField.Tag.Get(tagJson)
		if len(jsonName) > 0 {
			varName = jsonName
		}

		if parName != "" {
			varName = parName + "/" + varName
		}

		if tt == reflect.Struct {
			generateVar(t, apiCall, varName)
		} else if tt == reflect.Slice {
			if t.Len() > 0 {
				ti := t.Index(0)
				if ti.Kind() == reflect.Struct {
					generateVar(ti, apiCall, varName)
				}

			}
		}

		tag := typeField.Tag.Get(tagName)
		varType := typeField.Type.String()

		apiCall.RequestComment[varName] = []string{varType, tag}
	}
}

func GenerateHtml(apiCall *models.ApiCall) {
	shouldAddPathSpec := true
	for k, apiSpec := range spec.ApiSpecs {
		if apiSpec.Path == apiCall.CurrentPath && apiSpec.HttpVerb == apiCall.MethodType {
			if !config.ReqRepeat {
				spec.ApiSpecs[k].Calls = nil
				apiSpec.Calls = nil
			}
			
			shouldAddPathSpec = false
			apiCall.Id = count
			count += 1
			deleteCommonHeaders(apiCall)
			avoid := false
			for _, currentApiCall := range spec.ApiSpecs[k].Calls {
				if apiCall.RequestBody == currentApiCall.RequestBody &&
					apiCall.ResponseCode == currentApiCall.ResponseCode &&
					apiCall.ResponseBody == currentApiCall.ResponseBody {
					avoid = true
				}
			}
			if !avoid {
				spec.ApiSpecs[k].Calls = append(apiSpec.Calls, *apiCall)
			}
		}
	}

	if shouldAddPathSpec {
		apiSpec := models.ApiSpec{
			HttpVerb: apiCall.MethodType,
			Path:     apiCall.CurrentPath,
		}
		apiCall.Id = count
		count += 1
		deleteCommonHeaders(apiCall)
		apiSpec.Calls = append(apiSpec.Calls, *apiCall)
		spec.ApiSpecs = append(spec.ApiSpecs, apiSpec)
	}

	sort.Sort(spec.ApiSpecs)

	filePath, err := filepath.Abs(config.DocPath)
	dataFile, err := os.Create(filePath + ".json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dataFile.Close()
	data, err := json.Marshal(spec)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = dataFile.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	generateHtml()
}

func generateHtml() {
	funcs := template.FuncMap{"add": add, "mult": mult}
	t := template.New("API Documentation").Funcs(funcs)
	htmlString := Template
	t, err := t.Parse(htmlString)
	if err != nil {
		fmt.Println(err)
		return
	}
	filePath, err := filepath.Abs(config.DocPath)
	if err != nil {
		panic("Error while creating file path : " + err.Error())
	}
	homeHtmlFile, err := os.Create(filePath)
	defer homeHtmlFile.Close()
	if err != nil {
		panic("Error while creating documentation file : " + err.Error())
	}
	homeWriter := io.Writer(homeHtmlFile)
	t.Execute(homeWriter, map[string]interface{}{"array": spec.ApiSpecs,
		"baseUrls": config.BaseUrls, "Title": config.DocTitle})
}

func add(x, y int) int {
	return x + y
}

func mult(x, y int) int {
	return (x + 1) * y
}
type Config struct {
	On bool
	ReqRepeat bool

	BaseUrls map[string]string

	DocTitle string
	DocPath  string
}


var config *Config
var count int
var spec *models.Spec = &models.Spec{}

func Init(conf *Config) {
	config = conf
	// load the config file
	if conf.DocPath == "" {
		conf.DocPath = "apidoc.html"
	}


	filePath, err := filepath.Abs(conf.DocPath + ".json")
	dataFile, err := os.Open(filePath)
	defer dataFile.Close()
	if err == nil {
		json.NewDecoder(io.Reader(dataFile)).Decode(spec)
		generateHtml()
	}
}

func On() bool {
	return config.On
}

func deleteCommonHeaders(call *models.ApiCall) {
	delete(call.RequestHeader, "Accept")
	delete(call.RequestHeader, "Accept-Encoding")
	delete(call.RequestHeader, "Accept-Language")
	delete(call.RequestHeader, "Cache-Control")
	delete(call.RequestHeader, "Connection")
	delete(call.RequestHeader, "Cookie")
	delete(call.RequestHeader, "Origin")
	delete(call.RequestHeader, "User-Agent")
}
