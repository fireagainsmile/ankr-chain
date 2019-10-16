package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	abi2 "github.com/Ankr-network/ankr-chain/tool/compiler/abi"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/scanner"
)

var(
	//expression to match invoke_action and invoke_func labeled functions
	invokeRegexp = `(INVOKE_ACTION|INVOKE_FUNC) \( "\b[\w]+`
)

type InvokeType struct {
	name string
	invokeType []string
}

func GenAbi(file string) error {
	cc := abi2.NewContractClass()
	err := parseClassFromFile(file, cc)
	if err != nil {
		return err
	}
	contractContent := readContract(file)
	functions := getActionAndFunc(contractContent)
	m := getActionEntry(functions, cc)
	jsonByte, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		return err
	}
	abiFile := getAbiFileName(file)
	err = writeABI(jsonByte, abiFile)
	if err != nil {
		return err
	}
	// if no action entry is found, auto generate extern code
	if len(functions) == 0 {
		err = cc.GenCode(file)
	}
	return err
}

func getAbiFileName(srcFile string) string {
	// replace cpp or cc with json
	abiFile := strings.TrimRight(srcFile, "cpp")
	abiFile = strings.TrimRight(srcFile, "cc")
	return  fmt.Sprintf("%sabi",abiFile)
}

func writeABI(abi []byte, fileName string) error {
	err := ioutil.WriteFile(fileName, abi, 0600)
	return err
}

func parseClassFromFile(file string, cc *abi2.ContractClass) error {
	cl := searchClass(file)
	if len(cl) == 0 {
		return errors.New("no class found! ")
	}
	//fmt.Println(cl)
	err := cc.ParseClass(cl)
	return err
}

//read contract
func readContract(file string) (contract []string) {
	var sc scanner.Scanner
	fileBuffer, err := os.Open(file)
	if err != nil {
		return
	}
	defer fileBuffer.Close()
	sc.Init(fileBuffer)
	for tok := sc.Scan(); tok != scanner.EOF; tok = sc.Scan() {
		contract = append(contract, sc.TokenText())
	}
	return
}

// collect functions called in action entry, prepare for abi output
// require actionEntry defined by user
func getActionEntry(funcs []InvokeType, cc *abi2.ContractClass) []*abi2.Method {
	m := make([]*abi2.Method, 0)
	// if no action entry is defined in contract, output abi
	if len(funcs) == 0 {
		for _, v := range cc.FuncCache {
			m = append(m, cc.FuncSigs[v])
		}
		return m
	}

	for _, v := range funcs {
		funcSig := cc.FuncSigs[v.name]
		if funcSig != nil {
			if len(v.invokeType) != 0 {
				funcSig.Type = v.invokeType
			}
			m = append(m, funcSig)
		}
	}
	return m
}

//input shall be the content of a contract
func getActionAndFunc(contract []string) []InvokeType {
	actionAndEvent := make([]InvokeType, 0)
	reg, _ := regexp.Compile(invokeRegexp)
	//fmt.Println(contract)
	invokes := reg.FindAllString(strings.Join(contract, " "), -1)
	for _, invoke := range invokes {
		var invokeFunc InvokeType
		invoke = strings.TrimRight(invoke, " ")
		invokeSlice := strings.Split(invoke, " ")
		length := len(invokeSlice)
		 invokeFunc.name = strings.TrimLeft(invokeSlice[length-1], "\"")
		switch invokeSlice[0] {
		case "INVOKE_ACTION":
			invokeFunc.invokeType = []string{"action","event"}
		case "INVOKE_FUNC":
			invokeFunc.invokeType = []string{"action","event"}
			
		}
		actionAndEvent = append(actionAndEvent, invokeFunc)
	}
	return actionAndEvent
}

func searchClass(file string) (class []string) {
	var sc scanner.Scanner
	fileBuffer, err := os.Open(file)
	defer fileBuffer.Close()
	if err != nil {
		return
	}
	sc.Init(fileBuffer)
	for tok := sc.Scan(); tok != scanner.EOF; tok = sc.Scan(){
		switch sc.TokenText() {
		case tokens[CLASS]:
			classDeclare := readClass(sc)
			reg, _ := regexp.Compile(baseRegexp)
			if reg.MatchString(strings.Join(classDeclare, " ")) {
				return classDeclare
			}
		}
	}
	return
}

func readClass(sc scanner.Scanner) []string {
	var scope Scope
	args := make([]string,0)
	for !isOutScope(scope) {
		switch sc.TokenText() {
		case tokens[LBRACE]:
			scope.entered = true
			scope.subScope++
		case tokens[RBRACE]:
			scope.subScope--
		}
		args = append(args, sc.TokenText())
		sc.Scan()
	}
	return args
}
