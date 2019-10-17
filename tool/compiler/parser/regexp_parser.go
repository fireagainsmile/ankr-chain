package parser

import (
	"errors"
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

func (p *RegexpParser)Execute(args []string) error {
	sourceFile := args[0]
	p.ParseFile(sourceFile)
	return p.ValidContract()
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
}

func (r *RegexpParser)ValidContract() error {
	if len(r.errMsg) != 0 {
		return r.errMsg[0]
	}
	for k, v := range r.st.requiredRule {
		if !v {
			return errors.New(fmt.Sprintf("Paser encountered an error. Error code:%d", k))
		}
	}
	return nil
}