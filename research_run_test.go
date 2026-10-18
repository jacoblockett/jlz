package main
import (
 "bytes"
 "fmt"
 "os"
 "testing"
 "time"
 "strings"
 "path/filepath"
)
func TestResearchRun(t *testing.T) {
 paths:=strings.Split(os.Getenv("RESEARCH_INPUTS"),":")
 for _,p:=range paths {if p=="" {t.Skip("no inputs")}}
 type inp struct {name string; data []byte}
 var inputs []inp
 for _,p:=range paths {b,e:=os.ReadFile(p);if e!=nil {t.Fatal(e)}; inputs=append(inputs,inp{filepath.Base(p),b})}
 for _,f:=range inputs {
  start:=time.Now();out,s,e:=compress(f.data); dur:=time.Since(start).Seconds();if e!=nil {t.Fatal(e)}
  if os.Getenv("RESEARCH_ROUNDTRIP")=="1" {decoded,e:=decompress(out,defaultLimit);if e!=nil || !bytes.Equal(decoded,f.data) {t.Fatalf("roundtrip %s: %v",f.name,e)}}
  fmt.Printf("JLZROW,%s,%d,%d,%.9f,%s\n",f.name,len(f.data),len(out),dur,s.method)
 }
}
