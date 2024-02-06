CREATE TYPE bank_name AS ENUM('KBANK', 'BBL', 'KTB', 'BAY', 'CIMB', 'TTB', 'SCB', 'GSB');

CREATE TYPE registered_type AS ENUM('EMAIL', 'GOOGLE');

CREATE TABLE users
(
    user_id                             UUID PRIMARY KEY                DEFAULT gen_random_uuid(),
    registered_type                     registered_type                 NOT NULL,
    email                               VARCHAR(50)         UNIQUE      NOT NULL,
    password                            VARCHAR(64)                     DEFAULT NULL,
    first_name                          VARCHAR(50)                     NOT NULL,
    last_name                           VARCHAR(50)                     NOT NULL,
    phone_number                        VARCHAR(10)         UNIQUE      NOT NULL,
    profile_image_url                   VARCHAR(2000)                   DEFAULT NULL,
    credit_card_cardholder_name         VARCHAR(50)                     DEFAULT NULL,
    credit_card_number                  VARCHAR(16)                     DEFAULT NULL,
    credit_card_expiration_month        VARCHAR(2)                      DEFAULT NULL,
    credit_card_expiration_year         VARCHAR(4)                      DEFAULT NULL,
    credit_card_cvv                     VARCHAR(3)                      DEFAULT NULL,
    bank_name                           bank_name                       DEFAULT NULL,
    bank_account_number                 VARCHAR(10)                     DEFAULT NULL,
    citizen_id                          VARCHAR(13)         UNIQUE      DEFAULT NULL,
    citizen_image_url                   VARCHAR(2000)                   DEFAULT NULL,
    is_verified                         BOOLEAN                         DEFAULT FALSE,
    created_at                          TIMESTAMP(0)                    DEFAULT CURRENT_TIMESTAMP,
    updated_at                          TIMESTAMP(0)                    DEFAULT CURRENT_TIMESTAMP,
    deleted_at                          TIMESTAMP(0)                    DEFAULT NULL
);

CREATE TABLE property
(
    property_id              UUID PRIMARY KEY                                       DEFAULT gen_random_uuid(),
    owner_id                 UUID REFERENCES users (user_id)                        NOT NULL,
    description              TEXT                                                   NOT NULL,
    residential_type         VARCHAR(99)                                            NOT NULL,
    project_name             VARCHAR(50),
    address                  VARCHAR(50)                                            NOT NULL,
    alley                    VARCHAR(50),
    street                   VARCHAR(50)                                            NOT NULL,
    sub_district             VARCHAR(50)                                            NOT NULL,
    district                 VARCHAR(50)                                            NOT NULL,
    province                 VARCHAR(50)                                            NOT NULL,
    country                  VARCHAR(50)                                            NOT NULL,
    postal_code              CHAR(5)                                                NOT NULL,
    created_at               TIMESTAMP(0)                                           DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP(0)                                           DEFAULT CURRENT_TIMESTAMP,
    deleted_at               TIMESTAMP(0)                                           DEFAULT NULL
);

CREATE TABLE property_image
(
    property_id UUID REFERENCES property (property_id) ON DELETE CASCADE        NOT NULL,
    image_url       VARCHAR(2000)                                               NOT NULL,
    PRIMARY KEY (property_id, image_url)
);

CREATE TABLE selling_property
(
    property_id UUID PRIMARY KEY REFERENCES property (property_id) ON DELETE CASCADE    NOT NULL,
    price       DOUBLE PRECISION                                                        NOT NULL,
    is_sold     BOOLEAN                                                                 NOT NULL
);

CREATE TABLE renting_property
(
    property_id     UUID PRIMARY KEY REFERENCES property (property_id) ON DELETE CASCADE    NOT NULL,
    price_per_month DOUBLE PRECISION                                                        NOT NULL,
    is_occupied     BOOLEAN                                                                 NOT NULL
);

-------------------- VIEWS --------------------

ALTER TABLE users RENAME TO _users;
CREATE VIEW users AS SELECT * FROM _users WHERE deleted_at IS NULL;

ALTER TABLE property RENAME TO _property;
CREATE VIEW property AS SELECT * FROM _property WHERE deleted_at IS NULL;

-------------------- RULES --------------------

CREATE RULE soft_deletion AS ON DELETE TO users DO INSTEAD (
    UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE user_id = old.user_id and deleted_at IS NULL
);

CREATE RULE soft_deletion AS ON DELETE TO property DO INSTEAD (
    UPDATE property SET deleted_at = CURRENT_TIMESTAMP WHERE property_id = old.property_id and deleted_at IS NULL
);

CREATE RULE delete_users AS ON UPDATE TO users
    WHERE old.deleted_at IS NULL AND new.deleted_at IS NOT NULL
    DO ALSO UPDATE property SET deleted_at = new.deleted_at WHERE owner_id = old.user_id;

-------------------- DUMMY DATA --------------------

INSERT INTO users (user_id, registered_type, email, password, first_name, last_name, phone_number, citizen_id, profile_image_url, credit_card_cardholder_name, credit_card_number, credit_card_expiration_month, credit_card_expiration_year, credit_card_cvv, bank_name, bank_account_number, is_verified) VALUES
('f38f80b3-f326-4825-9afc-ebc331626555', 'EMAIL', 'johnd@email.com', '$2a$10$eEkTbe/JskFiociJ8U/bGOwwiea9dZ6sN7ac9ZvuiUgtrekZ7b.ya', 'John', 'Doe', '1234567890', '1234567890123', 'https://picsum.photos/200/300?random=1', 'JOHN DOE', '1234123412341234', '12', '2023', '123', 'KBANK', '1234567890', TRUE),
('bc5891ce-d6f2-d6f2-d6f2-ebc331626555', 'EMAIL', 'sams@email.com', '$2a$10$eEkTbe/JskFiociJ8U/bGOwwiea9dZ6sN7ac9Zvuhfkdle9405.ya', 'Sam', 'Smith', '0987654321', NULL, 'https://picsum.photos/200/300?random=3', 'SAM SMITH', '4321432143214321', '11', '2024', '321', 'BBL', '1234567890', FALSE),
('62dd40da-f326-4825-9afc-2d68e06e0282', 'GOOGLE', 'gmail@gmail.com', NULL, 'C', 'C', '3333333333', '3333333333333', 'https://picsum.photos/200/300?random=1', 'C C', '1234123412341234', '12', '2023', '123', 'SCB', '1234567890', TRUE);

INSERT INTO property (property_id, owner_id, description, residential_type, project_name, address, alley, street, sub_district, district, province, country, postal_code) VALUES
('f38f80b3-f326-4825-9afc-ebc331626875', 'f38f80b3-f326-4825-9afc-ebc331626555', 'Et sequi dolor praes', 'Sequi reiciendis odi', 'Anita', 'Quas iusto expedita ', 'Delisa', 'Grace', 'Michael', 'Christine', 'Anthony', 'Andrew', '53086'),
('41a448d4-43ec-411a-a692-2d68e06e0282', 'f38f80b3-f326-4825-9afc-ebc331626555', 'Impedit quae itaque ', 'Mollitia quidem quas', 'Rose', 'Sunt fuga quo perspi', 'Raquel', 'Brandy', 'Jacob', 'Lino', 'Edward', 'Reginald', '12894'),
('414854bf-bdee-45a5-929f-073aedaceea0', 'f38f80b3-f326-4825-9afc-ebc331626555', 'Architecto iure labo', 'Maiores magnam quaer', 'Michele', 'Pariatur temporibus ', 'Robert', 'Nancy', 'Barbara', 'David', 'Henry', 'David', '24264'),
('62dd40da-8238-4d21-b9a7-7f1c24efdd0c', '62dd40da-f326-4825-9afc-2d68e06e0282', 'Optio in asperiores ', 'Consectetur doloribu', 'Charles', 'Ea nobis mollitia ea', 'Tina', 'Linda', 'Ronald', 'Julia', 'Russell', 'William', '10287'),
('bc5891ce-6d5e-40d6-8563-f7cebe9667e8', '62dd40da-f326-4825-9afc-2d68e06e0282', 'Sunt at totam animi ', 'In ratione veritatis', 'Jonathan', 'Unde natus nesciunt ', 'Norma', 'Gregory', 'Donovan', 'Charles', 'Kevin', 'Tyrone', '10055'),
('3df779f2-1f72-44d1-9a31-51929ed130a2', '62dd40da-f326-4825-9afc-2d68e06e0282', 'Animi vero ipsa nihi', 'Itaque porro veniam ', 'Roger', 'Totam nam minus veni', 'Allen', 'Linda', 'Bobby', 'Nora', 'James', 'Lucinda', '01229'),
('a8329428-6971-42e8-974a-4df030cd27be', '62dd40da-f326-4825-9afc-2d68e06e0282', 'Numquam sit dicta be', 'Dignissimos corrupti', 'Diane', 'Consequatur incidunt', 'Cecil', 'David', 'Nancy', 'Brandon', 'John', 'Lillian', '48668'),
('f8eaf2fc-d6f2-4a8c-a714-5425cc76bbfa', '62dd40da-f326-4825-9afc-2d68e06e0282', 'Iure nostrum ab reru', 'Natus aliquid fuga, ', 'Matthew', 'Nisi officia nemo au', 'Keith', 'Joseph', 'Joseph', 'Goldie', 'Danika', 'Bernice', '47550'),
('b7c8ce65-8fa3-4759-bc4e-42a396ef4fc1', '62dd40da-f326-4825-9afc-2d68e06e0282', 'Aut nemo incidunt ul', 'Quasi facilis aliqui', 'Annie', 'Porro molestias rati', 'Brian', 'Gregory', 'Geraldine', 'Edward', 'Charles', 'James', '97186');

INSERT INTO property_image (property_id, image_url) VALUES
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

INSERT INTO selling_property (property_id, price, is_sold) VALUES
('f38f80b3-f326-4825-9afc-ebc331626875', 258883.7091280503, FALSE),
('62dd40da-8238-4d21-b9a7-7f1c24efdd0c', 128734.8123476912, FALSE);

INSERT INTO renting_property (property_id, price_per_month, is_occupied) VALUES
('f38f80b3-f326-4825-9afc-ebc331626875', 123423.2931847312, FALSE),
('f8eaf2fc-d6f2-4a8c-a714-5425cc76bbfa', 112302.9182347433, TRUE);
