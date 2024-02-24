CREATE TYPE bank_names AS ENUM('KBANK', 'BBL', 'KTB', 'BAY', 'CIMB', 'TTB', 'SCB', 'GSB');

CREATE TYPE registered_types AS ENUM('EMAIL', 'GOOGLE');

CREATE TYPE appointments_status AS ENUM('PENDING', 'APPROVED', 'REJECTED', 'REQUEST_CHANGE', 'CANCELLED', 'COMPLETED');

CREATE TYPE card_colors AS ENUM('LIGHT_BLUE', 'BLUE', 'DARK_BLUE', 'VERY_DARK_BLUE');

CREATE TYPE property_types AS ENUM('CONDOMINIUM', 'APARTMENT', 'SEMI_DETACHED_HOUSE', 'HOUSE', 'SERVICED_APARTMENT', 'TOWNHOUSE');

CREATE TYPE furnishing AS ENUM('UNFURNISHED', 'PARTIALLY_FURNISHED', 'FULLY_FURNISHED', 'READY_TO_MOVE_IN');

CREATE TYPE floor_size_units AS ENUM('SQM', 'SQFT');

CREATE TABLE email_verification_codes
(
    email                     VARCHAR(50) PRIMARY KEY           NOT NULL,
    code                      VARCHAR(20)                       NOT NULL,
    expired_at                TIMESTAMP(0) WITH TIME ZONE       NOT NULL
);

CREATE TABLE google_oauth_states (
    code       UUID PRIMARY KEY            NOT NULL,
    expired_at TIMESTAMP(0) WITH TIME ZONE NOT NULL
);

CREATE TABLE users
(
    user_id                             UUID PRIMARY KEY                DEFAULT gen_random_uuid(),
    registered_type                     registered_types                NOT NULL,
    email                               VARCHAR(50)                     NOT NULL,
    password                            VARCHAR(64)                     DEFAULT NULL,
    first_name                          VARCHAR(50)                     NOT NULL,
    last_name                           VARCHAR(50)                     NOT NULL,
    phone_number                        VARCHAR(10)                     NOT NULL,
    profile_image_url                   VARCHAR(2000)                   DEFAULT NULL,
    is_verified                         BOOLEAN                         DEFAULT FALSE,
    created_at                          TIMESTAMP(0) WITH TIME ZONE     DEFAULT CURRENT_TIMESTAMP,
    updated_at                          TIMESTAMP(0) WITH TIME ZONE     DEFAULT CURRENT_TIMESTAMP,
    deleted_at                          TIMESTAMP(0) WITH TIME ZONE     DEFAULT NULL,
    UNIQUE(email, deleted_at),
    UNIQUE(phone_number, deleted_at)
);

CREATE TABLE user_financial_informations
(
    user_id                             UUID PRIMARY KEY REFERENCES users(user_id)  NOT NULL,
    bank_name                           bank_names                                  DEFAULT NULL,
    bank_account_number                 VARCHAR(10)                                 DEFAULT NULL,
    created_at                          TIMESTAMP(0) WITH TIME ZONE                 DEFAULT CURRENT_TIMESTAMP,
    updated_at                          TIMESTAMP(0) WITH TIME ZONE                 DEFAULT CURRENT_TIMESTAMP,
    deleted_at                          TIMESTAMP(0) WITH TIME ZONE                 DEFAULT NULL
);

CREATE TABLE credit_cards
(
    user_id                             UUID REFERENCES users (user_id) ON DELETE CASCADE,
    tag_number                          INTEGER CHECK(1 <= tag_number AND tag_number <= 4) NOT NULL,
    card_nickname                       VARCHAR(50)                     NOT NULL,
    cardholder_name                     VARCHAR(50)                     NOT NULL,
    card_number                         VARCHAR(16)                     NOT NULL,
    expire_month                        VARCHAR(2)                      NOT NULL,
    expire_year                         VARCHAR(4)                      NOT NULL,
    cvv                                 VARCHAR(3)                      NOT NULL,
    card_color                          card_colors                     DEFAULT 'LIGHT_BLUE',
    PRIMARY KEY (user_id, tag_number)
);

CREATE TABLE user_verifications 
(
    user_id                 UUID PRIMARY KEY NOT NULL REFERENCES users(user_id),
    citizen_id              VARCHAR(13)      NOT NULL,
    citizen_card_image_url  VARCHAR(2000)    NOT NULL,
    verified_at             TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE properties
(
    property_id              UUID PRIMARY KEY                                       DEFAULT gen_random_uuid(),
    owner_id                 UUID REFERENCES users (user_id) ON DELETE CASCADE      NOT NULL,
    property_name            VARCHAR(50)                                            NOT NULL,
    property_description     TEXT                                                   NOT NULL,
    property_type            property_types                                         NOT NULL,
    address                  VARCHAR(50)                                            NOT NULL,
    alley                    VARCHAR(50)                                            DEFAULT NULL,
    street                   VARCHAR(50)                                            NOT NULL,
    sub_district             VARCHAR(50)                                            NOT NULL,
    district                 VARCHAR(50)                                            NOT NULL,
    province                 VARCHAR(50)                                            NOT NULL,
    country                  VARCHAR(50)                                            NOT NULL,
    postal_code              CHAR(5)                                                NOT NULL,
    bedrooms                 INTEGER                                                NOT NULL,
    bathrooms                INTEGER                                                NOT NULL,
    furnishing               furnishing                                             NOT NULL,
    floor                    INTEGER                                                NOT NULL,
    floor_size               DOUBLE PRECISION                                       NOT NULL,
    floor_size_unit          floor_size_units                                       DEFAULT 'SQM',
    unit_number              INTEGER                                                NOT NULL,
    created_at               TIMESTAMP(0) WITH TIME ZONE                            DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP(0) WITH TIME ZONE                            DEFAULT CURRENT_TIMESTAMP,
    deleted_at               TIMESTAMP(0) WITH TIME ZONE                            DEFAULT NULL
);

CREATE TABLE property_images
(
    property_id UUID REFERENCES properties (property_id) ON DELETE CASCADE          NOT NULL,
    image_url       VARCHAR(2000)                                                   NOT NULL,
    created_at               TIMESTAMP(0) WITH TIME ZONE                            DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP(0) WITH TIME ZONE                            DEFAULT CURRENT_TIMESTAMP,
    deleted_at               TIMESTAMP(0) WITH TIME ZONE                            DEFAULT NULL,
    PRIMARY KEY (property_id, image_url)
);

CREATE TABLE selling_properties
(
    property_id UUID PRIMARY KEY REFERENCES properties (property_id) ON DELETE CASCADE  NOT NULL,
    price       DOUBLE PRECISION                                                        NOT NULL,
    is_sold     BOOLEAN                                                                 NOT NULL,
    created_at               TIMESTAMP(0) WITH TIME ZONE                                DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP(0) WITH TIME ZONE                                DEFAULT CURRENT_TIMESTAMP,
    deleted_at               TIMESTAMP(0) WITH TIME ZONE                                DEFAULT NULL
);

CREATE TABLE renting_properties
(
    property_id     UUID PRIMARY KEY REFERENCES properties (property_id) ON DELETE CASCADE  NOT NULL,
    price_per_month DOUBLE PRECISION                                                        NOT NULL,
    is_occupied     BOOLEAN                                                                 NOT NULL,
    created_at               TIMESTAMP(0) WITH TIME ZONE                                    DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP(0) WITH TIME ZONE                                    DEFAULT CURRENT_TIMESTAMP,
    deleted_at               TIMESTAMP(0) WITH TIME ZONE                                    DEFAULT NULL
);

CREATE TABLE favorite_properties
(
    user_id         UUID REFERENCES users (user_id)             ON DELETE CASCADE   NOT NULL,
    property_id     UUID REFERENCES properties (property_id)    ON DELETE CASCADE   NOT NULL,
    PRIMARY KEY (user_id, property_id)
);

CREATE TABLE appointments
(
    appointment_id      UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    property_id         UUID REFERENCES properties (property_id)   NOT NULL,
    owner_user_id       UUID REFERENCES users (user_id)            NOT NULL,
    dweller_user_id     UUID REFERENCES users (user_id)            NOT NULL,
    appointments_status appointments_status DEFAULT 'PENDING'      NOT NULL,
    appointment_date    TIMESTAMP(0) WITH TIME ZONE                NOT NULL,
    created_at          TIMESTAMP(0) WITH TIME ZONE                DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP(0) WITH TIME ZONE                DEFAULT CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP(0) WITH TIME ZONE                DEFAULT NULL,
    UNIQUE (property_id, appointment_date, deleted_at)
);

CREATE TABLE agreements
(
    agreement_id   UUID PRIMARY KEY DEFAULT gen_random_uuid()   NOT NULL,
    property_id      UUID REFERENCES properties (property_id)   NOT NULL,
    owner_user_id    UUID REFERENCES users (user_id)            NOT NULL,
    dweller_user_id  UUID REFERENCES users (user_id)            NOT NULL,
    agreement_date TIMESTAMP(0) WITH TIME ZONE                  NOT NULL,
    created_at       TIMESTAMP(0) WITH TIME ZONE                DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP(0) WITH TIME ZONE                DEFAULT CURRENT_TIMESTAMP,
    deleted_at       TIMESTAMP(0) WITH TIME ZONE                DEFAULT NULL,
    UNIQUE (property_id, agreement_date)
);

-------------------- RULES --------------------

CREATE RULE soft_deletion AS ON DELETE TO users DO INSTEAD (
    UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE user_id = old.user_id and deleted_at IS NULL
);

CREATE RULE soft_deletion AS ON DELETE TO user_financial_informations DO INSTEAD (
    UPDATE user_financial_informations SET deleted_at = CURRENT_TIMESTAMP WHERE user_id = old.user_id and deleted_at IS NULL
);

CREATE RULE soft_deletion AS ON DELETE TO properties DO INSTEAD (
    UPDATE properties SET deleted_at = CURRENT_TIMESTAMP WHERE property_id = old.property_id and deleted_at IS NULL
);

CREATE RULE soft_deletion AS ON DELETE TO property_images DO INSTEAD (
    UPDATE property_images SET deleted_at = CURRENT_TIMESTAMP WHERE property_id = old.property_id and deleted_at IS NULL
);

CREATE RULE soft_deletion AS ON DELETE TO selling_properties DO INSTEAD (
    UPDATE selling_properties SET deleted_at = CURRENT_TIMESTAMP WHERE property_id = old.property_id and deleted_at IS NULL
);

CREATE RULE soft_deletion AS ON DELETE TO renting_properties DO INSTEAD (
    UPDATE renting_properties SET deleted_at = CURRENT_TIMESTAMP WHERE property_id = old.property_id and deleted_at IS NULL
);

CREATE RULE soft_deletion AS ON DELETE TO appointments DO INSTEAD (
    UPDATE appointments SET deleted_at = CURRENT_TIMESTAMP WHERE appointment_id = old.appointment_id and deleted_at IS NULL
);

CREATE RULE delete_users AS ON UPDATE TO users
    WHERE old.deleted_at IS NULL AND new.deleted_at IS NOT NULL
    DO ALSO (
        UPDATE properties SET deleted_at = new.deleted_at WHERE owner_id = old.user_id;
        UPDATE appointments SET deleted_at = new.deleted_at WHERE owner_user_id = old.user_id OR dweller_user_id = old.user_id;
        DELETE FROM favorite_properties WHERE user_id = old.user_id;
        DELETE FROM user_verifications WHERE user_id = old.user_id;
        UPDATE user_financial_informations SET deleted_at = new.deleted_at WHERE user_id = old.user_id;
    );

CREATE RULE delete_properties AS ON UPDATE TO properties
    WHERE old.deleted_at IS NULL AND new.deleted_at IS NOT NULL
    DO ALSO (
        UPDATE property_images SET deleted_at = new.deleted_at WHERE property_id = old.property_id;
        UPDATE selling_properties SET deleted_at = new.deleted_at WHERE property_id = old.property_id;
        UPDATE renting_properties SET deleted_at = new.deleted_at WHERE property_id = old.property_id;
        DELETE FROM favorite_properties WHERE property_id = old.property_id;
        UPDATE appointments SET deleted_at = new.deleted_at WHERE property_id = old.property_id;
    );

CREATE RULE create_email_verification_codes AS ON INSERT TO email_verification_codes
    WHERE new.email = (SELECT email FROM email_verification_codes WHERE email = new.email) DO INSTEAD(
        UPDATE email_verification_codes SET code = new.code, expired_at = new.expired_at WHERE email = new.email
    );

CREATE RULE creat_user_financial_informations AS ON INSERT TO users DO ALSO (
    INSERT INTO user_financial_informations (user_id) VALUES (new.user_id)
);

CREATE RULE update_user_verified AS ON INSERT TO user_verifications DO ALSO (
    UPDATE users SET is_verified = TRUE WHERE user_id = new.user_id
);

-------------------- DUMMY DATA --------------------

INSERT INTO users (user_id, registered_type, email, password, first_name, last_name, phone_number, profile_image_url, is_verified) VALUES
('f38f80b3-f326-4825-9afc-ebc331626555', 'EMAIL', 'johnd@email.com', '$2a$10$eEkTbe/JskFiociJ8U/bGOwwiea9dZ6sN7ac9ZvuiUgtrekZ7b.ya', 'John', 'Doe', '1234567890', 'https://picsum.photos/200/300?random=1', TRUE),
('bc5891ce-d6f2-d6f2-d6f2-ebc331626555', 'EMAIL', 'sams@email.com', '$2a$10$eEkTbe/JskFiociJ8U/bGOwwiea9dZ6sN7ac9Zvuhfkdle9405.ya', 'Sam', 'Smith', '0987654321', NULL, FALSE),
('62dd40da-f326-4825-9afc-2d68e06e0282', 'GOOGLE', 'gmail@gmail.com', NULL, 'C', 'C', '3333333333', 'https://picsum.photos/200/300?random=1', TRUE);

INSERT INTO properties (property_id, owner_id, property_name, property_description, property_type, address, alley, street, sub_district, district, province, country, postal_code, bedrooms, bathrooms, furnishing, floor, floor_size, floor_size_unit, unit_number) VALUES
('f38f80b3-f326-4825-9afc-ebc331626875', 'f38f80b3-f326-4825-9afc-ebc331626555', 'Et sequi dolor praes', 'sdfasdfdsalflvasdldk', 'HOUSE', 'Quas iusto expedita ', 'Delisa', 'Grace', 'Michael', 'Christine', 'Anthony', 'Andrew', '53086', 3, 2, 'UNFURNISHED', 20, 45.78, 'SQM', 1123),
('41a448d4-43ec-411a-a692-2d68e06e0282', 'f38f80b3-f326-4825-9afc-ebc331626555', 'Impedit quae itaque ', 'asludfowyegfubhsalas', 'APARTMENT', 'Sunt fuga quo perspi', 'Raquel', 'Brandy', 'Jacob', 'Lino', 'Edward', 'Reginald', '12894', 2, 1, 'FULLY_FURNISHED',18, 22.13, 'SQM', 1233),
('414854bf-bdee-45a5-929f-073aedaceea0', 'f38f80b3-f326-4825-9afc-ebc331626555', 'Architecto iure labo', 'asdasfhsfjdkaasdfjks', 'CONDOMINIUM', 'Pariatur temporibus ', 'Robert', 'Nancy', 'Barbara', 'David', 'Henry', 'David', '24264', 3, 2, 'READY_TO_MOVE_IN', 1, 200.00, 'SQFT', 555),
('62dd40da-8238-4d21-b9a7-7f1c24efdd0c', '62dd40da-f326-4825-9afc-2d68e06e0282', 'Optio in asperiores ', 'ioquwerewqpurwpqeruu', 'SEMI_DETACHED_HOUSE', 'Ea nobis mollitia ea', 'Tina', 'Linda', 'Ronald', 'Julia', 'Russell', 'William', '10287', 9, 9, 'FULLY_FURNISHED', 9, 90.99, 'SQM', 9909),
('bc5891ce-6d5e-40d6-8563-f7cebe9667e8', '62dd40da-f326-4825-9afc-2d68e06e0282', 'Sunt at totam animi ', 'iuwuerhihdfsiladfjas', 'TOWNHOUSE', 'Unde natus nesciunt ', 'Norma', 'Gregory', 'Donovan', 'Charles', 'Kevin', 'Tyrone', '10055', 1, 1, 'PARTIALLY_FURNISHED', 30, 90.99, 'SQFT', 1234),
('3df779f2-1f72-44d1-9a31-51929ed130a2', '62dd40da-f326-4825-9afc-2d68e06e0282', 'Animi vero ipsa nihi', 'hubgqewhbflasdhbfahs', 'HOUSE', 'Totam nam minus veni', 'Allen', 'Linda', 'Bobby', 'Nora', 'James', 'Lucinda', '01229', 7, 2, 'UNFURNISHED', 12, 127.27, 'SQFT', 1207),
('a8329428-6971-42e8-974a-4df030cd27be', '62dd40da-f326-4825-9afc-2d68e06e0282', 'Numquam sit dicta be', 'euyqrbdfhaivbhdbewjf', 'SERVICED_APARTMENT', 'Consequatur incidunt', 'Cecil', 'David', 'Nancy', 'Brandon', 'John', 'Lillian', '48668', 3, 2, 'FULLY_FURNISHED', 6, 66.00, 'SQFT', 6666),
('f8eaf2fc-d6f2-4a8c-a714-5425cc76bbfa', '62dd40da-f326-4825-9afc-2d68e06e0282', 'Iure nostrum ab reru', 'ewurblhdsfhladlhfdas', 'SEMI_DETACHED_HOUSE', 'Nisi officia nemo au', 'Keith', 'Joseph', 'Joseph', 'Goldie', 'Danika', 'Bernice', '47550', 1, 1, 'READY_TO_MOVE_IN', 4, 44.44, 'SQM', 4444),
('b7c8ce65-8fa3-4759-bc4e-42a396ef4fc1', '62dd40da-f326-4825-9afc-2d68e06e0282', 'Aut nemo incidunt ul', 'sldlfghewrvjdsbppppp', 'CONDOMINIUM', 'Porro molestias rati', 'Brian', 'Gregory', 'Geraldine', 'Edward', 'Charles', 'James', '97186', 3, 1, 'UNFURNISHED', 13, 1313.13, 'SQFT', 1313);

INSERT INTO property_images (property_id, image_url) VALUES
('f38f80b3-f326-4825-9afc-ebc331626875', 'https://picsum.photos/800/600?random=1'),
('f38f80b3-f326-4825-9afc-ebc331626875', 'https://picsum.photos/800/600?random=2'),
('f38f80b3-f326-4825-9afc-ebc331626875', 'https://picsum.photos/800/600?random=3'),
('41a448d4-43ec-411a-a692-2d68e06e0282', 'https://picsum.photos/800/600?random=1'),
('41a448d4-43ec-411a-a692-2d68e06e0282', 'https://picsum.photos/800/600?random=2'),
('414854bf-bdee-45a5-929f-073aedaceea0', 'https://picsum.photos/800/600?random=1'),
('414854bf-bdee-45a5-929f-073aedaceea0', 'https://picsum.photos/800/600?random=2'),
('62dd40da-8238-4d21-b9a7-7f1c24efdd0c', 'https://picsum.photos/800/600?random=1'),
('62dd40da-8238-4d21-b9a7-7f1c24efdd0c', 'https://picsum.photos/800/600?random=2'),
('62dd40da-8238-4d21-b9a7-7f1c24efdd0c', 'https://picsum.photos/800/600?random=3'),
('62dd40da-8238-4d21-b9a7-7f1c24efdd0c', 'https://picsum.photos/800/600?random=4'),
('bc5891ce-6d5e-40d6-8563-f7cebe9667e8', 'https://picsum.photos/800/600?random=1'),
('bc5891ce-6d5e-40d6-8563-f7cebe9667e8', 'https://picsum.photos/800/600?random=2'),
('bc5891ce-6d5e-40d6-8563-f7cebe9667e8', 'https://picsum.photos/800/600?random=3'),
('bc5891ce-6d5e-40d6-8563-f7cebe9667e8', 'https://picsum.photos/800/600?random=4'),
('3df779f2-1f72-44d1-9a31-51929ed130a2', 'https://picsum.photos/800/600?random=1'),
('a8329428-6971-42e8-974a-4df030cd27be', 'https://picsum.photos/800/600?random=1'),
('a8329428-6971-42e8-974a-4df030cd27be', 'https://picsum.photos/800/600?random=2'),
('a8329428-6971-42e8-974a-4df030cd27be', 'https://picsum.photos/800/600?random=3'),
('a8329428-6971-42e8-974a-4df030cd27be', 'https://picsum.photos/800/600?random=4'),
('a8329428-6971-42e8-974a-4df030cd27be', 'https://picsum.photos/800/600?random=5'),
('f8eaf2fc-d6f2-4a8c-a714-5425cc76bbfa', 'https://picsum.photos/800/600?random=1'),
('f8eaf2fc-d6f2-4a8c-a714-5425cc76bbfa', 'https://picsum.photos/800/600?random=2'),
('b7c8ce65-8fa3-4759-bc4e-42a396ef4fc1', 'https://picsum.photos/800/600?random=1');

INSERT INTO selling_properties (property_id, price, is_sold) VALUES
('f38f80b3-f326-4825-9afc-ebc331626875', 258883.7091280503, FALSE),
('62dd40da-8238-4d21-b9a7-7f1c24efdd0c', 128734.8123476912, FALSE);

INSERT INTO renting_properties (property_id, price_per_month, is_occupied) VALUES
('f38f80b3-f326-4825-9afc-ebc331626875', 123423.2931847312, FALSE),
('f8eaf2fc-d6f2-4a8c-a714-5425cc76bbfa', 112302.9182347433, TRUE);

-- mock data for appointments

-- Insert mock data into the users table
INSERT INTO users (user_id, registered_type, email, password, first_name, last_name, phone_number, profile_image_url, is_verified)
VALUES
('123e4567-e89b-12d3-a456-426614174001', 'EMAIL', 'user1@email.com', 'password123', 'User', 'One', '1234567890', 'https://example.com/image1.jpg', TRUE),
('123e4567-e89b-12d3-a456-426614174002', 'EMAIL', 'user2@email.com', 'password456', 'User', 'Two', '9876543210', 'https://example.com/image2.jpg', FALSE),
('123e4567-e89b-12d3-a456-426614174003', 'GOOGLE', 'user3@gmail.com', NULL, 'User', 'Three', '3333333333', 'https://example.com/image3.jpg', TRUE);

-- Insert mock data into the properties table
INSERT INTO properties (property_id, owner_id, property_name, property_description, property_type, address, alley, street, sub_district, district, province, country, postal_code, bedrooms, bathrooms, furnishing, floor, floor_size, unit_number) 
VALUES
('223e4567-e89b-12d3-a456-426614174001', '123e4567-e89b-12d3-a456-426614174001', 'Beautiful House', 'Dream House', 'HOUSE', '123 Main St', NULL, 'Dream Street', 'Dreamville', 'Dream District', 'Dream Province', 'Dream Country', '12345', 1, 1, 'UNFURNISHED', 1, 11.11, 1010),
('223e4567-e89b-12d3-a456-426614174002', '123e4567-e89b-12d3-a456-426614174002', 'Cozy Apartment', 'Sky Towers', 'APARTMENT', '456 Sky Blvd', 'Sky Alley', 'Cloud Street', 'Cloudsville', 'Cloud District', 'Cloud Province', 'Cloud Country', '56789', 2, 2, 'PARTIALLY_FURNISHED', 2, 22.22, 2020),
('223e4567-e89b-12d3-a456-426614174003', '123e4567-e89b-12d3-a456-426614174003', 'Luxury Condo', 'Golden Heights', 'CONDOMINIUM', '789 Gold Ave', 'Gold Alley', 'Golden Street', 'Goldenville', 'Gold District', 'Gold Province', 'Gold Country', '98765', 3, 3, 'FULLY_FURNISHED', 3, 33.33, 3030);

-- Insert mock data into the agreements table
INSERT INTO agreements (agreement_id, property_id, owner_user_id, dweller_user_id, agreement_date)
VALUES
('323e4567-e89b-12d3-a456-426614174001', '223e4567-e89b-12d3-a456-426614174001', '123e4567-e89b-12d3-a456-426614174001', '123e4567-e89b-12d3-a456-426614174002', CURRENT_TIMESTAMP),
('323e4567-e89b-12d3-a456-426614174002', '223e4567-e89b-12d3-a456-426614174002', '123e4567-e89b-12d3-a456-426614174002', '123e4567-e89b-12d3-a456-426614174003', CURRENT_TIMESTAMP),
('323e4567-e89b-12d3-a456-426614174003', '223e4567-e89b-12d3-a456-426614174003', '123e4567-e89b-12d3-a456-426614174001', '123e4567-e89b-12d3-a456-426614174003', CURRENT_TIMESTAMP);

-------------------- VIEWS --------------------

ALTER TABLE users RENAME TO _users;
CREATE VIEW users AS SELECT * FROM _users WHERE deleted_at IS NULL;

ALTER TABLE user_financial_informations RENAME TO _user_financial_informations;
CREATE VIEW user_financial_informations AS SELECT * FROM _user_financial_informations WHERE user_id IN (SELECT user_id FROM _users WHERE deleted_at IS NULL);

ALTER TABLE credit_cards RENAME TO _credit_cards;
CREATE VIEW credit_cards AS SELECT * FROM _credit_cards WHERE user_id IN (SELECT user_id FROM _users WHERE deleted_at IS NULL);

ALTER TABLE user_verifications RENAME TO _user_verifications;
CREATE VIEW user_verifications AS SELECT * FROM _user_verifications WHERE user_id IN (SELECT user_id FROM _users WHERE deleted_at IS NULL);

ALTER TABLE properties RENAME TO _properties;
CREATE VIEW properties AS SELECT *
    FROM _properties
    WHERE (
        deleted_at IS NULL AND
        owner_id IN (SELECT user_id FROM _users WHERE deleted_at IS NULL)
    );

ALTER TABLE property_images RENAME TO _property_images;
CREATE VIEW property_images AS SELECT * FROM _property_images WHERE property_id IN (SELECT property_id FROM properties WHERE deleted_at IS NULL);

ALTER TABLE selling_properties RENAME TO _selling_properties;
CREATE VIEW selling_properties AS SELECT * FROM _selling_properties WHERE property_id IN (SELECT property_id FROM properties WHERE deleted_at IS NULL);

ALTER TABLE renting_properties RENAME TO _renting_properties;
CREATE VIEW renting_properties AS SELECT * FROM _renting_properties WHERE property_id IN (SELECT property_id FROM properties WHERE deleted_at IS NULL);

ALTER TABLE appointments RENAME TO _appointments;
CREATE VIEW appointments AS SELECT *
    FROM _appointments
    WHERE (
     	deleted_at IS NULL AND
        property_id IN (SELECT property_id FROM properties WHERE deleted_at IS NULL) AND
        dweller_user_id IN (SELECT user_id FROM _users WHERE deleted_at IS NULL) AND
        owner_user_id IN (SELECT user_id FROM _users WHERE deleted_at IS NULL)
    );

-------------------- INDEX --------------------

CREATE INDEX idx_users_deleted_at                       ON _users (deleted_at);
CREATE INDEX idx_user_financial_information_deleted_at  ON _users_financial_informations (deleted_at);
CREATE INDEX idx_properties_deleted_at                  ON _properties (deleted_at);
CREATE INDEX idx_property_images_deleted_at             ON _property_images (deleted_at);
CREATE INDEX idx_selling_properties_deleted_at          ON _selling_properties (deleted_at);
CREATE INDEX idx_renting_properties_deleted_at          ON _renting_properties (deleted_at);
CREATE INDEX idx_appointments_deleted_at                ON _appointments (deleted_at);