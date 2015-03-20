package main

import (
        "os"
        "fmt"
        "bytes"
        "io/ioutil"
        //"regexp"
        git "github.com/libgit2/git2go" //"gopkg.in/libgit2/git2go.v22"
        "github.com/mattn/go-gtk/gtk"
        )

func file_chooser(url string, c chan string) {
        //var reponame = regexp.MustCompile(`[^/]*$`)
        gtk.Init(nil)
        dialog := gtk.NewFileChooserDialog(
                fmt.Sprintln("Clone", url, "to..."),
                nil,
                gtk.FILE_CHOOSER_ACTION_CREATE_FOLDER,
                "Clone",
                gtk.RESPONSE_ACCEPT)
        dialog.Response(func() {
                target := dialog.GetFilename()
                dialog.Destroy()

                fmt.Println("Destination folder: ", target)

                c <- target
        })
        //not in gtk2 dialog.SetCurrentName(reponame.MatchString(url))
        dialog.Run()
        dialog.Destroy()
}

func main() {
        var rep *git.Repository
        var e error
        var path string

        if len(os.Args[1:]) != 1 {
                fmt.Println("wrong number of arguments\n")
                os.Exit(1)
        }

        path = os.Args[1]

        if !bytes.Equal([]byte(path[0:6]), []byte("git://")) &&
           !bytes.Equal([]byte(path[0:8]), []byte("https://")) {
                fmt.Println("not a git url\n")
                os.Exit(1)
        }

        c := make(chan string)

        go file_chooser(path, c)

        fmt.Println("Cloning", path)

        var t string
        t, e = ioutil.TempDir("", "git-open-uri-")
        if e != nil {
                fmt.Printf("couldn't create temporary directory: %v\n", e)
                os.Exit(1)
        }

        /* use bare so it takes less memory if /tmp is on tmpfs */
        rep, e = git.Clone(path, t, &git.CloneOptions{nil, nil, true, "", nil, nil})
        if e != nil {
                fmt.Printf("couldn't clone repository: %v\n", e)
                os.Exit(1)
        }

        fmt.Println("Cloned", path, "to", t)

        target := <-c

        if target != nil {
                /* this un-bares it */
                rep, e = git.Clone(rep.Path(), target, nil)
                if e != nil {
                        fmt.Printf("couldn't clone repository: %v\n", e)
                        os.Exit(1)
                }
        }

        e = os.RemoveAll(t)
        if e != nil {
                fmt.Printf("failed to remove %s: %v\n", t, e)
                os.Exit(1)
        }
}
