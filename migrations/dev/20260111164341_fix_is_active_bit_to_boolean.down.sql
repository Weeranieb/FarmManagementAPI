-- Revert is_active from boolean to bit in clients table
ALTER TABLE clients 
ALTER COLUMN is_active TYPE bit(1) 
USING CASE WHEN is_active THEN B'1'::bit(1) ELSE B'0'::bit(1) END;

-- Revert is_active from boolean to bit in active_ponds table
ALTER TABLE active_ponds 
ALTER COLUMN is_active TYPE bit(1) 
USING CASE WHEN is_active THEN B'1'::bit(1) ELSE B'0'::bit(1) END;

-- Revert is_active from boolean to bit in workers table
ALTER TABLE workers 
ALTER COLUMN is_active TYPE bit(1) 
USING CASE WHEN is_active THEN B'1'::bit(1) ELSE B'0'::bit(1) END;
