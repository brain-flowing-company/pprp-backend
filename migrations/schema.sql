-- company table
CREATE TABLE company
(
    company_id      UUID PRIMARY KEY NOT NULL,
    company_name    VARCHAR(35)      NOT NULL,
    company_address VARCHAR(300)     NOT NULL,
    company_phone   CHAR(14)         NOT NULL
);

-- admin table
CREATE TABLE admin
(
    admin_id    UUID PRIMARY KEY NOT NULL,
    admin_fname VARCHAR(35)      NOT NULL,
    admin_lname VARCHAR(35)      NOT NULL,
    admin_phone CHAR(14)         NOT NULL
);

-- advertisement table
CREATE TABLE advertisement
(
    ad_id      UUID PRIMARY KEY                     NOT NULL,
    start_date TIMESTAMP WITHOUT TIME ZONE          NOT NULL,
    end_date   TIMESTAMP WITHOUT TIME ZONE,
    company_id UUID REFERENCES company (company_id) NOT NULL,
    admin_id   UUID REFERENCES admin (admin_id)     NOT NULL
);

-- user table
CREATE TABLE "user"
(
    email            VARCHAR(319) PRIMARY KEY NOT NULL,
    fname            VARCHAR(35)              NOT NULL,
    lname            VARCHAR(35)              NOT NULL,
    profile_image    VARCHAR(2000),
    expired_date_kyc TIMESTAMP WITHOUT TIME ZONE
);

-- issue table
CREATE TABLE issue
(
    issue_id         UUID PRIMARY KEY                       NOT NULL,
    issue_type       VARCHAR(100)                           NOT NULL,
    is_resolved      BOOLEAN                                NOT NULL,
    issue_time_stamp TIMESTAMP WITHOUT TIME ZONE            NOT NULL,
    issue_message    TEXT                                   NOT NULL,
    issue_image      VARCHAR(2000),
    user_email       VARCHAR(319) REFERENCES "user" (email) NOT NULL
);

-- issue_review table
CREATE TABLE issue_review
(
    admin_id      UUID REFERENCES admin (admin_id) NOT NULL,
    issue_id      UUID REFERENCES issue (issue_id) NOT NULL,
    date_resolved TIMESTAMP WITHOUT TIME ZONE,
    PRIMARY KEY (admin_id, issue_id)
);

-- owner table
CREATE TABLE owner
(
    owner_email VARCHAR(319) PRIMARY KEY REFERENCES "user" (email) NOT NULL
);

-- dweller table
CREATE TABLE dweller
(
    dweller_email VARCHAR(319) PRIMARY KEY REFERENCES "user" (email) NOT NULL
);

-- property_review table
CREATE TABLE property_review
(
    property_id       UUID PRIMARY KEY                                NOT NULL,
    dweller_email     VARCHAR(319) REFERENCES dweller (dweller_email) NOT NULL,
    review_time_stamp TIMESTAMP WITHOUT TIME ZONE                     NOT NULL,
    image             VARCHAR(2000),
    score             DOUBLE PRECISION                                NOT NULL CHECK (score >= 0 AND score <= 5),
    description       TEXT                                            NOT NULL
);

-- like table
CREATE TABLE "like"
(
    property_id   UUID REFERENCES property_review (property_id)   NOT NULL,
    dweller_email VARCHAR(319) REFERENCES dweller (dweller_email) NOT NULL,
    PRIMARY KEY (property_id, dweller_email)
);

-- property_listing table
-- Remove Image from property-listing table
CREATE TABLE property_listing
(
    property_id              UUID PRIMARY KEY                            NOT NULL,
    owner_email              VARCHAR(319) REFERENCES owner (owner_email) NOT NULL,
    owner_contact            VARCHAR(100)                                NOT NULL,
    description              TEXT                                        NOT NULL,
    residential_type         VARCHAR(99)                                 NOT NULL,
    project_name             VARCHAR(50),
    address                  VARCHAR(50)                                 NOT NULL,
    alley                    VARCHAR(50),
    street                   VARCHAR(50)                                 NOT NULL,
    sub_district             VARCHAR(50)                                 NOT NULL,
    district                 VARCHAR(50)                                 NOT NULL,
    province                 VARCHAR(50)                                 NOT NULL,
    country                  VARCHAR(50)                                 NOT NULL,
    postal_code              CHAR(5)                                     NOT NULL,
    property_list_time_stamp TIMESTAMP WITHOUT TIME ZONE                 NOT NULL
);

-- appointment table
CREATE TABLE appointment
(
    appointment_id         UUID PRIMARY KEY                                NOT NULL,
    property_id            UUID REFERENCES property_listing (property_id)  NOT NULL,
    dweller_email          VARCHAR(319) REFERENCES dweller (dweller_email) NOT NULL,
    appointment_time_stamp TIMESTAMP WITHOUT TIME ZONE                     NOT NULL,
    appointment_date       TIMESTAMP WITHOUT TIME ZONE                     NOT NULL,
    is_confirmed           BOOLEAN                                         NOT NULL,
    is_rejected            BOOLEAN                                         NOT NULL,
    is_met                 BOOLEAN                                         NOT NULL
);

-- image table
CREATE TABLE image
(
    property_id UUID REFERENCES property_listing (property_id) NOT NULL,
    image       VARCHAR(2000)                                            NOT NULL,
    PRIMARY KEY (property_id, image)
);

-- property_listing_for_sell table
CREATE TABLE property_listing_for_sell
(
    property_id UUID PRIMARY KEY REFERENCES property_listing (property_id) NOT NULL,
    price       DOUBLE PRECISION                                           NOT NULL,
    is_sold     BOOLEAN                                                    NOT NULL
);

-- property_listing_for_rent table
CREATE TABLE property_listing_for_rent
(
    property_id     UUID PRIMARY KEY REFERENCES property_listing (property_id) NOT NULL,
    price_per_month DOUBLE PRECISION                                           NOT NULL,
    is_occupied     BOOLEAN                                                    NOT NULL
);

-- transaction table
CREATE TABLE transaction
(
    transaction_id UUID PRIMARY KEY            NOT NULL,
    timestamp      TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    payment_method CHAR(50)                    NOT NULL,
    is_released    BOOLEAN                     NOT NULL,
    is_canceled    BOOLEAN                     NOT NULL,
    amount         DOUBLE PRECISION            NOT NULL
);

-- rent table
CREATE TABLE rent
(
    property_id     UUID REFERENCES property_listing_for_rent (property_id) NOT NULL,
    dweller_email   VARCHAR(319) REFERENCES dweller (dweller_email)         NOT NULL,
    transaction_id  UUID REFERENCES transaction (transaction_id)            NOT NULL,
    rent_start_date TIMESTAMP WITHOUT TIME ZONE                             NOT NULL,
    PRIMARY KEY (transaction_id)
);

-- sell table
CREATE TABLE sell
(
    property_id    UUID REFERENCES property_listing_for_sell (property_id) NOT NULL,
    dweller_email  VARCHAR(319) REFERENCES dweller (dweller_email)         NOT NULL,
    transaction_id UUID REFERENCES transaction (transaction_id)            NOT NULL,
    PRIMARY KEY (transaction_id)
);

CREATE INDEX idx_district ON property_listing(district);

CREATE INDEX idx_price_sell ON property_listing_for_sell(price);
CLUSTER property_listing_for_sell USING idx_price_sell;

CREATE INDEX idx_price_rent ON property_listing_for_rent(price_per_month);
CLUSTER property_listing_for_rent USING idx_price_rent;