import { test, expect } from '../fixtures';

test.describe('Tareas', () => {
  test('cada tarea muestra su estado como etiqueta', async ({ page, tasks }) => {
    await tasks([
      { title: 'Tarea pendiente', status: 'pendiente' },
      { title: 'Tarea en curso', status: 'en_curso' },
      { title: 'Tarea bloqueada', status: 'bloqueada' },
      { title: 'Tarea lista', status: 'lista' },
      { title: 'Tarea finalizada', status: 'finalizada' },
      { title: 'Tarea cancelada', status: 'cancelada' },
    ]);
    await page.reload();

    // 1. Las seis etiquetas, en el orden de la lista
    const labels = page.locator('.task-status');
    await expect(labels).toHaveText(['Pendiente', 'En curso', 'Bloqueada', 'Lista', 'Finalizada', 'Cancelada']);

    // 2. Todas son visibles, no solo para los lectores de pantalla
    for (const label of await labels.all()) {
      await expect(label).toBeVisible();
    }

    // 3. «Completada» ya no existe
    await expect(page.getByText('Completada')).toHaveCount(0);
  });
});
