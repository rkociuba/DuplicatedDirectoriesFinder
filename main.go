package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "path"
)

var DataByLevel []map[int64]map[*Dir]struct{}

func main() {
	FirstDir, _, _:= listDir("./dir")	
	fmt.Println(FirstDir)
	fmt.Printf("|%+v|", DataByLevel)
}

type Dir struct {
	ParentDir *Dir
	SubDirs []*Dir
	DirName string
	FileCount  int
	ReversedLevel int
	TotalSize int64
}

	var dirIDP = 0 

	func listDir(dirPath string) (*Dir, int, int64) {
		fmt.Printf(">>> %s\n", dirPath)
		files, err := ioutil.ReadDir(dirPath)
		if err != nil {
			log.Fatal(err)
		}

		if len(files) == 0 {
			/* empty dirs are not taken into account at all */
			fmt.Printf("<<< %s (empty)\n", dirPath)
			return nil, 0, 0
		}

		dir := Dir{
			DirName: path.Base(dirPath),
			SubDirs: []*Dir{}, 
		}
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

		fmt.Printf("<<< %s: %+v\n", dirPath, dir)
		return &dir, dir.ReversedLevel+1, dir.TotalSize	
	}
