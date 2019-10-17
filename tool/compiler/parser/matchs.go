package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// regexp to match target function definition
	baseRegexp   = `class [\w]+( : public akchain : : Contract)`
	derivedRegexp = ` : (public|private|protected) [\w]+`
	initFuncRegexp     = `char \* init \( \)`
	entryfuncRegexp = `char \* ContractEntry \( char \* [\w]+ , char \* [\w]+ \)`
	actionEntryRegexp = `char \* actionEntry \( const char \* [\w]+ , const char \* [\w]+ \)`
	virtualRexp = `virtual (void|int|char *) [\w]+`
)

const (
	BASE = iota
	DERIVED
	INIT
	DESTRUCT
	VIRTUALFUNC
	ACTIONENTRY
	CONTRACTENTRY
	TYPECHECK
)

type Matcher interface {
	Match(c *ContractState, str string) error
}

type initMather struct {
	regPatten1 string
	regPatten2 string
}

func NewInitMatcher() *initMather {
	return &initMather{
		regPatten1:initFuncRegexp,
	}
}

func (i *initMather)Match(c *ContractState, args string) error  {
	reg, _ := regexp.Compile(i.regPatten1)
	if reg.MatchString(args) {
		c.requiredRule[INIT] = true
		return nil
	}
	return nil
}

type derivedMatcher struct {
	basePatten string
	derivedPatten string
}

func NewDerivedMatcher() *derivedMatcher {
	return &derivedMatcher{
		basePatten:baseRegexp,
		derivedPatten:derivedRegexp,
	}
}

func (d *derivedMatcher)Match(c *ContractState, str string) error  {
	reg, _ := regexp.Compile(d.derivedPatten)
	matchedStrs := reg.FindAllString(str, -1)
	if len(matchedStrs) != 1 {
		return errors.New("one derived class from ankr base should be defined! ")
	}

	reg, _ = regexp.Compile(d.basePatten)
	matchedStrs = reg.FindAllString(str, -1)
	if len(matchedStrs) == 0 {
		return errors.New("user defined class should derived from ankr base. ")
	}
	c.requiredRule[DERIVED] = true
	strSlice := strings.Split(matchedStrs[0], " ")
	c.className = strSlice[1]
	return nil
}

type destructMatcher struct {
	//reg string
}

func NewDestructMatcher() *destructMatcher {
	return &destructMatcher{}
}

func (d *destructMatcher)Match(c *ContractState, str string) error  {
	if c.className == str {
		c.requiredRule[DESTRUCT] = true
	}
	return nil
}
type virtualMatcher struct {
	reg string
}

func NewVirtualMatcher() *virtualMatcher {
	return &virtualMatcher{
		virtualRexp,
	}
}
func (v *virtualMatcher)Match(c *ContractState, str string) error  {
	reg, _ := regexp.Compile(v.reg)
	if reg.MatchString(str) {
		return errors.New("virtual member function found")
	}
	return nil
}

type entryMatcher struct {
	reg string
}

func NewEntryMatcher() *entryMatcher  {
	return &entryMatcher{
		reg:entryfuncRegexp,
	}
}

func (e *entryMatcher)Match(c *ContractState, str string) error  {
	reg, _ := regexp.Compile(e.reg)
	if reg.MatchString(str) {
		c.requiredRule[CONTRACTENTRY] = true
	}
	return nil
}

type actionEntryMatcher struct {
	reg string
}

func NewActionEntryMatcher() *actionEntryMatcher {
	return &actionEntryMatcher{
		reg:actionEntryRegexp,
	}
}

func (a *actionEntryMatcher)Match(c *ContractState, str string) error  {
	reg, _ := regexp.Compile(a.reg)
	if reg.MatchString(str) {
		c.requiredRule[ACTIONENTRY] = true
		entry := NewEntryMatcher()
		return entry.Match(c, str)
	}
	return nil
}

type typeMatcher struct {
	limitedTypes []string
}

func NewTypeMatcher() *typeMatcher {
	return &typeMatcher{
		limitedTypes:[]string{"float", "complex"},
	}
}

func (t *typeMatcher)Match(c *ContractState, str string) error {
	for _, exp := range t.limitedTypes{
		reg, err := regexp.Compile(exp)
		if err != nil  {
			return err
		}
		if reg.MatchString(str){
			return errors.New(fmt.Sprintf("%s is not allowed type!", exp) )
		}
	}
	c.requiredRule[TYPECHECK] = true
	return nil
}