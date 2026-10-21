# Cambios documentados

Cada caso que cierra el workflow `adw` deja aquí su documento de cierre: qué se construyó, cómo se usa, cómo se probó y qué conviene saber para trabajar sobre ello. Lo escribe la skill `document` en la rama del caso, y viaja en el PR con el código.

Este índice es condicional: cada entrada dice cuándo conviene cargar su documento. Se llega a él desde [`docs/README.md`](../../README.md): quien planifica, como las skills `feature`, `bug` y `chore`, lee solo los documentos cuyas condiciones encajan con el encargo.

Un documento de cierre cuenta qué pasó en su caso, y no se reescribe cuando otro caso cambia el mismo código. Lo vigente está en `docs/knowledge/` y `docs/procedures/`.

## Formato de las entradas

```markdown
- [<tipo>-<N>-<slug>.md](<tipo>-<N>-<slug>.md): <qué cambió, en una línea>
  - Cuando se trabaje con <rutas, acciones o plantillas concretas>
  - Cuando se modifique <modelo, campo o regla relacionados>
  - Cuando se investigue <un problema concreto de esa área>
```

Entre tres y cinco condiciones por documento, y las entradas en orden por número de issue. Si un documento ya tiene entrada, se actualizan sus condiciones en lugar de añadir otra.

## Índice

- [feature-1-estado-de-tarea.md](feature-1-estado-de-tarea.md): cada tarea tiene un estado de seis, que sustituye a `completed`, con selector al editar, etiqueta, completar y reabrir y desplegable en la lista
  - Cuando se trabaje con `models/task.go`, `actions/tasks.go` (`TasksUpdate`, `TasksToggle`, `TasksUpdateStatus`) o la ruta `PUT /tasks/{task_id}/status`
  - Cuando se modifique la lista de tareas (`templates/tasks/index.plush.html`, `public/assets/app.js`, `public/assets/app.css`) o el selector de `templates/tasks/edit.plush.html`
  - Cuando se añada, renombre o filtre un estado de tarea, o se use la columna `status` o la antigua `completed`
  - Cuando se cree una migración de `tasks` o se revierta una (la vuelta atrás pierde estados) o se toquen `fixtures/tasks.toml` y `e2e/fixtures.ts`
  - Cuando se investigue un estado que no se guarda, una etiqueta que no coincide con la base o un 422 al cambiar de estado
