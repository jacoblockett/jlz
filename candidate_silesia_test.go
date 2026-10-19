package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"
)

// Candidate-only Silesia medians. Entire originals reside in RAM before timing.
func TestCandidateSilesiaMedian(t *testing.T) {
	root := os.Getenv("JLZ_CANDIDATE_CORPUS")
	if root == "" { t.Skip("candidate Silesia bench only") }
	names := []string{"dickens","mozilla","mr","nci","ooffice","osdb","reymont","samba","sao","webster","xml","x-ray"}
	expected := []string{
		"b24c37886142e11d0ee687db6ab06f936207aa7f2ea1fd1d9a36763c7a507e6a",
		"657fc3764b0c75ac9de9623125705831ebbfbe08fed248df73bc2dc66e2a963b",
		"68637ed52e3e4860174ed2dc0840ac77d5f1a60abbcb13770d5754e3774d53e6",
		"fc63a31770947b8c2062d3b19ca94c00485a232bb91b502021948fee983e1635",
		"e7ee013880d34dd5208283d0d3d91b07f442e067454276095ded14f322a656eb",
		"60f027179302ca3ad87c58ac90b6be72ec23588aaa7b3fe8ecc0f11def3fa3",
		"0eac0114a3dfe6e2ee1f345a0f79d653cb26c3bc9f0ed79238af4933422b7578",
		"93ba07bc44d8267789c1d911992f40b089ffa2140b4a160fac11ccae9a40e7b2",
		"c2d0ea2cc59d4c21b7fe43a71499342a00cbe530a1d5548770e91ecd6214adcc",
		"6a68f69b26daf09f9dd84f7470368553194a0b294fcfa80f1604efb11143a383",
		"0e82e54e695c1938e4193448022543845b33020c8be6bf3bf3ead2224903e08c",
		"7de9fce1405dc44ae5e6813ed21cd5751e761bd4265655a005d39b9685d1c9ad",
	}
	type item struct { name string; data []byte }
	files := make([]item,0,12)
	for i,name := range names {
		data,err := os.ReadFile(filepath.Join(root,name))
		if err!=nil {t.Fatal(err)}
		if fmt.Sprintf("%x",sha256.Sum256(data))!=expected[i] {t.Fatalf("source hash mismatch %s",name)}
		files=append(files,item{name,data})
	}
	encTimes:=make([]float64,0,10)
	decTimes:=make([]float64,0,10)
	referenceSize:=0
	for run:=0;run<10;run++ {
		enct,decs,total:=time.Duration(0),time.Duration(0),0
		for _,f:=range files {
			start:=time.Now()
			arc,_,err:=compress(f.data)
			enct+=time.Since(start)
			if err!=nil {t.Fatal(err)}
			total+=len(arc)
			start=time.Now()
			out,err:=decompress(arc,defaultLimit)
			decs+=time.Since(start)
			if err!=nil || !bytes.Equal(out,f.data) {t.Fatalf("roundtrip %s %v",f.name,err)}
			if run==0 {fmt.Printf("CANDIDATE FILE %s %d\n",f.name,len(arc))}
		}
		if run==0 {referenceSize=total}
		if total!=referenceSize {t.Fatalf("unstable size: %d vs %d",total,referenceSize)}
		encTimes=append(encTimes,enct.Seconds())
		decTimes=append(decTimes,decs.Seconds())
		fmt.Printf("CANDIDATE RUN %d SIZE %d ENCODE %.6f DECODE %.6f\n",run+1,total,enct.Seconds(),decs.Seconds())
	}
	sort.Float64s(encTimes);sort.Float64s(decTimes)
	fmt.Printf("CANDIDATE MEDIAN SIZE %d ENCODE %.6f DECODE %.6f\n",referenceSize,(encTimes[4]+encTimes[5])/2,(decTimes[4]+decTimes[5])/2)
}
