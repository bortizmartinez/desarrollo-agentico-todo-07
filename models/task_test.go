package models

import (
	"errors"
	"io/fs"
	"strings"

	"todo/migrations"

	"github.com/gobuffalo/pop/v6"
)

func (ms *ModelSuite) Test_Task_Create_Valid() {
	task := &Task{Title: "Comprar café"}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())

	ms.NotZero(task.ID)
	ms.Equal(StatusPending, task.Status)

	count, err := ms.DB.Count(&Task{})
	ms.NoError(err)
	ms.Equal(1, count)
}

func (ms *ModelSuite) Test_Task_Title_Required() {
	task := &Task{Title: ""}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.True(verrs.HasAny())
	ms.Contains(verrs.Get("title"), "El título es obligatorio.")

	count, err := ms.DB.Count(&Task{})
	ms.NoError(err)
	ms.Equal(0, count)
}

func (ms *ModelSuite) Test_Task_Title_OnlySpaces_Invalid() {
	task := &Task{Title: "   \t  "}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.True(verrs.HasAny())
	ms.Contains(verrs.Get("title"), "El título es obligatorio.")
}

func (ms *ModelSuite) Test_Task_Title_MaxLength() {
	task := &Task{Title: strings.Repeat("a", TitleMaxLength)}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())
}

func (ms *ModelSuite) Test_Task_Title_TooLong_Invalid() {
	task := &Task{Title: strings.Repeat("a", TitleMaxLength+1)}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.True(verrs.HasAny())
	ms.NotEmpty(verrs.Get("title"))
}

func (ms *ModelSuite) Test_Task_Toggle_Persists() {
	task := &Task{Title: "Revisar el correo"}
	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())

	task.Status = StatusDone
	ms.NoError(ms.DB.Update(task))

	reloaded := &Task{}
	ms.NoError(ms.DB.Find(reloaded, task.ID))
	ms.Equal(StatusDone, reloaded.Status)
}

func (ms *ModelSuite) Test_Task_Create_DefaultsToPending() {
	task := &Task{Title: "Sin estado"}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())

	reloaded := &Task{}
	ms.NoError(ms.DB.Find(reloaded, task.ID))
	ms.Equal(StatusPending, reloaded.Status)
}

func (ms *ModelSuite) Test_Task_Status_AllValid() {
	for _, s := range TaskStatuses {
		task := &Task{Title: "Tarea " + s.Label, Status: s.Value}

		verrs, err := ms.DB.ValidateAndCreate(task)
		ms.NoError(err)
		ms.False(verrs.HasAny(), s.Value)
	}

	count, err := ms.DB.Count(&Task{})
	ms.NoError(err)
	ms.Equal(len(TaskStatuses), count)
}

func (ms *ModelSuite) Test_Task_Status_Invalid() {
	for _, value := range []string{"archivada", "Finalizada"} {
		task := &Task{Title: "Estado inválido", Status: value}

		verrs, err := ms.DB.ValidateAndCreate(task)
		ms.NoError(err)
		ms.True(verrs.HasAny(), value)
		ms.Contains(verrs.Get("status"), "El estado no es válido.")
	}

	count, err := ms.DB.Count(&Task{})
	ms.NoError(err)
	ms.Equal(0, count)
}

func (ms *ModelSuite) Test_Task_Status_EmptyOnUpdate_Invalid() {
	task := &Task{Title: "Con estado"}
	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())

	task.Status = ""
	verrs, err = ms.DB.ValidateAndUpdate(task)
	ms.NoError(err)
	ms.True(verrs.HasAny())
	ms.Contains(verrs.Get("status"), "El estado no es válido.")

	reloaded := &Task{}
	ms.NoError(ms.DB.Find(reloaded, task.ID))
	ms.Equal(StatusPending, reloaded.Status)
}

func (ms *ModelSuite) Test_Task_Toggle() {
	for _, s := range TaskStatuses {
		task := &Task{Status: s.Value}
		task.Toggle()
		if s.Value == StatusDone {
			ms.Equal(StatusPending, task.Status)
		} else {
			ms.Equal(StatusDone, task.Status, s.Value)
		}
	}
}

func (ms *ModelSuite) Test_Task_StatusLabel() {
	labels := map[string]string{
		StatusPending:    "Pendiente",
		StatusInProgress: "En curso",
		StatusBlocked:    "Bloqueada",
		StatusReady:      "Lista",
		StatusDone:       "Finalizada",
		StatusCancelled:  "Cancelada",
	}
	ms.Len(labels, len(TaskStatuses))
	for value, label := range labels {
		ms.Equal(label, Task{Status: value}.StatusLabel())
	}
}

func (ms *ModelSuite) Test_Task_Migration_ExistingTasksGetStatus() {
	backfill, err := fs.ReadFile(migrations.FS(), "20261021000001_backfill_task_status.up.sql")
	ms.NoError(err)

	rollback := errors.New("deshacer")
	err = ms.DB.Transaction(func(tx *pop.Connection) error {
		ms.NoError(tx.RawQuery("ALTER TABLE tasks ADD COLUMN completed boolean NOT NULL DEFAULT false").Exec())
		ms.NoError(tx.RawQuery(`INSERT INTO tasks (id, title, completed, created_at, updated_at)
			VALUES (gen_random_uuid(), 'Completada', true, now(), now()),
			       (gen_random_uuid(), 'Abierta', false, now(), now())`).Exec())
		ms.NoError(tx.RawQuery(string(backfill)).Exec())

		statuses := map[string]string{}
		rows := []struct {
			Title  string `db:"title"`
			Status string `db:"status"`
		}{}
		ms.NoError(tx.RawQuery("SELECT title, status FROM tasks").All(&rows))
		for _, row := range rows {
			statuses[row.Title] = row.Status
		}
		ms.Equal(StatusDone, statuses["Completada"])
		ms.Equal(StatusPending, statuses["Abierta"])

		return rollback
	})
	ms.ErrorIs(err, rollback)
}
