package main

import (
	"testing"

	"github.com/gdamore/tcell/v3"
)

var gKeyTests = []struct {
	ev *tcell.EventKey
	s  string
}{
	{tcell.NewEventKey(tcell.KeyRune, "<", tcell.ModNone), "<lt>"},
	{tcell.NewEventKey(tcell.KeyRune, ">", tcell.ModNone), "<gt>"},
	{tcell.NewEventKey(tcell.KeyRune, " ", tcell.ModNone), "<space>"},
	{tcell.NewEventKey(tcell.KeyRune, "a", tcell.ModNone), "a"},
	{tcell.NewEventKey(tcell.KeyCtrlA, "", tcell.ModNone), "<c-a>"},
	{tcell.NewEventKey(tcell.KeyCtrlP, "", tcell.ModNone), "<c-p>"},
	{tcell.NewEventKey(tcell.KeyCtrlN, "", tcell.ModNone), "<c-n>"},
	{tcell.NewEventKey(tcell.KeyRune, "A", tcell.ModNone), "A"},
	{tcell.NewEventKey(tcell.KeyRune, "a", tcell.ModAlt), "<a-a>"},
	{tcell.NewEventKey(tcell.KeyLeft, "", tcell.ModNone), "<left>"},
	{tcell.NewEventKey(tcell.KeyLeft, "", tcell.ModCtrl), "<c-left>"},
	{tcell.NewEventKey(tcell.KeyLeft, "", tcell.ModShift), "<s-left>"},
	{tcell.NewEventKey(tcell.KeyLeft, "", tcell.ModAlt), "<a-left>"},
	{tcell.NewEventKey(tcell.KeyEsc, "", tcell.ModNone), "<esc>"},
	{tcell.NewEventKey(tcell.KeyF1, "", tcell.ModNone), "<f-1>"},
}

func TestReadKey(t *testing.T) {
	for _, test := range gKeyTests {
		if got := readKey(test.ev); got != test.s {
			t.Errorf("at input '%#v' expected '%s' but got '%s'", test.ev, test.s, got)
		}
	}
}

func TestParseKey(t *testing.T) {
	keyEqual := func(ev1, ev2 *tcell.EventKey) bool {
		return ev1.Key() == ev2.Key() && ev1.Modifiers() == ev2.Modifiers() && ev1.Str() == ev2.Str()
	}

	for _, test := range gKeyTests {
		if got := parseKey(test.s); !keyEqual(got, test.ev) {
			t.Errorf("at input '%s' expected '%#v' but got '%#v'", test.s, test.ev, got)
		}
	}
}

func TestCtrlPRemappable(t *testing.T) {
	// Verify that <c-p> is NOT bound by default in normal/visual mode
	// (only in command-line mode), allowing users to remap it freely
	
	// Check that <c-p> is in cmdkeys (command-line mode)
	if _, ok := gOpts.cmdkeys["<c-p>"]; !ok {
		t.Error("<c-p> should be bound in command-line mode by default")
	}
	
	// Check that <c-p> is NOT in nkeys (normal mode) by default
	if _, ok := gOpts.nkeys["<c-p>"]; ok {
		t.Error("<c-p> should NOT be bound in normal mode by default to allow remapping")
	}
	
	// Check that <c-p> is NOT in vkeys (visual mode) by default
	if _, ok := gOpts.vkeys["<c-p>"]; ok {
		t.Error("<c-p> should NOT be bound in visual mode by default to allow remapping")
	}
	
	// Verify that <c-n> has the same behavior
	if _, ok := gOpts.cmdkeys["<c-n>"]; !ok {
		t.Error("<c-n> should be bound in command-line mode by default")
	}
	if _, ok := gOpts.nkeys["<c-n>"]; ok {
		t.Error("<c-n> should NOT be bound in normal mode by default to allow remapping")
	}
	if _, ok := gOpts.vkeys["<c-n>"]; ok {
		t.Error("<c-n> should NOT be bound in visual mode by default to allow remapping")
	}
	
	// Test that we can remap <c-p> in normal mode
	testExpr := &callExpr{"paste", nil, 1}
	gOpts.nkeys["<c-p>"] = testExpr
	
	if expr, ok := gOpts.nkeys["<c-p>"]; !ok {
		t.Error("<c-p> should be mappable in normal mode")
	} else if expr != testExpr {
		t.Error("<c-p> should be mapped to the custom expression in normal mode")
	}
}
