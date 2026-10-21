package actions

import (
	"net/url"
	"strings"

	"todo/models"

	"github.com/gofrs/uuid"
)

func (as *ActionSuite) Test_TasksIndex_EmptyState() {
	res := as.HTML("/").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "No hay tareas todavía")
}

func (as *ActionSuite) Test_TasksIndex_ListsTasks() {
	as.LoadFixture("tres tareas")

	res := as.HTML("/").Get()
	as.Equal(200, res.Code)

	body := res.Body.String()
	as.Contains(body, "Comprar material de oficina")
	as.Contains(body, "Llamar al gestor")
	as.Contains(body, "Publicar la oferta de empleo")
	as.Contains(body, "Finalizada")
	as.Contains(body, "Pendiente")
	as.NotContains(body, "Completada")
}

func (as *ActionSuite) createTask(title, status string) *models.Task {
	task := &models.Task{Title: title, Status: status}
	verrs, err := as.DB.ValidateAndCreate(task)
	as.NoError(err)
	as.False(verrs.HasAny())
	return task
}

func (as *ActionSuite) Test_TasksIndex_ShowsEveryStatusLabel() {
	for _, s := range models.TaskStatuses {
		as.createTask("Tarea "+s.Value, s.Value)
	}

	res := as.HTML("/").Get()
	as.Equal(200, res.Code)
	for _, s := range models.TaskStatuses {
		as.Contains(res.Body.String(), `<span class="task-status">`+s.Label+`</span>`)
	}
}

func (as *ActionSuite) Test_TasksIndex_ToggleButtonByStatus() {
	as.createTask("En marcha", models.StatusInProgress)
	as.createTask("Terminada", models.StatusDone)

	body := as.HTML("/").Get().Body.String()
	as.Contains(body, "Completar tarea: En marcha")
	as.NotContains(body, "Reabrir tarea: En marcha")
	as.Contains(body, "Reabrir tarea: Terminada")
	as.NotContains(body, "Completar tarea: Terminada")
}

func (as *ActionSuite) Test_TasksIndex_StatusMenu() {
	task := as.createTask("Con menú", models.StatusPending)

	body := as.HTML("/").Get().Body.String()
	as.Contains(body, "Cambiar estado de la tarea: Con menú")
	as.Equal(4, strings.Count(body, `action="/tasks/`+task.ID.String()+`/status"`))
	for _, value := range []string{"en_curso", "bloqueada", "lista", "cancelada"} {
		as.Contains(body, `name="status" value="`+value+`"`)
	}
	as.NotContains(body, `name="status" value="pendiente"`)
	as.NotContains(body, `name="status" value="finalizada"`)
}

func (as *ActionSuite) Test_TasksCreate_Valid() {
	res := as.HTML("/tasks").Post(&models.Task{Title: "  Nueva tarea de prueba  "})
	as.Equal(303, res.Code)
	as.Equal("/", res.Header().Get("Location"))

	task := &models.Task{}
	as.NoError(as.DB.First(task))
	as.Equal("Nueva tarea de prueba", task.Title)
	as.Equal(models.StatusPending, task.Status)
}

func (as *ActionSuite) Test_TasksCreate_IgnoresStatus() {
	res := as.HTML("/tasks").Post(url.Values{"title": {"Con estado forzado"}, "status": {"cancelada"}})
	as.Equal(303, res.Code)

	task := &models.Task{}
	as.NoError(as.DB.First(task))
	as.Equal(models.StatusPending, task.Status)
}

func (as *ActionSuite) Test_TasksCreate_EmptyTitle_Invalid() {
	res := as.HTML("/tasks").Post(&models.Task{Title: ""})
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "El título es obligatorio.")

	count, err := as.DB.Count(&models.Task{})
	as.NoError(err)
	as.Equal(0, count)
}

func (as *ActionSuite) Test_TasksCreate_OnlySpaces_Invalid() {
	res := as.HTML("/tasks").Post(&models.Task{Title: "     "})
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "El título es obligatorio.")

	count, err := as.DB.Count(&models.Task{})
	as.NoError(err)
	as.Equal(0, count)
}

func (as *ActionSuite) Test_TasksToggle_CompletesAndReopens() {
	task := &models.Task{Title: "Pendiente de completar"}
	verrs, err := as.DB.ValidateAndCreate(task)
	as.NoError(err)
	as.False(verrs.HasAny())

	res := as.HTML("/tasks/%s/toggle", task.ID).Put(url.Values{})
	as.Equal(303, res.Code)

	as.NoError(as.DB.Reload(task))
	as.Equal(models.StatusDone, task.Status)

	res = as.HTML("/tasks/%s/toggle", task.ID).Put(url.Values{})
	as.Equal(303, res.Code)

	as.NoError(as.DB.Reload(task))
	as.Equal(models.StatusPending, task.Status)
}

func (as *ActionSuite) Test_TasksToggle_CompletesFromAnyStatus() {
	for _, status := range []string{models.StatusInProgress, models.StatusBlocked, models.StatusReady, models.StatusCancelled} {
		task := as.createTask("Desde "+status, status)

		res := as.HTML("/tasks/%s/toggle", task.ID).Put(url.Values{})
		as.Equal(303, res.Code)

		as.NoError(as.DB.Reload(task))
		as.Equal(models.StatusDone, task.Status, status)
	}
}

func (as *ActionSuite) Test_TasksToggle_NotFound() {
	id := uuid.Must(uuid.NewV4())
	res := as.HTML("/tasks/%s/toggle", id).Put(url.Values{})
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_TasksToggle_InvalidID() {
	res := as.HTML("/tasks/no-es-un-uuid/toggle").Put(url.Values{})
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_TasksDestroy() {
	task := &models.Task{Title: "Tarea que se elimina"}
	verrs, err := as.DB.ValidateAndCreate(task)
	as.NoError(err)
	as.False(verrs.HasAny())

	res := as.HTML("/tasks/%s", task.ID).Delete()
	as.Equal(303, res.Code)

	count, err := as.DB.Count(&models.Task{})
	as.NoError(err)
	as.Equal(0, count)
}

func (as *ActionSuite) Test_TasksDestroy_NotFound() {
	id := uuid.Must(uuid.NewV4())
	res := as.HTML("/tasks/%s", id).Delete()
	as.Equal(404, res.Code)
}

func (as *ActionSuite) taskForEditing(status string) *models.Task {
	task := &models.Task{Title: "Título original", Status: status}
	verrs, err := as.DB.ValidateAndCreate(task)
	as.NoError(err)
	as.False(verrs.HasAny())
	as.NoError(as.DB.Reload(task))
	return task
}

func (as *ActionSuite) Test_TasksIndex_EditLink() {
	task := as.taskForEditing(models.StatusPending)
	res := as.HTML("/").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "/tasks/"+task.ID.String()+"/edit")
	as.Contains(res.Body.String(), "Editar tarea: Título original")
}

func (as *ActionSuite) Test_TasksEdit_PrefillsTitle() {
	task := as.taskForEditing(models.StatusDone)
	res := as.HTML("/tasks/%s/edit", task.ID).Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), `value="Título original"`)
	as.Contains(res.Body.String(), `value="PUT"`)
	as.Contains(res.Body.String(), "Guardar cambios")
	as.Contains(res.Body.String(), "Cancelar")
	as.NoError(as.DB.Reload(task))
	as.Equal(models.StatusDone, task.Status)
}

func (as *ActionSuite) Test_TasksEdit_StatusSelect() {
	task := as.taskForEditing(models.StatusBlocked)
	body := as.HTML("/tasks/%s/edit", task.ID).Get().Body.String()
	as.Contains(body, `<select id="task-status" name="status"`)

	last := -1
	for _, s := range models.TaskStatuses {
		option := `<option value="` + s.Value + `"`
		i := strings.Index(body, option)
		as.True(i > last, s.Value)
		last = i
	}
	as.Contains(body, `<option value="bloqueada" selected>Bloqueada</option>`)
	as.Equal(1, strings.Count(body, "selected>"))
}

func (as *ActionSuite) Test_TasksEdit_NotFound() {
	res := as.HTML("/tasks/%s/edit", uuid.Must(uuid.NewV4())).Get()
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_TasksEdit_InvalidID() {
	res := as.HTML("/tasks/no-es-un-uuid/edit").Get()
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_TasksEdit_EscapesTitle() {
	task := &models.Task{Title: `<img src=x onerror="alert(1)">`}
	verrs, err := as.DB.ValidateAndCreate(task)
	as.NoError(err)
	as.False(verrs.HasAny())
	res := as.HTML("/tasks/%s/edit", task.ID).Get()
	as.Equal(200, res.Code)
	as.NotContains(res.Body.String(), `<img src=x`)
	as.Contains(res.Body.String(), "&lt;img")
}

func (as *ActionSuite) Test_TasksUpdate_PersistsTitleAndStatus() {
	task := as.taskForEditing(models.StatusDone)
	id, createdAt, updatedAt := task.ID, task.CreatedAt, task.UpdatedAt
	res := as.HTML("/tasks/%s", id).Put(url.Values{
		"title": {"  Título actualizado  "}, "status": {models.StatusBlocked},
		"id": {uuid.Must(uuid.NewV4()).String()}, "created_at": {"2000-01-01T00:00:00Z"},
	})
	as.Equal(303, res.Code)
	as.Equal("/", res.Header().Get("Location"))
	as.NoError(as.DB.Reload(task))
	as.Equal("Título actualizado", task.Title)
	as.Equal(models.StatusBlocked, task.Status)
	as.Equal(id, task.ID)
	as.True(task.CreatedAt.Equal(createdAt))
	as.False(task.UpdatedAt.Before(updatedAt))
	count, err := as.DB.Count(&models.Task{})
	as.NoError(err)
	as.Equal(1, count)
}

func (as *ActionSuite) Test_TasksUpdate_KeepsStatusSent() {
	task := as.taskForEditing(models.StatusPending)
	res := as.HTML("/tasks/%s", task.ID).Put(url.Values{"title": {"Nuevo título"}, "status": {models.StatusPending}})
	as.Equal(303, res.Code)
	as.NoError(as.DB.Reload(task))
	as.Equal("Nuevo título", task.Title)
	as.Equal(models.StatusPending, task.Status)
}

func (as *ActionSuite) Test_TasksUpdate_InvalidStatus() {
	task := as.taskForEditing(models.StatusPending)
	res := as.HTML("/tasks/%s", task.ID).Put(url.Values{"title": {"Otro título"}, "status": {"archivada"}})
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "El estado no es válido.")
	as.Contains(res.Body.String(), `aria-describedby="status-error"`)
	as.NoError(as.DB.Reload(task))
	as.Equal("Título original", task.Title)
	as.Equal(models.StatusPending, task.Status)
}

func (as *ActionSuite) Test_TasksUpdate_ShowsSavedStatusOnEdit() {
	task := as.taskForEditing(models.StatusPending)
	res := as.HTML("/tasks/%s", task.ID).Put(url.Values{"title": {"Título original"}, "status": {models.StatusReady}})
	as.Equal(303, res.Code)

	res = as.HTML("/tasks/%s/edit", task.ID).Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), `<option value="lista" selected>Lista</option>`)
}

func (as *ActionSuite) Test_TasksUpdateStatus_Valid() {
	for _, s := range models.MenuStatuses {
		task := as.createTask("Cambia a "+s.Value, models.StatusPending)
		res := as.HTML("/tasks/%s/status", task.ID).Put(url.Values{"status": {s.Value}})
		as.Equal(303, res.Code)
		as.Equal("/", res.Header().Get("Location"))
		as.NoError(as.DB.Reload(task))
		as.Equal(s.Value, task.Status)
	}
}

func (as *ActionSuite) Test_TasksUpdateStatus_Invalid() {
	task := as.taskForEditing(models.StatusInProgress)
	res := as.HTML("/tasks/%s/status", task.ID).Put(url.Values{"status": {"archivada"}})
	as.Equal(422, res.Code)
	as.NoError(as.DB.Reload(task))
	as.Equal(models.StatusInProgress, task.Status)
}

func (as *ActionSuite) Test_TasksUpdateStatus_NotFound() {
	res := as.HTML("/tasks/%s/status", uuid.Must(uuid.NewV4())).Put(url.Values{"status": {models.StatusReady}})
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_TasksUpdateStatus_InvalidID() {
	res := as.HTML("/tasks/no-es-un-uuid/status").Put(url.Values{"status": {models.StatusReady}})
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_TasksUpdate_EmptyTitle() {
	task := as.taskForEditing(models.StatusPending)
	res := as.HTML("/tasks/%s", task.ID).Put(url.Values{"title": {""}, "status": {models.StatusPending}})
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "El título es obligatorio.")
	as.Contains(res.Body.String(), `aria-describedby="title-error"`)
	as.NoError(as.DB.Reload(task))
	as.Equal("Título original", task.Title)
}

func (as *ActionSuite) Test_TasksUpdate_OnlySpaces() {
	task := as.taskForEditing(models.StatusDone)
	res := as.HTML("/tasks/%s", task.ID).Put(url.Values{"title": {"   "}, "status": {models.StatusDone}})
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "El título es obligatorio.")
	as.NoError(as.DB.Reload(task))
	as.Equal("Título original", task.Title)
	as.Equal(models.StatusDone, task.Status)
}

func (as *ActionSuite) Test_TasksUpdate_TitleTooLong() {
	task := as.taskForEditing(models.StatusPending)
	title := strings.Repeat("á", models.TitleMaxLength+1)
	res := as.HTML("/tasks/%s", task.ID).Put(url.Values{"title": {title}, "status": {models.StatusPending}})
	as.Equal(422, res.Code)
	as.Contains(res.Body.String(), "El título no puede superar los 200 caracteres.")
	as.Contains(res.Body.String(), title)
	as.NoError(as.DB.Reload(task))
	as.Equal("Título original", task.Title)
}

func (as *ActionSuite) Test_TasksUpdate_TitleAtLimit() {
	task := as.taskForEditing(models.StatusPending)
	title := strings.Repeat("á", models.TitleMaxLength)
	res := as.HTML("/tasks/%s", task.ID).Put(url.Values{"title": {title}, "status": {models.StatusPending}})
	as.Equal(303, res.Code)
	as.NoError(as.DB.Reload(task))
	as.Equal(title, task.Title)
}

func (as *ActionSuite) Test_TasksUpdate_NotFound() {
	res := as.HTML("/tasks/%s", uuid.Must(uuid.NewV4())).Put(url.Values{"title": {"Título nuevo"}, "status": {models.StatusPending}})
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_TasksUpdate_InvalidID() {
	res := as.HTML("/tasks/no-es-un-uuid").Put(url.Values{"title": {"Título nuevo"}, "status": {models.StatusPending}})
	as.Equal(404, res.Code)
}
