package env

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func init() {
	Log = log.New(os.Stderr, "[env] ", log.LstdFlags)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var (
	Log *log.Logger
)

type EnvSet struct {
	name   string
	envs   []string
	envMap map[string]*Env
	mu     *sync.Mutex
	Usage  func()
}

func (e *EnvSet) defaultUsage() {
	if e.name == "" {
		fmt.Fprintf(os.Stdout, "Usage:\n")
	} else {
		fmt.Fprintf(os.Stdout, "Usage of %s:\n", e.name)
	}
	e.PrintDefaults()
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		Environment.Usage()
	}
}

func NewEnvSet() *EnvSet {
	e := &EnvSet{
		// 		envs:   []*Env{},
		envMap: make(map[string]*Env),
		mu:     &sync.Mutex{},
	}
	e.Usage = func() {
		fmt.Println("Support ENV:")
		e.PrintDefaults()
	}
	return e
}

var Environment = NewEnvSet()

type Env struct {
	Name  string
	Usage string
	// 	Value  Value
	DefVal interface{}
}

// type Value interface {
// 	String() string
// 	Set(string) error
// }

func (e *Env) String() string {
	v := os.Getenv(e.Name)
	if len(v) > 0 {
		return v
	}
	if v, ok := e.DefVal.(string); ok {
		return v
	}
	return ""
}
func StringVar(p *string, name string, value string, usage string) {
	Environment.StringVar(p, name, value, usage)
}
func (es *EnvSet) StringVar(p *string, name string, value string, usage string) {
	e := es.New(name, value, usage)
	if p != nil {
		*p = e.String()
	}
	es.Add(e)
}

func (e *Env) Bool() (bool, error) {
	v := strings.ToLower(os.Getenv(e.Name))
	for _, s := range []string{"true", "1"} {
		if v == s {
			return true, nil
		}
	}
	for _, s := range []string{"false", "0"} {
		if v == s {
			return false, nil
		}
	}
	if d, ok := e.DefVal.(bool); ok {
		return d, nil
	}
	return false, fmt.Errorf("%s not bool (%s)", e.Name, os.Getenv(e.Name))
}
func (e *Env) MustBool() bool {
	x, err := e.Bool()
	if err != nil {
		Log.Println("env", e.Name, "must be bool", err.Error())
		if v, ok := e.DefVal.(bool); ok {
			x = v
		}
	}
	return x
}
func BoolVar(p *bool, name string, value bool, usage string) {
	Environment.BoolVar(p, name, value, usage)
}
func (es *EnvSet) BoolVar(p *bool, name string, value bool, usage string) {
	e := es.New(name, value, usage)
	if p != nil {
		*p = e.MustBool()
	}
	es.Add(e)
}

func (e *Env) Int() (int, error) {
	s := os.Getenv(e.Name)
	if len(s) > 0 {
		return strconv.Atoi(s)
	}
	if v, ok := e.DefVal.(int); ok {
		return v, nil
	}
	return 0, fmt.Errorf("DEFAULT VALUE FOR [%s] IS NOT INT (%T: %v)", e.Name, e.DefVal, e.DefVal)
}
func (e *Env) MustInt() int {
	x, err := e.Int()
	if err != nil {
		Log.Println("env", e.Name, "must be int", err.Error())
		if v, ok := e.DefVal.(int); ok {
			x = v
		}
	}
	return x
}
func IntVar(p *int, name string, value int, usage string) {
	Environment.IntVar(p, name, value, usage)
}
func (es *EnvSet) IntVar(p *int, name string, value int, usage string) {
	e := es.New(name, value, usage)
	if p != nil {
		*p = e.MustInt()
	}
	es.Add(e)
}

func (e *Env) Int64() (int64, error) {
	s := os.Getenv(e.Name)
	if len(s) > 0 {
		return strconv.ParseInt(s, 10, 64)
	}
	if v, ok := e.DefVal.(int64); ok {
		return v, nil
	}
	return 0, fmt.Errorf("DEFAULT VALUE FOR [%s] IS NOT INT64 (%T: %v)", e.Name, e.DefVal, e.DefVal)
}
func (e *Env) MustInt64() int64 {
	x, err := e.Int64()
	if err != nil {
		Log.Println("env", e.Name, "must be int64", err.Error())
		if v, ok := e.DefVal.(int64); ok {
			x = v
		}
	}
	return x
}
func Int64Var(p *int64, name string, value int64, usage string) {
	Environment.Int64Var(p, name, value, usage)
}
func (es *EnvSet) Int64Var(p *int64, name string, value int64, usage string) {
	e := es.New(name, value, usage)
	if p != nil {
		*p = e.MustInt64()
	}
	es.Add(e)
}

// --- Start Float64

func (e *Env) Float64() (float64, error) {
	s := os.Getenv(e.Name)
	if len(s) > 0 {
		return strconv.ParseFloat(s, 64)
	}
	if v, ok := e.DefVal.(float64); ok {
		return v, nil
	}
	return 0, fmt.Errorf("DEFAULT VALUE FOR [%s] IS NOT FLOAT64 (%T: %v)", e.Name, e.DefVal, e.DefVal)
}
func (e *Env) MustFloat64() float64 {
	x, err := e.Float64()
	if err != nil {
		Log.Println("env", e.Name, "must be float64", err.Error())
		if v, ok := e.DefVal.(float64); ok {
			x = v
		}
	}
	return x
}
func Float64Var(p *float64, name string, value float64, usage string) {
	Environment.Float64Var(p, name, value, usage)
}
func (es *EnvSet) Float64Var(p *float64, name string, value float64, usage string) {
	e := es.New(name, value, usage)
	if p != nil {
		*p = e.MustFloat64()
	}
	es.Add(e)
}

// --- Start Duration

func (e *Env) Duration() (time.Duration, error) {
	s := os.Getenv(e.Name)
	if len(s) > 0 {
		return time.ParseDuration(os.Getenv(e.Name))
	}
	if v, ok := e.DefVal.(time.Duration); ok {
		return v, nil
	}
	return 0, fmt.Errorf("DEFAULT VALUE FOR [%s] IS NOT DURATION (%T: %v)", e.Name, e.DefVal, e.DefVal)
}
func (e *Env) MustDuration() time.Duration {
	x, err := e.Duration()
	if err != nil {
		Log.Println("env", e.Name, "must be duration", err.Error())
		if v, ok := e.DefVal.(time.Duration); ok {
			x = v
		}
	}
	return x
}
func DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	Environment.DurationVar(p, name, value, usage)
}
func (es *EnvSet) DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	e := es.New(name, value, usage)
	if p != nil {
		*p = e.MustDuration()
	}
	es.Add(e)
}

// --- End Duration

func New(name string, defVal interface{}, usage ...string) *Env {
	return Environment.New(name, defVal, usage...)
}
func (es *EnvSet) New(name string, defVal interface{}, usage ...string) *Env {
	if len(usage) > 0 {
		return &Env{Name: name, DefVal: defVal, Usage: usage[0]}
	}
	return &Env{Name: name, DefVal: defVal}
}

func Add(e *Env) {
	Environment.Add(e)
}
func (es *EnvSet) Add(e *Env) {
	es.mu.Lock()
	defer es.mu.Unlock()
	if val, ok := es.envMap[e.Name]; !ok {
		es.envMap[e.Name] = e
		es.envs = append(es.envs, e.Name)
	} else {
		if val.String() != e.String() {
			Log.Printf("env defined already %s; exist: %s; overwrite: %s", e.Name, val.String(), e.String())
			es.envMap[e.Name] = e
		}
	}
}
func Get(name string) *Env {
	return Environment.Get(name)
}
func (es *EnvSet) Get(name string) *Env {
	es.mu.Lock()
	defer es.mu.Unlock()
	return es.envMap[name]
}

func PrintDefaults() {
	Environment.PrintDefaults()
}

func (es *EnvSet) PrintDefaults() {
	var s string
	for _, envName := range es.envs {
		v := es.envMap[envName]
		cur := os.Getenv(envName)
		if len(cur) == 0 {
			cur = fmt.Sprint(v.DefVal)
		}
		s += fmt.Sprintf("  - %s (current \"%s\")\n", v.Name, cur)
		if len(v.Usage) > 0 {
			s += fmt.Sprintf("    \t%s\n", v.Usage)
		}
	}
	fmt.Printf(s)
}
