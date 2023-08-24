CREATE TABLE Users (
    xnEntryID SERIAL PRIMARY KEY,
    sPassword VARCHAR(255) NOT NULL,
    sUsername VARCHAR(255) UNIQUE NOT NULL
);

-- Insert default user
--INSERT INTO Users (sPassword, sUsername) VALUES ('password', 'user');