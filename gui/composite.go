// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	//    "log"
	"os"
	"syscall"
)

import (
	"walk/drawing"
	. "walk/winapi/user32"
)

type Composite struct {
	Container
}

func NewCompositeWithStyle(parent IContainer, style uint) (*Composite, os.Error) {
	if parent == nil {
		return nil, newError("parent cannot be nil")
	}

	ensureMainWindowInitialized()

	hWnd := CreateWindowEx(
		WS_EX_CONTROLPARENT, syscall.StringToUTF16Ptr("Container_WindowClass"), nil,
		WS_CHILD|WS_VISIBLE|style,
		0, 0, 0, 0, parent.Handle(), 0, 0, nil)
	if hWnd == 0 {
		return nil, lastError("CreateWindowEx")
	}

	c := &Composite{Container: Container{Widget: Widget{hWnd: hWnd, parent: parent}}}

	c.children = newObservedWidgetList(c)

	c.SetFont(defaultFont)

	widgetsByHWnd[hWnd] = c

	parent.Children().Add(c)

	return c, nil
}

func NewComposite(parent IContainer) (*Composite, os.Error) {
	return NewCompositeWithStyle(parent, 0)
}

func (c *Composite) LayoutFlags() LayoutFlags {
	var flags LayoutFlags

	count := c.children.Len()
	if count == 0 {
		return ShrinkHorz | ShrinkVert
	} else {
		for i := 0; i < count; i++ {
			flags |= c.children.At(i).LayoutFlags()
		}
	}

	return flags
}

func (c *Composite) PreferredSize() drawing.Size {
	var maxW, maxH int

	count := c.children.Len()
	for i := 0; i < count; i++ {
		prefSize := c.children.At(i).PreferredSize()
		if prefSize.Width > maxW {
			maxW = prefSize.Width
		}
		if prefSize.Height > maxH {
			maxH = prefSize.Height
		}
	}

	if c.layout != nil {
		marg := c.layout.Margins()
		maxW += marg.Left + marg.Right
		maxH += marg.Top + marg.Bottom
	}

	return drawing.Size{maxW, maxH}
}

func (c *Composite) raiseEvent(msg *MSG) (err os.Error) {
	/*    switch msg.Message {
		case resizeMsgId:
	        log.Stdout("*Composite.raiseEvent: resizeMsgId")
	    }*/

	return c.Container.raiseEvent(msg)
}
