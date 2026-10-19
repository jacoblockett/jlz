package jlz

import (
 "bytes"
 "crypto/sha256"
 "fmt"
 "os"
 "path/filepath"
 "sort"
 "strconv"
 "strings"
 "testing"
 "time"
)
func TestResearchCorpus(t *testing.T) {
 root:=os.Getenv("JLZ_CORPUS"); if root=="" {t.Skip("research-only")}
 runs,_:=strconv.Atoi(os.Getenv("JLZ_RUNS")); if runs<1 {runs=1}
 names:=[]string{"dickens","mozilla","mr","nci","ooffice","osdb","reymont","samba","sao","webster","xml","x-ray"}
 if sub:=os.Getenv("JLZ_FILES");sub!="" {names=strings.Split(sub,",")}
 type f struct{name string;data []byte; hash [32]byte}; var all []f
 for _,n:=range names {b,e:=os.ReadFile(filepath.Join(root,n));if e!=nil {t.Fatal(e)};all=append(all,f{n,b,sha256.Sum256(b)})}
 times:=make([]float64,0,runs); sizes:=make([]int,0,runs)
 for i:=0;i<runs;i++ {totalTime:=time.Duration(0); total:=0
  for _,f:=range all {
   st:=time.Now(); arc,_,err:=compressFast(f.data); totalTime+=time.Since(st); if err!=nil {t.Fatal(err)}
   total+=len(arc)
   if i==0 {orig,e:=decompress(arc,1<<30);if e!=nil||!bytes.Equal(orig,f.data) {t.Fatalf("failed roundtrip %s: %v",f.name,e)};fmt.Printf("FILE %s %d %d\n",f.name,len(f.data),len(arc))}
  }
  sec:=totalTime.Seconds();times=append(times,sec);sizes=append(sizes,total);fmt.Printf("RUN %d %.6f %d\n",i+1,sec,total)
 }
 sort.Float64s(times);fmt.Printf("MEDIAN %.6f TOTAL %d\n", (times[(runs-1)/2]+times[runs/2])/2, sizes[0])
}
