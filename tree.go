package main

import (
        "fmt"
        "os"
        "io/fs"
        "sort"
        "path/filepath"
)

type DirEntryFlags uint64
const (
        Directory = DirEntryFlags(1 << iota)
        Expanded
        ReadError
)

type TreeData struct {
        RootEntry       *DirEntry
        Cache           map[DirEntry][]*DirEntry
        Lines           []*DirEntry
        ScrollIndex     int
        CurrentLine     int
}

type DirEntry struct {
        Name   string
        Parent *DirEntry
        Level  int
        Flags  DirEntryFlags
}

func (d *DirEntry) GetFullPath() string {
        if d.Parent == nil {
                path, _ := filepath.Abs(d.Name)
                return path
        }

        path := d.Parent.GetFullPath()
        path = filepath.Join(path, d.Name)
        return path
}

func (d *DirEntry) Expand() ([]*DirEntry, error) {
        path := d.GetFullPath()
        fsys := os.DirFS(path)
        file, err := fsys.Open(path)
        if err != nil {
                d.Flags |= ReadError
                return nil, err
        }

        dir, ok := file.(fs.ReadDirFile)
        if !ok {
                d.Flags |= ReadError
                return nil, fmt.Errorf("Can't read directory %s - not implemented", path)
        }

        list, err := dir.ReadDir(-1)
        if err != nil {
                d.Flags |= ReadError
                return nil, err
        }

        var out []*DirEntry
        for _, e := range list {
                de := &DirEntry{}
                de.Name   = e.Name()
                de.Parent = d.Parent
                de.Level  = d.Level + 1
                if e.IsDir() {
                        de.Flags |= Directory
                }
                out = append(out, de)
        }

        sort.Slice(out, func(i, j int) bool {
                if out[i].IsDir() && out[j].IsDir() {
                        return out[i].Name < out[j].Name
                } else if out[i].IsDir() {
                        return true
                } else if out[j].IsDir() {
                        return false
                } else {
                        return out[i].Name < out[j].Name
                }
        })
        d.Flags |= Expanded
        return out, nil
}

func (d *DirEntry) Collapse() {
        d.Flags &= ^Expanded
}

func (d *DirEntry) IsDir() bool {
        return (d.Flags & Directory) != 0
}

func (t *TreeData) ExpandEntry(d *DirEntry, refresh bool) {
        if !d.IsDir() {
                return
        }

        if !refresh {
                _, ok := t.Cache[*d]
                if ok {
                        return
                }
        }

        out, err := d.Expand()
        if err != nil {
                d.Flags |= ReadError
                return
        }
        t.Cache[*d] = out
        t.update(d)
}

func (t *TreeData) CollapseEntry(d *DirEntry) {
        if !d.IsDir() {
                return
        }
        d.Collapse()
        t.update(d)
}

func (t *TreeData) update(d *DirEntry) {
        if len(t.Lines) == 0 {
                t.Lines = append(t.Lines, d)
                if (d.Flags & Expanded) != 0 {
                        ch := t.Cache[*d]
                        t.Lines = append(t.Lines, ch...)
                }
        }

        idx := t.getEntryIndex(d)

        before := t.Lines[:idx+1]       //lines before entry in question including it
        after  := t.Lines[idx+1:]

        if (d.Flags & Expanded) == 0 {
                //entry collapsed
                ch := t.Cache[*d]
                gap := len(ch)
                if gap == len(after) {
                        //just cut the tail, there is nothing behind it
                        t.Lines = before
                } else {
                        after = after[gap:]
                        before = append(before, after...)
                        t.Lines = before
                }
        } else {
                //entry expanded
                ch := t.Cache[*d]
                before = append(before, ch...)
                if len(after) > 0 {
                        before = append(before, after...)
                }
                t.Lines = before
        }
}

func (t *TreeData) getEntryIndex(d *DirEntry) int {
        for i, e := range t.Lines {
                //first find correct level
                if e == d {
                        return i
                }
        }
        return -1
}
