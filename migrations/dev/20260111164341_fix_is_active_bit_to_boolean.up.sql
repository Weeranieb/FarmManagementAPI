-- Convert is_active from bit to boolean in clients table
ALTER TABLE clients 
ALTER COLUMN is_active TYPE boolean 
USING (is_active::text = '1');

-- Convert is_active from bit to boolean in active_ponds table
ALTER TABLE active_ponds 
ALTER COLUMN is_active TYPE boolean 
USING (is_active::text = '1');

-- Convert is_active from bit to boolean in workers table
ALTER TABLE workers 
ALTER COLUMN is_active TYPE boolean 
USING (is_active::text = '1');
