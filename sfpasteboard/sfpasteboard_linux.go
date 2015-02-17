// +build linux

package main

import (
	"fmt"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/gtk"
	"io/ioutil"
	"os"
)

// Linux GTK+ iteration of sfpasteboard clipboard cut/paste

func main() {
	gtk.Init(nil)

	clipboard := gtk.NewClipboardGetForDisplay(
		gdk.DisplayGetDefault(),
		gdk.SELECTION_CLIPBOARD)

	if len(os.Args) > 0 {
		// Write to temporary file, but use passed filename, otherwise the
		// GTK API will pass the wrong name
		tmpfile := os.TempDir() + string(os.PathSeparator) + os.Args[0]
		imgdata, _ := ioutil.ReadAll(os.Stdin)
		ioutil.WriteFile(tmpfile, imgdata, 0644)
		defer os.Remove(tmpfile)

		// Read from file
		pixbuf, _ := gdkpixbuf.NewPixbufFromFile(tmpfile)

		clipboard.SetImage(pixbuf)
		gtk.MainIterationDo(true)
		clipboard.Store()
		gtk.MainIterationDo(true)

		os.Exit(0)
	} else {
		// No argument, receive from clipboard

		// gtk_clipboard_request_image implemented in https://github.com/mattn/go-gtk/pull/224
		gtk.MainIterationDo(true)
		pixbuf := clipboard.WaitForText()
		gtk.MainIterationDo(true)

		// Save to file
		tmpfile := os.TempDir() + string(os.PathSeparator) + os.Args[0]
		// err :=
		pixbuf.Save(tmpfile, "image/png", "")
		// Attempt to replicate sfpasteboard behavior -- drop to stdout
		fmt.Println(tmpfile)

		// TODO: check for gtk error from pixbuf.Save, return 1

		// Since we're not doing anything, set it to fail.
		os.Exit(0)
	}
}
