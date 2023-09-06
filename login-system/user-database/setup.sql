CREATE TABLE users (
    xnEntryID SERIAL PRIMARY KEY,
    sPassword VARCHAR(255) NOT NULL,
    sUsername VARCHAR(255) UNIQUE NOT NULL
);

-- Insert default user
INSERT INTO users (sPassword, sUsername) VALUES ('$2a$10$gvAudYjysfI4zZTR.kRr5ezPV2qsYHCn.bzQtJ5ks4FSQYnb22gaq', 'user');