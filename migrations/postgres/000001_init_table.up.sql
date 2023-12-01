DROP TABLE IF EXISTS user_to_channel                CASCADE;
DROP TABLE IF EXISTS channel                CASCADE;
DROP TABLE IF EXISTS users                CASCADE;



CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE IF NOT EXISTS users (
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
   username VARCHAR (255) NOT NULL,
   email VARCHAR (255) UNIQUE NOT NULL,
   password VARCHAR (255) NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   login_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS channel (
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
   name VARCHAR (255) NOT NULL,
   user_id uuid REFERENCES users (id),
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_to_channel (
   user_id uuid REFERENCES users (id),
   channel_id uuid REFERENCES channel (id),
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   CONSTRAINT userid_channelid_pkey PRIMARY KEY (user_id,  channel_id)
);