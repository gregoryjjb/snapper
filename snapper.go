package snapper

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// UseAny instead of interface{}
const UseAny = true

// Snap outputs a snapshot of i to stdout
func Snap(i any, pkgAlias map[string]string) {
	Fsnap(os.Stdout, i, pkgAlias)
}

// Ssnap returns a string of a snapshot of i
func Ssnap(i any, pkgAlias map[string]string) string {
	buf := new(bytes.Buffer)
	Fsnap(buf, i, pkgAlias)
	return buf.String()
}

// Fsnap outputs a snapshot of i to the provided writer
func Fsnap(w io.Writer, i any, pkgAlias map[string]string) {
	var patterns []string
	for old, new := range pkgAlias {
		if new == "" {
			old = old + "."
		}

		patterns = append(patterns, old, new)
	}

	s := &snapper{
		w:            w,
		typeReplacer: strings.NewReplacer(patterns...),
	}

	s.snap(i, 0, false)
}

type snapper struct {
	w            io.Writer
	typeReplacer *strings.Replacer
}

func (s *snapper) write(str string) {
	io.WriteString(s.w, str)
}

func (s *snapper) snap(i any, indent int, omitStructName bool) {
	if i == nil {
		s.write("nil")
		return
	}

	// Easy literals that can just be printed
	switch v := i.(type) {
	case
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		bool:
		fmt.Fprintf(s.w, "%v", v)
		return

	case string:
		s.write(strconv.Quote(v))
		return
	}

	baseTabs := ""
	for i := 0; i < indent; i++ {
		baseTabs += "\t"
	}
	innerTabs := baseTabs + "\t"

	v := reflect.ValueOf(i)
	switch v.Kind() {

	// Structs
	case reflect.Struct:
		t := v.Type()

		if !omitStructName {
			name := t.String()
			name = s.typeReplacer.Replace(name)

			s.write(name)
		}
		s.write("{")

		printedSoFar := 0
		for i, field := range reflect.VisibleFields(t) {
			if !field.IsExported() {
				continue
			}

			name := t.Field(i).Name
			value := v.Field(i).Interface()

			// if printedSoFar > 0 {
			// 	fmt.Fprint(w, ", ")
			// }

			fmt.Fprintf(s.w, "\n%s%s: ", innerTabs, name)
			s.snap(value, indent+1, false)
			s.write(",")
			printedSoFar += 1
		}

		fmt.Fprintf(s.w, "\n%s}", baseTabs)

	// Slices
	case reflect.Slice, reflect.Array:
		t := v.Type()

		name := t.String()
		name = cleanEmptyInterface(name, UseAny)
		name = s.typeReplacer.Replace(name)

		s.write(name + "{")

		for i := 0; i < v.Len(); i++ {
			element := v.Index(i).Interface()
			fmt.Fprintf(s.w, "\n%s", innerTabs)
			s.snap(element, indent+1, true)
			s.write(",")
		}

		if v.Len() > 0 {
			s.write("\n")
			s.write(baseTabs)
		}
		s.write("}")

	// Maps
	case reflect.Map:
		t := v.Type()

		name := t.String()
		name = s.typeReplacer.Replace(name)
		s.write(name + "{")

		for _, key := range v.MapKeys() {
			// Render key
			s.write("\n" + innerTabs)
			s.snap(key.Interface(), indent+1, true)
			s.write(": ")

			// Render value
			element := v.MapIndex(key)
			s.snap(element.Interface(), indent+1, true)
			s.write(",")
		}

		if v.Len() > 0 {
			s.write("\n" + baseTabs)
		}
		s.write("}")

	// Pointers
	case reflect.Pointer:
		// We know it's not nil because we checked earlier
		s.write("&")
		s.snap(v.Elem().Interface(), indent, false) // Might be able to omit struct name here?
	}

	return
}

func cleanEmptyInterface(in string, useAny bool) string {
	if useAny {
		return strings.ReplaceAll(in, "interface {}", "any")
	}

	return strings.ReplaceAll(in, "interface {}", "interface{}")
}
