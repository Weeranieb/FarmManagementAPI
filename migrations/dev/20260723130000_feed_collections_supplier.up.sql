-- Optional supplier / source of a feed (ผู้ขาย / แหล่งที่มา) — free text, e.g.
-- "CP Foods", "เบทาโกร", "ตลาดสี่มุมเมือง". Nullable: not every feed records one.
ALTER TABLE feed_collections ADD COLUMN supplier TEXT;
