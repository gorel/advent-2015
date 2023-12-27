package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/golang-collections/collections/set"
)

type Operation string

const (
	CONSTANT Operation = ""
	AND      Operation = "AND"
	OR       Operation = "OR"
	LSHIFT   Operation = "LSHIFT"
	RSHIFT   Operation = "RSHIFT"
	NOT      Operation = "NOT"
)

func Ptr[T any](val T) *T {
	return &val
}

func OperationFromString(s string) *Operation {
	switch s {
	case "":
		return Ptr(CONSTANT)
	case "AND":
		return Ptr(AND)
	case "OR":
		return Ptr(OR)
	case "LSHIFT":
		return Ptr(LSHIFT)
	case "RSHIFT":
		return Ptr(RSHIFT)
	case "NOT":
		return Ptr(NOT)
	default:
		return nil
	}
}

type Input struct {
	op        Operation
	operators []string
}

func NewInput(s string) *Input {
	// Split on spaces
	operators := strings.Split(s, " ")
	if len(operators) == 1 {
		return &Input{CONSTANT, operators}
	}

	var op *Operation
	for _, opStr := range operators {
		op = OperationFromString(opStr)
		if op != nil {
			break
		}
	}
	return &Input{*op, operators}
}

type Wire struct {
	name  string
	input *Input
	deps  []string

	value *uint16
}

func NewWire(line string) *Wire {
	dashIndex := strings.Index(line, "->")
	input := NewInput(line[:dashIndex-1])
	name := line[dashIndex+3:]
	var deps []string
	for _, operator := range input.operators {
		if OperationFromString(operator) == nil {
			if _, err := strconv.Atoi(operator); err != nil {
				deps = append(deps, operator)
			}
		}
	}
	return &Wire{
		name:  name,
		input: input,
		deps:  deps,
	}
}

func (w *Wire) Ready(toCompute *set.Set) bool {
	for _, dep := range w.deps {
		if toCompute.Has(dep) {
			return false
		}
	}
	return w.value == nil
}

func (w *Wire) getValue(name string, wires map[string]*Wire) uint16 {
	if val, err := strconv.Atoi(name); err == nil {
		return uint16(val)
	}
	return *wires[name].value
}

func (w *Wire) ComputeValue(wires map[string]*Wire) {
	val := uint16(0)
	switch w.input.op {
	case CONSTANT:
		val = w.getValue(w.input.operators[0], wires)
	case AND:
		val1 := w.getValue(w.input.operators[0], wires)
		val2 := w.getValue(w.input.operators[2], wires)
		val = val1 & val2
	case OR:
		val1 := w.getValue(w.input.operators[0], wires)
		val2 := w.getValue(w.input.operators[2], wires)
		val = val1 | val2
	case LSHIFT:
		val1 := w.getValue(w.input.operators[0], wires)
		val2 := w.getValue(w.input.operators[2], wires)
		val = val1 << val2
	case RSHIFT:
		val1 := w.getValue(w.input.operators[0], wires)
		val2 := w.getValue(w.input.operators[2], wires)
		val = val1 >> val2
	case NOT:
		input := w.getValue(w.input.operators[1], wires)
		val = ^input
	}

	w.value = &val
}

func (w *Wire) Reset() {
	w.value = nil
}

func ComputeAll(wires map[string]*Wire) {
	toCompute := set.New()
	for _, wire := range wires {
		toCompute.Insert(wire.name)
		wire.Reset()
	}

	for toCompute.Len() > 0 {
		for _, wire := range wires {
			if wire.Ready(toCompute) {
				wire.ComputeValue(wires)
				toCompute.Remove(wire.name)
			}
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	wires := make(map[string]*Wire)
	for scanner.Scan() {
		line := scanner.Text()
		wire := NewWire(line)
		wires[wire.name] = wire
	}

	ComputeAll(wires)

	aValue := *wires["a"].value
	fmt.Printf("Part 1: %d\n", aValue)

	wires["b"].input.operators = []string{strconv.Itoa(int(aValue))}
	ComputeAll(wires)
	aValue = *wires["a"].value
	fmt.Printf("Part 2: %d\n", aValue)
}
