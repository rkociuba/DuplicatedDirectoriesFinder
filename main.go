package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "path"
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
	Duplicated bool
}

func deleteSubDirs(dir *Dir) {
	for _, sd := range dir.SubDirs {
		deleteSubDirs(sd)			
	}
	dir = nil
}

func markSubDirsDuplicated(dir *Dir) {
	dir.Duplicated = true
	for _, sd := range dir.SubDirs {
		markSubDirsDuplicated(sd)			
	}
}

func absPath(dir *Dir) string {
	if dir.ParentDir == nil {
		return ""
	} else {
		return path.Join(absPath(dir.ParentDir), dir.Name)
	}
}

func findLevelDuplicate(dlev map[int64]map[*Dir]struct{}){
	for _, pdirs := range dlev {
		if len(pdirs) > 1 {
			i := 0
			for d, _ := range pdirs {
				fmt.Printf("%d: Duplicate Found: %s: (%+v)\n", i, absPath(d), d)	
				if i == 0 {
					markSubDirsDuplicated(d)
				} else {
					deleteSubDirs(d)
				}
				i++
			}
			fmt.Println("")	
		} 	
	}
} 

func listDir(dirPath string) (*Dir, int, int64) {
	//fmt.Printf(">>> %s\n", dirPath)
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		/* empty dirs are not taken into account at all */
		//fmt.Printf("<<< %s (empty)\n", dirPath)
		return nil, 0, 0
	}

	dir := Dir{
		Name: path.Base(dirPath),
		SubDirs: []*Dir{}, 
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

	if dir.ReversedLevel >= len(DataByLevel) {
		for i:=len(DataByLevel); i<=dir.ReversedLevel; i++ {
			DataByLevel = append(DataByLevel, map[int64]map[*Dir]struct{}{})
		}
	}
	if _, exists := DataByLevel[dir.ReversedLevel][dir.TotalSize]; !exists {
		DataByLevel[dir.ReversedLevel][dir.TotalSize] = map[*Dir]struct{}{&dir:struct{}{}}
	} else {
		//duplicate dir - actually
		DataByLevel[dir.ReversedLevel][dir.TotalSize][&dir] = struct{}{}
	}

	//fmt.Printf("<<< %s: %+v\n", dirPath, dir)
	return &dir, dir.ReversedLevel+1, dir.TotalSize	
}

func main() {
	//FirstDir, _, _:= listDir("./dir")	
	listDir("/media/rkociuba/500")	
	//_, _, _= listDir("./dir")	
	fmt.Printf("|LEVELS: %+v|\n", len(DataByLevel))
	fmt.Printf("|DIRS: %d|\n", Dirs)
	//fmt.Println(FirstDir)
	//fmt.Printf("|%+v|", DataByLevel)

	//for i:= len(DataByLevel)-1; i >= 0; i-- {
	for i:= len(DataByLevel)-1; i >= len(DataByLevel)-10; i-- {
		findLevelDuplicate(DataByLevel[i])	
	}
}

