package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

type item struct {
	Task      string
	Done      bool
	CreatedAt time.Time
	CompletedAt time.Time
}

type Todos []item

func (t *Todos) Add(task string) {
	todo := item{
		Task: task,
		Done: false,
		CreatedAt: time.Now(),
		CompletedAt: time.Time{},
	}

	*t = append(*t, todo)
}

func (t *Todos) Complete(index int) error {
	list := *t
	if index <= 0 || index > len(list) {
		return errors.New("Invalid index!")
	}
	list[index-1].Done = true
	list[index-1].CompletedAt = time.Now()

	return nil
}

func (t *Todos) Delete(index int) error {
	list := *t
	if index <= 0 || index > len(list) {
		return errors.New("Invalid index!")
	}

	*t = append(list[:index-1], list[index:]...)

	return nil
}

func (t *Todos) Load(fileName string) error {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		} else {
			return err
		}
	}

	if len(file) == 0 {
		return err
	}

	err = json.Unmarshal(file, t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Todos) Store(fileName string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, data, 0644)
}

func (t *Todos) Print() {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell {
			{Align: simpletable.AlignCenter, Text: "#Id"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done?"},
			{Align: simpletable.AlignCenter, Text: "Created At"},
			{Align: simpletable.AlignCenter, Text: "Completed At"},
		},
	}

	var cells [][]*simpletable.Cell
	for index, task := range *t {
		index++
		cells = append(cells, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", index)},
			{Text: task.Task},
			{Text: fmt.Sprintf("%t", task.Done)},
			{Text: task.CreatedAt.Format(time.RFC822)},
			{Text: task.CompletedAt.Format(time.RFC822)},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell {
		{Align: simpletable.AlignCenter, Span: 5, Text: "Your todos are here!"},
	}}

	table.SetStyle(simpletable.StyleUnicode)

	table.Print()
}