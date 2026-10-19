package engine
import("bytes";"os";"strings";"path/filepath";"fmt";"testing";"time")
func TestResearchDecode(t *testing.T){
 arcpaths:=strings.Split(os.Getenv("RESEARCH_ARCHIVES"),":")
 files:=strings.Split(os.Getenv("RESEARCH_INPUTS"),":")
 if len(arcpaths)!=len(files) || len(files)<2 {t.Skip("no paths")}
 for i,p:=range arcpaths {
  archive,err:=os.ReadFile(p);if err!=nil {t.Fatal(err)}
  raw,err:=os.ReadFile(files[i]);if err!=nil{t.Fatal(err)}
  start:=time.Now();got,err:=decompress(archive,defaultLimit);duration:=time.Since(start).Seconds()
  if err!=nil || !bytes.Equal(got,raw) {t.Fatalf("decode %s: %v",p,err)}
  fmt.Printf("JLZDECODE,%s,%.9f\n",filepath.Base(files[i]),duration)
 }
}
