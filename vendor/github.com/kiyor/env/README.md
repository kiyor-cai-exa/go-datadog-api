Use with flag package

```
func init() {
	env.Add(env.New("HTTP_PROXY","proxy use for call"))
}
func main() {
	flag.Usage = func() {
	    fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	    flag.PrintDefaults()
	    fmt.Println("Support ENV:")                                                                                         
	    env.PrintDefaults()
	}
	flag.Parse()
}
```
