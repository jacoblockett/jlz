package main
import ("fmt"; "os";"strings";"testing"; "path/filepath")
func TestGateProbe(t *testing.T){p:=os.Getenv("SILESIA");if p=="" {t.Skip()};for _,name:=range []string{"dickens","mozilla","mr","nci","ooffice","osdb","reymont","samba","sao","webster","x-ray","xml"} {data,e:=os.ReadFile(filepath.Join(p,name));if e!=nil{t.Fatal(e)};on,off:=0,0;for pos:=0;pos<len(data);pos+=stage3BlockSize {end:=min(pos+stage3BlockSize,len(data));if likelyTransformedMatches(data[pos:end]){on++}else{off++}};fmt.Println("GATE",name,on,off,strings.Repeat(".",0))}}
