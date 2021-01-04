package main

import (
  "github.com/valyala/fastjson"
  "io/ioutil"
  "net/http"
  "strings"
  "regexp"
  "bufio"
  "log"
  "fmt"
  "os"
)

var pathFile = os.Args[1] //'/var/WWW/approot"
var DisSep = "/"
var pathes []string

type (
  ApiCred struct {
    Login string
    Password string
  }
)

func main() {
  
  if len(os.Args) != 2 {
    fmt.Println("Usage: approot path")
    return
  }
  
  // ## Find pathes (basic or advanced apps)
  pathes := []string{
    "controllers",
    "modules/admin/controllers",
  }

  // ## Collecting name of the controllers for current path
  // var controllers []string

  // ## Collecting names of actions of current controller file
  // var actions []string
  // var parameters []string // parameters for current controller

  lines, err := GetFileLines(pathFile + DisSep + pathes[0] + DisSep + "ZzsoController.php", func(s string)(string,bool){ 
    return s, true
  })

  errCheck(err, "parse file")

  for _, l := range lines {
    if isFunctionLine(l) {
      fmt.Println(l)
    }
  }

  // ## Construction the test-urls from modules/controller/action names

  // ## Make request foreach collected test-url and save responce code

  // ## Save result of bad statuses and bad results occured

}

func GetFileLines(filePath string, parse func(string) (string,bool)) ([]string, error) {
  inputFile, err := os.Open(filePath)
  if err != nil {
    return nil, err
  }
  defer inputFile.Close()

  scanner := bufio.NewScanner(inputFile)
  var results []string
  for scanner.Scan() {
    if output, add := parse(scanner.Text()); add {
      results = append(results, output)
    }
  }
  if err := scanner.Err(); err != nil {
    return nil, err
  }
  return results, nil
}

func errCheck(err error, msg string) {
  if err != nil {
    fmt.Println("Error while", msg, err)
    return
  }
}

func selectControllersItemName(str string) string {
    res := strings.Split(str, "Controller")
    return strings.ToLower(res[0])
}

func isFunctionLine(str string) bool {
    re:= regexp.MustCompile("(function)(.*)(action([a-zA-Z]+))[(]{1}(.*)[)]{1}")
    return re.MatchString(str)
}

func selectActionName(str string) string {
    re:= regexp.MustCompile("action([a-zA-Z]+)([/b]*)")
    res := re.FindStringSubmatch(str)
    if len(res) > 1 {
        return res[1]
    }
    return ""
}

func selectActionParameters(str string) string {
    re:= regexp.MustCompile("([(]{1})(.*)([)]{1})")
    res := re.FindStringSubmatch(str)
    if len(res) > 2 {
        return strings.TrimSpace(res[2])
    }
    return ""
}

func collectParameters(str string) []string {
    res := strings.Split(str, ",")
    for i, value := range res {
        splitted := strings.Split(value, "=")
        res[i] = strings.Trim(splitted[0], "$ ")
    }
    return res
}

func makeGetRequest(cr ApiCred, url string) (string, string, string) {
  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    log.Fatal(err)
  }
  req.SetBasicAuth(cr.Login, cr.Password)
  cli := &http.Client{}
  resp, err := cli.Do(req)

  bodyData, err := ioutil.ReadAll(resp.Body)
  msg := fastjson.GetString(bodyData, "tasks", "0", "result", "0", "login")
  statusCode := fastjson.GetString(bodyData, "status_code")
  errMsg := fastjson.GetString(bodyData, "tasks_error")
  return string(msg), statusCode, errMsg
}