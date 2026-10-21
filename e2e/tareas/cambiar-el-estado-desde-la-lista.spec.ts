import { test, expect } from '../fixtures';

test.describe('Tareas', () => {
  test('cambiar el estado desde la lista', async ({ page, tasks }) => {
    await tasks([{ title: 'Revisar el contrato', status: 'pendiente' }]);
    await page.reload();

    // 1. Despliega los estados con el botón aparte
    await page.getByLabel('Cambiar estado de la tarea: Revisar el contrato').click();

    // 2. Ofrece exactamente los cuatro estados restantes
    const options = page.locator('.status-menu-option');
    await expect(options).toHaveText(['En curso', 'Bloqueada', 'Lista', 'Cancelada']);
    for (const option of await options.all()) {
      await expect(option).toBeVisible();
    }

    // 3. Elige «Lista»: la etiqueta cambia y se conserva al recargar
    await page.getByRole('button', { name: 'Lista', exact: true }).click();
    await expect(page.locator('.task-status')).toHaveText('Lista');
    await page.reload();
    await expect(page.locator('.task-status')).toHaveText('Lista');
  });
});
