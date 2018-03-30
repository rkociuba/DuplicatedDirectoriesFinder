package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "path"
    "os"
    "sort"
)

type DataByLevelType []map[int64]map[*Dir]struct{}
var DataByLevel DataByLevelType
var Dirs uint64

type Dir struct {
	ParentDir *Dir
	SubDirs []*Dir
	Name string
	FileCount  int
	ReversedLevel int
	TotalSize int64
	HighestOriginal int 
}

func deleteSubDirs(dir *Dir) {
	for _, sd := range dir.SubDirs {
		deleteSubDirs(sd)			
	}
	dir.SubDirs = []*Dir{}
	DeleteFromDataByLevel(dir)
	dir = nil
}

func markSubDirsHighestOriginal(dir *Dir, level int) {
	dir.HighestOriginal = level
	for _, sd := range dir.SubDirs {
		markSubDirsHighestOriginal(sd, level)			
	}
}

func printSubDirs(dir *Dir) {
	fmt.Printf(" %+v\n", dir)
	for _, sd := range dir.SubDirs {
		printSubDirs(sd)					
	}
}

func absPath(dir *Dir) string {
	if dir.ParentDir == nil {
		return ""
	} else {
		return path.Join(absPath(dir.ParentDir), dir.Name)
	}
}

//results: trunk to duplicates(i.e. level to slice of duplicates)
var dups = map[*Dir][][]*Dir{}

type Pair struct {
	Key *Dir
	Value [][]*Dir
}
type PairList []Pair
func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Key.ReversedLevel < p[j].Key.ReversedLevel }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }

func SortByLevel(input map[*Dir][][]*Dir) PairList{
	pl := make(PairList, len(input))
	i := 0
	for k, v := range input {
		pl[i] = Pair{k, v}
		i++
  	}
  	sort.Sort(sort.Reverse(pl))
  	return pl
}


func findHighestRoot(d *Dir) *Dir{
	for {
		if d.ParentDir == nil || d.ParentDir.HighestOriginal == -1 {
			return d
		}
		d = d.ParentDir
	}
}

func findLevelDuplicate(dlev map[int64]map[*Dir]struct{}){
	for _, pdirs := range dlev {
		if len(pdirs) > 1 {
			// checik  if  one of (or more) dirs is subtree (duplicate tree) of bigger three *
			highest := -1
			var pHighest, d *Dir = nil, nil
			for d, _ = range pdirs {
				//find highest original
				if d.HighestOriginal > highest{
					highest = d.HighestOriginal
					pHighest = d
				}
			}

			var trunk *Dir 
			if pHighest == nil {
				//no trunk - this duplicate set is not subset of biggest duplicate set
				//select new trunk
				trunk = d //TODO: whichever for now
				dups[trunk] = [][]*Dir{}
				markSubDirsHighestOriginal(trunk, trunk.ReversedLevel)
				fmt.Printf("NEW TRUNK: %s (%+v)\n", absPath(trunk), d)
			} else {
				trunk = findHighestRoot(pHighest)		
				fmt.Printf("OLD TRUNK: %s\n", absPath(trunk))
			}

			trunkOneWildcard := true //to prevent marking and potentially remove all duplicates from trunk (remove totally all)
			group := []*Dir{}
			for d, _ := range pdirs {
				if d.HighestOriginal != -1 && trunkOneWildcard {
					trunkOneWildcard = false
					continue
				}
				fmt.Printf(" Duplicate: %d %s: (%+v)\n", d.ReversedLevel, absPath(d), d)	
				//
				
				group = append(group, d)
				//this is duplicate or sub-duplicate 
				deleteSubDirs(d)
			}
			dups[trunk] = append(dups[trunk], group)

			fmt.Println("")	
		} 	
	}
} 

func InsertToDataByLevel(dir *Dir) {
	if dir.ReversedLevel >= len(DataByLevel) {
		for i:=len(DataByLevel); i<=dir.ReversedLevel; i++ {
			DataByLevel = append(DataByLevel, map[int64]map[*Dir]struct{}{})
		}
	}
	if _, exists := DataByLevel[dir.ReversedLevel][dir.TotalSize]; !exists {
		DataByLevel[dir.ReversedLevel][dir.TotalSize] = map[*Dir]struct{}{dir:struct{}{}}
	} else {
		//duplicate dir - actually
		DataByLevel[dir.ReversedLevel][dir.TotalSize][dir] = struct{}{}
	}
}

func DeleteFromDataByLevel(dir *Dir) {
	delete(DataByLevel[dir.ReversedLevel][dir.TotalSize], dir)
}

func listDir(dirPath string) (*Dir, int, int64) {
	//fmt.Printf(">>> %s\n", dirPath)
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return nil, 0, 0
	}

	if len(files) == 0 {
		/* empty dirs are not taken into account at all */
		//fmt.Printf("<<< %s (empty)\n", dirPath)
		return nil, 0, 0
	}

	dir := Dir{
		Name: path.Base(dirPath),
		SubDirs: []*Dir{}, 
		HighestOriginal: -1,
	}
	Dirs++
	for _, f := range files {
		dirPath := path.Join(dirPath, f.Name())
	      if f.IsDir() {
		      subDirP, level, totalSize := listDir(dirPath)			
		      if subDirP != nil {
			      subDirP.ParentDir = &dir
			      dir.SubDirs = append(dir.SubDirs, subDirP)
			      if level > dir.ReversedLevel {
				    dir.ReversedLevel = level
			      }
			      dir.TotalSize += totalSize
		      }	
	      } else {
			dir.TotalSize += f.Size()			
			dir.FileCount += 1
	      }
	}

	InsertToDataByLevel(&dir)

	//fmt.Printf("<<< %s: %+v\n", dirPath, dir)
	return &dir, dir.ReversedLevel+1, dir.TotalSize	
}


func main() {
	listDir(os.Args[1])
	//listDir("/media/rkociuba/500")
	//listDir("d:/")
	//listDir("./dir")	

	fmt.Printf("|LEVELS: %+v|\n", len(DataByLevel))
	fmt.Printf("|DIRS: %d|\n", Dirs)
	//fmt.Printf("|%+v|\n", DataByLevel)

	for i:= len(DataByLevel)-1; i >= 0; i-- {
	//for i:= len(DataByLevel)-1; i >= len(DataByLevel)-2; i-- {
		findLevelDuplicate(DataByLevel[i])	
		var a int
		fmt.Printf("LEV: %d\n", i)
    		fmt.Scanf("%d", &a)
	}

	for _, v := range SortByLevel(dups) {
		fmt.Printf("%s\n", absPath(v.Key))
		for _, vv := range(v.Value) {
			fmt.Printf("  Duplicate Group: %d\n", vv[0].ReversedLevel)
			for _, vvv := range(vv) {
				fmt.Printf("    %s\n", absPath(vvv))
			}
		}
	}	
}

