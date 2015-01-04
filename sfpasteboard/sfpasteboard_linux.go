// +build linux

package main

import (
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

		// gtk_clipboard_request_image is not implemented. :(
		gtk.MainIterationDo(true)

		// Since we're not doing anything, set it to fail.
		os.Exit(1)
	}
}
