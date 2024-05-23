package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type TransformHandler[In, Out any] func(input In) (Out, error)

type ETLPipe[I, O any] struct {
	transform TransformHandler[I, O]
}

func NewETLPipe[I, O any](transform TransformHandler[I, O]) *ETLPipe[I, O] {
	return &ETLPipe[I, O]{transform: transform}
}

func (t *ETLPipe[I, O]) PipeBatchJson(w io.Writer, r io.Reader) error {
	dec := json.NewDecoder(r)
	enc := json.NewEncoder(w)

	if _, err := dec.Token(); err != nil {
		return fmt.Errorf("failed to read the start of the JSON array: %w", err)
	}

	_, err := w.Write([]byte("["))
	if err != nil {
		return fmt.Errorf("failed to write the start of the JSON array: %w", err)
	}

	var raw json.RawMessage
	first := true

	for dec.More() {
		if err := dec.Decode(&raw); err != nil {
			return fmt.Errorf("failed to decode JSON object: %w", err)
		}
		var input I
		if err := json.Unmarshal(raw, &input); err != nil {
			return fmt.Errorf("failed to unmarshal input JSON: %w", err)
		}
		output, err := t.transform(input)
		if err != nil {
			return fmt.Errorf("failed to perform ETL: %w", err)
		}
		if !first {
			_, err := w.Write([]byte(","))
			if err != nil {
				return fmt.Errorf("failed to write comma: %w", err)
			}
		} else {
			first = false
		}
		if err := enc.Encode(output); err != nil {
			return fmt.Errorf("failed to encode results JSON: %w", err)
		}
	}

	_, err = w.Write([]byte("]"))
	if err != nil {
		return fmt.Errorf("failed to write the end of the JSON array: %w", err)
	}

	return nil
}

type JsonBatchPipe[I, O any] struct {
	inner     io.Reader
	decoder   *json.Decoder
	transform TransformHandler[I, O]
	firstRead bool
	lastRead  bool
}

func NewJsonBatchPipe[I, O any](inner io.Reader, transform TransformHandler[I, O]) (*JsonBatchPipe[I, O], error) {
	dec := json.NewDecoder(inner)

	if _, err := dec.Token(); err != nil {
		return nil, fmt.Errorf("failed to read the start of the JSON array: %w", err)
	}

	return &JsonBatchPipe[I, O]{
		inner:     inner,
		decoder:   dec,
		transform: transform,
		firstRead: false,
		lastRead:  false,
	}, nil
}

func (p *JsonBatchPipe[I, O]) Read(b []byte) (int, error) {
	if !p.firstRead {
		p.firstRead = true
		b = []byte("[")
		return len(b), nil
	}

	if p.decoder.More() {
		var raw json.RawMessage
		if err := p.decoder.Decode(&raw); err != nil {
			return 0, fmt.Errorf("failed to decode JSON object: %w", err)
		}
		var input I
		if err := json.Unmarshal(raw, &input); err != nil {
			return 0, fmt.Errorf("failed to unmarshal input JSON: %w", err)
		}
		output, err := p.transform(input)
		if err != nil {
			return 0, fmt.Errorf("failed to perform ETL: %w", err)
		}
		//// may not be proper implementation to copy output to b
		//b, err = json.Marshal(output)

		encodedOutput, err := json.Marshal(output)
		if err != nil {
			return 0, fmt.Errorf("failed to encode results JSON: %w", err)
		}

		// Copy the data from the temporary buffer to the provided buffer (b)
		n := copy(b, encodedOutput)

		// Check if the provided buffer was large enough to hold the data
		if n < len(encodedOutput) {
			return n, io.ErrShortBuffer
		}

		return n, nil
	}

	if !p.lastRead {
		p.lastRead = true
		b = []byte("]")
		return len(b), nil
	}

	return 0, io.EOF
}
