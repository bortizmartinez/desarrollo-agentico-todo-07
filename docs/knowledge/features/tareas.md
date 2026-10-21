# Lista de tareas

Lo que ve y puede hacer el usuario con sus tareas. Todo pasa en dos pantallas: la lista (`/`) y la edición de una tarea (`/tasks/{task_id}/edit`). El código está en `actions/tasks.go`, `models/task.go`, `templates/tasks/` y `public/assets/app.js`.

## La lista

- La pestaña y el enlace de la navegación lateral dicen «Tareas». El título de la página también.
- Debajo del título, el número de tareas: «3 tareas en tu lista», y en singular «1 tarea en tu lista».
- Las tareas salen por orden de creación, de la más antigua a la más reciente.
- Cada tarea muestra:
  - un botón redondo para completarla o reabrirla;
  - su título;
  - la etiqueta de su estado, siempre visible;
  - un botón para cambiar su estado, un enlace para editarla y un botón para eliminarla.
- Los botones solo tienen icono, y su nombre accesible incluye el título: «Completar tarea: Comprar pan», «Reabrir tarea: …», «Cambiar estado de la tarea: …», «Editar tarea: …», «Eliminar tarea: …».
- Sin tareas, la lista muestra «No hay tareas todavía. Añade la primera con el formulario.», y el número dice «0 tareas en tu lista».

## Estados

Cada tarea tiene un solo estado, uno de estos seis, que `models.TaskStatuses` define en este orden:

| Estado | Valor guardado | Significado |
| --- | --- | --- |
| Pendiente | `pendiente` | Aún no se ha empezado |
| En curso | `en_curso` | Se está trabajando en ella |
| Bloqueada | `bloqueada` | No puede avanzar por algo externo |
| Lista | `lista` | Terminada y lista para probar y validar |
| Finalizada | `finalizada` | Se ha validado y el resultado ha sido satisfactorio |
| Cancelada | `cancelada` | No se va a hacer |

- Una tarea nueva se crea «Pendiente»: lo pone el modelo, y el formulario de creación no puede elegir otro estado.
- `Task.Validate` rechaza cualquier otro valor, también una etiqueta en lugar del valor, con «El estado no es válido.».
- Solo «Finalizada» lleva el estilo de tarea terminada: título tachado y círculo relleno.

## Reglas del título

Las comprueba `Task.Validate` al crear y al editar:

| Regla | Mensaje si no se cumple |
| --- | --- |
| Es obligatorio. Un título solo con espacios cuenta como vacío | «El título es obligatorio.» |
| Tiene como mucho 200 caracteres (`models.TitleMaxLength`), contados sin los espacios del principio y del final | «El título no puede superar los 200 caracteres.» |

Antes de guardar, la acción quita los espacios del principio y del final. Los campos de título del navegador tienen también `maxlength="200"`.

## Añadir

- El formulario está debajo de la lista: el campo «Nueva tarea», con el texto de ayuda «Añadir una tarea…», y el botón «Añadir».
- Si el título es válido, la tarea se crea pendiente, aparece al final de la lista y sale el aviso «Tarea creada.».
- Si no, la lista se vuelve a mostrar con `422`: el campo conserva lo escrito, queda marcado como inválido y el mensaje sale debajo.

## Editar

- Se llega desde el enlace de editar de cada tarea. La pantalla tiene el enlace «← Volver a la lista», el título «Editar tarea» y el campo «Título» con el título actual. Debajo, la ayuda: «Hasta 200 caracteres. Los demás datos de la tarea no cambian.».
- Debajo del título está el selector «Estado», con los seis estados y el actual elegido. Solo se puede elegir uno.
- «Guardar cambios» cambia el título y el estado. Después vuelve a la lista con el aviso «Tarea actualizada.».
- Con un título o un estado inválido, la pantalla se vuelve a mostrar con `422` y el mensaje junto al campo que falla, conservando lo enviado. «Cancelar» vuelve a la lista sin cambiar nada.
- Si la tarea no existe, `404`.

## Completar y reabrir

- El botón redondo de cada tarea actúa con un solo clic, sin desplegar nada. En cualquier estado salvo «Finalizada» completa la tarea y la pone en «Finalizada»; en «Finalizada» la reabre y la pone en «Pendiente» (`Task.Toggle`).
- Sin JavaScript, el formulario envía `PUT /tasks/{task_id}/toggle`, la acción cambia el estado y vuelve a la lista. No hay aviso.
- Con JavaScript, `app.js` cambia la tarea en la página sin recargarla (la etiqueta, que pasa a «Finalizada» o «Pendiente», el estado del botón y su nombre accesible) y envía el formulario de verdad: un `POST` con `_method=PUT` y el token.

## Cambiar el estado desde la lista

- Cada tarea tiene, aparte, el botón «Cambiar estado de la tarea: …», que despliega los cuatro estados que no pone el botón de completar o reabrir: «En curso», «Bloqueada», «Lista» y «Cancelada» (`models.MenuStatuses`). Se ofrecen siempre, también el que ya tiene la tarea.
- Es un `<details>` nativo y funciona sin JavaScript. Cada opción es un formulario que envía `PUT /tasks/{task_id}/status` con el estado elegido; la acción lo guarda y vuelve a la lista.
- La acción acepta cualquiera de los seis estados; con otro valor responde `422`. Si la tarea no existe, `404`.

## Eliminar

- El botón de eliminar pide confirmación: «¿Seguro que quieres eliminar esta tarea?». Si se cancela, no pasa nada. Sin JavaScript, elimina sin preguntar.
- Al eliminar, la tarea desaparece de la lista y sale el aviso «Tarea eliminada.».
- Si la tarea no existe, `404`.

## Avisos

Los avisos salen arriba, en la primera pantalla que se muestra después de la acción, y se anuncian a los lectores de pantalla (`role="status"`):

| Acción | Aviso |
| --- | --- |
| Añadir | «Tarea creada.» |
| Editar | «Tarea actualizada.» |
| Eliminar | «Tarea eliminada.» |
| Completar o reabrir | Ninguno |
| Cambiar el estado desde la lista | Ninguno |
