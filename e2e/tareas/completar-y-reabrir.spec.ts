import { test, expect } from '../fixtures';

test.describe('Tareas', () => {
  test('completar y reabrir una tarea con un solo clic', async ({ page, tasks }) => {
    await tasks([{ title: 'Preparar la demo', status: 'en_curso' }]);
    await page.reload();

    const status = page.locator('.task-status');
    await expect(status).toHaveText('En curso');

    // 1. Completar: queda «Finalizada» y el botón pasa a reabrir, sin desplegar nada
    await page.getByRole('button', { name: 'Completar tarea: Preparar la demo' }).click();
    await expect(status).toHaveText('Finalizada');
    await expect(page.getByRole('button', { name: 'Reabrir tarea: Preparar la demo' })).toBeVisible();
    for (const name of ['En curso', 'Bloqueada', 'Lista', 'Cancelada']) {
      await expect(page.getByRole('button', { name, exact: true })).toHaveCount(0);
    }

    // 2. Al recargar, el estado se ha guardado
    await page.reload();
    await expect(status).toHaveText('Finalizada');

    // 3. Reabrir: queda «Pendiente», también tras recargar
    await page.getByRole('button', { name: 'Reabrir tarea: Preparar la demo' }).click();
    await expect(status).toHaveText('Pendiente');
    await page.reload();
    await expect(status).toHaveText('Pendiente');
  });
});
