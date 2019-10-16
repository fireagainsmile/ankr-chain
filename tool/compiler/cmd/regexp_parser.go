package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/scanner"
)

type RegexpParser struct {
	scan           scanner.Scanner
	matcherMap map[int]Matcher
	st *ContractState
	errMsg []error
}

func (r *RegexpParser)doMatch(ruleId int, str string) {
	m := r.matcherMap[ruleId]
	if m == nil {
		fmt.Println("matcher not found!")
		return
	}
	err := m.Match(r.st, str)
	if err != nil {
		r.errMsg = append(r.errMsg, err)
	}
}

func NewRegexpParser() *RegexpParser {
	var p = &RegexpParser{
		matcherMap: make(map[int]Matcher),
		st:         NewContractClass(),
	}
	// add rules
	p.AddMatcher(INIT, NewInitMatcher(), true)
	p.AddMatcher(DERIVED, NewDerivedMatcher(), true)
	p.AddMatcher(VIRTUALFUNC, NewVirtualMatcher(), false)
	//p.AddMatcher(CONTRACTENTRY, NewEntryMatcher(), true)
	p.AddMatcher(ACTIONENTRY, NewActionEntryMatcher(), false)
	p.AddMatcher(TYPECHECK, NewTypeMatcher(), true)
	return p
}

func (r *RegexpParser)AddMatcher(id int, m Matcher, mustCheck bool) *RegexpParser {
	r.matcherMap[id] = m
	if mustCheck{
		r.st.requiredRule[id] = mustCheck
	}
	return r
}

func (r *RegexpParser) ParseFile(file string) {
	fileBuffer, err := os.Open(file)
	if err != nil {
		fmt.Println("file not found!")
		return
	}
	r.scan.Init(fileBuffer)
	defer fileBuffer.Close()
	allToken := make([]string, 0)
	for tok := r.scan.Scan(); tok != scanner.EOF; tok = r.scan.Scan(){
		allToken = append(allToken, r.scan.TokenText())
	}
	tokenString := strings.Join(allToken, " ")
	for k, _ := range r.matcherMap {
		r.doMatch(k, tokenString)
	}
	//r.scan.Init()
	//for tok := r.scan.Scan(); tok != scanner.EOF; tok = r.scan.Scan() {
	//	switch r.scan.TokenText() {
	//	case tokens[CLASS]:
	//		r.parseClassDeclare()
	//	case tokens[VOID], tokens[INT],tokens[FLOAT],tokens[CHAR]:
	//		r.parseFuncDefinition()
	//	case tokens[EXTERN]:
	//		r.parseExternFunc()
	//	}
	//}
}

func (r *RegexpParser)parseExternFunc()  {
	args := make([]string, 0)
	var sc Scope
	for !isOutScope(sc) {
		switch r.scan.TokenText() {
		case tokens[LBRACE]:
			sc.entered = true
			sc.subScope++
		case tokens[RBRACE]:
			sc.subScope--
		default:
			args = append(args, r.scan.TokenText())
		}
		r.scan.Scan()
	}
	argsString := strings.Join(args, " ")
	r.doMatch(CONTRACTENTRY, argsString)
}

//make sure how to define virtual and destruct function
func (r *RegexpParser) parseClassDeclare() {
	declareArgs := make([]string, 0)
	for r.scan.TokenText() != "{" {
		declareArgs = append(declareArgs, r.scan.TokenText())
		r.scan.Scan()
	}
	declaration := strings.Join(declareArgs, " ")
	r.doMatch(DERIVED, declaration)
	r.parseInsideClass()
}

// parse internal function definitions
// void init() and ~Class
func (r *RegexpParser) parseInsideClass() {
	var sc Scope
	for !isOutScope(sc) {
		switch r.scan.TokenText() {
		case tokens[LBRACE]:
			sc.entered = true
			sc.subScope++
		case tokens[RBRACE]:
			sc.subScope--
			//check void init() definition
		case tokens[VOID], tokens[INT], tokens[FLOAT], tokens[CHAR]:
			funcMember := r.tryGetFunctions()
			r.doMatch(INIT, funcMember)
			r.doMatch(ACTIONENTRY, funcMember)
		case "~":
			r.scan.Scan()
			r.doMatch(DESTRUCT, r.scan.TokenText())
			//check destruct function
		case tokens[VIRTUAL]:
			r.doMatch(VIRTUALFUNC, r.scan.TokenText())
		default:
		}
		r.scan.Scan()
	}
}


func (r *RegexpParser)tryGetFunctions() string {
	argList := make([]string,0)
	for  tok := r.scan.TokenText(); tok != ";" && tok != "{" ;  {
		argList = append(argList, tok)
		r.scan.Scan()
		tok = r.scan.TokenText()
	}

	if !isFunction(argList) {
		return ""
	}
	//fmt.Println(argList)
	return strings.Join(argList, " ")
}

func (r *RegexpParser)parseFuncDefinition()  {
	functionStr := r.tryGetFunctions()
	if functionStr == ""{
		return
	}
	r.doMatch(INIT, functionStr)
	r.doMatch(ACTIONENTRY, functionStr)
}

func (r *RegexpParser)ValidContract() bool {
	if len(r.errMsg) != 0 {
		for _, err := range r.errMsg {
			fmt.Println(err)
			return false
		}
	}
	for k, v := range r.st.requiredRule {
		if !v {
			fmt.Println("lack of definition of rule:", k)
			return false
		}
	}
	return true
}

type Scope struct {
	entered bool //mark if we entered a class scope
	subScope int //sub scope counter
}

func isOutScope(s Scope) bool {
	if s.entered && s.subScope == 0 {
		return true
	}
	return false
}

func isFunction(arg []string) bool {
	for _, tok := range arg {
		if tok == "(" {
			return  true
		}
	}
	return false
}