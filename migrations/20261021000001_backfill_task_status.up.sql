UPDATE tasks SET status = CASE WHEN completed THEN 'finalizada' ELSE 'pendiente' END;
