import { test, expect } from '../fixtures';

test.describe('Tareas', () => {
  test('cambiar el estado al editar', async ({ page, tasks }) => {
    await tasks([{ title: 'Escribir el informe' }]);
    await page.reload();

    // 1. Abre la edición: el selector tiene los seis estados y «Pendiente» elegida
    await page.getByRole('link', { name: 'Editar tarea: Escribir el informe' }).click();
    const select = page.getByRole('combobox', { name: 'Estado' });
    await expect(select.getByRole('option')).toHaveText([
      'Pendiente',
      'En curso',
      'Bloqueada',
      'Lista',
      'Finalizada',
      'Cancelada',
    ]);
    await expect(select).toHaveValue('pendiente');

    // 2. Elige «Bloqueada» y guarda
    await select.selectOption({ label: 'Bloqueada' });
    await page.getByRole('button', { name: 'Guardar cambios' }).click();

    // 3. La lista muestra la etiqueta nueva
    await expect(page.getByRole('status')).toHaveText('Tarea actualizada.');
    await expect(page.locator('.task-status')).toHaveText('Bloqueada');

    // 4. Al volver a editar, el selector conserva el estado
    await page.getByRole('link', { name: 'Editar tarea: Escribir el informe' }).click();
    await expect(page.getByRole('combobox', { name: 'Estado' })).toHaveValue('bloqueada');
    await expect(page.getByRole('combobox', { name: 'Estado' }).locator('option:checked')).toHaveText('Bloqueada');
  });
});
