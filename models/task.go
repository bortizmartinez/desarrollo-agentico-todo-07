package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// TitleMaxLength es el número máximo de caracteres permitido en el título.
const TitleMaxLength = 200

// Estados de una tarea, tal como se guardan en la columna status.
const (
	StatusPending    = "pendiente"
	StatusInProgress = "en_curso"
	StatusBlocked    = "bloqueada"
	StatusReady      = "lista"
	StatusDone       = "finalizada"
	StatusCancelled  = "cancelada"
)

// TaskStatus es un estado con la etiqueta que ve el usuario.
type TaskStatus struct {
	Value string
	Label string
}

// TaskStatuses son los seis estados, en el orden en que se muestran.
var TaskStatuses = []TaskStatus{
	{StatusPending, "Pendiente"},
	{StatusInProgress, "En curso"},
	{StatusBlocked, "Bloqueada"},
	{StatusReady, "Lista"},
	{StatusDone, "Finalizada"},
	{StatusCancelled, "Cancelada"},
}

// MenuStatuses son los estados que se eligen desde el desplegable de la lista:
// todos salvo Pendiente y Finalizada, que pone el botón de completar o reabrir.
var MenuStatuses = []TaskStatus{
	{StatusInProgress, "En curso"},
	{StatusBlocked, "Bloqueada"},
	{StatusReady, "Lista"},
	{StatusCancelled, "Cancelada"},
}

// Task representa una tarea de la lista.
type Task struct {
	ID        uuid.UUID `json:"id" db:"id" form:"-"`
	Title     string    `json:"title" db:"title" form:"title"`
	Status    string    `json:"status" db:"status" form:"-"`
	CreatedAt time.Time `json:"created_at" db:"created_at" form:"-"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" form:"-"`
}

// Tasks es una colección de Task.
type Tasks []Task

// BeforeValidate pone el estado inicial a una tarea nueva que no trae estado.
func (t *Task) BeforeValidate(tx *pop.Connection) error {
	if t.ID == uuid.Nil && t.Status == "" {
		t.Status = StatusPending
	}
	return nil
}

// Validate se ejecuta en cada llamada a los métodos "pop.Validate*".
func (t *Task) Validate(tx *pop.Connection) (*validate.Errors, error) {
	verrs := validate.NewErrors()

	title := strings.TrimSpace(t.Title)
	if title == "" {
		verrs.Add("title", "El título es obligatorio.")
	} else if len([]rune(title)) > TitleMaxLength {
		verrs.Add("title", fmt.Sprintf("El título no puede superar los %d caracteres.", TitleMaxLength))
	}

	if !validStatus(t.Status) {
		verrs.Add("status", "El estado no es válido.")
	}

	return verrs, nil
}

// StatusLabel devuelve la etiqueta del estado de la tarea.
func (t Task) StatusLabel() string {
	for _, s := range TaskStatuses {
		if s.Value == t.Status {
			return s.Label
		}
	}
	return t.Status
}

// IsDone indica si la tarea está finalizada.
func (t Task) IsDone() bool {
	return t.Status == StatusDone
}

// Toggle completa la tarea, o la reabre si ya estaba finalizada.
func (t *Task) Toggle() {
	if t.IsDone() {
		t.Status = StatusPending
	} else {
		t.Status = StatusDone
	}
}

func validStatus(value string) bool {
	for _, s := range TaskStatuses {
		if s.Value == value {
			return true
		}
	}
	return false
}
