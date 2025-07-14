CREATE TABLE public.users (
    id serial  PRIMARY KEY,
	version integer default '1',
    logging VARCHAR(255) UNIQUE NOT NULL,
    pwd_hash bytea NOT NULL, -- Stocker le hash du mot de passe
    nom VARCHAR(255),
    prenom VARCHAR(255),
    roles VARCHAR(255) DEFAULT 'user' -- Pour gérer les rôles (ex: admin, user)
);

INSERT INTO public.users (logging, pwd_hash, nom, prenom, roles) 
        VALUES ('admin', '$2a$10$ai/vI/n5obKKMA6ENfV.kO7qW0DJNkiYPPiMnshpTXH2QwwVErtjC', 'admin', 'admin', 'admin');


CREATE TABLE feedback (
    id serial PRIMARY KEY,
    message text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
