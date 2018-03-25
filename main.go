package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "path"
    "os"
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

func absPath(dir *Dir) string {
	if dir.ParentDir == nil {
		return ""
	} else {
		return path.Join(absPath(dir.ParentDir), dir.Name)
	}
}

type levelDup struct {
	level int
	duplicates string
}

type duplicate struct {
	trunk string
	levelDups map[int][]levelDup
}

var dups map[*Dir]duplicate //map: trunk-duplicates

func findHighestRoot(d *Dir) *Dir{
	for {
		if d.ParentDir == nil || p.HighestOriginal != -1 {
			return dir
		}
		d = d.ParentDir
	}
}

func findLevelDuplicate(dlev map[int64]map[*Dir]struct{}){
	for _, pdirs := range dlev {
		if len(pdirs) > 1 {
			// checik  if  one of (or more) dirs is subtree (duplicate tree) of bigger three *
			highest := -1
			pHihgst := nil
			for d, _ := range pdirs {
				//find highest original
				if highest > d.HighestOriginal {
					highest = d.HighestOriginal
					pHighest = d
				}
			}

			if pHighest == nil {
				//no trunk - this duplicate set is not subset of biggest duplicate set
				//choose new trunk
				trunk = pdirs[0] //TODO: for now 
				dups[trunk] = duplicate{trunk: absPath(dir)}
			} else {
				trunk = findHighestRoot(pHighest)		
			}

			i := 0
			for d, _ := range pdirs {
				if i == 0 && highest == -1 {
					// leave it only if its original (other words - mark one of idenitacl dirs - original)
					markSubDirsHighestOriginal(d, d.ReversedLevel)
				} else {
					deleteSubDirs(d)
				}
				fmt.Printf("%d: Duplicate Found: %s: (%+v)\n", i, absPath(d), d)	
				i++
			}

			dups = append(dups, duplicate)

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
	//for i:= len(DataByLevel)-1; i >= len(DataByLevel)-10; i-- {
		findLevelDuplicate(DataByLevel[i])	
		var a int
		fmt.Printf("LEV: %d\n", i)
    		fmt.Scanf("%d", &a)
	}
}

