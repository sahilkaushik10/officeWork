
CREATE SCHEMA IF NOT EXISTS my_schema;

CREATE TABLE IF NOT EXISTS my_schema.leave_table (
    id SERIAL,
    name VARCHAR(100) NOT NULL,
    leave_type VARCHAR(20) NOT NULL,
    leave_from VARCHAR(20) NOT NULL,
    leave_to VARCHAR(20) NOT NULL,
    leave_to VARCHAR(20),
    team VARCHAR(50) NOT NULL,
    file VARCHAR(50),
    reporter VARCHAR(100) NOT NULL
);


CREATE TABLE IF NOT EXISTS my_schema.notifications (
    leave_id INT NOT NULL,
    reporting_manager VARCHAR(100) NOT NULL,
    approved BOOLEAN DEFAULT FALSE
);

CREATE OR REPLACE FUNCTION insert_notification() RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO my_schema.notifications (leave_id, reporting_manager, approved)
    VALUES (NEW.id, NEW.reporter, FALSE);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS notify_table_update ON my_schema.leave_table;
CREATE TRIGGER notify_table_update
AFTER INSERT ON my_schema.leave_table
FOR EACH ROW
EXECUTE FUNCTION insert_notification();
