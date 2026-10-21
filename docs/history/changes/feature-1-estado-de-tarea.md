# Estado de tarea

**Issue:** #1 · **Tipo:** feature

## Resumen

Cada tarea tiene ahora un único estado, de entre seis (Pendiente, En curso, Bloqueada, Lista, Finalizada, Cancelada). Se elige al editar la tarea y se cambia desde la lista, con el botón de completar y reabrir de un clic y un desplegable aparte. La columna `status` sustituye a `completed`.

## Qué se construyó

- Columna `status` (obligatoria, por defecto `pendiente`) que reemplaza a `completed`; las tareas existentes se migran: completadas a «Finalizada», el resto a «Pendiente».
- Selector «Estado» con los seis estados en la pantalla de edición.
- Etiqueta con el estado de cada tarea en la lista; desaparece «Completada» y su lugar lo ocupa «Finalizada».
- Botón de completar (a «Finalizada») y reabrir (a «Pendiente») de un solo clic, que se mantiene.
- Botón aparte por tarea que despliega «En curso», «Bloqueada», «Lista» y «Cancelada».

## Implementación técnica

| Fichero | Papel |
| --- | --- |
| `migrations/20261021000000_add_status_to_tasks.{up,down}.fizz` | Añade `status varchar(20) NOT NULL default 'pendiente'`. |
| `migrations/20261021000001_backfill_task_status.{up,down}.sql` | Rellena `status` desde `completed`; la vuelta atrás hace `completed = (status = 'finalizada')`. |
| `migrations/20261021000002_remove_completed_from_tasks.{up,down}.fizz` | Elimina `completed`; la vuelta atrás la recrea con `default false`. |
| `models/task.go` | Constantes y listas `TaskStatuses` (seis) y `MenuStatuses` (cuatro), `BeforeValidate` (estado inicial), validación, `StatusLabel`, `IsDone`, `Toggle`. |
| `actions/tasks.go`, `actions/app.go` | `TasksUpdate` guarda también el estado; `TasksToggle` usa `Toggle()`; nueva `TasksUpdateStatus` en `PUT /tasks/{task_id}/status`. |
| `templates/tasks/edit.plush.html`, `templates/tasks/index.plush.html` | Selector, etiqueta y desplegable (`<details>`, un formulario por opción). |
| `public/assets/app.css`, `public/assets/app.js` | Estilos de etiqueta y desplegable; el JS actualiza la etiqueta al completar o reabrir. |
| `fixtures/tasks.toml`, `e2e/fixtures.ts` | Datos de prueba con `status` en lugar de `completed`. |
| `actions/tasks_test.go`, `models/task_test.go`, cuatro specs en `e2e/tareas/` | Tests de Go y de navegador. |
| `docs/knowledge/README.md`, `docs/knowledge/features/tareas.md`, `docs/knowledge/overview.md` | Documentación viva, actualizada por el propio cambio. |

Cambios que importan:

1. **`status` sustituye a `completed`.** Tres migraciones separadas (añadir, rellenar, quitar). Cualquier código o consulta que use `completed` deja de funcionar; se usa `Task.IsDone()`.
2. **La validez de los estados vive solo en la aplicación.** `Validate` rechaza cualquier valor fuera de `TaskStatuses` con «El estado no es válido.»; la base no tiene restricción `CHECK`.
3. **Estado inicial.** `BeforeValidate` pone `pendiente` solo a una tarea nueva (ID nulo) sin estado. La creación no lee el estado de la petición, así que crear con otro estado sigue dejando «Pendiente».
4. **El servidor acepta los seis estados en `PUT /tasks/{id}/status`; la página solo ofrece cuatro.** «Pendiente» y «Finalizada» las ponen el botón de completar o reabrir, pero la ruta no las excluye. Un estado inválido responde 422 con `c.Error` (sin formulario).
5. **Sin JavaScript todo funciona**: completar, reabrir y el desplegable son formularios `POST` con `_method=PUT` y token. Con JavaScript, `app.js` cambia la etiqueta al instante y envía el formulario con `fetch` (antes lo enviaba como GET y descartaba el resultado).

## Cómo se usa

- **Editar:** `/tasks/{id}/edit` muestra el selector «Estado»; al guardar se actualizan título y estado.
- **Completar / reabrir:** un clic en el botón de la tarea; «Finalizada» si estaba en cualquier otro estado, «Pendiente» si estaba finalizada.
- **Otro estado desde la lista:** botón del desplegable de la tarea y elegir «En curso», «Bloqueada», «Lista» o «Cancelada».

## Configuración

Hay tres migraciones nuevas; se aplican con el procedimiento habitual de migraciones (por ejemplo `buffalo pop migrate`). No hay variables ni ajustes nuevos.

**Vuelta atrás con pérdida.** Deshacer las migraciones recrea `completed` como `status = 'finalizada'`: las tareas en «En curso», «Bloqueada», «Lista» o «Cancelada» vuelven como «sin completar» y su estado se pierde. Al reaplicar, pasarían a «Pendiente».

## Cómo se probó

Acta de la verificación (intento 1 de 3, validado): <https://github.com/bortizmartinez/desarrollo-agentico-todo-07/issues/1#issuecomment-6323424341>. Los ocho criterios están en verde.

| Afirmación | Resultado |
| --- | --- |
| Estado único y un valor ajeno se rechaza (`archivada`, vacío) | 422 «El estado no es válido.» al editar y desde la lista; tarea intacta. |
| Tarea nueva en «Pendiente» | Verificado en el formulario y en un test con `cancelada` en la petición. |
| Migración de tareas existentes | Base con el esquema anterior y cinco tareas: completadas a «Finalizada», el resto a «Pendiente»; columna antigua eliminada. |
| Selector con seis estados, de elección única | Medido en la página. |
| El estado se guarda y se ve al reeditar | «Lista» y «Cancelada», también en la base. |
| Etiquetas en la lista, sin «Completada» | Seis estados correctos; 0 apariciones de «Completada». |
| Completar y reabrir con un clic | Desde los cinco estados; con y sin recargar, en la base y sin JavaScript. |
| Desplegable con los cuatro estados | Opciones ocultas hasta pulsar; guarda tras recargar y en la base. |
| Sin regresiones | Migraciones de test, `go test -p 1 ./...` y 9 tests de navegador en verde; `gofmt`, `go vet`, `go build` y `docker build` correctos. |

Límites de la prueba (según el acta):

- El único test automático de la migración ejecuta el SQL de relleno sobre una columna `completed` recreada a mano; no ejecuta los ficheros de migración ni la vuelta atrás.
- Los tests de Go en rojo sobre la base fallan solo por compilación.
- El test de navegador del desplegable solo elige «Lista»; las demás opciones las cubren los tests de Go y la medición manual.

Para repetir:

```
go test -p 1 ./...
npx playwright test
```

Los tests de Go necesitan la base de test con las migraciones aplicadas (puerta de migraciones de la CI; ver `docs/procedures/`).

## Notas

**Preguntas abiertas que el acta deja al humano** (ninguna bloquea):

1. **Texto de ayuda de la edición.** Debajo del título sigue «Hasta 200 caracteres. Los demás datos de la tarea no cambian.», y en esa pantalla ya se cambia el estado. La petición no fija el texto.
2. **Aspecto de las etiquetas y del botón del desplegable.** Hoy todas las etiquetas tienen el mismo aspecto y solo «Finalizada» se ve tachada; colores, orden e iconos son decisión de producto.
3. **Copia optimista del estado en `public/assets/app.js`.** El JS escribe «Finalizada» o «Pendiente» por su cuenta y descarta la respuesta de `fetch`. Si el guardado fallara, la página mostraría un estado no guardado. Un juicio lo marcó como bloqueo y otro lo refutó; queda como deuda: leer el estado de la respuesta del servidor.
4. **Lista del desplegable duplicada.** `MenuStatuses` en `models/task.go` repite a mano cuatro parejas valor/etiqueta de `TaskStatuses`; un renombrado puede desincronizarlas.

Otras observaciones del acta, menores:

- Las opciones del desplegable se llaman solo «Lista», «En curso»…, sin la tarea; un lector de pantalla no sabe a cuál se aplican.
- Un estado inválido enviado a mano desde la lista muestra la página de error de desarrollo de Buffalo; en la edición, el selector no marca ninguna opción.
- La columna no tiene restricción en la base de datos.
- Documentación fuera del alcance del caso y posiblemente desfasada: `docs/procedures/conventions.md:31` cita como ejemplo el comentario antiguo de `TasksToggle`. Además, `docs/knowledge/overview.md:63` generaliza que un error de validación vuelve a mostrar el formulario con 422, y la ruta nueva de la lista responde 422 sin formulario (`tareas.md` ya lo cuenta). No se han tocado.
- Para una persona (confianza baja en el acta): quitar la columna `completed`, que la petición no ordena pero implica «un solo estado», y corregir la llamada `fetch` de `app.js`, necesaria para que el criterio 7 guarde con JavaScript.

No hubo preguntas con opciones respondidas por el humano, así que no hay documento de decisiones.
